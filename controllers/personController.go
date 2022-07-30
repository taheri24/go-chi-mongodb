package controllers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/goava/di"
	"github.com/gorilla/mux"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/db"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/deps"
	middlewares "github.com/umangraval/Go-Mongodb-REST-boilerplate/handlers"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/models"
	"github.com/umangraval/Go-Mongodb-REST-boilerplate/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
	"os"
)

var database, _ = db.ConnectToDatabase()

// Auths -> get token
var Auths = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	validToken, err := middlewares.GenerateJWT()
	if err != nil {
		middlewares.ErrorResponse("Failed to generate token", response)
	}

	middlewares.SuccessResponse(string(validToken), response)
})

// CreatePersonEndpoint -> create person
var CreatePersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	var person models.Person
	err := json.NewDecoder(request.Body).Decode(&person)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	if ok, errors := validators.ValidateInputs(person); !ok {
		middlewares.ValidationResponse(errors, response)
		return
	}
	var collection *mongo.Collection
	deps.MainContainer.Resolve(&collection, di.Tags{"collection": "people"})

	result, err := collection.InsertOne(context.TODO(), person)

	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}

	id := models.SecureID("People", result.InsertedID.(primitive.ObjectID))
	middlewares.SuccessResponse(`Inserted at `+id, response)
})

// GetPeopleEndpoint -> get people
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	people := make([]models.Person, 0)

	var collection *mongo.Collection
	deps.MainContainer.Resolve(&collection, di.Tags{"collection": "people"})

	cursor, err := collection.Find(request.Context(), bson.D{{}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	for cursor.Next(request.Context()) {
		var person models.Person
		err := cursor.Decode(&person)
		if err != nil {
			log.Fatal(err)
		}
		person.ID = models.SecureID("People", person.ObjectID)
		people = append(people, person)

	}
	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(people)
}

// GetPersonEndpoint -> get person by id
var GetPersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person models.Person

	var collection *mongo.Collection
	deps.MainContainer.Resolve(&collection, di.Tags{"collection": "people"})

	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&person)
	if err != nil {
		middlewares.ErrorResponse("Person does not exist", response)
		return
	}
	middlewares.SuccessRespond(person, response)
})

// DeletePersonEndpoint -> delete person by id
var DeletePersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person models.Person

	collection := database.Collection("people")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&person)
	if err != nil {
		middlewares.ErrorResponse("Person does not exist", response)
		return
	}
	_, derr := collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if derr != nil {
		middlewares.ServerErrResponse(derr.Error(), response)
		return
	}
	middlewares.SuccessResponse("Deleted", response)
})

// UpdatePersonEndpoint -> update person by id
var UpdatePersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	type fname struct {
		Firstname string `json:"firstname"`
	}
	var fir fname
	json.NewDecoder(request.Body).Decode(&fir)
	collection := database.Collection("people")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "firstname", Value: fir.Firstname}}}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("Person does not exist", response)
		return
	}
	middlewares.SuccessResponse("Updated", response)
})

// UploadFileEndpoint -> upload file
var UploadFileEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	file, handler, err := request.FormFile("file")
	// fileName := request.FormValue("file_name")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	f, err := os.OpenFile("uploaded/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, _ = io.Copy(f, file)

	middlewares.SuccessResponse("Uploaded Successfully", response)
})

type PeopleApiController struct {
}

var _ ApiControler = &PeopleApiController{}

func (PeopleApiController) SetupRouter(r chi.Router) {
	r.Get("/", GetPeopleEndpoint)
	r.Post("/", CreatePersonEndpoint)
	r.Get("/{id}", GetPersonEndpoint)
	r.Delete("/{id}", DeletePersonEndpoint)
	r.Put("/{id}", UpdatePersonEndpoint)
}

func (PeopleApiController) GetPrefix() string {
	return "/people"
}

func (PeopleApiController) SetupOpenAPI() {

}
