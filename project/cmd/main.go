package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/app"
	"github.com/rickywei/sparrow/project/logger"
)

func main() {
	defer func() {
		logger.L().Sync()
	}()

	app, err := app.WireApp()
	if err != nil {
		logger.L().Fatal("wire app failed", zap.Error(err))
	}

	rg := run.Group{}
	{
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		rg.Add(
			func() error {
				select {
				case <-term:
					return nil
				}
			},
			func(err error) {},
		)
	}
	{
		rg.Add(func() error {
			return app.Run()
		}, func(err error) {
			app.Stop()
		})
	}

	if err := rg.Run(); err != nil {
		logger.L().Error("app stopped", zap.String("err", err.Error()))
	}
}
