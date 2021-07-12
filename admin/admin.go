package admin

import (
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
	ret, err := json.Marshal(pages)
	if err != nil {
		interfaces.Trail(interfaces.CRITICAL, "error while generating menu in admin", err)
	}
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
	Actions []*AdminModelAction `json:"-"`
	ActionsSelectionCounter bool `json:"-"`
	DateHierarchyField string `json:"-"`
	EmptyValueDisplay string `json:"-"`
	ExcludeFields interfaces.IFieldRegistry `json:"-"`
	FieldsToShow interfaces.IFieldRegistry `json:"-"`
	FilterHorizontal bool `json:"-"`
	FilterVertical bool `json:"-"`
	Form *form.Form `json:"-"`
	ShowAllFields bool `json:"-"`
	Validators []interfaces.IValidator `json:"-"`
	SortBy []*SortBy `json:"-"`
	Inlines []*AdminPageInlines `json:"-"`
	ListDisplay []*ListDisplay `json:"-"`
	ListFilter []*ListFilter `json:"-"`
	MaxShowAll int `json:"-"`
	PreserveFilters bool `json:"-"`
	SaveAndContinue bool `json:"-"`
	SaveOnTop bool `json:"-"`
	SearchFields []*SearchField `json:"-"`
	ShowFullResultCount bool `json:"-"`
	ViewOnSite bool `json:"-"`
	ListTemplate string `json:"-"`
	AddTemplate string `json:"-"`
	EditTemplate string `json:"-"`
	DeleteConfirmationTemplate string `json:"-"`
	DeleteSelectedConfirmationTemplate string `json:"-"`
	ObjectHistoryTemplate string `json:"-"`
	PopupResponseTemplate string `json:"-"`
	ExtraStatic *StaticFiles `json:"-"`
	Paginator *Paginator `json:"-"`
	SubPages *AdminPageRegistry `json:"-"`
	Ordering int
	PageName string
	Slug string
	ToolTip string
	Icon string
	ListHandler func (ctx *gin.Context) `json:"-"`
	EditHandler func (ctx *gin.Context) `json:"-"`
	AddHandler func (ctx *gin.Context) `json:"-"`
	DeleteHandler func (ctx *gin.Context) `json:"-"`
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