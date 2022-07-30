package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/controllers"
)

// Routes -> define endpoints
func Routes(apiControlers []controllers.ApiControler) *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3600, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/debug", middleware.Profiler())
	r.Get("/auth", controllers.Auths)
	r.Route("/api", func(r chi.Router) {
		for _, apiCtrl := range apiControlers {
			r.Route(apiCtrl.GetPrefix(), apiCtrl.SetupRouter)
		}
	})
	r.Post("/upload", controllers.UploadFileEndpoint)
	return r
}
