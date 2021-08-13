package admin

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	sessionsblueprint "github.com/uadmin/uadmin/blueprint/sessions"
	interfaces2 "github.com/uadmin/uadmin/blueprint/sessions/interfaces"
	"github.com/uadmin/uadmin/interfaces"
	"github.com/uadmin/uadmin/utils"
	"gorm.io/gorm"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func PopulateTemplateContextForAdminPanel(ctx *gin.Context, context interfaces.IAdminContext, adminRequestParams *interfaces.AdminRequestParams) {
	sessionAdapter, _ := sessionsblueprint.ConcreteBlueprint.SessionAdapterRegistry.GetDefaultAdapter()
	var cookieName string
	cookieName = interfaces.CurrentConfig.D.Uadmin.AdminCookieName
	cookie, _ := ctx.Cookie(cookieName)
	var session interfaces2.ISessionProvider
	if cookie != "" {
		session, _ = sessionAdapter.GetByKey(cookie)
	}
	if adminRequestParams.CreateSession && session == nil {
		session = sessionAdapter.Create()
		expiresOn := time.Now().Add(time.Duration(interfaces.CurrentConfig.D.Uadmin.SessionDuration)*time.Second)
		session.ExpiresOn(&expiresOn)
		ctx.SetCookie(interfaces.CurrentConfig.D.Uadmin.AdminCookieName, session.GetKey(), int(interfaces.CurrentConfig.D.Uadmin.SessionDuration), "/", ctx.Request.URL.Host, interfaces.CurrentConfig.D.Uadmin.SecureCookie, interfaces.CurrentConfig.D.Uadmin.HttpOnlyCookie)
		session.Save()
	}
	if adminRequestParams.GenerateCSRFToken {
		token := utils.GenerateCSRFToken()
		currentCsrfToken, _ := session.Get("csrf_token")
		if currentCsrfToken == "" {
			session.Set("csrf_token", token)
			session.Save()
		}
	}
	if session == nil {
		session.Save()
	}
	context.SetCurrentURL(ctx.Request.URL.Path)
	context.SetCurrentQuery(ctx.Request.URL.RawQuery)
	context.SetFullURL(ctx.Request.URL)
	context.SetSiteName(interfaces.CurrentConfig.D.Uadmin.SiteName)
	context.SetRootAdminURL(interfaces.CurrentConfig.D.Uadmin.RootAdminURL)
	if session != nil {
		context.SetSessionKey(session.GetKey())
	}
	context.SetRootURL(interfaces.CurrentConfig.D.Uadmin.RootAdminURL)
	context.SetLanguage(interfaces.GetLanguage(ctx))
	context.SetLogo(interfaces.CurrentConfig.D.Uadmin.Logo)
	context.SetFavIcon(interfaces.CurrentConfig.D.Uadmin.FavIcon)
	if adminRequestParams.NeedAllLanguages {
		context.SetLanguages(interfaces.GetActiveLanguages())
	}
	// context.SetDemo()
	if session != nil {
		user := session.GetUser()
		context.SetUser(user.Username)
		context.SetUserExists(user.ID != 0)
		if user.ID != 0 {
			context.SetUserPermissionRegistry(user.BuildPermissionRegistry())
		}
	}
}

type DashboardAdminPanel struct {
	AdminPages *interfaces.AdminPageRegistry
	ListHandler func (ctx *gin.Context)
}

