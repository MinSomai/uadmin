package models

import (
	"fmt"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	usermodel "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/colors"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/metrics"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"
	"strconv"
	"strings"
	"time"
)
// SettingCategory is a category for system settings
type SettingCategory struct {
	interfaces.Model
	Name string
	Icon string `uadmin:"image"`
}

// DataType is a list of data types used for settings
type DataType int

// String is a type
func (DataType) String() DataType {
	return 1
}

// Integer is a type
func (DataType) Integer() DataType {
	return 2
}

// Float is a type
func (DataType) Float() DataType {
	return 3
}

// Boolean is a type
func (DataType) Boolean() DataType {
	return 4
}

// File is a type
func (DataType) File() DataType {
	return 5
}

// Image is a type
func (DataType) Image() DataType {
	return 6
}

// DateTime is a type
func (DataType) DateTime() DataType {
	return 7
}

// Setting model stored system settings
type Setting struct {
	interfaces.Model
	Name         string `uadmin:"required;filter;search"`
	DefaultValue string
	DataType     DataType `uadmin:"required;filter"`
	Value        string
	Help         string          `uadmin:"search" sql:"type:text;"`
	Category     SettingCategory `uadmin:"required;filter"`
	CategoryID   uint
	Code         string `uadmin:"read_only;search"`
}

// Save overides save
func (s *Setting) Save() {
	//database.Preload(s)
	//s.Code = strings.Replace(s.Category.Name, " ", "", -1) + "." + strings.Replace(s.Name, " ", "", -1)
	//s.ApplyValue()
	//database.Save(s)
}

// ParseFormValue takes the value of a setting from an HTTP request and saves in the instance of setting
func (s *Setting) ParseFormValue(v []string) {
	switch s.DataType {
	case s.DataType.Boolean():
		tempV := len(v) == 1 && v[0] == "on"
		if tempV {
			s.Value = "1"
		} else {
			s.Value = "0"
		}
	case s.DataType.DateTime():
		if len(v) == 1 && v[0] != "" {
			s.Value = v[0] + ":00"
		} else {
			s.Value = ""
		}
	default:
		if len(v) == 1 && v[0] != "" {
			s.Value = v[0]
		} else {
			s.Value = ""
		}
	}
}

// GetValue returns an interface representing the value of the setting
func (s *Setting) GetValue() interface{} {
	var err error
	var v interface{}

	switch s.DataType {
	case s.DataType.String():
		if s.Value == "" {
			v = s.DefaultValue
		} else {
			v = s.Value
		}
	case s.DataType.Integer():
		if s.Value != "" {
			v, err = strconv.ParseInt(s.Value, 10, 64)
			v = int(v.(int64))
		}
		if err != nil {
			v, err = strconv.ParseInt(s.DefaultValue, 10, 64)
		}
		if err != nil {
			v = 0
		}
	case s.DataType.Float():
		if s.Value != "" {
			v, err = strconv.ParseFloat(s.Value, 64)
		}
		if err != nil {
			v, err = strconv.ParseFloat(s.DefaultValue, 64)
		}
		if err != nil {
			v = 0.0
		}
	case s.DataType.Boolean():
		if s.Value != "" {
			v = s.Value == "1"
		}
		if v == nil {
			v = s.DefaultValue == "1"
		}
	case s.DataType.File():
		if s.Value == "" {
			v = s.DefaultValue
		} else {
			v = s.Value
		}
	case s.DataType.Image():
		if s.Value == "" {
			v = s.DefaultValue
		} else {
			v = s.Value
		}
	case s.DataType.DateTime():
		if s.Value != "" {
			v, err = time.Parse("2006-01-02 15:04:05", s.Value)
		}
		if err != nil {
			v, err = time.Parse("2006-01-02 15:04:05", s.DefaultValue)
		}
		if err != nil {
			v = time.Now()
		}
	}
	return v
}

