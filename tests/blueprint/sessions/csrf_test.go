package sessions

import (
	"github.com/stretchr/testify/assert"
	"github.com/sergeyglazyrindev/uadmin"
	interfaces2 "github.com/sergeyglazyrindev/uadmin/blueprint/sessions/interfaces"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/sergeyglazyrindev/uadmin/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type CsrfTestSuite struct {
	uadmin.TestSuite
}

func (s *CsrfTestSuite) TestSuccessfulCsrfCheck() {
	session := interfaces2.NewSession()
	token := utils.GenerateCSRFToken()
	session.SetData("csrf_token", token)
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	uadminDatabase.Db.Create(session)
	req, _ := http.NewRequest("POST", "/testcsrf", nil)
	tokenmasked := utils.MaskCSRFToken(token)
	req.Header.Set("X-CSRF-TOKEN", tokenmasked)
	req.Header.Set("X-UADMIN-API", session.Key)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
	req, _ = http.NewRequest("POST", "/testcsrf", nil)
	req.Header.Set("X-CSRF-TOKEN", "dsadsada")
	req.Header.Set("X-UADMIN-API", session.Key)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		body := w.Body.String()
		assert.Equal(s.T(), body, "Incorrect length of csrf-token")
		return strings.EqualFold(body, "Incorrect length of csrf-token")
	})
}

func (s *CsrfTestSuite) TestIgnoreCsrfCheck() {
	req, _ := http.NewRequest("POST", "/ignorecsrfcheck", nil)
	uadmin.TestHTTPResponse(s.T(), s.App, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(s.T(), w.Code, 200)
		return w.Code == 200
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCsrf(t *testing.T) {
	uadmin.Run(t, new(CsrfTestSuite))
}
