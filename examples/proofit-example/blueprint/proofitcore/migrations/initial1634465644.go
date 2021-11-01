package migrations

import (
	"github.com/sergeyglazyrindev/proofit-example/blueprint/proofitcore/models"
	"github.com/sergeyglazyrindev/uadmin/core"
)

type initial1634465644 struct {
}

func (m initial1634465644) GetName() string {
	return "proofit_core.1634465644"
}

func (m initial1634465644) GetID() int64 {
	return 1634465644
}

func (m initial1634465644) Up(uadminDatabase *core.UadminDatabase) error {
	err := uadminDatabase.Db.AutoMigrate(&models.ScienceCategory{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.Expert{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.Discussion{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.ScienceTerm{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.DiscussionReview{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.DiscussionComment{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.DiscussionCommentReview{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.ProofItApp{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.ExpertScienceCategory{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.ScienceCategoryLocalized{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.DiscussionLocalized{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.ScienceTermLocalized{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.DiscussionCommentLocalized{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1634465644) Down(uadminDatabase *core.UadminDatabase) error {
	err := uadminDatabase.Db.AutoMigrate(&models.ScienceCategoryLocalized{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.ScienceCategory{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.DiscussionCommentLocalized{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.DiscussionCommentReview{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.DiscussionComment{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.DiscussionLocalized{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.DiscussionReview{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.Discussion{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.ExpertScienceCategory{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.AutoMigrate(&models.ScienceTermLocalized{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.ScienceTerm{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.Expert{})
	if err != nil {
		return err
	}
	err = uadminDatabase.Db.Migrator().DropTable(&models.ProofItApp{})
	if err != nil {
		return err
	}
	return nil
}

func (m initial1634465644) Deps() []string {
	return []string{"user.1621680132"}
}
