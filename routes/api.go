package routes

import (
	"response-std/libs/router"
	"response-std/routes/api"
)

func init() {
	router.Register("v1", api.SetupRoutesv1)
}
