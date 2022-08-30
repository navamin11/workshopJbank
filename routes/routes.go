package routes

import (
	"net/http"
	"time"
	"workshop/controllers"
	"workshop/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var r *chi.Mux

func SetupRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", controllers.HealthCheckGetHandler)
		r.Post("/register", controllers.RegisterPostHandler)
		r.Post("/login", controllers.LoginPostHandler)
	})

	r.Mount("/api/v1/app", memberRouter())
	r.Mount("/api/v1/thirdparty", thirdpartyRouter())

	return r
}

func memberRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.JwtAuthMiddleware)
	r.Get("/home", controllers.HomeGetHandler)
	r.Delete("/logout", controllers.LogoutDelHandler)
	r.Get("/account", controllers.DepositGetHandler)
	r.Get("/accounts", controllers.DepositListGetHandler)
	r.Post("/banking/{type}", controllers.TransferOutPostHandler)

	return r
}

func thirdpartyRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/{type}", controllers.TransferInPostHandler)

	return r
}
