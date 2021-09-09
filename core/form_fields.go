package core

import (
	"fmt"
	"gorm.io/gorm/schema"
	"mime/multipart"
	"reflect"
	"sort"
)

type IFieldRegistry interface {
	GetByName(name string) (*Field, error)
	AddField(field *Field)
	GetAllFields() map[string]*Field
	GetAllFieldsWithOrdering() []*Field
	GetPrimaryKey() (*Field, error)
}

type UadminFieldType string

const BigIntegerUadminFieldType UadminFieldType = "biginteger"
const BinaryUadminFieldType UadminFieldType = "binary"
const BooleanUadminFieldType UadminFieldType = "boolean"
const CharUadminFieldType UadminFieldType = "char"
const DateUadminFieldType UadminFieldType = "date"
const DateTimeUadminFieldType UadminFieldType = "datetime"
const DecimalUadminFieldType UadminFieldType = "decimal"
const DurationUadminFieldType UadminFieldType = "duration"
const EmailUadminFieldType UadminFieldType = "email"
const FileUadminFieldType UadminFieldType = "file"
const FilePathUadminFieldType UadminFieldType = "filepath"
const FloatUadminFieldType UadminFieldType = "float"
const ForeignKeyUadminFieldType UadminFieldType = "foreignkey"
const ImageFieldUadminFieldType UadminFieldType = "imagefield"
const IntegerUadminFieldType UadminFieldType = "integer"
const UintUadminFieldType UadminFieldType = "uint"
const IPAddressUadminFieldType UadminFieldType = "ipaddress"
const GenericIPAddressUadminFieldType UadminFieldType = "genericipaddress"
const ManyToManyUadminFieldType UadminFieldType = "manytomany"
const NullBooleanUadminFieldType UadminFieldType = "nullboolean"
const PositiveBigIntegerUadminFieldType UadminFieldType = "positivebiginteger"
const PositiveIntegerUadminFieldType UadminFieldType = "positiveinteger"
const PositiveSmallIntegerUadminFieldType UadminFieldType = "positivesmallinteger"
const SlugUadminFieldType UadminFieldType = "slug"
const SmallIntegerUadminFieldType UadminFieldType = "smallinteger"
const TextUadminFieldType UadminFieldType = "text"
const TimeUadminFieldType UadminFieldType = "time"
const URLUadminFieldType UadminFieldType = "url"
const UUIDUadminFieldType UadminFieldType = "uuid"

type FieldConfig struct {
	Widget                 IWidget
	AutocompleteURL        string
	DependsOnAnotherFields []string
}

type Field struct {
	schema.Field
	ReadOnly        bool
	UadminFieldType UadminFieldType
	FieldConfig     *FieldConfig
	Required        bool
	DisplayName     string
	HelpText        string
	Choices         *FieldChoiceRegistry
	Validators      *ValidatorRegistry
	SortingDisabled bool
	Populate        func(field *Field, m interface{}) interface{}
	Initial         interface{}
	WidgetType      string
	SetUpField      func(w IWidget, modelI interface{}, v interface{}, afo IAdminFilterObjects) error
	Ordering        int
}

func (f *Field) ProceedForm(form *multipart.Form, afo IAdminFilterObjects, renderContext *FormRenderContext) ValidationError {
	err := f.FieldConfig.Widget.ProceedForm(form, afo, renderContext)
	if err == nil {
		validationErrors := make(ValidationError, 0)
		for validator := range f.Validators.GetAllValidators() {
			validationErr := validator(f.FieldConfig.Widget.GetOutputValue(), form)
			if validationErr == nil {
				continue
			}
			validationErrors = append(validationErrors, validationErr)
		}
		if len(validationErrors) > 0 {
			f.FieldConfig.Widget.SetErrors(validationErrors)
		}
		return validationErrors
	}
	errors := ValidationError{err}
	f.FieldConfig.Widget.SetErrors(errors)
	return errors
}

type ValidationError []error

type FieldFormOptions struct {
	Name           string
	Initial        interface{}
	DisplayName    string
	Validators     *ValidatorRegistry
	Choices        *FieldChoiceRegistry
	HelpText       string
	WidgetType     string
	ReadOnly       bool
	Required       bool
	WidgetPopulate func(renderContext *FormRenderContext, currentField *Field) interface{}
	IsFk           bool
}

func (ffo *FieldFormOptions) GetName() string {
	return ffo.Name
}

func (ffo *FieldFormOptions) IsItFk() bool {
	return ffo.IsFk
}

func (ffo *FieldFormOptions) GetWidgetPopulate() func(renderContext *FormRenderContext, currentField *Field) interface{} {
	return ffo.WidgetPopulate
}

func (ffo *FieldFormOptions) GetInitial() interface{} {
	return ffo.Initial
}

func (ffo *FieldFormOptions) GetDisplayName() string {
	return ffo.DisplayName
}

func (ffo *FieldFormOptions) GetValidators() *ValidatorRegistry {
	if ffo.Validators == nil {
		return NewValidatorRegistry()
	}
	return ffo.Validators
}

func (ffo *FieldFormOptions) GetChoices() *FieldChoiceRegistry {
	return ffo.Choices
}

func (ffo *FieldFormOptions) GetHelpText() string {
	return ffo.HelpText
}

func (ffo *FieldFormOptions) GetWidgetType() string {
	return ffo.WidgetType
}

func (ffo *FieldFormOptions) GetReadOnly() bool {
	return ffo.ReadOnly
}

