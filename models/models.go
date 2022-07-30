package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Person Model
type Person struct {
	ObjectID  primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ID        string             `json:"id" bson:"-"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty" validate:"required,alpha"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty" validate:"required,alpha"`
}
