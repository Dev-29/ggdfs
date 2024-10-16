package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

const ACCESS_TOKEN_HEADER = "X-Access-Token"

func main() {
	var port string
	flag.StringVar(&port, "port", "9191", "pass the port number")
	flag.Parse()
	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"http://localhost"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, ACCESS_TOKEN_HEADER},
		ExposeHeaders: []string{echo.HeaderContentLength, echo.HeaderContentDisposition, echo.HeaderContentEncoding},
	}))
	router := e.Group("arbokcore")
	// Routes
	router.GET("/ping", hello)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: e,
	}
	go func() {
		log.Info().Str("port", port).Msg("server started at")
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()
	appCtx := context.Background()
	ctx, stop := signal.NotifyContext(appCtx, os.Interrupt)
	defer stop()
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("error during server shutdown")
	}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
