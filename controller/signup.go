package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/user/model"
	"github.com/RSOI/user/utils"
	"github.com/RSOI/user/view"
)

// UserPUT new user
func UserPUT(body []byte) (*model.User, error) {
	var err error

	var NewUser model.User
	err = json.Unmarshal(body, &NewUser)
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return nil, err
	}

	err = view.ValidateNewUser(NewUser)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return nil, err
	}

	NewUser, err = UserModel.AddUser(NewUser)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("New user added successfully")
	return &NewUser, nil
}
