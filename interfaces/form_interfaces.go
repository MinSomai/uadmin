package interfaces

import (
	"gorm.io/gorm/schema"
	"mime/multipart"
	"text/template"
)

type IFieldRegistry interface {
	GetByName(name string) (*Field, error)
	AddField(field *Field)
	GetAllFields() map[string]*Field
}

type WidgetType string

const TextInputWidgetType WidgetType = "text"
const NumberInputWidgetType WidgetType = "number"
const EmailInputWidgetType WidgetType = "email"
const URLInputWidgetType WidgetType = "url"
const PasswordInputWidgetType WidgetType = "password"
const HiddenInputWidgetType WidgetType = "hidden"
const DateInputWidgetType WidgetType = "date"
const DateTimeInputWidgetType WidgetType = "datetime"
const TimeInputWidgetType WidgetType = "time"
const TextareaInputWidgetType WidgetType = "textarea"
const CheckboxInputWidgetType WidgetType = "checkbox"
const SelectWidgetType WidgetType = "select"
const NullBooleanWidgetType WidgetType = "nullboolean"
const SelectMultipleWidgetType WidgetType = "selectmultiple"
const RadioSelectWidgetType WidgetType = "radioselect"
const RadioWidgetType WidgetType = "radio"
const CheckboxSelectMultipleWidgetType WidgetType = "checkboxselectmultiple"
const FileInputWidgetType WidgetType = "file"
const ClearableFileInputWidgetType WidgetType = "clearablefileinput"
const MultipleHiddenInputWidgetType WidgetType = "multiplehiddeninput"
const SplitDateTimeWidgetType WidgetType = "splitdatetime"
const SplitHiddenDateTimeWidgetType WidgetType = "splithiddendatetime"
const SelectDateWidgetType WidgetType = "selectdate"

type WidgetData map[string]interface{}
type IWidget interface {
	IdForLabel(model interface{}, F *Field) string
	GetName(model interface{}, F *Field) string
	GetWidgetType() WidgetType
	GetAttrs() map[string]string
	GetTemplateName() string
	RenderUsingRenderer(renderer ITemplateRenderer)
	// GetValue(v interface{}, model interface{}) interface{}
	Render() string
	SetValue(v interface{})
	SetName(name string)
	GetDataForRendering() WidgetData
	SetAttr(attrName string, value string)
	SetBaseFuncMap(baseFuncMap template.FuncMap)
	InitializeAttrs()
	SetFieldDisplayName (displayName string)
	SetReadonly(readonly bool)
	GetValue() interface{}
	ProceedForm(form *multipart.Form) error
	SetRequired()
	SetOutputValue(v interface{})
	GetOutputValue() interface{}
	SetErrors(validationErrors ValidationError)
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
const IpAddressUadminFieldType UadminFieldType = "ipaddress"
const GenericIpAddressUadminFieldType UadminFieldType = "genericipaddress"
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
	Widget IWidget
	AutocompleteURL string
	DependsOnAnotherFields []string
}

type Field struct {
	schema.Field
	ReadOnly bool
	UadminFieldType UadminFieldType
	FieldConfig *FieldConfig
	Required bool
	DisplayName string
	HelpText string
	Choices *FieldChoiceRegistry
	Validators []IValidator
	SortingDisabled bool
	Populate func(field *Field, m interface{}) interface{}
	Initial interface{}
	WidgetType string
}

func (f *Field) ProceedForm(form *multipart.Form) ValidationError {
	err := f.FieldConfig.Widget.ProceedForm(form)
	if err == nil {
		validationErrors := make(ValidationError, 0)
		for _, validator := range f.Validators {
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