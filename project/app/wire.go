//go:build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/rickywei/sparrow/project/api"
	"github.com/rickywei/sparrow/project/handler"
)

func WireApp() (*App, error) {
	panic(wire.Build(handler.ProviderSet, api.ProviderSet, NewApp))
}
