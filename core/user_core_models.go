package core

import (
	"fmt"
	"gorm.io/gorm/clause"
	"time"
)

type User struct {
	Model

	Username             string       `protobuf:"bytes,1,opt,name=Username,proto3" gorm:"uniqueIndex;not null" json:"Username,omitempty" uadmin:"list,search" uadminform:"UsernameOptions"`
	FirstName            string       `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName,omitempty" gorm:"default:''" uadmin:"list,search"`
	LastName             string       `protobuf:"bytes,3,opt,name=LastName,proto3" json:"LastName,omitempty" gorm:"default:''" uadmin:"list,search"`
	Password             string       `protobuf:"bytes,4,opt,name=Password,proto3" json:"Password,omitempty" uadminform:"PasswordOptions" gorm:"default:''"`
	IsPasswordUsable     bool         `gorm:"default:false"`
	Email                string       `protobuf:"bytes,5,opt,name=Email,proto3" gorm:"uniqueIndex;not null" json:"Email,omitempty" uadmin:"list,search"`
	Active               bool         `protobuf:"varint,6,opt,name=Active,proto3" json:"Active,omitempty" gorm:"default:false" uadmin:"list"`
	IsStaff              bool         `json:"IsStaff,omitempty" gorm:"default:false"`
	IsSuperUser          bool         `json:"IsSuperUser,omitempty" gorm:"default:false" uadmin:"list"`
	UserGroups           []UserGroup  `protobuf:"bytes,9,opt,name=UserGroup,proto3" json:"UserGroup,omitempty" gorm:"many2many:user_user_groups;foreignKey:ID;" uadminform:"ChooseFromSelectOptions"`
	Permissions          []Permission `protobuf:"bytes,9,opt,name=UserGroup,proto3" json:"UserGroup,omitempty" gorm:"many2many:user_permissions;foreignKey:ID;" uadminform:"ChooseFromSelectOptions"`
	Photo                string       `protobuf:"bytes,11,opt,name=Photo,proto3" json:"Photo,omitempty" uadminform:"UserPhotoFormOptions" gorm:"default:''"`
	LastLogin            *time.Time   `protobuf:"bytes,12,opt,name=LastLogin,proto3" json:"LastLogin,omitempty" uadminform:"ReadonlyField" uadmin:"list"`
	ExpiresOn            *time.Time   `protobuf:"bytes,13,opt,name=ExpiresOn,proto3" json:"ExpiresOn,omitempty" uadminform:"ReadonlyField"`
	GeneratedOTPToVerify string       `protobuf:"bytes,14,opt,name=GeneratedOTPToVerify,proto3" json:"GeneratedOTPToVerify,omitempty"`
	OTPSeed              string       `protobuf:"bytes,15,opt,name=OTPSeed,proto3" json:"OTPSeed,omitempty"`
	OTPRequired          bool         `protobuf:"bytes,15,opt,name=OTPRequired,proto3" json:"OTPRequired,omitempty" uadminform:"OTPRequiredOptions" gorm:"default:false"`
	Salt                 string       `protobuf:"bytes,16,opt,name=Salt,proto3" json:"Salt,omitempty"`
}

func (u *User) Reset() { *u = User{} }

