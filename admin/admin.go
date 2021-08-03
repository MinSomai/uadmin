package admin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uadmin/uadmin/form"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/template"
	"github.com/uadmin/uadmin/templatecontext"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AdminActionPlacement struct {
	DisplayToTheTop bool
	DisplayToTheBottom bool
	DisplayToTheRight bool
	DisplayToTheLeft bool
	ShowOnTheListPage bool
}

type IAdminModelActionInterface interface {
}

type AdminModelAction struct {
	IAdminModelActionInterface
	ActionName string
	Description string
	ShowFutureChanges bool
	RedirectToRootModelPage bool
	Placement *AdminActionPlacement
	PermName interfaces.CustomPermission
	Handler func (afo *AdminFilterObjects) (bool, int64)
	IsDisabled func (afo *AdminFilterObjects, ctx *gin.Context) bool
	SlugifiedActionName string
}

func prepareAdminModelActionName(adminModelAction string) string {
	slugifiedAdminModelAction := interfaces.AsciiRegex.ReplaceAllLiteralString(adminModelAction, "")
	slugifiedAdminModelAction = strings.Replace(strings.ToLower(slugifiedAdminModelAction), " ", "_", -1)
	slugifiedAdminModelAction = strings.Replace(strings.ToLower(slugifiedAdminModelAction), ".", "_", -1)
	return slugifiedAdminModelAction
}

func NewAdminModelAction(actionName string, placement *AdminActionPlacement) *AdminModelAction {
	return &AdminModelAction{
		RedirectToRootModelPage: true,
		ActionName: actionName,
		Placement: placement,
		SlugifiedActionName: prepareAdminModelActionName(actionName),
	}
}

type AdminFilterObjects struct {
	InitialGormQuerySet *gorm.DB
	GormQuerySet *gorm.DB
	PaginatedGormQuerySet *gorm.DB
	Model interface{}
	UadminDatabase *interfaces.UadminDatabase
	GenerateModelI func() interface{}
}

type IterateAdminObjects struct {
	Model interface {}
	Id uint
}

func (afo *AdminFilterObjects) GetPaginated() <- chan *IterateAdminObjects {
	chnl := make(chan *IterateAdminObjects)
	go func() {
		defer close(chnl)
		rows, _ := afo.PaginatedGormQuerySet.Rows()
		defer rows.Close()
		db := interfaces.NewUadminDatabase()
		for rows.Next() {
			model := afo.GenerateModelI()
			db.Db.ScanRows(rows, model)
			statement := &gorm.Statement{DB: afo.UadminDatabase.Db}
			statement.Parse(model)
			gormModelV := reflect.Indirect(reflect.ValueOf(model))
			Id := interfaces.TransformValueForWidget(gormModelV.FieldByName(statement.Schema.PrimaryFields[0].Name).Interface())
			IdN, _ := strconv.Atoi(Id.(string))
			yieldV := &IterateAdminObjects{
				Model: model,
				Id: uint(IdN),
			}
			chnl <- yieldV
		}
	}()
	return chnl
}

type ISortInterface interface {
	Order(afo *AdminFilterObjects)
}

type ISortBy interface {
	Sort (afo *AdminFilterObjects, direction int)
}

type SortBy struct {
	ISortBy
	Direction int // -1 descending order, 1 ascending order
	Field *interfaces.Field
}

func (sb *SortBy) Sort(afo *AdminFilterObjects, direction int) {
	sortBy := sb.Field.DBName
	if direction == -1 {
		sortBy += " desc"
	}
	afo.PaginatedGormQuerySet = afo.PaginatedGormQuerySet.Order(sortBy)
}

type ListDisplayRegistry struct {
	ListDisplayFields map[string]*ListDisplay
}

func (ldr *ListDisplayRegistry) AddField(ld *ListDisplay) {
	ldr.ListDisplayFields[ld.DisplayName] = ld
}

