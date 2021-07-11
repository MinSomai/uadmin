package migrations

import (
    usermodels "github.com/uadmin/uadmin/blueprint/user/models"
    "github.com/uadmin/uadmin/interfaces"
)

type adding_use_1623259185 struct {
}

func (m adding_use_1623259185) GetName() string {
    return "user.1623259185"
}

func (m adding_use_1623259185) GetId() int64 {
    return 1623259185
}

func (m adding_use_1623259185) Up() {
    db := interfaces.GetDB()
    var superuserGroup usermodels.UserGroup
    db.Model(&usermodels.UserGroup{}).Where(&usermodels.UserGroup{GroupName: "Superusers"}).First(&superuserGroup)
    if superuserGroup.ID == 0 {
        superuserGroup = usermodels.UserGroup{
            GroupName: "Superusers",
        }
        db.Create(&superuserGroup)
    }
}

func (m adding_use_1623259185) Down() {
    db := interfaces.GetDB()
    var superuserGroup usermodels.UserGroup
    db.Model(&usermodels.UserGroup{}).Where(&usermodels.UserGroup{GroupName: "Superusers"}).First(&superuserGroup)
    db.Delete(&superuserGroup)
}

func (m adding_use_1623259185) Deps() []string {
    return []string{"user.1621680132"}
}
