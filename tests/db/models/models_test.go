package models

import (
	"fmt"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type BuildRemovalTreeTestSuite struct {
	uadmin.TestSuite
	ContentType *core.ContentType
}

func (s *BuildRemovalTreeTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	uadminDatabase.Db.AutoMigrate(&UserGroupContentType{})
	uadminDatabase.Db.AutoMigrate(&UserContentType{})
	uadminDatabase.Db.AutoMigrate(&OneTimeActionContentType{})
	uadminDatabase.Db.AutoMigrate(&SessionContentType{})
	core.ProjectModels.RegisterModel(func() interface{} { return &SessionContentType{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &OneTimeActionContentType{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &UserContentType{} })
	core.ProjectModels.RegisterModel(func() interface{} { return &UserGroupContentType{} })
}

func (s *BuildRemovalTreeTestSuite) ConfigureData(uadminDatabase *core.UadminDatabase) {
	contentType := &core.ContentType{BlueprintName: "user", ModelName: "user"}
	uadminDatabase.Db.Create(contentType)
	s.ContentType = contentType
	permission := &core.Permission{Name: "user_read", ContentType: *contentType}
	uadminDatabase.Db.Create(permission)
	usergroup := &core.UserGroup{GroupName: "test"}
	uadminDatabase.Db.Create(usergroup)
	uadminDatabase.Db.Model(usergroup).Association("Permissions").Append(permission)
	uadminDatabase.Db.Save(usergroup)
	user := &core.User{Email: "admin@example.com"}
	uadminDatabase.Db.Create(user)
	uadminDatabase.Db.Model(user).Association("Permissions").Append(permission)
	uadminDatabase.Db.Model(user).Association("UserGroups").Append(usergroup)
	uadminDatabase.Db.Save(user)
	oneTimeAction := &core.OneTimeAction{User: *user, Code: "aaa"}
	uadminDatabase.Db.Create(oneTimeAction)
	session := &core.Session{User: *user, LoginTime: time.Now(), LastLogin: time.Now()}
	uadminDatabase.Db.Create(session)
	sessionContentType := &SessionContentType{Session: *session, ContentType: *contentType}
	uadminDatabase.Db.Save(sessionContentType)
	oneTimeActionContentType := &OneTimeActionContentType{OneTimeAction: *oneTimeAction, ContentType: *contentType}
	uadminDatabase.Db.Save(oneTimeActionContentType)
	userContentType := &UserContentType{User: *user, ContentType: *contentType}
	uadminDatabase.Db.Save(userContentType)
	userGroupContentType := &UserGroupContentType{UserGroup: *usergroup, ContentType: *contentType}
	uadminDatabase.Db.Save(userGroupContentType)
}

func (s *BuildRemovalTreeTestSuite) TearDownSuite() {
	uadminDatabase := core.NewUadminDatabase()
	uadminDatabase.Db.Migrator().DropTable(&UserGroupContentType{})
	uadminDatabase.Db.Migrator().DropTable(&UserContentType{})
	uadminDatabase.Db.Migrator().DropTable(&OneTimeActionContentType{})
	uadminDatabase.Db.Migrator().DropTable(&SessionContentType{})
	uadminDatabase.Close()
	s.TestSuite.TearDownSuite()
}

type UserGroupContentType struct {
	core.Model
	UserGroup     core.UserGroup
	UserGroupID   uint
	ContentType   core.ContentType
	ContentTypeID uint
}

func (ugct *UserGroupContentType) String() string {
	return fmt.Sprintf("dsadsa-usergroup-%d-%s", ugct.ID, ugct.ContentType.String())
}

type UserContentType struct {
	core.Model
	User          core.User
	UserID        uint
	ContentType   core.ContentType
	ContentTypeID uint
}

func (ugct *UserContentType) String() string {
	return fmt.Sprintf("dsadsa-user-%d-%s", ugct.ID, ugct.ContentType.String())
}

type OneTimeActionContentType struct {
	core.Model
	OneTimeAction   core.OneTimeAction
	OneTimeActionID uint
	ContentType     core.ContentType
	ContentTypeID   uint
}

func (ugct *OneTimeActionContentType) String() string {
	return fmt.Sprintf("dsadsa-onetimeaction-%d-%s", ugct.ID, ugct.ContentType.String())
}

type SessionContentType struct {
	core.Model
	Session       core.Session
	SessionID     uint
	ContentType   core.ContentType
	ContentTypeID uint
}

func (ugct *SessionContentType) String() string {
	return fmt.Sprintf("dsadsa-session-%d-%s", ugct.ID, ugct.ContentType.String())
}

func (s *BuildRemovalTreeTestSuite) TestRemovalStringified() {
	uadminDatabase := core.NewUadminDatabase()
	defer uadminDatabase.Close()
	s.ConfigureData(uadminDatabase)
	//spew.Dump("contentType", contentType.ID)
	//spew.Dump("permission", permission.ID)
	//spew.Dump("usergroup", usergroup.ID)
	//spew.Dump("user", user.ID)
	//spew.Dump("onetimeaction", oneTimeAction.ID)
	//spew.Dump("session", session.ID)
	//spew.Dump("usergroup permissions", len(usergroup.Permissions))
	//spew.Dump("user permissions", len(user.Permissions))
	//spew.Dump("user groups", len(user.UserGroups))
	removalTreeNode := core.BuildRemovalTree(uadminDatabase, s.ContentType)
	deletionStringified := removalTreeNode.BuildDeletionTreeStringified(uadminDatabase)
	assert.Equal(s.T(), len(deletionStringified), 15)
}

func (s *BuildRemovalTreeTestSuite) TestRemoval() {
	uadminDatabase := core.NewUadminDatabase()
	s.ConfigureData(uadminDatabase)
	defer uadminDatabase.Close()
	var c int64
	removalTreeNode := core.BuildRemovalTree(uadminDatabase, s.ContentType)
	removalTreeNode.RemoveFromDatabase(uadminDatabase)
	uadminDatabase.Db.Model(&core.Permission{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	uadminDatabase.Db.Model(&OneTimeActionContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	uadminDatabase.Db.Model(&UserContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	uadminDatabase.Db.Model(&UserGroupContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
	uadminDatabase.Db.Model(&SessionContentType{}).Count(&c)
	assert.Equal(s.T(), c, int64(0))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBuildRemovalTree(t *testing.T) {
	uadmin.Run(t, new(BuildRemovalTreeTestSuite))
}
