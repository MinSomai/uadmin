package interfaces

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func SetUpStructField(structF reflect.Value, v interface{}) error {
	switch structF.Kind() {
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		v := v.(int64)
		if !structF.OverflowInt(v) {
			structF.SetInt(v)
		} else {
			return fmt.Errorf("can't set field with value %d", v)
		}
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		v := v.(uint64)
		if !structF.OverflowUint(v) {
			structF.SetUint(v)
		} else {
			return fmt.Errorf("can't set field with value %d", v)
		}
	case reflect.Bool:
		v := v.(string)
		structF.SetBool(v != "")
	case reflect.String:
		v := v.(string)
		structF.SetString(v)
	case reflect.Float32:
	case reflect.Float64:
		v := v.(float64)
		structF.SetFloat(v)
	case reflect.Struct:
		switch structF.Interface().(type) {
		case time.Time:
			v := v.(time.Time)
			structF.Set(reflect.ValueOf(v))
		case gorm.DeletedAt:
			v := v.(gorm.DeletedAt)
			structF.Set(reflect.ValueOf(v))
		}
	}
	return nil
}

func GetUadminFieldTypeFromGormField(gormField *schema.Field) UadminFieldType {
	var t UadminFieldType
	switch gormField.FieldType.Kind() {
	case reflect.Bool:
		t = BooleanUadminFieldType
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		t = IntegerUadminFieldType
	case reflect.String:
		t = TextUadminFieldType
	case reflect.Float32:
	case reflect.Float64:
		t = FloatUadminFieldType
	case reflect.Struct:
		//s := gormField.
		//switch s. {
		//case time.Time:
		//	return value.(time.Time)
		//}
		//
	}
	return t
}

func IsTruthyValue(value interface{}) bool {
	r := reflect.TypeOf(value)
	if value == nil {
		return false
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		return s.Len() != 0
	} else if r.Kind() == reflect.Struct {
	} else if typeString == "string" {
		return value != ""
	} else if typeString == "int" {
		return value.(int) != 0
	} else if typeString == "Month" {
		return value.(int) != 0
	}
	return true
}

func TransformValueForWidget(value interface{}) interface{} {
	r := reflect.TypeOf(value)
	if value == nil {
		return ""
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		newSlice := make([]string, 0)
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			newSlice = append(newSlice, TransformValueForWidget(s.Index(i).Interface()).(string))
		}
		return newSlice
	} else if r.Kind() == reflect.Bool {
		return strconv.FormatBool(value.(bool))
	} else if r.Kind() == reflect.Struct {
		s := reflect.ValueOf(value)
		switch s.Interface().(type) {
		case time.Time:
			return value.(time.Time).Format(CurrentConfig.D.Uadmin.DateFormat)
		case gorm.DeletedAt:
			return value.(gorm.DeletedAt).Time.Format(CurrentConfig.D.Uadmin.DateFormat)
		}
		return ""
	} else if r.Kind() == reflect.Ptr {
		// @todo, handle pointer to time.Time
		s := reflect.Indirect(reflect.ValueOf(value))
		if !s.IsValid() {
			return nil
		}
		switch s.Interface().(type) {
		case time.Time:
			return value.(*time.Time)
		}
	} else if typeString == "string" {
		return value
	} else if typeString == "int" {
		return strconv.Itoa(value.(int))
	} else if typeString == "uint" {
		return fmt.Sprint(value.(uint))
	} else if typeString == "Month" {
		return strconv.Itoa(int(value.(time.Month)))
	}
	return value
}

func TransformValueForOperator(value interface{}) interface{} {
	r := reflect.TypeOf(value)
	if value == nil {
		return ""
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		newSlice := make([]string, 0)
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			newSlice = append(newSlice, TransformValueForWidget(s.Index(i).Interface()).(string))
		}
		return newSlice
	} else if r.Kind() == reflect.Bool {
		return strconv.FormatBool(value.(bool))
	} else if r.Kind() == reflect.Struct {
		s := reflect.ValueOf(value)
		switch s.Interface().(type) {
		case time.Time:
			return value.(time.Time).Format(CurrentConfig.D.Uadmin.DateFormat)
		case gorm.DeletedAt:
			return value.(gorm.DeletedAt).Time.Format(CurrentConfig.D.Uadmin.DateFormat)
		}
		return ""
	} else if r.Kind() == reflect.Ptr {
		// @todo, handle pointer to time.Time
		s := reflect.Indirect(reflect.ValueOf(value))
		if !s.IsValid() {
			return nil
		}
		switch s.Interface().(type) {
		case time.Time:
			return value.(*time.Time)
		}
	} else if typeString == "string" {
		boolean, err := strconv.ParseBool(value.(string))
		if err == nil {
			return boolean
		}
		return value
	} else if typeString == "int" {
		return strconv.Itoa(value.(int))
	} else if typeString == "uint" {
		return fmt.Sprint(value.(uint))
	} else if typeString == "Month" {
		return strconv.Itoa(int(value.(time.Month)))
	}
	return value
}

func TransformValueForListDisplay(value interface{}) string {
	r := reflect.TypeOf(value)
	if value == nil {
		return ""
	}
	var typeString string
	if r.Kind() == reflect.Ptr {
		typeString = r.Elem().Name()
	} else {
		typeString = r.Name()
	}
	if r.Kind() == reflect.Slice {
		newSlice := make([]string, 0)
		s := reflect.ValueOf(value)
		for i := 0; i < s.Len(); i++ {
			newSlice = append(newSlice, TransformValueForWidget(s.Index(i).Interface()).(string))
		}
		return strings.Join(newSlice, ",")
	} else if r.Kind() == reflect.Bool {
		return strconv.FormatBool(value.(bool))
	} else if r.Kind() == reflect.Struct {
		s := reflect.ValueOf(value)
		switch s.Interface().(type) {
		case time.Time:
			return value.(time.Time).Format(CurrentConfig.D.Uadmin.DateFormat)
		case gorm.DeletedAt:
			return value.(gorm.DeletedAt).Time.Format(CurrentConfig.D.Uadmin.DateFormat)
		}
		return ""
	} else if r.Kind() == reflect.Ptr {
		// @todo, handle pointer to time.Time
		s := reflect.Indirect(reflect.ValueOf(value))
		if !s.IsValid() {
			return ""
		}
		switch s.Interface().(type) {
		case time.Time:
			return value.(*time.Time).String()
		}
	} else if typeString == "string" {
		return value.(string)
	} else if typeString == "int" {
		return strconv.Itoa(value.(int))
	} else if typeString == "uint" {
		return fmt.Sprint(value.(uint))
	} else if typeString == "Month" {
		return strconv.Itoa(int(value.(time.Month)))
	}
	return value.(string)
}
