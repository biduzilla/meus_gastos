package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users       UserModel
	Permissions PermissionModel
	Categories  CategoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Categories:  CategoryModel{DB: db},
	}
}
