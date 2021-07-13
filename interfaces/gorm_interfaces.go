package interfaces

import "gorm.io/gorm"

// Model is the standard struct to be embedded
// in any other struct to make it a model for uadmin
type Model struct {
	gorm.Model
}

type ContentType struct {
	Model
	BlueprintName string
	ModelName string
}

