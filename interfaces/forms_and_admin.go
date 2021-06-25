package interfaces

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"
)

type PermissionRegistry interface {
	DoesUserHaveRightFor(permissionName string) bool
	AddCustomPermission(permission CustomPermission)
}
type CustomPermission string

type Perm struct {
	PermissionRegistry
	PermBitInteger PermBit
	CustomPermissions []CustomPermission
}

func (ap *Perm) HasReadPermission() bool {
	return (ap.PermBitInteger & ReadPermBit) == ReadPermBit
}

func (ap *Perm) DoesUserHaveRightFor(permissionName string) bool {
	return false
}

func (ap *Perm) AddCustomPermission(permission CustomPermission){
	ap.CustomPermissions = append(ap.CustomPermissions, permission)
}

func (ap *Perm) HasAddPermission() bool {
	return (ap.PermBitInteger & AddPermBit) == AddPermBit
}

func (ap *Perm) HasEditPermission() bool {
	return (ap.PermBitInteger & EditPermBit) == EditPermBit
}

func (ap *Perm) HasDeletePermission() bool {
	return (ap.PermBitInteger & DeletePermBit) == DeletePermBit
}

func (ap *Perm) HasPublishPermission() bool {
	return (ap.PermBitInteger & PublishPermBit) == PublishPermBit
}

func (ap *Perm) HasRevertPermission() bool {
	return (ap.PermBitInteger & RevertPermBit) == RevertPermBit
}

type PermBit int

var ReadPermBit PermBit = 0
var AddPermBit PermBit = 2
var EditPermBit PermBit = 4
var DeletePermBit PermBit = 8
var PublishPermBit PermBit = 16
var RevertPermBit PermBit = 32

type AdminActionPlacement struct {
	DisplayToTheTop bool
	DisplayToTheBottom bool
	DisplayToTheRight bool
	DisplayToTheLeft bool
}

type IAdminModelActionInterface interface {
	Handler (m interface{}, ctx *gin.Context)
	IsDisabled (m interface{}, ctx *gin.Context) bool
}

type AdminModelAction struct {
	IAdminModelActionInterface
	ActionName string
	Description string
	ShowFutureChanges bool
	RedirectToRootModelPage bool
	Permissions *Perm
	Placement *AdminActionPlacement
}

type UadminFieldType string

var BigIntegerUadminFieldType UadminFieldType = "biginteger"
var BinaryUadminFieldType UadminFieldType = "binary"
var BooleanUadminFieldType UadminFieldType = "boolean"
var CharUadminFieldType UadminFieldType = "char"
var DateUadminFieldType UadminFieldType = "date"
var DateTimeUadminFieldType UadminFieldType = "datetime"
var DecimalUadminFieldType UadminFieldType = "decimal"
var DurationUadminFieldType UadminFieldType = "duration"
var EmailUadminFieldType UadminFieldType = "email"
var FileUadminFieldType UadminFieldType = "file"
var FilePathUadminFieldType UadminFieldType = "filepath"
var FloatUadminFieldType UadminFieldType = "float"
var ForeignKeyUadminFieldType UadminFieldType = "foreignkey"
var ImageFieldUadminFieldType UadminFieldType = "imagefield"
var IntegerUadminFieldType UadminFieldType = "integer"
var IpAddressUadminFieldType UadminFieldType = "ipaddress"
var GenericIpAddressUadminFieldType UadminFieldType = "genericipaddress"
var ManyToManyUadminFieldType UadminFieldType = "manytomany"
var NullBooleanUadminFieldType UadminFieldType = "nullboolean"
var PositiveBigIntegerUadminFieldType UadminFieldType = "positivebiginteger"
var PositiveIntegerUadminFieldType UadminFieldType = "positiveinteger"
var PositiveSmallIntegerUadminFieldType UadminFieldType = "positivesmallinteger"
var SlugUadminFieldType UadminFieldType = "slug"
var SmallIntegerUadminFieldType UadminFieldType = "smallinteger"
var TextUadminFieldType UadminFieldType = "text"
var TimeUadminFieldType UadminFieldType = "time"
var URLUadminFieldType UadminFieldType = "url"
var UUIDUadminFieldType UadminFieldType = "uuid"