func (ldr *ListDisplayRegistry) GetAllFields() <- chan *ListDisplay {
	chnl := make(chan *ListDisplay)
	go func() {
		defer close(chnl)
		dFields := make([]*ListDisplay, 0)
		for _, dField := range ldr.ListDisplayFields {
			dFields = append(dFields, dField)
		}
		sort.Slice(dFields, func(i, j int) bool {
			if dFields[i].Ordering == dFields[j].Ordering {
				return dFields[i].DisplayName < dFields[j].DisplayName
			}
			return dFields[i].Ordering < dFields[j].Ordering
		})
		for _, dField := range dFields {
			chnl <- dField
		}
	}()
	return chnl
}

func (ldr *ListDisplayRegistry) GetFieldByDisplayName(displayName string) (*ListDisplay, error) {
	listField, exists := ldr.ListDisplayFields[displayName]
	if !exists {
		return nil, fmt.Errorf("found no display field with name %s", displayName)
	}
	return listField, nil
}

func NewAdminModelActionRegistry() *AdminModelActionRegistry {
	return &AdminModelActionRegistry{AdminModelActions: make(map[string]*AdminModelAction)}
}

type AdminModelActionRegistry struct {
	AdminModelActions map[string]*AdminModelAction
}

func (amar *AdminModelActionRegistry) AddModelAction(ma *AdminModelAction) {
	amar.AdminModelActions[ma.SlugifiedActionName] = ma
}

func (amar *AdminModelActionRegistry) GetAllModelActions() <- chan *AdminModelAction {
	chnl := make(chan *AdminModelAction)
	go func() {
		defer close(chnl)
		mActions := make([]*AdminModelAction, 0)
		for _, mAction := range amar.AdminModelActions {
			mActions = append(mActions, mAction)
		}
		sort.Slice(mActions, func(i, j int) bool {
			return mActions[i].ActionName < mActions[j].ActionName
		})
		for _, mAction := range mActions {
			chnl <- mAction
		}
	}()
	return chnl
}

func (amar *AdminModelActionRegistry) GetModelActionByName(actionName string) (*AdminModelAction, error) {
	mAction, exists := amar.AdminModelActions[actionName]
	if !exists {
		return nil, fmt.Errorf("found no model action with name %s", actionName)
	}
	return mAction, nil
}

type IListDisplayInterface interface {
	GetValue (m interface{}) string
}

type ListDisplay struct {
	IListDisplayInterface
	DisplayName string
	Field *interfaces.Field
	ChangeLink bool
	Editable bool
	Ordering int
	SortBy *SortBy
	Populate func (m interface{}) string
	MethodName string

}

func (ld *ListDisplay) GetOrderingName(initialOrdering []string) string {
	for _, part := range initialOrdering {
		negativeOrdering := false
		if strings.HasPrefix(part, "-") {
			part = part[1:]
			negativeOrdering = true
		}
		if part == ld.DisplayName {
			if negativeOrdering {
				return ld.DisplayName
			}
			return "-" + ld.DisplayName
		}
	}

	return ld.DisplayName
}

func (ld *ListDisplay) IsEligibleForOrdering() bool {
	return ld.SortBy != nil
}

func (ld *ListDisplay) GetValue(m interface{}) string {
	if ld.MethodName != "" {
		values := reflect.ValueOf(m).MethodByName(ld.MethodName).Call([]reflect.Value{})
		return values[0].String()
	}
	if ld.Populate != nil {
		return ld.Populate(m)
	}
	gormModelV := reflect.Indirect(reflect.ValueOf(m))
	return interfaces.TransformValueForListDisplay(gormModelV.FieldByName(ld.Field.Name).Interface())
}

func NewListDisplay(field *interfaces.Field) *ListDisplay {
	displayName := ""
	if field != nil {
		displayName = field.DisplayName
	}
	return &ListDisplay{
		DisplayName: displayName, Field: field, ChangeLink: true, Editable: true,
		SortBy: &SortBy{Field: field, Direction: 1},
	}
}

