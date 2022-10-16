package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/models"
	"github.com/ShikharY10/gbAPI/routes"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gossip API Test Home Route..."))
}

func main() {

	// LOADING ENVIRONMENT VARIABLES
	godotenv.Load()

	mongoIP, found := os.LookupEnv("MONGO_LOC_IP")
	if !found {
		panic("key -> MONGO_LOC_IP is not found in .env")
	}

	rabbitIP, found := os.LookupEnv("RABBITMQ_LOC_IP")
	if !found {
		panic("key -> RABBITMQ_LOC_IP is not found in .env")
	}

	redisIP, found := os.LookupEnv("REDIS_LOC_IP")
	if !found {
		panic("key -> REDIS_LOC_IP is not found in .env")
	}

	rabbitUsername, found := os.LookupEnv("RABBITMQ_USERNAME")
	if !found {
		panic("key -> RABBITMQ_USERNAME is not found in .env")
	}

	rabbitPassword, found := os.LookupEnv("RABBITMQ_PASSWORD")
	if !found {
		panic("key -> RABBITMQ_PASSWORD is not found in .env")
	}

	mongoUsername, found := os.LookupEnv("MONGO_USERNAME")
	if !found {
		panic("key -> MONGO_USERNAME is not found in .env")
	}

	mongoPassword, found := os.LookupEnv("MONGO_PASSWORD")
	if !found {
		panic("key -> MONGO_PASSWORD is not found in .env")
	}

	serverIP, found := os.LookupEnv("SERVER_LOC_IP")
	if !found {
		panic("key -> SERVER_LOC_IP is not found in .env")
	}

	mongoDB := config.MongoInit(mongoIP, mongoUsername, mongoPassword)
	redisDB := config.RedisInit(redisIP)
	rabbitMQ := config.RabbitInit(rabbitIP, rabbitUsername, rabbitPassword)

	var redisModel models.Redis = models.CreateMainRedisModel(redisDB.Client)
	var userModel models.UserModel = models.CreateUserModel(mongoDB.UserCollection)
	var msgModel models.MsgModel = models.CreateMsgModel(mongoDB.MsgCollection)

	var m config.MAIN
	m.MsgModel = &msgModel
	m.UserModel = &userModel
	m.RedisModel = &redisModel
	m.RMQ = rabbitMQ

	mainRouter := mux.NewRouter().StrictSlash(true)
	// Test Home Routes
	mainRouter.HandleFunc("/", home).Methods("GET")

	routes.Version1API(mainRouter, m)
	routes.Version2API(mainRouter, m)
	routes.Version3API(mainRouter, m)

	log.Fatal(http.ListenAndServe(serverIP+":8080", mainRouter))
}
