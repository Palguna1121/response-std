package router

import (
	v1 "response-std/v1/routes"
	v2 "response-std/v2/routes"

	"github.com/gin-gonic/gin"
)

type RouteSetupFunc func(*gin.Engine)

var RouteRegistry = map[string]RouteSetupFunc{
	"v1": v1.SetupRoutes,
	"v2": v2.SetupRoutes,
}
