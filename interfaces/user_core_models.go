package interfaces

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"gorm.io/gorm/clause"
	"io"
	math "math"
	math_bits "math/bits"
	"time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type User struct {
	Model

	Username             string     `protobuf:"bytes,1,opt,name=Username,proto3" gorm:"uniqueIndex;not null" json:"Username,omitempty" uadmin:"list,search" uadminform:"UsernameOptions"`
	FirstName            string     `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName,omitempty" gorm:"default:''" uadmin:"list,search"`
	LastName             string     `protobuf:"bytes,3,opt,name=LastName,proto3" json:"LastName,omitempty" gorm:"default:''" uadmin:"list,search"`
	Password             string     `protobuf:"bytes,4,opt,name=Password,proto3" json:"Password,omitempty" uadminform:"PasswordOptions" gorm:"default:''"`
	IsPasswordUsable bool `gorm:"default:false"`
	Email                string     `protobuf:"bytes,5,opt,name=Email,proto3" gorm:"uniqueIndex;not null" json:"Email,omitempty" uadmin:"list,search"`
	Active               bool       `protobuf:"varint,6,opt,name=Active,proto3" json:"Active,omitempty" gorm:"default:false" uadmin:"list"`
	IsStaff bool `json:"IsStaff,omitempty" gorm:"default:false"`
	IsSuperUser bool `json:"IsSuperUser,omitempty" gorm:"default:false" uadmin:"list"`
	UserGroups           []UserGroup  `protobuf:"bytes,9,opt,name=UserGroup,proto3" json:"UserGroup,omitempty" gorm:"many2many:user_user_groups;foreignKey:ID;" uadminform:"ChooseFromSelectOptions"`
	Permissions           []Permission  `protobuf:"bytes,9,opt,name=UserGroup,proto3" json:"UserGroup,omitempty" gorm:"many2many:user_permissions;foreignKey:ID;" uadminform:"ChooseFromSelectOptions"`
	Photo                string     `protobuf:"bytes,11,opt,name=Photo,proto3" json:"Photo,omitempty" uadminform:"UserPhotoFormOptions" gorm:"default:''"`
	LastLogin            *time.Time `protobuf:"bytes,12,opt,name=LastLogin,proto3" json:"LastLogin,omitempty" uadminform:"ReadonlyField" uadmin:"list"`
	ExpiresOn            *time.Time `protobuf:"bytes,13,opt,name=ExpiresOn,proto3" json:"ExpiresOn,omitempty" uadminform:"ReadonlyField"`
	GeneratedOTPToVerify string     `protobuf:"bytes,14,opt,name=GeneratedOTPToVerify,proto3" json:"GeneratedOTPToVerify,omitempty"`
	OTPSeed              string     `protobuf:"bytes,15,opt,name=OTPSeed,proto3" json:"OTPSeed,omitempty"`
	OTPRequired          bool     `protobuf:"bytes,15,opt,name=OTPRequired,proto3" json:"OTPRequired,omitempty" uadminform:"OTPRequiredOptions" gorm:"default:false"`
	Salt                 string     `protobuf:"bytes,16,opt,name=Salt,proto3" json:"Salt,omitempty"`
}

func (m *User) Reset()         { *m = User{} }

func (m *User) String() string {
	return fmt.Sprintf("User %s - %s", m.Email, m.FullName())
}
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ba69147a0c9d872, []int{0}
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return m.Size()
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *User) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *User) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

func (m *User) GetPhoto() string {
	if m != nil {
		return m.Photo
	}
	return ""
}

func (m *User) GetLastLogin() *time.Time {
	if m != nil {
		return m.LastLogin
	}
	return nil
}

func (m *User) GetExpiresOn() *time.Time {
	if m != nil {
		return m.ExpiresOn
	}
	return nil
}

func (m *User) GetGeneratedOTPToVerify() string {
	if m != nil {
		return m.GeneratedOTPToVerify
	}
	return ""
}

func (m *User) GetOTPSeed() string {
	if m != nil {
		return m.OTPSeed
	}
	return ""
}

func (m *User) GetSalt() string {
	if m != nil {
		return m.Salt
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "models.User")
}

func init() {
	proto.RegisterFile("blueprint/user/models/generatemodels.proto", fileDescriptor_5ba69147a0c9d872)
}

