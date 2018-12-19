package model

import (
	"github.com/jackc/pgx"
)

// User interface
type User struct {
	ID       int     `json:"id"`
	Nickname *string `json:"nickname"`
	Rating   int     `json:"rating"`
}

// UService connection holder
type UService struct {
	Conn *pgx.ConnPool
}

// UServiceInterface user methods interface
type UServiceInterface interface {
	AddUser(a User) (User, error)
	GetUserByID(aID int) (User, error)
	UpdateUser(a User) (User, error)
	GetUsageStatistic(host string) (ServiceStatus, error)
	LogStat(request []byte, responseStatus int, responseError string)
}
