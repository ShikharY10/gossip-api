package handler

import (
	"github.com/ShikharY10/gbAPI/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostHandler struct {
	posts  *mongo.Collection
	images *mongo.Collection
	logger *logger.Logger
}

func CreatePostHandler(posts *mongo.Collection, images *mongo.Collection, logger *logger.Logger) *PostHandler {
	postHandler := &PostHandler{
		posts:  posts,
		images: images,
		logger: logger,
	}
	return postHandler
}