type IListFilterInterface interface {
	FilterQs (afo *AdminFilterObjects, filterString string)
}

type ListFilter struct {
	IListFilterInterface
	Title string
	UrlFilteringParam string
	OptionsToShow []*interfaces.FieldChoice
	FetchOptions func(m interface{}) []*interfaces.FieldChoice
	CustomFilterQs func(afo *AdminFilterObjects, filterString string)
	Template string
	Ordering int
}

func (lf *ListFilter) FilterQs (afo *AdminFilterObjects, filterString string) {
	if lf.CustomFilterQs != nil {
		lf.CustomFilterQs(afo, filterString)
	} else {
		statement := &gorm.Statement{DB: afo.UadminDatabase.Db}
		statement.Parse(afo.Model)
		schema1 := statement.Schema
		operatorContext := interfaces.FilterGormModel(afo.UadminDatabase.Adapter, afo.GormQuerySet, schema1, []string{filterString}, afo.Model)
		afo.GormQuerySet = operatorContext.Tx
		operatorContext = interfaces.FilterGormModel(afo.UadminDatabase.Adapter, afo.PaginatedGormQuerySet, schema1, []string{filterString}, afo.Model)
		afo.PaginatedGormQuerySet = operatorContext.Tx
	}
}

func (lf *ListFilter) IsItActive (fullUrl *url.URL) bool {
	return strings.Contains(fullUrl.String(), lf.UrlFilteringParam)
}

func (lf *ListFilter) GetURLToClearFilter(fullUrl *url.URL) string {
	clonedUrl := interfaces.CloneNetUrl(fullUrl)
	qs := clonedUrl.Query()
	qs.Del(lf.UrlFilteringParam)
	clonedUrl.RawQuery = qs.Encode()
	return clonedUrl.String()
}

func (lf *ListFilter) IsThatOptionActive(option *interfaces.FieldChoice, fullUrl *url.URL) bool {
	qs := fullUrl.Query()
	value := qs.Get(lf.UrlFilteringParam)
	if value != "" {
		optionValue := interfaces.TransformValueForListDisplay(option.Value)
		if optionValue == value {
			return true
		}
	}
	return false
}

func (lf *ListFilter) GetURLForOption(option *interfaces.FieldChoice, fullUrl *url.URL) string {
	clonedUrl := interfaces.CloneNetUrl(fullUrl)
	qs := clonedUrl.Query()
	qs.Set(lf.UrlFilteringParam, interfaces.TransformValueForListDisplay(option.Value))
	clonedUrl.RawQuery = qs.Encode()
	return clonedUrl.String()
}

type ListFilterRegistry struct {
	ListFilter []*ListFilter
}

type ListFilterList [] *ListFilter

func (apl ListFilterList) Len() int { return len(apl) }
func (apl ListFilterList) Less(i, j int) bool {
	return apl[i].Ordering < apl[j].Ordering
}
func (apl ListFilterList) Swap(i, j int){ apl[i], apl[j] = apl[j], apl[i] }



func (lfr *ListFilterRegistry) Iterate() <- chan *ListFilter {
	chnl := make(chan *ListFilter)
	go func() {
		lfList := make(ListFilterList, 0)
		defer close(chnl)
		for _, lF := range lfr.ListFilter {
			lfList = append(lfList, lF)
		}
		sort.Sort(lfList)
		for _, lf := range lfList {
			chnl <- lf
		}
	}()
	return chnl
}

func (lfr *ListFilterRegistry) IsEmpty() bool {
	return !(len(lfr.ListFilter) > 0)
}

func (lfr *ListFilterRegistry) Add(lf *ListFilter) {
	lfr.ListFilter = append(lfr.ListFilter, lf)
}

type ISearchFieldInterface interface {
	Search (afo *AdminFilterObjects, searchString string)
}

