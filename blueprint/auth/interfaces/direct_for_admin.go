package interfaces

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	utils2 "github.com/uadmin/uadmin/blueprint/auth/utils"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	sessioninterfaces "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	user2 "github.com/uadmin/uadmin/blueprint/user"
	"github.com/uadmin/uadmin/core"
	"github.com/uadmin/uadmin/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"time"
)

// Binding from JSON
type LoginParamsForUadminAdmin struct {
	// SigninByField     string `form:"signinbyfield" json:"signinbyfield" xml:"signinbyfield"  binding:"required"`
	SigninField string `form:"signinfield" json:"signinfield" xml:"signinfield"  binding:"required"`
	Password    string `form:"password" json:"password" xml:"password" binding:"required"`
	OTP         string `form:"otp" json:"otp" xml:"otp" binding:"omitempty"`
}

type SignupParamsForUadminAdmin struct {
	Username          string `form:"username" json:"username" xml:"username"  binding:"required" valid:"username-unique"`
	Email             string `form:"email" json:"email" xml:"email"  binding:"required" valid:"email,email-unique"`
	Password          string `form:"password" json:"password" xml:"password" binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password" binding:"required"`
}

type DirectAuthForAdminProvider struct {
}

func (ap *DirectAuthForAdminProvider) GetUserFromRequest(c *gin.Context) *core.User {
	session := ap.GetSession(c)
	if session != nil {
		return session.GetUser()
	}
	return nil
}

func (ap *DirectAuthForAdminProvider) Signin(c *gin.Context) {
	var json LoginParamsForUadminAdmin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	var user core.User
	db.Model(core.User{}).Where(&core.User{Username: json.SigninField}).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login credentials are incorrect."})
		return
	}
	if !user.Active {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this user is inactive"})
		return
	}
	if !user.IsPasswordUsable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this user doesn't have a password"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password+user.Salt))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login credentials are incorrect."})
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	cookieName := core.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, err := c.Cookie(cookieName)
	sessionDuration := time.Duration(core.CurrentConfig.D.Uadmin.SessionDuration) * time.Second
	sessionExpirationTime := time.Now().Add(sessionDuration)
	if cookie != "" {
		sessionAdapter, _ = sessionAdapter.GetByKey(cookie)
		if sessionAdapter.IsExpired() {
			sessionAdapter.ExpiresOn(&sessionExpirationTime)
		}
	} else {
		sessionAdapter = sessionAdapter.Create()
		sessionAdapter.ExpiresOn(&sessionExpirationTime)
		c.SetCookie(core.CurrentConfig.D.Uadmin.AdminCookieName, sessionAdapter.GetKey(), int(core.CurrentConfig.D.Uadmin.SessionDuration), "/", c.Request.URL.Host, core.CurrentConfig.D.Uadmin.SecureCookie, core.CurrentConfig.D.Uadmin.HTTPOnlyCookie)
	}
	sessionAdapter.SetUser(&user)
	sessionAdapter.Save()
	c.JSON(http.StatusOK, getUserForUadminPanel(sessionAdapter.GetUser()))
}

func (ap *DirectAuthForAdminProvider) Signup(c *gin.Context) {
	var json SignupParamsForUadminAdmin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := govalidator.ValidateStruct(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	passwordValidationStruct := &user2.PasswordValidationStruct{
		Password:          json.Password,
		ConfirmedPassword: json.ConfirmedPassword,
	}
	_, err = govalidator.ValidateStruct(passwordValidationStruct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//if utils.Contains(config.CurrentConfig.D.Auth.Twofactor_auth_required_for_signin_adapters, ap.GetName()) {
	//	if json.OTP == "" {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "otp is required"})
	//		return
	//	}
	//	if user.GeneratedOTPToVerify != json.OTP {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "otp provided by user is wrong"})
	//		return
	//	}
	//	user.GeneratedOTPToVerify = ""
	//	db.Save(&user)
	//}
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	salt := utils.RandStringRunes(core.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass(json.Password, salt)
	user := core.User{
		Username:         json.Username,
		Email:            json.Email,
		Password:         hashedPassword,
		Active:           true,
		Salt:             salt,
		IsPasswordUsable: true,
	}
	db.Create(&user)
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter = sessionAdapter.Create()
	sessionAdapter.SetUser(&user)
	sessionDuration := time.Duration(core.CurrentConfig.D.Uadmin.SessionDuration) * time.Second
	sessionExpirationTime := time.Now().Add(sessionDuration)
	sessionAdapter.ExpiresOn(&sessionExpirationTime)
	sessionAdapter.Save()
	c.SetCookie(core.CurrentConfig.D.Uadmin.AdminCookieName, sessionAdapter.GetKey(), int(core.CurrentConfig.D.Uadmin.SessionDuration), "/", c.Request.URL.Host, core.CurrentConfig.D.Uadmin.SecureCookie, core.CurrentConfig.D.Uadmin.HTTPOnlyCookie)
	c.JSON(http.StatusOK, getUserForUadminPanel(sessionAdapter.GetUser()))
}

func (ap *DirectAuthForAdminProvider) Logout(c *gin.Context) {
	var cookie string
	var err error
	var cookieName string
	cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, err = c.Cookie(cookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.APIBadResponse(err.Error()))
		return
	}
	if cookie == "" {
		c.JSON(http.StatusBadRequest, utils.APIBadResponse("empty cookie passed"))
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.APIBadResponse(err.Error()))
		return
	}
	sessionAdapter.Delete()
	timeInPast := time.Now().Add(-10 * time.Minute)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(""),
		MaxAge:   0,
		Path:     "/",
		Domain:   c.Request.URL.Host,
		SameSite: http.SameSiteDefaultMode,
		Secure:   core.CurrentConfig.D.Uadmin.SecureCookie,
		HttpOnly: core.CurrentConfig.D.Uadmin.HTTPOnlyCookie,
		Expires:  timeInPast,
	})
	c.Status(http.StatusNoContent)
}

func (ap *DirectAuthForAdminProvider) IsAuthenticated(c *gin.Context) {
	var cookieName string
	cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.APIBadResponse(err.Error()))
		return
	}
	if cookie == "" {
		c.JSON(http.StatusBadRequest, utils.APIBadResponse("empty cookie passed"))
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.APIBadResponse(err.Error()))
		return
	}
	if !sessionAdapter.IsExpired() {
		c.JSON(http.StatusBadRequest, utils.APIBadResponse("session expired"))
		return
	}
	c.JSON(http.StatusOK, getUserForUadminPanel(sessionAdapter.GetUser()))
}

func getUserForUadminPanel(user *core.User) *gin.H {
	return &gin.H{"name": user.Username, "id": user.ID, "for-uadmin-panel": true}
}

func (ap *DirectAuthForAdminProvider) GetSession(c *gin.Context) sessioninterfaces.ISessionProvider {
	var cookieName string
	cookieName = core.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return nil
	}
	if cookie == "" {
		return nil
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err != nil {
		return nil
	}
	if !sessionAdapter.IsExpired() {
		return nil
	}
	return sessionAdapter
}

func (ap *DirectAuthForAdminProvider) GetName() string {
	return "direct-for-admin"
}
