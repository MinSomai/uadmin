package testsold

import (
	"github.com/uadmin/uadmin/config"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func setupFunction() {
	Register(
		TestStruct1{},
		TestModelA{},
		TestModelB{},
		TestApproval{},
	)

	schema := Schema["testmodelb"]
	schema.ListTheme = "default"
	schema.FormTheme = "default"
	Schema["testmodelb"] = schema

	EmailFrom = "uadmin@example.com"
	EmailPassword = "password"
	EmailUsername = "uadmin@example.com"
	EmailSMTPServer = "localhost"
	EmailSMTPServerPort = 2525

	model.RegisterInlines(TestModelA{}, map[string]string{"TestModelB": "OtherModelID"})

	ErrorHandleFunc = func(level int, err string, stack string) {
		if level >= utils.ERROR {
			utils.Trail(utils.DEBUG, stack)
		}
	}

	go StartServer(config.NewConfig("configs/test.yaml"))
	//time.Sleep(time.Second * 10)
	for !dbOK {
		time.Sleep(time.Millisecond * 100)
	}
	RateLimit = 1000000
	RateLimitBurst = 1000000
	go startEmailServer()

	PasswordAttempts = 1000000
	AllowedHosts += ",example.com"
}

func teardownFunction() {
	// Remove Generated Files
	os.Remove("uadmin.db")
	os.Remove(".key")
	os.Remove(".salt")
	os.Remove(".uproj")
	os.Remove(".bindip")

	// Delete temp media file
	os.RemoveAll("./media")
	os.RemoveAll("./static/i18n")
}

func TestMain(t *testing.M) {
	teardownFunction()
	setupFunction()
	//te := testing.T{}
	//TestSendEmail(&te)
	//time.Sleep(time.Second * 20)
	retCode := t.Run()
	teardownFunction()
	os.Exit(retCode)
}

// TestMainHandler is a unit testing function for mainHandler() function
func TestMainHandler(t *testing.T) {
	allowed := AllowedIPs
	blocked := BlockedIPs

	s1 := &Session{
		UserID: 1,
		Active: true,
	}
	s1.GenerateKey()
	s1.Save()

	u1 := &User{
		Username: "u1",
		Password: "u1",
		Active:   true,
	}
	u1.Save()

	s2 := &Session{
		UserID: u1.ID,
		Active: true,
	}
	s2.GenerateKey()
	s2.Save()

	examples := []struct {
		r       *http.Request
		ip      string
		allowed string
		blocked string
		session *Session
		code    int
		title   string
		errMsg  string
	}{
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/", nil), "", "", "", nil, 200, "uAdmin - Login", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/", nil), "10.0.0.1", "10.0.0.0/24", "", nil, 200, "uAdmin - Login", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/", nil), "10.0.0.1", "10.0.1.0/24", "", nil, 403, "uAdmin - 403", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/resetpassword", nil), "", "", "", nil, 404, "uAdmin - 404", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/", nil), "1.1.1.1", "", "", s2, 404, "uAdmin - 404", "Remote Access Denied"},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/", nil), "10.0.0.1", "", "", s2, 200, "uAdmin - Dashboard", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/export/?m=user", nil), "", "", "", s1, 303, "", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/cropper", nil), "", "", "", s1, 200, "", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/profile", nil), "10.0.0.1", "", "", s2, 200, "uAdmin - u1's Profile", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/settings", nil), "10.0.0.1", "", "", s1, 200, "uAdmin - Settings", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/user", nil), "10.0.0.1", "", "", s1, 200, "uAdmin - User", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/users/1", nil), "10.0.0.1", "", "", s1, 200, "uAdmin - User", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/users/1/1", nil), "10.0.0.1", "", "", s1, 404, "uAdmin - 404", ""},
		{httptest.NewRequest("GET", "http://0.0.0.0:5000/logout", nil), "10.0.0.1", "", "", s1, 303, "", ""},
	}

	for i, e := range examples {
		w := httptest.NewRecorder()
		if e.session != nil {
			e.r.AddCookie(&http.Cookie{Name: "session", Value: e.session.Key})
		}
		if e.ip != "" {
			e.r.RemoteAddr = e.ip + ":1234"
		}
		if e.allowed == "" {
			AllowedIPs = allowed
			BlockedIPs = blocked
		} else {
			AllowedIPs = e.allowed
			BlockedIPs = e.blocked
		}
		mainHandler(w, e.r)

		if w.Code != e.code {
			t.Errorf("mainHandler returned invalid code on example %d. Requesting %s. got %d, expected %d", i, e.r.URL.Path, w.Code, e.code)
		}

		doc, err := parseHTML(w.Result().Body, t)
		if err != nil {
			continue
		}

		if e.title != "" {
			_, content, _ := tagSearch(doc, "title", "", 0)
			if len(content) == 0 || content[0] != e.title {
				t.Errorf("mainHandler returned invalid title on example %d. Requesting %s. got %s, expected %s", i, e.r.URL.Path, content, e.title)
			}
		}
		if e.errMsg != "" {
			_, content, _ := tagSearch(doc, "h3", "", 0)
			if len(content) == 0 || content[0] != e.errMsg {
				t.Errorf("mainHandler returned invalid error message on example %d. Requesting %s. got %s, expected %s", i, e.r.URL.Path, content, e.errMsg)
			}
		}
	}

	// Test rate limit
	RateLimit = 1
	RateLimitBurst = 1
	rateLimitMap = map[string]int64{}

	for i := 0; i < 3; i++ {
		r := httptest.NewRequest("GET", "http://0.0.0.0:5000/", nil)
		w := httptest.NewRecorder()
		mainHandler(w, r)

		if i == 2 {
			if w.Body.String() != "Slow down. You are going too fast!" {
				t.Errorf("mainHandler is not rate limiting")
			}
		}
	}

	// Clean up
	AllowedIPs = allowed
	BlockedIPs = blocked

	RateLimit = 1000000
	RateLimitBurst = 1000000
	rateLimitMap = map[string]int64{}

	Delete(s1)
	Delete(s2)
	Delete(u1)
}
