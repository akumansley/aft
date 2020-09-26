package lib

import (
	"awans.org/aft/internal/db"
)

type Module interface {
	ProvideRoutes() []Route
	ProvideMiddleware() []Middleware
	ProvideModels() []db.ModelL
	ProvideDatatypes() []db.DatatypeL
	ProvideFunctions() []db.FunctionL
	ProvideHandlers() []interface{}
	ProvideFunctionLoaders() []db.FunctionLoader
	ProvideLiterals() []db.Literal
}

type BlankModule struct {
}

func (bm *BlankModule) ProvideRoutes() []Route {
	return []Route{}
}

func (bm *BlankModule) ProvideMiddleware() []Middleware {
	return []Middleware{}
}

func (bm *BlankModule) ProvideModels() []db.ModelL {
	return []db.ModelL{}
}

func (bm *BlankModule) ProvideHandlers() []interface{} {
	return []interface{}{}
}

func (bm *BlankModule) ProvideDatatypes() []db.DatatypeL {
	return []db.DatatypeL{}
}

func (bm *BlankModule) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{}
}

func (bm *BlankModule) ProvideFunctionLoaders() []db.FunctionLoader {
	return []db.FunctionLoader{}
}

func (bm *BlankModule) ProvideLiterals() []db.Literal {
	return []db.Literal{}
}
