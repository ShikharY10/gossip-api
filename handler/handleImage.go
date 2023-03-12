package handler

import (
	"context"
	"time"

	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageHandler struct {
	postImage   *mongo.Collection
	avatarImage *mongo.Collection
	logger      *logger.Logger
}

func CreateImageHandler(postImage *mongo.Collection, avatarImage *mongo.Collection, logger *logger.Logger) *ImageHandler {
	handler := &ImageHandler{
		postImage:   postImage,
		avatarImage: avatarImage,
		logger:      logger,
	}
	return handler
}

func (iH *ImageHandler) SaveUserAvatar(userId primitive.ObjectID, imageData string, imageExt string) (string, error) {
	var image models.Image
	image.Id = primitive.NewObjectID()
	image.UserId = userId
	image.ImageData = imageData
	image.ImageExt = imageExt

	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	result, err := iH.avatarImage.InsertOne(ctx, image)
	if err != nil {
		return "", err
	} else {
		return result.InsertedID.(primitive.ObjectID).Hex(), nil
	}
}
