package api

import (
	"context"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	logmodel "github.com/uadmin/uadmin/blueprint/logging/models"
	usermodel "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/metrics"
	"github.com/uadmin/uadmin/utils"
	authservices "github.com/uadmin/uadmin/blueprint/auth/services"
)

// IsAuthenticated returns if the http.Request is authenticated or not
func IsAuthenticated(r *http.Request) *sessionmodel.Session {
	key := GetSession(r)

	if strings.HasPrefix(key, "nouser:") {
		return nil
	}

	s := sessionmodel.Session{}
	if preloaded.CacheSessions {
		s = authservices.CachedSessions[key]
	} else {
		dialect1 := dialect.GetDialectForDb()
		dialect1.Equals("key", key)
		database.Get(&s, dialect1.ToString(), key)
	}
	if isValidSession(r, &s) {
		return &s
	}
	return nil
}

// SetSessionCookie sets the session cookie value, The the value passed in
// session is nil, then the session assiged will be a no user session
func SetSessionCookie(w http.ResponseWriter, r *http.Request, s *sessionmodel.Session) {
	if s == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "nouser:" + authservices.GenerateBase64(24),
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().AddDate(0, 0, 1),
		})
	} else {
		exDate := time.Time{}
		if s.ExpiresOn != nil {
			exDate = *s.ExpiresOn
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    s.Key,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  exDate,
		})
	}
}

func isValidSession(r *http.Request, s *sessionmodel.Session) bool {
	if s != nil && s.ID != 0 {
		if s.Active && !s.PendingOTP && (s.ExpiresOn == nil || s.ExpiresOn.After(time.Now())) {
			if s.User.ID != s.UserID {
				database.Get(&s.User, "id = ?", s.UserID)
			}
			if s.User.Active && (s.User.ExpiresOn == nil || s.User.ExpiresOn.After(time.Now())) {
				// Check for IP restricted session
				if preloaded.RestrictSessionIP {
					ip, _, _ := net.SplitHostPort(r.RemoteAddr)
					return ip == s.IP
				}
				return true
			}
		}
	}
	return false
}

// GetUserFromRequest returns a user from a request
func GetUserFromRequest(r *http.Request) *usermodel.User {
	s := GetSessionFromRequest(r)
	if s != nil {
		if s.User.ID != 0 {
			database.Get(&s.User, "id = ?", s.UserID)
		}
		if s.User.ID != 0 {
			return &s.User
		}
	}
	return nil
}

// getUserFromRequest returns a session from a request
func GetSessionFromRequest(r *http.Request) *sessionmodel.Session {
	key := GetSession(r)
	s := sessionmodel.Session{}

	if preloaded.CacheSessions {
		s = authservices.CachedSessions[key]
	} else {
		dialect1 := dialect.GetDialectForDb()
		dialect1.Equals("key", key)
		database.Get(&s, dialect1.ToString(), key)
	}

	if s.ID != 0 {
		return &s
	}
	return nil
}

// Login return *User and a bool for Is OTP Required
func Login(r *http.Request, username string, password string) (*sessionmodel.Session, bool) {
	// Get the user from DB
	user := usermodel.User{}
	database.Get(&user, "username = ?", username)
	if user.ID == 0 {
		metrics.IncrementMetric("uadmin/security/invalidlogin")
		go func() {
			log := &logmodel.Log{}
			if r.Form == nil {
				r.ParseForm()
			}
			ctx := context.WithValue(r.Context(), preloaded.CKey("login-status"), "invalid username")
			r = r.WithContext(ctx)
			log.SignIn(username, log.Action.LoginDenied(), r)
			log.Save()
		}()
		return nil, false
	}
	s := user.Login(password, "")
	if s != nil && s.ID != 0 {
		s.IP, _, _ = net.SplitHostPort(r.RemoteAddr)
		s.Save()
		if s.Active && (s.ExpiresOn == nil || s.ExpiresOn.After(time.Now())) {
			s.User = user
			if s.User.Active && (s.User.ExpiresOn == nil || s.User.ExpiresOn.After(time.Now())) {
				metrics.IncrementMetric("uadmin/security/validlogin")
				// Store login successful to the user log
				go func() {
					log := &logmodel.Log{}
					if r.Form == nil {
						r.ParseForm()
					}
					log.SignIn(user.Username, log.Action.LoginSuccessful(), r)
					log.Save()
				}()
				return s, s.User.OTPRequired
			}
		}
	} else {
		go func() {
			log := &logmodel.Log{}
			if r.Form == nil {
				r.ParseForm()
			}
			ctx := context.WithValue(r.Context(), preloaded.CKey("login-status"), "invalid password or inactive user")
			r = r.WithContext(ctx)
			log.SignIn(username, log.Action.LoginDenied(), r)
			log.Save()
		}()
	}

	// Increment password attempts and check if it reached
	// the maximum invalid password attempts
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	authservices.InvalidAttempts[ip]++

	if authservices.InvalidAttempts[ip] >= preloaded.PasswordAttempts {
		utils.RateLimitLock.Lock()
		utils.RateLimitMap[ip] = time.Now().Add(time.Duration(preloaded.PasswordTimeout)*time.Minute).Unix() * preloaded.RateLimit
		utils.RateLimitLock.Unlock()
	}

	// Record metrics
	metrics.IncrementMetric("uadmin/security/invalidlogin")
	return nil, false
}

