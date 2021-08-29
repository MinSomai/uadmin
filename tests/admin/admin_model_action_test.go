package admin

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/interfaces"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AdminModelActionTestSuite struct {
	uadmin.UadminTestSuite
}

func (suite *AdminModelActionTestSuite) TestAdminModelAction() {
	userModel := &interfaces.User{Username: "admin", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	uadminDatabase.Db.Create(userModel)
	adminUserBlueprintPage, _ := interfaces.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	adminModelAction := interfaces.NewAdminModelAction(
		"TurnSuperusersToNormalUsers", &interfaces.AdminActionPlacement{},
	)
	adminModelAction.Handler = func (ap *interfaces.AdminPage, afo interfaces.IAdminFilterObjects, ctx *gin.Context) (bool, int64) {
		tx := afo.GetFullQuerySet().Update("IsSuperUser", false).Commit()
		return tx.(*interfaces.GormPersistenceStorage).Db.Error == nil, tx.(*interfaces.GormPersistenceStorage).Db.RowsAffected
	}
	adminUserPage.ModelActionsRegistry.AddModelAction(adminModelAction)
	var jsonStr = []byte(fmt.Sprintf(`{"object_ids": "%d"}`, userModel.ID))
	req, _ := http.NewRequest("POST", "/admin/users/user/turnsuperuserstonormalusers/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	for adminModelAction := range adminUserPage.ModelActionsRegistry.GetAllModelActions() {
		suite.App.Router.Any(fmt.Sprintf("%s/%s/%s/%s/", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, "users", adminUserPage.ModelName, adminModelAction.SlugifiedActionName), func(adminPage *interfaces.AdminPage, slugifiedModelActionName string) func (ctx *gin.Context) {
			return func(ctx *gin.Context) {
				adminPage.HandleModelAction(slugifiedModelActionName, ctx)
			}
		}(adminUserPage, adminModelAction.SlugifiedActionName))

	}
	adminContext := &interfaces.AdminContext{}
	userForm := interfaces.NewFormFromModelFromGinContext(adminContext, &interfaces.User{}, make([]string, 0), []string{"ID"}, true, "")
	adminUserPage.Form = userForm
	uadmin.TestHTTPResponse(suite.T(), suite.App, req, func(w *httptest.ResponseRecorder) bool {
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var user interfaces.User
		db.Model(&interfaces.User{}).First(&user)
		assert.False(suite.T(), user.IsSuperUser)
		return user.IsSuperUser == false
	})
	// adminUserPage.HandleModelAction("TurnSuperusersToNormalUsers", &gin.Context{})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAdminModelAction(t *testing.T) {
	uadmin.Run(t, new(AdminModelActionTestSuite))
}
