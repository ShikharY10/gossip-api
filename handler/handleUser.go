package handler

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserHandler struct {
	userCollection      *mongo.Collection
	frequencyCollection *mongo.Collection
	logger              *logger.Logger
}

func CreateUserHandler(user *mongo.Collection, images *mongo.Collection, frequency *mongo.Collection, logger *logger.Logger) *UserHandler {
	userHandler := &UserHandler{
		userCollection:      user,
		frequencyCollection: frequency,
		logger:              logger,
	}
	return userHandler
}

// func (um UserHandler) AddUser(user utils.UserData) (string, error) {

// 	result, err := um.userCollection.InsertOne(context.TODO(), user)
// 	if err != nil {
// 		log.Fatal(err)
// 		return "", err
// 	}
// 	fmt.Println("Adduser")
// 	id := result.InsertedID.(primitive.ObjectID)
// 	return id.Hex(), nil
// }

func (uM *UserHandler) CreateNewUser(user models.User) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	_, err := uM.userCollection.InsertOne(ctx, user)
	return err
}

func (um *UserHandler) UpdateUserName(id string, name string) bool {
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"name": name}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	i := result.MatchedCount
	return i == 1
}

func (um *UserHandler) UpdateUserAvatar(id string, data models.Avatar) bool {
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"avatar": data}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	i := result.MatchedCount
	return i == 1
}

func (um UserHandler) UpdateUsername(id string, newUsername string) bool {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"username": newUsername}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	i := result.MatchedCount
	return i == 1
}

func (um UserHandler) UpdateUserEmail(id string, newEmail string) bool {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"email": newEmail}}
	result, err := um.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}
	i := result.MatchedCount
	return i == 1
}

func (um *UserHandler) GetUserEmail(id string) (string, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	filter := bson.M{"_id": _id}
	cursor, err := um.userCollection.Find(context.TODO(), filter)
	if err != nil {
		return "", err
	}

	var users []models.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return "", err
	}

	var user models.User = users[0]
	return user.Email, nil
}

func (um *UserHandler) IsPartners(id1 string, id2 string) bool {
	_id, err1 := primitive.ObjectIDFromHex(id1)
	p_id, err2 := primitive.ObjectIDFromHex(id2)
	if err1 != nil || err2 != nil {
		um.logger.LogError(err1)
		return false
	} else {
		result := um.userCollection.FindOne(
			context.TODO(),
			bson.M{"_id": _id, "partners": bson.A{p_id}},
		)
		if result.Err() != nil {
			return false
		} else {
			return true
		}
	}
}

func (um *UserHandler) GetUserDetails(id string, level string) (map[string]interface{}, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": _id}
	cursor, err := um.userCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return map[string]interface{}{}, errors.New("no user found")
	}
	if level == "L1" { // Only accessed by admins
		return map[string]interface{}{
			"user": users[0],
		}, nil
	} else if level == "L2" { // only accessed by partners
		return map[string]interface{}{
			"id":        users[0].ID,
			"name":      users[0].Name,
			"username":  users[0].Username,
			"avatar":    users[0].Avatar,
			"massageId": users[0].MessageID,
			"partners":  users[0].Partners,
			"posts":     users[0].Posts,
			"updatedAt": users[0].UpdatedAt,
		}, nil
	} else if level == "L3" { // accessed by anyone as soon as they are singed up.
		return map[string]interface{}{
			"id":           users[0].ID,
			"name":         users[0].Name,
			"username":     users[0].Username,
			"avatar":       users[0].Avatar.SecureUrl,
			"partnerCount": len(users[0].Partners),
			"postCount":    len(users[0].Posts),
		}, nil
	} else {
		return nil, errors.New("level not specified")
	}
}

// https://res.cloudinary.com/shikhar-lco/image/upload/h_800,w_600/c_scale,w_0.50/v1677949704/gb-profile_pic/pcmqsrqqrrn68skfx55k.jpg
// https://res.cloudinary.com/shikhar-lco/image/ upload /v1677949704/gb-profile_pic/pcmqsrqqrrn68skfx55k.jpg
func (um *UserHandler) GetUserAvatar(id string, width string, height string, scale string) ([]byte, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor := um.userCollection.FindOne(
		context.TODO(),
		bson.M{"_id": _id},
	)

	var avatarPicker models.AvatarPicker
	err = cursor.Decode(&avatarPicker)
	if err != nil {
		return nil, err
	}

	fmt.Println("SecureUrl: ", avatarPicker.Avatar.SecureUrl)

	splited := strings.Split(avatarPicker.Avatar.SecureUrl, "upload")
	url := splited[0] + "upload"

	if width != "" && height != "" {
		resizeParams := "/h_" + height + ",w_" + width
		url = url + resizeParams
	}

	if scale != "" {
		scaleParam := "/c_scale,w_" + scale
		url = url + scaleParam
	}

	url = url + splited[1]

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	return bodyBytes, nil

}

func (um *UserHandler) SetNewFollowRequest(id string, key string, detail models.PartnerRequest) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": _id}
	result, err := um.userCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.M{"$push": bson.M{key: detail}},
		// bson.M{key: bson.A{detail}},
	)
	if err != nil {
		return err
	}
	if result.ModifiedCount > int64(0) {
		return nil
	} else {
		return errors.New("error while updating user")
	}
}

