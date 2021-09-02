package core

import (
	"bytes"
	"database/sql"
	"fmt"
	excelize1 "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var PopulateTemplateContextForAdminPanel func(ctx *gin.Context, context IAdminContext, adminRequestParams *AdminRequestParams)

type DashboardAdminPanel struct {
	AdminPages  *AdminPageRegistry
	ListHandler func(ctx *gin.Context)
}

func (dap *DashboardAdminPanel) FindPageForGormModel(m interface{}) *AdminPage {
	mDescription := ProjectModels.GetModelFromInterface(m)
	for adminPage := range dap.AdminPages.GetAll() {
		for subPage := range adminPage.SubPages.GetAll() {
			modelDescription := ProjectModels.GetModelFromInterface(subPage.Model)
			if modelDescription.Statement.Table == mDescription.Statement.Table {
				return subPage
			}
		}
	}
	return nil
}

func (dap *DashboardAdminPanel) RegisterHTTPHandlers(router *gin.Engine) {
	if dap.ListHandler != nil {
		router.GET(CurrentConfig.D.Uadmin.RootAdminURL, dap.ListHandler)
	}
	for adminPage := range dap.AdminPages.GetAll() {
		router.GET(fmt.Sprintf("%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug), func(pageTitle string, adminPageRegistry *AdminPageRegistry) func(ctx *gin.Context) {
			return func(ctx *gin.Context) {
				type Context struct {
					AdminContext
					Menu        string
					CurrentPath string
				}

				c := &Context{}
				PopulateTemplateContextForAdminPanel(ctx, c, NewAdminRequestParams())
				menu := string(adminPageRegistry.PreparePagesForTemplate(c.UserPermissionRegistry))
				c.Menu = menu
				c.CurrentPath = ctx.Request.URL.Path
				tr := NewTemplateRenderer(pageTitle)
				tr.Render(ctx, CurrentConfig.TemplatesFS, CurrentConfig.GetPathToTemplate("home"), c, FuncMap)
			}
		}(adminPage.PageName, adminPage.SubPages))
		for subPage := range adminPage.SubPages.GetAll() {
			if subPage.RegisteredHTTPHandlers {
				continue
			}
			router.Any(fmt.Sprintf("%s/%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.Slug), func(adminPage *AdminPage) func(ctx *gin.Context) {
				return func(ctx *gin.Context) {
					if adminPage.ListHandler != nil {
						adminPage.ListHandler(ctx)
					} else {
						type Context struct {
							AdminContext
							AdminFilterObjects       IAdminFilterObjects
							ListDisplay              *ListDisplayRegistry
							PermissionForBlueprint   *UserPerm
							ListFilter               *ListFilterRegistry
							InitialOrder             string
							InitialOrderList         []string
							Search                   string
							TotalRecords             int64
							TotalPages               int64
							ListEditableFormError    bool
							AdminModelActionRegistry *AdminModelActionRegistry
							Message                  string
							CurrentAdminContext      IAdminContext
							NoPermissionToAddNew     bool
							NoPermissionToEdit       bool
						}

						c := &Context{}
						c.NoPermissionToAddNew = adminPage.NoPermissionToAddNew
						adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
						PopulateTemplateContextForAdminPanel(ctx, c, NewAdminRequestParams())
						c.Message = ctx.Query("message")
						c.NoPermissionToEdit = adminPage.NoPermissionToEdit
						c.PermissionForBlueprint = c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
						c.AdminFilterObjects = adminPage.GetQueryset(adminPage, adminRequestParams)
						c.AdminModelActionRegistry = adminPage.ModelActionsRegistry
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.BlueprintName, URL: fmt.Sprintf("%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.ParentPage.Slug)})
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.ModelName, IsActive: true})
						if ctx.Request.Method == "POST" {
							c.AdminFilterObjects.WithTransaction(func(afo1 IAdminFilterObjects) error {
								postForm, _ := ctx.MultipartForm()
								ids := postForm.Value["object_id"]
								for _, objectID := range ids {
									objectModel, _ := c.AdminFilterObjects.GenerateModelInterface()
									IDInt, _ := strconv.Atoi(objectID)
									IDUint := uint(IDInt)
									afo1.LoadDataForModelByID(IDUint, objectModel)
									modelI, _ := c.AdminFilterObjects.GenerateModelInterface()
									listEditableForm := NewFormListEditableFromListDisplayRegistry(c, "", IDUint, modelI, adminPage.ListDisplay)
									formListEditableErr := listEditableForm.ProceedRequest(postForm, objectModel, ctx)
									if formListEditableErr.IsEmpty() {
										dbRes := afo1.SaveModel(objectModel)
										if dbRes != nil {
											c.ListEditableFormError = true
											return dbRes
										}
									}
								}
								return nil
							})
						}
						c.AdminFilterObjects.GetFullQuerySet().Count(&c.TotalRecords)
						c.TotalPages = int64(math.Ceil(float64(c.TotalRecords / int64(adminPage.Paginator.PerPage))))
						c.ListDisplay = adminPage.ListDisplay
						c.Search = adminRequestParams.Search
						c.ListFilter = adminPage.ListFilter
						c.InitialOrder = adminRequestParams.GetOrdering()
						c.InitialOrderList = adminRequestParams.Ordering
						c.CurrentAdminContext = c
						tr := NewTemplateRenderer(adminPage.PageName)
						tr.Render(ctx, CurrentConfig.TemplatesFS, CurrentConfig.GetPathToTemplate("list"), c, FuncMap)
					}
				}
			}(subPage))
			router.POST(fmt.Sprintf("%s/%s/%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.Slug, "export"), func(adminPage *AdminPage) func(ctx *gin.Context) {
				return func(ctx *gin.Context) {
					type Context struct {
						AdminContext
					}
					c := &Context{}
					adminRequestParams := NewAdminRequestParamsFromGinContext(ctx)
					PopulateTemplateContextForAdminPanel(ctx, c, NewAdminRequestParams())
					// permissionForBlueprint := c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
					adminFilterObjects := adminPage.GetQueryset(adminPage, adminRequestParams)
					rows, _ := adminFilterObjects.GetFullQuerySet().Rows()
					defer rows.Close()
					db := NewUadminDatabase()
					defer db.Close()
					f := excelize1.NewFile()
					i := 1
					currentColumn := 'A'
					for listDisplay := range adminPage.ListDisplay.GetAllFields() {
						f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", currentColumn, i), listDisplay.DisplayName)
						currentColumn++
					}
					i++
					for rows.Next() {
						model, _ := adminFilterObjects.GenerateModelInterface()
						db.Db.ScanRows(rows.(*sql.Rows), model)
						// db.Db.ScanRows(rows, model)
						currentColumn = 'A'
						for listDisplay := range adminPage.ListDisplay.GetAllFields() {
							f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", currentColumn, i), listDisplay.GetValue(model, true))
							currentColumn++
						}
						i++
					}

					//f.SetCellValue("Sheet1", "B2", 100)
					//f.SetCellValue("Sheet1", "A1", 50)
					b, _ := f.WriteToBuffer()
					downloadName := adminPage.PageName + ".xlsx"
					ctx.Header("Content-Description", "File Transfer")
					ctx.Header("Content-Disposition", "attachment; filename="+downloadName)
					ctx.Data(http.StatusOK, "application/octet-stream", b.Bytes())
				}
			}(subPage))
			router.Any(fmt.Sprintf("%s/%s/%s/edit/:id", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.Slug), func(adminPage *AdminPage) func(ctx *gin.Context) {
				return func(ctx *gin.Context) {
					id := ctx.Param("id")
					type Context struct {
						AdminContext
						AdminModelActionRegistry    *AdminModelActionRegistry
						Message                     string
						PermissionForBlueprint      *UserPerm
						Form                        *Form
						Model                       interface{}
						ID                          uint
						IsNew                       bool
						ListURL                     string
						AdminPageInlineRegistry     *AdminPageInlineRegistry
						AdminRequestParams          *AdminRequestParams
						CurrentAdminContext         IAdminContext
						ListEditableFormsForInlines *FormListEditableCollection
					}

					c := &Context{}
					c.ListURL = fmt.Sprintf("%s/%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug)
					c.PageTitle = adminPage.ModelName
					c.CurrentAdminContext = c
					c.ListEditableFormsForInlines = NewFormListEditableCollection()
					modelI, _ := adminPage.GenerateModelI()
					if id != "new" {
						uadminDatabase := NewUadminDatabase()
						idI, _ := strconv.Atoi(id)
						uadminDatabase.Db.Preload(clause.Associations).First(modelI, idI)
						uadminDatabase.Close()
					}
					adminRequestParams := NewAdminRequestParams()
					c.AdminRequestParams = adminRequestParams
					PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)
					form := adminPage.GenerateForm(modelI, c)
					field, _ := form.FieldRegistry.GetByName("ID")
					ID, _ := field.FieldConfig.Widget.GetValue().(uint)
					c.ID = ID
					form.TemplateName = "admin/form_edit"
					form.RequestContext["ID"] = c.GetID()
					c.Model = modelI
					form.DontGenerateFormTag = true
					c.IsNew = true
					c.AdminPageInlineRegistry = adminPage.InlineRegistry
					form.ForAdminPanel = true
					if ctx.Request.Method == "POST" {
						requestForm, _ := ctx.MultipartForm()
						var modelToSave interface{}
						if id != "new" {
							modelToSave = modelI
						} else {
							modelToSave, _ = adminPage.GenerateModelI()
						}
						uadminDatabase := NewUadminDatabase()
						transactionerror := uadminDatabase.Db.Transaction(func(tx *gorm.DB) error {
							afo := &AdminFilterObjects{UadminDatabase: &UadminDatabase{
								Adapter: uadminDatabase.Adapter,
								Db:      tx,
							}}
							formError := form.ProceedRequest(requestForm, modelToSave, ctx, afo)
							if formError.IsEmpty() {
								if adminPage.SaveModel != nil {
									modelToSave = adminPage.SaveModel(modelToSave, ID, afo)
								} else {
									tx.Save(modelToSave)
								}
								mID := GetID(reflect.ValueOf(modelToSave))
								successfulInline := true
								for inline := range adminPage.InlineRegistry.GetAll() {
									afo1 := &AdminFilterObjects{UadminDatabase: &UadminDatabase{
										Adapter: uadminDatabase.Adapter,
										Db:      tx,
									}}
									inlineListEditableCollection, formError1 := inline.ProceedRequest(afo1, ctx, requestForm, modelToSave, adminRequestParams, c)
									if formError1 != nil {
										successfulInline = false
									}
									c.ListEditableFormsForInlines.AddForInlineWholeCollection(inline.Prefix, inlineListEditableCollection)
								}
								if !successfulInline {
									return fmt.Errorf("error while submitting inlines")
								}
								if ctx.Query("_popup") == "1" {
									data := make(map[string]interface{})
									data["Link"] = ctx.Request.URL.String()
									data["ID"] = mID
									data["Name"] = reflect.ValueOf(modelToSave).MethodByName("String").Call([]reflect.Value{})[0].Interface().(string)
									htmlResponseWriter := bytes.NewBuffer(make([]byte, 0))
									AddedObjectInPopup.ExecuteTemplate(htmlResponseWriter, "addedobjectinpopup", data)
									ctx.Data(http.StatusOK, "text/html; charset=utf-8", htmlResponseWriter.Bytes())
								} else if len(requestForm.Value["save_add_another"]) > 0 {
									ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s/edit/new", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug))
								} else if len(requestForm.Value["save_continue"]) > 0 {
									ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s/edit/%d", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug, mID))
								} else {
									ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug))
								}
								return nil
							}
							return fmt.Errorf("not successful form validation")
						})
						uadminDatabase.Close()
						if transactionerror == nil {
							return
						}
					} else {
						for inline := range adminPage.InlineRegistry.GetAll() {
							if id == "new" {
								continue
							}
							for iterateAdminObjects := range inline.GetAll(c.Model, c.AdminRequestParams) {
								listEditable := inline.ListDisplay.BuildFormForListEditable(c, iterateAdminObjects.ID, iterateAdminObjects.Model)
								c.ListEditableFormsForInlines.AddForInline(inline.Prefix, strconv.Itoa(int(iterateAdminObjects.ID)), listEditable)
							}
						}
					}
					c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.BlueprintName, URL: fmt.Sprintf("%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.ParentPage.Slug)})
					c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{Name: adminPage.ModelName, URL: fmt.Sprintf("%s/%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.ParentPage.Slug, adminPage.Slug)})
					if id != "new" {
						values := reflect.ValueOf(modelI).MethodByName("String").Call([]reflect.Value{})
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{IsActive: true, Name: values[0].String()})
					} else {
						c.BreadCrumbs.AddBreadCrumb(&AdminBreadcrumb{IsActive: true, Name: "New"})
					}
					c.Form = form
					c.PermissionForBlueprint = c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
					c.Message = ctx.Query("message")
					c.AdminModelActionRegistry = adminPage.ModelActionsRegistry
					tr := NewTemplateRenderer(adminPage.PageName)
					tr.Render(ctx, CurrentConfig.TemplatesFS, CurrentConfig.GetPathToTemplate("change"), c, FuncMap)
				}
			}(subPage))
			for adminModelAction := range subPage.ModelActionsRegistry.GetAllModelActions() {
				router.Any(fmt.Sprintf("%s/%s/%s/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.ModelName, adminModelAction.SlugifiedActionName), func(adminPage *AdminPage, slugifiedModelActionName string) func(ctx *gin.Context) {
					return func(ctx *gin.Context) {
						adminPage.HandleModelAction(slugifiedModelActionName, ctx)
					}
				}(subPage, adminModelAction.SlugifiedActionName))
			}
			for pageInline := range subPage.InlineRegistry.GetAll() {
				for inlineAdminModelAction := range pageInline.Actions.GetAllModelActions() {
					router.Any(fmt.Sprintf("%s/%s/%s/edit/:id/%s", CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.ModelName, inlineAdminModelAction.SlugifiedActionName), func(adminPage *AdminPage, adminPageInline *AdminPageInline, slugifiedModelActionName string) func(ctx *gin.Context) {
						return func(ctx *gin.Context) {
							adminPage.HandleModelAction(slugifiedModelActionName, ctx)
						}
					}(subPage, pageInline, inlineAdminModelAction.SlugifiedActionName))
				}
			}
			subPage.RegisteredHTTPHandlers = true
		}
	}
}

