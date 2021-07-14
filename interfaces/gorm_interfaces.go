package interfaces

import "gorm.io/gorm"

// Model is the standard struct to be embedded
// in any other struct to make it a model for uadmin
type Model struct {
	gorm.Model
}

type ContentType struct {
	Model
	BlueprintName string `sql:"unique_index:idx_contenttype_content_type"`
	ModelName string `sql:"unique_index:idx_contenttype_content_type"`
}

