package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/preloaded"
	"golang.org/x/net/html"
	"gorm.io/gorm/schema"
	"net/http"
	"strconv"
	"strings"
)

/*
	admin Tags
	read_only:TRUE
	email:TRUE
	hidden:TRUE
	html:TRUE
	fk:"ModelName"
	list:TRUE
	list_filter:TRUE
	search:TRUE
	dontCache:TRUE
  	required:TRUE
  	help:TRUE
  	pattern:TRUE
  	pattern_msg:"Message"
	max:"int"
	min:"int"
	link:TRUE
	file:TRUE
	dependsOn:""
	linkerObj:""
	linkerParentField:""
	linkerChildField:""
	childObj:""
	upload_to:"path"
	code:"true"
	money:"true" use on float
	defaultValue:""
*/

// commaf is a function to format number with thousand separator
// and two decimal points
func Commaf(j interface{}) string {
	v, _ := strconv.ParseFloat(fmt.Sprint(j), 64)
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}
	s := fmt.Sprintf("%.2f", v)

	comma := []byte{','}

	parts := strings.Split(s, ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func IsLocal(Addr string) bool {
	if strings.Contains(Addr, ":") && !strings.Contains(Addr, ".") {
		Addr = strings.TrimPrefix(Addr, "[")
		if strings.HasPrefix(Addr, "::") || strings.HasPrefix(Addr, "fc") || strings.HasPrefix(Addr, "fd") {
			return true
		}
	}
	p := strings.Split(strings.Split(Addr, ":")[0], ".")
	if len(p) != 4 {
		return false
	}
	_, err := strconv.ParseInt(p[2], 10, 64)
	if err != nil {
		return false
	}
	_, err = strconv.ParseInt(p[3], 10, 64)
	if err != nil {
		return false
	}
	v1, err := strconv.ParseInt(p[0], 10, 64)
	if err != nil {
		return false
	}
	v2, err := strconv.ParseInt(p[1], 10, 64)
	if err != nil {
		return false
	}
	if v1 == 10 {
		return true
	}
	if v1 == 172 {
		if v2 >= 16 && v2 <= 31 {
			return true
		}
	}
	if v1 == 192 && v2 == 168 {
		return true
	}
	if v1 == 127 {
		return true
	}
	return false
}

// saver is an interface to deal with form froms
type saver interface {
	Save()
}

// counter !
type Counter interface {
	Count(interface{}, interface{}, ...interface{}) int
}

type AdminPager interface {
	AdminPage(string, bool, int, int, interface{}, interface{}, ...interface{}) error
}

func PaginationHandler(itemCount int, PageLength int) (i int) {
	pCount := (float64(itemCount) / float64(PageLength))
	if pCount > float64(int(pCount)) {
		pCount++
	}
	i = int(pCount)
	if i == 1 {
		i--
	}
	return i
}

// toSnakeCase !
func toSnakeCase(str string) string {
	snake := preloaded.MatchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = preloaded.MatchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// JSONMarshal Generates JSON format from an object
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	// b, err := json.Marshal(v)
	b, err := json.MarshalIndent(v, "", " ")

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

// ReturnJSON returns json to the client
func ReturnJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		response := map[string]interface{}{
			"status":    "error",
			"error_msg": fmt.Sprintf("unable to encode JSON. %s", err),
		}
		b, _ = json.MarshalIndent(response, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}
	w.Write(b)
}

func StripHTMLScriptTag(v string) string {
	doc, err := html.Parse(strings.NewReader(v))
	if err != nil {
		return ""
	}
	removeScript(doc)
	b := bytes.NewBuffer([]byte{})
	if err := html.Render(b, doc); err != nil {
		return ""
	}
	return b.String()
}

func removeScript(n *html.Node) {
	// if note is script tag
	if n.Type == html.ElementNode && (strings.Contains(n.Data, "script") || strings.Contains(n.Data, "frame")) {
		n.Parent.RemoveChild(n)
		return // script tag is gone...
	}
	// traverse DOM
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeScript(c)
	}
}

type AdminPerm struct {
	Add bool
	Edit bool
	Delete bool
}

type ActionPlacement struct {
	PlacementName string
}

type AdminModelAction struct {
	ActionName string
	Description string
	Handler func(ctx *gin.Context)
	ShowFutureChanges bool
	RedirectToRootModelPage bool
	PermissionsRequired []*AdminPerm
	Placements []*ActionPlacement
}

type UadminFieldType string

func (uft UadminFieldType) BigIntegerField() string {
	return "biginteger"
}

func (uft UadminFieldType) BinaryField() string {
	return "binary"
}

