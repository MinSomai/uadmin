package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func NewGormAdminPage(parentPage *AdminPage, genModelI func() (interface{}, interface{}), generateForm func(modelI interface{}, ctx IAdminContext) *Form) *AdminPage {
	modelI4, _ := genModelI()
	modelName := ""
	if modelI4 != nil {
		uadminDatabase := NewUadminDatabaseWithoutConnection()
		stmt := &gorm.Statement{DB: uadminDatabase.Db}
		stmt.Parse(modelI4)
		modelName = strings.ToLower(stmt.Schema.Name)
	}
	var form *Form
	var listDisplay *ListDisplayRegistry
	var searchFieldRegistry *SearchFieldRegistry
	if modelI4 != nil {
		form = NewFormFromModelFromGinContext(&AdminContext{}, modelI4, make([]string, 0), []string{"ID"}, true, "")
		listDisplay = NewListDisplayRegistryFromGormModel(modelI4)
		searchFieldRegistry = NewSearchFieldRegistryFromGormModel(modelI4)
	}
	return &AdminPage{
		Form:           form,
		SubPages:       NewAdminPageRegistry(),
		GenerateModelI: genModelI,
		ParentPage:     parentPage,
		GetQueryset: func(adminPage *AdminPage, adminRequestParams *AdminRequestParams) IAdminFilterObjects {
			uadminDatabase := NewUadminDatabase()
			db := uadminDatabase.Db
			var paginatedQuerySet IPersistenceStorage
			var perPage int
			modelI, _ := genModelI()
			modelI1, _ := genModelI()
			modelI2, _ := genModelI()
			modelI3, _ := genModelI()
			ret := &GormAdminFilterObjects{
				InitialGormQuerySet:   NewGormPersistenceStorage(db.Model(modelI)),
				GormQuerySet:          NewGormPersistenceStorage(db.Model(modelI1)),
				PaginatedGormQuerySet: NewGormPersistenceStorage(db.Model(modelI2)),
				Model:                 modelI3,
				UadminDatabase:        uadminDatabase,
				GenerateModelI:        genModelI,
			}
			if adminRequestParams != nil && adminRequestParams.RequestURL != "" {
				url1, _ := url.Parse(adminRequestParams.RequestURL)
				queryParams, _ := url.ParseQuery(url1.RawQuery)
				for filter := range adminPage.ListFilter.Iterate() {
					filterValue := queryParams.Get(filter.URLFilteringParam)
					if filterValue != "" {
						filter.FilterQs(ret, fmt.Sprintf("%s=%s", filter.URLFilteringParam, filterValue))
					}
				}
			}
			if adminRequestParams != nil && adminRequestParams.Search != "" {
				searchFilterObjects := &GormAdminFilterObjects{
					InitialGormQuerySet:   NewGormPersistenceStorage(db),
					GormQuerySet:          NewGormPersistenceStorage(db),
					PaginatedGormQuerySet: NewGormPersistenceStorage(db),
					Model:                 modelI3,
					UadminDatabase:        uadminDatabase,
					GenerateModelI:        genModelI,
				}
				for filter := range adminPage.SearchFields.GetAll() {
					filter.Search(searchFilterObjects, adminRequestParams.Search)
				}
				ret.SetPaginatedQuerySet(ret.GetPaginatedQuerySet().Where(searchFilterObjects.GetPaginatedQuerySet().GetCurrentDB()))
				ret.SetFullQuerySet(ret.GetFullQuerySet().Where(searchFilterObjects.GetFullQuerySet().GetCurrentDB()))
			}
			if adminRequestParams != nil && adminRequestParams.Paginator.PerPage > 0 {
				perPage = adminRequestParams.Paginator.PerPage
			} else {
				perPage = adminPage.Paginator.PerPage
			}
			if adminRequestParams != nil {
				paginatedQuerySet = ret.GetPaginatedQuerySet().Offset(adminRequestParams.Paginator.Offset)
				if adminPage.Paginator.ShowLastPageOnPreviousPage {
					var countRecords int64
					ret.GetFullQuerySet().Count(&countRecords)
					if countRecords > int64(adminRequestParams.Paginator.Offset+(2*perPage)) {
						paginatedQuerySet = paginatedQuerySet.Limit(perPage)
					} else {
						paginatedQuerySet = paginatedQuerySet.Limit(int(countRecords - int64(adminRequestParams.Paginator.Offset)))
					}
				} else {
					paginatedQuerySet = paginatedQuerySet.Limit(perPage)
				}
				ret.SetPaginatedQuerySet(paginatedQuerySet)
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
		Model:                   modelI4,
		ModelName:               modelName,
		Validators:              NewValidatorRegistry(),
		ExcludeFields:           NewFieldRegistry(),
		FieldsToShow:            NewFieldRegistry(),
		ModelActionsRegistry:    NewAdminModelActionRegistry(),
		InlineRegistry:          NewAdminPageInlineRegistry(),
		ListDisplay:             listDisplay,
		ListFilter:              &ListFilterRegistry{ListFilter: make([]*ListFilter, 0)},
		SearchFields:            searchFieldRegistry,
		Paginator:               &Paginator{PerPage: CurrentConfig.D.Uadmin.AdminPerPage, ShowLastPageOnPreviousPage: true},
		ActionsSelectionCounter: true,
		FilterOptions:           NewFilterOptionsRegistry(),
		GenerateForm:            generateForm,
	}
}

var CurrentAdminPageRegistry *AdminPageRegistry

type AdminPagesList []*AdminPage

func (apl AdminPagesList) Len() int { return len(apl) }
func (apl AdminPagesList) Less(i, j int) bool {
	if apl[i].Ordering == apl[j].Ordering {
		return apl[i].PageName < apl[j].PageName
	}
	return apl[i].Ordering < apl[j].Ordering
}
func (apl AdminPagesList) Swap(i, j int) { apl[i], apl[j] = apl[j], apl[i] }

type AdminPageRegistry struct {
	AdminPages map[string]*AdminPage
}

func (apr *AdminPageRegistry) GetByModelName(modelName string) *AdminPage {
	for adminPage := range apr.GetAll() {
		for subPage := range adminPage.SubPages.GetAll() {
			projectModel := ProjectModels.GetModelFromInterface(subPage.Model)
			// spew.Dump("AAAAAAAAAA", projectModel.Statement.Schema.Name)
			if projectModel.Statement.Schema.Name == modelName {
				return subPage
			}
			for subPage1 := range subPage.SubPages.GetAll() {
				projectModel1 := ProjectModels.GetModelFromInterface(subPage1.Model)
				// spew.Dump("AAAAAAAAAA1111111111", projectModel1.Statement.Schema.Name)
				if projectModel1.Statement.Schema.Name == modelName {
					return subPage1
				}
			}
		}
		if adminPage.Model == nil {
			continue
		}
		projectModel := ProjectModels.GetModelFromInterface(adminPage.Model)
		// spew.Dump("AAAAAAAAAA222", projectModel.Statement.Schema.Name)
		if projectModel.Statement.Schema.Name == modelName {
			return adminPage
		}
	}
	return nil
}

func (apr *AdminPageRegistry) GetBySlug(slug string) (*AdminPage, error) {
	adminPage, ok := apr.AdminPages[slug]
	if !ok {
		return nil, fmt.Errorf("No admin page with alias %s", slug)
	}
	return adminPage, nil
}

func (apr *AdminPageRegistry) AddAdminPage(adminPage *AdminPage) error {
	apr.AdminPages[adminPage.Slug] = adminPage
	return nil
}

func (apr *AdminPageRegistry) GetAll() <-chan *AdminPage {
	chnl := make(chan *AdminPage)
	go func() {
		defer close(chnl)
		sortedPages := make(AdminPagesList, 0)

		for _, adminPage := range apr.AdminPages {
			sortedPages = append(sortedPages, adminPage)
		}
		sort.Slice(sortedPages, func(i int, j int) bool {
			if sortedPages[i].Ordering == sortedPages[j].Ordering {
				return sortedPages[i].PageName < sortedPages[j].PageName
			}
			return sortedPages[i].Ordering < sortedPages[j].Ordering

		})
		for _, adminPage := range sortedPages {
			chnl <- adminPage
		}

	}()
	return chnl
}

func (apr *AdminPageRegistry) PreparePagesForTemplate(permRegistry *UserPermRegistry) []byte {
	pages := make([]*AdminPage, 0)

	for page := range apr.GetAll() {
		blueprintName := page.BlueprintName
		modelName := page.ModelName
		var userPerm *UserPerm
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
		Trail(CRITICAL, "error while generating menu in admin", err)
	}
	return ret
}

type AdminPage struct {
	Model                              interface{}                                               `json:"-"`
	GenerateModelI                     func() (interface{}, interface{})                         `json:"-"`
	GenerateForm                       func(modelI interface{}, ctx IAdminContext) *Form         `json:"-"`
	GetQueryset                        func(*AdminPage, *AdminRequestParams) IAdminFilterObjects `json:"-"`
	ModelActionsRegistry               *AdminModelActionRegistry                                 `json:"-"`
	FilterOptions                      *FilterOptionsRegistry                                    `json:"-"`
	ActionsSelectionCounter            bool                                                      `json:"-"`
	BlueprintName                      string
	EmptyValueDisplay                  string                   `json:"-"`
	ExcludeFields                      IFieldRegistry           `json:"-"`
	FieldsToShow                       IFieldRegistry           `json:"-"`
	Form                               *Form                    `json:"-"`
	ShowAllFields                      bool                     `json:"-"`
	Validators                         *ValidatorRegistry       `json:"-"`
	InlineRegistry                     *AdminPageInlineRegistry `json:"-"`
	ListDisplay                        *ListDisplayRegistry     `json:"-"`
	ListFilter                         *ListFilterRegistry      `json:"-"`
	MaxShowAll                         int                      `json:"-"`
	PreserveFilters                    bool                     `json:"-"`
	SaveAndContinue                    bool                     `json:"-"`
	SaveOnTop                          bool                     `json:"-"`
	SearchFields                       *SearchFieldRegistry     `json:"-"`
	ShowFullResultCount                bool                     `json:"-"`
	ViewOnSite                         bool                     `json:"-"`
	ListTemplate                       string                   `json:"-"`
	AddTemplate                        string                   `json:"-"`
	EditTemplate                       string                   `json:"-"`
	DeleteConfirmationTemplate         string                   `json:"-"`
	DeleteSelectedConfirmationTemplate string                   `json:"-"`
	ObjectHistoryTemplate              string                   `json:"-"`
	PopupResponseTemplate              string                   `json:"-"`
	Paginator                          *Paginator               `json:"-"`
	SubPages                           *AdminPageRegistry       `json:"-"`
	Ordering                           int
	PageName                           string
	ModelName                          string
	Slug                               string
	ToolTip                            string
	Icon                               string
	ListHandler                        func(ctx *gin.Context)                                                 `json:"-"`
	EditHandler                        func(ctx *gin.Context)                                                 `json:"-"`
	AddHandler                         func(ctx *gin.Context)                                                 `json:"-"`
	DeleteHandler                      func(ctx *gin.Context)                                                 `json:"-"`
	Router                             *gin.Engine                                                            `json:"-"`
	ParentPage                         *AdminPage                                                             `json:"-"`
	SaveModel                          func(modelI interface{}, ID uint, afo IAdminFilterObjects) interface{} `json:"-"`
	RegisteredHTTPHandlers             bool
	NoPermissionToAddNew               bool
	NoPermissionToEdit                 bool
}

type ModelActionRequestParams struct {
	ObjectIds     string `form:"object_ids" json:"object_ids" xml:"object_ids"  binding:"required"`
	RealObjectIds []uint
}

func (ap *AdminPage) GenerateLinkToEditModel(gormModelV reflect.Value) string {
	ID := GetID(gormModelV)
	return fmt.Sprintf("%s/%s/%s/edit/%d", CurrentConfig.D.Uadmin.RootAdminURL, ap.ParentPage.Slug, ap.Slug, ID)
}

func (ap *AdminPage) GenerateLinkToAddNewModel() string {
	return fmt.Sprintf("%s/%s/%s/edit/new?_to_field=id&_popup=1", CurrentConfig.D.Uadmin.RootAdminURL, ap.ParentPage.Slug, ap.Slug)
}

func (ap *AdminPage) HandleModelAction(modelActionName string, ctx *gin.Context) {
	afo := ap.GetQueryset(ap, nil)
	var json ModelActionRequestParams
	if ctx.GetHeader("Content-Type") == "application/json" {
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := ctx.ShouldBind(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	objectIds := strings.Split(json.ObjectIds, ",")
	objectUintIds := make([]uint, 0)
	for _, objectID := range objectIds {
		idV, err := strconv.Atoi(objectID)
		if err == nil {
			objectUintIds = append(objectUintIds, uint(idV))
		}
	}
	json.RealObjectIds = objectUintIds
	if len(json.RealObjectIds) > 0 {
		primaryKeyField, _ := ap.Form.FieldRegistry.GetPrimaryKey()
		afo.SetFullQuerySet(afo.GetFullQuerySet().Where(fmt.Sprintf("%s IN ?", primaryKeyField.DBName), json.RealObjectIds))
		modelAction, _ := ap.ModelActionsRegistry.GetModelActionByName(modelActionName)
		if ctx.GetHeader("Content-Type") == "application/json" {
			_, affectedRows := modelAction.Handler(ap, afo, ctx)
			ctx.JSON(http.StatusOK, gin.H{"Affected": strconv.Itoa(int(affectedRows))})
		} else {
			modelAction.Handler(ap, afo, ctx)
		}
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

func NewAdminPageRegistry() *AdminPageRegistry {
	return &AdminPageRegistry{
		AdminPages: make(map[string]*AdminPage),
	}
}

type AdminPageInlineRegistry struct {
	Inlines []*AdminPageInline
}

func (apir *AdminPageInlineRegistry) Add(pageInline *AdminPageInline) {
	apir.Inlines = append(apir.Inlines, pageInline)
}

func (apir *AdminPageInlineRegistry) GetAll() <-chan *AdminPageInline {
	chnl := make(chan *AdminPageInline)
	go func() {
		defer close(chnl)
		for _, inline := range apir.Inlines {
			chnl <- inline
		}
	}()
	return chnl
}