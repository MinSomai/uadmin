package models

import (
	"fmt"
	"github.com/uadmin/uadmin/dialect"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"
	"net/http"
	"strconv"
	"sync"
	"time"

	usermodels "github.com/uadmin/uadmin/blueprint/user/models"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/database"

)

// ABTestType is the type of the AB testing: model or static
type ABTestType int

// Static is used to do AB testing for static assets (images, js, css, ...)
func (ABTestType) Static() ABTestType {
	return 1
}

// Model is used to do AB testing for model values coming from database
func (ABTestType) Model() ABTestType {
	return 2
}

// ModelList a list of registered models
type ModelList int

// FieldList is a list of fields from schema for a registered model
type FieldList int

var StaticABTests map[string][]struct {
	v     string
	vid   uint
	imp   uint
	click uint
	group string
}

var ModelABTests map[string][]struct {
	v     string
	vid   uint
	fname int
	pk    uint
	imp   uint
	click uint
	group string
}

var AbTestsMutex = sync.Mutex{}
var AbTestCount = 0

// ABTest is a model that stores an A/B test
type ABTest struct {
	model.Model
	Name        string     `uadmin:"required"`
	Type        ABTestType `uadmin:"required"`
	StaticPath  string
	ModelName   ModelList
	Field       FieldList
	PrimaryKey  int
	Active      bool
	Group       string
	ResetABTest string `uadmin:"link"`
}

// Save !
func (a *ABTest) Save() {
	if a.ResetABTest == "" {
		a.ResetABTest = preloaded.RootURL + "api/d/abtest/method/Reset/" + fmt.Sprint(a.ID) + "/?$next=$back"
	}
	database.Save(a)
	AbTestCount = database.Count([]ABTest{}, "active = ?", true)
}

func loadModels(a interface{}, u *usermodels.User) []model.Choice {
	c := []model.Choice{}
	for i, m := range model.ModelList {
		c = append(c, model.Choice{K: uint(i), V: model.GetModelName(m)})
	}
	return c
}

func loadFields(a interface{}, u *usermodels.User) []model.Choice {
	m, ok := a.(ABTest)
	if !ok {
		mp, ok := a.(*ABTest)
		if !ok {
			utils.Trail(utils.ERROR, "loadFields Unable to cast a to ABTest")
			return []model.Choice{}
		}
		m = *mp
	}

	if m.Type != m.Type.Model() {
		return []model.Choice{}
	}

	s := model.Schema[model.GetModelName(model.ModelList[int(m.ModelName)])]
	c := []model.Choice{}
	for i, f := range s.Fields {
		c = append(c, model.Choice{K: uint(i), V: f.Name})
	}
	return c
}

func SyncABTests() {
	// Check if there are stats to save to the DB
	AbTestsMutex.Lock()
	if StaticABTests != nil {
		tx := dialect.GetDB().Begin()
		for _, v := range StaticABTests {
			for i := range v {
				if v[i].imp != 0 || v[i].click != 0 {
					// store results to DB
					tx.Exec("UPDATE ab_test_values SET impressions = impressions + ?, clicks = clicks + ? WHERE id = ?", v[i].imp, v[i].click, v[i].vid)
				}
			}
		}

		for _, v := range ModelABTests {
			for i := range v {
				if v[i].imp != 0 || v[i].click != 0 {
					// store results to DB
					tx.Exec("UPDATE ab_test_values SET impressions = impressions + ?, clicks = clicks + ? WHERE id = ?", v[i].imp, v[i].click, v[i].vid)
				}
			}
		}
		tx.Commit()
	}
	StaticABTests = map[string][]struct {
		v     string
		vid   uint
		imp   uint
		click uint
		group string
	}{}

	ModelABTests = map[string][]struct {
		v     string
		vid   uint
		fname int
		pk    uint
		imp   uint
		click uint
		group string
	}{}

	tests := []ABTest{}
	database.Filter(&tests, "active = ?", true)

	// Process Static AB Tests
	for _, t := range tests {
		if t.Type != t.Type.Static() {
			continue
		}
		values := []ABTestValue{}
		database.Filter(&values, "ab_test_id = ? AND active = ?", t.ID, true)
		tempList := []struct {
			v     string
			vid   uint
			imp   uint
			click uint
			group string
		}{}
		for _, v := range values {
			tempList = append(tempList, struct {
				v     string
				vid   uint
				imp   uint
				click uint
				group string
			}{v: v.Value, vid: v.ID, group: t.Group})
		}
		StaticABTests[t.StaticPath] = tempList
	}

	// Process Models AB Tests
	for _, t := range tests {
		if t.Type != t.Type.Model() {
			continue
		}
		schema := model.Schema[model.GetModelName(model.ModelList[int(t.ModelName)])]
		fName := schema.Fields[int(t.Field)].Name
		values := []ABTestValue{}
		database.Filter(&values, "ab_test_id = ? AND active = ?", t.ID, true)
		tempList := []struct {
			v     string
			vid   uint
			fname int
			pk    uint
			imp   uint
			click uint
			group string
		}{}
		for _, v := range values {
			tempList = append(tempList, struct {
				v     string
				vid   uint
				fname int
				pk    uint
				imp   uint
				click uint
				group string
			}{v: v.Value, vid: v.ID, group: t.Group, pk: uint(t.PrimaryKey), fname: int(t.Field)})
		}
		ModelABTests[schema.ModelName+"__"+fName+"__"+fmt.Sprint(t.PrimaryKey)] = tempList
	}
	AbTestsMutex.Unlock()
}

// ABTestClick is a function to register a click for an ABTest group
func ABTestClick(r *http.Request, group string) {
	go func() {
		abt := GetABT(r)
		var index int
		AbTestsMutex.Lock()
		for k, v := range StaticABTests {
			if len(v) != 0 && v[0].group == group {
				index = abt % len(v)
				v[index].click++
				StaticABTests[k] = v
			}
		}
		for k, v := range ModelABTests {
			if len(v) != 0 && v[0].group == group {
				index = abt % len(v)
				v[index].click++
				ModelABTests[k] = v
			}
		}
		AbTestsMutex.Unlock()
	}()
}

func GetABT(r *http.Request) int {
	c, err := r.Cookie("abt")
	if err != nil || c == nil {
		now := time.Now().AddDate(0, 0, 1)
		/*http.SetCookie(&http.Cookie{
			Name:    "abt",
			Value:   fmt.Sprint(now.Second()),
			Path:    "/",
			Expires: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		})
		*/
		return now.Second()
	}

	v, _ := strconv.ParseInt(c.Value, 10, 64)
	return int(v)
}

// Reset resets the impressions and clicks to 0 based on a specified
// AB Test ID
func (a ABTest) Reset() {
	AbTestsMutex.Lock()
	abtestValue := ABTestValue{}
	database.Update(&abtestValue, "Impressions", 0, "ab_test_id = ?", a.ID)
	database.Update(&abtestValue, "Clicks", 0, "ab_test_id = ?", a.ID)
	AbTestsMutex.Unlock()
}
