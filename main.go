package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/goava/di"
	"github.com/sirupsen/logrus"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/controllers"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/db"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/deps"
	middlewares "github.com/umangraval/Go-Mongodb-REST-boilerplate/handlers"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	checkHealth()
	di.SetTracer(&di.StdTracer{})
	mainContainer, _ := di.New(
		controllers.Options,
		di.Provide(db.ConnectToDatabase),
		db.Collections,
		di.Provide(routes.Routes),
		di.Provide(deps.NewDefaultLogger),
	)
	deps.SetDiContainer(mainContainer)
	mainContainer.Invoke(startHttpServer)
}
func checkHealth() {
	logger := log.Default()
	log.Println("Connecting to Database...")
	_, err := db.ConnectToDatabase()
	if err != nil {
		logger.Fatalf("Connection Failed to Database", err)
		os.Exit(100)
		return
	}
}
func startHttpServer(logger *logrus.Logger, mux *chi.Mux) {
	port := middlewares.DotEnvVariable("PORT")

	logger.WithField("port", port).Info("Server running on localhost ")
	http.ListenAndServe(":"+port, mux)
}
