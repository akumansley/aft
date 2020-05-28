package auth

import (
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func GetModule() lib.Module {
	return Module{}
}
