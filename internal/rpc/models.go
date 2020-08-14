package rpc

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
	"io/ioutil"
	"os"
	"strings"
)

var RPCModel = db.MakeModel(
	db.MakeID("29209517-1c39-4be9-9808-e1ed8e40c566"),
	"rpc",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("6ec81a63-0406-4d13-aacf-070a26c2adbc"),
			"name",
			db.String,
		),
	},
	[]db.RelationshipL{RPCCode},
	[]db.ConcreteInterfaceL{},
)

var RPCCode = db.MakeConcreteRelationship(
	db.MakeID("9221119b-495a-449c-b2b3-2c6610f89d7b"),
	"code",
	false,
	db.FunctionInterface,
)

func reactForm() string {
    wd, _ := os.Getwd()
    base := strings.Split(wd, "aft")[0]
	fileB, _ := ioutil.ReadFile(base + "aft/internal/rpc/reactForm.star")
	return string(fileB)
}

var reactFormRPC = starlark.MakeStarlarkFunction(
	db.MakeID("d8179f1f-d94e-4b81-953b-6c170d3de9b7"),
	"reactForm",
	db.RPC, reactForm())

func validateForm() string {
    wd, _ := os.Getwd()
    base := strings.Split(wd, "aft")[0]
	fileB, _ := ioutil.ReadFile(base + "aft/internal/rpc/validateForm.star")
	return string(fileB)
}

var validateFormRPC = starlark.MakeStarlarkFunction(
	db.MakeID("d7633de5-9fa2-4409-a1b2-db96a59be52b"),
	"validateForm",
	db.RPC,validateForm())

func terminal() string {
    wd, _ := os.Getwd()
    base := strings.Split(wd, "aft")[0]
	fileB, _ := ioutil.ReadFile(base + "aft/internal/rpc/terminal.star")
	return string(fileB)
}

var terminalRPC = starlark.MakeStarlarkFunction(
	db.MakeID("180bb262-8835-4ec5-9c2b-3f455615be9a"),
	"terminal",
	db.RPC,terminal())

func lint() string {
    wd, _ := os.Getwd()
    base := strings.Split(wd, "aft")[0]
	fileB, _ := ioutil.ReadFile(base + "aft/internal/rpc/lint.star")
	return string(fileB)
}

var lintRPC = starlark.MakeStarlarkFunction(
	db.MakeID("e4be72dc-9462-49f7-bba9-3543cc6bf6c2"),
	"lint",
	db.RPC,lint())
