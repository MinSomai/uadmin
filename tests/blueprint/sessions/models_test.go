package sessions

import (
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"testing"
	"time"
)

type SessionTestSuite struct {
	uadmin.UadminTestSuite
}

func (s *SessionTestSuite) TestSavingSession() {
	session := sessionmodel.NewSession()
	session.SetData("testkey", "testvalue")
	s.Db.Create(session)
	var loadedsession sessionmodel.Session
	s.Db.Model(&sessionmodel.Session{}).First(&loadedsession)
	val, _ := loadedsession.GetData("testkey")
	assert.Equal(s.T(), val, "testvalue")
}

func (s *SessionTestSuite) TestTransactionConsistencyInTests() {
	var loadedsession sessionmodel.Session
	s.Db.Model(&sessionmodel.Session{}).First(&loadedsession)
	val, _ := loadedsession.GetData("testkey")
	assert.Equal(s.T(), val, "")
}

func (s *SessionTestSuite) TestDbSessionAdapter() {
	blueprint, _ := s.App.BlueprintRegistry.GetByName("sessions")
	dbadapter, _ := blueprint.(sessionsblueprint.Blueprint).SessionAdapterRegistry.GetAdapter("db")
	dbadapter = dbadapter.Create()
	assert.Equal(s.T(), dbadapter.GetUser().ID, uint(0))
	dbadapter.Set("testkey", "testvalue")
	dbadapter.Save()
	dbadapter, _ = dbadapter.GetByKey(dbadapter.GetKey())
	val, _ := dbadapter.Get("testkey")
	assert.Equal(s.T(), val, "testvalue")
	dbadapter.ClearAll()
	dbadapter.Save()
	dbadapter, _ = dbadapter.GetByKey(dbadapter.GetKey())
	_, err := dbadapter.Get("testkey")
	if err == nil {
		assert.True(s.T(), false)
	}
	expiresOn := time.Time{}
	expiresOn = expiresOn.Add(10 * time.Minute)
	dbadapter.ExpiresOn(&expiresOn)
	assert.False(s.T(), dbadapter.IsExpired())
	sessionKey := dbadapter.GetKey()
	isRemoved := dbadapter.Delete()
	if !isRemoved {
		assert.True(s.T(), false)
	}
	dbadapter, err = dbadapter.GetByKey(sessionKey)
	if err == nil {
		assert.True(s.T(), false)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSessions(t *testing.T) {
	uadmin.Run(t, new(SessionTestSuite))
}
