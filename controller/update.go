package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/user/model"
	"github.com/RSOI/user/utils"
	"github.com/RSOI/user/view"
)

// UserPATCH remove user
func UserPATCH(body []byte) (*model.User, error) {
	var err error

	var UserToUpdate model.User
	var UpdatedUser model.User
	err = json.Unmarshal(body, &UserToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return nil, err
	}

	err = view.ValidateUpdateUserRating(UserToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return nil, err
	}

	UpdatedUser, err = UserModel.UpdateUser(UserToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("User rating updated successfully")
	return &UpdatedUser, nil
}
