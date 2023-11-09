package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/docs"
	"github.com/rickywei/sparrow/project/handler"
	"github.com/rickywei/sparrow/project/middleware"
)

var (
	ProviderSet = wire.NewSet(NewApi)

	routesFuncs = []func(*API){}
)

type API struct {
	engine *gin.Engine
	srv    *http.Server
	ctx    context.Context
	cancel context.CancelFunc

	userHandler *handler.UserHandler
}

func NewApi(userHandler *handler.UserHandler) *API {
	engine := gin.New()
	engine.Use(middleware.Logger(), middleware.Recover())
	// swagger
	if conf.Bool("app.doc.enable") {
		// https://github.com/swaggo/swag#swag
		docs.SwaggerInfo.Title = "App api"
		docs.SwaggerInfo.Version = ""
		docs.SwaggerInfo.BasePath = "/"
		engine.GET("/swagger/*any",
			gin.BasicAuth(gin.Accounts{conf.String("app.doc.account"): conf.String("app.doc.password")}),
			ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.String("app.ip"), conf.Int("app.port")),
		Handler: engine,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	api := &API{
		engine: engine,
		srv:    srv,
		ctx:    ctx,
		cancel: cancel,

		userHandler: userHandler,
	}

	for _, f := range routesFuncs {
		f(api)
	}

	return api
}

func (a *API) Run() (err error) {
	return a.srv.ListenAndServe()
}

func (a *API) Stop() {
	defer a.cancel()

	a.srv.Shutdown(a.ctx)
}
