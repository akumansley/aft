package lib

import (
	"awans.org/aft/internal/db"
)

type Module interface {
	ProvideRoutes() []Route
	ProvideMiddleware() []Middleware
	ProvideModels() []db.Model
	ProvideRelationships() []db.Relationship
	ProvideRecords(db.RWTx) error
	ProvideDatatypes() []db.Datatype
	ProvideCode() []db.Code
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

func (bm *BlankModule) ProvideRelationships() []db.Relationship {
	return []db.Relationship{}
}

func (bm *BlankModule) ProvideRecords(db.RWTx) error {
	return nil
}

func (bm *BlankModule) ProvideHandlers() []interface{} {
	return []interface{}{}
}

func (bm *BlankModule) ProvideDatatypes() []db.Datatype {
	return []db.Datatype{}
}

func (bm *BlankModule) ProvideCode() []db.Code {
	return []db.Code{}
}
