package routes

import (
	"github.com/ShikharY10/gbAPI/config"
	controllers "github.com/ShikharY10/gbAPI/controllers"
	"github.com/gorilla/mux"
)

func Version1API(mainRouter *mux.Router, m config.MAIN) {

	var api_v1 controllers.API_V1
	api_v1.MsgModel = m.MsgModel
	api_v1.UserModel = m.UserModel
	api_v1.RedisModel = m.RedisModel
	api_v1.RMQ = m.RMQ

	apiV1Router := mainRouter.PathPrefix("/api/v1").Subrouter()

	apiV1Router.HandleFunc("/", api_v1.Apiv1).Methods("GET")
	apiV1Router.HandleFunc("/newuser", api_v1.NewUser).Methods("POST")
	apiV1Router.HandleFunc("/sendotp", api_v1.VerifyNumber).Methods("POST")
	apiV1Router.HandleFunc("/verifyotp", api_v1.VarifyNumberOTP).Methods("POST")
	apiV1Router.HandleFunc("/login", api_v1.LoginUser).Methods("POST")
	apiV1Router.HandleFunc("/logout", api_v1.LogOut).Methods("POST")
	apiV1Router.HandleFunc("/delete", api_v1.DeleteUser).Methods("POST")
	apiV1Router.HandleFunc("/toggleblock", api_v1.ToggleBlock).Methods("POST")
	apiV1Router.HandleFunc("/checkUser", api_v1.CheckAwailibity).Methods("POST")
	apiV1Router.HandleFunc("/removehs", api_v1.RemoveHandshake).Methods("POST")
	apiV1Router.HandleFunc("/uprofile", api_v1.UpdateProfilePicture).Methods("POST")
	apiV1Router.HandleFunc("/unumber", api_v1.UpdateNumber).Methods("POST")
	apiV1Router.HandleFunc("/uemail", api_v1.UpdateEmail).Methods("POST")
}
