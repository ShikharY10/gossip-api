package main

import (
	"log"
	"net/http"
	"os"

	userAPI "github.com/ShikharY10/goAPI/api"
	"github.com/ShikharY10/goAPI/middleware"
	"github.com/ShikharY10/goAPI/mongoAction"
	"github.com/ShikharY10/goAPI/redisAction"
	"github.com/ShikharY10/goAPI/rmq"
	"github.com/gorilla/mux"
)

type MAIN struct {
	RedisDB *redisAction.Redis
	MongoDB *mongoAction.Mongo
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gossip API Test Home Route..."))
}

func runAPIs(serverIP string, m *mongoAction.Mongo, r *redisAction.Redis, rmq *rmq.RMQ) {
	var api_v1 userAPI.API_V1
	api_v1.Mongo = m
	api_v1.Redis = r
	api_v1.RMQ = rmq

	var api_v2 userAPI.API_V2
	api_v2.Mongo = m
	api_v2.Redis = r
	api_v2.RMQ = rmq

	mainRouter := mux.NewRouter().StrictSlash(true)

	// Main Routes
	mainRouter.HandleFunc("/", home).Methods("GET")

	// API VERSION 1 Routes
	apiV1Router := mainRouter.PathPrefix("/api/v1").Subrouter()
	apiV1Router.HandleFunc("/", api_v1.Apiv1).Methods("GET")
	apiV1Router.HandleFunc("/newuser", api_v1.NewUser).Methods("POST")
	apiV1Router.HandleFunc("/sendotp", api_v1.VerifyNumber).Methods("POST")
	apiV1Router.HandleFunc("/verifyotp", api_v1.VarifyNumberOTP).Methods("POST")
	apiV1Router.HandleFunc("/login", api_v1.LoginUser).Methods("POST")
	apiV1Router.HandleFunc("/logout", api_v1.LogOut).Methods("POST")
	apiV1Router.HandleFunc("/delete", api_v1.DeleteUser).Methods("POST")
	apiV1Router.HandleFunc("/toggleblock", api_v1.ToggleBlock).Methods("POST")
	apiV1Router.HandleFunc("/checkUser", api_v1.CheckAwailibity).Methods("POST")
	apiV1Router.HandleFunc("/removehs", api_v1.RemoveHandshake).Methods("POST")
	apiV1Router.HandleFunc("/uprofile", api_v1.UpdateProfilePicture).Methods("POST")
	apiV1Router.HandleFunc("/unumber", api_v1.UpdateNumber).Methods("POST")
	apiV1Router.HandleFunc("/uemail", api_v1.UpdateEmail).Methods("POST")

	// API VERSION 2 Routes
	apiV2Router := mainRouter.PathPrefix("/api/v2").Subrouter()
	apiV2Router.HandleFunc("/", api_v2.Apiv2).Methods("GET")
	apiV2Router.HandleFunc("/sendotp/{number}", api_v2.SendOTP).Methods("GET")
	apiV2Router.HandleFunc("/varifyotp", api_v2.VarifyOTP).Methods("POST")
	apiV2Router.HandleFunc("/createuser", api_v2.CreateNewUser).Methods("POST")
	apiV2Router.HandleFunc("/login", api_v2.Login).Methods("POST")

	// API VERSION 2 Secured Routes
	apiV2RouterSecure := apiV2Router.PathPrefix("/secure").Subrouter()
	var v2_authenticator middleware.JWT
	v2_authenticator.Mongo = m
	v2_authenticator.Redis = r
	apiV2RouterSecure.Use(v2_authenticator.APIV2Auth)
	apiV2RouterSecure.HandleFunc("/logout", api_v2.Logout).Methods("GET")
	apiV2RouterSecure.HandleFunc("/dashboard", api_v2.Dashboard).Methods("GET")
	apiV2RouterSecure.HandleFunc("/toggleblock", api_v2.ToggleBlock).Methods("POST")
	apiV2RouterSecure.HandleFunc("/checkuser", api_v2.CheckAwailibity).Methods("POST")
	apiV2RouterSecure.HandleFunc("/removehs", api_v2.RemoveFromHandshake).Methods("POST")
	apiV2RouterSecure.HandleFunc("/updatepic", api_v2.UpdateProfilePicture).Methods("POST")
	apiV2RouterSecure.HandleFunc("/updatenumber", api_v2.UpdateNumber).Methods("POST")
	apiV2RouterSecure.HandleFunc("/updateemail", api_v2.UpdateEmail).Methods("POST")

	log.Fatal(http.ListenAndServe(serverIP+":8080", mainRouter))
}

func main() {

	var mongoIP string = "127.0.0.1"
	var rabbitIP string = "127.0.0.1"
	var redisIP string = "127.0.0.1"
	var serverIP string = "0.0.0.0"

	var rabbitUsername string = "guest"
	var rabbitPassword string = "guest"

	var mongoUsername string = "rootuser"
	var mongoPassword string = "rootpass"

	// ---------------------------------------
	val, found := os.LookupEnv("MONGO_LOC_IP")
	if found {
		mongoIP = val
	}

	val, found = os.LookupEnv("RABBITMQ_LOC_IP")
	if found {
		rabbitIP = val
	}

	val, found = os.LookupEnv("REDIS_LOC_IP")
	if found {
		redisIP = val
	}
	// -----------------------------------------

	// -----------------------------------------
	val, found = os.LookupEnv("RABBITMQ_USERNAME")
	if found {
		rabbitUsername = val
	}

	val, found = os.LookupEnv("RABBITMQ_PASSWORD")
	if found {
		rabbitPassword = val
	}
	// -------------------------------------------

	// -------------------------------------------
	val, found = os.LookupEnv("MONGO_USERNAME")
	if found {
		mongoUsername = val
	}

	val, found = os.LookupEnv("MONGO_PASSWORD")
	if found {
		mongoPassword = val
	}
	// -------------------------------------------

	val, found = os.LookupEnv("SERVER_LOC_IP")
	if found {
		serverIP = val
	}

	var redisDB redisAction.Redis
	redisDB.Init(redisIP)
	var mongoDB mongoAction.Mongo
	mongoDB.Init(mongoIP, mongoUsername, mongoPassword)

	var RMQ rmq.RMQ
	RMQ.RedisDB = &redisDB
	RMQ.Init(rabbitIP, rabbitUsername, rabbitPassword)

	runAPIs(serverIP, &mongoDB, &redisDB, &RMQ)
}

// TOken: "MID+++EMAIL"
