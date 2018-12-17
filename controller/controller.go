package controller

import (
	"github.com/RSOI/user/model"
	"github.com/RSOI/user/utils"
	"github.com/jackc/pgx"
)

var (
	// UserModel interface with methods
	UserModel model.UServiceInterface
)

// Init Init model with pgx connection
func Init(db *pgx.ConnPool) {
	utils.LOG("Setup model...")
	UserModel = &model.UService{
		Conn: db,
	}
}
