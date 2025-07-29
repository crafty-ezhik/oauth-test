package routes

import (
	"github.com/crafty-ezhik/oauth-test/handlers"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(router *chi.Mux, handler handlers.AuthHandler) {
	router.Route("/auth", func(r chi.Router) {
		r.Get("/google/url", handler.GetGoogleAuthRedirectURI())
		r.Post("/google/callback", handler.GoogleCode())
	})
}
