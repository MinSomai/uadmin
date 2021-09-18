package uadmin

import (
	"fmt"
	abtestmodel "github.com/sergeyglazyrindev/uadmin/blueprint/abtest/models"
	"github.com/sergeyglazyrindev/uadmin/blueprint/approval/models"
	logmodel "github.com/sergeyglazyrindev/uadmin/blueprint/logging/models"
	"github.com/sergeyglazyrindev/uadmin/core"
	"strconv"
	"time"
)

type CreateFakedDataCommand struct {
}

func (c CreateFakedDataCommand) Proceed(subaction string, args []string) error {
	uadminDatabase := core.NewUadminDatabase()
	for i := range core.GenerateNumberSequence(1, 100) {
		userModel := &core.User{
			Email:     fmt.Sprintf("admin_%d@example.com", i),
			Username:  "admin_" + strconv.Itoa(i),
			FirstName: "firstname_" + strconv.Itoa(i),
			LastName:  "lastname_" + strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&userModel)
		oneTimeAction := &core.OneTimeAction{
			User:      *userModel,
			ExpiresOn: time.Now(),
			Code:      strconv.Itoa(i),
		}
		uadminDatabase.Db.Create(&oneTimeAction)
		session := &core.Session{
			User:      userModel,
			LoginTime: time.Now(),
			LastLogin: time.Now(),
		}
		uadminDatabase.Db.Create(&session)
	}
	var contentTypes []*core.ContentType
	uadminDatabase.Db.Find(&contentTypes)
	for _, contentType := range contentTypes {
		logModel := logmodel.Log{
			Username:      "admin",
			ContentTypeID: contentType.ID,
		}
		uadminDatabase.Db.Create(&logModel)
		approvalModel := models.Approval{
			ContentTypeID:       contentType.ID,
			ModelPK:             uint(1),
			ColumnName:          "Email",
			OldValue:            "admin@example.com",
			NewValue:            "admin1@example.com",
			NewValueDescription: "changing email",
			ChangedBy:           "superuser",
			ChangeDate:          time.Now(),
		}
		uadminDatabase.Db.Create(&approvalModel)
		abTestModel := abtestmodel.ABTest{
			Name:          "test_1",
			ContentTypeID: contentType.ID,
		}
		uadminDatabase.Db.Create(&abTestModel)
		for i := range core.GenerateNumberSequence(0, 100) {
			abTestValueModel := abtestmodel.ABTestValue{
				ABTest: abTestModel,
				Value:  strconv.Itoa(i),
			}
			uadminDatabase.Db.Create(&abTestValueModel)
		}
	}
	for i := range core.GenerateNumberSequence(1, 20) {
		groupModel := core.UserGroup{
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