// ApplyValue changes uAdmin global variables' value based in the setting value
func (s *Setting) ApplyValue() {
	v := s.GetValue()

	switch s.Code {
	case "uAdmin.Theme":
		preloaded.Theme = strings.Replace(v.(string), "/", "_", -1)
		preloaded.Theme = strings.Replace(preloaded.Theme, "\\", "_", -1)
		preloaded.Theme = strings.Replace(preloaded.Theme, "..", "_", -1)
	case "uAdmin.SiteName":
		preloaded.SiteName = v.(string)
	case "uAdmin.ReportingLevel":
		utils.ReportingLevel = v.(int)
	case "uAdmin.ReportTimeStamp":
		utils.ReportTimeStamp = v.(bool)
	case "uAdmin.DebugDB":
		if preloaded.DebugDB != v.(bool) {
			preloaded.DebugDB = v.(bool)
		}
	case "uAdmin.PageLength":
		preloaded.PageLength = v.(int)
	case "uAdmin.MaxImageHeight":
		preloaded.MaxImageHeight = v.(int)
	case "uAdmin.MaxImageWidth":
		preloaded.MaxImageWidth = v.(int)
	case "uAdmin.MaxUploadFileSize":
		preloaded.MaxUploadFileSize = int64(v.(int))
	case "uAdmin.Port":
		// Port = v.(int)
	case "uAdmin.EmailFrom":
		preloaded.EmailFrom = v.(string)
	case "uAdmin.EmailUsername":
		preloaded.EmailUsername = v.(string)
	case "uAdmin.EmailPassword":
		preloaded.EmailPassword = v.(string)
	case "uAdmin.EmailSMTPServer":
		preloaded.EmailSMTPServer = v.(string)
	case "uAdmin.EmailSMTPServerPort":
		preloaded.EmailSMTPServerPort = v.(int)
	case "uAdmin.RootURL":
		preloaded.RootURL = v.(string)
	case "uAdmin.OTPAlgorithm":
		preloaded.OTPAlgorithm = v.(string)
	case "uAdmin.OTPDigits":
		preloaded.OTPDigits = v.(int)
	case "uAdmin.OTPPeriod":
		preloaded.OTPPeriod = uint(v.(int))
	case "uAdmin.OTPSkew":
		preloaded.OTPSkew = uint(v.(int))
	case "uAdmin.PublicMedia":
		preloaded.PublicMedia = v.(bool)
	case "uAdmin.LogDelete":
		preloaded.LogDelete = v.(bool)
	case "uAdmin.LogAdd":
		preloaded.LogAdd = v.(bool)
	case "uAdmin.LogEdit":
		preloaded.LogEdit = v.(bool)
	case "uAdmin.LogRead":
		preloaded.LogRead = v.(bool)
	case "uAdmin.CacheTranslation":
		preloaded.CacheTranslation = v.(bool)
	case "uAdmin.AllowedIPs":
		preloaded.AllowedIPs = v.(string)
	case "uAdmin.BlockedIPs":
		preloaded.BlockedIPs = v.(string)
	case "uAdmin.RestrictSessionIP":
		preloaded.RestrictSessionIP = v.(bool)
	case "uAdmin.RetainMediaVersions":
		preloaded.RetainMediaVersions = v.(bool)
	case "uAdmin.RateLimit":
		if preloaded.RateLimit != int64(v.(int)) {
			preloaded.RateLimit = int64(v.(int))
			utils.RateLimitMap = map[string]int64{}
		}
	case "uAdmin.RateLimitBurst":
		preloaded.RateLimitBurst = int64(v.(int))
	case "uAdmin.OptimizeSQLQuery":
		preloaded.OptimizeSQLQuery = v.(bool)
	case "uAdmin.APILogRead":
		preloaded.APILogRead = v.(bool)
	case "uAdmin.APILogEdit":
		preloaded.APILogEdit = v.(bool)
	case "uAdmin.APILogAdd":
		preloaded.APILogAdd = v.(bool)
	case "uAdmin.APILogDelete":
		preloaded.APILogDelete = v.(bool)
	case "uAdmin.APILogSchema":
		preloaded.APILogSchema = v.(bool)
	case "uAdmin.LogHTTPRequests":
		preloaded.LogHTTPRequests = v.(bool)
	case "uAdmin.HTTPLogFormat":
		preloaded.HTTPLogFormat = v.(string)
	case "uAdmin.LogTrail":
		preloaded.LogTrail = v.(bool)
	case "uAdmin.TrailLoggingLevel":
		interfaces.TrailLoggingLevel = v.(int)
	case "uAdmin.SystemMetrics":
		metrics.SystemMetrics = v.(bool)
	case "uAdmin.UserMetrics":
		metrics.UserMetrics = v.(bool)
	case "uAdmin.CacheSessions":
		preloaded.CacheSessions = v.(bool)
		if preloaded.CacheSessions {
			sessionmodel.LoadSessions()
		}
	case "uAdmin.CachePermissions":
		preloaded.CachePermissions = v.(bool)
		if preloaded.CachePermissions {
			usermodel.LoadPermissions()
		}
	case "uAdmin.PasswordAttempts":
		preloaded.PasswordAttempts = v.(int)
	case "uAdmin.PasswordTimeout":
		preloaded.PasswordTimeout = v.(int)
	case "uAdmin.AllowedHosts":
		preloaded.AllowedHosts = v.(string)
	case "uAdmin.Logo":
		preloaded.Logo = v.(string)
	case "uAdmin.FavIcon":
		preloaded.FavIcon = v.(string)
	}
}

