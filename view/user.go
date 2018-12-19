package view

import (
	"github.com/RSOI/user/model"
	"github.com/RSOI/user/ui"
)

// ValidateNewUser returns nil if all the required form values are passed
func ValidateNewUser(data model.User) error {
	if data.Nickname == nil || *data.Nickname == "" {
		return ui.ErrFieldsRequired
	}
	return nil
}

// ValidateUpdateUser returns nil if parameter to delete found
func ValidateUpdateUserRating(data model.User) error {
	if data.ID != 0 {
		return nil
	}
	return ui.ErrFieldsRequired
}
