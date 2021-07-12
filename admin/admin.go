package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/form"
	"github.com/uadmin/uadmin/interfaces"
	"sort"
)

type AdminActionPlacement struct {
	DisplayToTheTop bool
	DisplayToTheBottom bool
	DisplayToTheRight bool
	DisplayToTheLeft bool
	ShowOnTheListPage bool
}

type IAdminModelActionInterface interface {
	Handler (m interface{}, ctx *gin.Context)
	IsDisabled (m interface{}, ctx *gin.Context) bool
}

type AdminModelAction struct {
	IAdminModelActionInterface
	ActionName string
	HttpMethod string
	Description string
	ShowFutureChanges bool
	RedirectToRootModelPage bool
	Permissions *interfaces.Perm
	Placement *AdminActionPlacement
}

func NewAdminModelAction(actionName string, httpMethod string, perm *interfaces.Perm, placement *AdminActionPlacement) *AdminModelAction {
	return &AdminModelAction{
		RedirectToRootModelPage: true,
		ActionName: actionName,
		HttpMethod: httpMethod,
		Permissions: perm,
		Placement: placement,
	}
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
	Field *interfaces.Field
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
	Field *interfaces.Field
	UrlFilteringParam string
	OptionsToShow []*interfaces.FieldChoice
	FetchOptions func(m interface{}) []*interfaces.FieldChoice
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
	Field *interfaces.Field
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

type AdminPagesList []*AdminPage

func (apl AdminPagesList) Len() int { return len(apl) }
func (apl AdminPagesList) Less(i, j int) bool {
	if apl[i].Ordering == apl[j].Ordering {
		return apl[i].PageName < apl[j].PageName
	}
	return apl[i].Ordering < apl[j].Ordering
}
func (apl AdminPagesList) Swap(i, j int){ apl[i], apl[j] = apl[j], apl[i] }

type AdminPageRegistry struct {
	AdminPages map[string]*AdminPage
}

func (apr *AdminPageRegistry) GetBySlug(slug string) (*AdminPage, error){
	adminPage, ok := apr.AdminPages[slug]
	if !ok {
		return nil, fmt.Errorf("No admin page with alias %s", slug)
	}
	return adminPage, nil
}

func (apr *AdminPageRegistry) AddAdminPage(adminPage *AdminPage) error{
	apr.AdminPages[adminPage.Slug] = adminPage
	return nil
}

func (apr *AdminPageRegistry) GetAll() <- chan *AdminPage{
	chnl := make(chan *AdminPage)
	go func() {
		defer close(chnl)
		sortedPages := make(AdminPagesList, 0)

		for _, adminPage := range apr.AdminPages {
			sortedPages = append(sortedPages, adminPage)
		}
		sort.Reverse(sortedPages)
		for _, adminPage := range sortedPages {
			chnl <- adminPage
		}

	}()
	return chnl
}

func (apr *AdminPageRegistry) PreparePagesForTemplate() []byte {
	pages := make([]*AdminPage, 0)

	for page := range apr.GetAll() {
		pages = append(pages, page)
	}
	ret, _ := json.Marshal(pages)
	return ret
}

type DashboardAdminPanel struct {
	AdminPages *AdminPageRegistry
}

var CurrentDashboardAdminPanel *DashboardAdminPanel

func NewDashboardAdminPanel() *DashboardAdminPanel {
	return &DashboardAdminPanel{
		AdminPages: NewAdminPageRegistry(),
	}
}

type AdminPage struct {
	Actions []*AdminModelAction
	ActionsSelectionCounter bool
	DateHierarchyField string
	EmptyValueDisplay string
	ExcludeFields interfaces.IFieldRegistry
	FieldsToShow interfaces.IFieldRegistry
	FilterHorizontal bool
	FilterVertical bool
	Form *form.Form
	ShowAllFields bool
	Validators []interfaces.IValidator
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
	SubPages *AdminPageRegistry
	Ordering int
	PageName string
	Slug string
	ToolTip string
	Icon string
	ListHandler func (ctx *gin.Context)
	EditHandler func (ctx *gin.Context)
	AddHandler func (ctx *gin.Context)
	DeleteHandler func (ctx *gin.Context)
}

func (ap *AdminPage) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, _ := json.Marshal(ap.PageName)
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "PageName", string(jsonValue)))
	buffer.WriteString(",")
	jsonValue, _ = json.Marshal(ap.Slug)
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "Slug", string(jsonValue)))
	buffer.WriteString(",")
	jsonValue, _ = json.Marshal(ap.Icon)
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "Icon", string(jsonValue)))
	buffer.WriteString(",")
	jsonValue, _ = json.Marshal(ap.ToolTip)
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "ToolTip", string(jsonValue)))
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

type ConnectionToParentModel struct {
	FieldNameToValue map[string]interface{}
}

type AdminPageInlines struct {
	Ordering int
	Actions []*AdminModelAction
	EmptyValueDisplay string
	ExcludeFields interfaces.IFieldRegistry
	FieldsToShow interfaces.IFieldRegistry
	Form *form.Form
	ShowAllFields bool
	Validators []interfaces.IValidator
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
	Permissions *interfaces.Perm
}

func NewAdminPageRegistry() *AdminPageRegistry {
	return &AdminPageRegistry{AdminPages: make(map[string]*AdminPage)}
}
func NewAdminPage() *AdminPage {
	return &AdminPage{
		SubPages: NewAdminPageRegistry(),
		Validators: make([]interfaces.IValidator, 0),
		ExcludeFields: interfaces.NewFieldRegistry(),
		FieldsToShow: interfaces.NewFieldRegistry(),
		Actions: make([]*AdminModelAction, 0),
		SortBy: make([]*SortBy, 0),
		Inlines: make([]*AdminPageInlines, 0),
		ListDisplay: make([]*ListDisplay, 0),
		ListFilter: make([]*ListFilter, 0),
		SearchFields: make([]*SearchField, 0),
		ExtraStatic: &StaticFiles{
			ExtraCSS: make([]string, 0),
			ExtraJS: make([]string, 0),
		},
		Paginator: &Paginator{},
	}
}