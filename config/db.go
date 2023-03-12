package config

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	PostImage   *mongo.Collection
	AvatarImage *mongo.Collection
	Users       *mongo.Collection
	Posts       *mongo.Collection
	Frequnecy   *mongo.Collection
	Chats       *mongo.Collection
}

func MongoDBConnectV3(env *ENV) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var mongoClient *mongo.Client
	var err error

	if env.MongoDBConnectionMethod == "manual" {
		credential := options.Credential{
			Username: env.MongoDBUsername,
			Password: env.MongoDBPassword,
		}

		clientOptions := options.Client().ApplyURI("mongodb://" + env.MongoDBHost + ":" + env.MongoDBPort).SetAuth(credential)
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			defer cancel()
			return nil, err
		}
	} else if env.MongoDBConnectionMethod == "auto" {
		serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI(env.MongoDBConnectionString).SetServerAPIOptions(serverAPIOptions)
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			defer cancel()
			return nil, err
		}
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		defer cancel()
		return nil, err
	}

	var mongodb MongoDB

	images := mongoClient.Database("images")
	mongodb.PostImage = images.Collection("post")
	mongodb.AvatarImage = images.Collection("avatar")

	storage := mongoClient.Database("storage")
	mongodb.Users = storage.Collection("users")
	mongodb.Posts = storage.Collection("posts")
	mongodb.Frequnecy = storage.Collection("userFrequencyTable")

	messaging := mongoClient.Database("messages")
	mongodb.Chats = messaging.Collection("chats")

	defer cancel()
	return &mongodb, nil
}

func RedisInit(env *ENV) (*redis.Client, error) {
	options := redis.Options{
		Addr:     env.RedisHost + ":" + env.RedisPort,
		Password: "",
		DB:       0,
	}
	client := redis.NewClient(&options)
	ping := client.Ping()
	if ping.Err() != nil {
		return nil, ping.Err()
	}
	return client, nil
}
