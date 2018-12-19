package controller

import (
	"encoding/json"
	"testing"

	"github.com/RSOI/user/model"
	"github.com/RSOI/user/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedUService struct {
	mock.Mock
}

func getMock() *MockedUService {
	UserModel = &MockedUService{}
	return UserModel.(*MockedUService)
}

func (s *MockedUService) AddUser(u model.User) (model.User, error) {
	args := s.Mock.Called(u)
	return args.Get(0).(model.User), args.Error(1)
}
func (s *MockedUService) GetUserByID(uID int) (model.User, error) {
	args := s.Mock.Called(uID)
	return args.Get(0).(model.User), args.Error(1)
}
func (s *MockedUService) UpdateUser(u model.User) (model.User, error) {
	args := s.Mock.Called(u)
	return args.Get(0).(model.User), args.Error(1)
}
func (s *MockedUService) GetUsageStatistic(host string) (model.ServiceStatus, error) {
	args := s.Mock.Called(host)
	return args.Get(0).(model.ServiceStatus), args.Error(1)
}
func (s *MockedUService) LogStat(request []byte, responseStatus int, responseError string) {
	// nothing interesting here, just store data without affecting main thread
}

var (
	defaultUserContent = "My User Content"
	nickname           = "MyNickname"
	defaultUser        = model.User{
		Nickname: &nickname,
	}
	createdUser = model.User{
		ID:       1,
		Nickname: &nickname,
		Rating:   0,
	}
	userToUpdate = model.User{
		ID:     1,
		Rating: -10,
	}
	updatedUser = model.User{
		ID:       1,
		Nickname: &nickname,
		Rating:   -10,
	}
)

/*
********************************************************************
TESTS FOR SIGNUP ***************************************************
********************************************************************
*/

func TestAskAddCorrectData(t *testing.T) {
	body, _ := json.Marshal(&defaultUser)

	cMock := getMock()
	cMock.On("AddUser", defaultUser).Return(createdUser, nil)

	data, err := UserPUT(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, createdUser.ID, data.ID)
		assert.Equal(t, *createdUser.Nickname, *data.Nickname)
		assert.Equal(t, createdUser.Rating, data.Rating)
	}
}

func TestUserMissedField(t *testing.T) {
	body := []byte("{}")

	data, err := UserPUT(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
	assert.Nil(t, data)
}

func TestUserBrokenBody(t *testing.T) {
	body := []byte("{nickname: test}")

	data, err := UserPUT(body)
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

// /*
// ********************************************************************
// TESTS FOR GET USER BY ID *******************************************
// ********************************************************************
// */

func TestUserGetByIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("GetUserByID", 1).Return(createdUser, nil)

	data, err := UserGET("1")
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, createdUser.ID, data.ID)
		assert.Equal(t, *createdUser.Nickname, *data.Nickname)
		assert.Equal(t, createdUser.Rating, data.Rating)
	}
}

func TestUserGetByIDNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("GetUserByID", 0).Return(model.User{}, ui.ErrNoResult)

	data, err := UserGET("0")
	if assert.NotNil(t, err) {
		cMock.AssertExpectations(t)

		assert.Nil(t, data)
		assert.Equal(t, ui.ErrNoResult, err)
	}
}

/*
********************************************************************
TESTS FOR UPDATE USER **********************************************
********************************************************************
*/

func TestUpdateCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("UpdateUser", userToUpdate).Return(updatedUser, nil)

	body, _ := json.Marshal(userToUpdate)
	response, err := UserPATCH(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, updatedUser.ID, response.ID)
		assert.Equal(t, updatedUser.Rating, response.Rating)
	}
}

func TestUpdateNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("UpdateUser", userToUpdate).Return(model.User{}, ui.ErrNoDataToUpdate)

	body, _ := json.Marshal(userToUpdate)
	data, err := UserPATCH(body)
	if assert.NotNil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, ui.ErrNoDataToUpdate, err)
		assert.Nil(t, data)
	}
}

func TestUpdateMissedID(t *testing.T) {
	body := []byte("{\"rating\": -10}")

	response, err := UserPATCH(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
	assert.Equal(t, (*model.User)(nil), response)
}

func TestUpdateBrokenBody(t *testing.T) {
	data, err := UserPATCH([]byte("{id: 1}"))

	if assert.NotNil(t, err) {
		assert.Nil(t, data)
	}
}