type SearchField struct {
	ISearchFieldInterface
	Field *schema.Field
	CustomSearch func(afo *AdminFilterObjects, searchString string)
}

func (sf *SearchField) Search(afo *AdminFilterObjects, searchString string) {
	if sf.CustomSearch != nil {
		sf.CustomSearch(afo, searchString)
	} else {
		operator := interfaces.ExactGormOperator{}
		gormOperatorContext := interfaces.NewGormOperatorContext(afo.GormQuerySet, afo.Model)
		operator.Build(afo.UadminDatabase.Adapter, gormOperatorContext, sf.Field, searchString)
		afo.GormQuerySet = gormOperatorContext.Tx
		gormOperatorContext = interfaces.NewGormOperatorContext(afo.PaginatedGormQuerySet, afo.Model)
		operator.Build(afo.UadminDatabase.Adapter, gormOperatorContext, sf.Field, searchString)
		afo.PaginatedGormQuerySet = gormOperatorContext.Tx
	}
}

type PaginationType string

var LimitPaginationType PaginationType = "limit"
var CursorPaginationType PaginationType = "cursor"

type IPaginationInterface interface {
	Paginate (afo *AdminFilterObjects)
}

type Paginator struct {
	IPaginationInterface
	PerPage int
	AllowEmptyFirstPage bool
	ShowLastPageOnPreviousPage bool
	Count int
	NumPages int
	Offset int
	Template string
	PaginationType PaginationType
}

func (p *Paginator) Paginate(afo *AdminFilterObjects) {

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

func (apr *AdminPageRegistry) PreparePagesForTemplate(permRegistry *interfaces.UserPermRegistry) []byte {
	pages := make([]*AdminPage, 0)

	for page := range apr.GetAll() {
		blueprintName := page.BlueprintName
		modelName := page.ModelName
		var userPerm *interfaces.UserPerm
		if modelName != "" {
			userPerm = permRegistry.GetPermissionForBlueprint(blueprintName, modelName)
			if !userPerm.HasReadPermission() {
				continue
			}
		} else {
			existsAnyPermission := permRegistry.IsThereAnyPermissionForBlueprint(blueprintName)
			if !existsAnyPermission {
				continue
			}
		}
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
	ListHandler func (ctx *gin.Context)
}

func (dap *DashboardAdminPanel) RegisterHttpHandlers(router *gin.Engine) {
	if dap.ListHandler != nil {
		router.GET(interfaces.CurrentConfig.D.Uadmin.RootAdminURL, dap.ListHandler)
	}
	for adminPage := range dap.AdminPages.GetAll() {
		router.GET(fmt.Sprintf("%s/%s", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug), func(pageTitle string, adminPageRegistry *AdminPageRegistry) func (ctx *gin.Context) {
			return func(ctx *gin.Context) {
				type Context struct {
					templatecontext.AdminContext
					Menu string
					CurrentPath string
				}

				c := &Context{}
				templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
				menu := string(adminPageRegistry.PreparePagesForTemplate(c.UserPermissionRegistry))
				c.Menu = menu
				c.CurrentPath = ctx.Request.URL.Path
				tr := interfaces.NewTemplateRenderer(pageTitle)
				tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("home"), c, template.FuncMap)
			}
		}(adminPage.PageName, adminPage.SubPages))
		for subPage := range adminPage.SubPages.GetAll() {
			router.GET(fmt.Sprintf("%s/%s/%s", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.Slug), func(adminPage *AdminPage) func (ctx *gin.Context) {
				return func(ctx *gin.Context) {
					if adminPage.ListHandler != nil {
						adminPage.ListHandler(ctx)
					} else {
						type Context struct {
							templatecontext.AdminContext
							AdminFilterObjects *AdminFilterObjects
							ListDisplay *ListDisplayRegistry
							PermissionForBlueprint *interfaces.UserPerm
							ListFilter *ListFilterRegistry
							InitialOrder string
							InitialOrderList []string
							Search string
							TotalRecords int64
						}

						c := &Context{}
						adminRequestParams := interfaces.NewAdminRequestParamsFromGinContext(ctx)
						templatecontext.PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
						c.PermissionForBlueprint = c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
						c.AdminFilterObjects = adminPage.GetQueryset(adminPage, adminRequestParams)
						c.AdminFilterObjects.GormQuerySet.Count(&c.TotalRecords)
						c.ListDisplay = adminPage.ListDisplay
						c.Search = adminRequestParams.Search
						c.ListFilter = adminPage.ListFilter
						c.InitialOrder = adminRequestParams.GetOrdering()
						c.InitialOrderList = adminRequestParams.Ordering
						tr := interfaces.NewTemplateRenderer(adminPage.PageName)
						tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("list"), c, template.FuncMap)
					}
				}
			}(subPage))
		}
		for adminModelAction := range adminPage.ModelActionsRegistry.GetAllModelActions() {
			router.POST(fmt.Sprintf("%s/%s/%s/%s/", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, adminPage.ModelName, adminModelAction.SlugifiedActionName), func(adminPage *AdminPage, slugifiedModelActionName string) func (ctx *gin.Context) {
				return func(ctx *gin.Context) {
					adminPage.HandleModelAction(slugifiedModelActionName, ctx)
				}
			}(adminPage, adminModelAction.SlugifiedActionName))

		}
	}
}