var CurrentDashboardAdminPanel *DashboardAdminPanel

func NewDashboardAdminPanel() *DashboardAdminPanel {
	adminPageRegistry := NewAdminPageRegistry()
	CurrentAdminPageRegistry = adminPageRegistry
	return &DashboardAdminPanel{
		AdminPages: adminPageRegistry,
	}
}

var GlobalModelActionRegistry *AdminModelActionRegistry

type RemovalTreeList []*RemovalTreeNodeStringified

func NewAdminModelActionRegistry() *AdminModelActionRegistry {
	adminModelActions := make(map[string]*AdminModelAction)
	ret := &AdminModelActionRegistry{AdminModelActions: adminModelActions}
	if GlobalModelActionRegistry != nil {
		for adminModelAction := range GlobalModelActionRegistry.GetAllModelActions() {
			ret.AddModelAction(adminModelAction)
		}
	}
	return ret
}

var AddedObjectInPopup *template.Template

func init() {
	AddedObjectInPopup, _ = template.New("addedobjectinpopup").Parse(`{{define "addedobjectinpopup"}}<html><head></head><body>
<script type="text/javascript">
	var link = "{{ .Link }}";
	var ID = "{{ .ID }}";
	var Name = "{{ .Name }}";
	var newOption = window.opener.$('<select><option value=""></option></select>');
	newOption.find('option').attr('value', ID);
	newOption.find('option').text(Name);
	newOption.find('option').attr('selected', 'selected');
	var select = window.opener.$("a[href='{{ .Link }}']").parent().parent().find('.related-target select');
	select.find('option:selected').removeAttr('selected');
	select.append(newOption.html());
	select.trigger('change');
	window.close();
</script>
</body></html>{{end}}
`)

	GlobalModelActionRegistry = NewAdminModelActionRegistry()
	removalModelAction := NewAdminModelAction(
		"Delete permanently", &AdminActionPlacement{
			ShowOnTheListPage: true,
		},
	)
	removalModelAction.RequiresExtraSteps = true
	removalModelAction.Description = "Delete users permanently"
	removalModelAction.Handler = func(ap *AdminPage, afo IAdminFilterObjects, ctx *gin.Context) (bool, int64) {
		removalPlan := make([]RemovalTreeList, 0)
		removalConfirmed := ctx.PostForm("removal_confirmed")
		afo.GetUadminDatabase().Db.Transaction(func(tx *gorm.DB) error {
			uadminDatabase := &UadminDatabase{Db: tx, Adapter: afo.GetUadminDatabase().Adapter}
			for modelIterated := range afo.IterateThroughWholeQuerySet() {
				removalTreeNode := BuildRemovalTree(uadminDatabase, modelIterated.Model)
				if removalConfirmed == "" {
					deletionStringified := removalTreeNode.BuildDeletionTreeStringified(uadminDatabase)
					removalPlan = append(removalPlan, deletionStringified)
				} else {
					err := removalTreeNode.RemoveFromDatabase(uadminDatabase)
					if err != nil {
						return err
					}
				}
			}
			if removalConfirmed != "" {
				truncateLastPartOfPath := regexp.MustCompile("/[^/]+/?$")
				newPath := truncateLastPartOfPath.ReplaceAll([]byte(ctx.Request.URL.RawPath), []byte(""))
				clonedURL := CloneNetURL(ctx.Request.URL)
				clonedURL.RawPath = string(newPath)
				clonedURL.Path = string(newPath)
				query := clonedURL.Query()
				query.Set("message", "Objects were removed succesfully")
				clonedURL.RawQuery = query.Encode()
				ctx.Redirect(http.StatusFound, clonedURL.String())
				return nil
			}
			type Context struct {
				AdminContext
				RemovalPlan []RemovalTreeList
				AdminPage   *AdminPage
				ObjectIds   string
			}
			c := &Context{}
			adminRequestParams := NewAdminRequestParams()
			c.RemovalPlan = removalPlan
			c.AdminPage = ap
			c.ObjectIds = ctx.PostForm("object_ids")
			PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)

			tr := NewTemplateRenderer(fmt.Sprintf("Remove %s ?", ap.ModelName))
			tr.Render(ctx, CurrentConfig.TemplatesFS, CurrentConfig.GetPathToTemplate("remove_objects"), c, FuncMap)
			return nil
		})
		return true, 1
	}
	GlobalModelActionRegistry.AddModelAction(removalModelAction)
}

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
			ret := &AdminFilterObjects{
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
				searchFilterObjects := &AdminFilterObjects{
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
