package server

import (
	"awans.org/aft/internal/server/services"
	"net/http"
)

type Route struct {
	Name        string
	Service     string
	Method      string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{
		"ListObjects",
		"objects",
		"list",
		services.ListObjects,
	},
	Route{
		"InfoObjects",
		"objects",
		"info",
		services.InfoObjects,
	},
}