// GetSetting return the value of a setting based on its code
func GetSetting(code string) interface{} {
	s := Setting{}
	// database.Get(&s, "code = ?", code)

	if s.ID == 0 {
		return nil
	}
	return s.GetValue()
}

func syncSystemSettings() {
	// Check if the uAdmin category is not there and add it
	cat := SettingCategory{}
	// database.Get(&cat, "Name = ?", "uAdmin")
	if cat.ID == 0 {
		cat = SettingCategory{Name: "uAdmin"}
		// database.Save(&cat)
	}

	t := DataType(0)

	settings := []Setting{
		{
			Name:         "Theme",
			Value:        preloaded.Theme,
			DefaultValue: "default",
			DataType:     t.String(),
			Help:         "is the name of the theme used in uAdmin",
		},
		{
			Name:         "Site Name",
			Value:        preloaded.SiteName,
			DefaultValue: "uAdmin",
			DataType:     t.String(),
			Help:         "is the name of the website that shows on title and dashboard",
		},
		{
			Name:         "Reporting Level",
			Value:        fmt.Sprint(utils.ReportingLevel),
			DefaultValue: "0",
			DataType:     t.Integer(),
			Help:         "Reporting level. DEBUG=0, WORKING=1, INFO=2, OK=3, WARNING=4, ERROR=5",
		},
		{
			Name:         "Report Time Stamp",
			Value:        fmt.Sprint(utils.ReportTimeStamp),
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
			}(preloaded.DebugDB),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "prints all SQL statements going to DB",
		},
		{
			Name:         "Page Length",
			Value:        fmt.Sprint(preloaded.PageLength),
			DefaultValue: "100",
			DataType:     t.Integer(),
			Help:         "is the list view max number of records",
		},
		{
			Name:         "Max Image Height",
			Value:        fmt.Sprint(preloaded.MaxImageHeight),
			DefaultValue: "600",
			DataType:     t.Integer(),
			Help:         "sets the maximum height of an Image",
		},
		{
			Name:         "Max Image Width",
			Value:        fmt.Sprint(preloaded.MaxImageWidth),
			DefaultValue: "800",
			DataType:     t.Integer(),
			Help:         "sets the maximum width of an image",
		},
		{
			Name:         "Max Upload File Size",
			Value:        fmt.Sprint(preloaded.MaxUploadFileSize),
			DefaultValue: "26214400",
			DataType:     t.Integer(),
			Help:         "is the maximum upload file size in bytes. 1MB = 1024 * 1024",
		},
		{
			Name:         "Port",
			Value:        fmt.Sprint(8080),
			DefaultValue: "8080",
			DataType:     t.Integer(),
			Help:         "is the port used for http or https server",
		},
		{
			Name:         "Email From",
			Value:        preloaded.EmailFrom,
			DefaultValue: "",
			DataType:     t.String(),
			Help:         "identifies where the email is coming from",
		},
		{
			Name:         "Email Username",
			Value:        preloaded.EmailUsername,
			DefaultValue: "",
			DataType:     t.String(),
			Help:         "sets the username of an email",
		},
		{
			Name:         "Email Password",
			Value:        preloaded.EmailPassword,
			DefaultValue: "",
			DataType:     t.String(),
			Help:         "sets the password of an email",
		},
		{
			Name:         "Email SMTP Server",
			Value:        preloaded.EmailSMTPServer,
			DefaultValue: "",
			DataType:     t.String(),
			Help:         "sets the name of the SMTP Server in an email",
		},
		{
			Name:         "Email SMTP Server Port",
			Value:        fmt.Sprint(preloaded.EmailSMTPServerPort),
			DefaultValue: "0",
			DataType:     t.Integer(),
			Help:         "sets the port number of an SMTP Server in an email",
		},
		{
			Name:         "Root URL",
			Value:        preloaded.RootURL,
			DefaultValue: "/",
			DataType:     t.String(),
			Help:         "is where the listener is mapped to",
		},
		{
			Name:         "OTP Algorithm",
			Value:        preloaded.OTPAlgorithm,
			DefaultValue: "sha1",
			DataType:     t.String(),
			Help:         "is the hashing algorithm of OTP. Other options are sha256 and sha512",
		},
		{
			Name:         "OTP Digits",
			Value:        fmt.Sprint(preloaded.OTPDigits),
			DefaultValue: "6",
			DataType:     t.Integer(),
			Help:         "is the number of digits for the OTP",
		},
		{
			Name:         "OTP Period",
			Value:        fmt.Sprint(preloaded.OTPPeriod),
			DefaultValue: "30",
			DataType:     t.Integer(),
			Help:         "the number of seconds for the OTP to change",
		},
		{
			Name:         "OTP Skew",
			Value:        fmt.Sprint(preloaded.OTPSkew),
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
			}(preloaded.PublicMedia),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "allows public access to media handler without authentication",
		},
		{
			Name: "Log Delete",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.LogDelete),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "adds a log when a record is deleted",
		},
		{
			Name: "Log Add",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.LogAdd),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "adds a log when a record is added",
		},
		{
			Name: "Log Edit",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.LogEdit),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "adds a log when a record is edited",
		},
		{
			Name: "Log Read",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.LogRead),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "adds a log when a record is read",
		},
		{
			Name: "Cache Translation",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.CacheTranslation),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "allows a translation to store data in a cache memory",
		},
		{
			Name:         "Allowed IPs",
			Value:        preloaded.AllowedIPs,
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
			Value:        preloaded.BlockedIPs,
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
			}(preloaded.RestrictSessionIP),
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
			}(preloaded.RetainMediaVersions),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "is to allow the system to keep files uploaded even after they are changed. This allows the system to \"Roll Back\" to an older version of the file",
		},
		{
			Name:         "Rate Limit",
			Value:        fmt.Sprint(preloaded.RateLimit),
			DefaultValue: "3",
			DataType:     t.Integer(),
			Help:         "is the maximum number of requests/second for any unique IP",
		},
		{
			Name:         "Rate Limit Burst",
			Value:        fmt.Sprint(preloaded.RateLimitBurst),
			DefaultValue: "3",
			DataType:     t.Integer(),
			Help:         "is the maximum number of requests for an idle user",
		},
		{
			Name: "Optimize SQL Query",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.OptimizeSQLQuery),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "OptimizeSQLQuery selects columns during rendering a form a list to visible fields.",
		},
		{
			Name: "API Log Read",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.APILogRead),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "APILogRead controls the data API's logging for read commands.",
		},
		{
			Name: "API Log Edit",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.APILogEdit),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "APILogEdit controls the data API's logging for edit commands.",
		},
		{
			Name: "API Log Add",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.APILogAdd),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "APILogAdd controls the data API's logging for add commands.",
		},
		{
			Name: "API Log Delete",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.APILogDelete),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "APILogDelete controls the data API's logging for delete commands.",
		},
		{
			Name: "API Log Schema",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.APILogSchema),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "APILogSchema controls the data API's logging for schema commands.",
		},
		{
			Name: "Log HTTP Requests",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.LogHTTPRequests),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "Logs http requests to syslog",
		},
		{
			Name:         "HTTP Log Format",
			Value:        preloaded.HTTPLogFormat,
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
			}(preloaded.LogTrail),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "Stores Trail logs to syslog",
		},
		{
			Name:         "Trail Logging Level",
			Value:        fmt.Sprint(interfaces.TrailLoggingLevel),
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
			}(metrics.SystemMetrics),
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
			}(metrics.UserMetrics),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "Enables the user metrics to be recorded",
		},
		{
			Name: "Cache Sessions",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.CacheSessions),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "Allows uAdmin to store sessions data in memory",
		},
		{
			Name: "Cache Permissions",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(preloaded.CachePermissions),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "Allows uAdmin to store permissions data in memory",
		},
		{
			Name:         "Password Attempts",
			Value:        fmt.Sprint(preloaded.PasswordAttempts),
			DefaultValue: "5",
			DataType:     t.Integer(),
			Help:         "The maximum number of invalid password attempts before the IP address is blocked for some time from usig the system",
		},
		{
			Name:         "Password Timeout",
			Value:        fmt.Sprint(preloaded.PasswordTimeout),
			DefaultValue: "5",
			DataType:     t.Integer(),
			Help:         "The maximum number of invalid password attempts before the IP address is blocked for some time from usig the system",
		},
		{
			Name:         "Allowed Hosts",
			Value:        preloaded.AllowedHosts,
			DefaultValue: "0.0.0.0,127.0.0.1,localhost,::1",
			DataType:     t.String(),
			Help:         "A comma seprated list of allowed hosts for the server to work. The default value if only for development and production domain should be added before deployment",
		},
		{
			Name:         "Logo",
			Value:        preloaded.Logo,
			DefaultValue: "/static/uadmin/logo.png",
			DataType:     t.Image(),
			Help:         "the main logo that shows on uAdmin UI",
		},
		{
			Name:         "Fav Icon",
			Value:        preloaded.FavIcon,
			DefaultValue: "/static/uadmin/favicon.ico",
			DataType:     t.File(),
			Help:         "the fav icon that shows on uAdmin UI",
		},
	}

	// Prepare uAdmin Settings
	for i := range settings {
		settings[i].CategoryID = cat.ID
		settings[i].Code = "uAdmin." + strings.Replace(settings[i].Name, " ", "", -1)
	}

	// Check if the settings exist in the DB
	var s Setting
	sList := []Setting{}
	// database.Filter(&sList, "category_id = ?", cat.ID)
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	tx := db.Begin()
	for i, setting := range settings {
		interfaces.Trail(interfaces.WORKING, "Synching System Settings: [%s%d/%d%s]", colors.FGGreenB, i+1, len(settings), colors.FGNormal)
		s = Setting{}
		for c := range sList {
			if sList[c].Code == setting.Code {
				s = sList[c]
			}
		}
		if s.ID == 0 {
			tx.Create(&setting)
			//setting.Save()
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
	interfaces.Trail(interfaces.OK, "Synching System Settings: [%s%d/%d%s]", colors.FGGreenB, len(settings), len(settings), colors.FGNormal)
	applySystemSettings()
	preloaded.SettingsSynched = true
}

func applySystemSettings() {
	_ = SettingCategory{}
	settings := []Setting{}

	//database.Get(&cat, "name = ?", "uAdmin")
	//database.Filter(&settings, "category_id = ?", cat.ID)

	for _, setting := range settings {
		setting.ApplyValue()
	}
}
