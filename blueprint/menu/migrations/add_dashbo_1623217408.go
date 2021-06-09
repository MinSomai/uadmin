package migrations

import (
    menumodel "github.com/uadmin/uadmin/blueprint/menu/models"
    "github.com/uadmin/uadmin/dialect"
)

type Adddashbo_1623217408 struct {
}

func (m Adddashbo_1623217408) GetName() string {
    return "menu.1623217408"
}

func (m Adddashbo_1623217408) GetId() int64 {
    return 1623217408
}

func (m Adddashbo_1623217408) Up() {
    db := dialect.GetDB()
    dashboardmenu := menumodel.DashboardMenu{
        MenuName: "Dashboard Menus",
        URL:      "dashboardmenu",
        Hidden:   false,
        Cat:      "System",
    }
    res := db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Users",
        URL:      "user",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "User Groups",
        URL:      "usergroup",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Sessions",
        URL:      "session",
        Hidden:   true,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "User Permissions",
        URL:      "userpermission",
        Hidden:   true,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Group Permissions",
        URL:      "grouppermission",
        Hidden:   true,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Languages",
        URL:      "language",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Logs",
        URL:      "log",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Settings",
        URL:      "setting",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Setting Categories",
        URL:      "settingcategory",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "Approvals",
        URL:      "approval",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "AB Tests",
        URL:      "abtest",
        Hidden:   false,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
    dashboardmenu = menumodel.DashboardMenu{
        MenuName: "AB Test Values",
        URL:      "abtestvalue",
        Hidden:   true,
        Cat:      "System",
    }
    res = db.Create(&dashboardmenu)
    if res.Error != nil {
        panic(res.Error)
    }
}

func (m Adddashbo_1623217408) Down() {
    db := dialect.GetDB()
    db.Unscoped().Where("1 = 1").Delete(&menumodel.DashboardMenu{})
}

func (m Adddashbo_1623217408) Deps() []string {
    return []string{"menu.1623081544"}
}
