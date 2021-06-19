package interfaces

import (
	"fmt"
	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"time"
)

type ISessionProvider interface {
	GetKey() string
	Create() ISessionProvider
	GetByKey(key string) (ISessionProvider, error)
	GetName() string
	IsExpired() bool
	Delete() bool
	Set(name string, value string)
	Get(name string) (string, error)
	ClearAll() bool
	GetUser() *usermodels.User
	SetUser(user *usermodels.User)
	Save() bool
	ExpiresOn(*time.Time)
}

type SessionProviderRegistry struct {
	registeredSessionAdapters map[string]ISessionProvider
	defaultAdapter string
}

func (r *SessionProviderRegistry) RegisterNewAdapter(adapter ISessionProvider, defaultAdapter bool) {
	r.registeredSessionAdapters[adapter.GetName()] = adapter
	if defaultAdapter {
		r.defaultAdapter = adapter.GetName()
	}
}

func (r *SessionProviderRegistry) GetAdapter(name string) (ISessionProvider, error) {
	adapter, ok := r.registeredSessionAdapters[name]
	if ok {
		return adapter, nil
	} else {
		return nil, fmt.Errorf("adapter with name %s not found", name)
	}
}

func (r *SessionProviderRegistry) GetDefaultAdapter() (ISessionProvider, error) {
	adapter, ok := r.registeredSessionAdapters[r.defaultAdapter]
	if ok {
		return adapter, nil
	} else {
		return nil, fmt.Errorf("no default session adapter configured")
	}
}

func NewSessionRegistry() *SessionProviderRegistry {
	return &SessionProviderRegistry{
		registeredSessionAdapters: make(map[string]ISessionProvider),
		defaultAdapter: "",
	}
}