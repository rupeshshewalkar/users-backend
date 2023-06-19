package main

import (
	"context"
	ctrl "github.com/rupeshshewalkar/users-backend/controllers"
	md "github.com/rupeshshewalkar/users-backend/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	if os.Getenv("MONGO_CONNECTION_URL") == "" {
		os.Setenv("MONGO_CONNECTION_URL", "mongodb://localhost:27017")
	}
}
func main() {
	var wait time.Duration
	mux := http.NewServeMux()

	getAllUsersHandler := http.HandlerFunc(ctrl.GetAllUsers)
	createNewUserHandler := http.HandlerFunc(ctrl.AddUser)
	updateUserByUserNameHandler := http.HandlerFunc(ctrl.UpdateUserByUserName)
	deleteUserByUserNameHandler := http.HandlerFunc(ctrl.DeleteUserByUserName)

	healthcheckHandler := http.HandlerFunc(ctrl.Healthcheck)
	readinessHandler := http.HandlerFunc(ctrl.Readiness)
	// Main Apis
	mux.Handle("/users/getAllUsers", md.LoggerMW(getAllUsersHandler))
	mux.Handle("/users/createNewUser", md.LoggerMW(md.HeaderValidatorMW(createNewUserHandler)))
	mux.Handle("/users/updateUserByUserName", md.LoggerMW(md.HeaderValidatorMW(updateUserByUserNameHandler)))
	mux.Handle("/users/deleteUserByUserName", md.LoggerMW(md.HeaderValidatorMW(deleteUserByUserNameHandler)))
	//container orch. support
	mux.Handle("/healthcheck", md.LoggerMW(healthcheckHandler))
	mux.Handle("/readiness", md.LoggerMW(readinessHandler))
	//server config
	srv := &http.Server{
		Addr:         ":8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}
	//start server
	go func() {
		log.Fatalln(srv.ListenAndServe())
	}()

	//graceful server shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
