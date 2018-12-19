package controller

import (
	"fmt"
	"strconv"

	"github.com/RSOI/user/model"
	"github.com/RSOI/user/utils"
)

// UserGET get user by id
func UserGET(id string) (*model.User, error) {
	aID, _ := strconv.Atoi(id)

	data, err := UserModel.GetUserByID(aID)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("User was found successfully")
	return &data, nil
}
