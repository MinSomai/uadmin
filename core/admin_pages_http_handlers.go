package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

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

type ConnectionToParentModel struct {
	FieldNameToValue map[string]interface{}
}

type InlineType string

var TabularInline InlineType
var StackedInline InlineType

func init() {
	TabularInline = "tabular"
	StackedInline = "stacked"
}

type AdminPageInline struct {
	Ordering          int
	GenerateModelI    func(m interface{}) (interface{}, interface{})
	GetQueryset       func(afo IAdminFilterObjects, model interface{}, rp *AdminRequestParams) IAdminFilterObjects
	Actions           *AdminModelActionRegistry
	EmptyValueDisplay string
	ExcludeFields     IFieldRegistry
	FieldsToShow      IFieldRegistry
	ShowAllFields     bool
	Validators        *ValidatorRegistry
	Classes           []string
	Extra             int
	MaxNum            int
	MinNum            int
	VerboseName       string
	VerboseNamePlural string
	ShowChangeLink    bool
	Template          string
	ContentType       *ContentType
	Permission        CustomPermission
	InlineType        InlineType
	Prefix            string
	ListDisplay       *ListDisplayRegistry `json:"-"`
}

func (api *AdminPageInline) RenderExampleForm(adminContext IAdminContext) string {
	type Context struct {
		AdminContext
		AdminContextInitial IAdminContext
		Inline              *AdminPageInline
	}
	c := &Context{}
	c.AdminContextInitial = adminContext
	c.Inline = api
	templateRenderer := NewTemplateRenderer("")
	func1 := make(template.FuncMap)
	path := "admin/inlineexampleform"
	templateName := CurrentConfig.GetPathToTemplate(path)
	return templateRenderer.RenderAsString(
		CurrentConfig.TemplatesFS, templateName,
		c, FuncMap, func1,
	)
}

func (api *AdminPageInline) GetFormForExample(adminContext IAdminContext) *FormListEditable {
	modelI, _ := api.GenerateModelI(nil)
	return api.ListDisplay.BuildListEditableFormForNewModel(adminContext, "toreplacewithid", modelI)
}

func (api *AdminPageInline) GetFormIdenForNewItems() string {
	return fmt.Sprintf("example-%s", api.Prefix)
}

func (api *AdminPageInline) GetInlineID() string {
	return PrepareStringToBeUsedForHTMLID(api.VerboseNamePlural)
}

func (api *AdminPageInline) GetAll(model interface{}, rp *AdminRequestParams) <-chan *IterateAdminObjects {
	qs := api.GetQueryset(nil, model, rp)
	return qs.IterateThroughWholeQuerySet()
}

func (api *AdminPageInline) ProceedRequest(afo IAdminFilterObjects, ctx *gin.Context, f *multipart.Form, model interface{}, rp *AdminRequestParams, adminContext IAdminContext) (InlineFormListEditableCollection, error) {
	collection := make(InlineFormListEditableCollection)
	var firstEditableField *ListDisplay
	qs := api.GetQueryset(afo, model, rp)
	for ld := range api.ListDisplay.GetAllFields() {
		if ld.IsEditable {
			firstEditableField = ld
			break
		}
	}
	if firstEditableField == nil {
		return collection, nil
	}
	var form *FormListEditable
	err := false
	var removalError error
	for fieldName := range f.Value {
		if !strings.HasSuffix(fieldName, firstEditableField.Field.FieldConfig.Widget.GetHTMLInputName()) {
			continue
		}
		if !strings.HasPrefix(fieldName, firstEditableField.Prefix) {
			continue
		}
		if strings.Contains(fieldName, "toreplacewithid") {
			continue
		}
		removalError = nil
		form = nil
		inlineID := strings.TrimPrefix(fieldName, firstEditableField.Prefix+"-")
		inlineID = strings.TrimSuffix(inlineID, "-"+firstEditableField.Field.FieldConfig.Widget.GetHTMLInputName())
		realInlineID := strings.Split(inlineID, "_")
		modelI, _ := api.GenerateModelI(model)
		inlineIDToRemove := f.Value[firstEditableField.Prefix+"-"+"object_id-to-remove-"+realInlineID[0]]
		isNew := false
		if !strings.Contains(inlineID, "new") {
			IDI, _ := strconv.Atoi(realInlineID[0])
			qs.LoadDataForModelByID(uint(IDI), modelI)
			form = api.ListDisplay.BuildFormForListEditable(adminContext, uint(IDI), modelI)
			collection[realInlineID[0]] = form
			if len(inlineIDToRemove) > 0 {
				removalError = qs.RemoveModelPermanently(modelI)
			}
		} else {
			form = api.ListDisplay.BuildListEditableFormForNewModel(adminContext, realInlineID[0], modelI)
			collection[realInlineID[0]] = form
			isNew = true
		}
		if len(inlineIDToRemove) > 0 {
			if removalError != nil {
				form.FormError = &FormError{
					FieldError:    make(map[string]ValidationError),
					GeneralErrors: make(ValidationError, 0),
				}
				form.FormError.AddGeneralError(removalError)
			}
		} else {
			formError := form.ProceedRequest(f, modelI, ctx)
			if removalError != nil {
				formError.AddGeneralError(formError)
			}
			if !formError.IsEmpty() {
				err = true
			} else {
				if isNew {
					error1 := afo.CreateNew(modelI)
					if error1 != nil {
						formError.AddGeneralError(error1)
						err = true
					}
				} else {
					error1 := afo.SaveModel(modelI)
					if error1 != nil {
						formError.AddGeneralError(error1)
						err = true
					}
				}
			}
		}
	}
	if err {
		return collection, fmt.Errorf("error while validating inlines")
	}
	return collection, nil
}

func NewAdminPageInline(
	inlineIden string,
	inlineType InlineType,
	generateModelI func(m interface{}) (interface{}, interface{}),
	getQuerySet func(afo IAdminFilterObjects, model interface{}, rp *AdminRequestParams) IAdminFilterObjects,
) *AdminPageInline {
	modelI, _ := generateModelI(nil)
	ld := NewListDisplayRegistryFromGormModelForInlines(modelI)
	ld.SetPrefix(PrepareStringToBeUsedForHTMLID(inlineIden))
	ret := &AdminPageInline{
		Actions:           NewAdminModelActionRegistry(),
		ExcludeFields:     NewFieldRegistry(),
		FieldsToShow:      NewFieldRegistry(),
		Validators:        NewValidatorRegistry(),
		Classes:           make([]string, 0),
		InlineType:        inlineType,
		ListDisplay:       ld,
		GenerateModelI:    generateModelI,
		GetQueryset:       getQuerySet,
		VerboseNamePlural: inlineIden,
	}
	return ret
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

func NewAdminPageInlineRegistry() *AdminPageInlineRegistry {
	return &AdminPageInlineRegistry{
		Inlines: make([]*AdminPageInline, 0),
	}
}

type SearchFieldRegistry struct {
	Fields []*SearchField
}

func (sfr *SearchFieldRegistry) GetAll() <-chan *SearchField {
	chnl := make(chan *SearchField)
	go func() {
		defer close(chnl)
		for _, field := range sfr.Fields {
			chnl <- field
		}

	}()
	return chnl
}

func (sfr *SearchFieldRegistry) AddField(sf *SearchField) {
	sfr.Fields = append(sfr.Fields, sf)
}
