package router

import (
	"whats-app-clone-service/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func AuthRouter(router *mux.Router) {
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", middleware.RegisterUser).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/login", middleware.Login).Methods("POST", "OPTIONS")
}
