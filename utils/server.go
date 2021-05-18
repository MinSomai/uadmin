package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiHandlerRequest struct {
	Methods       []string
	Handler       func(*gin.Context)
	Detail        bool
	Subroute_name string
}

func HandleApiRequest(apiHandlers map[string]func(c *gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		path := c.Param("regex")
		regex := UrlRemainderRegex.Copy()
		if router_matches := regex.FindStringSubmatch(path); len(router_matches) > 0 {
			api_handler, exists := apiHandlers[router_matches[regex.SubexpIndex("Path")]]
			if !exists {
				c.JSON(404, ApiNoMethodFound())
				return
			}
			api_handler(c)
		}
	}
}

func InitializeRouter(r *gin.Engine, route_prefix string, apiHandlers []ApiHandlerRequest) {
	router_group := r.Group(fmt.Sprintf("/%s", route_prefix))
	var route_path string
	for _, apiHandler := range apiHandlers {
		for _, method := range apiHandler.Methods {
			method = strings.ToLower(method)
			route_path = "/"
			if apiHandler.Detail {
				route_path = "/:id/"
			}
			if len(apiHandler.Subroute_name) > 0 {
				route_path += apiHandler.Subroute_name + "/"
			}
			if method == "get" {
				router_group.GET(route_path, apiHandler.Handler)
			} else if method == "post" {
				router_group.POST(route_path, apiHandler.Handler)
			} else if method == "delete" {
				router_group.DELETE(route_path, apiHandler.Handler)
			} else if method == "patch" {
				router_group.PATCH(route_path, apiHandler.Handler)
			} else if method == "put" {
				router_group.PUT(route_path, apiHandler.Handler)
			} else if method == "head" {
				router_group.HEAD(route_path, apiHandler.Handler)
			}
		}
	}
}
