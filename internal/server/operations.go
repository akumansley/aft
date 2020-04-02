package server

import (
	"awans.org/aft/internal/server/services"
	"net/http"
)

type Server interface {
	Parse(*http.Request) interface{}
	Serve(http.ResponseWriter, interface{})
}

type Operation struct {
	Name    string
	Service string
	Method  string
	Write   bool
	Server  Server
}

var operations = []Operation{
	Operation{
		"ListObjects",
		"objects",
		"list",
		false,
		services.ListObjectsServer{},
	},
	Operation{
		"InfoObjects",
		"objects",
		"info",
		false,
		services.InfoObjectsServer{},
	},
	Operation{
		"Query",
		"{object}",
		"query",
		false,
		services.QueryServer{},
	},
	Operation{
		"Create",
		"{object}",
		"create",
		true,
		services.CreateServer{},
	},
}
