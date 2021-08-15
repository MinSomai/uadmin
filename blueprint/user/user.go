package user

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/admin"
	utils2 "github.com/uadmin/uadmin/blueprint/auth/utils"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	"github.com/uadmin/uadmin/blueprint/user/migrations"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"text/template"
	"time"
)

type Blueprint struct {
	interfaces.Blueprint
}

type PasswordValidationStruct struct {
	Password string `valid:"password-uadmin"`
	ConfirmedPassword string
}

type ForgotPasswordHandlerParams struct {
	Email string    `form:"email" json:"email" xml:"email"  binding:"required" valid:"email"`
}

type ResetPasswordHandlerParams struct {
	Code string    `form:"code" json:"code" xml:"code"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password"  binding:"required"`
}

type ChangePasswordHandlerParams struct {
	OldPassword string    `form:"old_password" json:"old_password" xml:"old_password"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password"  binding:"required"`
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	mainRouter.GET("/reset-password", func(ctx *gin.Context) {
		type Context struct {
			interfaces.AdminContext
		}
		c := &Context{}
		admin.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
		tr := interfaces.NewTemplateRenderer("Reset Password")
		tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("resetpassword"), c, interfaces.FuncMap)
	})
	group.POST("/api/forgot", func(ctx *gin.Context) {
		var json ForgotPasswordHandlerParams
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var err1 error
		_, err1 = govalidator.ValidateStruct(&json)
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
			return
		}
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var user interfaces.User
		db.Model(interfaces.User{}).Where(&interfaces.User{Email: json.Email}).First(&user)
		if user.ID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with this email not found"})
			return
		}
		templateWriter := bytes.NewBuffer([]byte{})
		template1, err := template.ParseFS(interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("email/forgot"))
		if err != nil {
			interfaces.Trail(interfaces.ERROR, "RenderHTML unable to parse %s. %s", interfaces.CurrentConfig.GetPathToTemplate("email/forgot"), err)
			ctx.JSON(http.StatusBadRequest, utils.ApiBadResponse(err.Error()))
			return
		}
		type Context struct {
			Name    string
			Website string
			URL     string
		}

		c := Context{}
		c.Name = user.Username
		c.Website = interfaces.CurrentConfig.D.Uadmin.SiteName
		host := interfaces.CurrentConfig.D.Uadmin.PoweredOnSite
		// @todo, generate code to restore access
		actionExpiresAt := time.Now()
		actionExpiresAt = actionExpiresAt.Add(time.Duration(interfaces.CurrentConfig.D.Uadmin.ForgotCodeExpiration)*time.Minute)
		var oneTimeAction = interfaces.OneTimeAction{
			User:       user,
			ExpiresOn: actionExpiresAt,
			Code: utils.RandStringRunesForOneTimeAction(32),
			ActionType: 1,

		}

		db.Model(interfaces.OneTimeAction{}).Save(&oneTimeAction)
		link := host + interfaces.CurrentConfig.D.Uadmin.RootAdminURL + "/resetpassword?key=" + oneTimeAction.Code
		c.URL = link
		err = template1.Execute(templateWriter, c)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ApiBadResponse(err.Error()))
			interfaces.Trail(interfaces.ERROR, "RenderHTML unable to parse %s. %s", interfaces.CurrentConfig.GetPathToTemplate("email/forgot"), err)
			return
		}
		subject := "Password reset for admin panel on the " + interfaces.CurrentConfig.D.Uadmin.SiteName
		err = utils.SendEmail(interfaces.CurrentConfig.D.Uadmin.EmailFrom, []string{user.Email}, []string{}, []string{}, subject, templateWriter.String())
		return
	})
	group.POST("/api/reset-password", func(ctx *gin.Context) {
		var json ResetPasswordHandlerParams
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var err1 error
		_, err1 = govalidator.ValidateStruct(&json)
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
			return
		}
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var oneTimeAction interfaces.OneTimeAction
		db.Model(interfaces.OneTimeAction{}).Where(&interfaces.OneTimeAction{Code: json.Code, IsUsed: false}).Preload("User").First(&oneTimeAction)
		if oneTimeAction.ID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No such code found"})
			return
		}
		if oneTimeAction.ExpiresOn.Before(time.Now()) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code is expired"})
			return
		}
		passwordValidationStruct := &PasswordValidationStruct{
			Password: json.Password,
			ConfirmedPassword: json.ConfirmedPassword,
		}
		_, err := govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashedPassword, err := utils2.HashPass(json.Password, oneTimeAction.User.Salt)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		oneTimeAction.User.Password = hashedPassword
		oneTimeAction.User.IsPasswordConfigured = true
		oneTimeAction.IsUsed = true
		db.Save(&oneTimeAction.User)
		db.Save(&oneTimeAction)
	})
	group.POST("/api/change-password", func(ctx *gin.Context) {
		var json ChangePasswordHandlerParams
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var err1 error
		_, err1 = govalidator.ValidateStruct(&json)
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
			return
		}
		passwordValidationStruct := &PasswordValidationStruct{
			Password: json.Password,
			ConfirmedPassword: json.ConfirmedPassword,
		}
		_, err := govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		hashedPassword, err := utils2.HashPass(json.OldPassword, user.Salt)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// @todo, get it back once stabilize pass api
		//if hashedPassword != user.Password {
		//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password doesn't match current one"})
		//	return
		//}
		hashedPassword, err = utils2.HashPass(json.Password, user.Salt)
		user.Password = hashedPassword
		user.IsPasswordConfigured = true
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	group.POST("/api/disable-2fa", func(ctx *gin.Context) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		user.OTPRequired = false
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	group.POST("/api/enable-2fa", func(ctx *gin.Context) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		user.OTPRequired = true
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	mainRouter.NoRoute(func(ctx *gin.Context) {
		// ctx.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		type Context struct {
			interfaces.AdminContext
			Menu     string
		}

		c := &Context{}
		admin.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
		//
		//if r.Form.Get("err_msg") != "" {
		//	c.ErrMsg = r.Form.Get("err_msg")
		//}
		//if code, err := strconv.ParseUint(r.Form.Get("err_code"), 10, 16); err == nil {
		//	c.ErrCode = int(code)
		//}NewAdminPage
		ctx.Status(404)
		tr := interfaces.NewTemplateRenderer("Page not found")
		tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("404"), c, interfaces.FuncMap)
	})
	usersAdminPage := admin.NewGormAdminPage(nil, func() (interface{}, interface{}) {return nil, nil}, "")
	usersAdminPage.PageName = "Users"
	usersAdminPage.Slug = "users"
	usersAdminPage.BlueprintName = "user"
	usersAdminPage.Router = mainRouter
	err := admin.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(usersAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
	usermodelAdminPage := admin.NewGormAdminPage(usersAdminPage, func() (interface{}, interface{}) {return &interfaces.User{}, &[]*interfaces.User{}}, "user")
	usermodelAdminPage.PageName = "Users"
	usermodelAdminPage.Slug = "user"
	usermodelAdminPage.BlueprintName = "user"
	usermodelAdminPage.Router = mainRouter
	adminContext := &interfaces.AdminContext{}
	userForm := interfaces.NewFormFromModelFromGinContext(adminContext, &interfaces.User{}, make([]string, 0), []string{"Username", "FirstName", "LastName", "Email", "Photo", "LastLogin"}, true, "")
	usermodelAdminPage.Form = userForm
	usermodelAdminPage.ListDisplay.ClearAllFields()
	usernameField, _ := userForm.FieldRegistry.GetByName("Username")
	usernameListDisplay := interfaces.NewListDisplay(usernameField)
	usermodelAdminPage.ListDisplay.AddField(usernameListDisplay)
	firstNameField, _ := userForm.FieldRegistry.GetByName("FirstName")
	usermodelAdminPage.ListDisplay.AddField(interfaces.NewListDisplay(firstNameField))
	lastNameField, _ := userForm.FieldRegistry.GetByName("LastName")
	usermodelAdminPage.ListDisplay.AddField(interfaces.NewListDisplay(lastNameField))
	emailField, _ := userForm.FieldRegistry.GetByName("Email")
	emailListDisplay := interfaces.NewListDisplay(emailField)
	usermodelAdminPage.ListDisplay.AddField(emailListDisplay)
	photoField, _ := userForm.FieldRegistry.GetByName("Photo")
	usermodelAdminPage.ListDisplay.AddField(interfaces.NewListDisplay(photoField))
	lastLoginField, _ := userForm.FieldRegistry.GetByName("LastLogin")
	usermodelAdminPage.ListDisplay.AddField(interfaces.NewListDisplay(lastLoginField))
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	userModelDescription := interfaces.ProjectModels.GetModelFromInterface(&interfaces.User{})
	statementUser := userModelDescription.Statement
	fieldEmail, _ := statementUser.Schema.FieldsByDBName["email"]
	searchField := &interfaces.SearchField{
		Field: fieldEmail,
	}
	usermodelAdminPage.SearchFields = append(usermodelAdminPage.SearchFields, searchField)
	listFilter := &interfaces.ListFilter{
		UrlFilteringParam: "IsSuperUser__exact",
		Title: "Is super user ?",
	}
	listFilter.OptionsToShow = append(listFilter.OptionsToShow, &interfaces.FieldChoice{DisplayAs: "Yes", Value: true})
	listFilter.OptionsToShow = append(listFilter.OptionsToShow, &interfaces.FieldChoice{DisplayAs: "No", Value: false})
	usermodelAdminPage.ListFilter.Add(listFilter)
	err = usersAdminPage.SubPages.AddAdminPage(usermodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
	usergroupsAdminPage := admin.NewGormAdminPage(usersAdminPage, func() (interface{}, interface{}) {return &interfaces.UserGroup{}, &[]*interfaces.UserGroup{}}, "usergroup")
	usergroupsAdminPage.PageName = "User groups"
	usergroupsAdminPage.Slug = "usergroup"
	usergroupsAdminPage.BlueprintName = "user"
	usergroupsAdminPage.Router = mainRouter
	adminContext = &interfaces.AdminContext{}
	userGroupForm := interfaces.NewFormFromModelFromGinContext(adminContext, &interfaces.UserGroup{}, make([]string, 0), []string{}, true, "")
	usergroupsAdminPage.ListDisplay.ClearAllFields()
	usergroupsAdminPage.Form = userGroupForm
	groupNameField, _ := userGroupForm.FieldRegistry.GetByName("GroupName")
	groupNameListDisplay := interfaces.NewListDisplay(groupNameField)
	usergroupsAdminPage.ListDisplay.AddField(groupNameListDisplay)
	err = usersAdminPage.SubPages.AddAdminPage(usergroupsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
}

type UsernameFormOptions struct {
	interfaces.FieldFormOptions
}

type UserPhotoOptions struct {
	interfaces.FieldFormOptions
}

type OtpRequiredOptions struct {
	interfaces.FieldFormOptions
}

type LastLoginOptions struct {
	interfaces.FieldFormOptions
}

type ExpiresOnOptions struct {
	interfaces.FieldFormOptions
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.OneTimeAction{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.User{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.UserGroup{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.Permission{}})
	fieldChoiceRegistry := interfaces.FieldChoiceRegistry{}
	fieldChoiceRegistry.Choices = make([]*interfaces.FieldChoice, 0)
	formOptions := &UsernameFormOptions{
		FieldFormOptions: interfaces.FieldFormOptions{
			Name: "UsernameOptions",
			Initial: "InitialUsername",
			DisplayName: "Username",
			Validators: make([]interfaces.IValidator, 0),
			Choices: &fieldChoiceRegistry,
			HelpText: "help for username",
		},
	}
	interfaces.CurrentConfig.AddFieldFormOptions(formOptions)
	userPhotoOptions := &UserPhotoOptions{
		FieldFormOptions: interfaces.FieldFormOptions{
			Name: "UserPhotoOptions",
			WidgetType: "image",
		},
	}
	interfaces.CurrentConfig.AddFieldFormOptions(userPhotoOptions)
	otpRequiredOptions := &OtpRequiredOptions{
		FieldFormOptions: interfaces.FieldFormOptions{
			Name: "OTPRequiredOptions",
			WidgetType: "hidden",
		},
	}
	interfaces.CurrentConfig.AddFieldFormOptions(otpRequiredOptions)
	lastLoginOptions := &LastLoginOptions{
		FieldFormOptions: interfaces.FieldFormOptions{
			Name: "LastLoginOptions",
			ReadOnly: true,
		},
	}
	interfaces.CurrentConfig.AddFieldFormOptions(lastLoginOptions)
	expiresOnOptions := &ExpiresOnOptions{
		FieldFormOptions: interfaces.FieldFormOptions{
			Name: "ExpiresOnOptions",
			ReadOnly: true,
		},
	}
	interfaces.CurrentConfig.AddFieldFormOptions(expiresOnOptions)

	interfaces.UadminValidatorRegistry.AddValidator("username-unique", func (i interface{}, o interface{}) error {
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var cUsers int64
		db.Model(&interfaces.User{}).Where(&interfaces.User{Username: i.(string)}).Count(&cUsers)
		if cUsers == 0 {
			return nil
		}
		return fmt.Errorf("user with name %s is already registered", i.(string))
	})

	interfaces.UadminValidatorRegistry.AddValidator("email-unique", func (i interface{}, o interface{}) error {
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var cUsers int64
		db.Model(&interfaces.User{}).Where(&interfaces.User{Email: i.(string)}).Count(&cUsers)
		if cUsers == 0 {
			return nil
		}
		return fmt.Errorf("user with email %s is already registered", i.(string))
	})

	interfaces.UadminValidatorRegistry.AddValidator("username-uadmin", func (i interface{}, o interface{}) error {
		minLength := interfaces.CurrentConfig.D.Auth.MinUsernameLength
		maxLength := interfaces.CurrentConfig.D.Auth.MaxUsernameLength
		currentUsername := i.(string)
		if maxLength < len(currentUsername) || len(currentUsername) < minLength {
			return fmt.Errorf("length of the username has to be between %d and %d symbols", minLength, maxLength)
		}
		return nil
	})

	interfaces.UadminValidatorRegistry.AddValidator("password-uadmin", func (i interface{}, o interface{}) error {
		passwordStruct := o.(PasswordValidationStruct)
		if passwordStruct.Password != passwordStruct.ConfirmedPassword {
			return fmt.Errorf("password doesn't equal to confirmed password")
		}
		if len(passwordStruct.Password) < interfaces.CurrentConfig.D.Auth.MinPasswordLength {
			return fmt.Errorf("length of the password has to be at least %d symbols", interfaces.CurrentConfig.D.Auth.MinPasswordLength)
		}
		return nil
	})

	govalidator.CustomTypeTagMap.Set("username-unique", func(i interface{}, o interface{}) bool {
		validator, _ := interfaces.UadminValidatorRegistry.GetValidator("username-unique")
		userExists := validator(i, o)
		return userExists == nil
	})
	govalidator.CustomTypeTagMap.Set("email-unique", func(i interface{}, o interface{}) bool {
		validator, _ := interfaces.UadminValidatorRegistry.GetValidator("email-unique")
		emailExists := validator(i, o)
		return emailExists == nil
	})
	govalidator.CustomTypeTagMap.Set("username-uadmin", func(i interface{}, o interface{}) bool {
		validator, _ := interfaces.UadminValidatorRegistry.GetValidator("username-uadmin")
		isValidUsername := validator(i, o)
		return isValidUsername == nil
	})
	govalidator.CustomTypeTagMap.Set("password-uadmin", func(i interface{}, o interface{}) bool {
		validator, _ := interfaces.UadminValidatorRegistry.GetValidator("password-uadmin")
		isValidPassword := validator(i, o)
		return isValidPassword == nil
	})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "user",
		Description:       "this blueprint is about users",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
