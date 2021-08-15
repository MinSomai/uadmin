package interfaces

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AdminRequestPaginator struct {
	PerPage int
	Offset int
}

type AdminRequestParams struct {
	CreateSession bool
	GenerateCSRFToken bool
	NeedAllLanguages bool
	Paginator *AdminRequestPaginator
	RequestURL string
	Search string
	Ordering []string
}

func (arp *AdminRequestParams) GetOrdering() string {
	return strings.Join(arp.Ordering, ",")
}

func NewAdminRequestParams() *AdminRequestParams {
	return &AdminRequestParams{
		CreateSession: true,
		GenerateCSRFToken: true,
		NeedAllLanguages: false,
		Paginator: &AdminRequestPaginator{},
	}
}

func NewAdminRequestParamsFromGinContext(ctx *gin.Context) *AdminRequestParams {
	ret := &AdminRequestParams{
		CreateSession: true,
		GenerateCSRFToken: true,
		NeedAllLanguages: false,
		Paginator: &AdminRequestPaginator{},
	}
	if ctx.Query("perpage") != "" {
		perPage, _ := strconv.Atoi(ctx.Query("perpage"))
		ret.Paginator.PerPage = perPage
 	} else {
		ret.Paginator.PerPage = CurrentConfig.D.Uadmin.AdminPerPage
	}
	if ctx.Query("offset") != "" {
		offset, _ := strconv.Atoi(ctx.Query("offset"))
		ret.Paginator.Offset = offset
	}
	if ctx.Query("p") != "" {
		page, _ := strconv.Atoi(ctx.Query("p"))
		if page > 1 {
			ret.Paginator.Offset = (page - 1) * ret.Paginator.PerPage
		}
	}
	ret.RequestURL = ctx.Request.URL.String()
	ret.Search = ctx.Query("search")
	orderingParts := strings.Split(ctx.Query("initialOrder"), ",")
	currentOrder := ctx.Query("o")
	currentOrderNameWithoutDirection := currentOrder
	if strings.HasPrefix(currentOrderNameWithoutDirection, "-") {
		currentOrderNameWithoutDirection = currentOrderNameWithoutDirection[1:]
	}
	foundNewOrder := false
	for i, part := range orderingParts {
		if strings.HasPrefix(part, "-") {
			part = part[1:]
		}
		if part == currentOrderNameWithoutDirection {
			orderingParts[i] = currentOrder
			foundNewOrder = true
		}
	}
	if !foundNewOrder {
		orderingParts = append(orderingParts, currentOrder)
	}
	finalOrderingParts := make([]string, 0)
	for _, part := range orderingParts {
		if part != "" {
			finalOrderingParts = append(finalOrderingParts, part)
		}
	}
	ret.Ordering = finalOrderingParts
	return ret
}

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
	PermName CustomPermission
	Handler func (adminPage *AdminPage, afo *AdminFilterObjects, ctx *gin.Context) (bool, int64)
	IsDisabled func (afo *AdminFilterObjects, ctx *gin.Context) bool
	SlugifiedActionName string
	RequestMethod string
	RequiresExtraSteps bool
}

func prepareAdminModelActionName(adminModelAction string) string {
	slugifiedAdminModelAction := AsciiRegex.ReplaceAllLiteralString(adminModelAction, "")
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
		RequestMethod: "POST",
	}
}

type AdminFilterObjects struct {
	InitialGormQuerySet *gorm.DB
	GormQuerySet *gorm.DB
	PaginatedGormQuerySet *gorm.DB
	Model interface{}
	UadminDatabase *UadminDatabase
	GenerateModelI func() (interface{}, interface{})
}

type IterateAdminObjects struct {
	Model interface {}
	Id uint
}

func (afo *AdminFilterObjects) WithTransaction(handler func(afo1 *AdminFilterObjects) error) {
	afo.UadminDatabase.Db.Transaction(func(tx *gorm.DB) error {
		return handler(&AdminFilterObjects{UadminDatabase: &UadminDatabase{Db: tx}, GenerateModelI: afo.GenerateModelI})
	})
}

func (afo *AdminFilterObjects) LoadDataForModelById(Id interface{}, model interface{}) {
	modelI, _ := afo.GenerateModelI()
	afo.UadminDatabase.Db.Model(modelI).First(model, Id)
}

func (afo *AdminFilterObjects) SaveModel(model interface{}) error {
	res := afo.UadminDatabase.Db.Save(model)
	return res.Error
}

