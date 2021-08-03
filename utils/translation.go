package utils

import (
	"encoding/json"
	"fmt"
	"github.com/uadmin/uadmin/interfaces"

	// "encoding/json"
	"github.com/gin-gonic/gin"
	langmodel "github.com/uadmin/uadmin/blueprint/language/models"
	"strings"
)

// GetLanguage returns the language of the request
func GetLanguage(c *gin.Context) *langmodel.Language {
	langCookie, err := c.Cookie("language")
	if err != nil || langCookie == "" {
		return GetDefaultLanguage()
	}
	var lang langmodel.Language
	uadminDatabase := interfaces.NewUadminDatabase()
	db := uadminDatabase.Db
	db.Model(langmodel.Language{}).Where(&langmodel.Language{Code: langCookie}).First(&lang)
	uadminDatabase.Close()
	return &lang
}

// GetDefaultLanguage returns the default language
func GetDefaultLanguage() *langmodel.Language {
	if defaultLang != nil {
		return defaultLang
	}
	var lang langmodel.Language
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	db.Model(langmodel.Language{}).Where(&langmodel.Language{Default: true}).First(&lang)
	defaultLang = &lang
	return &lang
}

// GetActiveLanguages returns a list of active langages
func GetActiveLanguages() []langmodel.Language {
	var langs []langmodel.Language
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	db := uadminDatabase.Db
	db.Model(langmodel.Language{}).Where(&langmodel.Language{Active: true}).Find(&langs)
	return langs
}

// DefaultLang is the default language of the system.
var defaultLang *langmodel.Language
type translationLoaded map[string]string
var langMapCache map[string]translationLoaded
const translateMe = "Translate me ---> "

// @todo, redo
// Tf is a function for translating strings into any given language
// Parameters:
// ===========
//   path (string): This is where to get the translation from. It is in the
//                  format of "GROUPNAME/FILENAME" for example: "uadmin/system"
//   lang (string): Is the language code. If empty string is passed we will use
//                  the default language.
//   term (string): The term to translate.
//   args (...interface{}): Is a list of args to fill the term with place holders
func Tf(path string, lang string, term string, args ...interface{}) string {
	if lang == "" {
		lang = GetDefaultLanguage().Code
	}

	// Check if the path if for an existing model schema
	pathParts := strings.Split(path, "/")
	isSchemaFile := false
	if len(pathParts) > 2 {
		path = strings.Join(pathParts[0:2], "/")
		isSchemaFile = true
	}
	if langMapCache == nil {
		langMapCache = make(map[string]translationLoaded)
	}
	langMap, ok := langMapCache[lang]
	if !ok {
		langFile, err := interfaces.CurrentConfig.LocalizationFS.ReadFile(fmt.Sprintf("localization/%s.json", lang))
		if err != nil {
			Trail(ERROR, "Unable to unmarshal json file with language (%s)", err)
		} else {
			err = json.Unmarshal(langFile, &langMap)
			if err != nil {
				Trail(ERROR, "Unable to unmarshal json file with language (%s)", err)
			} else {
				langMapCache[lang] = langMap
			}
		}
	}
	// If the term exists, then return it
	if val, ok := langMap[term]; ok {
		return strings.TrimPrefix(val, translateMe)
	}
	if !isSchemaFile {
		// If the term exists, then return it
		if val, ok := langMap[term]; ok {
			return strings.TrimPrefix(val, translateMe)
		}

		// If it doesn't exist then add it to the file
		if lang != "en" {
			langMap[term] = translateMe + term
			Trail(WARNING, "Unknown term %s", term)
			return translateMe + term
		} else {
			langMap[term] = term
			return term
		}
	} else {
	}
	return term
}

func Translate(c *gin.Context, raw string, lang string, args ...bool) string {
	var langParser map[string]json.RawMessage
	err := json.Unmarshal([]byte(raw), &langParser)
	if err != nil {
		return raw
	}
	transtedStr := string(langParser[lang])

	if len(transtedStr) > 2 {
		return transtedStr[1 : len(transtedStr)-1]
	}
	if len(args) > 0 && !args[0] {
		return ""
	}
	language := GetLanguage(c)
	transtedStr = string(langParser[language.Code])

	if len(transtedStr) > 2 {
		return transtedStr[1 : len(transtedStr)-1]
	}
	return ""
}
