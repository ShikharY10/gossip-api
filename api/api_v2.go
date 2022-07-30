package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ShikharY10/goAPI/gbp"
	"github.com/ShikharY10/goAPI/mongoAction"
	"github.com/ShikharY10/goAPI/redisAction"
	"github.com/ShikharY10/goAPI/rmq"
	"github.com/ShikharY10/goAPI/utils"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

type API_V2 struct {
	Mongo *mongoAction.Mongo
	Redis *redisAction.Redis
	RMQ   *rmq.RMQ
}

func (a *API_V2) SendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("service", "Gossip API")
	params := mux.Vars(r)
	number := params["number"]

	id, otp := a.Redis.RegisterOTP()
	var otpData map[string]string = map[string]string{
		"otp":    otp,
		"number": number,
	}
	b, _ := json.Marshal(otpData)
	a.RMQ.Produce("OTPd3hdzl8", b)

	var response gbp.Response
	response.Status = true
	response.Disc = "otp send"
	response.Data = id
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	}
	w.Write(responseBytes)
}

func (a *API_V2) VarifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("service", "Gossip API")
	var response gbp.Response

	if r.Body == nil {
		response.Status = false
		response.Disc = "empty body"
	} else {
		var __otp utils.VOTP
		_ = json.NewDecoder(r.Body).Decode(&__otp)
		res := a.Redis.VarifyOTP(__otp.Id, __otp.Otp)

		if res {
			response.Status = true
			response.Disc = "number varified"
		} else {
			response.Status = false
			response.Disc = "wrong otp"
		}
	}
	response.Data = ""
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	}
	w.Write(responseBytes)
}

func (a *API_V2) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("service", "Gossip API")

	var response gbp.Response
	response.Data = ""
	var body []byte = make([]byte, 10)
	i, _ := r.Body.Read(body)

	if i == 0 {
		response.Status = false
		response.Disc = "empty body"
	} else {
		var newUserdata utils.NewUser

		err := json.NewDecoder(r.Body).Decode(&newUserdata)
		if err != nil {
			response.Status = false
			response.Disc = "bad json data"
		}
		id, _ := a.Mongo.AddUserMsgField()
		newUserdata.MsgId = id

		aes_key := utils.GenerateAesKey(32)
		publicKey, err := utils.LoadKey([]byte(newUserdata.MainKey))
		if err != nil {
			response.Status = false
			response.Disc = "bad main key format"
		} else {
			cipherText, err := utils.RsaEncrypt(*publicKey, []byte(utils.Encode(aes_key)))

			if err != nil {
				response.Status = false
				response.Disc = "error while rsa encryption"
			} else {
				newUserdata.MainKey = string(aes_key)
				var uD utils.UserData
				uD.Age = newUserdata.Age
				uD.Blocked = map[string]int{newUserdata.MsgId: 1}
				uD.Connections = map[string]int{}
				uD.Email = newUserdata.Email
				uD.Gender = newUserdata.Gender
				uD.MainKey = utils.Encode(aes_key)
				uD.MsgId = newUserdata.MsgId
				uD.Name = newUserdata.Name
				uD.Password = newUserdata.Password
				uD.PhoneNo = newUserdata.PhoneNo
				uD.ProfilePic = newUserdata.ProfilePic
				uD.Logout = false

				uid, err := a.Mongo.AddUser(uD)
				if err != nil {
					response.Status = false
					response.Disc = "unable to add user data to the database"
				} else {
					var responseMap map[string]string = map[string]string{
						"uid":     uid,
						"mid":     newUserdata.MsgId,
						"Eaeskey": utils.Encode(cipherText),
					}

					res, err := json.Marshal(responseMap)
					if err != nil {
						response.Status = false
						response.Disc = "error while preparing response"
					} else {
						response.Status = true
						response.Disc = "success"
						response.Data = string(res)
					}
				}
			}
		}
	}

	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	}
	w.Write(responseBytes)
}