var fileDescriptor_5ba69147a0c9d872 = []byte{
	// 377 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcb, 0x6e, 0xe2, 0x30,
	0x14, 0x86, 0xc9, 0x00, 0x01, 0x0c, 0x73, 0x91, 0x85, 0x46, 0xd6, 0x68, 0x14, 0x45, 0xac, 0xd0,
	0x2c, 0x88, 0x34, 0xf3, 0x04, 0x8c, 0x4a, 0x69, 0x25, 0x54, 0xa2, 0x40, 0xbb, 0xe8, 0x2e, 0x24,
	0xa7, 0x60, 0x29, 0x89, 0x23, 0xdb, 0xe9, 0xe5, 0x2d, 0xfa, 0x58, 0x5d, 0xb2, 0xe8, 0xa2, 0xcb,
	0x0a, 0x5e, 0xa4, 0xb2, 0x4d, 0xd2, 0x56, 0x62, 0xe5, 0x7c, 0xff, 0x77, 0xec, 0x1c, 0x27, 0x07,
	0xfd, 0x59, 0x25, 0x05, 0xe4, 0x9c, 0x66, 0xd2, 0x2b, 0x04, 0x70, 0x2f, 0x65, 0x31, 0x24, 0xc2,
	0x5b, 0x43, 0x06, 0x3c, 0x94, 0x60, 0x70, 0x94, 0x73, 0x26, 0x19, 0xb6, 0x0d, 0x0d, 0x9e, 0xeb,
	0xa8, 0x71, 0x29, 0x80, 0xe3, 0x5f, 0xa8, 0xad, 0xd6, 0x2c, 0x4c, 0x81, 0x58, 0xae, 0x35, 0xec,
	0x04, 0x15, 0xe3, 0xdf, 0xa8, 0x73, 0x4a, 0xb9, 0x90, 0x17, 0x4a, 0x7e, 0xd1, 0xf2, 0x3d, 0x50,
	0x3b, 0x67, 0xe1, 0x41, 0xd6, 0xcd, 0xce, 0x92, 0x95, 0xf3, 0x43, 0x21, 0xee, 0x18, 0x8f, 0x49,
	0xc3, 0xb8, 0x92, 0x71, 0x1f, 0x35, 0x27, 0x69, 0x48, 0x13, 0xd2, 0xd4, 0xc2, 0x00, 0xfe, 0x89,
	0xec, 0x71, 0x24, 0xe9, 0x2d, 0x10, 0xdb, 0xb5, 0x86, 0xed, 0xe0, 0x40, 0xaa, 0x7a, 0x1c, 0xa7,
	0x34, 0x23, 0x2d, 0x1d, 0x1b, 0xc0, 0x03, 0xd4, 0x0b, 0x20, 0x65, 0x12, 0xc6, 0x51, 0x04, 0x42,
	0x90, 0xb6, 0x96, 0x9f, 0x32, 0xd5, 0xbd, 0xba, 0xc9, 0x94, 0xb3, 0x22, 0x27, 0x1d, 0xd3, 0x7d,
	0x15, 0x60, 0x17, 0x75, 0x2b, 0x38, 0x3f, 0x21, 0xc8, 0xb5, 0x86, 0xcd, 0xe0, 0x63, 0xa4, 0xde,
	0xec, 0x6f, 0x98, 0x64, 0xa4, 0x6b, 0xfa, 0xd4, 0xa0, 0x4e, 0x55, 0xb7, 0x9c, 0xb1, 0x35, 0xcd,
	0x48, 0xcf, 0x9c, 0x5a, 0x05, 0xca, 0x4e, 0xee, 0x73, 0xca, 0x41, 0xcc, 0x33, 0xf2, 0xd5, 0xd8,
	0x2a, 0xc0, 0x7f, 0x51, 0x7f, 0x7a, 0xf8, 0x29, 0xf1, 0x7c, 0xe9, 0x2f, 0xd9, 0x15, 0x70, 0x7a,
	0xf3, 0x40, 0xbe, 0xe9, 0xc2, 0xa3, 0x0e, 0x13, 0xd4, 0x9a, 0x2f, 0xfd, 0x05, 0x40, 0x4c, 0xbe,
	0xeb, 0xb2, 0x12, 0x31, 0x46, 0x8d, 0x45, 0x98, 0x48, 0xf2, 0x43, 0xc7, 0xfa, 0xf9, 0xff, 0xd9,
	0xd3, 0xce, 0xb1, 0xb6, 0x3b, 0xc7, 0x7a, 0xdd, 0x39, 0xd6, 0xe3, 0xde, 0xa9, 0x6d, 0xf7, 0x4e,
	0xed, 0x65, 0xef, 0xd4, 0xae, 0x47, 0x6b, 0x2a, 0x37, 0xc5, 0x6a, 0x14, 0xb1, 0xd4, 0x2b, 0x42,
	0xf5, 0x11, 0xcb, 0xe5, 0xe8, 0xf4, 0xac, 0x6c, 0x3d, 0x2f, 0xff, 0xde, 0x02, 0x00, 0x00, 0xff,
	0xff, 0xc6, 0x32, 0x62, 0x43, 0x5d, 0x02, 0x00, 0x00,
}

func encodeVarintGeneratemodels(dAtA []byte, offset int, v uint64) int {
	offset -= sovGeneratemodels(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *User) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Username)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	l = len(m.FirstName)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	l = len(m.LastName)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	l = len(m.Password)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	if m.Active {
		n += 2
	}
	l = len(m.Photo)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	//l = len(m.LastLogin)
	//if l > 0 {
	//	n += 1 + l + sovGeneratemodels(uint64(l))
	//}
	//l = len(m.ExpiresOn)
	//if l > 0 {
	//	n += 1 + l + sovGeneratemodels(uint64(l))
	//}
	l = len(m.GeneratedOTPToVerify)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	l = len(m.OTPSeed)
	if l > 0 {
		n += 1 + l + sovGeneratemodels(uint64(l))
	}
	l = len(m.Salt)
	if l > 0 {
		n += 2 + l + sovGeneratemodels(uint64(l))
	}
	return n
}

func sovGeneratemodels(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGeneratemodels(x uint64) (n int) {
	return sovGeneratemodels(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

func skipGeneratemodels(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGeneratemodels
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGeneratemodels
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGeneratemodels
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGeneratemodels
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGeneratemodels
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGeneratemodels
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGeneratemodels        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGeneratemodels          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGeneratemodels = fmt.Errorf("proto: unexpected end of group")
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
	GroupName string `uadmin:"list" gorm:"uniqueIndex;not null"`
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
	Name string
	ContentType ContentType
	ContentTypeID uint
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
	ExpiresOn time.Time `gorm:"index"`
	Code string `gorm:"uniqueIndex"`
	ActionType OneTimeActionType
	IsUsed bool `gorm:"default:false"`
}

func (m *OneTimeAction) String() string {
	return fmt.Sprintf("One time action for user %s ", m.User.String())
}

type UserAuthToken struct {
	Model
	User       User
	UserID     uint
	Token      string
	SessionDuration  time.Duration
}
