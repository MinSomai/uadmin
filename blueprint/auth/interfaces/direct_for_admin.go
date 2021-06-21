package interfaces

import (
	"github.com/asaskevich/govalidator"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	utils2 "github.com/uadmin/uadmin/blueprint/auth/utils"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	sessioninterfaces "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	user2 "github.com/uadmin/uadmin/blueprint/user"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"time"
)

// Binding from JSON
type LoginParamsForUadminAdmin struct {
	// SigninByField     string `form:"signinbyfield" json:"signinbyfield" xml:"signinbyfield"  binding:"required"`
	SigninField     string `form:"signinfield" json:"signinfield" xml:"signinfield"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	OTP string `form:"otp" json:"otp" xml:"otp" binding:"omitempty"`
}

type SignupParamsForUadminAdmin struct {
	Username string    `form:"username" json:"username" xml:"username"  binding:"required" valid:"username-unique"`
	Email string    `form:"email" json:"email" xml:"email"  binding:"required" valid:"email,email-unique"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
	ConfirmedPassword string `form:"confirm_password" json:"confirm_password" xml:"confirm_password" binding:"required"`
}

type DirectAuthForAdminProvider struct {
}

func (ap *DirectAuthForAdminProvider) GetUserFromRequest(c *gin.Context) *usermodels.User {
	session := ap.GetSession(c)
	if session != nil {
		spew.Dump(session)
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
	db := dialect.GetDB()
	var user usermodels.User
	db.Model(usermodels.User{}).Where(&usermodels.User{Username: json.SigninField}).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login credentials are incorrect."})
		return
	}
	if !user.Active {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this user is inactive"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password + user.Salt))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login credentials are incorrect."})
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter = sessionAdapter.Create()
	sessionAdapter.SetUser(&user)
	sessionDuration := time.Duration(config.CurrentConfig.D.Uadmin.SessionDuration)*time.Second
	sessionExpirationTime := time.Now().Add(sessionDuration)
	sessionAdapter.ExpiresOn(&sessionExpirationTime)
	sessionAdapter.Save()
	c.SetCookie(config.CurrentConfig.D.Uadmin.AdminCookieName, sessionAdapter.GetKey(), int(config.CurrentConfig.D.Uadmin.SessionDuration), "/", c.Request.URL.Host, config.CurrentConfig.D.Uadmin.SecureCookie, config.CurrentConfig.D.Uadmin.HttpOnlyCookie)
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
		Password: json.Password,
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
	db := dialect.GetDB()
	salt := utils.RandStringRunes(config.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass(json.Password, salt)
	user := usermodels.User{
		Username:     json.Username,
		Email: json.Email,
		Password:     hashedPassword,
		Active:       true,
		Salt: salt,
	}
	db.Create(&user)
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter = sessionAdapter.Create()
	sessionAdapter.SetUser(&user)
	sessionDuration := time.Duration(config.CurrentConfig.D.Uadmin.SessionDuration)*time.Second
	sessionExpirationTime := time.Now().Add(sessionDuration)
	sessionAdapter.ExpiresOn(&sessionExpirationTime)
	sessionAdapter.Save()
	c.SetCookie(config.CurrentConfig.D.Uadmin.AdminCookieName, sessionAdapter.GetKey(), int(config.CurrentConfig.D.Uadmin.SessionDuration), "/", c.Request.URL.Host, config.CurrentConfig.D.Uadmin.SecureCookie, config.CurrentConfig.D.Uadmin.HttpOnlyCookie)
	c.JSON(http.StatusOK, getUserForUadminPanel(sessionAdapter.GetUser()))
}

func (ap *DirectAuthForAdminProvider) Logout(c *gin.Context) {
	var cookie string
	var err error
	var cookieName string
	cookieName = config.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, err = c.Cookie(cookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ApiBadResponse(err.Error()))
		return
	}
	if cookie == "" {
		c.JSON(http.StatusBadRequest, utils.ApiBadResponse("empty cookie passed"))
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ApiBadResponse(err.Error()))
		return
	}
	sessionAdapter.Delete()
	timeInPast := time.Now().Add(-10*time.Minute)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(""),
		MaxAge:   0,
		Path:     "/",
		Domain:   c.Request.URL.Host,
		SameSite: http.SameSiteDefaultMode,
		Secure:   config.CurrentConfig.D.Uadmin.SecureCookie,
		HttpOnly: config.CurrentConfig.D.Uadmin.HttpOnlyCookie,
		Expires: timeInPast,
	})
	c.Status(http.StatusNoContent)
}

func (ap *DirectAuthForAdminProvider) IsAuthenticated(c *gin.Context) {
	var cookieName string
	cookieName = config.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ApiBadResponse(err.Error()))
		return
	}
	if cookie == "" {
		c.JSON(http.StatusBadRequest, utils.ApiBadResponse("empty cookie passed"))
		return
	}
	sessionAdapterRegistry := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry
	sessionAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	sessionAdapter, err = sessionAdapter.GetByKey(cookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ApiBadResponse(err.Error()))
		return
	}
	if !sessionAdapter.IsExpired() {
		c.JSON(http.StatusBadRequest, utils.ApiBadResponse("session expired"))
		return
	}
	c.JSON(http.StatusOK, getUserForUadminPanel(sessionAdapter.GetUser()))
}

func getUserForUadminPanel(user *usermodels.User) *gin.H {
	return &gin.H{"name": user.Username, "id": user.ID, "for-uadmin-panel": true}
}

func (ap *DirectAuthForAdminProvider) GetSession(c *gin.Context) sessioninterfaces.ISessionProvider {
	var cookieName string
	cookieName = config.CurrentConfig.D.Uadmin.AdminCookieName
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
