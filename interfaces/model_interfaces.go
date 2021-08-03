package interfaces

import (
	"fmt"
	"reflect"
)

type ProjectModelRegistry struct {
	models map[string]interface{}
}

func (pmr *ProjectModelRegistry) RegisterModel(model interface{}) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	modelName := v.Type().Name()
	pmr.models[modelName] = model
}

func (pmr *ProjectModelRegistry) GetModelByName(modelName string) interface{} {
	model, exists := pmr.models[modelName]
	if !exists {
		panic(fmt.Errorf("no model with name %s registered in the project", modelName))
	}
	return model
}

func (pmr *ProjectModelRegistry) GetModelFromInterface(model interface{}) interface{} {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	modelName := v.Type().Name()
	modelI, _ := pmr.models[modelName]
	return modelI
}

var ProjectModels *ProjectModelRegistry

func init() {
	ProjectModels = &ProjectModelRegistry{
		models: make(map[string]interface{}),
	}
	ProjectModels.RegisterModel(&ContentType{})
}