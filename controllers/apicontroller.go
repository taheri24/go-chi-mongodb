package controllers

import "github.com/go-chi/chi/v5"

type ApiControler interface {
	GetPrefix() string
	SetupRouter(router chi.Router)
	SetupOpenAPI()
}
