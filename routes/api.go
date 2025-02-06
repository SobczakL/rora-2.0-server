package routes

import (
	"net/http"

	"rora-server/controllers"
)

func RegisterAPIRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/data", controllers.GetData)
}
