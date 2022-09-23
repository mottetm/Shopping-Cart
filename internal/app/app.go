package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type App struct {
	Echo   *echo.Echo
	Config Config
}

func New(c Config) *App {
	e := echo.New()

	e.GET("/items", c.Server.GetItems)
	e.POST("/items", c.Server.PostItem)

	return &App{e, c}
}

func (a *App) Start() {
	go func() {
		err := a.Echo.Start(fmt.Sprintf("%s:%s", a.Config.Host, a.Config.Port))
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (a *App) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := a.Echo.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
