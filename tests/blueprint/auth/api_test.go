package auth

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	utils2 "github.com/uadmin/uadmin/blueprint/auth/utils"
	"github.com/uadmin/uadmin/blueprint/otp/services"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type AuthProviderTestSuite struct {
	uadmin.UadminTestSuite
}

func (s *AuthProviderTestSuite) TestDirectAuthProviderForUadminAdmin() {
	req, _ := http.NewRequest("GET", "/auth/direct/status/?for-uadmin-panel=1", nil)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "http: named cookie not present")
		return strings.Contains(w.Body.String(), "http: named cookie not present")
	})
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.AdminCookieName, ""),
	)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "empty cookie passed")
		return strings.Contains(w.Body.String(), "empty cookie passed")
	})
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.AdminCookieName, "test"),
	)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "no session with key test found")
		return strings.Contains(w.Body.String(), "no session with key test found")
	})
	sessionsblueprint1, _ := s.App.BlueprintRegistry.GetByName("sessions")
	sessionAdapterRegistry := sessionsblueprint1.(sessionsblueprint.Blueprint).SessionAdapterRegistry
	defaultAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	defaultAdapter = defaultAdapter.Create()
	expiresOn := time.Now().Add(-5*time.Minute)
	defaultAdapter.ExpiresOn(&expiresOn)
	defaultAdapter.Save()
	// directProvider.
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.AdminCookieName, defaultAdapter.GetKey()),
	)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "session expired")
		return strings.Contains(w.Body.String(), "session expired")
	})
	expiresOn = time.Now()
	expiresOn = expiresOn.Add(10*time.Minute)
	defaultAdapter.ExpiresOn(&expiresOn)
	defaultAdapter.Save()
	req.URL = &url.URL{
		Path:"/auth/direct/status/",
		RawQuery: "for-uadmin-panel=1",
	}
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "for-uadmin-panel")
		return strings.Contains(w.Body.String(), "for-uadmin-panel")
	})
	var jsonStr = []byte(`{"username":"test", "password": "123456"}`)
	req, _ = http.NewRequest("POST", "/auth/direct/signin/?for-uadmin-panel=1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "login credentials are incorrect")
		return strings.Contains(w.Body.String(), "login credentials are incorrect")
	})
	salt := utils.RandStringRunes(config.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass("123456", salt)
	user := usermodels.User{
		FirstName:    "testuser-firstname",
		LastName:     "testuser-lastname",
		Username:     "test",
		Password:     hashedPassword,
		Active:       false,
		Salt: salt,
	}
	db := dialect.GetDB()
	db.Create(&user)
	req, _ = http.NewRequest("POST", "/auth/direct/signin/?for-uadmin-panel=1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "this user is inactive")
		return strings.Contains(w.Body.String(), "this user is inactive")
	})
	user.Active = true
	secretString, _ := services.GenerateOTPSeed(config.CurrentConfig.D.Uadmin.OTPDigits, config.CurrentConfig.D.Uadmin.OTPAlgorithm, config.CurrentConfig.D.Uadmin.OTPSkew, config.CurrentConfig.D.Uadmin.OTPPeriod, &user)
	user.OTPSeed = secretString
	otpPassword := services.GetOTP(user.OTPSeed, config.CurrentConfig.D.Uadmin.OTPDigits, config.CurrentConfig.D.Uadmin.OTPAlgorithm, config.CurrentConfig.D.Uadmin.OTPSkew, config.CurrentConfig.D.Uadmin.OTPPeriod)
	user.GeneratedOTPToVerify = otpPassword
	var jsonStrForSignup = []byte(fmt.Sprintf(`{"username":"test", "password": "123456", "otp": "%s"}`, otpPassword))
	db.Save(&user)
	req, _ = http.NewRequest("POST", "/auth/direct/signin/?for-uadmin-panel=1", bytes.NewBuffer(jsonStrForSignup))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Header().Get("Set-Cookie"), "uadmin-admin=")
		sessionKey := strings.Split(strings.Split(w.Header().Get("Set-Cookie"), ";")[0], "=")[1]
		req1, _ := http.NewRequest("GET", "/auth/direct/status/?for-uadmin-panel=1", nil)
		req1.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.AdminCookieName, sessionKey),
		)
		uadmin.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			assert.Contains(s.T(), w.Body.String(), "for-uadmin-panel")
			return strings.Contains(w.Body.String(), "for-uadmin-panel")
		})
		return strings.Contains(w.Header().Get("Set-Cookie"), "uadmin-admin=")
	})
}

func (s *AuthProviderTestSuite) TestSignupForUadminAdmin() {
	// hashedPassword, err := utils2.HashPass(password, salt)
	var jsonStr = []byte(`{"username":"test", "confirm_password": "12345678", "password": "12345678", "email": "uadmin@example.com"}`)
	req, _ := http.NewRequest("POST", "/auth/direct/signup/?for-uadmin-panel=1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Header().Get("Set-Cookie"), "uadmin-admin=")
		return strings.Contains(w.Header().Get("Set-Cookie"), "uadmin-admin=")
	})
}

