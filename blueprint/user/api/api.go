package api

import (
	"fmt"
	authservices "github.com/uadmin/uadmin/blueprint/auth/services"
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	usermodel "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/metrics"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	// "github.com/uadmin/uadmin/translation"
	"github.com/uadmin/uadmin/utils"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	c.JSON(http.StatusOK, utils.ApiSuccessResp())
}

var apiHandlers = []utils.ApiHandlerRequest{
	utils.ApiHandlerRequest{
		Methods: []string{"post"},
		Handler: CreateUser,
		Detail:  false,
	},
}

// profileHandler !
func ProfileHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	/*
		http://domain.com/admin/profile/
	*/
	r.ParseMultipartForm(32 << 20)
	type Context struct {
		User         string
		ID           uint
		Schema       model.ModelSchema
		Status       bool
		IsUpdated    bool
		Notif        string
		Demo         bool
		SiteName     string
		Language     langmodel.Language
		RootURL      string
		ProfilePhoto string
		OTPImage     string
		OTPRequired  bool
		Logo         string
		FavIcon      string
	}

	c := Context{}
	c.RootURL = preloaded.RootURL
	// @todo, redo
	// c.Language = translation.GetLanguage(r)
	c.SiteName = preloaded.SiteName
	user := session.User
	c.User = user.Username
	c.ProfilePhoto = session.User.Photo
	// c.OTPImage = "/media/otp/" + session.User.OTPSeed + ".png"
	c.Logo = preloaded.Logo
	c.FavIcon = preloaded.FavIcon

	// Check if OTP Required has been changed
	if r.URL.Query().Get("otp_required") != "" {
		//if r.URL.Query().Get("otp_required") == "1" {
		//	user.OTPRequired = true
		//} else if r.URL.Query().Get("otp_required") == "0" {
		//	user.OTPRequired = false
		//}
		r.URL.RawQuery = ""
		(&user).Save()
		// c.OTPImage = "/media/otp/" + user.OTPSeed + ".png"
	}

	// c.OTPRequired = user.OTPRequired

	c.Schema, _ = model.GetSchema(user)
	r.Form.Set("ModelID", fmt.Sprint(user.ID))
	// @todo probably, return
	// model.GetFormData(user, r, session, &c.Schema, &user)

	if r.Method == preloaded.CPOST {
		c.IsUpdated = true
		if r.FormValue("save") == "" {
			user.Username = r.FormValue("Username")
			user.FirstName = r.FormValue("FirstName")
			user.LastName = r.FormValue("LastName")
			user.Email = r.FormValue("Email")
			// @todo, redo
			// f := c.Schema.FieldByName("Photo")
			//if _, _, err := r.FormFile("Photo"); err == nil {
			//	user.Photo = imageapi.ProcessUpload(r, f, "user", session, &c.Schema)
			//}
			(&user).Save()
			c.ProfilePhoto = user.Photo
		}
		if r.FormValue("save") == "password" {
			// @todo, redo
			//oldPassword := r.FormValue("oldPassword")
			//newPassword := r.FormValue("newPassword")
			//confirmPassword := r.FormValue("confirmPassword")
			//_session := user.Login(oldPassword, "")
			//
			//if _session == nil || !user.Active {
			//	c.Status = true
			//	c.Notif = "Incorrent old password."
			//} else if newPassword != confirmPassword {
			//	c.Status = true
			//	c.Notif = "New password and confirm password do not match."
			//} else {
			//	user.Password = authservices.HashPass(newPassword)
			//	user.Save()
			//
			//	// To logout
			//	authapi.Logout(r)
			//
			//	return
			//}
		}
	}

	// @todo, redo
	// uadminhttp.RenderHTML(w, r, "./templates/uadmin/"+preloaded.Theme+"/profile.html", c)
}

func PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	type Context struct {
		Err       string
		ErrExists bool
		SiteName  string
		Language  langmodel.Language
		RootURL   string
		Logo      string
		FavIcon   string
	}

	c := Context{}
	c.SiteName = preloaded.SiteName
	c.RootURL = preloaded.RootURL
	// @todo, redo
	// c.Language = translation.GetLanguage(r)
	c.Logo = preloaded.Logo
	c.FavIcon = preloaded.FavIcon

	// Get the user and the code and verify them
	userID := r.FormValue("u")

	user := &usermodel.User{}
	database.Get(user, "id = ?", userID)
	if user.ID == 0 {
		go func() {
			log := &logmodel.Log{}
			r.Form.Set("reset-status", "invalid user id")
			log.PasswordReset(userID, log.Action.PasswordResetDenied(), r)
			log.Save()
		}()
		// @todo, redo
		// uadminhttp.PageErrorHandler(w, r, nil)
		return
	}
	otpCode := r.FormValue("key")
	if !user.VerifyOTP(otpCode) {
		go func() {
			log := &logmodel.Log{}
			r.Form.Set("reset-status", "invalid otp code: "+otpCode)
			log.PasswordReset(user.Username, log.Action.PasswordResetDenied(), r)
			log.Save()
		}()
		// @todo, redo
		// uadminhttp.PageErrorHandler(w, r, nil)
		return
	}

	if r.Method == preloaded.CPOST {
		if r.FormValue("password") != r.FormValue("confirm_password") {
			c.ErrExists = true
			c.Err = "Password does not match the confirm password"
		} else {
			user.Password = authservices.HashPass(r.FormValue("password"))
			user.Save()
			//log successful password reset
			go func() {
				log := &logmodel.Log{}
				r.Form.Set("reset-status", "Successfully changed the password")
				log.PasswordReset(user.Username, log.Action.PasswordResetSuccessful(), r)
				log.Save()
			}()
			http.Redirect(w, r, preloaded.RootURL, http.StatusSeeOther)
			return
		}
	}
	// @todo, redo
	// uadminhttp.RenderHTML(w, r, "./templates/uadmin/"+preloaded.Theme+"/resetpassword.html", c)
}

