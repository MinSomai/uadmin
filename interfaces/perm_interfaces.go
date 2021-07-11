package interfaces

type PermissionRegistry interface {
	DoesUserHaveRightFor(permissionName string) bool
	AddCustomPermission(permission CustomPermission)
}
type CustomPermission string

type Perm struct {
	PermissionRegistry
	PermBitInteger PermBitInteger
	CustomPermissions []CustomPermission
}

func (ap *Perm) HasReadPermission() bool {
	return (ap.PermBitInteger & ReadPermBit) == ReadPermBit
}

func (ap *Perm) DoesUserHaveRightFor(permissionName CustomPermission) bool {
	for _, a := range ap.CustomPermissions {
		if a == permissionName {
			return true
		}
	}
	return false
}

func (ap *Perm) AddCustomPermission(permission CustomPermission){
	ap.CustomPermissions = append(ap.CustomPermissions, permission)
}

func (ap *Perm) HasAddPermission() bool {
	return (ap.PermBitInteger & AddPermBit) == AddPermBit
}

func (ap *Perm) HasEditPermission() bool {
	return (ap.PermBitInteger & EditPermBit) == EditPermBit
}

func (ap *Perm) HasDeletePermission() bool {
	return (ap.PermBitInteger & DeletePermBit) == DeletePermBit
}

func (ap *Perm) HasPublishPermission() bool {
	return (ap.PermBitInteger & PublishPermBit) == PublishPermBit
}

func (ap *Perm) HasRevertPermission() bool {
	return (ap.PermBitInteger & RevertPermBit) == RevertPermBit
}

type PermBitInteger int

const ReadPermBit PermBitInteger = 0
const AddPermBit PermBitInteger = 2
const EditPermBit PermBitInteger = 4
const DeletePermBit PermBitInteger = 8
const PublishPermBit PermBitInteger = 16
const RevertPermBit PermBitInteger = 32

func NewPerm(permBitInteger PermBitInteger, customPermissions ...CustomPermission) *Perm {
	return &Perm{
		PermBitInteger: permBitInteger,
		CustomPermissions: customPermissions,
	}
}