// Login2FA login using username, password and otp for users with OTPRequired = true
func Login2FA(r *http.Request, username string, password string, otpPass string) *sessionmodel.Session {
	s, otpRequired := Login(r, username, password)
	if s != nil {
		if otpRequired && s.User.VerifyOTP(otpPass) {
			s.PendingOTP = false
			s.Save()
		}
		return s
	}
	return nil
}

// Logout logs out a user
func Logout(r *http.Request) {
	s := GetSessionFromRequest(r)
	if s.ID == 0 {
		return
	}

	// Store Logout to the user log
	func() {
		log := &logmodel.Log{}
		log.SignIn(s.User.Username, log.Action.Logout(), r)
		log.Save()
	}()

	s.Logout()

	// Delete the cookie from memory if we sessions are cached
	if preloaded.CacheSessions {
		delete(authservices.CachedSessions, s.Key)
	}

	metrics.IncrementMetric("uadmin/security/logout")
}

// ValidateIP is a function to check if the IP in the request is allowed in the allowed based on allowed
// and block strings
func ValidateIP(r *http.Request, allow string, block string) bool {
	allowed := false
	allowSize := uint32(0)

	allowList := strings.Split(allow, ",")
	for _, net := range allowList {
		if v, size := requestInNet(r, net); v {
			allowed = true
			if size > allowSize {
				allowSize = size
			}
		}
	}

	blockList := strings.Split(block, ",")
	for _, net := range blockList {
		if v, size := requestInNet(r, net); v {
			if size > allowSize {
				allowed = false
				break
			}
		}
	}
	if !allowed {
		metrics.IncrementMetric("uadmin/security/blockedip")
	}
	return allowed
}