// logoutHandler !
func LogoutHandler(w http.ResponseWriter, r *http.Request, session *sessionmodel.Session) {
	// authapi.Logout(r)

	// Expire all cookies on logout
	for _, cookie := range r.Cookies() {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, preloaded.RootURL, http.StatusSeeOther)
}

// loginHandler HTTP handeler for verifying login data and creating sessions for users
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Err         string
		ErrExists   bool
		SiteName    string
		Languages   []langmodel.Language
		RootURL     string
		OTPRequired bool
		Language    langmodel.Language
		Username    string
		Password    string
		Logo        string
		FavIcon     string
	}

	c := Context{}
	c.SiteName = preloaded.SiteName
	c.RootURL = preloaded.RootURL
	// @todo, redo
	// c.Language = translation.GetLanguage(r)
	c.Logo = preloaded.Logo
	c.FavIcon = preloaded.FavIcon

	if r.Method == preloaded.CPOST {
		if r.FormValue("save") == "Send Request" {
			// This is a password reset request
			metrics.IncrementMetric("uadmin/security/passwordreset/request")
			email := r.FormValue("email")
			user := usermodel.User{}
			database.Get(&user, "Email = ?", email)
			if user.ID != 0 {
				metrics.IncrementMetric("uadmin/security/passwordreset/emailsent")
				c.ErrExists = true
				c.Err = "Password recovery request sent. Please check email to reset your password"
				forgotPasswordHandler(&user, r)
			} else {
				metrics.IncrementMetric("uadmin/security/passwordreset/invalidemail")
				c.ErrExists = true
				c.Err = "Please check email address. Email address must be associated with the account to be recovered."
			}
		} else {
			// This is a login request
			username := r.PostFormValue("username")
			username = strings.TrimSpace(strings.ToLower(username))
			//password := r.PostFormValue("password")
			//otp := r.PostFormValue("otp")
			//lang := r.PostFormValue("language")

			//session := authapi.Login2FA(r, username, password, otp)
			//if session == nil || !session.User.Active {
			//	c.ErrExists = true
			//	c.Err = "Invalid username/password or inactive user"
			//} else {
			//	if session.PendingOTP {
			//		utils.Trail(utils.INFO, "User: %s OTP: %s", session.User.Username, session.User.GetOTP())
			//	}
			//	cookie, _ := r.Cookie("session")
			//	if cookie == nil {
			//		cookie = &http.Cookie{}
			//	}
			//	cookie.Name = "session"
			//	cookie.Value = session.Key
			//	cookie.Path = "/"
			//	cookie.SameSite = http.SameSiteStrictMode
			//	http.SetCookie(w, cookie)
			//
			//	// set language cookie
			//	cookie, _ = r.Cookie("language")
			//	if cookie == nil {
			//		cookie = &http.Cookie{}
			//	}
			//	cookie.Name = "language"
			//	cookie.Value = lang
			//	cookie.Path = "/"
			//	http.SetCookie(w, cookie)
			//
			//	// Check for OTP
			//	if session.PendingOTP {
			//		c.Username = username
			//		c.Password = password
			//		c.OTPRequired = true
			//	} else {
			//		if r.URL.Query().Get("next") == "" {
			//			http.Redirect(w, r, strings.TrimSuffix(r.RequestURI, "logout"), http.StatusSeeOther)
			//			return
			//		}
			//		http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusSeeOther)
			//		return
			//	}
			//}
		}
	}
	c.Languages = langmodel.ActiveLangs
	// @todo, redo
	// uadminhttp.RenderHTML(w, r, "./templates/uadmin/"+preloaded.Theme+"/login.html", c)
}

// forgotPasswordHandler !
func forgotPasswordHandler(u *usermodel.User, r *http.Request) error {
	if u.Email == "" {
		return fmt.Errorf("unable to reset password, the user does not have an email")
	}
	msg := `Dear {NAME},

Have you forgotten your password to access {WEBSITE}. Don't worry we got your back. Please follow the link below to reset your password.

If you want to reset your password, click this link:
<a href="{URL}">{URL}</a>

If you didn't request a password reset, you can ignore this message.

Regards,
{WEBSITE} Support
`
	// Check if the host name is in the allowed hosts list
	allowed := false
	var host string
	var allowedHost string
	var err error
	if host, _, err = net.SplitHostPort(r.Host); err != nil {
		host = r.Host
	}
	for _, v := range strings.Split(preloaded.AllowedHosts, ",") {
		if allowedHost, _, err = net.SplitHostPort(v); err != nil {
			allowedHost = v
		}
		if allowedHost == host {
			allowed = true
		}
	}
	if !allowed {
		utils.Trail(utils.CRITICAL, "Reset password request for host: (%s) which is not in AllowedHosts settings", host)
		return nil
	}

	urlParts := strings.Split(r.Header.Get("origin"), "://")
	link := urlParts[0] + "://" + r.Host + preloaded.RootURL + "resetpassword?u=" + fmt.Sprint(u.ID) + "&key=" + u.GetOTP()
	msg = strings.Replace(msg, "{NAME}", u.String(), -1)
	msg = strings.Replace(msg, "{WEBSITE}", preloaded.SiteName, -1)
	msg = strings.Replace(msg, "{URL}", link, -1)
	subject := "Password reset for " + preloaded.SiteName
	err = utils.SendEmail([]string{u.Email}, []string{}, []string{}, subject, msg)

	return err
}
