package model

import (
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// Model is the standard struct to be embedded
// in any other struct to make it a model for uadmin
type Model struct {
	gorm.Model
}

// Models is where we keep all registered models
var models map[string]interface{} = make(map[string]interface{})

// getModelName returns the name of a model
func GetModelName(a interface{}) string {
	if val, ok := a.(reflect.Value); ok {
		return GetModelName(val.Interface())
	}
	if val, ok := a.(*reflect.Value); ok {
		return GetModelName(val.Elem().Interface())
	}
	if reflect.TypeOf(a).Kind() == reflect.Ptr {
		return GetModelName(reflect.ValueOf(a).Elem().Interface())
	}
	if reflect.TypeOf(a).Kind() == reflect.Slice {
		return GetModelName(reflect.New(reflect.TypeOf(a).Elem()))
	}
	//if val, ok := a.(reflect.Type); ok {
	//	return strings.ToLower(val.Name())
	//}
	return strings.ToLower(reflect.TypeOf(a).Name())
}

var ModelList []interface{}
