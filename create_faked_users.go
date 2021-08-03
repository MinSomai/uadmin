package uadmin

import (
	"fmt"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"strconv"
)

type CreateFakedUsersCommand struct {
}

func (c CreateFakedUsersCommand) Proceed(subaction string, args []string) error {
	uadminDatabase := interfaces.NewUadminDatabase()
	for i := range interfaces.GenerateNumberSequence(1, 100) {
		userModel := &usermodels.User{
			Email: fmt.Sprintf("admin_%d@example.com", i),
			Username: "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName: "lastname_" + strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&userModel)
	}
	uadminDatabase.Close()
	return nil
}

func (c CreateFakedUsersCommand) GetHelpText() string {
	return "Create faked users"
}
