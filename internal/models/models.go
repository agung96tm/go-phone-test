package models

import "database/sql"

type Models struct {
	Phone PhoneModel
}

func New(db *sql.DB) *Models {
	return &Models{
		PhoneModel{DB: db},
	}
}
