package models

import (
	"encoding/json"
	"fmt"
	menumodel "github.com/uadmin/uadmin/blueprint/menu/models"
	"github.com/uadmin/uadmin/database"
	// "github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"

	// "time"
)

//type User struct {
//	model.Model
//	Username     string    `uadmin:"required;filter;search" gorm:"uniqueIndex" json:"username"`
//	FirstName    string    `uadmin:"filter;search" json:"first_name"`
//	LastName     string    `uadmin:"filter;search" json:"last_name"`
//	Password     string    `uadmin:"required;password;help:To reset password, clear the field and type a new password.;list_exclude" json:"password"`
//	Email        string    `uadmin:"email;search" gorm:"uniqueIndex" json:"email"`
//	Active       bool      `uadmin:"filter" json:"active"`
//	Admin        bool      `uadmin:"filter" json:"admin"`
//	RemoteAccess bool      `uadmin:"filter" json:"remote_access"`
//	UserGroup    UserGroup `uadmin:"filter"`
//	UserGroupID  uint `json:"user_group_id"`
//	Photo        string `uadmin:"image" json:"photo"`
//	//Language     []Language `gorm:"many2many:user_languages" listExclude:"true"`
//	LastLogin   *time.Time `uadmin:"read_only" json:"read_only"`
//	ExpiresOn   *time.Time `json:"expires_on"`
//	GeneratedOTPToVerify     string `uadmin:"list_exclude;hidden;read_only" json:"generated_otp_to_verify"`
//	OTPSeed     string `uadmin:"list_exclude;hidden;read_only" json:"otp_seed"`
//	Salt        string `json:"salt"`
//}

// String return string
//func (u User) String() string {
//	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
//}

// Save !
func (u *User) Save() {
	// @todo, redo
	//if !strings.HasPrefix(u.Password, "$2a$") && len(u.Password) != 60 {
	//	u.Password = authservices.HashPass(u.Password)
	//}
	//if u.OTPSeed == "" {
	//	// @todo, redo
	//	// u.OTPSeed, _ = otpservices.GenerateOTPSeed(preloaded.OTPDigits, preloaded.OTPAlgorithm, preloaded.OTPSkew, preloaded.OTPPeriod, u)
	//} else if u.ID != 0 {
	//	oldUser := User{}
	//	database.Get(&oldUser, "id = ?", u.ID)
	//	if !oldUser.OTPRequired && u.OTPRequired {
	//		// @todo, redo
	//		// u.OTPSeed, _ = otpservices.GenerateOTPSeed(preloaded.OTPDigits, preloaded.OTPAlgorithm, preloaded.OTPSkew, preloaded.OTPPeriod, u)
	//	}
	//}
	// u.Username = strings.ToLower(u.Username)
	// database.Save(u)
}

// @todo, redo
//// GetActiveSession !
//func (u *User) GetActiveSession() *sessionmodel.Session {
//	s := sessionmodel.Session{}
//	dialect1 := dialect.GetDialectForDb()
//	database.Get(&s, dialect1.Quote("user_id")+" = ? AND "+dialect1.Quote("active")+" = ?", u.ID, true)
//	if s.ID == 0 {
//		return nil
//	}
//	return &s
//}

// @todo, redo
//// Login Logs in user using password and otp. If there is no OTP, just pass an empty string
//func (u *User) Login(pass string, otp string) *sessionmodel.Session {
//	if u == nil {
//		return nil
//	}
//
//	password := []byte(pass + authservices.Salt)
//	hashedPassword := []byte(u.Password)
//	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
//	if err == nil && u.ID != 0 {
//		s := u.GetActiveSession()
//		if s == nil {
//			s = &sessionmodel.Session{}
//			s.Active = true
//			s.UserID = u.ID
//			s.LoginTime = time.Now()
//			s.GenerateKey()
//			if authservices.CookieTimeout > -1 {
//				ExpiresOn := s.LoginTime.Add(time.Second * time.Duration(authservices.CookieTimeout))
//				s.ExpiresOn = &ExpiresOn
//			}
//		}
//		s.LastLogin = time.Now()
//		if u.OTPRequired {
//			if otp == "" {
//				s.PendingOTP = true
//			} else {
//				s.PendingOTP = !u.VerifyOTP(otp)
//			}
//		}
//		u.LastLogin = &s.LastLogin
//		u.Save()
//		s.Save()
//		return s
//	}
//	return nil
//}

