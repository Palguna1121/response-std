package router

import "github.com/gin-gonic/gin"

type RouteSetupFunc func(*gin.Engine)

var RouteRegistry = make(map[string]RouteSetupFunc)

func Register(version string, setup RouteSetupFunc) {
	RouteRegistry[version] = setup
}