func (afo *AdminFilterObjects) GetPaginated() <- chan *IterateAdminObjects {
	chnl := make(chan *IterateAdminObjects)
	go func() {
		defer close(chnl)
		modelI, models := afo.GenerateModelI()
		modelDescription := ProjectModels.GetModelFromInterface(modelI)
		afo.PaginatedGormQuerySet.Preload(clause.Associations).Find(models)
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			gormModelV := reflect.Indirect(reflect.ValueOf(model))
			Id := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
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

func (afo *AdminFilterObjects) IterateThroughModelActionsSelected() <- chan *IterateAdminObjects {
	chnl := make(chan *IterateAdminObjects)
	go func() {
		defer close(chnl)
		modelI, models := afo.GenerateModelI()
		modelDescription := ProjectModels.GetModelFromInterface(modelI)
		afo.GormQuerySet.Preload(clause.Associations).Find(models)
		s := reflect.Indirect(reflect.ValueOf(models))
		for i := 0; i < s.Len(); i++ {
			model := s.Index(i).Interface()
			gormModelV := reflect.Indirect(reflect.ValueOf(model))
			Id := TransformValueForWidget(gormModelV.FieldByName(modelDescription.Statement.Schema.PrimaryFields[0].Name).Interface())
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
	Field *Field
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
	MaxOrdering int
}

func (ldr *ListDisplayRegistry) ClearAllFields() {
	ldr.MaxOrdering = 0
	ldr.ListDisplayFields = make(map[string]*ListDisplay)
}

func (ldr *ListDisplayRegistry) IsThereAnyEditable() bool {
	for ld := range ldr.GetAllFields() {
		if ld.IsEditable {
			return true
		}
	}
	return false
}
func (ldr *ListDisplayRegistry) AddField(ld *ListDisplay) {
	ldr.ListDisplayFields[ld.DisplayName] = ld
	ldr.MaxOrdering = int(math.Max(float64(ldr.MaxOrdering + 1), float64(ld.Ordering + 1)))
	ld.Ordering = ldr.MaxOrdering
}

func (ldr *ListDisplayRegistry) BuildFormForListEditable(ID uint, model interface{}) *FormListEditable {
	return NewFormListEditableFromListDisplayRegistry(ID, model, ldr)
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

type AdminModelActionRegistry struct {
	AdminModelActions map[string]*AdminModelAction
}

func (amar *AdminModelActionRegistry) AddModelAction(ma *AdminModelAction) {
	amar.AdminModelActions[ma.SlugifiedActionName] = ma
}

func (amar *AdminModelActionRegistry) IsThereAnyActions() bool {
	return len(amar.AdminModelActions) > 0
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
	Field *Field
	ChangeLink bool
	Ordering int
	SortBy *SortBy
	Populate func (m interface{}) string
	MethodName string
	IsEditable bool
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
	if reflect.ValueOf(m).IsZero() || gormModelV.IsZero() || gormModelV.FieldByName(ld.Field.Name).IsZero() {
		return ""
	}
	return TransformValueForListDisplay(gormModelV.FieldByName(ld.Field.Name).Interface())
}

func NewListDisplay(field *Field) *ListDisplay {
	displayName := ""
	if field != nil {
		displayName = field.DisplayName
	}
	return &ListDisplay{
		DisplayName: displayName, Field: field, ChangeLink: true,
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
	OptionsToShow []*FieldChoice
	FetchOptions func(m interface{}) []*FieldChoice
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
		operatorContext := FilterGormModel(afo.UadminDatabase.Adapter, afo.GormQuerySet, schema1, []string{filterString}, afo.Model)
		afo.GormQuerySet = operatorContext.Tx
		operatorContext = FilterGormModel(afo.UadminDatabase.Adapter, afo.PaginatedGormQuerySet, schema1, []string{filterString}, afo.Model)
		afo.PaginatedGormQuerySet = operatorContext.Tx
	}
}

func (lf *ListFilter) IsItActive (fullUrl *url.URL) bool {
	return strings.Contains(fullUrl.String(), lf.UrlFilteringParam)
}

func (lf *ListFilter) GetURLToClearFilter(fullUrl *url.URL) string {
	clonedUrl := CloneNetUrl(fullUrl)
	qs := clonedUrl.Query()
	qs.Del(lf.UrlFilteringParam)
	clonedUrl.RawQuery = qs.Encode()
	return clonedUrl.String()
}

func (lf *ListFilter) IsThatOptionActive(option *FieldChoice, fullUrl *url.URL) bool {
	qs := fullUrl.Query()
	value := qs.Get(lf.UrlFilteringParam)
	if value != "" {
		optionValue := TransformValueForListDisplay(option.Value)
		if optionValue == value {
			return true
		}
	}
	return false
}

func (lf *ListFilter) GetURLForOption(option *FieldChoice, fullUrl *url.URL) string {
	clonedUrl := CloneNetUrl(fullUrl)
	qs := clonedUrl.Query()
	qs.Set(lf.UrlFilteringParam, TransformValueForListDisplay(option.Value))
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
		operator := ExactGormOperator{}
		gormOperatorContext := NewGormOperatorContext(afo.GormQuerySet, afo.Model)
		operator.Build(afo.UadminDatabase.Adapter, gormOperatorContext, sf.Field, searchString)
		afo.GormQuerySet = gormOperatorContext.Tx
		gormOperatorContext = NewGormOperatorContext(afo.PaginatedGormQuerySet, afo.Model)
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
	uadminDatabase := NewUadminDatabase()
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

var CurrentAdminPageRegistry *AdminPageRegistry