func (u *User) String() string {
	return fmt.Sprintf("User %s - %s", u.Email, u.FullName())
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

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

func (u *User) BuildPermissionRegistry() *UserPermRegistry {
	userPermRegistry := NewUserPermRegistry()
	userPermRegistry.IsSuperUser = u.IsSuperUser
	if u.IsSuperUser {
		return userPermRegistry
	}
	uadminDatabase := NewUadminDatabase()
	db := uadminDatabase.Db
	var permissions []Permission
	var userGroups []UserGroup
	db.Preload(clause.Associations).Model(u).Association("UserGroups").Find(&userGroups)
	for _, group := range userGroups {
		db.Preload(clause.Associations).Model(&group).Association("Permissions").Find(&permissions)
		for _, permission := range permissions {
			blueprintName := permission.ContentType.BlueprintName
			modelName := permission.ContentType.ModelName
			permBits := permission.PermissionBits
			blueprintPerms := userPermRegistry.GetPermissionForBlueprint(blueprintName, modelName)
			blueprintPerms.AddPermission(permBits)
		}
	}
	db.Preload(clause.Associations).Model(u).Association("Permissions").Find(&permissions)
	for _, permission := range permissions {
		blueprintName := permission.ContentType.BlueprintName
		modelName := permission.ContentType.ModelName
		permBits := permission.PermissionBits
		blueprintPerms := userPermRegistry.GetPermissionForBlueprint(blueprintName, modelName)
		blueprintPerms.AddPermission(permBits)
	}
	uadminDatabase.Close()
	return userPermRegistry
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

// HasAccess returns the user level permission to a model. The modelName
// the the URL of the model
func (u *User) HasAccess(modelName string) Permission {
	Trail(WARNING, "User.HasAccess will be deprecated in version 0.6.0. Use User.GetAccess instead.")
	return u.hasAccess(modelName)
}

// hasAccess returns the user level permission to a model. The modelName
// the the URL of the model
func (u *User) hasAccess(modelName string) Permission {
	up := Permission{}
	//dm := menumodel.DashboardMenu{}
	//if preloaded.CachePermissions {
	//	modelID := uint(0)
	//	for _, m := range cachedModels {
	//		if m.URL == modelName {
	//			modelID = m.ID
	//			break
	//		}
	//	}
	//	for _, p := range cacheUserPerms {
	//		if p.UserID == u.ID && p.DashboardMenuID == modelID {
	//			up = p
	//			break
	//		}
	//	}
	//} else {
	//	database.Get(&dm, "url = ?", modelName)
	//	database.Get(&up, "user_id = ? and dashboard_menu_id = ?", u.ID, dm.ID)
	//}
	return up
}

// GetAccess returns the user's permission to a dashboard menu based on
// their admin status, group and user permissions
func (u *User) GetAccess(modelName string) Permission {
	// Check if the user has permission to a model
	//if u.UserGroup.ID != u.UserGroupID {
	//	database.Preload(u)
	//}
	//uPerm := u.hasAccess(modelName)
	//gPerm := u.UserGroup.hasAccess(modelName)
	perm := Permission{}

	//if gPerm.ID != 0 {
	//	perm.Read = gPerm.Read
	//	perm.Edit = gPerm.Edit
	//	perm.Add = gPerm.Add
	//	perm.Delete = gPerm.Delete
	//	perm.Approval = gPerm.Approval
	//}
	//if uPerm.ID != 0 {
	//	perm.Read = uPerm.Read
	//	perm.Edit = uPerm.Edit
	//	perm.Add = uPerm.Add
	//	perm.Delete = uPerm.Delete
	//	perm.Approval = uPerm.Approval
	//}
	//if u.Admin {
	//	perm.Read = true
	//	perm.Edit = true
	//	perm.Add = true
	//	perm.Delete = true
	//	perm.Approval = true
	//}
	return perm
}

// Validate user when saving from uadmin
func (u User) Validate() (ret map[string]string) {
	//ret = map[string]string{}
	//if u.ID == 0 {
	//	database.Get(&u, "username=?", u.Username)
	//	if u.ID > 0 {
	//		ret["Username"] = "Username is already Taken."
	//	}
	//}
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
	Model
	GroupName   string       `uadmin:"list" gorm:"uniqueIndex;not null"`
	Permissions []Permission `gorm:"foreignKey:ID;many2many:usergroup_permissions;" uadminform:"ChooseFromSelectOptions"`
}

func (u UserGroup) String() string {
	return u.GroupName
}

// Save !
func (u *UserGroup) Save() {
	// database.Save(u)
}

// HasAccess !
func (u *UserGroup) HasAccess(modelName string) Permission {
	// utils.Trail(utils.WARNING, "UserGroup.HasAccess will be deprecated in version 0.6.0. Use User.GetAccess instead.")
	return u.hasAccess(modelName)
}

// hasAccess !
func (u *UserGroup) hasAccess(modelName string) Permission {
	up := Permission{}
	//dm := menumodel.DashboardMenu{}
	//if preloaded.CachePermissions {
	//	modelID := uint(0)
	//	for _, m := range cachedModels {
	//		if m.URL == modelName {
	//			modelID = m.ID
	//			break
	//		}
	//	}
	//	for _, g := range cacheGroupPerms {
	//		if g.UserGroupID == u.ID && g.DashboardMenuID == modelID {
	//			up = g
	//			break
	//		}
	//	}
	//} else {
	//	database.Get(&dm, "url = ?", modelName)
	//	database.Get(&up, "user_group_id = ? AND dashboard_menu_id = ?", u.ID, dm.ID)
	//}
	return up
}

var cacheUserPerms []Permission

// UserPermission !
type Permission struct {
	Model
	Name           string
	ContentType    ContentType
	ContentTypeID  uint
	PermissionBits PermBitInteger
	//Read            bool          `uadmin:"filter"`
	//Add             bool          `uadmin:"filter"`
	//Edit            bool          `uadmin:"filter"`
	//Delete          bool          `uadmin:"filter"`
	//Approval        bool          `uadmin:"filter"`
}

func (m *Permission) String() string {
	return fmt.Sprintf("Permission name %s for content type %s", m.Name, m.ContentType.String())
}

func (m *Permission) ShortDescription() string {
	permission := ProjectPermRegistry.GetPermissionName(m.PermissionBits)
	return fmt.Sprintf("blueprint-%s-model-%s-%s", m.ContentType.BlueprintName, m.ContentType.ModelName, permission)
}

// HideInDashboard to return false and auto hide this from dashboard
func (Permission) HideInDashboard() bool {
	return true
}

func LoadPermissions() {
	cacheUserPerms = []Permission{}
	//database.All(&cacheUserPerms)
	//database.All(&cacheGroupPerms)
	//database.All(&cachedModels)
}

// Action !
type OneTimeActionType int

func (a OneTimeActionType) ResetPassword() OneTimeActionType {
	return 1
}

type OneTimeAction struct {
	Model
	User       User
	UserID     uint
	ExpiresOn  time.Time `gorm:"index"`
	Code       string    `gorm:"uniqueIndex"`
	ActionType OneTimeActionType
	IsUsed     bool `gorm:"default:false"`
}

func (m *OneTimeAction) String() string {
	return fmt.Sprintf("One time action for user %s ", m.User.String())
}

type UserAuthToken struct {
	Model
	User            User
	UserID          uint
	Token           string
	SessionDuration time.Duration
}
