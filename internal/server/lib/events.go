package lib

import (
	"awans.org/aft/internal/db"
)

type ParseRequest struct {
	Request interface{}
}

type DatabaseReady struct {
	DB db.DB
}
