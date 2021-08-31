package models

import (
	"fmt"
	"github.com/uadmin/uadmin/core"
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

func HumanizeAbTestType(abTestType ABTestType) string {
	switch abTestType {
	case 1:
		return "static"
	case 2:
		return "model"
	default:
		return "unknown"
	}
}

// ABTest is a model that stores an A/B test
type ABTest struct {
	core.Model
	Name          string           `uadminform:"RequiredFieldOptions" uadmin:"list"`
	Type          ABTestType       `uadminform:"RequiredSelectFieldOptions" uadmin:"list"`
	StaticPath    string           `uadmin:"list"`
	ContentType   core.ContentType `uadmin:"list" uadminform:"ContentTypeFieldOptions"`
	ContentTypeID uint
	Field         string `uadmin:"list" uadminform:"RequiredSelectFieldOptions"`
	PrimaryKey    uint   `uadmin:"list" gorm:"default:0"`
	Active        bool   `gorm:"default:false" uadmin:"list"`
	Group         string `uadmin:"list"`
}

func (m *ABTest) String() string {
	return fmt.Sprintf("ABTest %s", m.Name)
}
