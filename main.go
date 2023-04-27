package main

import (
	"flag"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var (
	maxInvitees         int
	defaultVideoQuality string
	defaultAudioQuality string
	addr                string
)

func main() {
	flag.IntVar(&maxInvitees, "m", 100, "set the maximum number of users permitted for an event")
	flag.StringVar(&defaultVideoQuality, "v", "720p", "set the default video quality served to clients")
	flag.StringVar(&defaultAudioQuality, "a", "mid", "set the default audio quality served to clients")
	flag.StringVar(&addr, "p", ":5555", "set the port the server listens to")
	flag.Parse()
	initRepo()
	e := echo.New()
	e.Validator = &customValidator{
		validator: validator.New(),
	}
	setRoutes(e)

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))
	// Run
	e.Logger.Fatal(e.Start(addr))
}