func (uft UadminFieldType) BooleanField() string {
	return "boolean"
}

func (uft UadminFieldType) CharField() string {
	return "char"
}

func (uft UadminFieldType) DateField() string {
	return "date"
}

func (uft UadminFieldType) DateTimeField() string {
	return "datetime"
}

func (uft UadminFieldType) DecimalField() string {
	return "decimal"
}

func (uft UadminFieldType) DurationField() string {
	return "duration"
}

func (uft UadminFieldType) EmailField() string {
	return "email"
}

func (uft UadminFieldType) FileField() string {
	return "file"
}

func (uft UadminFieldType) FilePathField() string {
	return "filepath"
}

func (uft UadminFieldType) FloatField() string {
	return "float"
}

func (uft UadminFieldType) ForeignKey() string {
	return "foreignkey"
}

func (uft UadminFieldType) ImageField() string {
	return "imagefield"
}

func (uft UadminFieldType) IntegerField() string {
	return "integer"
}

func (uft UadminFieldType) IpAddressField() string {
	return "ipaddress"
}

func (uft UadminFieldType) GenericIpAddressField() string {
	return "genericipaddress"
}

func (uft UadminFieldType) JSONField() string {
	return "json"
}

func (uft UadminFieldType) ManyToManyField() string {
	return "manytomany"
}

func (uft UadminFieldType) NullBooleanField() string {
	return "nullboolean"
}

func (uft UadminFieldType) PositiveBigIntegerField() string {
	return "positivebiginteger"
}

func (uft UadminFieldType) PositiveIntegerField() string {
	return "positiveinteger"
}

func (uft UadminFieldType) PositiveSmallIntegerField() string {
	return "positivesmallinteger"
}

func (uft UadminFieldType) SlugField() string {
	return "slug"
}

func (uft UadminFieldType) SmallIntegerField() string {
	return "smallinteger"
}

func (uft UadminFieldType) TextField() string {
	return "text"
}

func (uft UadminFieldType) TimeField() string {
	return "time"
}

func (uft UadminFieldType) URLField() string {
	return "url"
}

func (uft UadminFieldType) UUIDField() string {
	return "uuid"
}

type WidgetType string

func (w WidgetType) TextInput() string {
	return "text"
}

func (w WidgetType) NumberInput() string {
	return "number"
}

func (w WidgetType) EmailInput() string {
	return "email"
}

func (w WidgetType) URLInput() string {
	return "url"
}

func (w WidgetType) PasswordInput() string {
	return "password"
}

func (w WidgetType) HiddenInput() string {
	return "hidden"
}

func (w WidgetType) DateInput() string {
	return "date"
}

func (w WidgetType) DateTimeInput() string {
	return "datetime"
}

func (w WidgetType) TimeInput() string {
	return "time"
}

func (w WidgetType) Textarea() string {
	return "textarea"
}

func (w WidgetType) CheckboxInput() string {
	return "checkbox"
}

func (w WidgetType) Select() string {
	return "select"
}

func (w WidgetType) NullBooleanSelect() string {
	return "nullboolean"
}

func (w WidgetType) SelectMultiple() string {
	return "selectmultiple"
}

func (w WidgetType) RadioSelect() string {
	return "radioselect"
}

func (w WidgetType) CheckboxSelectMultiple() string {
	return "checkboxselectmultiple"
}

func (w WidgetType) FileInput() string {
	return "fileinput"
}

func (w WidgetType) ClearableFileInput() string {
	return "clearablefileinput"
}

func (w WidgetType) MultipleHiddenInput() string {
	return "multiplehiddeninput"
}

func (w WidgetType) SplitDateTime() string {
	return "splitdatetime"
}

func (w WidgetType) SplitHiddenDateTime() string {
	return "splithiddendatetime"
}

func (w WidgetType) SelectDateWidget() string {
	return "selectdatewidget"
}

type IWidget interface {
	IdForLabel(model interface{}, f interface{})
	FormatValue(v string, model interface{})
	GetWidgetType() WidgetType
	GetAttrs() map[string]string
	GetTemplateName() string
	RenderUsingRenderer(renderer ITemplateRenderer)
}
type Widget struct {
	WidgetType WidgetType
	Attrs map[string]string
	CustomFormatValue func(v string, model interface{})
	TemplateName string
}


type FieldConfig struct {
	Widget *Widget
	AutocompleteURL string
	DependsOnAnotherFields []string
}

type Choice struct {
	DisplayAs string
	Value interface{}
}

type Field struct {
	schema.Field
	ReadOnly bool
	UadminFieldType UadminFieldType
	FieldConfig *FieldConfig
	Required bool
	DisplayName string
	HelpText string
	Choices []*Choice
	Validators []govalidator.CustomTypeValidator
}

