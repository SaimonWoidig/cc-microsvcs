package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/SaimonWoidig/cc-microsvcs/service.auth/pkg/service"
)

func main() {
	c := service.NewContainer()
	c.Logger.Info("container initialized")

	c.Logger.Debug("dumping config", "config", c.Config)

	interruptCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	<-interruptCtx.Done()
	c.Logger.Info("got interrupt, shutting down")

	defer func() {
		if err := c.Shutdown(); err != nil {
			c.Logger.Error("error while shutting down container", "error", err.Error())
			panic(err.Error())
		}
	}()
}
