package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	config "github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/ShikharY10/gbAPI/models"
)

type Response struct {
	Status int    `json:"status"`
	Disc   string `json:"disc"`
	Data   string `json:"data"`
}

type API_V3 struct {
	MsgModel   *models.MsgModel
	UserModel  *models.UserModel
	RedisModel *models.Redis
	RMQ        *config.RMQ
	AuthJwt    *middleware.JWT
	Cloudinary *config.Cloudinary
}

func (v3 *API_V3) Apiv3(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Version 3 APIs..."))
}

type imageUpload struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Ext   string `json:"ext"`
}

func (v3 *API_V3) UploadChatImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[V3] : Chat Image Upload")
	var response Response
	response.Data = ""
	// var body map[string]string
	var body imageUpload
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		response.Status = 1
		response.Disc = err.Error()
		fmt.Println("Error")
	} else {
		response.Data = "working"
		secureUrl, publicId, err := v3.Cloudinary.UploadChatImage(body.Image, body.Ext)
		if err != nil {
			fmt.Println("err: ", err)
			response.Status = 1
			response.Disc = err.Error()
			fmt.Println("Upload Error")
		} else {
			var resMap map[string]string = make(map[string]string)
			resMap["secureUrl"] = secureUrl
			resMap["publicId"] = publicId
			resBytes, err := json.Marshal(&resMap)
			if err != nil {
				response.Status = 1
				response.Disc = err.Error()
				fmt.Println("Response Parse Error")
			} else {
				response.Data = string(resBytes)
				fmt.Println("Image Upload Response Written!")
			}
		}
	}
	json.NewEncoder(w).Encode(response)
}

func (v3 *API_V3) UploadChatVideo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Version 3 APIs..."))
}

func (v3 *API_V3) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Version 3 APIs..."))
}
