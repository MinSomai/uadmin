package config

import (
	"github.com/go-openapi/loads"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
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

// Info from config file
type UadminConfig struct {
	ApiSpec *loads.Document
	D struct {
		Uadmin struct {
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
			OTPAlgorithm string `yaml:"otp_algorithm"`
			OTPDigits int `yaml:"otp_digits"`
			OTPPeriod uint `yaml:"otp_period"`
			OTPSkew uint `yaml:"otp_skew"`
			PublicMedia bool `yaml:"public_media"`
			LogDelete bool `yaml:"log_delete"`
			LogAdd bool `yaml:"log_add"`
			LogEdit bool `yaml:"log_edit"`
			LogRead bool `yaml:"log_read"`
			CacheTranslation bool `yaml:"cache_translation"`
			AllowedIPs string `yaml:"allowed_ips"`
			BlockedIPs string `yaml:"blocked_ips"`
			RestrictSessionIP bool `yaml:"restrict_session_ip"`
			RetainMediaVersions bool `yaml:"retain_media_versions"`
			RateLimit uint `yaml:"rate_limit"`
			RateLimitBurst uint `yaml:"rate_limit_burst"`
			APILogRead bool `yaml:"api_log_read"`
			APILogDelete bool `yaml:"api_log_delete"`
			APILogAdd bool `yaml:"api_log_add"`
			APILogEdit bool `yaml:"api_log_edit"`
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
		} `yaml:"uadmin"`
		Test string `yaml:"test"`
		Db   struct {
			Default *DBSettings
		} `yaml:"db"`
		Auth struct {
			JWT_SECRET_TOKEN string `yaml:"jwt_secret_token"`
			MinUsernameLength int `yaml:"min_username_length"`
			MaxUsernameLength int `yaml:"max_username_length"`
			MinPasswordLength int `yaml:"min_password_length"`
			SaltLength int `yaml:"salt_length"`
		} `yaml:"auth"`
		Admin struct {
			ListenPort int `yaml:"listen_port"`
			SSL        struct {
				ListenPort int `yaml:"listen_port"`
			} `yaml:"ssl"`
		} `yaml:"admin"`
		Api struct {
			ListenPort int `yaml:"listen_port"`
			SSL        struct {
				ListenPort int `yaml:"listen_port"`
			} `yaml:"ssl"`
		} `yaml:"api"`
		Swagger struct {
			ListenPort int `yaml:"listen_port"`
			SSL        struct {
				ListenPort int `yaml:"listen_port"`
			} `yaml:"ssl"`
			PathToSpec string `yaml:"path_to_spec"`
			ApiEditorListenPort int `yaml:"api_editor_listen_port"`
		} `yaml:"swagger"`
	}
}

var CurrentConfig *UadminConfig

// Reads info from config file
func NewConfig(file string) *UadminConfig {
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
	if c.D.Auth.SaltLength == 0 {
		c.D.Auth.SaltLength = 16
	}
	if c.D.Uadmin.Theme == "" {
		c.D.Uadmin.Theme = "default"
	}
	if c.D.Uadmin.SiteName == "" {
		c.D.Uadmin.SiteName = "uAdmin"
	}
	if c.D.Uadmin.PageLength == 0 {
		c.D.Uadmin.PageLength = 100
	}
	if c.D.Uadmin.MaxImageHeight == 0 {
		c.D.Uadmin.MaxImageHeight = 600
	}
	if c.D.Uadmin.MaxImageWidth == 0 {
		c.D.Uadmin.MaxImageWidth = 800
	}
	if c.D.Uadmin.MaxUploadFileSize == 0 {
		c.D.Uadmin.MaxUploadFileSize = int64(25 * 1024 * 1024)
	}
	if c.D.Uadmin.RootURL == "" {
		c.D.Uadmin.RootURL = "/"
	}
	if c.D.Uadmin.OTPAlgorithm == "" {
		c.D.Uadmin.OTPAlgorithm = "sha1"
	}
	if c.D.Uadmin.OTPDigits == 0 {
		c.D.Uadmin.OTPDigits = 6
	}
	if c.D.Uadmin.OTPPeriod == 0 {
		c.D.Uadmin.OTPPeriod = uint(30)
	}
	if c.D.Uadmin.OTPSkew == 0 {
		c.D.Uadmin.OTPSkew = uint(5)
	}
	if c.D.Uadmin.AllowedIPs == "" {
		c.D.Uadmin.AllowedIPs = "*"
	}
	if c.D.Uadmin.RateLimit == 0 {
		c.D.Uadmin.RateLimit = uint(3)
	}
	if c.D.Uadmin.RateLimitBurst == 0 {
		c.D.Uadmin.RateLimitBurst = uint(3)
	}
	if c.D.Uadmin.HTTPLogFormat == "" {
		c.D.Uadmin.HTTPLogFormat = "%a %>s %B %U %D"
	}
	if c.D.Uadmin.PasswordAttempts == 0 {
		c.D.Uadmin.PasswordAttempts = 5
	}
	if c.D.Uadmin.PasswordTimeout == 0 {
		c.D.Uadmin.PasswordTimeout = 15
	}
	if c.D.Uadmin.AllowedHosts == "" {
		c.D.Uadmin.AllowedHosts = "0.0.0.0,127.0.0.1,localhost,::1"
	}
	if c.D.Uadmin.Logo == "" {
		c.D.Uadmin.Logo = "/static/uadmin/logo.png"
	}
	if c.D.Uadmin.FavIcon == "" {
		c.D.Uadmin.FavIcon = "/static/uadmin/favicon.ico"
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
