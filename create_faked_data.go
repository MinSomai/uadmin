package uadmin

import (
	"fmt"
	abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
	"github.com/uadmin/uadmin/blueprint/approval/models"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	"github.com/uadmin/uadmin/interfaces"
	"strconv"
	"time"
)

type CreateFakedDataCommand struct {
}

func (c CreateFakedDataCommand) Proceed(subaction string, args []string) error {
	uadminDatabase := interfaces.NewUadminDatabase()
	for i := range interfaces.GenerateNumberSequence(1, 100) {
		userModel := &interfaces.User{
			Email: fmt.Sprintf("admin_%d@example.com", i),
			Username: "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName: "lastname_" + strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&userModel)
		oneTimeAction := &interfaces.OneTimeAction{
			User: *userModel,
			ExpiresOn: time.Now(),
			Code: strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&oneTimeAction)
		session := &interfaces.Session{
			User: *userModel,
			LoginTime: time.Now(),
			LastLogin: time.Now(),
		}
		uadminDatabase.Db.Create(&session)
	}
	for i := range interfaces.GenerateNumberSequence(1, 100) {
		abTestModel := abtestmodel.ABTest{
			Name: fmt.Sprintf("test_%d", i),
		}
		uadminDatabase.Db.Create(&abTestModel)
		abTestValueModel := abtestmodel.ABTestValue{
			ABTest: abTestModel,
		}
		uadminDatabase.Db.Create(&abTestValueModel)
	}
	for i := range interfaces.GenerateNumberSequence(1, 100) {
		approvalModel := models.Approval{
			ModelName: "user",
			ModelPK: uint(i),
			ColumnName: "Email",
			OldValue: "admin@example.com",
			NewValue: "admin1@example.com",
			NewValueDescription: "changing email",
			ChangedBy: "superuser",
			ChangeDate: time.Now(),
		}
		uadminDatabase.Db.Create(&approvalModel)
	}
	for i := range interfaces.GenerateNumberSequence(1, 100) {
		logModel := logmodel.Log{
			Username: "admin",
			TableName: "user",
			TableID: i,
		}
		uadminDatabase.Db.Create(&logModel)
	}
	for i := range interfaces.GenerateNumberSequence(1, 20) {
		groupModel := interfaces.UserGroup{
			GroupName: fmt.Sprintf("Group name %d", i),
		}
		uadminDatabase.Db.Create(&groupModel)
	}
	uadminDatabase.Close()
	return nil
}

func (c CreateFakedDataCommand) GetHelpText() string {
	return "Create fake data for testing uadmin"
}
