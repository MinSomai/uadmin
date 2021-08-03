package interfaces

import (
	"container/list"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/loads"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

// DBSettings !
type DBSettings struct {
	Type     string `json:"type"` // sqlite, mysql
	Name     string `json:"name"` // File/DB name
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type UadminConfigOptions struct {
	Theme string `yaml:"theme"`
	SiteName string `yaml:"site_name"`
	ReportingLevel int `yaml:"reporting_level"`
	ReportTimeStamp bool `yaml:"report_timestamp"`
	DebugDB bool `yaml:"debug_db"`
	PageLength int `yaml:"page_length"`
	MaxImageHeight int `yaml:"max_image_height"`
	MaxImageWidth int `yaml:"max_image_width"`
	MaxUploadFileSize int64 `yaml:"max_upload_file_size"`
	EmailFrom string `yaml:"email_from"`
	EmailUsername string `yaml:"email_username"`
	EmailPassword string `yaml:"email_password"`
	EmailSmtpServer string `yaml:"email_smtp_server"`
	EmailSmtpServerPort int `yaml:"email_smtp_server_port"`
	RootURL string `yaml:"root_url"`
	RootAdminURL string `yaml:"root_admin_url"`
	OTPAlgorithm string `yaml:"otp_algorithm"`
	OTPDigits int `yaml:"otp_digits"`
	OTPPeriod uint `yaml:"otp_period"`
	OTPSkew uint `yaml:"otp_skew"`
	PublicMedia bool `yaml:"public_media"`
	//LogDelete bool `yaml:"log_delete"`
	//LogAdd bool `yaml:"log_add"`
	//LogEdit bool `yaml:"log_edit"`
	//LogRead bool `yaml:"log_read"`
	//CacheTranslation bool `yaml:"cache_translation"`
	AllowedIPs string `yaml:"allowed_ips"`
	BlockedIPs string `yaml:"blocked_ips"`
	RestrictSessionIP bool `yaml:"restrict_session_ip"`
	RetainMediaVersions bool `yaml:"retain_media_versions"`
	RateLimit uint `yaml:"rate_limit"`
	RateLimitBurst uint `yaml:"rate_limit_burst"`
	//APILogRead bool `yaml:"api_log_read"`
	//APILogDelete bool `yaml:"api_log_delete"`
	//APILogAdd bool `yaml:"api_log_add"`
	//APILogEdit bool `yaml:"api_log_edit"`
	LogHTTPRequests bool `yaml:"log_http_requests"`
	HTTPLogFormat string `yaml:"http_log_format"`
	LogTrail bool `yaml:"log_trail"`
	TrailLoggingLevel int `yaml:"trail_logging_level"`
	SystemMetrics bool `yaml:"system_metrics"`
	UserMetrics bool `yaml:"user_metrics"`
	PasswordAttempts int `yaml:"password_attempts"`
	PasswordTimeout int `yaml:"password_timeout"`
	AllowedHosts string `yaml:"allowed_hosts"`
	Logo string `yaml:"logo"`
	FavIcon string `yaml:"fav_icon"`
	AdminCookieName string `yaml:"admin_cookie_name"`
	ApiCookieName string `yaml:"api_cookie_name"`
	SessionDuration int64 `yaml:"session_duration"`
	SecureCookie bool `yaml:"secure_cookie"`
	HttpOnlyCookie bool `yaml:"http_only_cookie"`
	DirectApiSigninByField string `yaml:"direct_api_signin_by_field"`
	DebugTests bool `yaml:"debug_tests"`
	PoweredOnSite string `yaml:"powered_on_site"`
	ForgotCodeExpiration int `yaml:"forgot_code_expiration"`
	DateFormat string `yaml:"date_format"`
	UploadPath string `yaml:"upload_path"`
	DateTimeFormat string `yaml:"datetime_format"`
	TimeFormat string `yaml:"time_format"`
	DateFormatOrder string `yaml:"date_format_order"`
	AdminPerPage int `yaml:"admin_per_page"`
}

type UadminDbOptions struct {
	Default *DBSettings
}

type UadminAuthOptions struct {
	JWT_SECRET_TOKEN string `yaml:"jwt_secret_token"`
	MinUsernameLength int `yaml:"min_username_length"`
	MaxUsernameLength int `yaml:"max_username_length"`
	MinPasswordLength int `yaml:"min_password_length"`
	SaltLength int `yaml:"salt_length"`
	Twofactor_auth_required_for_signin_adapters []string `yaml:"twofactor_auth_required_for_signin_adapters"`
}

type UadminAdminOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
	BindIP string `yaml:"bind_ip"`
}

type UadminApiOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
}

type UadminSwaggerOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
	PathToSpec string `yaml:"path_to_spec"`
	ApiEditorListenPort int `yaml:"api_editor_listen_port"`
}

type UadminConfigurableConfig struct {
	Uadmin *UadminConfigOptions `yaml:"uadmin"`
	Test string `yaml:"test"`
	Db *UadminDbOptions `yaml:"db"`
	Auth *UadminAuthOptions `yaml:"auth"`
	Admin *UadminAdminOptions `yaml:"admin"`
	Api *UadminApiOptions `yaml:"api"`
	Swagger *UadminSwaggerOptions `yaml:"swagger"`
}

type FieldChoice struct {
	DisplayAs string
	Value interface{}
}

type IFieldChoiceRegistryInterface interface {
	IsValidChoice (v interface{}) bool
}

