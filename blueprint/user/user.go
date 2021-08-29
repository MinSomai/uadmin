package user

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	utils2 "github.com/uadmin/uadmin/blueprint/auth/utils"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	"github.com/uadmin/uadmin/blueprint/user/migrations"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"strconv"
	"strings"
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
		interfaces.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
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
		oneTimeAction.User.IsPasswordUsable = true
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
		user.IsPasswordUsable = true
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
		interfaces.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
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
	usersAdminPage := interfaces.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) {return nil, nil},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {return nil},
	)
	usersAdminPage.PageName = "Users"
	usersAdminPage.Slug = "users"
	usersAdminPage.BlueprintName = "user"
	usersAdminPage.Router = mainRouter
	err := interfaces.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(usersAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
	var usermodelAdminPage *interfaces.AdminPage
	usermodelAdminPage = interfaces.NewGormAdminPage(
		usersAdminPage,
		func() (interface{}, interface{}) {return &interfaces.User{}, &[]*interfaces.User{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {
			fields := []string{"Username", "FirstName", "LastName", "Email", "Active", "IsStaff", "IsSuperUser", "Password", "Photo", "LastLogin", "ExpiresOn"}
			if ctx.GetUserObject().IsSuperUser {
				fields = append(fields, "UserGroups")
				fields = append(fields, "Permissions")
			}
			form := interfaces.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			if ctx.GetUserObject().IsSuperUser {
				usergroupsField, _ := form.FieldRegistry.GetByName("UserGroups")
				usergroupsField.SetUpField = func(w interfaces.IWidget, modelI interface{}, v interface{}, afo interfaces.IAdminFilterObjects) error {
					model := modelI.(*interfaces.User)
					vTmp := v.([]string)
					var usergroup *interfaces.UserGroup
					if model.ID != 0 {
						afo.GetUadminDatabase().Db.Model(model).Association("UserGroups").Clear()
						model.UserGroups = make([]interfaces.UserGroup, 0)
					}
					for _, ID := range vTmp {
						afo.GetUadminDatabase().Db.First(&usergroup, ID)
						if usergroup.ID != 0 {
							model.UserGroups = append(model.UserGroups, *usergroup)
						}
						usergroup = nil
					}
					return nil
				}
				userGroupsWidget := usergroupsField.FieldConfig.Widget.(*interfaces.ChooseFromSelectWidget)
				userGroupsWidget.AddNewLink = fmt.Sprintf("%s/%s/usergroup/edit/%s?_to_field=id&_popup=1", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, usersAdminPage.Slug, "new")
				userGroupsWidget.AddNewTitle = "Add another group"
				userGroupsWidget.PopulateLeftSide = func()[]*interfaces.SelectOptGroup {
					var groups []*interfaces.UserGroup
					uadminDatabase := interfaces.NewUadminDatabase()
					uadminDatabase.Db.Find(&groups)
					ret := make([]*interfaces.SelectOptGroup, 0)
					for _, group := range groups {
						ret = append(ret, &interfaces.SelectOptGroup{
							OptLabel: group.GroupName,
							Value: group.ID,
						})
					}
					uadminDatabase.Close()
					return ret
				}
				userGroupsWidget.PopulateRightSide = func()[]*interfaces.SelectOptGroup {
					ret := make([]*interfaces.SelectOptGroup, 0)
					user := modelI.(*interfaces.User)
					if user.ID != 0 {
						var groups []*interfaces.UserGroup
						uadminDatabase := interfaces.NewUadminDatabase()
						uadminDatabase.Db.Model(user).Association("UserGroups").Find(&groups)
						ret = make([]*interfaces.SelectOptGroup, 0)
						for _, group := range groups {
							ret = append(ret, &interfaces.SelectOptGroup{
								OptLabel: group.GroupName,
								Value:    group.ID,
							})
						}
						uadminDatabase.Close()
						return ret
					} else {
						formD := ctx.GetPostForm()
						if formD != nil {
							Ids := strings.Split(formD.Value["UserGroups"][0], ",")
							IdI := make([]uint, 0)
							for _, tmp := range Ids {
								tmpI, _ := strconv.Atoi(tmp)
								IdI = append(IdI, uint(tmpI))
							}
							if len(IdI) > 0 {
								var groups []*interfaces.UserGroup
								uadminDatabase := interfaces.NewUadminDatabase()
								uadminDatabase.Db.Find(&groups, IdI)
								ret = make([]*interfaces.SelectOptGroup, 0)
								for _, group := range groups {
									ret = append(ret, &interfaces.SelectOptGroup{
										OptLabel: group.GroupName,
										Value:    group.ID,
									})
								}
								uadminDatabase.Close()
								return ret
							}
						}
					}
					return ret
				}
				userGroupsWidget.LeftSelectTitle = "Available groups"
				userGroupsWidget.LeftSelectHelp = "This is the list of available groups. You may choose some by selecting them in the box below and then clicking the \"Choose\" arrow between the two boxes."
				userGroupsWidget.LeftSearchSelectHelp = "Type into this box to filter down the list of available groups."
				userGroupsWidget.LeftHelpChooseAll = "Click to choose all groups at once."
				userGroupsWidget.RightSelectTitle = "Chosen groups"
				userGroupsWidget.RightSelectHelp = "This is the list of chosen groups. You may remove some by selecting them in the box below and then clicking the \"Remove\" arrow between the two boxes."
				userGroupsWidget.RightSearchSelectHelp = ""
				userGroupsWidget.RightHelpChooseAll = "Click to remove all chosen groups at once."
				userGroupsWidget.HelpText = "The groups this user belongs to. A user will get all permissions granted to each of their groups. Hold down \"Control\", or \"Command\" on a Mac, to select more than one."
				permissionsField, _ := form.FieldRegistry.GetByName("Permissions")
				permissionsField.SetUpField = func(w interfaces.IWidget, modelI interface{}, v interface{}, afo interfaces.IAdminFilterObjects) error {
					model := modelI.(*interfaces.User)
					vTmp := v.([]string)
					var permission *interfaces.Permission
					if model.ID != 0 {
						afo.GetUadminDatabase().Db.Model(model).Association("Permissions").Clear()
						model.Permissions = make([]interfaces.Permission, 0)
					}
					for _, ID := range vTmp {
						afo.GetUadminDatabase().Db.First(&permission, ID)
						if permission.ID != 0 {
							model.Permissions = append(model.Permissions, *permission)
						}
						permission = nil
					}
					return nil
				}
				permissionsWidget := permissionsField.FieldConfig.Widget.(*interfaces.ChooseFromSelectWidget)
				permissionsWidget.PopulateLeftSide = func()[]*interfaces.SelectOptGroup {
					var permissions []*interfaces.Permission
					uadminDatabase := interfaces.NewUadminDatabase()
					uadminDatabase.Db.Preload("ContentType").Find(&permissions)
					ret := make([]*interfaces.SelectOptGroup, 0)
					for _, permission := range permissions {
						ret = append(ret, &interfaces.SelectOptGroup{
							OptLabel: permission.ShortDescription(),
							Value: permission.ID,
						})
					}
					uadminDatabase.Close()
					return ret
				}
				permissionsWidget.LeftSelectTitle = "Available user permissions"
				permissionsWidget.LeftSelectHelp = "This is the list of available user permissions. You may choose some by selecting them in the box below and then clicking the \"Choose\" arrow between the two boxes."
				permissionsWidget.LeftSearchSelectHelp = "Type into this box to filter down the list of available user permissions."
				permissionsWidget.LeftHelpChooseAll = "Click to choose all user permissions at once."
				permissionsWidget.RightSelectTitle = "Chosen user permissions"
				permissionsWidget.RightSelectHelp = "This is the list of chosen user permissions. You may remove some by selecting them in the box below and then clicking the \"Remove\" arrow between the two boxes."
				permissionsWidget.RightSearchSelectHelp = ""
				permissionsWidget.RightHelpChooseAll = "Click to remove all chosen user permissions at once."
				permissionsWidget.HelpText = "Specific permissions for this user. Hold down \"Control\", or \"Command\" on a Mac, to select more than one."
				permissionsWidget.PopulateRightSide = func()[]*interfaces.SelectOptGroup {
					ret := make([]*interfaces.SelectOptGroup, 0)
					user := modelI.(*interfaces.User)
					if user.ID != 0 {
						var permissions []*interfaces.Permission
						uadminDatabase := interfaces.NewUadminDatabase()
						uadminDatabase.Db.Model(user).Association("Permissions").Find(&permissions)
						ret = make([]*interfaces.SelectOptGroup, 0)
						for _, permission := range permissions {
							ret = append(ret, &interfaces.SelectOptGroup{
								OptLabel: permission.ShortDescription(),
								Value:    permission.ID,
							})
						}
						uadminDatabase.Close()
						return ret
					} else {
						formD := ctx.GetPostForm()
						if formD != nil {
							Ids := strings.Split(formD.Value["Permissions"][0], ",")
							IdI := make([]uint, 0)
							for _, tmp := range Ids {
								tmpI, _ := strconv.Atoi(tmp)
								IdI = append(IdI, uint(tmpI))
							}
							var permissions []*interfaces.Permission
							if len(IdI) > 0 {
								uadminDatabase := interfaces.NewUadminDatabase()
								uadminDatabase.Db.Preload("ContentType").Find(&permissions, IdI)
								ret = make([]*interfaces.SelectOptGroup, 0)
								for _, permission := range permissions {
									ret = append(ret, &interfaces.SelectOptGroup{
										OptLabel: permission.ShortDescription(),
										Value:    permission.ID,
									})
								}
								uadminDatabase.Close()
								return ret
							}
						}
					}
					return ret
				}
			}
			return form
		},
	)
	usermodelAdminPage.SaveModel = func(modelI interface{}, ID uint, afo interfaces.IAdminFilterObjects) interface{} {
		user := modelI.(*interfaces.User)
		if user.Salt == "" && user.Password != "" {
			user.Salt = utils.RandStringRunes(interfaces.CurrentConfig.D.Auth.SaltLength)
		}
		if ID != 0 {
			userM := &interfaces.User{}
			afo.GetUadminDatabase().Db.First(userM, ID)
			if userM.Password != user.Password && user.Password != "" {
				// hashedPassword, err := utils2.HashPass(password, salt)
				hashedPassword, _ := utils2.HashPass(user.Password, user.Salt)
				user.IsPasswordUsable = true
				user.Password = hashedPassword
			}
		} else {
			if user.Password != "" {
				// hashedPassword, err := utils2.HashPass(password, salt)
				hashedPassword, _ := utils2.HashPass(user.Password, user.Salt)
				user.Password = hashedPassword
				user.IsPasswordUsable = true
			}
		}
		afo.GetUadminDatabase().Db.Save(user)
		return user
	}
	usermodelAdminPage.PageName = "Users"
	usermodelAdminPage.Slug = "user"
	usermodelAdminPage.BlueprintName = "user"
	usermodelAdminPage.Router = mainRouter
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
	usergroupsAdminPage := interfaces.NewGormAdminPage(
		usersAdminPage,
		func() (interface{}, interface{}) {return &interfaces.UserGroup{}, &[]*interfaces.UserGroup{}},
		func(modelI interface{}, ctx interfaces.IAdminContext) *interfaces.Form {
			fields := []string{"GroupName"}
			if ctx.GetUserObject().IsSuperUser {
				fields = append(fields, "Permissions")
			}
			form := interfaces.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			if ctx.GetUserObject().IsSuperUser {
				permissionsField, _ := form.FieldRegistry.GetByName("Permissions")
				permissionsField.SetUpField = func(w interfaces.IWidget, modelI interface{}, v interface{}, afo interfaces.IAdminFilterObjects) error {
					model := modelI.(*interfaces.UserGroup)
					vTmp := v.([]string)
					var permission *interfaces.Permission
					if model.ID != 0 {
						afo.GetUadminDatabase().Db.Model(model).Association("Permissions").Clear()
						model.Permissions = make([]interfaces.Permission, 0)
					}
					for _, ID := range vTmp {
						afo.GetUadminDatabase().Db.First(&permission, ID)
						if permission.ID != 0 {
							model.Permissions = append(model.Permissions, *permission)
						}
						permission = nil
					}
					return nil
				}
				permissionsWidget := permissionsField.FieldConfig.Widget.(*interfaces.ChooseFromSelectWidget)
				permissionsWidget.PopulateLeftSide = func()[]*interfaces.SelectOptGroup {
					var permissions []*interfaces.Permission
					uadminDatabase := interfaces.NewUadminDatabase()
					uadminDatabase.Db.Preload("ContentType").Find(&permissions)
					ret := make([]*interfaces.SelectOptGroup, 0)
					for _, permission := range permissions {
						ret = append(ret, &interfaces.SelectOptGroup{
							OptLabel: permission.ShortDescription(),
							Value: permission.ID,
						})
					}
					uadminDatabase.Close()
					return ret
				}
				permissionsWidget.LeftSelectTitle = "Available permissions"
				permissionsWidget.LeftSelectHelp = "This is the list of available permissions. You may choose some by selecting them in the box below and then clicking the \"Choose\" arrow between the two boxes."
				permissionsWidget.LeftSearchSelectHelp = "Type into this box to filter down the list of available user permissions."
				permissionsWidget.LeftHelpChooseAll = "Click to choose all user permissions at once."
				permissionsWidget.RightSelectTitle = "Chosen permissions"
				permissionsWidget.RightSelectHelp = "This is the list of chosen permissions. You may remove some by selecting them in the box below and then clicking the \"Remove\" arrow between the two boxes."
				permissionsWidget.RightSearchSelectHelp = ""
				permissionsWidget.RightHelpChooseAll = "Click to remove all chosen permissions at once."
				permissionsWidget.HelpText = "Specific permissions for this user. Hold down \"Control\", or \"Command\" on a Mac, to select more than one."
				permissionsWidget.PopulateRightSide = func()[]*interfaces.SelectOptGroup {
					ret := make([]*interfaces.SelectOptGroup, 0)
					user := modelI.(*interfaces.UserGroup)
					if user.ID != 0 {
						var permissions []*interfaces.Permission
						uadminDatabase := interfaces.NewUadminDatabase()
						uadminDatabase.Db.Model(user).Association("Permissions").Find(&permissions)
						ret = make([]*interfaces.SelectOptGroup, 0)
						for _, permission := range permissions {
							ret = append(ret, &interfaces.SelectOptGroup{
								OptLabel: permission.ShortDescription(),
								Value:    permission.ID,
							})
						}
						uadminDatabase.Close()
						return ret
					} else {
						formD := ctx.GetPostForm()
						if formD != nil {
							Ids := strings.Split(formD.Value["Permissions"][0], ",")
							IdI := make([]uint, 0)
							for _, tmp := range Ids {
								tmpI, _ := strconv.Atoi(tmp)
								IdI = append(IdI, uint(tmpI))
							}
							var permissions []*interfaces.Permission
							if len(IdI) > 0 {
								uadminDatabase := interfaces.NewUadminDatabase()
								uadminDatabase.Db.Preload("ContentType").Find(&permissions, IdI)
								ret = make([]*interfaces.SelectOptGroup, 0)
								for _, permission := range permissions {
									ret = append(ret, &interfaces.SelectOptGroup{
										OptLabel: permission.ShortDescription(),
										Value:    permission.ID,
									})
								}
								uadminDatabase.Close()
								return ret
							}
						}
					}
					return ret
				}
			}
			return form
		},
	)
	usergroupsAdminPage.PageName = "User groups"
	usergroupsAdminPage.Slug = "usergroup"
	usergroupsAdminPage.BlueprintName = "user"
	usergroupsAdminPage.Router = mainRouter
	err = usersAdminPage.SubPages.AddAdminPage(usergroupsAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
}

func (b Blueprint) Init() {
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.OneTimeAction{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.User{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.UserGroup{}})
	interfaces.ProjectModels.RegisterModel(func() interface{}{return &interfaces.Permission{}})

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
	fsStorage := interfaces.NewFsStorage()
	interfaces.UadminFormCongirurableOptionInstance.AddFieldFormOptions(&interfaces.FieldFormOptions{
		WidgetType: "image",
		Name: "UserPhotoFormOptions",
		WidgetPopulate: func(m interface{}, currentField *interfaces.Field) interface{} {
			photo := m.(*interfaces.User).Photo
			if photo == "" {
				return ""
			}
			return fmt.Sprintf("%s%s", fsStorage.GetUploadUrl(), photo)
		},
	})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "user",
		Description:       "this blueprint is about users",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