// GetDashboardMenu !
func (u *User) GetDashboardMenu() (menus []menumodel.DashboardMenu) {
	allItems := []menumodel.DashboardMenu{}
	database.All(&allItems)

	userItems := []UserPermission{}
	database.Filter(&userItems, "user_id = ?", u.ID)

	groupItems := []GroupPermission{}
	database.Filter(&groupItems, "user_group_id = ?", u.UserGroupID)

	var groupItemIndex int
	var userItemIndex int
	dashboardItems := []menumodel.DashboardMenu{}
	for _, item := range allItems {
		groupItemIndex = -1
		userItemIndex = -1
		for i, groupItem := range groupItems {
			if groupItem.DashboardMenuID == item.ID {
				groupItemIndex = i
				break
			}
		}
		for i, userItem := range userItems {
			if userItem.DashboardMenuID == item.ID {
				userItemIndex = i
				break
			}
		}
		// Permission exists for group and user: overide group with user
		if groupItemIndex != -1 && userItemIndex != -1 {
			groupItems[groupItemIndex].Read = userItems[userItemIndex].Read
			groupItems[groupItemIndex].Add = userItems[userItemIndex].Add
			groupItems[groupItemIndex].Edit = userItems[userItemIndex].Edit
			groupItems[groupItemIndex].Delete = userItems[userItemIndex].Delete
		}
		// User permission exists but no group, add it to permessions
		if groupItemIndex == -1 && userItemIndex != -1 {
			groupItems = append(groupItems, GroupPermission{
				DashboardMenuID: userItems[userItemIndex].DashboardMenuID,
				Read:            userItems[userItemIndex].Read,
				Add:             userItems[userItemIndex].Add,
				Edit:            userItems[userItemIndex].Edit,
				Delete:          userItems[userItemIndex].Delete,
			})
			groupItemIndex = len(groupItems) - 1
		}
		// Reconstruct the dashboard list
		if u.Admin || groupItemIndex != -1 || userItemIndex != -1 {
			if u.Admin || groupItems[groupItemIndex].Read {
				dashboardItems = append(dashboardItems, item)
			}
		}
	}
	return dashboardItems
}

// HasAccess returns the user level permission to a model. The modelName
// the the URL of the model
func (u *User) HasAccess(modelName string) UserPermission {
	utils.Trail(utils.WARNING, "User.HasAccess will be deprecated in version 0.6.0. Use User.GetAccess instead.")
	return u.hasAccess(modelName)
}

// hasAccess returns the user level permission to a model. The modelName
// the the URL of the model
func (u *User) hasAccess(modelName string) UserPermission {
	up := UserPermission{}
	dm := menumodel.DashboardMenu{}
	if preloaded.CachePermissions {
		modelID := uint(0)
		for _, m := range cachedModels {
			if m.URL == modelName {
				modelID = m.ID
				break
			}
		}
		for _, p := range cacheUserPerms {
			if p.UserID == u.ID && p.DashboardMenuID == modelID {
				up = p
				break
			}
		}
	} else {
		database.Get(&dm, "url = ?", modelName)
		database.Get(&up, "user_id = ? and dashboard_menu_id = ?", u.ID, dm.ID)
	}
	return up
}

// GetAccess returns the user's permission to a dashboard menu based on
// their admin status, group and user permissions
func (u *User) GetAccess(modelName string) UserPermission {
	// Check if the user has permission to a model
	if u.UserGroup.ID != u.UserGroupID {
		database.Preload(u)
	}
	uPerm := u.hasAccess(modelName)
	gPerm := u.UserGroup.hasAccess(modelName)
	perm := UserPermission{}

	if gPerm.ID != 0 {
		perm.Read = gPerm.Read
		perm.Edit = gPerm.Edit
		perm.Add = gPerm.Add
		perm.Delete = gPerm.Delete
		perm.Approval = gPerm.Approval
	}
	if uPerm.ID != 0 {
		perm.Read = uPerm.Read
		perm.Edit = uPerm.Edit
		perm.Add = uPerm.Add
		perm.Delete = uPerm.Delete
		perm.Approval = uPerm.Approval
	}
	if u.Admin {
		perm.Read = true
		perm.Edit = true
		perm.Add = true
		perm.Delete = true
		perm.Approval = true
	}
	return perm
}

