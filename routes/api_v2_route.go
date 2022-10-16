package routes

import (
	"github.com/ShikharY10/gbAPI/config"
	controllers "github.com/ShikharY10/gbAPI/controllers"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/gorilla/mux"
)

func Version2API(mainRouter *mux.Router, m config.MAIN) {

	var v2_authenticator middleware.JWT
	v2_authenticator.MsgModel = m.MsgModel
	v2_authenticator.Redis = m.RedisModel

	var api_v2 controllers.API_V2
	api_v2.MsgModel = m.MsgModel
	api_v2.UserModel = m.UserModel
	api_v2.RedisModel = m.RedisModel
	api_v2.RMQ = m.RMQ
	api_v2.AuthJwt = &v2_authenticator

	// API VERSION 2 Routes
	apiV2Router := mainRouter.PathPrefix("/api/v2").Subrouter()
	apiV2Router.HandleFunc("/", api_v2.Apiv2).Methods("GET")
	apiV2Router.HandleFunc("/sendotp/{number}", api_v2.SendOTP).Methods("GET")
	apiV2Router.HandleFunc("/varifyotp", api_v2.VarifyOTP).Methods("POST")
	apiV2Router.HandleFunc("/createuser", api_v2.CreateNewUser).Methods("POST")
	apiV2Router.HandleFunc("/login", api_v2.Login).Methods("POST")

	// secure
	apiV2RouterSecure := apiV2Router.PathPrefix("/secure").Subrouter()

	apiV2RouterSecure.Use(v2_authenticator.APIV2Auth)
	apiV2RouterSecure.HandleFunc("/logout", api_v2.Logout).Methods("GET")
	apiV2RouterSecure.HandleFunc("/dashboard", api_v2.Dashboard).Methods("GET")
	apiV2RouterSecure.HandleFunc("/toggleblock", api_v2.ToggleBlock).Methods("POST")
	apiV2RouterSecure.HandleFunc("/checkuser", api_v2.CheckAwailibity).Methods("POST")
	apiV2RouterSecure.HandleFunc("/removehs", api_v2.RemoveFromHandshake).Methods("POST")
	apiV2RouterSecure.HandleFunc("/updatepic", api_v2.UpdateProfilePicture).Methods("POST")
	apiV2RouterSecure.HandleFunc("/updatenumber", api_v2.UpdateNumber).Methods("POST")
	apiV2RouterSecure.HandleFunc("/updateemail", api_v2.UpdateEmail).Methods("POST")

}
