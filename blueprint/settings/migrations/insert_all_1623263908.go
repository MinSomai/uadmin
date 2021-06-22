package migrations

import (
    "fmt"
    settingmodel "github.com/uadmin/uadmin/blueprint/settings/models"
    "github.com/uadmin/uadmin/config"
    "github.com/uadmin/uadmin/dialect"
    "strings"
)

type insert_all_1623263908 struct {
}

func (m insert_all_1623263908) GetName() string {
    return "settings.1623263908"
}

func (m insert_all_1623263908) GetId() int64 {
    return 1623263908
}

func (m insert_all_1623263908) Up() {
    // Check if the uAdmin category is not there and add it
    db := dialect.GetDB()
    var uadminSettingcategory settingmodel.SettingCategory
    db.Model(&settingmodel.SettingCategory{}).Where(&settingmodel.SettingCategory{Name: "uAdmin"}).First(&uadminSettingcategory)
    if uadminSettingcategory.ID == 0 {
        uadminSettingcategory = settingmodel.SettingCategory{Name: "uAdmin"}
        db.Create(&uadminSettingcategory)
    }
    t := settingmodel.DataType(0)

    settings := []settingmodel.Setting{
        {
            Name:         "Theme",
            Value:        config.CurrentConfig.D.Uadmin.Theme,
            DefaultValue: "default",
            DataType:     t.String(),
            Help:         "is the name of the theme used in uAdmin",
        },
        {
            Name:         "Site Name",
            Value:        config.CurrentConfig.D.Uadmin.SiteName,
            DefaultValue: "uAdmin",
            DataType:     t.String(),
            Help:         "is the name of the website that shows on title and dashboard",
        },
        {
            Name:         "Reporting Level",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.ReportingLevel),
            DefaultValue: "0",
            DataType:     t.Integer(),
            Help:         "Reporting level. DEBUG=0, WORKING=1, INFO=2, OK=3, WARNING=4, ERROR=5",
        },
        {
            Name:         "Report Time Stamp",
            Value: fmt.Sprint(config.CurrentConfig.D.Uadmin.ReportTimeStamp),
            DefaultValue: "0",
            DataType:     t.Boolean(),
            Help:         "set this to true to have a time stamp in your logs",
        },
        {
            Name: "Debug DB",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.DebugDB),
            DefaultValue: "0",
            DataType:     t.Boolean(),
            Help:         "prints all SQL statements going to DB",
        },
        {
            Name:         "Page Length",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.PageLength),
            DefaultValue: "100",
            DataType:     t.Integer(),
            Help:         "is the list view max number of records",
        },
        {
            Name:         "Max Image Height",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.MaxImageHeight),
            DefaultValue: "600",
            DataType:     t.Integer(),
            Help:         "sets the maximum height of an Image",
        },
        {
            Name:         "Max Image Width",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.MaxImageWidth),
            DefaultValue: "800",
            DataType:     t.Integer(),
            Help:         "sets the maximum width of an image",
        },
        {
            Name:         "Max Upload File Size",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.MaxUploadFileSize),
            DefaultValue: "26214400",
            DataType:     t.Integer(),
            Help:         "is the maximum upload file size in bytes. 1MB = 1024 * 1024",
        },
        {
            Name:         "Email From",
            Value:        config.CurrentConfig.D.Uadmin.EmailFrom,
            DefaultValue: "",
            DataType:     t.String(),
            Help:         "identifies where the email is coming from",
        },
        {
            Name:         "Email Username",
            Value:        config.CurrentConfig.D.Uadmin.EmailUsername,
            DefaultValue: "",
            DataType:     t.String(),
            Help:         "sets the username of an email",
        },
        {
            Name:         "Email Password",
            Value:        config.CurrentConfig.D.Uadmin.EmailPassword,
            DefaultValue: "",
            DataType:     t.String(),
            Help:         "sets the password of an email",
        },
        {
            Name:         "Email SMTP Server",
            Value:        config.CurrentConfig.D.Uadmin.EmailSmtpServer,
            DefaultValue: "",
            DataType:     t.String(),
            Help:         "sets the name of the SMTP Server in an email",
        },
        {
            Name:         "Email SMTP Server Port",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.EmailSmtpServerPort),
            DefaultValue: "0",
            DataType:     t.Integer(),
            Help:         "sets the port number of an SMTP Server in an email",
        },
        {
            Name:         "Root URL",
            Value:        config.CurrentConfig.D.Uadmin.RootURL,
            DefaultValue: "/",
            DataType:     t.String(),
            Help:         "is where the listener is mapped to",
        },
        {
            Name:         "OTP Algorithm",
            Value:        config.CurrentConfig.D.Uadmin.OTPAlgorithm,
            DefaultValue: "sha1",
            DataType:     t.String(),
            Help:         "is the hashing algorithm of OTP. Other options are sha256 and sha512",
        },
        {
            Name:         "OTP Digits",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.OTPDigits),
            DefaultValue: "6",
            DataType:     t.Integer(),
            Help:         "is the number of digits for the OTP",
        },
        {
            Name:         "OTP Period",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.OTPPeriod),
            DefaultValue: "30",
            DataType:     t.Integer(),
            Help:         "the number of seconds for the OTP to change",
        },
        {
            Name:         "OTP Skew",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.OTPSkew),
            DefaultValue: "5",
            DataType:     t.Integer(),
            Help:         "is the number of minutes to search around the OTP",
        },
        {
            Name: "Public Media",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.PublicMedia),
            DefaultValue: "0",
            DataType:     t.Boolean(),
            Help:         "allows public access to media handler without authentication",
        },
        //{
        //    Name: "Log Delete",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.LogDelete),
        //    DefaultValue: "1",
        //    DataType:     t.Boolean(),
        //    Help:         "adds a log when a record is deleted",
        //},
        //{
        //    Name: "Log Add",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.LogAdd),
        //    DefaultValue: "1",
        //    DataType:     t.Boolean(),
        //    Help:         "adds a log when a record is added",
        //},
        //{
        //    Name: "Log Edit",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.LogEdit),
        //    DefaultValue: "1",
        //    DataType:     t.Boolean(),
        //    Help:         "adds a log when a record is edited",
        //},
        //{
        //    Name: "Log Read",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.LogRead),
        //    DefaultValue: "0",
        //    DataType:     t.Boolean(),
        //    Help:         "adds a log when a record is read",
        //},
        //{
        //    Name: "Cache Translation",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.CacheTranslation),
        //    DefaultValue: "0",
        //    DataType:     t.Boolean(),
        //    Help:         "allows a translation to store data in a cache memory",
        //},
        {
            Name:         "Allowed IPs",
            Value:        config.CurrentConfig.D.Uadmin.AllowedIPs,
            DefaultValue: "*",
            DataType:     t.String(),
            Help: `is a list of allowed IPs to access uAdmin interfrace in one of the following formats:
										- * = Allow all
										- "" = Allow none
							 			- "192.168.1.1" Only allow this IP
										- "192.168.1.0/24" Allow all IPs from 192.168.1.1 to 192.168.1.254
											You can also create a list of the above formats using comma to separate them.
											For example: "192.168.1.1,192.168.1.2,192.168.0.0/24`,
        },
        {
            Name:         "Blocked IPs",
            Value:        config.CurrentConfig.D.Uadmin.BlockedIPs,
            DefaultValue: "",
            DataType:     t.String(),
            Help: `is a list of blocked IPs from accessing uAdmin interfrace in one of the following formats:
										 - "*" = Block all
										 - "" = Block none
										 - "192.168.1.1" Only block this IP
										 - "192.168.1.0/24" Block all IPs from 192.168.1.1 to 192.168.1.254
										 		You can also create a list of the above formats using comma to separate them.
												For example: "192.168.1.1,192.168.1.2,192.168.0.0/24`,
        },
        {
            Name: "Restrict Session IP",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.RestrictSessionIP),
            DefaultValue: "0",
            DataType:     t.Boolean(),
            Help:         "is to block access of a user if their IP changes from their original IP during login",
        },
        {
            Name: "Retain Media Versions",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.RetainMediaVersions),
            DefaultValue: "1",
            DataType:     t.Boolean(),
            Help:         "is to allow the system to keep files uploaded even after they are changed. This allows the system to \"Roll Back\" to an older version of the file",
        },
        {
            Name:         "Rate Limit",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.RateLimit),
            DefaultValue: "3",
            DataType:     t.Integer(),
            Help:         "is the maximum number of requests/second for any unique IP",
        },
        {
            Name:         "Rate Limit Burst",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.RateLimitBurst),
            DefaultValue: "3",
            DataType:     t.Integer(),
            Help:         "is the maximum number of requests for an idle user",
        },
        //{
        //    Name: "API Log Read",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.APILogRead),
        //    DefaultValue: "0",
        //    DataType:     t.Boolean(),
        //    Help:         "APILogRead controls the data API's logging for read commands.",
        //},
        //{
        //    Name: "API Log Edit",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.APILogEdit),
        //    DefaultValue: "1",
        //    DataType:     t.Boolean(),
        //    Help:         "APILogEdit controls the data API's logging for edit commands.",
        //},
        //{
        //    Name: "API Log Add",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.APILogAdd),
        //    DefaultValue: "1",
        //    DataType:     t.Boolean(),
        //    Help:         "APILogAdd controls the data API's logging for add commands.",
        //},
        //{
        //    Name: "API Log Delete",
        //    Value: func(v bool) string {
        //        n := 0
        //        if v {
        //            n = 1
        //        }
        //        return fmt.Sprint(n)
        //    }(config.CurrentConfig.D.Uadmin.APILogDelete),
        //    DefaultValue: "1",
        //    DataType:     t.Boolean(),
        //    Help:         "APILogDelete controls the data API's logging for delete commands.",
        //},
        {
            Name: "Log HTTP Requests",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.LogHTTPRequests),
            DefaultValue: "1",
            DataType:     t.Boolean(),
            Help:         "Logs http requests to syslog",
        },
        {
            Name:         "HTTP Log Format",
            Value:        config.CurrentConfig.D.Uadmin.HTTPLogFormat,
            DefaultValue: "",
            DataType:     t.String(),
            Help: `Is the format used to log HTTP access
									%a: Client IP address
									%{remote}p: Client port
									%A: Server hostname/IP
									%{local}p: Server port
									%U: Path
									%c: All coockies
									%{NAME}c: Cookie named 'NAME'
									%{GET}f: GET request parameters
									%{POST}f: POST request parameters
									%B: Response length
									%>s: Response code
									%D: Time taken in microseconds
									%T: Time taken in seconds
									%I: Request length`,
        },
        {
            Name: "Log Trail",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.LogTrail),
            DefaultValue: "0",
            DataType:     t.Boolean(),
            Help:         "Stores Trail logs to syslog",
        },
        {
            Name:         "Trail Logging Level",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.TrailLoggingLevel),
            DefaultValue: "2",
            DataType:     t.Integer(),
            Help:         "Is the minimum level to be logged into syslog.",
        },
        {
            Name: "System Metrics",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.SystemMetrics),
            DefaultValue: "0",
            DataType:     t.Boolean(),
            Help:         "Enables uAdmin system metrics to be recorded",
        },
        {
            Name: "User Metrics",
            Value: func(v bool) string {
                n := 0
                if v {
                    n = 1
                }
                return fmt.Sprint(n)
            }(config.CurrentConfig.D.Uadmin.UserMetrics),
            DefaultValue: "0",
            DataType:     t.Boolean(),
            Help:         "Enables the user metrics to be recorded",
        },
        {
            Name:         "Password Attempts",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.PasswordAttempts),
            DefaultValue: "5",
            DataType:     t.Integer(),
            Help:         "The maximum number of invalid password attempts before the IP address is blocked for some time from usig the system",
        },
        {
            Name:         "Password Timeout",
            Value:        fmt.Sprint(config.CurrentConfig.D.Uadmin.PasswordTimeout),
            DefaultValue: "5",
            DataType:     t.Integer(),
            Help:         "The maximum number of invalid password attempts before the IP address is blocked for some time from usig the system",
        },
        {
            Name:         "Allowed Hosts",
            Value:        config.CurrentConfig.D.Uadmin.AllowedHosts,
            DefaultValue: "0.0.0.0,127.0.0.1,localhost,::1",
            DataType:     t.String(),
            Help:         "A comma seprated list of allowed hosts for the server to work. The default value if only for development and production domain should be added before deployment",
        },
        {
            Name:         "Logo",
            Value:        config.CurrentConfig.D.Uadmin.Logo,
            DefaultValue: "/static-inbuilt/uadmin/logo.png",
            DataType:     t.Image(),
            Help:         "the main logo that shows on uAdmin UI",
        },
        {
            Name:         "Fav Icon",
            Value:        config.CurrentConfig.D.Uadmin.FavIcon,
            DefaultValue: "/static-inbuilt/uadmin/favicon.ico",
            DataType:     t.File(),
            Help:         "the fav icon that shows on uAdmin UI",
        },
    }

    // Prepare uAdmin Settings
    for i := range settings {
        settings[i].CategoryID = uadminSettingcategory.ID
        settings[i].Code = "uAdmin." + strings.Replace(settings[i].Name, " ", "", -1)
    }
    // Check if the settings exist in the DB
    var s settingmodel.Setting
    sList := []settingmodel.Setting{}
    db.Model(&settingmodel.Setting{}).Where(&settingmodel.Setting{CategoryID: uadminSettingcategory.ID}).Find(&sList)
    tx := dialect.GetDB().Begin()
    for _, setting := range settings {
        s = settingmodel.Setting{}
        for c := range sList {
            if sList[c].Code == setting.Code {
                s = sList[c]
            }
        }
        if s.ID == 0 {
            tx.Create(&setting)
        } else {
            if s.DefaultValue != setting.DefaultValue || s.Help != setting.Help {
                if s.Help != setting.Help {
                    s.Help = setting.Help
                }
                if s.Value == s.DefaultValue {
                    s.Value = setting.DefaultValue
                }
                s.DefaultValue = setting.DefaultValue
                tx.Save(s)
                //s.Save()
            }
        }
    }
    tx.Commit()
}

func (m insert_all_1623263908) Down() {
    db := dialect.GetDB()
    db.Unscoped().Where("1 = 1").Delete(&settingmodel.Setting{})
}

func (m insert_all_1623263908) Deps() []string {
    return []string{"settings.1623082592"}
}
