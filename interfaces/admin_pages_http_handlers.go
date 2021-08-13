package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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
	Model interface{} `json:"-"`
	GetQueryset func(*AdminPage, *AdminRequestParams) *AdminFilterObjects `json:"-"`
	ModelActionsRegistry *AdminModelActionRegistry `json:"-"`
	FilterOptions *FilterOptionsRegistry `json:"-"`
	ActionsSelectionCounter bool `json:"-"`
	BlueprintName string
	EmptyValueDisplay string `json:"-"`
	ExcludeFields IFieldRegistry `json:"-"`
	FieldsToShow IFieldRegistry `json:"-"`
	Form *Form `json:"-"`
	ShowAllFields bool `json:"-"`
	Validators []IValidator `json:"-"`
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
	ParentPage *AdminPage `json:"-"`
}

type ModelActionRequestParams struct {
	ObjectIds string `form:"object_ids" json:"object_ids" xml:"object_ids"  binding:"required"`
	RealObjectIds []uint
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
	objectIds := strings.Split(json.ObjectIds,",")
	objectUintIds := make([]uint, 0)
	for _, objectId := range objectIds {
		idV, err := strconv.Atoi(objectId)
		if err == nil {
			objectUintIds = append(objectUintIds, uint(idV))
		}
	}
	json.RealObjectIds = objectUintIds
	if len(json.RealObjectIds) > 0 {
		primaryKeyField, _ := ap.Form.FieldRegistry.GetPrimaryKey()
		afo.GormQuerySet = afo.GormQuerySet.Where(fmt.Sprintf("%s IN ?", primaryKeyField.DBName), json.RealObjectIds)
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
	Validators []IValidator
	Classes []string
	Extra int
	MaxNum int
	MinNum int
	VerboseName string
	VerboseNamePlural string
	ShowChangeLink bool
	ConnectionToParentModel ConnectionToParentModel
	Template string
	Permissions *UserPerm
}

func NewAdminPageRegistry() *AdminPageRegistry {
	return &AdminPageRegistry{AdminPages: make(map[string]*AdminPage)}
}
