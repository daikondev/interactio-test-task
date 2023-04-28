package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// Errors

var (
	errBadRequest          = echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	errInvalidID           = echo.NewHTTPError(http.StatusBadRequest, "invalid event id")
	errNotFound            = echo.NewHTTPError(http.StatusNotFound, "Event Not Found")
	errInternalServerError = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
)

// Handlers

func handleHello(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the test exercise\n")
}

func handleEventCreate(c echo.Context) error {
	ev := event{}
	if err := c.Bind(&ev); err != nil {
		return err
	}
	if err := c.Validate(ev); err != nil {
		return errBadRequest.SetInternal(err)
	}
	if err := validateDate(ev.Date); err != nil {
		return errBadRequest.SetInternal(err)
	}
	if err := validateInvitees(ev.Invitees); err != nil {
		return errBadRequest.SetInternal(err)
	}
	res, err := EventRepo.create(ev)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func handleGetOneEvent(c echo.Context) error {
	qs := eventQueryStringFields{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errInvalidID
	}

	qs.VideoQuality = c.QueryParam("videoQuality")
	if qs.VideoQuality == "" {
		qs.VideoQuality = defaultVideoQuality
	}

	qs.AudioQuality = c.QueryParam("audioQuality")
	if qs.AudioQuality == "" {
		qs.AudioQuality = defaultAudioQuality
	}

	res, err := EventRepo.getOneEvent(id, qs)
	if err != nil {
		return errNotFound
	}
	return c.JSON(http.StatusOK, res)
}

func handleGetAllEvents(c echo.Context) error {
	all, err := EventRepo.getAllEvents()
	if err != nil {
		return errInternalServerError
	}
	return c.JSON(http.StatusOK, all)
}
