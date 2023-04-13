package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/ShikharY10/gbAPI/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MsgHandler struct {
	msgCollection *mongo.Collection
	logger        *logger.Logger
}

func CreateMsgHandler(collection *mongo.Collection, logger *logger.Logger) *MsgHandler {
	return &MsgHandler{
		msgCollection: collection,
		logger:        logger,
	}

}

func (hm *MsgHandler) AddUserField() (*primitive.ObjectID, error) {
	b := bson.M{
		"msg": bson.M{},
	}
	AM := hm.msgCollection
	res, err := AM.InsertOne(context.TODO(), b)
	if err != nil {
		return nil, err
	}
	_id := res.InsertedID.(primitive.ObjectID)
	return &_id, nil
}

func (hm *MsgHandler) DeleteMsg(Tid string, MsgId string) {
	_id, _ := primitive.ObjectIDFromHex(Tid)
	r, err := hm.msgCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": _id},
		bson.M{"$unset": bson.M{"msg." + MsgId: 1}},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.MatchedCount)
}