func (ffo *FieldFormOptions) GetIsRequired() bool {
	return ffo.Required
}

func NewFieldFromGormField(gormField *schema.Field, forcedWidgetType string) *Field {
	var widget IWidget
	if gormField.PrimaryKey {
		widget = GetWidgetByWidgetType("hidden")
	} else if forcedWidgetType != "" {
		widget = GetWidgetByWidgetType(forcedWidgetType)
	} else {
		uadminFieldType := GetUadminFieldTypeFromGormField(gormField)
		widget = GetWidgetFromUadminFieldTypeAndGormField(uadminFieldType, gormField)
	}
	widget.InitializeAttrs()
	widget.SetName(gormField.Name)
	if gormField.NotNull && !gormField.HasDefaultValue {
		widget.SetRequired()
	}
	if gormField.Unique {
		widget.SetRequired()
	}
	if !gormField.PrimaryKey {
		widget.SetValue(gormField.DefaultValueInterface)
	}
	field := &Field{
		Field:           *gormField,
		UadminFieldType: GetUadminFieldTypeFromGormField(gormField),
		FieldConfig:     &FieldConfig{Widget: widget},
		Required:        gormField.NotNull && !gormField.HasDefaultValue,
		DisplayName:     gormField.Name,
	}
	return field
}

func NewUadminFieldFromGormField(gormModelV reflect.Value, field *schema.Field, r ITemplateRenderer, renderForAdmin bool) *Field {
	uadminformtag := field.Tag.Get("uadminform")
	var fieldOptions IFieldFormOptions
	var uadminField *Field
	if uadminformtag != "" {
		fieldOptions = UadminFormCongirurableOptionInstance.GetFieldFormOptions(uadminformtag)
		uadminField = NewFieldFromGormField(field, fieldOptions.GetWidgetType())
	} else {
		if field.PrimaryKey {
			fieldOptions = UadminFormCongirurableOptionInstance.GetFieldFormOptions("ReadonlyField")
			uadminField = NewFieldFromGormField(field, fieldOptions.GetWidgetType())
		} else {
			uadminField = NewFieldFromGormField(field, "")
		}
	}
	uadminField.DisplayName = field.Name
	if renderForAdmin {
		uadminField.FieldConfig.Widget.RenderForAdmin()
	}
	if fieldOptions != nil {
		uadminField.Initial = fieldOptions.GetInitial()
		if fieldOptions.GetDisplayName() != "" {
			uadminField.DisplayName = fieldOptions.GetDisplayName()
		}
		if fieldOptions.GetWidgetPopulate() != nil {
			uadminField.FieldConfig.Widget.SetPopulate(fieldOptions.GetWidgetPopulate())
		}
		uadminField.Validators = fieldOptions.GetValidators()
		uadminField.Choices = fieldOptions.GetChoices()
		uadminField.HelpText = fieldOptions.GetHelpText()
		uadminField.WidgetType = fieldOptions.GetWidgetType()
		uadminField.ReadOnly = fieldOptions.GetReadOnly()
		uadminField.FieldConfig.Widget.SetReadonly(uadminField.ReadOnly)
		if fieldOptions.GetIsRequired() {
			uadminField.FieldConfig.Widget.SetRequired()
		}
		if fieldOptions.GetHelpText() != "" {
			uadminField.FieldConfig.Widget.SetHelpText(fieldOptions.GetHelpText())
		}
	}
	uadminField.FieldConfig.Widget.RenderUsingRenderer(r)
	uadminField.FieldConfig.Widget.SetFieldDisplayName(field.Name)
	isTruthyValue := IsTruthyValue(gormModelV.FieldByName(field.Name).Interface())
	if isTruthyValue {
		uadminField.FieldConfig.Widget.SetValue(gormModelV.FieldByName(field.Name).Interface())
	}
	return uadminField
}

type FieldRegistry struct {
	Fields      map[string]*Field
	MaxOrdering int
}

func (fr *FieldRegistry) GetByName(name string) (*Field, error) {
	f, ok := fr.Fields[name]
	if !ok {
		return nil, fmt.Errorf("no field %s found", name)
	}
	return f, nil
}

func (fr *FieldRegistry) GetAllFields() map[string]*Field {
	return fr.Fields
}

func (fr *FieldRegistry) GetAllFieldsWithOrdering() []*Field {
	allFields := make([]*Field, 0)
	for _, field := range fr.Fields {
		allFields = append(allFields, field)
	}
	sort.Slice(allFields, func(i int, j int) bool {
		if allFields[i].Ordering == allFields[j].Ordering {
			return allFields[i].Name < allFields[j].Name
		}
		return allFields[i].Ordering < allFields[j].Ordering
	})
	return allFields
}

func (fr *FieldRegistry) GetPrimaryKey() (*Field, error) {
	for _, field := range fr.Fields {
		if field.PrimaryKey {
			return field, nil
		}
	}
	return nil, fmt.Errorf("no primary key found for model")
}

func (fr *FieldRegistry) AddField(field *Field) {
	if _, err := fr.GetByName(field.Name); err == nil {
		panic(fmt.Errorf("field %s already in the field registry", field.Name))
	}
	fr.Fields[field.Name] = field
	ordering := fr.MaxOrdering + 1
	field.Ordering = ordering
	fr.MaxOrdering = ordering
}

func NewFieldRegistry() *FieldRegistry {
	return &FieldRegistry{Fields: make(map[string]*Field)}
}