type ColumnSchema struct {
	ShowLabel bool
	Field *Field
}

type RowSchema struct {
	Columns []*ColumnSchema
}

type GrouppedFields struct {
	RowSchema *RowSchema
	ExtraCssClasses []string
	Description string
}

type AdminForm struct {
	FieldsToShow []string
	CustomDefinedFields map[string]*Field
}

var GloballyAvailableActionsForAdmin = make([]*AdminModelAction, 0)

type SortBy struct {
	Direction int // -1 descending order, 1 ascending order
	FieldName string
}

type ListDisplay struct {
	DisplayName string
	Field *Field
	Callable func(m interface{}) string
	UadminFieldType UadminFieldType
	ChangeLink bool
	Editable bool
}

type ListFilter struct {
	Title string
	Field *Field
	Callable func(m interface{}) string
	UadminFieldType UadminFieldType
	UrlFilteringParam string
	OptionsToShow []*Choice
	Template string
}

type PrepopulatedField struct {
	Field *Field
	// @todo, make it a part of interface
	// Populate(m interface{}) string
}

type SearchField struct {
	Field *Field
}

type Paginator struct {
	PerPage int
	AllowEmptyFirstPage bool
	ShowLastPageOnPreviousPage bool
	Count int
	NumPages int
	Template string
}

type AdminPage struct {
	Actions []*AdminModelAction
	DisabledActions []string
	ActionsSelectionCounter bool
	DateHierarchyField string
	EmptyValueDisplay string
	ExcludeFields []string
	FieldsToShow []string
	GroupsOfTheFields []*GrouppedFields
	FilterHorizontal bool
	FilterVertical bool
	Form *AdminForm
	ShowAllFields bool
	Validators []govalidator.CustomTypeValidator
	InitialData map[string]interface{}
	SortBy []*SortBy
	Inlines []*AdminPageInlines
	ListDisplay []*ListDisplay
	ListFilter []*ListFilter
	MaxShowAll int
	PerPage int
	OrderBy []*SortBy
	PrepopulatedFields []*PrepopulatedField
	PreserveFilters bool
	SaveAndContinue bool
	SaveOnTop bool
	SearchFields []*SearchField
	ShowFullResultCount bool
	ListSortingDisabledForFields []*ListDisplay
	ViewOnSite bool
	ListTemplate string
	AddTemplate string
	EditTemplate string
	DeleteConfirmationTemplate string
	DeleteSelectedConfirmationTemplate string
	ObjectHistoryTemplate string
	PopupResponseTemplate string
	ExtraCSS []string
	ExtraJS []string
	Ordering int
	Paginator *Paginator
}

type ConnectionToParentModel struct {
	FieldNameToValue map[string]interface{}
}

type AdminPageInlines struct {
	Actions []*AdminModelAction
	DisabledActions []string
	EmptyValueDisplay string
	ExcludeFields []string
	FieldsToShow []string
	GroupsOfTheFields []*GrouppedFields
	FilterHorizontal bool
	FilterVertical bool
	Form *AdminForm
	ShowAllFields bool
	Validators []govalidator.CustomTypeValidator
	InitialData map[string]interface{}
	SortBy []*SortBy
	PrepopulatedFields []*PrepopulatedField
	Classes []string
	Extra int
	MaxNum int
	MinNum int
	VerboseName string
	VerboseNamePlural string
	CanDelete bool
	ShowChangeLink bool
	ConnectionToParentModel ConnectionToParentModel
	Template string
	Perm AdminPerm
}

// @todo, DO IT! analyze later
// permissions
// Calling save_m2m() is only required if you use save(commit=False). When you use a save() on a form, all data – including many-to-many data – is saved without the need for any additional method calls. For example:
// In addition, Django applies the following rule: if you set editable=False on the model field, any form created from the model via ModelForm will not include that field.
// After calling save(), your model formset will have three new attributes containing the formset’s changes:
//
//models.BaseModelFormSet.changed_objects¶
//models.BaseModelFormSet.deleted_objects¶
//models.BaseModelFormSet.new_objects¶
// raw id widget
//type PrepopulatedField struct {
//	Field *Field
//	// @todo, make it a part of interface
//	// Populate(m interface{}) string
//}
// type AdminPage struct {
// 	ViewOnSiteHandler func
// }
// https://docs.djangoproject.com/en/3.2/ref/contrib/admin/#modeladmin-methods

func BuildAdminModelSchema(modelI interface{}, adminFormSchema *AdminPage) {
	// modelType := reflect.TypeOf(modelI)
	// nFields := modelType.NumField()
}