type FieldChoiceRegistry struct {
	IFieldChoiceRegistryInterface
	Choices []*FieldChoice
}

func (fcr *FieldChoiceRegistry) IsValidChoice(v interface{}) bool {
	return false
}

type IFieldFormOptions interface {
	GetName() string
	GetInitial() interface{}
	GetDisplayName() string
	GetValidators() []IValidator
	GetChoices() *FieldChoiceRegistry
	GetHelpText() string
	GetWidgetType() string
	GetReadOnly() bool
}

// Info from config file
type UadminConfig struct {
	ApiSpec *loads.Document
	D *UadminConfigurableConfig
	TemplatesFS embed.FS
	LocalizationFS embed.FS
	RequiresCsrfCheck func(c *gin.Context) bool
	PatternsToIgnoreCsrfCheck *list.List
	ErrorHandleFunc func(int, string, string)
	InTests bool
	FieldFormOptions map[string]IFieldFormOptions
}

func (c *UadminConfig) GetPathToTemplate(templateName string) string {
	return fmt.Sprintf("templates/uadmin/%s/%s.html", c.D.Uadmin.Theme, templateName)
}

func (c *UadminConfig) GetPathToUploadDirectory() string {
	return fmt.Sprintf("%s/%s", os.Getenv("UADMIN_PATH"), c.D.Uadmin.UploadPath)
}

func (c *UadminConfig) AddFieldFormOptions(formOptions IFieldFormOptions) {
	c.FieldFormOptions[formOptions.GetName()] = formOptions
}

func (c *UadminConfig) GetFieldFormOptions(formOptionsName string) IFieldFormOptions {
	ret, _ := c.FieldFormOptions[formOptionsName]
	return ret
}

func (ucc *UadminConfigurableConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawStuff UadminConfigurableConfig
	raw := rawStuff{
		Admin: &UadminAdminOptions{BindIP: "0.0.0.0"},
		Auth: &UadminAuthOptions{SaltLength: 16, Twofactor_auth_required_for_signin_adapters: []string{}},
		Uadmin: &UadminConfigOptions{
			Theme: "default",
			SiteName: "uadmin",
			ReportingLevel: 0,
			ReportTimeStamp: false,
			DebugDB: false,
			PageLength: 100,
			MaxImageHeight: 600,
			MaxImageWidth: 800,
			MaxUploadFileSize: int64(25 * 1024 * 1024),
			RootURL: "/",
			RootAdminURL: "/admin",
			OTPAlgorithm: "sha1",
			OTPDigits: 6,
			OTPPeriod: uint(30),
			OTPSkew: uint(5),
			PublicMedia: false,
			//LogDelete: true,
			//LogAdd: true,
			//LogEdit: true,
			//LogRead: false,
			//CacheTranslation: false,
			AllowedIPs: "*",
			BlockedIPs: "",
			RestrictSessionIP: false,
			RetainMediaVersions: true,
			RateLimit: uint(3),
			RateLimitBurst: uint(3),
			//APILogRead: false,
			//APILogEdit: true,
			//APILogAdd: true,
			//APILogDelete: true,
			LogHTTPRequests: true,
			HTTPLogFormat: "%a %>s %B %U %D",
			LogTrail: false,
			TrailLoggingLevel: 2,
			SystemMetrics: false,
			UserMetrics: false,
			PasswordAttempts: 5,
			PasswordTimeout: 15,
			AllowedHosts: "0.0.0.0,127.0.0.1,localhost,::1",
			Logo: "/static-inbuilt/uadmin/logo.png",
			FavIcon: "/static-inbuilt/uadmin/favicon.ico",
			AdminCookieName: "uadmin-admin",
			ApiCookieName: "uadmin-api",
			SessionDuration: 3600,
			SecureCookie: false,
			HttpOnlyCookie: true,
			DirectApiSigninByField: "username",
			DebugTests: false,
			ForgotCodeExpiration: 10,
			DateFormat: "01/_2/2006",
			DateTimeFormat: "01/_2/2006 15:04",
			TimeFormat: "15:04",
			UploadPath: "uploads",
			DateFormatOrder: "mm/dd/yyyy",
			AdminPerPage: 10,
		},
	}
	// Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*ucc = UadminConfigurableConfig(raw)
	return nil

}

var CurrentConfig *UadminConfig

// Reads info from config file
func NewConfig(file string) *UadminConfig {
	file = os.Getenv("UADMIN_PATH")+"/"+file
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal("Config file is missing: ", file)
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	c := UadminConfig{}
	err = yaml.Unmarshal([]byte(content), &c.D)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	c.FieldFormOptions = make(map[string]IFieldFormOptions)
	c.PatternsToIgnoreCsrfCheck = list.New()
	c.PatternsToIgnoreCsrfCheck.PushBack("/ignorecsrfcheck")
	c.RequiresCsrfCheck = func(c *gin.Context) bool {
		for e := CurrentConfig.PatternsToIgnoreCsrfCheck.Front(); e != nil; e = e.Next() {
			pathToIgnore := e.Value.(string)
			if c.Request.URL.Path == pathToIgnore {
				return false
			}
		}
		return true
	}
	CurrentConfig = &c
	return &c
}

// Reads info from config file
func NewSwaggerSpec(file string) *loads.Document {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal("Config file is missing: ", file)
	}
	doc, err := loads.Spec(file)
	if err != nil {
		panic(err)
	}
	return doc
}

