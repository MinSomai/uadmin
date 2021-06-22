package services

import (
	"crypto/rand"
	"github.com/uadmin/uadmin/debug"
	"golang.org/x/crypto/bcrypt"
	// sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"math/big"
)

// CookieTimeout is the timeout of a login cookie in seconds.
// If the value is -1, then the session cookie will not have
// an expiry date.
var CookieTimeout = -1

// Salt is extra salt added to password hashing
var Salt = ""

// bcryptDiff
var bcryptDiff = 12

// cachedSessions is variable for keeping active sessions
// var CachedSessions map[string]sessionmodel.Session

// invalidAttemps keeps track of invalid password attempts
// per IP address
var InvalidAttempts = map[string]int{}

// GenerateBase64 generates a base64 string of length length
func GenerateBase64(length int) string {
	base := new(big.Int)
	base.SetString("64", 10)

	base64 := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	tempKey := ""
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, base)
		tempKey += string(base64[int(index.Int64())])
	}
	return tempKey
}

//// GenerateBase32 generates a base64 string of length length
//func GenerateBase32(length int) string {
//	base := new(big.Int)
//	base.SetString("32", 10)
//
//	base32 := "234567abcdefghijklmnopqrstuvwxyz"
//	tempKey := ""
//	for i := 0; i < length; i++ {
//		index, _ := rand.Int(rand.Reader, base)
//		tempKey += string(base32[int(index.Int64())])
//	}
//	return tempKey
//}
//
// hashPass Generates a hash from a password and salt
func HashPass(pass string) string {
	password := []byte(pass + Salt)
	hash, err := bcrypt.GenerateFromPassword(password, bcryptDiff)
	if err != nil {
		debug.Trail(debug.ERROR, "uadmin.auth.hashPass.GenerateFromPassword: %s", err)
		return ""
	}
	return string(hash)
}

//func getSessionByKey(key string) *sessionmodel.Session {
//	s := sessionmodel.Session{}
//	if preloaded.CacheSessions {
//		s = CachedSessions[key]
//	} else {
//		database.Get(&s, "`key` = ?", key)
//	}
//	if s.ID == 0 {
//		return nil
//	}
//	return &s
//}