func (s *AuthProviderTestSuite) TestSignupForApi() {
	// hashedPassword, err := utils2.HashPass(password, salt)
	var jsonStr = []byte(`{"username":"test", "confirm_password": "12345678", "password": "12345678", "email": "uadmin@example.com"}`)
	req, _ := http.NewRequest("POST", "/auth/direct/signup/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
}

func (s *AuthProviderTestSuite) TestDirectAuthProviderForApi() {
	req, _ := http.NewRequest("GET", "/auth/direct/status/", nil)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "http: named cookie not present")
		return strings.Contains(w.Body.String(), "http: named cookie not present")
	})
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.ApiCookieName, ""),
	)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "empty cookie passed")
		return strings.Contains(w.Body.String(), "empty cookie passed")
	})
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.ApiCookieName, "test"),
	)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "no session with key test found")
		return strings.Contains(w.Body.String(), "no session with key test found")
	})
	sessionsblueprint1, _ := s.App.BlueprintRegistry.GetByName("sessions")
	sessionAdapterRegistry := sessionsblueprint1.(sessionsblueprint.Blueprint).SessionAdapterRegistry
	defaultAdapter, _ := sessionAdapterRegistry.GetDefaultAdapter()
	defaultAdapter = defaultAdapter.Create()
	expiresOn := time.Now().Add(-5*time.Minute)
	defaultAdapter.ExpiresOn(&expiresOn)
	defaultAdapter.Save()
	// directProvider.
	req.Header.Set(
		"Cookie",
		fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.ApiCookieName, defaultAdapter.GetKey()),
	)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "session expired")
		return strings.Contains(w.Body.String(), "session expired")
	})
	expiresOn = time.Now()
	expiresOn = expiresOn.Add(10*time.Minute)
	defaultAdapter.ExpiresOn(&expiresOn)
	defaultAdapter.Save()
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "\"id\":0")
		return strings.Contains(w.Body.String(), "\"id\":0")
	})
	var jsonStr = []byte(`{"username":"test", "password": "123456"}`)
	req, _ = http.NewRequest("POST", "/auth/direct/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "login credentials are incorrect")
		return strings.Contains(w.Body.String(), "login credentials are incorrect")
	})
	salt := utils.RandStringRunes(config.CurrentConfig.D.Auth.SaltLength)
	// hashedPassword, err := utils2.HashPass(password, salt)
	hashedPassword, _ := utils2.HashPass("123456", salt)
	user := usermodels.User{
		FirstName:    "testuser-firstname",
		LastName:     "testuser-lastname",
		Username:     "test",
		Password:     hashedPassword,
		Active:       false,
		Salt: salt,
	}
	db := dialect.GetDB()
	db.Create(&user)
	req, _ = http.NewRequest("POST", "/auth/direct/signin/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Contains(s.T(), w.Body.String(), "this user is inactive")
		return strings.Contains(w.Body.String(), "this user is inactive")
	})
	user.Active = true
	secretString, _ := services.GenerateOTPSeed(config.CurrentConfig.D.Uadmin.OTPDigits, config.CurrentConfig.D.Uadmin.OTPAlgorithm, config.CurrentConfig.D.Uadmin.OTPSkew, config.CurrentConfig.D.Uadmin.OTPPeriod, &user)
	user.OTPSeed = secretString
	otpPassword := services.GetOTP(user.OTPSeed, config.CurrentConfig.D.Uadmin.OTPDigits, config.CurrentConfig.D.Uadmin.OTPAlgorithm, config.CurrentConfig.D.Uadmin.OTPSkew, config.CurrentConfig.D.Uadmin.OTPPeriod)
	user.GeneratedOTPToVerify = otpPassword
	var jsonStrForSignup = []byte(fmt.Sprintf(`{"username":"test", "password": "123456", "otp": "%s"}`, otpPassword))
	db.Save(&user)
	req, _ = http.NewRequest("POST", "/auth/direct/signin/", bytes.NewBuffer(jsonStrForSignup))
	req.Header.Set("Content-Type", "application/json")
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		sessionKey := strings.Split(strings.Split(w.Header().Get("Set-Cookie"), ";")[0], "=")[1]
		req1, _ := http.NewRequest("GET", "/auth/direct/status/", nil)
		req1.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.ApiCookieName, sessionKey),
		)
		uadmin.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			assert.Contains(s.T(), w.Body.String(), "\"id\":0")
			return strings.Contains(w.Body.String(), "\"id\":0")
		})
		req2, _ := http.NewRequest("POST", "/auth/direct/logout/", bytes.NewBuffer([]byte("")))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set(
			"Cookie",
			fmt.Sprintf("%s=%s", config.CurrentConfig.D.Uadmin.ApiCookieName, sessionKey),
		)
		uadmin.TestHTTPResponse(s.T(), s.App, req2, func(w *httptest.ResponseRecorder) bool {
			assert.Equal(s.T(), w.Result().StatusCode, 204)
			return w.Result().StatusCode == 204
		})
		uadmin.TestHTTPResponse(s.T(), s.App, req1, func(w *httptest.ResponseRecorder) bool {
			assert.Contains(s.T(), w.Body.String(), "no session with key")
			return strings.Contains(w.Body.String(), "no session with key")
		})
		return w.Code == 200
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAuthAdapters(t *testing.T) {
	uadmin.Run(t, new(AuthProviderTestSuite))
}
