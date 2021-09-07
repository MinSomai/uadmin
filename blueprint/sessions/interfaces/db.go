package interfaces

import (
	"fmt"
	"github.com/sergeyglazyrindev/uadmin/core"
	"time"
)

type DbSession struct {
	session *core.Session
}

func (s *DbSession) GetName() string {
	return "db"
}

func (s *DbSession) Set(name string, value string) {
	s.session.SetData(name, value)
}

func (s *DbSession) SetUser(user *core.User) {
	s.session.UserID = user.ID
	s.session.User = *user
}

func (s *DbSession) Get(name string) (string, error) {
	return s.session.GetData(name)
}

func (s *DbSession) ClearAll() bool {
	return s.session.ClearAll()
}

func (s *DbSession) ExpiresOn(expiresOn *time.Time) {
	s.session.ExpiresOn = expiresOn
}

func (s *DbSession) GetKey() string {
	return s.session.Key
}

func (s *DbSession) GetUser() *core.User {
	if s.session == nil {
		return nil
	}
	return &s.session.User
}

func (s *DbSession) GetByKey(sessionKey string) (ISessionProvider, error) {
	db := core.NewUadminDatabase()
	defer db.Close()
	var session core.Session
	db.Db.Model(&core.Session{}).Where(&core.Session{Key: sessionKey}).Preload("User").First(&session)
	if session.ID == 0 {
		return nil, fmt.Errorf("no session with key %s found", sessionKey)
	}
	return &DbSession{
		session: &session,
	}, nil
}

func (s *DbSession) Create() ISessionProvider {
	session := NewSession()
	db := core.NewUadminDatabase()
	defer db.Close()
	db.Db.Create(session)
	return &DbSession{
		session: session,
	}
}

func (s *DbSession) Delete() bool {
	db := core.NewUadminDatabase()
	defer db.Close()
	db.Db.Unscoped().Delete(&core.Session{}, s.session.ID)
	return db.Db.Error == nil
}

func (s *DbSession) IsExpired() bool {
	return s.session.ExpiresOn.After(time.Now())
}

func (s *DbSession) Save() bool {
	db := core.NewUadminDatabase()
	defer db.Close()
	res := db.Db.Save(s.session)
	return res.Error == nil
}