var CurrentDashboardAdminPanel *DashboardAdminPanel

func NewDashboardAdminPanel() *DashboardAdminPanel {
	return &DashboardAdminPanel{
		AdminPages: NewAdminPageRegistry(),
	}
}

type DisplayFilterOption struct {
	FilterField string
	FilterValue string
	DisplayAs string
}

type FilterOption struct {
	FieldName string
	FetchOptions func(afo *AdminFilterObjects) []*DisplayFilterOption
}

type FilterOptionsRegistry struct {
	FilterOption []*FilterOption
}

func (for1 *FilterOptionsRegistry) AddFilterOption(fo *FilterOption) {
	for1.FilterOption = append(for1.FilterOption, fo)
}

func (for1 *FilterOptionsRegistry) GetAll() <- chan *FilterOption {
	chnl := make(chan *FilterOption)
	go func() {
		defer close(chnl)
		for _, fo := range for1.FilterOption {
			chnl <- fo
		}
	}()
	return chnl
}

func NewFilterOptionsRegistry() *FilterOptionsRegistry {
	return &FilterOptionsRegistry{FilterOption: make([]*FilterOption, 0)}
}

func NewFilterOption() *FilterOption {
	return &FilterOption{}
}

func FetchOptionsFromGormModelFromDateTimeField(afo *AdminFilterObjects, filterOptionField string) []*DisplayFilterOption {
	ret := make([]*DisplayFilterOption, 0)
	uadminDatabase := interfaces.NewUadminDatabase()
	defer uadminDatabase.Close()
	filterString := uadminDatabase.Adapter.GetStringToExtractYearFromField(filterOptionField)
	rows, _ := afo.InitialGormQuerySet.Select(filterString + " as year, count(*) as total").Group(filterString).Rows()
	var filterValue uint
	var filterCount uint
	for rows.Next() {
		rows.Scan(&filterValue, &filterCount)
		filterString := strconv.Itoa(int(filterValue))
		ret = append(ret, &DisplayFilterOption{
			FilterField: filterOptionField,
			FilterValue: filterString,
			DisplayAs: filterString,
		})
	}
	if len(ret) < 2 {
		ret = make([]*DisplayFilterOption, 0)
		filterString := uadminDatabase.Adapter.GetStringToExtractMonthFromField(filterOptionField)
		rows, _ := afo.InitialGormQuerySet.Select(filterString + " as month, count(*) as total").Group(filterString).Rows()
		var filterValue uint
		var filterCount uint
		for rows.Next() {
			rows.Scan(&filterValue, &filterCount)
			filterString := strconv.Itoa(int(filterValue))
			filteredMonth, _ := strconv.Atoi(filterString)
			ret = append(ret, &DisplayFilterOption{
				FilterField: filterOptionField,
				FilterValue: filterString,
				DisplayAs: time.Month(filteredMonth).String(),
			})
		}
	}
	return ret
}