// Validate user when saving from uadmin
func (u User) Validate() (ret map[string]string) {
	ret = map[string]string{}
	if u.ID == 0 {
		database.Get(&u, "username=?", u.Username)
		if u.ID > 0 {
			ret["Username"] = "Username is already Taken."
		}
	}
	return
}

// GetOTP !
func (u *User) GetOTP() string {
	return ""
	// return otpservices.GetOTP(u.OTPSeed, preloaded.OTPDigits, preloaded.OTPAlgorithm, preloaded.OTPSkew, preloaded.OTPPeriod)
}

// VerifyOTP !
func (u *User) VerifyOTP(pass string) bool {
	return false
	// return otpservices.VerifyOTP(pass, u.OTPSeed, preloaded.OTPDigits, preloaded.OTPAlgorithm, preloaded.OTPSkew, preloaded.OTPPeriod)
}

// UserGroup !
type UserGroup struct {
	model.Model
	GroupName string `uadmin:"filter"`
}

func (u UserGroup) String() string {
	return u.GroupName
}

// Save !
func (u *UserGroup) Save() {
	database.Save(u)
}

// HasAccess !
func (u *UserGroup) HasAccess(modelName string) GroupPermission {
	utils.Trail(utils.WARNING, "UserGroup.HasAccess will be deprecated in version 0.6.0. Use User.GetAccess instead.")
	return u.hasAccess(modelName)
}

// hasAccess !
func (u *UserGroup) hasAccess(modelName string) GroupPermission {
	up := GroupPermission{}
	dm := menumodel.DashboardMenu{}
	if preloaded.CachePermissions {
		modelID := uint(0)
		for _, m := range cachedModels {
			if m.URL == modelName {
				modelID = m.ID
				break
			}
		}
		for _, g := range cacheGroupPerms {
			if g.UserGroupID == u.ID && g.DashboardMenuID == modelID {
				up = g
				break
			}
		}
	} else {
		database.Get(&dm, "url = ?", modelName)
		database.Get(&up, "user_group_id = ? AND dashboard_menu_id = ?", u.ID, dm.ID)
	}
	return up
}

var cacheUserPerms []UserPermission
var cacheGroupPerms []GroupPermission
var cachedModels []menumodel.DashboardMenu

// UserPermission !
type UserPermission struct {
	model.Model
	DashboardMenu   menumodel.DashboardMenu `uadmin:"filter"`
	DashboardMenuID uint          ``
	User            User          `uadmin:"filter"`
	UserID          uint          ``
	Read            bool          `uadmin:"filter"`
	Add             bool          `uadmin:"filter"`
	Edit            bool          `uadmin:"filter"`
	Delete          bool          `uadmin:"filter"`
	Approval        bool          `uadmin:"filter"`
}

func (u UserPermission) String() string {
	return fmt.Sprint(u.ID)
}

// HideInDashboard to return false and auto hide this from dashboard
func (UserPermission) HideInDashboard() bool {
	return true
}

func LoadPermissions() {
	cacheUserPerms = []UserPermission{}
	cacheGroupPerms = []GroupPermission{}
	cachedModels = []menumodel.DashboardMenu{}
	database.All(&cacheUserPerms)
	database.All(&cacheGroupPerms)
	database.All(&cachedModels)
}

// GroupPermission !
type GroupPermission struct {
	model.Model
	DashboardMenu   menumodel.DashboardMenu `uadmin:"required;filter"`
	DashboardMenuID uint
	UserGroup       UserGroup `uadmin:"required;filter"`
	UserGroupID     uint
	Read            bool `uadmin:"filter"`
	Add             bool `uadmin:"filter"`
	Edit            bool `uadmin:"filter"`
	Delete          bool `uadmin:"filter"`
	Approval        bool `uadmin:"filter"`
}

func (g GroupPermission) String() string {
	return fmt.Sprint(g.ID)
}

// HideInDashboard to return false and auto hide this from dashboard
func (GroupPermission) HideInDashboard() bool {
	return true
}

func NewUserModelFromJson(jsonForUser ...[]byte) (*User, error) {
	if len(jsonForUser) == 0 {
		return &User{}, nil
	} else {
		u := &User{}
		err := json.Unmarshal(jsonForUser[0], u)
		return u, err
	}
}