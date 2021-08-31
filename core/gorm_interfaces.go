package core

import (
	"fmt"
	"gorm.io/gorm"
)

// Model is the standard struct to be embedded
// in any other struct to make it a model for uadmin
type Model struct {
	gorm.Model
}

func (m *Model) GetID() uint {
	return m.ID
}

type ContentType struct {
	Model
	BlueprintName string `sql:"unique_index:idx_contenttype_content_type"`
	ModelName     string `sql:"unique_index:idx_contenttype_content_type"`
}

func (ct *ContentType) String() string {
	return fmt.Sprintf("Content type for blueprint %s and model name %s", ct.BlueprintName, ct.ModelName)
}
