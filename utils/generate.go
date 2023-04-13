package utils

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/ShikharY10/gbAPI/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Generator struct{}

func (f *Generator) GenerateUserData(userCollection *mongo.Collection, quantity int) {
	objectIDs := generateObjectID(quantity)
	for i := 0; i < quantity; i++ {
		createdAt := time.Now().Local().Format(time.RFC822)
		var user models.User
		user.Avatar = models.Avatar{
			PublicId:  randomString(10),
			SecureUrl: GenerateRandomId() + "/" + randomString(6),
		}
		user.CreatedAt = createdAt
		user.DeletedAt = ""
		user.Email = randomString(8) + "@gmail.com"
		user.PartnerRequested = []models.PartnerRequest{}
		user.PartnerRequests = []models.PartnerRequest{}
		user.Partners = []primitive.ObjectID{}
		user.ID = objectIDs[i]
		// user.Logout = rand.Int()%2 == 0
		user.Name = randomString(6) + " " + randomString(4)
		user.Role = "user"
		// user.Token = randomString(20)
		user.UpdatedAt = createdAt
		user.Username = randomString(5)
		userCollection.InsertOne(context.TODO(), user)
	}
}

func generateObjectID(quantity int) []primitive.ObjectID {
	var objectIDs []primitive.ObjectID = make([]primitive.ObjectID, quantity)
	for i := 0; i < quantity; i++ {
		objectID := primitive.NewObjectID()
		objectIDs = append(objectIDs, objectID)
	}
	return objectIDs
}

func randomString(size int) string {
	var str []string = make([]string, size)
	for i := 0; i < size; i++ {
		index := rand.Intn(25)
		str = append(str, string(small[index]))
	}
	return strings.Join(str, "")
}

var small string = "abcdefghijklmnopqrstuvwxyz"

// var generator utils.Generator
// generator.GenerateUserData(mongoDB_v3.Users, 100)
