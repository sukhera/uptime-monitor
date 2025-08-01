package routes

import (
	"github.com/gorilla/mux"
	"github.com/sukhera/uptime-monitor/internal/api/handlers"
	"github.com/sukhera/uptime-monitor/internal/api/middleware"
	"github.com/sukhera/uptime-monitor/internal/database"
)

func Setup(db *database.DB) *mux.Router {
	router := mux.NewRouter()
	
	statusHandler := handlers.NewStatusHandler(db)
	
	router.HandleFunc("/api/status", statusHandler.GetStatus).Methods("GET")
	router.HandleFunc("/api/health", statusHandler.HealthCheck).Methods("GET")
	
	return router
}

func WithMiddleware(router *mux.Router) *mux.Router {
	cors := middleware.NewCORS()
	handler := cors.Handler(router)
	
	wrappedRouter := mux.NewRouter()
	wrappedRouter.PathPrefix("/").Handler(handler)
	
	return wrappedRouter
}