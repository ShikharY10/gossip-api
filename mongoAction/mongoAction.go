package mongoAction

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ShikharY10/goAPI/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Ctx              context.Context
	Client           *mongo.Client
	UserCollection   *mongo.Collection
	MsgCollection    *mongo.Collection
	ServerCollection *mongo.Collection
}

func (m *Mongo) Init(mongoIP string, username string, password string) {
	var cred options.Credential
	cred.Username = username
	cred.Password = password

	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI("mongodb://" + mongoIP + ":27017").SetAuth(cred)
	client, err := mongo.Connect(m.Ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(m.Ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	var sKey string = "aeofh2983rh293hf2infd20i3fjd023j"

	uCollection := client.Database("Users").Collection("UserDatas")
	sCollection := client.Database("Users").Collection("info")
	mCollection := client.Database("messages").Collection("userMsg")
	m.Client = client
	m.Ctx = ctx
	m.UserCollection = uCollection
	m.ServerCollection = sCollection
	m.MsgCollection = mCollection

	m.AddServerDetails(sKey)

	fmt.Println("Mongo client connected!")
}

func (m *Mongo) AddServerDetails(key string) {
	_, err := m.Secretekey()
	if err != nil {
		var sData utils.ServerData
		sData.ID = "serverid"
		sData.SecreteKey = key

		m.ServerCollection.InsertOne(m.Ctx, sData)
	}
}

func (m *Mongo) AddUserMsgField() (string, error) {
	b := bson.M{
		"msg": bson.M{},
	}
	AM := m.MsgCollection
	res, err := AM.InsertOne(context.Background(), b)
	if err != nil {
		return "", err
	}
	_id := res.InsertedID.(primitive.ObjectID)
	return _id.Hex(), nil
}

func (m *Mongo) DeleteMsg(Tid string, MsgId string) {
	_id, _ := primitive.ObjectIDFromHex(Tid)
	r, err := m.MsgCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": _id},
		bson.M{"$unset": bson.M{"msg." + MsgId: 1}},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.MatchedCount)
}

func (m *Mongo) AddUser(user utils.UserData) (string, error) {

	result, err := m.UserCollection.InsertOne(m.Ctx, user)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("Adduser")
	id := result.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (m *Mongo) DeleteUser(filter primitive.M) error {
	_, err := m.UserCollection.DeleteOne(m.Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) DeleteUserByEmail(email string) error {
	filter := bson.M{"email": email}
	err := m.DeleteUser(filter)
	return err
}

func (m *Mongo) DeleteUserByPhoneNo(phoneno string) error {
	filter := bson.M{"phoneno": phoneno}
	err := m.DeleteUser(filter)
	return err
}

func (m *Mongo) UpdateUserName(id string, name string) (int64, error) {
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"name": name}}
	result, err := m.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return -1, nil
	}
	i := result.MatchedCount
	return i, nil
}

func (m *Mongo) UpdateUserAge(id int, age string) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"age": age}}
	result, err := m.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return -1, nil
	}
	i := result.MatchedCount
	return i, nil
}

func (m *Mongo) UpdateUserNumber(mid string, number string) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"phone_no": number}}
	result, err := m.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (m *Mongo) UpdateUserEmail(mid string, email string) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"email": email}}
	result, err := m.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (m *Mongo) UpdateUserProfilePic(mid string, profilePic string) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"profile_pic": profilePic}}
	result, err := m.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (m *Mongo) ReadUserData(filter primitive.M) (*utils.UserData, error) {
	cursor, err := m.UserCollection.Find(m.Ctx, filter)
	if err != nil {
		return nil, err
	}
	var userData []utils.UserData
	cursor.All(m.Ctx, &userData)
	return &userData[0], nil
}

func (m *Mongo) ReadUserDataById(id string) (*utils.UserData, error) {
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	udata, err := m.ReadUserData(filter)
	if err != nil {
		return nil, err
	}
	return udata, nil
}

func (m *Mongo) ReadUserDataByMNo(number string) (utils.UserData, error) {
	cursor, err := m.UserCollection.Find(
		context.TODO(),
		bson.M{"phone_no": number},
	)
	if err != nil {
		log.Println("[MONGOGETERROR] : ", err.Error())
	}
	var userd []utils.UserData
	err = cursor.All(context.TODO(), &userd)
	if err != nil {
		log.Println("[MONGOCURSORERROR] : ", err.Error())
	}
	return userd[0], nil
}

func (m *Mongo) ReadUserDataByMID(mid string) (*utils.UserData, error) {
	filter := bson.M{"msgid": mid}
	udata, err := m.ReadUserData(filter)
	if err != nil {
		return nil, err
	}
	return udata, nil
}

func (m *Mongo) GetUserDataByMID(target string) (*utils.UserData, error) {
	cursor, err := m.UserCollection.Find(context.TODO(), bson.M{"msgid": target})
	if err != nil {
		return nil, err
	}
	var ud []utils.UserData
	cursor.All(context.TODO(), &ud)
	return &ud[0], err
}

func (m *Mongo) AddTOBlocking(user string, target string) int {
	r, err := m.UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"msgid": user},
		bson.M{"$set": bson.M{"blocked." + target: 1}},
	)
	if err != nil {
		panic(err)
	}
	return int(r.MatchedCount)
}

func (m *Mongo) DeleteFromBlocking(user string, target string) int {
	r, err := m.UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"msgid": user},
		bson.M{"$unset": bson.M{"blocked." + target: 0}},
	)
	if err != nil {
		log.Fatal(err)
	}
	return int(r.MatchedCount)
}

func (m *Mongo) CheckBlocking(user string, target string) bool {
	cursor, err := m.UserCollection.Find(
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

func (m *Mongo) CheckAccountPresence(number string) bool {
	r := m.UserCollection.FindOne(
		context.TODO(),
		bson.M{"phone_no": number},
	)
	if r.Err() == nil {
		return true
	} else {
		return false
	}
}

func (m *Mongo) RemoveFromConnection(userMID string, connMID string) bool {
	r, err := m.UserCollection.UpdateOne(
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

func (m *Mongo) GetMsgIdByNum(mNum string) string {
	cursor, err := m.UserCollection.Find(
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

func (m *Mongo) GetNUMIdByMsgId(mid string) string {
	cursor, err := m.UserCollection.Find(
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

func (m *Mongo) UpdateLogoutStatus(mid string, status bool) bool {
	filter := bson.M{"msgid": mid}
	update := bson.M{"$set": bson.M{"logout": status}}
	result, err := m.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	if result.MatchedCount == 1 {
		return true
	}
	return false
}

func (m *Mongo) Secretekey() (string, error) {
	cursor, err := m.ServerCollection.Find(
		context.TODO(),
		bson.M{"id": "serverid"},
	)
	if err != nil {
		log.Println("[MONGOGETERROR] : ", err.Error())
	}
	var serverD []utils.ServerData
	err = cursor.All(context.TODO(), &serverD)
	if err != nil {
		log.Println("[MONGOCURSORERROR] : ", err.Error())
	}
	if len(serverD) > 0 {
		return serverD[0].SecreteKey, nil
	} else {
		return "", errors.New("no data found")
	}

}

func (m *Mongo) CheckUserExistence(email string) bool {
	cursor, err := m.UserCollection.Find(
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
