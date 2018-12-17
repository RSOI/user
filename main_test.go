package main

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/RSOI/user/controller"
	"github.com/RSOI/user/model"
	"github.com/RSOI/user/ui"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedUService struct {
	mock.Mock
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
	HOST        = "http://localhost"
	nickname    = "MyNickname"
	defaultUser = model.User{
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

func initServer() (*fasthttp.Client, *fasthttp.Request, *fasthttp.Response, *MockedUService) {
	listener := fasthttputil.NewInmemoryListener()
	server := &fasthttp.Server{
		Handler: initRoutes().Handler,
	}
	go server.Serve(listener)

	client := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return listener.Dial()
		},
	}

	controller.UserModel = &MockedUService{}
	cMock := controller.UserModel.(*MockedUService)
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	return client, req, res, cMock
}

/*
********************************************************************
TESTS FOR SIGNUP ***************************************************
********************************************************************
*/

func TestSignupCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	a, _ := json.Marshal(&defaultUser)

	req.SetRequestURI(HOST + "/signup")
	req.Header.SetMethod("PUT")
	req.SetBody(a)

	cMock.On("AddUser", defaultUser).Return(createdUser, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 201, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 201, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, createdUser.ID, int(responseData["id"].(float64)))
		assert.Equal(t, *createdUser.Nickname, responseData["nickname"])
		assert.Equal(t, createdUser.Rating, int(responseData["rating"].(float64)))
	}
}

func TestUserMissedField(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/signup")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, ui.ErrFieldsRequired.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestUserBrokenBody(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/signup")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{id: 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 500, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 500, response.Status)
		assert.NotEqual(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR GET USER *************************************************
********************************************************************
*/

func TestUserGetByIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/user/id1")
	req.Header.SetMethod("GET")

	cMock.On("GetUserByID", 1).Return(createdUser, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, createdUser.ID, int(responseData["id"].(float64)))
		assert.Equal(t, *createdUser.Nickname, responseData["nickname"])
		assert.Equal(t, createdUser.Rating, int(responseData["rating"].(float64)))
	}
}

func TestUserGetByIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/user/id0")
	req.Header.SetMethod("GET")

	cMock.On("GetUserByID", 0).Return(model.User{}, ui.ErrNoResult)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoResult.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR UPDATE QUESTION ******************************************
********************************************************************
*/

func TestUpdateCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	source, _ := json.Marshal(userToUpdate)

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBody(source)

	cMock.On("UpdateUser", userToUpdate).Return(updatedUser, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)

		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, userToUpdate.ID, int(responseData["id"].(float64)))
		assert.Equal(t, userToUpdate.Rating, int(responseData["rating"].(float64)))
	}
}

func TestUpdateNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	source, _ := json.Marshal(userToUpdate)

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBody(source)

	cMock.On("UpdateUser", userToUpdate).Return(model.User{}, ui.ErrNoDataToUpdate)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoDataToUpdate.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestUpdateMissedID(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBodyString("{}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, "missed required field(s)", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestUpdateBrokenBody(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/update")
	req.Header.SetMethod("PATCH")
	req.SetBodyString("{id: 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 500, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 500, response.Status)
		assert.NotEqual(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}
