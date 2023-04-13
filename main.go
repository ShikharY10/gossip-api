package main

import (
	"log"

	"github.com/ShikharY10/gbAPI/api"
	"github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/handler"
	"github.com/ShikharY10/gbAPI/logger"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	ENV := config.LoadENV()

	logging, err := logger.InitializeLogger(ENV, "API")
	if err != nil {
		log.Fatal(err)
	}

	mongoDB, err := config.MongoDBConnectV3(ENV)
	if err != nil {
		log.Fatal(err)
	}

	redisDB, err := config.RedisInit(ENV)
	if err != nil {
		log.Fatal(err)
	}

	rabbitMQ, err := config.RabbitInit(ENV)
	if err != nil {
		log.Fatal(err)
	}

	handle := &handler.Handler{
		CacheHandler: handler.CreateCacheHandler(redisDB, logging),
		MsgHandler:   handler.CreateMsgHandler(mongoDB.Payloads, logging),
		QueueHandler: handler.CreateQueueHandler(rabbitMQ.Channel, logging),
		PostHandler:  handler.CreatePostHandler(mongoDB.Posts, mongoDB.PostImage, logging),
		UserHandler:  handler.CreateUserHandler(mongoDB.Users, mongoDB.AvatarImage, mongoDB.Frequnecy, logging),
		ImageHandler: handler.CreateImageHandler(mongoDB.PostImage, mongoDB.AvatarImage, logging),
	}

	api.StartVersionThreeAPIs(handle, ENV, logging)
}
