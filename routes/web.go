package routes

import (
	"response-std/libs/router"
	"response-std/routes/web"
)

func init() {
	router.Register("web", web.SetupWebRoutes)
}