func requestInNet(r *http.Request, net string) (bool, uint32) {
	// Check if the IP is V4
	if strings.Contains(r.RemoteAddr, ".") {
		var ip uint32
		var subnet uint32
		var oct uint64
		var mask uint32

		// check if the net is IPv4
		if !strings.Contains(net, ".") && net != "*" && net != "" {
			return false, 0
		}

		// Convert the IP to uint32
		ipParts := strings.Split(strings.Split(r.RemoteAddr, ":")[0], ".")
		for i, o := range ipParts {
			oct, _ = strconv.ParseUint(o, 10, 8)
			ip += uint32(oct << ((3 - uint(i)) * 8))
		}

		// convert the net to uint32
		// but first convert standard nets to IPv4 format
		if net == "*" {
			net = "0.0.0.0/0"
		} else if net == "" {
			net = "255.255.255.255/32"
		} else if !strings.Contains(net, "/") {
			net += "/32"
		}
		ipParts = strings.Split(strings.Split(net, "/")[0], ".")
		for i, o := range ipParts {
			oct, _ = strconv.ParseUint(o, 10, 8)
			subnet += uint32(oct << ((3 - uint(i)) * 8))
		}

		maskLength := getNetSize(r, net)
		mask -= uint32(math.Pow(2, float64(32-maskLength)))
		return ((ip & mask) ^ subnet) == 0, uint32(maskLength)
	}
	// Process IPV6
	var ip1 uint64
	var ip2 uint64
	var subnet1 uint64
	var subnet2 uint64
	var oct uint64
	var mask1 uint64
	var mask2 uint64

	// check if the net is IPv6
	if strings.Contains(net, ".") && net != "*" && net != "" {
		return false, 0
	}

	// Normalize IP
	ipS := r.RemoteAddr              // [::1]:10000
	ipS = strings.Trim(ipS, "[")     // ::1]:10000
	ipS = strings.Split(ipS, "]")[0] // ::1
	if strings.HasPrefix(ipS, "::") {
		ipS = "0" + ipS
	} else if strings.HasSuffix(ipS, "::") {
		ipS = ipS + "0"
	}
	// find and replace ::
	ipParts := strings.Split(ipS, ":")
	ipFinalParts := []uint16{}
	processedDC := false
	for i := range ipParts {
		if ipParts[i] == "" && !processedDC {
			processedDC = true
			for counter := 0; counter < 8-i-(len(ipParts)-(i+1)); counter++ {
				//oct, _ = strconv.ParseUint(ipParts[i], 16, 16)
				ipFinalParts = append(ipFinalParts, uint16(0))
			}
		} else {
			oct, _ = strconv.ParseUint(ipParts[i], 16, 16)
			ipFinalParts = append(ipFinalParts, uint16(oct))
		}
	}

	// Parse the IP into two uint64 variables
	for i := 0; i < 4; i++ {
		oct = uint64(ipFinalParts[i])
		ip1 += uint64((oct << ((3 - uint(i)) * 16)))
	}
	for i := 0; i < 4; i++ {
		oct = uint64(ipFinalParts[i+4])
		ip2 += uint64((oct << ((3 - uint(i)) * 16)))
	}

	subnetv6 := net
	if subnetv6 == "*" {
		subnetv6 = "0::0/0"
	} else if subnetv6 == "" {
		subnetv6 = "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128"
	} else if !strings.Contains(subnetv6, "/") {
		subnetv6 = subnetv6 + "/128"
	}
	maskS := strings.Split(subnetv6, "/")[1]
	subnetv6 = strings.Split(subnetv6, "/")[0]
	if strings.HasPrefix(subnetv6, "::") {
		subnetv6 = "0" + subnetv6
	} else if strings.HasSuffix(subnetv6, "::") {
		subnetv6 = subnetv6 + "0"
	}
	// find and replace ::
	ipParts = strings.Split(subnetv6, ":")
	ipFinalParts = []uint16{}
	processedDC = false
	for i := range ipParts {
		if ipParts[i] == "" && !processedDC {
			processedDC = true
			for counter := 0; counter < 8-i-(len(ipParts)-(i+1)); counter++ {
				//oct, _ = strconv.ParseUint(ipParts[i], 16, 16)
				ipFinalParts = append(ipFinalParts, uint16(0))
			}
		} else {
			oct, _ = strconv.ParseUint(ipParts[i], 16, 16)
			ipFinalParts = append(ipFinalParts, uint16(oct))
		}
	}

	for i := 0; i < 4; i++ {
		oct = uint64(ipFinalParts[i])
		subnet1 += uint64((oct << ((3 - uint(i)) * 16)))
	}
	for i := 0; i < 4; i++ {
		oct = uint64(ipFinalParts[i+4])
		subnet2 += uint64((oct << ((3 - uint(i)) * 16)))
	}

	oct, _ = strconv.ParseUint(maskS, 10, 8)
	maskLength := int(oct)

	maskLength2 := math.Max(float64(maskLength-64), 0)
	maskLength1 := float64(maskLength) - maskLength2

	mask1 -= uint64(math.Pow(2, 64-maskLength1))
	mask2 -= uint64(math.Pow(2, 64-maskLength2))
	if maskLength1 == 0 {
		mask1 = 0
	}
	if maskLength2 == 0 {
		mask2 = 0
	}

	xored1 := (ip1 & mask1) ^ subnet1
	xored2 := (ip2 & mask2) ^ subnet2

	return xored1 == 0 && xored2 == 0, uint32(maskLength)
}

func getNetSize(r *http.Request, net string) int {
	var maskLength int
	var oct uint64

	// Check if the IP is V4
	if strings.Contains(r.RemoteAddr, ".") {
		// Get the Netmask
		oct, _ = strconv.ParseUint(strings.Split(net, "/")[1], 10, 8)
		maskLength = int(oct)
	}
	return maskLength
}

func GetSession(r *http.Request) string {
	key, err := r.Cookie("session")
	if err == nil && key != nil {
		return key.Value
	}
	if r.Method == "GET" && r.FormValue("session") != "" {
		return r.FormValue("session")
	}
	if r.Method == "POST" {
		r.ParseForm()
		if r.FormValue("session") != "" {
			return r.FormValue("session")
		}
	}
	return ""
}

// GetRemoteIP is a function that returns the IP for a remote
// user from a request
func GetRemoteIP(r *http.Request) string {
	var ip string
	var err error

	if ip, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
		return ip
	}
	return r.RemoteAddr
}


