package models

import (
	"fmt"
	"github.com/uadmin/uadmin/database"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/translation"
)

// Language !
type Language struct {
	model.Model
	EnglishName    string `uadmin:"required;read_only;filter;search"`
	Name           string `uadmin:"required;read_only;filter;search"`
	Flag           string `uadmin:"image;list_exclude"`
	Code           string `uadmin:"filter;read_only;list_exclude"`
	RTL            bool   `uadmin:"list_exclude"`
	Default        bool   `uadmin:"help:Set as the default language;list_exclude"`
	Active         bool   `uadmin:"help:To show this in available languages;filter"`
	AvailableInGui bool   `uadmin:"help:The App is available in this language;read_only"`
}

// DefaultLang is the default language of the system.
var DefaultLang Language

// Global active languages
var ActiveLangs []Language


// String !
func (l Language) String() string {
	return l.Code
}

// Save !
func (l *Language) Save() {
	if l.Default {
		database.Update([]Language{}, "default", false, "\"default\" = ?", true)
		DefaultLang = *l
	}
	database.Save(l)
	tempActiveLangs := []Language{}
	dialect := dialect.GetDialectForDb()
	dialect.Equals("active", true)
	database.Filter(&tempActiveLangs, dialect.ToString(), true)
	ActiveLangs = tempActiveLangs

	tanslationList := []translation.Translation{}
	for i := range ActiveLangs {
		tanslationList = append(tanslationList, translation.Translation{
			Active:  ActiveLangs[i].Active,
			Default: ActiveLangs[i].Default,
			Code:    ActiveLangs[i].Code,
			Name:    fmt.Sprintf("%s (%s)", ActiveLangs[i].Name, ActiveLangs[i].EnglishName),
		})
	}

	for modelName := range model.Schema {
		for i := range model.Schema[modelName].Fields {
			if model.Schema[modelName].Fields[i].Type == preloaded.CMULTILINGUAL {
				// @todo, redo
				// model.Schema[modelName].Fields[i].Translations = tanslationList
			}
		}
	}
}

// GetDefaultLanguage returns the default language
func GetDefaultLanguage() Language {
	return DefaultLang
}

// GetActiveLanguages returns a list of active langages
func GetActiveLanguages() []Language {
	return ActiveLangs
}
