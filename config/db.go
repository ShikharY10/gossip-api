package config

import (
	"context"
	"fmt"
	"log"

	"github.com/ShikharY10/gbAPI/models"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MAIN struct {
	MsgModel   *models.MsgModel
	UserModel  *models.UserModel
	RedisModel *models.Redis
	RMQ        *RMQ
}

type mongodb struct {
	// Ctx              context.Context
	// Client           *mongo.Client
	UserCollection   *mongo.Collection
	MsgCollection    *mongo.Collection
	ServerCollection *mongo.Collection
}

func MongoInit(mongoIP string, username string, password string) *mongodb {
	var m mongodb
	var cred options.Credential
	cred.Username = username
	cred.Password = password

	// ctx := context.TODO()
	clientOptions := options.Client().ApplyURI("mongodb://" + mongoIP + ":27017").SetAuth(cred)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// var sKey string = "aeofh2983rh293hf2infd20i3fjd023j"

	uCollection := client.Database("Users").Collection("UserDatas")
	mCollection := client.Database("messages").Collection("userMsg")
	// m.Client = client
	// m.Ctx = ctx
	m.UserCollection = uCollection
	m.MsgCollection = mCollection

	// m.AddServerDetails(sKey)

	fmt.Println("Mongo client connected!")

	return &m
}

type redisdb struct {
	Client *redis.Client
}

func RedisInit(redisIP string) *redisdb {
	var r redisdb
	client := redis.NewClient(&redis.Options{
		Addr:     redisIP + ":6379",
		Password: "",
		DB:       0,
	})
	s := client.Ping()
	fmt.Println(s.String())

	r.Client = client
	fmt.Println("Redis client connected!")
	return &r
}
