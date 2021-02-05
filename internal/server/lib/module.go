package lib

import (
	"awans.org/aft/internal/db"
)

type Module interface {
	ID() db.ID
	Name() string
	Package() string
	ProvideRoutes() []Route
	ProvideMiddleware() []Middleware
	ProvideInterfaces() []db.InterfaceL
	ProvideDatatypes() []db.DatatypeL
	ProvideFunctions() []db.FunctionL
	ProvideHandlers() []interface{}
	ProvideFunctionLoaders() []db.FunctionLoader
	ProvideLiterals() []db.Literal
}

type BlankModule struct{}

func (bm *BlankModule) ProvideRoutes() []Route {
	return []Route{}
}

func (bm *BlankModule) ProvideMiddleware() []Middleware {
	return []Middleware{}
}

func (bm *BlankModule) ProvideInterfaces() []db.InterfaceL {
	return []db.InterfaceL{}
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

func ToLiteral(mod Module) db.ModuleL {
	var modLits []db.ModuleLiteral
	for _, lit := range mod.ProvideLiterals() {
		ml, ok := lit.(db.ModuleLiteral)
		if !ok {
			continue
		}
		modLits = append(modLits, ml)
	}

	ifaces := mod.ProvideInterfaces()
	for _, fl := range mod.ProvideFunctionLoaders() {
		ifaces = append(ifaces, fl.ProvideModel())
	}

	return db.MakeModule(
		mod.ID(),
		mod.Name(),
		mod.Package(),
		ifaces,
		mod.ProvideFunctions(),
		mod.ProvideDatatypes(),
		modLits,
	)
}
