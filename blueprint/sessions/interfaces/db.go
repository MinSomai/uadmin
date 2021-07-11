package interfaces

import (
	"fmt"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/interfaces"
	"time"
)

type DbSession struct {
	session *sessionmodel.Session
}

func (s *DbSession) GetName() string {
	return "db"
}

func (s *DbSession) Set(name string, value string) {
	s.session.SetData(name, value)
}

func (s *DbSession) SetUser(user *usermodels.User) {
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

func (s *DbSession) GetUser() *usermodels.User {
	if s.session == nil {
		return nil
	}
	return &s.session.User
}

func (s *DbSession) GetByKey(sessionKey string) (ISessionProvider, error) {
	db := interfaces.GetDB()
	var session sessionmodel.Session
	db.Model(&sessionmodel.Session{}).Where(&sessionmodel.Session{Key: sessionKey}).Preload("User").First(&session)
	if session.ID == 0 {
		return nil, fmt.Errorf("no session with key %s found", sessionKey)
	}
	return &DbSession{
		session: &session,
	}, nil
}

func (s *DbSession) Create() ISessionProvider {
	session := sessionmodel.NewSession()
	db := interfaces.GetDB()
	db.Create(session)
	return &DbSession{
		session: session,
	}
}

func (s *DbSession) Delete() bool {
	db := interfaces.GetDB()
	db.Unscoped().Delete(&sessionmodel.Session{}, s.session.ID)
	return db.Error == nil
}

func (s *DbSession) IsExpired() bool {
	return s.session.ExpiresOn.After(time.Now())
}

func (s *DbSession) Save() bool {
	db := interfaces.GetDB()
	res := db.Save(s.session)
	return res.Error == nil
}
