package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name             string               `bson:"name,omitempty" json:"name"`
	Username         string               `bson:"username,omitempty" json:"username"`
	Email            string               `bson:"email,omitempty" json:"email"`
	Avatar           Avatar               `bson:"avatar,omitempty" json:"avatar"`
	DeliveryId       primitive.ObjectID   `bson:"deliveryId,omitempty" json:"deliveryId"`
	Posts            []primitive.ObjectID `bson:"posts,omitempty" json:"posts"`
	Partners         []primitive.ObjectID `bson:"partners,omitempty" json:"partners"`
	PartnerRequests  []PartnerRequest     `bson:"partnerrequests,omitempty" json:"partnerrequests"`
	PartnerRequested []PartnerRequest     `bson:"partnerrequested,omitempty" json:"partnerrequested"`
	AccessToken      string               `bson:"accessToken,omitempty" json:"accessToken,omitempty"`
	Role             string               `bson:"role,omitempty" json:"role,omitempty"`
	CreatedAt        string               `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt        string               `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedAt        string               `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type Post struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	Media         Media              `bson:"media" json:"media"`
	Tags          []string           `bson:"tags" json:"tags"`
	Mentions      []string           `bson:"mentions" json:"mentions"`
	CreatedBy     primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt     string             `bson:"createdAt" json:"createdAt"`
	Likes         []Like             `bson:"likes" json:"likes"`
	NumOfLikes    int                `bson:"numOfLikes" json:"numOfLikes"`
	Comments      []Comment          `bson:"comments" json:"comments"`
	NumOfComments int                `bson:"numOfComments" json:"numOfComment"`
}

type PartnerRequest struct {
	ID                string `bson:"id" json:"id"`
	RequesterId       string `bson:"requesterId" json:"requesterId"`
	RequesterUsername string `bson:"requesterUsername" json:"requesterUsername"`
	RequesterName     string `bson:"requesterName" json:"requesterName"`
	TargetId          string `bson:"targetId" json:"targetId"`
	TargetUsername    string `bson:"targetUsername" json:"targetUsername"`
	TargetName        string `bson:"targetName" json:"targetName"`
	PublicKey         string `bson:"publicKey" json:"publicKey"`
	CreatedAt         string `bson:"createdAt" json:"createdAt"`
}

type PartnerResponse struct {
	ID          string `json:"id"`
	IsAccepted  bool   `json:"isAccepted"`
	ResponserId string `json:"responderId"`
	TargetId    string `json:"targetId"`
	SharedKey   string `json:"key"`
}

type RemovePartnerNotify struct {
	ID         string `json:"id"`
	NotifierId string `json:"notifierId"`
	TargetId   string `json:"targetId"`
}

type Partner struct {
	Id        string `bson:"Id" json:"Id"`
	Username  string `bson:"username" json:"username"`
	Name      string `bson:"name" json:"name"`
	CreatedAt string `bson:"createdAt" json:"creaetdAt"`
}

type Avatar struct {
	PublicId  string `json:"publicId" bson:"publicId"`
	SecureUrl string `json:"secureUrl" bson:"secureUrl"`
}

type AvatarPicker struct {
	Avatar Avatar `bson:"avatar" json:"avatar,omitempty"`
}

type Media struct {
	PublicId  string `bson:"publicId" json:"publicId"`
	FileName  string `bson:"fileName" json:"fileName"`
	FileType  string `bson:"fileType" json:"fileType"`
	SecureUrl string `bson:"secureUrl" json:"secureUrl"`
}

type Like struct {
	UserId string `bson:"userId" json:"userId"`
	Name   string `bson:"name" json:"name"`
	Avatar string `bson:"avatar" json:"avatar"`
}

type Comment struct {
	UserId  string `bson:"userId" json:"userId"`
	Content string `bson:"content" json:"content"`
	Name    string `bson:"name" json:"name"`
	Avatar  string `bson:"avatar" json:"avatar"`
}

type Image struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	UserId    primitive.ObjectID `bson:"userid" json:"userid"`
	ImageData string             `bson:"imagedata" json:"imagedata"`
	ImageExt  string             `bson:"imageext" json:"imageext"`
}

type FrequencyTable struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Username  string             `bson:"username" json:"username"`
	Frequency int                `bson:"frequency" json:"frequency"`
}

type Message struct {
	Id    primitive.ObjectID `bson:"_id" json:"_id"`
	Chats []Chat             `bson:"chats" json:"chats"`
}

type Chat struct {
	Id   string `bson:"chatkey" json:"chatkey"`
	Data string `bson:"data" json:"data"`
}

type Packet struct {
	NodeName string `json:"name"`
	Type     string `json:"type"`
	Message  string `json:"message"`
}

type Log struct {
	TimeStamp   string `bson:"timeStamp" json:"timeStamp"`
	ServiceType string `bson:"serviceType" json:"serviceType"`
	Type        string `bson:"type" json:"type"`
	FileName    string `bson:"fileName" json:"fileName"`
	LineNumber  int    `bson:"lineNumber" json:"lineNumber"`
	Message     string `bson:"errorMessage" json:"errorMessage"`
}
