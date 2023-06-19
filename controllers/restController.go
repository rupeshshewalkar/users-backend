package controllers

import (
	"context"
	"encoding/json"
	"github.com/rupeshshewalkar/users-backend/helper"
	"github.com/rupeshshewalkar/users-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	collection := helper.ConnectToDB()

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.SetError(err, http.StatusInternalServerError, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	collection := helper.ConnectToDB()
	//decode req body into emp
	_ = json.NewDecoder(r.Body).Decode(&user)

	result, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		helper.SetError(err, http.StatusInternalServerError, w)
		return
	}
	//json return type
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
func UpdateUserByUserName(w http.ResponseWriter, r *http.Request) {
	var username models.User
	collection := helper.ConnectToDB()
	//decode req body into emp
	_ = json.NewDecoder(r.Body).Decode(&username)
	//except employeeID, all values can be updated
	filter := bson.M{"username": username.Username}
	result, err := collection.ReplaceOne(context.TODO(), filter, username)

	if err != nil {
		helper.SetError(err, http.StatusInternalServerError, w)
		return
	}
	//json return type
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func DeleteUserByUserName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("username")
	collection := helper.ConnectToDB()
	filter := bson.M{"username": name}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.SetError(err, http.StatusInternalServerError, w)
		return
	}
	//json return type
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Readiness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
