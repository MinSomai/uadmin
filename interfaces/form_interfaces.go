package interfaces

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"mime/multipart"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type IFieldRegistry interface {
	GetByName(name string) (*Field, error)
	AddField(field *Field)
	GetAllFields() map[string]*Field
	GetPrimaryKey() (*Field, error)
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
	SetShowOnlyHtmlInput()
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

type FieldRegistry struct {
	IFieldRegistry
	Fields map[string]*Field
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
}

func NewFieldRegistry() *FieldRegistry {
	return &FieldRegistry{Fields: make(map[string]*Field)}
}

type FieldFormOptions struct {
	IFieldFormOptions
	Name string
	Initial interface{}
	DisplayName string
	Validators []IValidator
	Choices *FieldChoiceRegistry
	HelpText string
	WidgetType string
	ReadOnly bool
}

func (ffo *FieldFormOptions) GetName() string {
	return ffo.Name
}

func (ffo *FieldFormOptions) GetInitial() interface{} {
	return ffo.Initial
}

func (ffo *FieldFormOptions) GetDisplayName() string {
	return ffo.DisplayName
}

func (ffo *FieldFormOptions) GetValidators() []IValidator {
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

type ColumnSchema struct {
	ShowLabel bool
	Fields []*Field
}

type FormRow struct {
	Columns []*ColumnSchema
}

type IGrouppedFieldsRegistry interface {
	AddGroup(grouppedFields *GrouppedFields)
	GetGroupByName(name string) *GrouppedFields
}

type GrouppedFieldsRegistry struct {
	IGrouppedFieldsRegistry
	GrouppedFields map[string]*GrouppedFields
}

func (tfr *GrouppedFieldsRegistry) GetGroupByName(name string) (*GrouppedFields, error) {
	gf, ok := tfr.GrouppedFields[name]
	if !ok {
		return nil, fmt.Errorf("No field %s found", name)
	}
	return gf, nil
}

func (tfr *GrouppedFieldsRegistry) AddGroup(grouppedFields *GrouppedFields) {
	if _, err := tfr.GetGroupByName(grouppedFields.Name); err != nil {
		panic(err)
	}
	tfr.GrouppedFields[grouppedFields.Name] = grouppedFields
}

type GrouppedFields struct {
	Rows []*FormRow
	ExtraCssClasses []string
	Description string
	Name string
}

type StaticFiles struct {
	ExtraCSS []string
	ExtraJS []string
}

type FormError struct{
	FieldError map[string]ValidationError
	GeneralErrors ValidationError
}

func (fe *FormError) AddGeneralError(err error) {
	fe.GeneralErrors = append(fe.GeneralErrors, err)
}
func (fe *FormError) IsEmpty() bool {
	return len(fe.FieldError) == 0 && len(fe.GeneralErrors) == 0
}

func (e *FormError) Error() string {
	return "Form validation not successful"
}

func GetWidgetByWidgetType(widgetType string) IWidget {
	var widget IWidget
	switch(widgetType) {
	case "image":
		widget = &FileWidget{}
		widget.SetAttr("accept", "image/*")
	case "hidden":
		widget = &HiddenWidget{}
	}
	return widget
}

type Widget struct {
	IWidget
	Attrs map[string]string
	TemplateName string
	Renderer ITemplateRenderer
	Value interface{}
	Name string
	FieldDisplayName string
	BaseFuncMap template.FuncMap
	ReadOnly bool
	ShowOnlyHtmlInput bool
	Required bool
	OutputValue interface{}
	ValidationErrors ValidationError
}

func (w *Widget) SetRequired() {
	w.Required = true
}

func (w *Widget) SetShowOnlyHtmlInput() {
	w.ShowOnlyHtmlInput = true
}

func (w *Widget) SetOutputValue(v interface{}) {
	w.OutputValue = v
}

func (w *Widget) GetOutputValue() interface{} {
	return w.OutputValue
}

func (w *Widget) SetErrors(validationErrors ValidationError) {
	w.ValidationErrors = validationErrors
}

func (w *Widget) InitializeAttrs() {
	w.Attrs = make(map[string]string)
}

func (w *Widget) SetBaseFuncMap(baseFuncMap template.FuncMap) {
	w.BaseFuncMap = baseFuncMap
}

func (w *Widget) IdForLabel(model interface{}, F *Field) string {
	return ""
}

func (w *Widget) SetFieldDisplayName(fieldDisplayName string) {
	w.FieldDisplayName = fieldDisplayName
}

func (w *Widget) SetReadonly(readonly bool) {
	w.ReadOnly = readonly
}

func (w *Widget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

func (w *Widget) RenderUsingRenderer(r ITemplateRenderer) {
	w.Renderer = r
}

func (w *Widget) SetAttr(attrName string, value string) {
	if w.Attrs == nil {
		w.InitializeAttrs()
	}
	w.Attrs[attrName] = value
}

func (w *Widget) GetName(model interface{}, F *Field) string {
	return w.Name
}

func (w *Widget) SetName(name string) {
	w.Name = name
}

func (w *Widget) GetAttrs() map[string]string {
	if w.Attrs != nil {
		return w.Attrs
	}
	return make(map[string]string)
}

func (w *Widget) SetValue(v interface{}) {
	w.Value = v
}

func (w *Widget) GetValue() interface{} {
	return w.Value
}

func (w *Widget) Render() string {
	data := w.GetDataForRendering()
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *Widget) GetDataForRendering() WidgetData {
	value := TransformValueForWidget(w.Value)
	return map[string]interface{}{
		"Attrs": w.GetAttrs(), "Value": template.HTMLEscapeString(value.(string)),
		"Name": w.Name, "FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func RenderWidget(renderer ITemplateRenderer, templateName string, data map[string]interface{}, baseFuncMap template.FuncMap) string {
	if renderer == nil {
		r := NewTemplateRenderer("")
		return r.RenderAsString(CurrentConfig.TemplatesFS, templateName, data, baseFuncMap)
	} else {
		return renderer.RenderAsString(
			CurrentConfig.TemplatesFS, templateName,
			data, baseFuncMap,
		)
	}
}

type TextWidget struct {
	Widget
}

func (tw *TextWidget) GetWidgetType() WidgetType {
	return TextInputWidgetType
}

func (tw *TextWidget) GetTemplateName() string {
	if tw.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/text")
	}
	return CurrentConfig.GetPathToTemplate(tw.TemplateName)
}

func (tw *TextWidget) Render() string {
	data := tw.Widget.GetDataForRendering()
	data["Type"] = tw.GetWidgetType()
	data["ShowOnlyHtmlInput"] = tw.ShowOnlyHtmlInput
	return RenderWidget(tw.Renderer, tw.GetTemplateName(), data, tw.BaseFuncMap) // tw.Value, tw.Widget.GetAttrs()
}

type NumberWidget struct {
	Widget
}

func (w *NumberWidget) GetWidgetType() WidgetType {
	return NumberInputWidgetType
}

func (w *NumberWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/number")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *NumberWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *NumberWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	if !govalidator.IsInt(v[0]) {
		return fmt.Errorf("should be a number")
	}
	w.SetOutputValue(v[0])
	return nil
}

type EmailWidget struct {
	Widget
}

func (w *EmailWidget) GetWidgetType() WidgetType {
	return EmailInputWidgetType
}

func (w *EmailWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/email")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *EmailWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *EmailWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	if !govalidator.IsEmail(v[0]) {
		return fmt.Errorf("should be an email")
	}
	w.SetOutputValue(v[0])
	return nil
}

type URLWidget struct {
	Widget
	UrlValid bool
	CurrentLabel string
	Href string
	Value string
	ChangeLabel string
	AppendHttpsAutomatically bool
}

func (w *URLWidget) GetWidgetType() WidgetType {
	return URLInputWidgetType
}

func (w *URLWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/url")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *URLWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	data["UrlValid"] = w.UrlValid
	if w.CurrentLabel == "" {
		data["CurrentLabel"] = "Url"
	} else {
		data["CurrentLabel"] = w.CurrentLabel
	}
	data["Href"] = w.Href
	data["Value"] = w.Widget.Value
	if w.ChangeLabel == "" {
		data["ChangeLabel"] = "Change"
	} else {
		data["ChangeLabel"] = w.ChangeLabel
	}
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *URLWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	url := v[0]
	if w.AppendHttpsAutomatically {
		urlInitialRegex := regexp.MustCompile(`^http(s)?://.*`)
		if !urlInitialRegex.Match([]byte(v[0])) {
			url = "https://" + url
		}
	}
	if !govalidator.IsURL(url) {
		return fmt.Errorf("should be an url")
	}
	w.SetOutputValue(v[0])
	return nil
}

type PasswordWidget struct {
	Widget
}

func (w *PasswordWidget) GetWidgetType() WidgetType {
	return PasswordInputWidgetType
}

func (w *PasswordWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/password")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *PasswordWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	data["Value"] = ""
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *PasswordWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	if len(v[0]) < CurrentConfig.D.Auth.MinPasswordLength {
		return fmt.Errorf("length of the password has to be at least %d symbols", CurrentConfig.D.Auth.MinPasswordLength)
	}
	w.SetOutputValue(v[0])
	return nil
}


type HiddenWidget struct {
	Widget
}

func (w *HiddenWidget) GetWidgetType() WidgetType {
	return HiddenInputWidgetType
}

func (w *HiddenWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/hidden")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *HiddenWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *HiddenWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type DateWidget struct {
	Widget
	DateValue string
}

func (w *DateWidget) GetWidgetType() WidgetType {
	return DateInputWidgetType
}

func (w *DateWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/date")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *DateWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	if w.DateValue != "" {
		data["Value"] = w.DateValue
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *DateWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.DateValue = v[0]
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(CurrentConfig.D.Uadmin.DateFormat, v[0])
	if err != nil {
		return err
	}
	w.SetOutputValue(&d)
	return nil
}

type DateTimeWidget struct {
	Widget
	DateTimeValue string
}

func (w *DateTimeWidget) GetWidgetType() WidgetType {
	return DateTimeInputWidgetType
}

func (w *DateTimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/datetime")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *DateTimeWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	if w.DateTimeValue != "" {
		data["Value"] = w.DateTimeValue
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *DateTimeWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.DateTimeValue = v[0]
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(CurrentConfig.D.Uadmin.DateTimeFormat, v[0])
	if err != nil {
		return err
	}
	w.SetOutputValue(&d)
	return nil
}

type TimeWidget struct {
	Widget
	TimeValue string
}

func (w *TimeWidget) GetWidgetType() WidgetType {
	return TimeInputWidgetType
}

func (w *TimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/time")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *TimeWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	if w.TimeValue != "" {
		data["Value"] = w.TimeValue
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *TimeWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.TimeValue = v[0]
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(CurrentConfig.D.Uadmin.TimeFormat, v[0])
	if err != nil {
		return err
	}
	w.SetOutputValue(&d)
	return nil
}

type TextareaWidget struct {
	Widget
}

func (w *TextareaWidget) GetWidgetType() WidgetType {
	return TextareaInputWidgetType
}

func (w *TextareaWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/textarea")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *TextareaWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *TextareaWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v[0])
	if w.Required && v[0] == "" {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type CheckboxWidget struct {
	Widget
}

func (w *CheckboxWidget) GetWidgetType() WidgetType {
	return CheckboxInputWidgetType
}

func (w *CheckboxWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/checkbox")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *CheckboxWidget) Render() string {
	value := TransformValueForWidget(w.Value)
	if value != "" {
		w.Attrs["checked"] = "checked"
	}
	w.Value = nil
	data := w.Widget.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *CheckboxWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	_, ok := form.Value[w.Name]
	w.SetValue(ok == true)
	w.SetOutputValue(ok == true)
	return nil
}

type SelectOptGroup struct {
	OptLabel string
	Value interface{}
	Selected bool
}

type SelectOptGroupStringified struct {
	OptLabel string
	Value string
	Selected bool
	OptionTemplateName string
}

type SelectWidget struct {
	Widget
	OptGroups map[string][]*SelectOptGroup
}

func (w *SelectWidget) GetWidgetType() WidgetType {
	return SelectWidgetType
}

func (w *SelectWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/select")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SelectWidget) GetDataForRendering() WidgetData {
	value := TransformValueForWidget(w.Value)
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
				OptLabel: optGroup.OptLabel,
				Value: value1,
				Selected: value1 == value,
				OptionTemplateName: "widgets/select.option",
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name, "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *SelectWidget) Render() string {
	data := w.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SelectWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v[0])
	if foundNotExistent {
		return fmt.Errorf("value %s is not valid for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type NullBooleanWidget struct {
	Widget
	OptGroups map[string][]*SelectOptGroup
}

func (w *NullBooleanWidget) GetWidgetType() WidgetType {
	return NullBooleanWidgetType
}

func (w *NullBooleanWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/nullboolean")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *NullBooleanWidget) GetOptGroups() map[string][]*SelectOptGroup {
	if w.OptGroups == nil {
		defaultOptGroups := make(map[string][]*SelectOptGroup)
		defaultOptGroups[""] = make([]*SelectOptGroup, 0)
		defaultOptGroups[""] = append(defaultOptGroups[""], &SelectOptGroup{
			OptLabel: "Yes",
			Value: "yes",
		})
		defaultOptGroups[""] = append(defaultOptGroups[""], &SelectOptGroup{
			OptLabel: "No",
			Value: "no",
		})
		return defaultOptGroups
	}
	return w.OptGroups
}

func (w *NullBooleanWidget) GetDataForRendering() WidgetData {
	value := TransformValueForWidget(w.Value)
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	for optGroupName, optGroups := range w.GetOptGroups() {
		optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
				OptLabel: optGroup.OptLabel,
				Value: value1,
				Selected: value1 == value,
				OptionTemplateName: "widgets/select.option",
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name, "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *NullBooleanWidget) Render() string {
	data := w.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *NullBooleanWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v[0])
	if foundNotExistent {
		return fmt.Errorf("value %s is not valid for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type SelectMultipleWidget struct {
	Widget
	OptGroups map[string][]*SelectOptGroup
}

func (w *SelectMultipleWidget) GetWidgetType() WidgetType {
	return SelectMultipleWidgetType
}

func (w *SelectMultipleWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/select")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SelectMultipleWidget) GetDataForRendering() WidgetData {
	w.Attrs["multiple"] = "true"
	value := TransformValueForWidget(w.Value).([]string)
	optGroupSstringified := make(map[string][]*SelectOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*SelectOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &SelectOptGroupStringified{
				OptLabel: optGroup.OptLabel,
				Value: value1,
				Selected: Contains(value, value1),
				OptionTemplateName: "widgets/select.option",
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name, "OptGroups": optGroupSstringified,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *SelectMultipleWidget) Render() string {
	data := w.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SelectMultipleWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v)
	if foundNotExistent {
		return fmt.Errorf("value %s is not valid for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v)
	return nil
}

type RadioOptGroup struct {
	OptLabel string
	Value interface{}
	Selected bool
	Label string
}

type RadioOptGroupStringified struct {
	OptLabel string
	Value string
	Selected bool
	OptionTemplateName string
	WrapLabel bool
	ForId string
	Label string
	Type string
	Name string
	Attrs map[string]string
	FieldDisplayName string
	ReadOnly bool
}

type RadioSelectWidget struct {
	Widget
	OptGroups map[string][]*RadioOptGroup
	Id string
	WrapLabel bool
}

func (w *RadioSelectWidget) GetWidgetType() WidgetType {
	return RadioWidgetType
}

func (w *RadioSelectWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/radioselect")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *RadioSelectWidget) SetId(id string) {
	w.Id = id
}

func (w *RadioSelectWidget) GetDataForRendering() WidgetData {
	value := TransformValueForWidget(w.Value).(string)
	optGroupSstringified := make(map[string][]*RadioOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*RadioOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &RadioOptGroupStringified{
				OptLabel: optGroup.OptLabel,
				Value: value1,
				Selected: value == value1,
				OptionTemplateName: "widgets/radio.option",
				Label: optGroup.Label,
				WrapLabel: w.WrapLabel,
				ForId: w.Id,
				Type: "radio",
				Name: w.Name,
				Attrs: w.Widget.GetAttrs(),
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name, "OptGroups": optGroupSstringified, "Id": w.Id,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *RadioSelectWidget) Render() string {
	data := w.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *RadioSelectWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v[0])
	if foundNotExistent {
		return fmt.Errorf("value %s is not valid for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v[0])
	return nil
}

type CheckboxSelectMultipleWidget struct {
	Widget
	OptGroups map[string][]*RadioOptGroup
	Id string
	WrapLabel bool
}

func (w *CheckboxSelectMultipleWidget) GetWidgetType() WidgetType {
	return CheckboxSelectMultipleWidgetType
}

func (w *CheckboxSelectMultipleWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/checkboxselectmultiple")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *CheckboxSelectMultipleWidget) SetId(id string) {
	w.Id = id
}

func (w *CheckboxSelectMultipleWidget) GetDataForRendering() WidgetData {
	value := TransformValueForWidget(w.Value).([]string)
	optGroupSstringified := make(map[string][]*RadioOptGroupStringified)
	for optGroupName, optGroups := range w.OptGroups {
		optGroupSstringified[optGroupName] = make([]*RadioOptGroupStringified, 0)
		for _, optGroup := range optGroups {
			value1 := TransformValueForWidget(optGroup.Value).(string)
			optGroupSstringified[optGroupName] = append(optGroupSstringified[optGroupName], &RadioOptGroupStringified{
				OptLabel: optGroup.OptLabel,
				Value: value1,
				Selected: Contains(value, value1),
				OptionTemplateName: "widgets/checkbox.option",
				Label: optGroup.Label,
				WrapLabel: w.WrapLabel,
				ForId: w.Id,
				Type: "checkbox",
				Name: w.Name,
				Attrs: w.Widget.GetAttrs(),
			})
		}
	}
	return map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name, "OptGroups": optGroupSstringified, "Id": w.Id,
		"FieldDisplayName": w.FieldDisplayName, "ReadOnly": w.ReadOnly,
	}
}

func (w *CheckboxSelectMultipleWidget) Render() string {
	data := w.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *CheckboxSelectMultipleWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	foundNotExistent := false
	optValues := []string{}
	for _, optGroup := range w.OptGroups {
		for _, optGroupOption := range optGroup {
			optValues = append(optValues, optGroupOption.Value.(string))
		}
	}
	var notExistentValue string
	for _, v1 := range v {
		if !Contains(optValues, v1) {
			foundNotExistent = true
			notExistentValue = v1
			break
		}
	}
	w.SetValue(v)
	if foundNotExistent {
		return fmt.Errorf("value %s is not valid for the field %s", notExistentValue, w.FieldDisplayName)
	}
	w.SetOutputValue(v)
	return nil
}

type FileWidget struct {
	Widget
	Storage IStorageInterface
	UploadPath string
	Multiple bool
}

func (w *FileWidget) GetWidgetType() WidgetType {
	return FileInputWidgetType
}

func (w *FileWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/file")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *FileWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["Value"] = ""
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *FileWidget) ProceedForm(form *multipart.Form) error {
	files := form.File[w.Name]
	storage := w.Storage
	if storage == nil {
		storage = NewFsStorage()
	}
	ret  := make([]string, 0)
	var filename string
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return err
		}
		var bytecontent []byte
		_, err = f.Read(bytecontent)
		if err != nil {
			return err
		}
		filename, err = storage.Save(&FileForStorage{
			Content: bytecontent,
			PatternForTheFile: "*." + strings.Split(file.Filename, ".")[1],
			Filename: file.Filename,
		})
		if err != nil {
			return err
		}
		ret = append(ret, filename)
	}
	if w.Multiple {
		w.SetOutputValue(ret)
	} else if len(ret)  > 0 {
		w.SetOutputValue(ret[0])
	} else {
		w.SetOutputValue("")
	}
	return nil
}

type URLValue struct {
	URL string
}

type ClearableFileWidget struct {
	Widget
	InitialText string
	CurrentValue *URLValue
	Required bool
	Id string
	ClearCheckboxLabel string
	InputText string
	Storage IStorageInterface
	UploadPath string
	Multiple bool
}

func (w *ClearableFileWidget) GetWidgetType() WidgetType {
	return FileInputWidgetType
}

func (w *ClearableFileWidget) SetId(id string) {
	w.Id = id
}

func (w *ClearableFileWidget) IsInitial() bool {
	return w.CurrentValue == nil
}

func (w *ClearableFileWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/clearablefile")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *ClearableFileWidget) Render() string {
	data := w.Widget.GetDataForRendering()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	data["IsInitial"] = w.IsInitial()
	data["InitialText"] = w.InitialText
	data["CurrentValue"] = w.CurrentValue
	data["Required"] = w.Required
	data["Id"] = w.Id
	data["ClearCheckboxLabel"] = w.ClearCheckboxLabel
	data["InputText"] = w.InputText
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *ClearableFileWidget) ProceedForm(form *multipart.Form) error {
	files := form.File[w.Name]
	storage := w.Storage
	if storage == nil {
		storage = NewFsStorage()
	}
	ret  := make([]string, 0)
	var err error
	var filename string
	for _, file := range files {
		f, _ := file.Open()
		var bytecontent []byte
		_, err = f.Read(bytecontent)
		filename, err = storage.Save(&FileForStorage{
			Content: bytecontent,
			PatternForTheFile: "*." + strings.Split(file.Filename, ".")[1],
			Filename: file.Filename,
		})
		if err != nil {
			return err
		}
		ret = append(ret, filename)
	}
	if w.Multiple {
		w.SetOutputValue(ret)
	} else if len(ret)  > 0 {
		w.SetOutputValue(ret[0])
	} else {
		w.SetOutputValue("")
		w.SetValue("")
	}
	return nil
}

type MultipleInputHiddenWidget struct {
	Widget
}

func (w *MultipleInputHiddenWidget) GetWidgetType() WidgetType {
	return MultipleHiddenInputWidgetType
}

func (w *MultipleInputHiddenWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/multipleinputhidden")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *MultipleInputHiddenWidget) Render() string {
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name,
	}
	data["Type"] = w.GetWidgetType()
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	subwidgets := make([]WidgetData, 0)
	value := TransformValueForWidget(w.Value).([]string)
	for _, v := range value {
		w1 := HiddenWidget{}
		w1.Name = w.Name
		w1.SetValue(v)
		w1.Attrs = make(map[string]string)
		for attrName, attrValue := range w.Attrs {
			w1.Attrs[attrName] = attrValue
		}
		vd := w1.GetDataForRendering()
		vd["Type"] = w1.GetWidgetType()
		vd["TemplateName"] = "widgets/hidden"
		subwidgets = append(subwidgets, vd)
	}
	data["Subwidgets"] = subwidgets
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *MultipleInputHiddenWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	v, ok := form.Value[w.Name]
	if !ok {
		return fmt.Errorf("no field with name %s has been submitted", w.FieldDisplayName)
	}
	w.SetValue(v)
	w.SetOutputValue(v)
	return nil
}

type SplitDateTimeWidget struct {
	Widget
	DateAttrs map[string]string
	TimeAttrs map[string]string
	DateFormat string
	TimeFormat string
	DateLabel string
	TimeLabel string
	DateValue string
	TimeValue string
}

func (w *SplitDateTimeWidget) GetWidgetType() WidgetType {
	return SplitDateTimeWidgetType
}

func (w *SplitDateTimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/splitdatetime")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SplitDateTimeWidget) Render() string {
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name,
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	if w.DateLabel == "" {
		data["DateLabel"] = "Date:"
	} else {
		data["DateLabel"] = w.DateLabel
	}
	if w.TimeLabel == "" {
		data["TimeLabel"] = "Time:"
	} else {
		data["TimeLabel"] = w.TimeLabel
	}
	subwidgets := make([]WidgetData, 0)
	value := TransformValueForWidget(w.Value).(*time.Time)
	w1 := DateWidget{}
	w1.Name = w.Name + "_date"
	if w.DateValue != "" {
		w1.SetValue(w.DateValue)
	} else {
		w1.SetValue(value.Format(w.DateFormat))
	}
	w1.Attrs = w.DateAttrs
	vd := w1.Widget.GetDataForRendering()
	vd["Type"] = w1.GetWidgetType()
	vd["TemplateName"] = "widgets/date"
	subwidgets = append(subwidgets, vd)
	w2 := TimeWidget{}
	w2.Name = w.Name + "_time"
	if w.TimeValue != "" {
		w2.SetValue(w.TimeValue)
	} else {
		w2.SetValue(value.Format(w.TimeFormat))
	}
	w2.Attrs = w.TimeAttrs
	vd1 := w2.Widget.GetDataForRendering()
	vd1["Type"] = w2.GetWidgetType()
	vd1["TemplateName"] = "widgets/time"
	subwidgets = append(subwidgets, vd1)
	data["Subwidgets"] = subwidgets
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SplitDateTimeWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	vDate, ok := form.Value[w.Name + "_date"]
	if !ok {
		return fmt.Errorf("no date has been submitted for field %s", w.FieldDisplayName)
	}
	w.DateValue = vDate[0]
	vTime, ok := form.Value[w.Name + "_time"]
	if !ok {
		return fmt.Errorf("no time has been submitted for field %s", w.FieldDisplayName)
	}
	w.TimeValue = vTime[0]
	if w.Required && (vDate[0] == "" || vTime[0] == "") {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(w.DateFormat, vDate[0])
	if err != nil {
		return err
	}
	t, err := time.Parse(w.TimeFormat, vTime[0])
	if err != nil {
		return err
	}
	newT := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	w.SetOutputValue(&newT)
	return nil
}

type SplitHiddenDateTimeWidget struct {
	Widget
	DateAttrs map[string]string
	TimeAttrs map[string]string
	DateFormat string
	TimeFormat string
	DateValue string
	TimeValue string
}

func (w *SplitHiddenDateTimeWidget) GetWidgetType() WidgetType {
	return SplitHiddenDateTimeWidgetType
}

func (w *SplitHiddenDateTimeWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/splithiddendatetime")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SplitHiddenDateTimeWidget) Render() string {
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name,
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	subwidgets := make([]WidgetData, 0)
	value := TransformValueForWidget(w.Value).(*time.Time)
	w1 := DateWidget{}
	w1.Name = w.Name + "_date"
	if w.DateValue != "" {
		w1.SetValue(w.DateValue)
	} else {
		w1.SetValue(value.Format(w.DateFormat))
	}
	w1.Attrs = w.DateAttrs
	vd := w1.Widget.GetDataForRendering()
	vd["Type"] = "hidden"
	vd["TemplateName"] = "widgets/date"
	subwidgets = append(subwidgets, vd)
	w2 := TimeWidget{}
	w2.Name = w.Name + "_time"
	if w.TimeValue != "" {
		w2.SetValue(w.TimeValue)
	} else {
		w2.SetValue(value.Format(w.TimeFormat))
	}
	w2.Attrs = w.TimeAttrs
	vd1 := w2.Widget.GetDataForRendering()
	vd1["Type"] = "hidden"
	vd1["TemplateName"] = "widgets/time"
	subwidgets = append(subwidgets, vd1)
	data["Subwidgets"] = subwidgets
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SplitHiddenDateTimeWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	vDate, ok := form.Value[w.Name + "_date"]
	if !ok {
		return fmt.Errorf("no date has been submitted for field %s", w.FieldDisplayName)
	}
	w.DateValue = vDate[0]
	vTime, ok := form.Value[w.Name + "_time"]
	if !ok {
		return fmt.Errorf("no time has been submitted for field %s", w.FieldDisplayName)
	}
	w.TimeValue = vTime[0]
	if w.Required && (vDate[0] == "" || vTime[0] == "") {
		return fmt.Errorf("field %s is required", w.FieldDisplayName)
	}
	d, err := time.Parse(w.DateFormat, vDate[0])
	if err != nil {
		return err
	}
	t, err := time.Parse(w.TimeFormat, vTime[0])
	if err != nil {
		return err
	}
	newT := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	w.SetOutputValue(&newT)
	return nil
}

type SelectDateWidget struct {
	Widget
	Years []int
	Months []*SelectOptGroup
	EmptyLabel []*SelectOptGroup
	EmptyLabelString string
	IsRequired bool
	YearValue string
	MonthValue string
	DayValue string
}

func (w *SelectDateWidget) GetWidgetType() WidgetType {
	return SelectDateWidgetType
}

func (w *SelectDateWidget) GetTemplateName() string {
	if w.TemplateName == "" {
		return CurrentConfig.GetPathToTemplate("widgets/selectdate")
	}
	return CurrentConfig.GetPathToTemplate(w.TemplateName)
}

func (w *SelectDateWidget) Render() string {
	value := TransformValueForWidget(w.Value).(*time.Time)
	data := map[string]interface{}{
		"Attrs": w.GetAttrs(),
		"Name": w.Name,
	}
	data["ShowOnlyHtmlInput"] = w.ShowOnlyHtmlInput
	data["Type"] = w.GetWidgetType()
	dateParts := []string{}
	for _, formatChar := range CurrentConfig.D.Uadmin.DateFormatOrder {
		if formatChar == 'y' {
			if Contains(dateParts, "year") {
				continue
			}
			dateParts = append(dateParts, "year")
		} else if formatChar == 'd' {
			if Contains(dateParts, "day") {
				continue
			}
			dateParts = append(dateParts, "day")
		} else if formatChar == 'm' {
			if Contains(dateParts, "month") {
				continue
			}
			dateParts = append(dateParts, "month")
		}
	}
	if w.Years == nil {
		initialYear := time.Now().Year()
		w.Years = make([]int, 0)
		for i := initialYear; i <= initialYear + 10; i++ {
			w.Years = append(w.Years, i)
		}
	}
	var yearNoneValue *SelectOptGroup
	var monthNoneValue *SelectOptGroup
	var dayNoneValue *SelectOptGroup
	if w.EmptyLabel == nil {
		noneValue := &SelectOptGroup{
			OptLabel: w.EmptyLabelString,
			Value: "",
		}
		dayNoneValue = noneValue
		yearNoneValue = noneValue
		monthNoneValue = noneValue
	} else {
		if len(w.EmptyLabel) != 3 {
			panic("empty_label slice must have 3 elements.")
		}
		dayNoneValue = w.EmptyLabel[2]
		yearNoneValue = w.EmptyLabel[0]
		monthNoneValue = w.EmptyLabel[1]
	}
	if w.Months == nil {
		w.Months = MakeMonthsSelect()
		if !w.IsRequired {
			w.Months = append(w.Months, monthNoneValue)
			copy(w.Months[1:], w.Months)
			w.Months[0] = monthNoneValue
		}
	}
	var yearChoices []*SelectOptGroup
	if !w.IsRequired {
		yearChoices = append(yearChoices, yearNoneValue)
	}
	for _, year := range w.Years {
		yearChoices = append(yearChoices, &SelectOptGroup{
			OptLabel: strconv.Itoa(year),
			Value: strconv.Itoa(year),
		})
	}
	var dayChoices []*SelectOptGroup
	for i := 1; i < 32; i++ {
		dayChoices = append(dayChoices, &SelectOptGroup{
			OptLabel: strconv.Itoa(i),
			Value: strconv.Itoa(i),
		})
		if !w.IsRequired {
			dayChoices = append(dayChoices, dayNoneValue)
			copy(dayChoices[1:], dayChoices)
			dayChoices[0] = dayNoneValue
		}
	}
	subwidgets := make([]WidgetData, 0)
	w1 := SelectWidget{}
	w1.OptGroups = make(map[string][]*SelectOptGroup)
	w1.OptGroups[""] = yearChoices
	w1.Name = w.Name + "_year"
	if w.YearValue != "" {
		w1.SetValue(w.YearValue)
	} else {
		w1.SetValue(value.Year())
	}
	w1.Attrs = w.GetAttrs()
	vd := w1.GetDataForRendering()
	vd["Type"] = "select"
	vd["TemplateName"] = "widgets/select"
	yearWd := vd
	w2 := SelectWidget{}
	w2.OptGroups = make(map[string][]*SelectOptGroup)
	w2.OptGroups[""] = w.Months
	w2.Name = w.Name + "_month"
	if w.YearValue != "" {
		w2.SetValue(w.MonthValue)
	} else {
		w2.SetValue(value.Month())
	}
	w2.Attrs = w.GetAttrs()
	vd2 := w2.GetDataForRendering()
	vd2["Type"] = "select"
	vd2["TemplateName"] = "widgets/select"
	w3 := SelectWidget{}
	w3.OptGroups = make(map[string][]*SelectOptGroup)
	w3.OptGroups[""] = dayChoices
	w3.Name = w.Name + "_day"
	if w.DayValue != "" {
		w3.SetValue(w.DayValue)
	} else {
		w3.SetValue(value.Day())
	}
	w3.Attrs = w.GetAttrs()
	vd3 := w3.GetDataForRendering()
	vd3["Type"] = "select"
	vd3["TemplateName"] = "widgets/select"
	dayWd := vd3
	monthWd := vd2
	for _, datePart := range dateParts {
		if datePart == "year" {
			subwidgets = append(subwidgets, yearWd)
		} else if datePart == "month" {
			subwidgets = append(subwidgets, monthWd)
		} else if datePart == "day" {
			subwidgets = append(subwidgets, dayWd)
		}
	}
	data["Subwidgets"] = subwidgets
	return RenderWidget(w.Renderer, w.GetTemplateName(), data, w.BaseFuncMap)
}

func (w *SelectDateWidget) ProceedForm(form *multipart.Form) error {
	if w.ReadOnly {
		return nil
	}
	vYear, ok := form.Value[w.Name + "_year"]
	if !ok {
		return fmt.Errorf("no year has been submitted for field %s", w.FieldDisplayName)
	}
	w.YearValue = vYear[0]
	vMonth, ok := form.Value[w.Name + "_month"]
	if !ok {
		return fmt.Errorf("no month has been submitted for field %s", w.FieldDisplayName)
	}
	w.MonthValue = vMonth[0]
	vDay, ok := form.Value[w.Name + "_day"]
	if !ok {
		return fmt.Errorf("no month has been submitted for field %s", w.FieldDisplayName)
	}
	w.DayValue = vDay[0]
	if w.Required && (w.YearValue == "" || w.MonthValue == "" || w.DayValue == "") {
		return fmt.Errorf("either year, month, value is empty")
	}
	day, err := strconv.Atoi(w.DayValue)
	if err != nil {
		return fmt.Errorf("incorrect day")
	}
	month, err := strconv.Atoi(w.MonthValue)
	if err != nil {
		return fmt.Errorf("incorrect month")
	}
	year, err := strconv.Atoi(w.YearValue)
	if err != nil {
		return fmt.Errorf("incorrect year")
	}
	d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	w.SetOutputValue(&d)
	return nil
}

func MakeMonthsSelect() []*SelectOptGroup {
	ret := make([]*SelectOptGroup, 0)
	ret = append(ret, &SelectOptGroup{
		Value: "1",
		OptLabel: "January",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "2",
		OptLabel: "February",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "3",
		OptLabel: "March",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "4",
		OptLabel: "April",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "5",
		OptLabel: "May",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "6",
		OptLabel: "June",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "7",
		OptLabel: "July",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "8",
		OptLabel: "August",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "9",
		OptLabel: "September",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "10",
		OptLabel: "October",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "11",
		OptLabel: "November",
	})
	ret = append(ret, &SelectOptGroup{
		Value: "12",
		OptLabel: "December",
	})
	return ret
}

type Form struct {
	ExcludeFields IFieldRegistry
	FieldsToShow IFieldRegistry
	FieldRegistry IFieldRegistry
	GroupsOfTheFields *GrouppedFieldsRegistry
	TemplateName string
	FormTitle string
	Renderer ITemplateRenderer
	RequestContext map[string]interface{}
	ErrorMessage string
	ExtraStatic *StaticFiles `json:"-"`
}

func (f *Form) Render() string {
	RenderFieldGroups := func(funcs1 template.FuncMap) func () string {
		return func () string {
			templateWriter := bytes.NewBuffer([]byte{})
			ret := make([]string, 0)
			for _, group := range f.GroupsOfTheFields.GrouppedFields {
				for _, row := range group.Rows {
					data2 := row
					templateWriter.Reset()
					err := RenderHTMLAsString(templateWriter, CurrentConfig.TemplatesFS, CurrentConfig.GetPathToTemplate("form/grouprow"), data2, FuncMap, funcs1)
					if err != nil {
						Trail(CRITICAL, "Error while parsing include of the template %s", "form/grouprow")
						panic(err)
					}
					ret = append(ret, templateWriter.String())
				}
			}
			return strings.Join(ret, "\n")
		}
	}
	FieldValue := func (fieldName string) interface{} {
		field, _ := f.FieldRegistry.GetByName(fieldName)
		return field.FieldConfig.Widget.GetValue()
	}
	func1 := make(template.FuncMap)
	func1["RenderFieldGroups"] = RenderFieldGroups(FuncMap)
	func1["FormFieldValue"] = FieldValue
	templateName := CurrentConfig.GetPathToTemplate("form")
	if f.TemplateName != "" {
		templateName = f.TemplateName
	}
	return f.Renderer.RenderAsString(
		CurrentConfig.TemplatesFS, templateName,
		f, FuncMap, func1,
	)
}

func (f *Form) ProceedRequest(form *multipart.Form, gormModel interface{}) *FormError {
	formError := &FormError{
		FieldError: make(map[string]ValidationError),
		GeneralErrors: make(ValidationError, 0),
	}
	for fieldName, field := range f.FieldRegistry.GetAllFields() {
		if field.Name == "ID" {
			continue
		}
		errors := field.ProceedForm(form)
		if len(errors) == 0 {
			continue
		}
		formError.FieldError[fieldName] = errors
	}
	if formError.IsEmpty() {
		valueOfModel := reflect.ValueOf(gormModel)
		model := valueOfModel.Elem()
		for _, field := range f.FieldRegistry.GetAllFields() {
			if field.Name == "ID" {
				continue
			}
			modelF := model.FieldByName(field.Name)
			if !modelF.IsValid() {
				formError.AddGeneralError(fmt.Errorf("not valid field %s for model", field.Name))
				continue
			}
			if !modelF.CanSet() {
				formError.AddGeneralError(fmt.Errorf("can't set field %s for model", field.Name))
				continue
			}
			err := SetUpStructField(modelF, field.FieldConfig.Widget.GetOutputValue())
			if err != nil {
				formError.AddGeneralError(err)
			}
		}
	}
	return formError
}

func NewFormFromModel(gormModel interface{}, excludeFields []string, fieldsToShow []string, buildFieldPlacement bool, formTitle string) *Form {
	fieldRegistry := NewFieldRegistry()
	fieldsToShowRegistry := NewFieldRegistry()
	excludeFieldsRegistry := NewFieldRegistry()
	uadminDatabase := NewUadminDatabase()
	defer uadminDatabase.Close()
	statement := &gorm.Statement{DB: uadminDatabase.Db}
	statement.Parse(gormModel)
	r := NewTemplateRenderer(formTitle)
	fields := statement.Schema.Fields
	gormModelV := reflect.Indirect(reflect.ValueOf(gormModel))
	grouppedFields := make(map[string]*GrouppedFields)
	grouppedFields["default"] = &GrouppedFields{
		Rows: make([]*FormRow, 0),
		ExtraCssClasses: make([]string, 0),
		Name: "Default",
	}
	for _, field := range fields {
		if len(fieldsToShow) > 0 && !Contains(fieldsToShow, field.Name) {
			if !field.PrimaryKey {
				continue
			}
		}
		fieldToBeExcluded := Contains(excludeFields, field.Name)
		if len(excludeFields) >0 && fieldToBeExcluded {
			continue
		}
		uadminformtag := field.Tag.Get("uadminform")
		var fieldOptions IFieldFormOptions
		var uadminField *Field
		if uadminformtag != "" {
			fieldOptions = CurrentConfig.GetFieldFormOptions(uadminformtag)
			uadminField = NewFieldFromGormField(field, fieldOptions.GetWidgetType())
		} else {
			uadminField = NewFieldFromGormField(field, "")
		}
		uadminField.DisplayName = field.Name
		if uadminformtag != "" {
			uadminField.Initial = fieldOptions.GetInitial()
			if fieldOptions.GetDisplayName() != "" {
				uadminField.DisplayName = fieldOptions.GetDisplayName()
			}
			uadminField.Validators = fieldOptions.GetValidators()
			uadminField.Choices = fieldOptions.GetChoices()
			uadminField.HelpText = fieldOptions.GetHelpText()
			uadminField.WidgetType = fieldOptions.GetWidgetType()
			uadminField.ReadOnly = fieldOptions.GetReadOnly()
			uadminField.FieldConfig.Widget.SetReadonly(uadminField.ReadOnly)
		}
		uadminField.FieldConfig.Widget.RenderUsingRenderer(r)
		uadminField.FieldConfig.Widget.SetFieldDisplayName(field.Name)
		isTruthyValue := IsTruthyValue(gormModelV.FieldByName(field.Name).Interface())
		if isTruthyValue {
			v := TransformValueForWidget(gormModelV.FieldByName(field.Name).Interface())
			uadminField.FieldConfig.Widget.SetValue(v)
		}
		fieldRegistry.AddField(uadminField)
		formRow := &FormRow{
			Columns: make([]*ColumnSchema, 0),
		}
		formRow.Columns = append(formRow.Columns, &ColumnSchema{
			Fields: []*Field{uadminField},
		})
		if len(fieldsToShow) > 0 && Contains(fieldsToShow, field.Name) {
			fieldsToShowRegistry.AddField(uadminField)
			if !fieldToBeExcluded && buildFieldPlacement {
				grouppedFields["default"].Rows = append(grouppedFields["default"].Rows, formRow)
			}
		} else {
			if !fieldToBeExcluded {
				fieldsToShowRegistry.AddField(uadminField)
				grouppedFields["default"].Rows = append(grouppedFields["default"].Rows, formRow)
			}
		}
	}
	form := &Form{
		ExcludeFields: excludeFieldsRegistry,
		FieldsToShow: fieldsToShowRegistry,
		FieldRegistry: fieldRegistry,
		GroupsOfTheFields: &GrouppedFieldsRegistry{},
		Renderer: r,
		ExtraStatic: &StaticFiles{
			ExtraCSS: make([]string, 0),
			ExtraJS: make([]string, 0),
		},
	}
	form.GroupsOfTheFields.GrouppedFields = grouppedFields
	return form
}

func NewFormFromModelFromGinContext(contextFromGin IAdminContext, gormModel interface{}, excludeFields []string, fieldsToShow []string, buildFieldPlacement bool, formTitle string) *Form {
	form := NewFormFromModel(gormModel, excludeFields, fieldsToShow, buildFieldPlacement, formTitle)
	form.RequestContext = make(map[string]interface{})
	form.RequestContext["Language"] = contextFromGin.GetLanguage()
	form.RequestContext["RootURL"] = contextFromGin.GetRootURL()
	form.RequestContext["OTPImage"] = ""
	contextFromGin.SetForm(form)
	return form
}

func NewFieldFromGormField(gormField *schema.Field, forcedWidgetType string) *Field {
	var widget IWidget
	if forcedWidgetType != "" {
		widget = GetWidgetByWidgetType(forcedWidgetType)
		widget.InitializeAttrs()
		widget.SetName(gormField.Name)
		widget.SetValue(gormField.DefaultValueInterface)
	} else {
		uadminFieldType := GetUadminFieldTypeFromGormField(gormField)
		widget = GetWidgetFromUadminFieldTypeAndGormField(uadminFieldType, gormField)
	}
	field := &Field{
		Field: *gormField,
		UadminFieldType: GetUadminFieldTypeFromGormField(gormField),
		FieldConfig: &FieldConfig{Widget: widget},
		Required: gormField.NotNull && !gormField.HasDefaultValue,
		DisplayName: gormField.Name,
	}
	field.FieldConfig.Widget.SetRequired()
	return field
}


func GetWidgetFromUadminFieldTypeAndGormField(uadminFieldType UadminFieldType, gormField *schema.Field) IWidget {
	var widget IWidget
	switch uadminFieldType {
	case "biginteger":
	case "integer":
	case "positivebiginteger":
	case "positiveinteger":
	case "positivesmallinteger":
	case "smallinteger":
		widget = &NumberWidget{}
	case "binary":
		widget = &TextareaWidget{}
	case "char":
		widget = &TextWidget{}
		widget.SetAttr("maxlength", "1")
	case "boolean":
		widget = &CheckboxWidget{}
	case "decimal":
	case "float":
		widget = &NumberWidget{}
		widget.SetAttr("step", "0.1")
	case "email":
		widget = &EmailWidget{}
	case "file":
		widget = &FileWidget{}
	case "filepath":
		widget = &TextWidget{}
	case "text":
		widget = &TextWidget{}
	case "time":
		widget = &TimeWidget{}
	case "nullboolean":
		widget = &NullBooleanWidget{}
	case "slug":
		widget = &TextWidget{}
	case "url":
		widget = &URLWidget{}
	case "uuid":
		widget = &TextWidget{}
	case "date":
		widget = &DateWidget{}
	case "datetime":
		widget = &DateTimeWidget{}
	case "duration":
		widget = &TextWidget{}
	case "foreignkey":
		// @todo, integrate autocomplate widget
		widget = &TextWidget{}
	case "imagefield":
		widget = &FileWidget{}
		widget.SetAttr("accept", "image/*")
	case "ipaddress":
	case "genericipaddress":
		widget = &TextWidget{}
		widget.SetAttr("minlength", "7")
		widget.SetAttr("maxlength", "15")
		widget.SetAttr("size", "15")
		widget.SetAttr("pattern", "^((\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])$")
		// @todo, make sure we handle many to many widget type
		// const ManyToManyUadminFieldType UadminFieldType = "manytomany"
	default:
		widget = &TextWidget{}
	}
	widget.InitializeAttrs()
	widget.SetBaseFuncMap(FuncMap)
	widget.SetName(gormField.Name)
	widget.SetValue(gormField.DefaultValueInterface)
	return widget
}

type FormListEditable struct {
	FieldRegistry IFieldRegistry
}

func (f *FormListEditable) ExistsField(ld *ListDisplay) bool {
	_, err := f.FieldRegistry.GetByName(ld.Field.Name)
	return err == nil
}

func (f *FormListEditable) ProceedRequest(form *multipart.Form, gormModel interface{}) *FormError {
	formError := &FormError{
		FieldError: make(map[string]ValidationError),
		GeneralErrors: make(ValidationError, 0),
	}
	for fieldName, field := range f.FieldRegistry.GetAllFields() {
		errors := field.ProceedForm(form)
		if len(errors) == 0 {
			continue
		}
		formError.FieldError[fieldName] = errors
	}
	if formError.IsEmpty() {
		valueOfModel := reflect.ValueOf(gormModel)
		model := valueOfModel.Elem()
		for _, field := range f.FieldRegistry.GetAllFields() {
			modelF := model.FieldByName(field.Name)
			if !modelF.IsValid() {
				formError.AddGeneralError(fmt.Errorf("not valid field %s for model", field.Name))
				continue
			}
			if !modelF.CanSet() {
				formError.AddGeneralError(fmt.Errorf("can't set field %s for model", field.Name))
				continue
			}
			err := SetUpStructField(modelF, field.FieldConfig.Widget.GetOutputValue())
			if err != nil {
				formError.AddGeneralError(err)
			}
		}
	}
	return formError
}

func NewFormListEditableFromListDisplayRegistry(ID uint, model interface{}, listDisplayRegistry *ListDisplayRegistry) *FormListEditable {
	modelForm := NewFormFromModel(model, []string{}, []string{}, false, "")
	ret := &FormListEditable{FieldRegistry: NewFieldRegistry()}
	for ld := range listDisplayRegistry.GetAllFields() {
		if ld.IsEditable {
			fieldFromNewForm, _ := modelForm.FieldRegistry.GetByName(ld.Field.Name)
			name := fieldFromNewForm.FieldConfig.Widget.GetName(model, fieldFromNewForm)
			fieldFromNewForm.FieldConfig.Widget.SetName(fmt.Sprintf("%d_%s", ID, name))
			fieldFromNewForm.FieldConfig.Widget.SetShowOnlyHtmlInput()
			ret.FieldRegistry.AddField(fieldFromNewForm)
		}
	}
	return ret
}