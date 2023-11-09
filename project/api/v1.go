package api

import "github.com/rickywei/sparrow/project/middleware"

func init() {
	routesFuncs = append(routesFuncs, func(api *API) {
		v1 := api.engine.Group("api").Group("v1")
		{
			v1.POST("user", api.userHandler.Create)
			v1.POST("user/login", api.userHandler.Login)
			v1.POST("user/refresh", api.userHandler.Refresh)
			v1.DELETE("user/:id", api.userHandler.Delete, middleware.Auth())
			v1.PUT("user/:id", api.userHandler.Update, middleware.Auth())
			v1.GET("user/:id", api.userHandler.Query, middleware.Auth())
			v1.GET("user", api.userHandler.QueryList, middleware.Auth())
		}
	})
}
