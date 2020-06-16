package lib

import (
	"awans.org/aft/internal/db"
)

type Module interface {
	ProvideRoutes() []Route
	ProvideMiddleware() []Middleware
	ProvideModels() []db.Model
	ProvideRecords() []db.Record
	ProvideHandlers() []interface{}
}

type BlankModule struct {
}

func (bm *BlankModule) ProvideRoutes() []Route {
	return []Route{}
}

func (bm *BlankModule) ProvideMiddleware() []Middleware {
	return []Middleware{}
}

func (bm *BlankModule) ProvideModels() []db.Model {
	return []db.Model{}
}

func (bm *BlankModule) ProvideRecords() []db.Record {
	return []db.Record{}
}

func (bm *BlankModule) ProvideHandlers() []interface{} {
	return []interface{}{}
}
