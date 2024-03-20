package router

import (
	"net/http"
	"whats-app-clone-service/middleware"
	"whats-app-clone-service/utils"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	AuthRouter(router)
	router.Handle("/api/user", utils.JWTMiddleware(http.HandlerFunc(middleware.GetAllUser))).Methods("GET")

	return router
}
