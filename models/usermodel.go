package models

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ShikharY10/gbAPI/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserData struct {
	MsgId       string         `bson:"msgid,omitempty"`
	Name        string         `bson:"name,omitempty"`
	Age         string         `bson:"age,omitempty"`
	PhoneNo     string         `bson:"phone_no,omitempty"`
	Email       string         `bson:"email,omitempty"`
	ProfilePic  string         `bson:"profile_pic,omitempty"`
	MainKey     string         `bson:"main_key,omitempty"`
	Gender      string         `bson:"gender,omitempty"`
	Password    string         `bson:"password,omitempty"`
	Logout      bool           `bson.A:"logout,omitempty"`
	Blocked     map[string]int `bson.M:"blocked,omitempty"`
	Connections map[string]int `bson.A:"connections,omitempty"`
}

type UserModel struct {
	userCollection *mongo.Collection
}

func CreateUserModel(collection *mongo.Collection) UserModel {
	s := UserModel{
		userCollection: collection,
	}
	return s
}

func (um *UserModel) AddUser(user utils.UserData) (string, error) {

	result, err := um.userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("Adduser")
	id := result.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (um *UserModel) UpdateUserName(id string, name string) (int64, error) {
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"name": name}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return -1, nil
	}
	i := result.MatchedCount
	return i, nil
}

func (um *UserModel) UpdateUserAge(id int, age string) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"age": age}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return -1, nil
	}
	i := result.MatchedCount
	return i, nil
}

func (um *UserModel) UpdateUserNumber(mid string, number string) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"phone_no": number}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (um *UserModel) UpdateUserEmail(mid string, email string) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"email": email}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (um *UserModel) UpdateUserProfilePic(mid string, profilePic string) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"profile_pic": profilePic}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (um *UserModel) GetUserDataByNumber(number string) (*UserData, error) {
	filter := bson.M{"phone_no": number}
	userData, err := readUserData(um, filter)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (um *UserModel) GetUserDataByMID(mid string) (*UserData, error) {
	filter := bson.M{"msgid": mid}
	userData, err := readUserData(um, filter)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func readUserData(um *UserModel, filter primitive.M) (*UserData, error) {
	cursor, err := um.userCollection.Find(
		context.TODO(),
		filter,
	)
	if err != nil {
		log.Println("[MONGOGETERROR] : ", err.Error())
		return nil, errors.New("no user found")
	}
	var userd []UserData
	err = cursor.All(context.TODO(), &userd)
	if err != nil {
		log.Println("[MONGOCURSORERROR] : ", err.Error())
		return nil, errors.New("no user found")
	}

	return &userd[0], nil
}

func (um *UserModel) AddTOBlocking(user string, target string) int {
	r, err := um.userCollection.UpdateOne(
		context.TODO(),
		bson.M{"msgid": user},
		bson.M{"$set": bson.M{"blocked." + target: 1}},
	)
	if err != nil {
		panic(err)
	}
	return int(r.MatchedCount)
}

func (um *UserModel) DeleteFromBlocking(user string, target string) int {
	r, err := um.userCollection.UpdateOne(
		context.TODO(),
		bson.M{"msgid": user},
		bson.M{"$unset": bson.M{"blocked." + target: 0}},
	)
	if err != nil {
		log.Fatal(err)
	}
	return int(r.MatchedCount)
}

func (um *UserModel) CheckBlocking(user string, target string) bool {
	cursor, err := um.userCollection.Find(
		context.TODO(),
		bson.M{"msgid": user},
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	var userData []utils.UserData
	cursor.All(
		context.TODO(),
		&userData,
	)
	if userData[0].Blocked[target] != 0 {
		return true
	} else {
		return false
	}
}

func (um *UserModel) CheckAccountPresence(number string) bool {
	r := um.userCollection.FindOne(
		context.TODO(),
		bson.M{"phone_no": number},
	)
	if r.Err() == nil {
		return true
	} else {
		return false
	}
}

func (um *UserModel) RemoveFromConnection(userMID string, connMID string) bool {
	r, err := um.userCollection.UpdateOne(
		context.TODO(),
		bson.M{"msgid": userMID},
		bson.M{"$unset": bson.M{"connections." + connMID: 1}},
	)
	if err != nil {
		log.Fatal(err)
	}
	if r.MatchedCount == 1 {
		return true
	} else {
		return false
	}
}

func (um *UserModel) GetMsgIdByNumber(mNum string) string {
	cursor, err := um.userCollection.Find(
		context.TODO(),
		bson.M{"phone_no": mNum},
	)
	if err != nil {
		log.Println("[MONGOGETERROR] : ", err.Error())
	}
	var userd []utils.UserData
	err = cursor.All(context.TODO(), &userd)
	if err != nil {
		log.Println("[MONGOCURSORERROR] : ", err.Error())
	}
	return userd[0].MsgId
}

func (um *UserModel) GetNumberByMsgId(mid string) string {
	cursor, err := um.userCollection.Find(
		context.TODO(),
		bson.M{"msgid": mid},
	)
	if err != nil {
		log.Println("[MONGOGETERROR] : ", err.Error())
	}
	var userd []utils.UserData
	err = cursor.All(context.TODO(), &userd)
	if err != nil {
		log.Println("[MONGOCURSORERROR] : ", err.Error())
	}
	return userd[0].PhoneNo
}

func (um *UserModel) UpdateLogoutStatus(mid string, status bool) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"logout": status}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (um *UserModel) CheckUserExistence(email string) bool {
	cursor, err := um.userCollection.Find(
		context.TODO(),
		bson.M{"email": email},
	)
	if err != nil {
		log.Println("[MONGOGETERROR] : ", err.Error())
	}
	var userd []utils.UserData
	err = cursor.All(context.TODO(), &userd)
	if err != nil {
		log.Println("[MONGOCURSORERROR] : ", err.Error())
	}
	if len(userd) > 0 {
		return userd[0].PhoneNo != ""
	}
	return false

}
