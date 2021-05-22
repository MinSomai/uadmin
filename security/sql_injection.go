package security

import (
	authapi "github.com/uadmin/uadmin/blueprint/auth/api"
	usermodel "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/utils"
	"net"
	"net/http"
	"strings"
)

// SQLInjection is the function to check for SQL injection attacks.
// Parameters:
//  -key: column_name, table name
//  -value: WHERE key(OP)value, SET key=value, VALUES (key,key...)
// return true for sql injection attempt and false for safe requests
func SQLInjection(r *http.Request, key, value string) bool {
	var err error

	user := authapi.GetUserFromRequest(r)
	if user == nil {
		user = &usermodel.User{}
	}
	ip := r.RemoteAddr
	if ip, _, err = net.SplitHostPort(ip); err != nil {
		ip = r.RemoteAddr
	}
	errMsg := "SQL Injection attempt (%s '%s'). User:" + user.Username + " IP:" + ip
	if key != "" {
		// Case 1 - Comment injection
		if strings.Contains(key, "--") || strings.Contains(key, "#") {
			utils.Trail(utils.CRITICAL, errMsg, "comment injection", key)
			return true
		}
		// Case 2 - Comment injection
		if strings.Contains(key, "/*") || strings.Contains(key, "*/") {
			utils.Trail(utils.CRITICAL, errMsg, "comment injection", key)
			return true
		}
		// Case 3 - Stacking
		if strings.Contains(key, ";") {
			utils.Trail(utils.CRITICAL, errMsg, "stacking", key)
			return true
		}
		// Case 4 - HEX Injection
		if strings.Contains(key, "0x") {
			utils.Trail(utils.CRITICAL, errMsg, "hex injection", key)
			return true
		}
		// Case 5 - Concatenation
		if strings.Contains(key, "+") || strings.Contains(key, "||") {
			utils.Trail(utils.CRITICAL, errMsg, "concatenation", key)
			return true
		}
		// Case 6 - Functions
		if strings.Contains(key, "(") || strings.Contains(key, ")") {
			utils.Trail(utils.CRITICAL, errMsg, "functions", key)
			return true
		}
		// Case 7 - Sapce
		if strings.Contains(key, " ") {
			utils.Trail(utils.CRITICAL, errMsg, "space", key)
			return true
		}
		// Case 8 - Escaping
		if strings.Contains(key, "'") || strings.Contains(key, "`") {
			utils.Trail(utils.CRITICAL, errMsg, "escaping", key)
			return true
		}
		// Case 9 - Escaping
		if strings.Contains(key, "'") || strings.Contains(key, "`") {
			utils.Trail(utils.CRITICAL, errMsg, "escaping", key)
			return true
		}
	}
	// if value != ""
	// We are depending on gorm in here
	return false
}