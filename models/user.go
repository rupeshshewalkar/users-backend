package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty"`
	Alias    string             `json:"alias" bson:"alias,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Username string             `json:"username" bson:"username,omitempty"`
}
