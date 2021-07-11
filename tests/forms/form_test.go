package forms

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/form"
	"github.com/uadmin/uadmin/interfaces"
	"testing"
)

type UsernameFormOptions struct {
	form.FieldFormOptions
}

type FormTestSuite struct {
	uadmin.UadminTestSuite
}

func (s *FormTestSuite) TestFormBuilder() {
	fieldChoiceRegistry := interfaces.FieldChoiceRegistry{}
	fieldChoiceRegistry.Choices = make([]*interfaces.FieldChoice, 0)
	formOptions := &UsernameFormOptions{
		FieldFormOptions: form.FieldFormOptions{
			Name: "UsernameOptions",
			Initial: "InitialUsername",
			DisplayName: "Display name",
			Validators: make([]interfaces.IValidator, 0),
			Choices: &fieldChoiceRegistry,
			HelpText: "help for username",
		},
	}
	interfaces.CurrentConfig.AddFieldFormOptions(formOptions)
	// initial=\"test\",displayname=\"uname\",validators=\"password-uadmin\",choices=UsernameChoices,helptext=\"HELPPPPPPPPPP\"
	user := &usermodels.User{}
	form1 := form.NewFormFromModel(user, make([]string, 0), []string{"Username", "FirstName", "LastName", "Email", "Photo", "LastLogin", "ExpiresOn", "OTPRequired"}, true, "")
	res := form1.Render()
	assert.Contains(s.T(), res, "<form")
	form2 := NewTestForm()
	form2.Value["Username"] = []string{"username"}
	form2.Value["FirstName"] = []string{"first name"}
	form2.Value["LastName"] = []string{"last name"}
	form2.Value["Email"] = []string{"email@example.com"}
	form2.Value["OTPRequired"] = []string{"yes"}
	formError := form1.ProceedRequest(form2, user)
	assert.Equal(s.T(), user.Username, "username")
	assert.True(s.T(), formError.IsEmpty())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestForm(t *testing.T) {
	uadmin.Run(t, new(FormTestSuite))
}


