package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	userAPI "github.com/ShikharY10/goAPI/api"
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

func apiv1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Version 1 APIs..."))
}

func apiv2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Version 2 APIs..."))
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

	router := mux.NewRouter().StrictSlash(true)

	// Routes
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/api/v1", apiv1).Methods("GET")
	router.HandleFunc("/api/v1/newuser", api_v1.NewUser).Methods("POST")
	router.HandleFunc("/api/v1/sendotp", api_v1.VerifyNumber).Methods("POST")
	router.HandleFunc("/api/v1/verifyotp", api_v1.VarifyNumberOTP).Methods("POST")
	router.HandleFunc("/api/v1/login", api_v1.LoginUser).Methods("POST")
	router.HandleFunc("/api/v1/logout", api_v1.LogOut).Methods("POST")
	router.HandleFunc("/api/v1/delete", api_v1.DeleteUser).Methods("POST")
	router.HandleFunc("/api/v1/toggleblock", api_v1.ToggleBlock).Methods("POST")
	router.HandleFunc("/api/v1/checkUser", api_v1.CheckAwailibity).Methods("POST")
	router.HandleFunc("/api/v1/removehs", api_v1.RemoveHandshake).Methods("POST")
	router.HandleFunc("/api/v1/uprofile", api_v1.UpdateProfilePicture).Methods("POST")
	router.HandleFunc("/api/v1/unumber", api_v1.UpdateNumber).Methods("POST")
	router.HandleFunc("/api/v1/uemail", api_v1.UpdateEmail).Methods("POST")

	router.HandleFunc("/api/v2", apiv2).Methods("GET")
	router.HandleFunc("/api/v2/sendotp/{number}", api_v2.SendOTP).Methods("GET")
	router.HandleFunc("/api/v2/varifyotp", api_v2.VarifyOTP).Methods("POST")
	router.HandleFunc("/api/v2/createnewuser", api_v2.CreateNewUser).Methods("POST")

	log.Fatal(http.ListenAndServe(serverIP+":8080", router))
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
	fmt.Println("Server Started")

}
