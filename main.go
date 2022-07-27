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

func runAPIs(serverIP string, m *mongoAction.Mongo, r *redisAction.Redis, rmq *rmq.RMQ) {
	var api userAPI.API
	api.Mongo = m
	api.Redis = r
	api.RMQ = rmq
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/newuser", api.NewUser).Methods("POST")
	router.HandleFunc("/sendotp", api.VerifyNumber).Methods("POST")
	router.HandleFunc("/verifyotp", api.VarifyNumberOTP).Methods("POST")
	router.HandleFunc("/login", api.LoginUser).Methods("POST")
	router.HandleFunc("/logout", api.LogOut).Methods("POST")
	router.HandleFunc("/delete", api.DeleteUser).Methods("POST")
	router.HandleFunc("/toggleblock", api.ToggleBlock).Methods("POST")
	router.HandleFunc("/checkUser", api.CheckAwailibity).Methods("POST")
	router.HandleFunc("/removehs", api.RemoveHandshake).Methods("POST")
	router.HandleFunc("/uprofile", api.UpdateProfilePicture).Methods("POST")
	router.HandleFunc("/unumber", api.UpdateNumber).Methods("POST")
	router.HandleFunc("/uemail", api.UpdateEmail).Methods("POST")
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
