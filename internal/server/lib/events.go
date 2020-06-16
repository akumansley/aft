package lib

import (
	"awans.org/aft/internal/db"
)

type ParseRequest struct {
	Request interface{}
}

type DatabaseReady struct {
	Db db.DB
}
