package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MsgModel struct {
	msgCollection *mongo.Collection
}

func CreateMsgModel(collection *mongo.Collection) MsgModel {
	mm := MsgModel{
		msgCollection: collection,
	}
	return mm
}

func (mm *MsgModel) AddUserMsgField() (string, error) {
	b := bson.M{
		"msg": bson.M{},
	}
	AM := mm.msgCollection
	res, err := AM.InsertOne(context.Background(), b)
	if err != nil {
		return "", err
	}
	_id := res.InsertedID.(primitive.ObjectID)
	return _id.Hex(), nil
}

func (mm *MsgModel) DeleteMsg(Tid string, MsgId string) {
	_id, _ := primitive.ObjectIDFromHex(Tid)
	r, err := mm.msgCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": _id},
		bson.M{"$unset": bson.M{"msg." + MsgId: 1}},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.MatchedCount)
}