type AdminPage struct {
	Model interface{} `json:"-"`
	GetQueryset func(*AdminPage, *interfaces.AdminRequestParams) *AdminFilterObjects `json:"-"`
	ModelActionsRegistry *AdminModelActionRegistry `json:"-"`
	FilterOptions *FilterOptionsRegistry `json:"-"`
	ActionsSelectionCounter bool `json:"-"`
	BlueprintName string
	EmptyValueDisplay string `json:"-"`
	ExcludeFields interfaces.IFieldRegistry `json:"-"`
	FieldsToShow interfaces.IFieldRegistry `json:"-"`
	Form *form.Form `json:"-"`
	ShowAllFields bool `json:"-"`
	Validators []interfaces.IValidator `json:"-"`
	Inlines []*AdminPageInlines `json:"-"`
	ListDisplay *ListDisplayRegistry `json:"-"`
	ListFilter *ListFilterRegistry `json:"-"`
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
	ModelName string
	Slug string
	ToolTip string
	Icon string
	ListHandler func (ctx *gin.Context) `json:"-"`
	EditHandler func (ctx *gin.Context) `json:"-"`
	AddHandler func (ctx *gin.Context) `json:"-"`
	DeleteHandler func (ctx *gin.Context) `json:"-"`
	Router *gin.Engine `json:"-"`
}

type ModelActionRequestParams struct {
	ObjectIds []uint `form:"object_ids" json:"object_ids" xml:"object_ids"  binding:"required"`
}

func (ap *AdminPage) HandleModelAction(modelActionName string, ctx *gin.Context) {
	afo := ap.GetQueryset(ap, nil)
	var json ModelActionRequestParams
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(json.ObjectIds) > 0 {
		primaryKeyField, _ := ap.Form.FieldRegistry.GetPrimaryKey()
		afo.GormQuerySet = afo.GormQuerySet.Where(fmt.Sprintf("%s IN ?", primaryKeyField.DBName), json.ObjectIds)
		modelAction, _ := ap.ModelActionsRegistry.GetModelActionByName(modelActionName)
		_, affectedRows := modelAction.Handler(afo)
		ctx.JSON(http.StatusOK, gin.H{"Affected": strconv.Itoa(int(affectedRows))})
	} else {
		ctx.Status(400)
	}
}

func (ap *AdminPage) FetchFilterOptions() []*DisplayFilterOption {
	afo := ap.GetQueryset(ap, nil)
	filterOptions := make([]*DisplayFilterOption, 0)
	for filterOption := range ap.FilterOptions.GetAll() {
		filterOptions = append(filterOptions, filterOption.FetchOptions(afo)...)
	}
	return filterOptions
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
	Classes []string
	Extra int
	MaxNum int
	MinNum int
	VerboseName string
	VerboseNamePlural string
	ShowChangeLink bool
	ConnectionToParentModel ConnectionToParentModel
	Template string
	Permissions *interfaces.UserPerm
}

func NewAdminPageRegistry() *AdminPageRegistry {
	return &AdminPageRegistry{AdminPages: make(map[string]*AdminPage)}
}

