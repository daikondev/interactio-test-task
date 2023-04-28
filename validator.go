package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type customValidator struct {
	validator *validator.Validate
}

var (
	dateRegex  = regexp.MustCompile(`\d{4}(.\d{2}){2}(\s|T)(\d{2}.){2}\d{2}`)
	emailRegex = regexp.MustCompile(`^(?P<local>[a-zA-Z0-9.!#$%&'*+/=?^_\x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)
)

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// validateDate ensures the date passed in the request adheres to the ISO 8601 format
func validateDate(date string) error {
	if !dateRegex.MatchString(date) {
		return errors.New("incorrect datetime format")
	}
	return nil
}

// validateInvitees ensures the number of invitees specified in the request does not exceed the specified maximum
func validateInvitees(invitees []string) error {
	if len(invitees) > maxInvitees {
		err := fmt.Errorf(`
	invitees exceed capacity\n
	current invitees: %c, capacity: %c\n
	`, len(invitees), maxInvitees)
		return err
	}
	for _, invitee := range invitees {
		if !emailRegex.MatchString(invitee) {
			err := fmt.Errorf("invalid Email address: %s", invitee)
			return err
		}
	}
	return nil
}

// validateUserIsInvited is used to validate if a user's email is in the invitee list, in order to be granted access
// to the event details
func validateUserIsInvited(userEmail string, eventID int64) (bool, error) {
	userIsInvited, err := EventRepo.checkIfInvited(userEmail, eventID)
	if err != nil {
		return false, err
	}
	return userIsInvited, nil
}