type WidgetType string

var TextInputWidgetType WidgetType = "text"
var NumberInputWidgetType WidgetType = "number"
var EmailInputWidgetType WidgetType = "email"
var URLInputWidgetType WidgetType = "url"
var PasswordInputWidgetType WidgetType = "password"
var HiddenInputWidgetType WidgetType = "hidden"
var DateInputWidgetType WidgetType = "date"
var DateTimeInputWidgetType WidgetType = "datetime"
var TimeInputWidgetType WidgetType = "time"
var TextareaInputWidgetType WidgetType = "textarea"
var CheckboxInputWidgetType WidgetType = "checkbox"
var SelectWidgetType WidgetType = "select"
var NullBooleanWidgetType WidgetType = "nullboolean"
var SelectMultipleWidgetType WidgetType = "selectmultiple"
var RadioSelectWidgetType WidgetType = "radioselect"
var CheckboxSelectMultipleWidgetType WidgetType = "checkboxselectmultiple"
var FileInputWidgetType WidgetType = "fileinput"
var ClearableFileInputWidgetType WidgetType = "clearablefileinput"
var MultipleHiddenInputWidgetType WidgetType = "multiplehiddeninput"
var SplitDateTimeWidgetType WidgetType = "splitdatetime"
var SplitHiddenDateTimeWidgetType WidgetType = "splithiddendatetime"
var SelectDateWidgetType WidgetType = "selectdate"

type IWidget interface {
	IdForLabel(model interface{}, F *Field)
	FormatValue(v interface{}, model interface{})
	GetWidgetType() WidgetType
	GetAttrs() map[string]string
	GetTemplateName() string
	RenderUsingRenderer(renderer ITemplateRenderer)
}

type IWidgetInterface interface {
	CustomFormatValue (v string, model interface{})
}

type Widget struct {
	IWidgetInterface
	ConcreteWidget IWidget
	Attrs map[string]string
	TemplateName string
}

type FieldConfig struct {
	Widget *Widget
	AutocompleteURL string
	DependsOnAnotherFields []string
}

type FieldChoice struct {
	DisplayAs string
	Value interface{}
}

type IFieldChoiceRegistryInterface interface {
	IsValidChoice (v interface{}) bool
}

type FieldChoiceRegistry struct {
	IFieldChoiceRegistryInterface
	Choices []*FieldChoice
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
	Validators []govalidator.CustomTypeValidator
	SortingDisabled bool
	Populate func(field *Field, m interface{}) interface{}
}

type IFieldRegistry interface {
	GetByName(name string) (*Field, error)
	AddField(field *Field)
	GetAllFields() map[string]*Field
}

type FieldRegistry struct {
	IFieldRegistry
	Fields map[string]*Field
}

func (fr *FieldRegistry) GetByName(name string) (*Field, error) {
	f, ok := fr.Fields[name]
	if !ok {
		return nil, fmt.Errorf("No field %s found", name)
	}
	return f, nil
}

func (fr *FieldRegistry) GetAllFields() map[string]*Field {
	return fr.Fields
}

func (fr *FieldRegistry) AddField(field *Field) {
	if _, err := fr.GetByName(field.Name); err != nil {
		panic(err)
	}
	fr.Fields[field.Name] = field
}