func (dap *DashboardAdminPanel) RegisterHttpHandlers(router *gin.Engine) {
	if dap.ListHandler != nil {
		router.GET(interfaces.CurrentConfig.D.Uadmin.RootAdminURL, dap.ListHandler)
	}
	for adminPage := range dap.AdminPages.GetAll() {
		router.GET(fmt.Sprintf("%s/%s", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug), func(pageTitle string, adminPageRegistry *interfaces.AdminPageRegistry) func (ctx *gin.Context) {
			return func(ctx *gin.Context) {
				type Context struct {
					interfaces.AdminContext
					Menu string
					CurrentPath string
				}

				c := &Context{}
				PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
				menu := string(adminPageRegistry.PreparePagesForTemplate(c.UserPermissionRegistry))
				c.Menu = menu
				c.CurrentPath = ctx.Request.URL.Path
				tr := interfaces.NewTemplateRenderer(pageTitle)
				tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("home"), c, interfaces.FuncMap)
			}
		}(adminPage.PageName, adminPage.SubPages))
		for subPage := range adminPage.SubPages.GetAll() {
			router.Any(fmt.Sprintf("%s/%s/%s", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.Slug), func(adminPage *interfaces.AdminPage) func (ctx *gin.Context) {
				return func(ctx *gin.Context) {
					if adminPage.ListHandler != nil {
						adminPage.ListHandler(ctx)
					} else {
						type Context struct {
							interfaces.AdminContext
							AdminFilterObjects *interfaces.AdminFilterObjects
							ListDisplay *interfaces.ListDisplayRegistry
							PermissionForBlueprint *interfaces.UserPerm
							ListFilter *interfaces.ListFilterRegistry
							InitialOrder string
							InitialOrderList []string
							Search string
							TotalRecords int64
							TotalPages int64
							ListEditableFormError bool
							AdminModelActionRegistry *interfaces.AdminModelActionRegistry
							Message string
						}

						c := &Context{}
						adminRequestParams := interfaces.NewAdminRequestParamsFromGinContext(ctx)
						PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
						c.Message = ctx.Query("message")
						c.PermissionForBlueprint = c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
						c.AdminFilterObjects = adminPage.GetQueryset(adminPage, adminRequestParams)
						c.AdminModelActionRegistry = adminPage.ModelActionsRegistry
						c.AdminModelActionRegistry.GetAllModelActions()
						if ctx.Request.Method == "POST" {
							c.AdminFilterObjects.WithTransaction(func(afo1 *interfaces.AdminFilterObjects) error {
								postForm, _ := ctx.MultipartForm()
								ids := postForm.Value["object_id"]
								for _, objectId := range ids {
									objectModel, _ := c.AdminFilterObjects.GenerateModelI()
									IDInt, _ := strconv.Atoi(objectId)
									IDUint := uint(IDInt)
									afo1.LoadDataForModelById(IDUint, objectModel)
									modelI, _ := c.AdminFilterObjects.GenerateModelI()
									listEditableForm := interfaces.NewFormListEditableFromListDisplayRegistry(IDUint, modelI, adminPage.ListDisplay)
									formListEditableErr := listEditableForm.ProceedRequest(postForm, objectModel)
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
						c.AdminFilterObjects.GormQuerySet.Count(&c.TotalRecords)
						c.TotalPages = int64(math.Ceil(float64(c.TotalRecords / int64(adminPage.Paginator.PerPage))))
						c.ListDisplay = adminPage.ListDisplay
						c.Search = adminRequestParams.Search
						c.ListFilter = adminPage.ListFilter
						c.InitialOrder = adminRequestParams.GetOrdering()
						c.InitialOrderList = adminRequestParams.Ordering
						tr := interfaces.NewTemplateRenderer(adminPage.PageName)
						tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("list"), c, interfaces.FuncMap)
					}
				}
			}(subPage))
			router.POST(fmt.Sprintf("%s/%s/%s/%s", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.Slug, "export"), func(adminPage *interfaces.AdminPage) func (ctx *gin.Context) {
				return func(ctx *gin.Context) {
					type Context struct {
						interfaces.AdminContext
					}
					c := &Context{}
					adminRequestParams := interfaces.NewAdminRequestParamsFromGinContext(ctx)
					PopulateTemplateContextForAdminPanel(ctx, c, interfaces.NewAdminRequestParams())
					// permissionForBlueprint := c.UserPermissionRegistry.GetPermissionForBlueprint(adminPage.BlueprintName, adminPage.ModelName)
					adminFilterObjects := adminPage.GetQueryset(adminPage, adminRequestParams)
					rows, _ := adminFilterObjects.GormQuerySet.Rows()
					defer rows.Close()
					db := interfaces.NewUadminDatabase()
					defer db.Close()
					f := excelize.NewFile()
					i := 1
					currentColumn := 'A'
					for listDisplay := range adminPage.ListDisplay.GetAllFields() {
						f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", currentColumn, i), listDisplay.DisplayName)
						currentColumn += 1
					}
					i += 1
					for rows.Next() {
						model, _ := adminFilterObjects.GenerateModelI()
						db.Db.ScanRows(rows, model)
						// db.Db.ScanRows(rows, model)
						currentColumn = 'A'
						for listDisplay := range adminPage.ListDisplay.GetAllFields() {
							f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", currentColumn, i), listDisplay.GetValue(model))
							currentColumn += 1
						}
						i += 1
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
			for adminModelAction := range subPage.ModelActionsRegistry.GetAllModelActions() {
				router.Any(fmt.Sprintf("%s/%s/%s/%s", interfaces.CurrentConfig.D.Uadmin.RootAdminURL, adminPage.Slug, subPage.ModelName, adminModelAction.SlugifiedActionName), func(adminPage *interfaces.AdminPage, slugifiedModelActionName string) func (ctx *gin.Context) {
					return func(ctx *gin.Context) {
						adminPage.HandleModelAction(slugifiedModelActionName, ctx)
					}
				}(subPage, adminModelAction.SlugifiedActionName))
			}
		}
	}
}

var CurrentDashboardAdminPanel *DashboardAdminPanel

func NewDashboardAdminPanel() *DashboardAdminPanel {
	adminPageRegistry := interfaces.NewAdminPageRegistry()
	interfaces.CurrentAdminPageRegistry = adminPageRegistry
	return &DashboardAdminPanel{
		AdminPages: adminPageRegistry,
	}
}

var GlobalModelActionRegistry *interfaces.AdminModelActionRegistry

type RemovalTreeList []*interfaces.RemovalTreeNodeStringified

func NewAdminModelActionRegistry() *interfaces.AdminModelActionRegistry {
	adminModelActions := make(map[string]*interfaces.AdminModelAction)
	ret := &interfaces.AdminModelActionRegistry{AdminModelActions: adminModelActions}
	if GlobalModelActionRegistry != nil {
		for adminModelAction := range GlobalModelActionRegistry.GetAllModelActions() {
			ret.AddModelAction(adminModelAction)
		}
	}
	return ret
}

func init() {
	GlobalModelActionRegistry = NewAdminModelActionRegistry()
	removalModelAction := interfaces.NewAdminModelAction(
		"Delete permanently", &interfaces.AdminActionPlacement{
			ShowOnTheListPage: true, DisplayToTheLeft: true, DisplayToTheTop: true,
		},
	)
	removalModelAction.RequiresExtraSteps = true
	removalModelAction.Description = "Delete users permanently"
	removalModelAction.Handler = func (ap *interfaces.AdminPage, afo *interfaces.AdminFilterObjects, ctx *gin.Context) (bool, int64) {
		removalPlan := make([]RemovalTreeList, 0)
		removalConfirmed := ctx.PostForm("removal_confirmed")
		afo.UadminDatabase.Db.Transaction(func(tx *gorm.DB) error {
			uadminDatabase := &interfaces.UadminDatabase{Db: tx, Adapter: afo.UadminDatabase.Adapter}
			for modelIterated := range afo.IterateThroughModelActionsSelected() {
				removalTreeNode := interfaces.BuildRemovalTree(uadminDatabase, modelIterated.Model)
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
				clonedUrl := interfaces.CloneNetUrl(ctx.Request.URL)
				clonedUrl.RawPath = string(newPath)
				clonedUrl.Path = string(newPath)
				query := clonedUrl.Query()
				query.Set("message", "Objects were removed succesfully")
				clonedUrl.RawQuery = query.Encode()
				ctx.Redirect(http.StatusFound, clonedUrl.String())
				return nil
			}
			type Context struct {
				interfaces.AdminContext
				RemovalPlan []RemovalTreeList
				AdminPage *interfaces.AdminPage
				ObjectIds string
			}
			c := &Context{}
			adminRequestParams := interfaces.NewAdminRequestParams()
			c.RemovalPlan = removalPlan
			c.AdminPage = ap
			c.ObjectIds = ctx.PostForm("object_ids")
			PopulateTemplateContextForAdminPanel(ctx, c, adminRequestParams)

			tr := interfaces.NewTemplateRenderer(fmt.Sprintf("Remove %s ?", ap.ModelName))
			tr.Render(ctx, interfaces.CurrentConfig.TemplatesFS, interfaces.CurrentConfig.GetPathToTemplate("remove_objects"), c, interfaces.FuncMap)
			return nil
		})
		return true, 1
	}
	GlobalModelActionRegistry.AddModelAction(removalModelAction)
}

func NewGormAdminPage(parentPage *interfaces.AdminPage, genModelI func() (interface{}, interface{}), modelName string) *interfaces.AdminPage {
	modelI4, _ := genModelI()
	return &interfaces.AdminPage{
		SubPages: interfaces.NewAdminPageRegistry(),
		ParentPage: parentPage,
		GetQueryset: func(adminPage *interfaces.AdminPage, adminRequestParams *interfaces.AdminRequestParams) *interfaces.AdminFilterObjects {
			uadminDatabase := interfaces.NewUadminDatabase()
			db := uadminDatabase.Db
			var paginatedQuerySet *gorm.DB
			var perPage int
			modelI, _ := genModelI()
			modelI1, _ := genModelI()
			modelI2, _ := genModelI()
			modelI3, _ := genModelI()
			ret := &interfaces.AdminFilterObjects{
				InitialGormQuerySet: db.Model(modelI),
				GormQuerySet: db.Model(modelI1),
				PaginatedGormQuerySet: db.Model(modelI2),
				Model: modelI3,
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
		Model: modelI4,
		ModelName: modelName,
		Validators: make([]interfaces.IValidator, 0),
		ExcludeFields: interfaces.NewFieldRegistry(),
		FieldsToShow: interfaces.NewFieldRegistry(),
		ModelActionsRegistry: NewAdminModelActionRegistry(),
		Inlines: make([]*interfaces.AdminPageInlines, 0),
		ListDisplay: &interfaces.ListDisplayRegistry{ListDisplayFields: make(map[string]*interfaces.ListDisplay)},
		ListFilter: &interfaces.ListFilterRegistry{ListFilter: make([]*interfaces.ListFilter, 0)},
		SearchFields: make([]*interfaces.SearchField, 0),
		Paginator: &interfaces.Paginator{PerPage: interfaces.CurrentConfig.D.Uadmin.AdminPerPage, ShowLastPageOnPreviousPage: true},
		ActionsSelectionCounter: true,
		FilterOptions: interfaces.NewFilterOptionsRegistry(),
	}
}