func (um *UserHandler) SetNewPartner(id string, partnerId string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": _id}
	result, err := um.userCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.M{"$push": bson.M{"partners": partnerId}},
	)
	if err != nil {
		return err
	}
	if result.ModifiedCount > int64(0) {
		return nil
	} else {
		return errors.New("error while updating partner")
	}
}

func (um *UserHandler) RemovePartner(id string, partnerId string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": _id}
	result, err := um.userCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.M{"$pull": bson.M{"partners": partnerId}},
	)
	if err != nil {
		return err
	}
	if result.ModifiedCount > int64(0) {
		return nil
	} else {
		return errors.New("error while updating partner")
	}
}

func (um *UserHandler) AdminGetOneUser(username string) (models.User, bool) {
	var user models.User
	filter := bson.M{"username": username}
	cursor, err := um.userCollection.Find(context.TODO(), filter)
	if err != nil {
		return user, false
	}

	var users []models.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return user, false
	}

	user = users[0]
	return user, true
}

func (um *UserHandler) AdminGetAllUser() ([]models.User, error) {
	//find records
	//pass these options to the Find method
	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(30)
	//Define an array in which you can store the decoded documents
	var results []models.User

	//Passing the bson.D{{}} as the filter matches  documents in the collection

	cursor, err := um.userCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models.User
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	//Close the cursor once finished
	cursor.Close(context.TODO())
	return results, nil
}

// -----------------------------

func (um UserHandler) UpdateUserProfilePic(mid string, profilePic string) bool {
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

func (um UserHandler) UpdateUserNumber(mid string, number string) bool {
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

func (um UserHandler) GetUserDataByNumber(number string) (*models.User, error) {
	filter := bson.M{"phone_no": number}
	userData, err := readUserData(um, filter)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (um UserHandler) GetUserDataByMID(mid string) (*models.User, error) {
	filter := bson.M{"msgid": mid}
	userData, err := readUserData(um, filter)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func readUserData(um UserHandler, filter primitive.M) (*models.User, error) {
	cursor, err := um.userCollection.Find(
		context.TODO(),
		filter,
	)
	if err != nil {
		log.Println("[MONGOGETERROR] : ", err.Error())
		return nil, errors.New("no user found")
	}
	var userd []models.User
	err = cursor.All(context.TODO(), &userd)
	if err != nil {
		log.Println("[MONGOCURSORERROR] : ", err.Error())
		return nil, errors.New("no user found")
	}

	return &userd[0], nil
}

func (um UserHandler) AddTOBlocking(user string, target string) int {
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

func (um UserHandler) DeleteFromBlocking(user string, target string) int {
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

func (um UserHandler) CheckAccountPresence(username string) bool {
	r := um.userCollection.FindOne(
		context.TODO(),
		bson.M{"username": username},
	)
	if r.Err() == nil {
		return true
	} else {
		return false
	}
}

func (um UserHandler) RemoveFromConnection(userMID string, connMID string) bool {
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

func (um UserHandler) UpdateLogoutStatus(username string, status bool) bool {
	filter := bson.M{"username": username}
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

func (um *UserHandler) GetUserRole(username string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	result := um.userCollection.FindOne(
		ctx,
		bson.M{"username": username},
	)
	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return "", errors.New("no user found")
	} else {
		return user.Role, nil
	}
}

func (um *UserHandler) SearchUsername(username string) ([]models.FrequencyTable, error) {
	cursor, err := um.frequencyCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var users []models.FrequencyTable
	for cursor.Next(context.TODO()) {
		var elem models.FrequencyTable
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, elem)
	}

	var finded []models.FrequencyTable
	for _, user := range users {
		if strings.Contains(user.Username, username) {
			finded = append(finded, user)
		}
	}
	sort.Slice(finded, func(i, j int) bool {
		return finded[i].Frequency > finded[j].Frequency
	})

	return finded, nil
}

func (um *UserHandler) InsetUserInFrequencyTable(id primitive.ObjectID, username string) error {
	fTable := models.FrequencyTable{
		Id:        id,
		Username:  username,
		Frequency: 0,
	}
	_, err := um.frequencyCollection.InsertOne(
		context.TODO(),
		fTable,
	)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (um *UserHandler) IncrementUserInFrequencyTable(username string) error {
	filter := bson.M{"username": username}
	result, err := um.frequencyCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.M{"$inc": bson.M{"frequency": 1}},
	)
	if err != nil {
		return err
	}
	if result.ModifiedCount > int64(0) {
		return nil
	}
	return errors.New("no user found for this username, " + username)
}

func (um *UserHandler) IsUsernameAwailable(username string) error {
	cursor, err := um.frequencyCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		um.logger.LogError(err)
		return err
	}
	var users []models.FrequencyTable
	for cursor.Next(context.TODO()) {
		var elem models.FrequencyTable
		err := cursor.Decode(&elem)
		if err != nil {
			um.logger.LogError(err)
			return err
		}
		users = append(users, elem)
	}
	for _, user := range users {
		if user.Username == username {
			return errors.New("username already present")
		}
	}
	return nil
}