func NewGormAdminPage(genModelI func() interface{}, modelName string) *AdminPage {
	return &AdminPage{
		SubPages: NewAdminPageRegistry(),
		GetQueryset: func(adminPage *AdminPage, adminRequestParams *interfaces.AdminRequestParams) *AdminFilterObjects {
			uadminDatabase := interfaces.NewUadminDatabase()
			db := uadminDatabase.Db
			var paginatedQuerySet *gorm.DB
			var perPage int
			ret := &AdminFilterObjects{
				InitialGormQuerySet: db.Model(genModelI()),
				GormQuerySet: db.Model(genModelI()),
				PaginatedGormQuerySet: db.Model(genModelI()),
				Model: genModelI(),
				UadminDatabase: uadminDatabase,
				GenerateModelI: genModelI,
			}
			if adminRequestParams != nil && adminRequestParams.RequestURL != "" {
				url1, _ := url.Parse(adminRequestParams.RequestURL)
				queryParams, _ := url.ParseQuery(url1.RawQuery)
				for filter := range adminPage.ListFilter.Iterate() {
					filterValue := queryParams.Get(filter.UrlFilteringParam)
					if filterValue != "" {
						filter.FilterQs(ret, fmt.Sprintf("%s=%s", filter.UrlFilteringParam, filterValue))
					}
				}
			}
			if adminRequestParams != nil && adminRequestParams.Search != "" {
				for _, filter := range adminPage.SearchFields {
					filter.Search(ret, adminRequestParams.Search)
				}
			}
			if adminRequestParams != nil && adminRequestParams.Paginator.PerPage > 0 {
				perPage = adminRequestParams.Paginator.PerPage
			} else {
				perPage = adminPage.Paginator.PerPage
			}
			if adminRequestParams != nil {
				paginatedQuerySet = ret.PaginatedGormQuerySet.Offset(adminRequestParams.Paginator.Offset)
				if adminPage.Paginator.ShowLastPageOnPreviousPage {
					var countRecords int64
					ret.GormQuerySet.Count(&countRecords)
					if countRecords > int64(adminRequestParams.Paginator.Offset+(2*perPage)) {
						paginatedQuerySet = paginatedQuerySet.Limit(perPage)
					} else {
						paginatedQuerySet = paginatedQuerySet.Limit(int(countRecords - int64(adminRequestParams.Paginator.Offset)))
					}
				} else {
					paginatedQuerySet = paginatedQuerySet.Limit(perPage)
				}
				ret.PaginatedGormQuerySet = paginatedQuerySet
				for listDisplay := range adminPage.ListDisplay.GetAllFields() {
					direction := listDisplay.SortBy.Direction
					if len(adminRequestParams.Ordering) > 0 {
						for _, ordering := range adminRequestParams.Ordering {
							directionSort := 1
							if strings.HasPrefix(ordering, "-") {
								directionSort = -1
								ordering = ordering[1:]
							}
							if ordering == listDisplay.DisplayName {
								direction = directionSort
								listDisplay.SortBy.Sort(ret, direction)
							}
						}
					}
				}
			}
			return ret
		},
		Model: genModelI(),
		ModelName: modelName,
		Validators: make([]interfaces.IValidator, 0),
		ExcludeFields: interfaces.NewFieldRegistry(),
		FieldsToShow: interfaces.NewFieldRegistry(),
		ModelActionsRegistry: NewAdminModelActionRegistry(),
		Inlines: make([]*AdminPageInlines, 0),
		ListDisplay: &ListDisplayRegistry{ListDisplayFields: make(map[string]*ListDisplay)},
		ListFilter: &ListFilterRegistry{ListFilter: make([]*ListFilter, 0)},
		SearchFields: make([]*SearchField, 0),
		ExtraStatic: &StaticFiles{
			ExtraCSS: make([]string, 0),
			ExtraJS: make([]string, 0),
		},
		Paginator: &Paginator{PerPage: interfaces.CurrentConfig.D.Uadmin.AdminPerPage, ShowLastPageOnPreviousPage: true},
		ActionsSelectionCounter: true,
		FilterOptions: NewFilterOptionsRegistry(),
	}
}
