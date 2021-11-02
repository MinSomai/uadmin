package proofit_example

import (
	nestedset "github.com/griffinqiu/go-nested-set"
	models2 "github.com/sergeyglazyrindev/proofit-example/blueprint/proofitcore/models"
	utils2 "github.com/sergeyglazyrindev/uadmin/blueprint/auth/utils"
	"github.com/sergeyglazyrindev/uadmin/core"
	"strconv"
)

type CreateFakedDataCommand struct {
}

func (c CreateFakedDataCommand) Proceed(subaction string, args []string) error {
	uadminDatabase := core.NewUadminDatabase()
	for i := range core.GenerateNumberSequence(1, 10) {
		uadminDatabase.Db.Create(&models2.ScienceCategory{Name: "category_" + strconv.Itoa(i)})
	}
	for i := range core.GenerateNumberSequence(1, 100) {
		salt := core.GenerateRandomString(currentApp.Config.D.Auth.SaltLength)
		// hashedPassword, err := utils2.HashPass(password, salt)
		hashedPassword, _ := utils2.HashPass("password_"+strconv.Itoa(i), salt)
		user := core.GenerateUserModel()
		user.SetFirstName("First name " + strconv.Itoa(i))
		user.SetLastName("Last name " + strconv.Itoa(i))
		user.SetUsername("username-" + strconv.Itoa(i))
		user.SetEmail("username-" + strconv.Itoa(i) + "@proofit.com")
		user.SetPassword(hashedPassword)
		user.SetActive(true)
		user.SetIsStaff(true)
		user.SetSalt(salt)
		user.SetIsPasswordUsable(true)
		uadminDatabase.Db.Create(user)
		uadminDatabase.Db.Create(&models2.Expert{UserID: user.GetID()})
	}
	for i := range core.GenerateNumberSequence(1, 10) {
		uadminDatabase.Db.Create(&models2.ScienceTerm{Alias: "axiom_" + strconv.Itoa(i), Type: 1, ScienceCategoryID: uint(i)})
	}
	for i := range core.GenerateNumberSequence(1, 10) {
		uadminDatabase.Db.Create(&models2.ScienceTerm{Alias: "law_" + strconv.Itoa(i), Type: 2, ScienceCategoryID: uint(i)})
	}
	for i := range core.GenerateNumberSequence(1, 10) {
		uadminDatabase.Db.Create(&models2.ScienceTerm{Alias: "deduction_" + strconv.Itoa(i), Type: 3, ScienceCategoryID: uint(i)})
	}
	for i := range core.GenerateNumberSequence(1, 10) {
		uadminDatabase.Db.Create(&models2.ScienceTerm{Alias: "induction_" + strconv.Itoa(i), Type: 4, ScienceCategoryID: uint(i)})
	}
	for i := range core.GenerateNumberSequence(1, 100) {
		discussion := &models2.Discussion{AuthorID: uint(i)}
		uadminDatabase.Db.Create(discussion)
		discussionComment := &models2.DiscussionComment{DiscussionID: discussion.ID, AuthorID: uint(i%10) + 1}
		nestedset.Create(uadminDatabase.Db, discussionComment, nil)
		discussionComment1 := &models2.DiscussionComment{DiscussionID: discussion.ID, AuthorID: uint(((i + 1) % 10) + 1)}
		nestedset.Create(uadminDatabase.Db, discussionComment1, discussionComment)
	}
	uadminDatabase.Close()
	return nil
}

func (c CreateFakedDataCommand) GetHelpText() string {
	return "Create fake data for testing proofit"
}
