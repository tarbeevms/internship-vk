package route

import (
	"myapp/api/controllers"
	"myapp/api/middleware"

	"github.com/gorilla/mux"
)

func CreateNewRoute(h *controllers.HandlerLayer) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/api/login", h.LoginHandler).Methods("POST")
	router.HandleFunc("/api/write", middleware.Authenticated(h.WriteDataHandler, h.LogicLayer)).Methods("POST")
	router.HandleFunc("/api/read", middleware.Authenticated(h.ReadDataHandler, h.LogicLayer)).Methods("POST")
	return router
}
