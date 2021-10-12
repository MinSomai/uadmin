package user

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	utils2 "github.com/sergeyglazyrindev/uadmin/blueprint/auth/utils"
	sessionsblueprint "github.com/sergeyglazyrindev/uadmin/blueprint/sessions"
	"github.com/sergeyglazyrindev/uadmin/blueprint/user/migrations"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/sergeyglazyrindev/uadmin/utils"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Blueprint struct {
	core.Blueprint
}

type PasswordValidationStruct struct {
	Password          string `valid:"password-uadmin"`
	ConfirmedPassword string
}

type ForgotPasswordHandlerParams struct {
	Email string `form:"email" json:"email" xml:"email"  binding:"required" valid:"email"`
}

type ResetPasswordHandlerParams struct {
	Code              string `form:"code" json:"code" xml:"code"  binding:"required"`
	Password          string `form:"password" json:"password" xml:"password"  binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password"  binding:"required"`
}

type ChangePasswordHandlerParams struct {
	OldPassword       string `form:"old_password" json:"old_password" xml:"old_password"  binding:"required"`
	Password          string `form:"password" json:"password" xml:"password"  binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password"  binding:"required"`
}

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	mainRouter.GET("/reset-password/", func(ctx *gin.Context) {
		type Context struct {
			core.AdminContext
		}
		c := &Context{}
		core.PopulateTemplateContextForAdminPanel(ctx, c, core.NewAdminRequestParams())
		tr := core.NewTemplateRenderer("Reset Password")
		tr.Render(ctx, core.CurrentConfig.TemplatesFS, core.CurrentConfig.GetPathToTemplate("resetpassword"), c, core.FuncMap)
	})
	group.POST("/api/forgot/", func(ctx *gin.Context) {
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
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		user := core.GenerateUserModel()
		db.Model(core.GenerateUserModel()).Where(&core.User{Email: json.Email}).First(user)
		if user.GetID() == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with this email not found"})
			return
		}
		templateWriter := bytes.NewBuffer([]byte{})
		template1, err := template.ParseFS(core.CurrentConfig.TemplatesFS, core.CurrentConfig.GetPathToTemplate("email/forgot"))
		if err != nil {
			core.Trail(core.ERROR, "RenderHTML unable to parse %s. %s", core.CurrentConfig.GetPathToTemplate("email/forgot"), err)
			ctx.JSON(http.StatusBadRequest, utils.APIBadResponse(err.Error()))
			return
		}
		type Context struct {
			Name    string
			Website string
			URL     string
		}

		c := Context{}
		c.Name = user.GetUsername()
		c.Website = core.CurrentConfig.D.Uadmin.SiteName
		host := core.CurrentConfig.D.Uadmin.PoweredOnSite
		// @todo, generate code to restore access
		actionExpiresAt := time.Now()
		actionExpiresAt = actionExpiresAt.Add(time.Duration(core.CurrentConfig.D.Uadmin.ForgotCodeExpiration) * time.Minute)
		var oneTimeAction = core.OneTimeAction{
			User:       *user.(*core.User),
			ExpiresOn:  actionExpiresAt,
			Code:       utils.RandStringRunesForOneTimeAction(32),
			ActionType: 1,
		}

		db.Model(core.OneTimeAction{}).Save(&oneTimeAction)
		link := host + core.CurrentConfig.D.Uadmin.RootAdminURL + "/resetpassword/?key=" + oneTimeAction.Code
		c.URL = link
		err = template1.Execute(templateWriter, c)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.APIBadResponse(err.Error()))
			core.Trail(core.ERROR, "RenderHTML unable to parse %s. %s", core.CurrentConfig.GetPathToTemplate("email/forgot"), err)
			return
		}
		subject := "Password reset for admin panel on the " + core.CurrentConfig.D.Uadmin.SiteName
		err = utils.SendEmail(core.CurrentConfig.D.Uadmin.EmailFrom, []string{user.GetEmail()}, []string{}, []string{}, subject, templateWriter.String())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.APIBadResponse(err.Error()))
		}
		return
	})
	group.POST("/api/reset-password/", func(ctx *gin.Context) {
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
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var oneTimeAction core.OneTimeAction
		db.Model(core.OneTimeAction{}).Where(&core.OneTimeAction{Code: json.Code, IsUsed: false}).Preload("User").First(&oneTimeAction)
		if oneTimeAction.ID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No such code found"})
			return
		}
		if oneTimeAction.ExpiresOn.Before(time.Now()) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code is expired"})
			return
		}
		passwordValidationStruct := &PasswordValidationStruct{
			Password:          json.Password,
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
	group.POST("/api/change-password/", func(ctx *gin.Context) {
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
			Password:          json.Password,
			ConfirmedPassword: json.ConfirmedPassword,
		}
		_, err := govalidator.ValidateStruct(passwordValidationStruct)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		hashedPassword, err := utils2.HashPass(json.OldPassword, user.GetSalt())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// @todo, get it back once stabilize pass api
		//if hashedPassword != user.Password {
		//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password doesn't match current one"})
		//	return
		//}
		if user.GetSalt() == "" {
			user.SetSalt(utils.RandStringRunes(core.CurrentConfig.D.Auth.SaltLength))
		}
		hashedPassword, err = utils2.HashPass(json.Password, user.GetSalt())
		user.SetPassword(hashedPassword)
		user.SetIsPasswordUsable(true)
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		user1 := user.(*core.User)
		db.Save(user1)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	group.POST("/api/disable-2fa/", func(ctx *gin.Context) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		user.SetOTPRequired(false)
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	group.POST("/api/enable-2fa/", func(ctx *gin.Context) {
		sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
		var cookieName string
		cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
		cookie, _ := ctx.Cookie(cookieName)
		session, _ := sessionAdapter.GetByKey(cookie)
		user := session.GetUser()
		user.SetOTPRequired(true)
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})
	mainRouter.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.RequestURI, "/static-inbuilt/") || strings.HasSuffix(ctx.Request.RequestURI, ".css") ||
			strings.HasSuffix(ctx.Request.RequestURI, ".js") || strings.HasSuffix(ctx.Request.RequestURI, ".map") {
			ctx.Abort()
			return
		}
		// ctx.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		type Context struct {
			core.AdminContext
			Menu string
		}
		c := &Context{}
		core.PopulateTemplateContextForAdminPanel(ctx, c, core.NewAdminRequestParams())
		//
		//if r.Form.Get("err_msg") != "" {
		//	c.ErrMsg = r.Form.Get("err_msg")
		//}
		//if code, err := strconv.ParseUint(r.Form.Get("err_code"), 10, 16); err == nil {
		//	c.ErrCode = int(code)
		//}NewAdminPage
		ctx.Status(404)
		tr := core.NewTemplateRenderer("Page not found")
		tr.Render(ctx, core.CurrentConfig.TemplatesFS, core.CurrentConfig.GetPathToTemplate("404"), c, core.FuncMap)
	})
	usersAdminPage := core.NewGormAdminPage(
		nil,
		func() (interface{}, interface{}) { return nil, nil },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form { return nil },
	)
	usersAdminPage.PageName = "Users"
	usersAdminPage.Slug = "users"
	usersAdminPage.BlueprintName = "user"
	usersAdminPage.Router = mainRouter
	err := core.CurrentDashboardAdminPanel.AdminPages.AddAdminPage(usersAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
	var usermodelAdminPage *core.AdminPage
	usermodelAdminPage = core.NewGormAdminPage(
		usersAdminPage,
		func() (interface{}, interface{}) { return &core.User{}, &[]*core.User{} },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"Username", "FirstName", "LastName", "Email", "Active", "IsStaff", "IsSuperUser", "Password", "Photo", "LastLogin", "ExpiresOn"}
			if ctx.GetUserObject().GetIsSuperUser() {
				fields = append(fields, "UserGroups")
				fields = append(fields, "Permissions")
			}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			if ctx.GetUserObject().GetIsSuperUser() {
				usergroupsField, _ := form.FieldRegistry.GetByName("UserGroups")
				usergroupsField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
					model := modelI.(*core.User)
					vTmp := v.([]string)
					var usergroup *core.UserGroup
					if model.ID != 0 {
						afo.GetUadminDatabase().Db.Model(model).Association("UserGroups").Clear()
						model.UserGroups = make([]core.UserGroup, 0)
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
				userGroupsWidget := usergroupsField.FieldConfig.Widget.(*core.ChooseFromSelectWidget)
				userGroupsWidget.AddNewLink = fmt.Sprintf("%s/%s/usergroup/edit/%s?_to_field=id&_popup=1", core.CurrentConfig.D.Uadmin.RootAdminURL, usersAdminPage.Slug, "new")
				userGroupsWidget.AddNewTitle = "Add another group"
				userGroupsWidget.PopulateLeftSide = func() []*core.SelectOptGroup {
					var groups []*core.UserGroup
					uadminDatabase := core.NewUadminDatabase()
					uadminDatabase.Db.Find(&groups)
					ret := make([]*core.SelectOptGroup, 0)
					for _, group := range groups {
						ret = append(ret, &core.SelectOptGroup{
							OptLabel: group.GroupName,
							Value:    group.ID,
						})
					}
					uadminDatabase.Close()
					return ret
				}
				userGroupsWidget.PopulateRightSide = func() []*core.SelectOptGroup {
					ret := make([]*core.SelectOptGroup, 0)
					user := modelI.(*core.User)
					if user.ID != 0 {
						var groups []*core.UserGroup
						uadminDatabase := core.NewUadminDatabase()
						uadminDatabase.Db.Model(user).Association("UserGroups").Find(&groups)
						ret = make([]*core.SelectOptGroup, 0)
						for _, group := range groups {
							ret = append(ret, &core.SelectOptGroup{
								OptLabel: group.GroupName,
								Value:    group.ID,
							})
						}
						uadminDatabase.Close()
						return ret
					}
					formD := ctx.GetPostForm()
					if formD != nil {
						Ids := strings.Split(formD.Value["UserGroups"][0], ",")
						IDI := make([]uint, 0)
						for _, tmp := range Ids {
							tmpI, _ := strconv.Atoi(tmp)
							IDI = append(IDI, uint(tmpI))
						}
						if len(IDI) > 0 {
							var groups []*core.UserGroup
							uadminDatabase := core.NewUadminDatabase()
							uadminDatabase.Db.Find(&groups, IDI)
							ret = make([]*core.SelectOptGroup, 0)
							for _, group := range groups {
								ret = append(ret, &core.SelectOptGroup{
									OptLabel: group.GroupName,
									Value:    group.ID,
								})
							}
							uadminDatabase.Close()
							return ret
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
				permissionsField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
					model := modelI.(*core.User)
					vTmp := v.([]string)
					var permission *core.Permission
					if model.ID != 0 {
						afo.GetUadminDatabase().Db.Model(model).Association("Permissions").Clear()
						model.Permissions = make([]core.Permission, 0)
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
				permissionsWidget := permissionsField.FieldConfig.Widget.(*core.ChooseFromSelectWidget)
				permissionsWidget.PopulateLeftSide = func() []*core.SelectOptGroup {
					var permissions []*core.Permission
					uadminDatabase := core.NewUadminDatabase()
					uadminDatabase.Db.Preload("ContentType").Find(&permissions)
					ret := make([]*core.SelectOptGroup, 0)
					for _, permission := range permissions {
						ret = append(ret, &core.SelectOptGroup{
							OptLabel: permission.ShortDescription(),
							Value:    permission.ID,
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
				permissionsWidget.PopulateRightSide = func() []*core.SelectOptGroup {
					ret := make([]*core.SelectOptGroup, 0)
					user := modelI.(*core.User)
					if user.ID != 0 {
						var permissions []*core.Permission
						uadminDatabase := core.NewUadminDatabase()
						uadminDatabase.Db.Model(user).Association("Permissions").Find(&permissions)
						ret = make([]*core.SelectOptGroup, 0)
						for _, permission := range permissions {
							ret = append(ret, &core.SelectOptGroup{
								OptLabel: permission.ShortDescription(),
								Value:    permission.ID,
							})
						}
						uadminDatabase.Close()
						return ret
					}
					formD := ctx.GetPostForm()
					if formD != nil {
						Ids := strings.Split(formD.Value["Permissions"][0], ",")
						IDI := make([]uint, 0)
						for _, tmp := range Ids {
							tmpI, _ := strconv.Atoi(tmp)
							IDI = append(IDI, uint(tmpI))
						}
						var permissions []*core.Permission
						if len(IDI) > 0 {
							uadminDatabase := core.NewUadminDatabase()
							uadminDatabase.Db.Preload("ContentType").Find(&permissions, IDI)
							ret = make([]*core.SelectOptGroup, 0)
							for _, permission := range permissions {
								ret = append(ret, &core.SelectOptGroup{
									OptLabel: permission.ShortDescription(),
									Value:    permission.ID,
								})
							}
							uadminDatabase.Close()
							return ret
						}
					}
					return ret
				}
			}
			passwordField, _ := form.FieldRegistry.GetByName("Password")
			passwordField.SetUpField = func(w core.IWidget, m interface{}, v interface{}, afo core.IAdminFilterObjects) error {
				user := m.(*core.User)
				vI, _ := v.(string)
				if vI != "" {
					if user.Salt == "" {
						user.Salt = utils.RandStringRunes(core.CurrentConfig.D.Auth.SaltLength)
					}
					hashedPassword, _ := utils2.HashPass(vI, user.Salt)
					user.IsPasswordUsable = true
					user.Password = hashedPassword
				}

				return nil
			}
			return form
		},
	)
	usermodelAdminPage.PageName = "Users"
	usermodelAdminPage.Slug = "user"
	usermodelAdminPage.BlueprintName = "user"
	usermodelAdminPage.Router = mainRouter
	listFilter := &core.ListFilter{
		URLFilteringParam: "IsSuperUser__exact",
		Title:             "Is super user ?",
	}
	listFilter.OptionsToShow = append(listFilter.OptionsToShow, &core.FieldChoice{DisplayAs: "Yes", Value: true})
	listFilter.OptionsToShow = append(listFilter.OptionsToShow, &core.FieldChoice{DisplayAs: "No", Value: false})
	usermodelAdminPage.ListFilter.Add(listFilter)
	err = usersAdminPage.SubPages.AddAdminPage(usermodelAdminPage)
	if err != nil {
		panic(fmt.Errorf("error initializing user blueprint: %s", err))
	}
	usergroupsAdminPage := core.NewGormAdminPage(
		usersAdminPage,
		func() (interface{}, interface{}) { return &core.UserGroup{}, &[]*core.UserGroup{} },
		func(modelI interface{}, ctx core.IAdminContext) *core.Form {
			fields := []string{"GroupName"}
			if ctx.GetUserObject().GetIsSuperUser() {
				fields = append(fields, "Permissions")
			}
			form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
			if ctx.GetUserObject().GetIsSuperUser() {
				permissionsField, _ := form.FieldRegistry.GetByName("Permissions")
				permissionsField.SetUpField = func(w core.IWidget, modelI interface{}, v interface{}, afo core.IAdminFilterObjects) error {
					model := modelI.(*core.UserGroup)
					vTmp := v.([]string)
					var permission *core.Permission
					if model.ID != 0 {
						afo.GetUadminDatabase().Db.Model(model).Association("Permissions").Clear()
						model.Permissions = make([]core.Permission, 0)
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
				permissionsWidget := permissionsField.FieldConfig.Widget.(*core.ChooseFromSelectWidget)
				permissionsWidget.PopulateLeftSide = func() []*core.SelectOptGroup {
					var permissions []*core.Permission
					uadminDatabase := core.NewUadminDatabase()
					uadminDatabase.Db.Preload("ContentType").Find(&permissions)
					ret := make([]*core.SelectOptGroup, 0)
					for _, permission := range permissions {
						ret = append(ret, &core.SelectOptGroup{
							OptLabel: permission.ShortDescription(),
							Value:    permission.ID,
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
				permissionsWidget.PopulateRightSide = func() []*core.SelectOptGroup {
					ret := make([]*core.SelectOptGroup, 0)
					user := modelI.(*core.UserGroup)
					if user.ID != 0 {
						var permissions []*core.Permission
						uadminDatabase := core.NewUadminDatabase()
						uadminDatabase.Db.Model(user).Association("Permissions").Find(&permissions)
						ret = make([]*core.SelectOptGroup, 0)
						for _, permission := range permissions {
							ret = append(ret, &core.SelectOptGroup{
								OptLabel: permission.ShortDescription(),
								Value:    permission.ID,
							})
						}
						uadminDatabase.Close()
						return ret
					}
					formD := ctx.GetPostForm()
					if formD != nil {
						Ids := strings.Split(formD.Value["Permissions"][0], ",")
						IDI := make([]uint, 0)
						for _, tmp := range Ids {
							tmpI, _ := strconv.Atoi(tmp)
							IDI = append(IDI, uint(tmpI))
						}
						var permissions []*core.Permission
						if len(IDI) > 0 {
							uadminDatabase := core.NewUadminDatabase()
							uadminDatabase.Db.Preload("ContentType").Find(&permissions, IDI)
							ret = make([]*core.SelectOptGroup, 0)
							for _, permission := range permissions {
								ret = append(ret, &core.SelectOptGroup{
									OptLabel: permission.ShortDescription(),
									Value:    permission.ID,
								})
							}
							uadminDatabase.Close()
							return ret
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
	core.ProjectModels.RegisterModel(func() interface{} { return &core.OneTimeAction{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &core.User{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &core.UserGroup{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &core.Permission{} })

	core.UadminValidatorRegistry.AddValidator("username-unique", func(i interface{}, o interface{}) error {
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var cUsers int64
		db.Model(&core.User{}).Where(&core.User{Username: i.(string)}).Count(&cUsers)
		if cUsers == 0 {
			return nil
		}
		return fmt.Errorf("user with name %s is already registered", i.(string))
	})

	core.UadminValidatorRegistry.AddValidator("email-unique", func(i interface{}, o interface{}) error {
		uadminDatabase := core.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var cUsers int64
		db.Model(&core.User{}).Where(&core.User{Email: i.(string)}).Count(&cUsers)
		if cUsers == 0 {
			return nil
		}
		return fmt.Errorf("user with email %s is already registered", i.(string))
	})

	core.UadminValidatorRegistry.AddValidator("username-uadmin", func(i interface{}, o interface{}) error {
		minLength := core.CurrentConfig.D.Auth.MinUsernameLength
		maxLength := core.CurrentConfig.D.Auth.MaxUsernameLength
		currentUsername := i.(string)
		if maxLength < len(currentUsername) || len(currentUsername) < minLength {
			return fmt.Errorf("length of the username has to be between %d and %d symbols", minLength, maxLength)
		}
		return nil
	})

	core.UadminValidatorRegistry.AddValidator("password-uadmin", func(i interface{}, o interface{}) error {
		passwordStruct := o.(PasswordValidationStruct)
		if passwordStruct.Password != passwordStruct.ConfirmedPassword {
			return fmt.Errorf("password doesn't equal to confirmed password")
		}
		if len(passwordStruct.Password) < core.CurrentConfig.D.Auth.MinPasswordLength {
			return fmt.Errorf("length of the password has to be at least %d symbols", core.CurrentConfig.D.Auth.MinPasswordLength)
		}
		return nil
	})

	govalidator.CustomTypeTagMap.Set("username-unique", func(i interface{}, o interface{}) bool {
		validator, _ := core.UadminValidatorRegistry.GetValidator("username-unique")
		userExists := validator(i, o)
		return userExists == nil
	})
	govalidator.CustomTypeTagMap.Set("email-unique", func(i interface{}, o interface{}) bool {
		validator, _ := core.UadminValidatorRegistry.GetValidator("email-unique")
		emailExists := validator(i, o)
		return emailExists == nil
	})
	govalidator.CustomTypeTagMap.Set("username-uadmin", func(i interface{}, o interface{}) bool {
		validator, _ := core.UadminValidatorRegistry.GetValidator("username-uadmin")
		isValidUsername := validator(i, o)
		return isValidUsername == nil
	})
	govalidator.CustomTypeTagMap.Set("password-uadmin", func(i interface{}, o interface{}) bool {
		validator, _ := core.UadminValidatorRegistry.GetValidator("password-uadmin")
		isValidPassword := validator(i, o)
		return isValidPassword == nil
	})
	fsStorage := core.NewFsStorage()
	core.UadminFormCongirurableOptionInstance.AddFieldFormOptions(&core.FieldFormOptions{
		WidgetType: "image",
		Name:       "UserPhotoFormOptions",
		WidgetPopulate: func(renderContext *core.FormRenderContext, currentField *core.Field) interface{} {
			photo := renderContext.Model.(*core.User).Photo
			if photo == "" {
				return ""
			}
			return fmt.Sprintf("%s%s", fsStorage.GetUploadURL(), photo)
		},
	})
}

var ConcreteBlueprint = Blueprint{
	core.Blueprint{
		Name:              "user",
		Description:       "this blueprint is about users",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