type ColumnSchema struct {
	ShowLabel bool
	FieldRegistry *FieldRegistry
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

type Form struct {
	FieldsToShow *FieldChoiceRegistry
	FieldRegistry *FieldRegistry
	GroupsOfTheFields *GrouppedFieldsRegistry
}

var GloballyAvailableActionsForAdmin = make([]*AdminModelAction, 0)

func AddGloballyAvailableAction(modelAction *AdminModelAction) {
	GloballyAvailableActionsForAdmin = append(GloballyAvailableActionsForAdmin, modelAction)
}

type ISortBy interface {
	IsApplicableTo (ctx *gin.Context, m interface{}) bool
	Sort (m interface{})
}

type SortBy struct {
	ISortBy
	Direction int // -1 descending order, 1 ascending order
	FieldName string
}

func (sb *SortBy) IsApplicableTo(ctx *gin.Context, m interface {}) bool {
	return false
}

func (sb *SortBy) Sort(m interface {}) {
}

type IListDisplayInterface interface {
	GetValue (m interface{}) interface{}
}

type ListDisplay struct {
	IListDisplayInterface
	DisplayName string
	Field *Field
	ChangeLink bool
	Editable bool
	Ordering int
}

func (ld *ListDisplay) GetValue(m interface{}) interface {} {
	return 1
}

type IListFilterInterface interface {
	FilterQs (m interface{}) interface{}
}

type ListFilter struct {
	IListFilterInterface
	Title string
	Field *Field
	UrlFilteringParam string
	OptionsToShow []*FieldChoice
	FetchOptions func(m interface{}) []*FieldChoice
	Template string
	Ordering int
}

func (lf *ListFilter) FilterQs (m interface{}) interface{} {
	return true
}

type ISearchFieldInterface interface {
	Search (m interface{})
}

type SearchField struct {
	ISearchFieldInterface
	Field *Field
}

func (sf *SearchField) Search(m interface{}) {

}

type PaginationType string

var LimitPaginationType PaginationType = "limit"
var CursorPaginationType PaginationType = "cursor"

type IPaginationInterface interface {
	Paginate (m interface{})
}

type Paginator struct {
	IPaginationInterface
	PerPage int
	AllowEmptyFirstPage bool
	ShowLastPageOnPreviousPage bool
	Count int
	NumPages int
	Template string
	PaginationType PaginationType
}

func (p *Paginator) Paginate(m interface{}) {

}

type StaticFiles struct {
	ExtraCSS []string
	ExtraJS []string
}

type AdminPage struct {
	Actions []*AdminModelAction
	ActionsSelectionCounter bool
	DateHierarchyField string
	EmptyValueDisplay string
	ExcludeFields IFieldRegistry
	FieldsToShow IFieldRegistry
	FilterHorizontal bool
	FilterVertical bool
	Form *Form
	ShowAllFields bool
	Validators []govalidator.CustomTypeValidator
	SortBy []*SortBy
	Inlines []*AdminPageInlines
	ListDisplay []*ListDisplay
	ListFilter []*ListFilter
	MaxShowAll int
	PreserveFilters bool
	SaveAndContinue bool
	SaveOnTop bool
	SearchFields []*SearchField
	ShowFullResultCount bool
	ViewOnSite bool
	ListTemplate string
	AddTemplate string
	EditTemplate string
	DeleteConfirmationTemplate string
	DeleteSelectedConfirmationTemplate string
	ObjectHistoryTemplate string
	PopupResponseTemplate string
	ExtraStatic *StaticFiles
	Paginator *Paginator
}

type ConnectionToParentModel struct {
	FieldNameToValue map[string]interface{}
}

type AdminPageInlines struct {
	Ordering int
	Actions []*AdminModelAction
	EmptyValueDisplay string
	ExcludeFields IFieldRegistry
	FieldsToShow IFieldRegistry
	Form *Form
	ShowAllFields bool
	Validators []govalidator.CustomTypeValidator
	SortBy []*SortBy
	Classes []string
	Extra int
	MaxNum int
	MinNum int
	VerboseName string
	VerboseNamePlural string
	ShowChangeLink bool
	ConnectionToParentModel ConnectionToParentModel
	Template string
	Permissions *Perm
}
