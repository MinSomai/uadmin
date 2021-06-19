package sessions

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type CsrfTestSuite struct {
	uadmin.UadminTestSuite
}

func (s *CsrfTestSuite) TestSuccessfulCsrfCheck() {
	session := sessionmodel.NewSession()
	token := utils.GenerateCSRFToken()
	session.SetData("csrf_token", token)
	s.Db.Create(session)
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
		assert.Equal(s.T(), body, "Incorrect csrf-token")
		return strings.EqualFold(body, "Incorrect csrf-token")
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