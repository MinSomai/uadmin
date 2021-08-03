package admin

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/admin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/form"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/templatecontext"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AdminModelActionTestSuite struct {
	uadmin.UadminTestSuite
}

func (suite *AdminModelActionTestSuite) TestAdminModelAction() {
	userModel := &usermodels.User{Username: "admin", FirstName: "firstname", LastName: "lastname", IsSuperUser: true}
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	uadminDatabase.Db.Create(userModel)
	adminUserBlueprintPage, _ := admin.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	adminUserPage, _ := adminUserBlueprintPage.SubPages.GetBySlug("user")
	adminModelAction := admin.NewAdminModelAction(
		"TurnSuperusersToNormalUsers", &admin.AdminActionPlacement{},
	)
	adminModelAction.Handler = func (afo *admin.AdminFilterObjects) (bool, int64) {
		tx := afo.GormQuerySet.Update("IsSuperUser", false).Commit()
		return tx.Error == nil, tx.RowsAffected
	}
	adminUserPage.ModelActionsRegistry.AddModelAction(adminModelAction)
	var jsonStr = []byte(fmt.Sprintf(`{"object_ids": [%d]}`, userModel.ID))
	req, _ := http.NewRequest("POST", "/admin/users/user/turnsuperuserstonormalusers/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	for adminModelAction := range adminUserPage.ModelActionsRegistry.GetAllModelActions() {
		suite.App.Router.POST(fmt.Sprintf("%s/%s/%s/%s/", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, "users", adminUserPage.ModelName, adminModelAction.SlugifiedActionName), func(adminPage *admin.AdminPage, slugifiedModelActionName string) func (ctx *gin.Context) {
			return func(ctx *gin.Context) {
				adminPage.HandleModelAction(slugifiedModelActionName, ctx)
			}
		}(adminUserPage, adminModelAction.SlugifiedActionName))

	}
	adminContext := &templatecontext.AdminContext{}
	userForm := form.NewFormFromModelFromGinContext(adminContext, &usermodels.User{}, make([]string, 0), []string{"ID"}, true, "")
	adminUserPage.Form = userForm
	uadmin.TestHTTPResponse(suite.T(), suite.App, req, func(w *httptest.ResponseRecorder) bool {
		uadminDatabase := interfaces.NewUadminDatabase()
		defer uadminDatabase.Close()
		db := uadminDatabase.Db
		var user usermodels.User
		db.Model(&usermodels.User{}).First(&user)
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
