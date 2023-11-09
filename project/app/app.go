package app

import (
	"github.com/google/uuid"

	"github.com/rickywei/sparrow/project/api"
)

type App struct {
	ID      string
	Version string
	Api     *api.API
}

func NewApp(api *api.API) (*App, error) {
	return &App{
		ID:  uuid.New().String(),
		Api: api,
	}, nil
}

func (app *App) Run() (err error) {
	return app.Api.Run()
}

func (app *App) Stop() {
	app.Api.Stop()
}
