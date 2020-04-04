package db

import (
	"awans.org/aft/internal/db"
)

func InitDB() {
	DB = db.New()
}

var DB db.DB
