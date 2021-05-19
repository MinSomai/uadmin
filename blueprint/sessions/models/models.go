package models

import (
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/dialect"
	"time"
)

// Session !
type Session struct {
	model.Model
	Key        string
	User       models.User `uadmin:"filter"`
	UserID     uint
	LoginTime  time.Time
	LastLogin  time.Time
	Active     bool   `uadmin:"filter"`
	IP         string `uadmin:"filter"`
	PendingOTP bool   `uadmin:"filter"`
	ExpiresOn  *time.Time
}

// String return string
func (s Session) String() string {
	return s.Key
}

// Save !
func (s *Session) Save() {
	u := s.User
	s.User = models.User{}
	database.Save(s)
	s.User = u
	if preloaded.CacheSessions {
		if s.Active {
			database.Preload(s)
			// @todo, redo
			// services.CachedSessions[s.Key] = *s
		} else {
			// @todo, redo
			// delete(services.CachedSessions, s.Key)
		}
	}
}

// GenerateKey !
func (s *Session) GenerateKey() {
	session := Session{}
	for {
		// TODO: Increase the session length to 124 and add 4 bytes for User.ID
		// @todo, redo
		// s.Key = services.GenerateBase64(24)
		dialect1 := dialect.GetDialectForDb()
		dialect1.Equals("key", s.Key)
		database.Get(&session, dialect1.ToString(), s.Key)
		if session.ID == 0 {
			break
		}
	}
}

// Logout deactivates a session
func (s *Session) Logout() {
	s.Active = false
	s.Save()
}

// HideInDashboard to return false and auto hide this from dashboard
func (Session) HideInDashboard() bool {
	return true
}

func LoadSessions() {
	sList := []Session{}
	database.Filter(&sList, "active = ? AND (expires_on IS NULL OR expires_on > ?)", true, time.Now())
	// @todo, redo
	// services.CachedSessions = map[string]Session{}
	for _, s := range sList {
		database.Preload(&s)
		database.Preload(&s.User)
		// @todo, redo
		// services.CachedSessions[s.Key] = s
	}
}
