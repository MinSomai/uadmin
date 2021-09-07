package core

import (
	"encoding/json"
	"fmt"
	"github.com/sergeyglazyrindev/uadmin/preloaded"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	Model
	Key        string
	User       User
	UserID     uint
	LoginTime  time.Time
	LastLogin  time.Time
	Active     bool `gorm:"default:false"`
	IP         string
	PendingOTP bool
	ExpiresOn  *time.Time
	Data       string `json:"data"`
	_data      map[string]string
}

// String return string
func (s *Session) String() string {
	return fmt.Sprintf("Session for user %s", s.User.String())
}

// Save !
func (s *Session) Save() {
	u := s.User
	s.User = User{}
	// database.Save(s)
	s.User = u
	if preloaded.CacheSessions {
		if s.Active {
			// database.Preload(s)
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
		dialect1 := GetAdapterForDb("default")
		dialect1.Equals("key", s.Key)
		// database.Get(&session, dialect1.ToString(), s.Key)
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

func (s *Session) GetData(name string) (string, error) {
	if s._data == nil {
		s.ClearAll()
	}
	val, ok := s._data[name]
	if ok {
		return val, nil
	}
	return "", fmt.Errorf("no key with name %s in this session", name)
}

func (s *Session) ClearAll() bool {
	s._data = make(map[string]string)
	return true
}

func (s *Session) BeforeSave(tx *gorm.DB) error {
	if s._data == nil {
		s.ClearAll()
	}
	var byteData []byte
	byteData, err := json.Marshal(s._data)
	if err != nil {
		return err
	}
	s.Data = string(byteData)
	return nil
}

func (s *Session) AfterFind(tx *gorm.DB) (err error) {
	if s._data == nil {
		s.ClearAll()
	}
	s._data = make(map[string]string)
	if err := json.Unmarshal([]byte(s.Data), &s._data); err != nil {
		return err
	}
	return nil
}

func (s *Session) SetData(name string, value string) error {
	if s._data == nil {
		s.ClearAll()
	}
	s._data[name] = value
	return nil
}

// HideInDashboard to return false and auto hide this from dashboard
func (Session) HideInDashboard() bool {
	return true
}

func LoadSessions() {
	// sList := []Session{}
	//database.Filter(&sList, "active = ? AND (expires_on IS NULL OR expires_on > ?)", true, time.Now())
	//// @todo, redo
	//// services.CachedSessions = map[string]Session{}
	//for _, s := range sList {
	//	database.Preload(&s)
	//	database.Preload(&s.User)
	//	// @todo, redo
	//	// services.CachedSessions[s.Key] = s
	//}
}