package model

import (
	"github.com/RSOI/user/ui"
	"github.com/RSOI/user/utils"
	"github.com/jackc/pgx"
)

// AddUser add new user
func (service *UService) AddUser(u User) (User, error) {
	var err error

	utils.LOG("Accessing database...")
	row := service.Conn.QueryRow(
		`INSERT INTO quser.user (nickname) VALUES ($1) RETURNING id`,
		u.Nickname,
	)

	err = row.Scan(&u.ID)
	return u, err
}

// GetUserByID get user data by it's id
func (service *UService) GetUserByID(uID int) (User, error) {
	var err error
	var u User

	utils.LOG("Accessing database...")
	row := service.Conn.QueryRow(`SELECT * FROM quser.user WHERE id = $1`, uID)

	err = row.Scan(
		&u.ID,
		&u.Nickname,
		&u.Rating)

	return u, err
}

// UpdateUser Update user rating
func (service *UService) UpdateUser(u User) (User, error) {
	currentUserData, err := service.GetUserByID(u.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = ui.ErrNoDataToUpdate
		}
		return u, err
	}

	utils.LOG("Accessing database...")
	res, err := service.Conn.Exec(`UPDATE quser.user SET rating = rating + $1`, u.Rating)
	if err == nil && res.RowsAffected() != 1 {
		err = ui.ErrNoDataToUpdate
	}
	currentUserData.Rating += u.Rating
	return currentUserData, err
}
