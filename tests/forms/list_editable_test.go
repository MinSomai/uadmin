package forms

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/admin"
	"github.com/uadmin/uadmin/interfaces"
	"mime/multipart"
	"testing"
)

func NewTestForm1() *multipart.Form {
	form1 := multipart.Form{
		Value: make(map[string][]string),
	}
	return &form1
}

type ListEditableFormTestSuite struct {
	uadmin.UadminTestSuite
}

func (s *ListEditableFormTestSuite) TestFormBuilder() {
	// userBlueprintRegistry, _ := s.App.BlueprintRegistry.GetByName("user")
	// NewFormListEditableFromListDisplayRegistry
	adminPanel, _ := admin.CurrentDashboardAdminPanel.AdminPages.GetBySlug("users")
	userAdminPage, _ := adminPanel.SubPages.GetBySlug("user")
	ld, _ := userAdminPage.ListDisplay.GetFieldByDisplayName("Email")
	ld.IsEditable = true
	listEditableForm := interfaces.NewFormListEditableFromListDisplayRegistry(10, &interfaces.User{}, userAdminPage.ListDisplay)
	form := NewTestForm1()
	userTest := &interfaces.User{}
	err := listEditableForm.ProceedRequest(form, userTest)
	assert.False(s.T(), err.IsEmpty())
	form.Value["10_Email"] = []string{"admin@example.com"}
	err = listEditableForm.ProceedRequest(form, userTest)
	assert.True(s.T(), err.IsEmpty())
	assert.Equal(s.T(), userTest.Email, "admin@example.com")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestListEditableForm(t *testing.T) {
	uadmin.Run(t, new(ListEditableFormTestSuite))
}