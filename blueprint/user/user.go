package user

import (
	"bytes"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	utils2 "github.com/uadmin/uadmin/blueprint/auth/utils"
	"github.com/uadmin/uadmin/blueprint/user/migrations"
	"github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/debug"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/templatecontext"
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

func (b Blueprint) InitRouter(mainRouter *gin.Engine, group *gin.RouterGroup) {
	mainRouter.GET("/reset-password", func(ctx *gin.Context) {
		type Context struct {
			templatecontext.AdminContext
		}
		c := &Context{}
		templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, &templatecontext.AdminRequestParams{
			CreateSession: true, GenerateCSRFToken: true,
		})
		tr := utils.NewTemplateRenderer("Reset Password")
		tr.Render(ctx, config.CurrentConfig.TemplatesFS, config.CurrentConfig.GetPathToTemplate("resetpassword"), c)
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
		db := dialect.GetDB()
		var user models.User
		db.Model(models.User{}).Where(&models.User{Email: json.Email}).First(&user)
		if user.ID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with this email not found"})
			return
		}
		templateWriter := bytes.NewBuffer([]byte{})
		template1, err := template.ParseFS(config.CurrentConfig.TemplatesFS, config.CurrentConfig.GetPathToTemplate("email/forgot"))
		if err != nil {
			debug.Trail(debug.ERROR, "RenderHTML unable to parse %s. %s", config.CurrentConfig.GetPathToTemplate("email/forgot"), err)
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
		c.Website = config.CurrentConfig.D.Uadmin.SiteName
		host := config.CurrentConfig.D.Uadmin.PoweredOnSite
		// @todo, generate code to restore access
		actionExpiresAt := time.Now()
		actionExpiresAt = actionExpiresAt.Add(time.Duration(config.CurrentConfig.D.Uadmin.ForgotCodeExpiration)*time.Minute)
		var oneTimeAction =models.OneTimeAction{
			User:       user,
			ExpiresOn: &actionExpiresAt,
			Code: utils.RandStringRunesForOneTimeAction(32),
			ActionType: 1,

		}

		db.Model(models.OneTimeAction{}).Save(&oneTimeAction)
		link := host + config.CurrentConfig.D.Uadmin.RootAdminURL + "/resetpassword?key=" + oneTimeAction.Code
		c.URL = link
		err = template1.Execute(templateWriter, c)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ApiBadResponse(err.Error()))
			debug.Trail(debug.ERROR, "RenderHTML unable to parse %s. %s", config.CurrentConfig.GetPathToTemplate("email/forgot"), err)
			return
		}
		subject := "Password reset for admin panel on the " + config.CurrentConfig.D.Uadmin.SiteName
		err = utils.SendEmail(config.CurrentConfig.D.Uadmin.EmailFrom, []string{user.Email}, []string{}, []string{}, subject, templateWriter.String())
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
		db := dialect.GetDB()
		var oneTimeAction models.OneTimeAction
		db.Model(models.OneTimeAction{}).Where(&models.OneTimeAction{Code: json.Code, IsUsed: false}).Preload("User").First(&oneTimeAction)
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
		oneTimeAction.IsUsed = true
		db.Save(&oneTimeAction.User)
		db.Save(&oneTimeAction)
	})
	mainRouter.NoRoute(func(ctx *gin.Context) {
		// ctx.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		type Context struct {
			templatecontext.AdminContext
			Demo     bool
			Menu     string
		}

		c := &Context{}
		templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, &templatecontext.AdminRequestParams{
			CreateSession: true, GenerateCSRFToken: true,
		})
		//
		//if r.Form.Get("err_msg") != "" {
		//	c.ErrMsg = r.Form.Get("err_msg")
		//}
		//if code, err := strconv.ParseUint(r.Form.Get("err_code"), 10, 16); err == nil {
		//	c.ErrCode = int(code)
		//}
		ctx.Status(404)
		tr := utils.NewTemplateRenderer("Page not found")
		tr.Render(ctx, config.CurrentConfig.TemplatesFS, config.CurrentConfig.GetPathToTemplate("404"), c)
	})
}


func (b Blueprint) Init(config *config.UadminConfig) {
	govalidator.CustomTypeTagMap.Set("username-unique", func(i interface{}, o interface{}) bool {
		db := dialect.GetDB()
		var cUsers int64
		db.Model(&models.User{}).Where(&models.User{Username: i.(string)}).Count(&cUsers)
		return cUsers == 0
	})
	govalidator.CustomTypeTagMap.Set("email-unique", func(i interface{}, o interface{}) bool {
		db := dialect.GetDB()
		var cUsers int64
		db.Model(&models.User{}).Where(&models.User{Email: i.(string)}).Count(&cUsers)
		return cUsers == 0
	})
	govalidator.CustomTypeTagMap.Set("username-uadmin", func(i interface{}, o interface{}) bool {
		minLength := config.D.Auth.MinUsernameLength
		maxLength := config.D.Auth.MaxUsernameLength
		currentUsername := i.(string)
		if maxLength < len(currentUsername) || len(currentUsername) < minLength {
			return false
		}
		return true
	})
	govalidator.CustomTypeTagMap.Set("password-uadmin", func(i interface{}, o interface{}) bool {
		passwordStruct := o.(PasswordValidationStruct)
		if passwordStruct.Password != passwordStruct.ConfirmedPassword {
			return false
		}
		if len(passwordStruct.Password) < config.D.Auth.MinPasswordLength {
			return false
		}
		return true
	})
}

var ConcreteBlueprint = Blueprint{
	interfaces.Blueprint{
		Name:              "user",
		Description:       "this blueprint is about users",
		MigrationRegistry: migrations.BMigrationRegistry,
	},
}
