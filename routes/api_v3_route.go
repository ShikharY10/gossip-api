package routes

import (
	config "github.com/ShikharY10/gbAPI/config"
	controllers "github.com/ShikharY10/gbAPI/controllers"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/gorilla/mux"
)

func Version3API(mainRouter *mux.Router, m config.MAIN) {

	var cloudinary config.Cloudinary = config.InitCloudinary()

	var v2_authenticator middleware.JWT
	v2_authenticator.MsgModel = m.MsgModel
	v2_authenticator.Redis = m.RedisModel

	var api_v3 controllers.API_V3
	api_v3.MsgModel = m.MsgModel
	api_v3.UserModel = m.UserModel
	api_v3.RedisModel = m.RedisModel
	api_v3.RMQ = m.RMQ
	api_v3.AuthJwt = &v2_authenticator
	api_v3.Cloudinary = &cloudinary

	// API VERSION 2 Routes
	apiV3Router := mainRouter.PathPrefix("/api/v3").Subrouter()
	apiV3Router.HandleFunc("/", api_v3.Apiv3).Methods("GET")
	apiV3Router.HandleFunc("/upload/chat/image", api_v3.UploadChatImage).Methods("POST")
	apiV3Router.HandleFunc("/upload/chat/video", api_v3.UploadChatVideo).Methods("POST")
	apiV3Router.HandleFunc("/upload/profilepic", api_v3.UploadProfilePicture).Methods("POST")

}
