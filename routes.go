package main

import "github.com/labstack/echo/v4"

func setRoutes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/", handleHello)
	e.GET("/events", handleGetAllEvents)
	e.POST("/events", handleEventCreate)
	e.GET("/events/:id", handleGetOneEvent)
	return e
}
