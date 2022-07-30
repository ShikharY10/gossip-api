package api

import (
	"encoding/json"
	"fmt"
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

type API_V1 struct {
	Mongo *mongoAction.Mongo
	Redis *redisAction.Redis
	RMQ   *rmq.RMQ
}

func (a *API_V1) VerifyNumber(w http.ResponseWriter, r *http.Request) {
	mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Status = "unsuccessful"
		s.Disc = "Please send mobile number!"
		json.NewEncoder(w).Encode(s)
	}
	var mn utils.MobileNo
	_ = json.NewDecoder(r.Body).Decode(&mn)
	fmt.Println("mn: ", mn)
	id, otp := a.Redis.RegisterOTP()
	var otpData map[string]string = map[string]string{
		"otp":    otp,
		"number": mn.Number,
	}
	b, _ := json.Marshal(otpData)
	a.RMQ.Produce("OTPd3hdzl8", b)
	// utils.SendOTP(mn.Number, otp)
	var opid utils.OperationId
	opid.Id = id
	json.NewEncoder(w).Encode(opid)
}

func (a *API_V1) VarifyNumberOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Status = "unsuccessful"
		s.Disc = "No OTP is provided"
		json.NewEncoder(w).Encode(s)
	}
	var __otp utils.VOTP
	_ = json.NewDecoder(r.Body).Decode(&__otp)
	res := a.Redis.VarifyOTP(__otp.Id, __otp.Otp)
	if res {
		var s utils.SuccessStruct
		s.Status = "successful"
		s.Disc = ""
		json.NewEncoder(w).Encode(s)
	}
}

func (a *API_V1) NewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send user data!"
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	var newUserdata utils.NewUser
	err := json.NewDecoder(r.Body).Decode(&newUserdata)
	if err != nil {
		log.Println(err.Error())
	}
	id, _ := a.Mongo.AddUserMsgField()
	newUserdata.MsgId = id

	aes_key := utils.GenerateAesKey(32)
	fmt.Println([]byte(newUserdata.MainKey))
	publicKey, err := utils.LoadKey([]byte(newUserdata.MainKey))
	if err != nil {
		log.Println("[PUBKEYLDERROR] : ", err.Error())
	}
	cipherText, err := utils.RsaEncrypt(*publicKey, []byte(utils.Encode(aes_key)))

	if err != nil {
		log.Println("[RSAENCRYPTERROR] : ", err.Error())
	}
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

	// fmt.Println("uD: ", uD)
	uid, err := a.Mongo.AddUser(uD)
	if err != nil {
		log.Println("[MONGOADDUSERERROR] : ", err.Error())
	}

	var response map[string]string = map[string]string{
		"uid":     uid,
		"mid":     newUserdata.MsgId,
		"Eaeskey": utils.Encode(cipherText),
	}

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("[JSONMARSHALERROR] : ", err.Error())
	}

	var s utils.SuccessStruct
	s.Status = "successful"
	s.Disc = utils.Encode(res)
	json.NewEncoder(w).Encode(s)
	a.Redis.Client.Del(newUserdata.MsgId)
	// fmt.Println("Len: ", len(s.Disc), " | Ciphertext: ", s.Disc)
}

func (a *API_V1) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send user data!"
		s.Status = "Unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	var loginData utils.LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&loginData)

	myData, err := a.Mongo.ReadUserDataByMNo(loginData.Number)
	if err != nil {
		var s utils.SuccessStruct
		s.Status = "unsuccessful"
		s.Disc = "Mobile number not matched"
		json.NewEncoder(w).Encode(s)
		return
	}

	if myData.Password == loginData.Password {
		var loginResponsePayload gbp.LoginResponse
		var loginEnginePayload gbp.LoginEnginePayload

		var myDataPayload gbp.UserData
		myDataPayload.Dob = myData.Age
		myDataPayload.Email = myData.Email
		myDataPayload.Gender = myData.Gender
		myDataPayload.MainKey = myData.MainKey
		myDataPayload.Mid = myData.MsgId
		myDataPayload.Name = myData.Name
		myDataPayload.Number = myData.PhoneNo
		myDataPayload.ProfilePic = myData.ProfilePic

		loginResponsePayload.MyData = &myDataPayload

		loginEnginePayload.SenderMid = myData.MsgId
		loginEnginePayload.PublicKey = loginData.PublicKey

		var allConn []string
		var connDatalist []*gbp.ConnectionData = []*gbp.ConnectionData{}
		for mid := range myData.Connections {

			connData, err := a.Mongo.GetUserDataByMID(mid)
			if err != nil {
				continue
			}

			if !connData.Logout {
				allConn = append(allConn, mid)
			}

			var connDataPayload gbp.ConnectionData
			connDataPayload.Mid = connData.MsgId
			connDataPayload.Name = connData.Name
			connDataPayload.Number = connData.PhoneNo
			connDataPayload.ProfilePic = connData.ProfilePic
			connDataPayload.Logout = connData.Logout
			connDatalist = append(connDatalist, &connDataPayload)
		}

		loginResponsePayload.ConnData = connDatalist
		loginEnginePayload.AllConn = allConn

		enginePayloadBytes, err := proto.Marshal(&loginEnginePayload)
		if err != nil {
			var s utils.SuccessStruct
			s.Status = "unsuccessful"
			s.Disc = "something went wrong"
			json.NewEncoder(w).Encode(s)
			return
		}

		var trans gbp.Transport
		trans.Id = ""
		trans.Msg = enginePayloadBytes
		trans.Tp = 10
		transBytes, err := proto.Marshal(&trans)
		if err != nil {
			var s utils.SuccessStruct
			s.Disc = "proto marshal error"
			s.Status = "unsuccessful"
			json.NewEncoder(w).Encode(s)
			return
		}
		engineName := a.RMQ.GetEngineChannel()
		a.RMQ.Produce(engineName, transBytes)

		userPayloadBytes, err := proto.Marshal(&loginResponsePayload)
		if err != nil {
			var s utils.SuccessStruct
			s.Status = "unsuccessful"
			s.Disc = "something went wrong"
			json.NewEncoder(w).Encode(s)
			return
		}

		var response gbp.Response
		response.Status = true
		response.Disc = "ok"
		response.Data = utils.Encode(userPayloadBytes)
		responseBytes, err := proto.Marshal(&response)
		if err != nil {
			log.Println("[marshal error]", err.Error())
		}
		w.Write(responseBytes)
		a.Mongo.UpdateLogoutStatus(myData.MsgId, false)
		return
	}
}

func (a *API_V1) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send user data!"
		s.Status = "Unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	var opid utils.DeletePayload
	_ = json.NewDecoder(r.Body).Decode(&opid)
	o := a.Redis.Client.Get(opid.OpId)
	if o.Err() != nil {
		var s utils.SuccessStruct
		s.Status = "unsuccessful"
		s.Disc = "Session Timeout"
		json.NewEncoder(w).Encode(s)
		return
	}
	var err error = nil
	if opid.TargetType == "email" {
		err = a.Mongo.DeleteUserByEmail(opid.Target)
	} else if opid.TargetType == "phoneno" {
		err = a.Mongo.DeleteUserByPhoneNo(opid.Target)
	}
	var s utils.SuccessStruct
	if err != nil {
		s.Status = "Unsuccessful"
	}
	s.Status = "Successfull"
	json.NewEncoder(w).Encode(s)
	a.Redis.Client.Del(opid.OpId)
}

func (a *API_V1) ToggleBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send relevent data!"
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	var blockRequest utils.ToggleBlocking
	err := json.NewDecoder(r.Body).Decode(&blockRequest)
	if err != nil {
		log.Println(err.Error())
	}

	userData, err := a.Mongo.GetUserDataByMID(blockRequest.SenderMID)
	if err != nil {
		var s utils.SuccessStruct
		s.Disc = "invalid senderMID / user not found"
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	tp := userData.Blocked[blockRequest.TargetNUM]

	if blockRequest.Type == -1 {
		if tp == 1 {
			res := a.Mongo.DeleteFromBlocking(blockRequest.SenderMID, blockRequest.TargetNUM)
			var s utils.SuccessStruct
			if res == 1 {
				s.Disc = "Deleted"
				s.Status = "successful"
				json.NewEncoder(w).Encode(s)
				return
			} else if res == 0 {
				s.Disc = "Something went wrong!"
				s.Status = "unsuccessful"
				json.NewEncoder(w).Encode(s)
				return
			}
		}
	} else if blockRequest.Type == 1 {
		if tp != 1 {
			res := a.Mongo.AddTOBlocking(blockRequest.SenderMID, blockRequest.TargetNUM)
			var s utils.SuccessStruct
			if res == 1 {
				s.Disc = "Added"
				s.Status = "successful"
				json.NewEncoder(w).Encode(s)
				return
			} else if res == 0 {
				s.Disc = "Something went wrong!"
				s.Status = "unsuccessful"
				json.NewEncoder(w).Encode(s)
				return
			}
		}
	}
}

func (a *API_V1) CheckAwailibity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send relevent data!"
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	var checkAccReq utils.CheckAccPayload
	err := json.NewDecoder(r.Body).Decode(&checkAccReq)
	if err != nil {
		log.Println(err.Error())
	}
	found := a.Mongo.CheckAccountPresence(checkAccReq.TargetNUM)
	if !found {
		var s utils.SuccessStruct
		s.Disc = "no account found related to this number"
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	} else {
		bFound := a.Mongo.CheckBlocking(checkAccReq.SenderMID, checkAccReq.TargetNUM)
		if bFound {
			var s utils.SuccessStruct
			s.Disc = "you are blocked!"
			s.Status = "unsuccessful"
			json.NewEncoder(w).Encode(s)
			return
		} else {
			var s utils.SuccessStruct
			s.Disc = ""
			s.Status = "successful"
			json.NewEncoder(w).Encode(s)
			return
		}
	}
}

func (a *API_V1) RemoveHandshake(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send relevent data!"
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	var hsRemove utils.HandshakeDeletePaylaod
	err := json.NewDecoder(r.Body).Decode(&hsRemove)
	if err != nil {
		var s utils.SuccessStruct
		s.Disc = "Wrong data format, expecting json but not found."
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
	targetMId := a.Mongo.GetMsgIdByNum(hsRemove.TargetNUM)
	senderNUM := a.Mongo.GetNUMIdByMsgId(hsRemove.UserMID)
	res1 := a.Mongo.RemoveFromConnection(hsRemove.UserMID, targetMId)
	res2 := a.Mongo.RemoveFromConnection(targetMId, hsRemove.UserMID)
	if res1 && res2 {
		var notify gbp.HandshakeDeleteNotify
		notify.Number = senderNUM
		notify.SenderMID = hsRemove.UserMID
		notify.TargetMID = targetMId
		notify.Mloc = "mloc"
		payload, err := proto.Marshal(&notify)
		if err != nil {
			var s utils.SuccessStruct
			s.Disc = "proto marshal error"
			s.Status = "unsuccessful"
			json.NewEncoder(w).Encode(s)
			return
		}
		var trans gbp.Transport
		trans.Id = hsRemove.UserMID
		trans.Msg = payload
		trans.Tp = 7
		transBytes, err := proto.Marshal(&trans)
		if err != nil {
			var s utils.SuccessStruct
			s.Disc = "proto marshal error"
			s.Status = "unsuccessful"
			json.NewEncoder(w).Encode(s)
			return
		}
		engineName := a.RMQ.GetEngineChannel()
		a.RMQ.Produce(engineName, transBytes)
		var s utils.SuccessStruct
		s.Disc = targetMId
		s.Status = "successful"
		json.NewEncoder(w).Encode(s)
		return
	} else {
		var s utils.SuccessStruct
		s.Disc = "Error while removing."
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}
}

func (a *API_V1) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send relevent data!"
		s.Status = "Unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}

	var changePic utils.ChangeProfilePayload
	err := json.NewDecoder(r.Body).Decode(&changePic)
	if err != nil {
		log.Println(err.Error())
	}

	userData, err := a.Mongo.ReadUserDataByMID(changePic.SenderMID)
	if err != nil {
		var s utils.SuccessStruct
		s.Disc = "no data found, bad sendermid"
		s.Status = "unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}

	if len(userData.Connections) > 0 {

		var notify gbp.ChangeProfilePayloads
		notify.PicData = changePic.PicData
		notify.SenderMID = changePic.SenderMID

		for key := range userData.Connections {
			notify.All = append(notify.All, key)
		}

		notifyBytes, err := proto.Marshal(&notify)
		if err != nil {
			var s utils.SuccessStruct
			s.Disc = "ERROR: " + err.Error()
			s.Status = "unsuccessful"
			json.NewEncoder(w).Encode(s)
			return
		}

		var trans gbp.Transport
		trans.Id = ""
		trans.Msg = notifyBytes
		trans.Tp = 8
		transBytes, err := proto.Marshal(&trans)
		if err != nil {
			var s utils.SuccessStruct
			s.Disc = "proto marshal error"
			s.Status = "unsuccessful"
			json.NewEncoder(w).Encode(s)
			return
		}
		engineName := a.RMQ.GetEngineChannel()
		a.RMQ.Produce(engineName, transBytes)
		res := a.Mongo.UpdateUserProfilePic(changePic.SenderMID, changePic.PicData)
		if res {
			var s utils.SuccessStruct
			s.Disc = ""
			s.Status = "successful"
			json.NewEncoder(w).Encode(s)
			fmt.Println("Pic changed with multiple notifies...")
			return
		}
	} else {
		res := a.Mongo.UpdateUserProfilePic(changePic.SenderMID, changePic.PicData)
		if res {
			var s utils.SuccessStruct
			s.Disc = ""
			s.Status = "successful"
			json.NewEncoder(w).Encode(s)
			fmt.Println("Pic changed")
			return
		}
	}

}

func (a *API_V1) UpdateNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Status = "unsuccessful"
		s.Disc = "No OTP is provided"
		json.NewEncoder(w).Encode(s)
	}
	var __otp utils.UpdateNumber
	_ = json.NewDecoder(r.Body).Decode(&__otp)
	res := a.Redis.VarifyOTP(__otp.OtpId, __otp.Otp)
	if res {

		if __otp.Notify == 1 {
			userData, err := a.Mongo.ReadUserDataByMID(__otp.MID)
			if err != nil {
				var s utils.SuccessStruct
				s.Disc = "no data found, bad sendermid"
				s.Status = "unsuccessful"
				json.NewEncoder(w).Encode(s)
				return
			}
			var notify gbp.NotifyChangeNumber
			notify.Number = __otp.Number
			notify.SenderMID = __otp.MID
			for mid := range userData.Connections {
				notify.All = append(notify.All, mid)
			}

			notifyBytes, err := proto.Marshal(&notify)
			if err != nil {
				var s utils.SuccessStruct
				s.Disc = "ERROR: " + err.Error()
				s.Status = "unsuccessful"
				json.NewEncoder(w).Encode(s)
				return
			}

			var trans gbp.Transport
			trans.Id = ""
			trans.Msg = notifyBytes
			trans.Tp = 9
			transBytes, err := proto.Marshal(&trans)
			if err != nil {
				var s utils.SuccessStruct
				s.Disc = "proto marshal error"
				s.Status = "unsuccessful"
				json.NewEncoder(w).Encode(s)
				return
			}
			engineName := a.RMQ.GetEngineChannel()
			err = a.RMQ.Produce(engineName, transBytes)

			if err != nil {
				res := a.Mongo.UpdateUserNumber(__otp.MID, __otp.Number)
				if res {
					var s utils.SuccessStruct
					s.Disc = ""
					s.Status = "successful"
					json.NewEncoder(w).Encode(s)
					fmt.Println("Numner changed with notify: 1")
					return
				}

			} else {
				var s utils.SuccessStruct
				s.Disc = "error while producing to rabbitmq"
				s.Status = "unsuccessful"
				json.NewEncoder(w).Encode(s)
				return
			}
		} else if __otp.Notify == 0 {
			res := a.Mongo.UpdateUserNumber(__otp.MID, __otp.Number)
			if res {
				var s utils.SuccessStruct
				s.Disc = ""
				s.Status = "successful"
				json.NewEncoder(w).Encode(s)
				fmt.Println("Number changed with notify: 0")
				return
			}
		}
	}
}

func (a *API_V1) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send relevent data!"
		s.Status = "Unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}

	var changeEmail utils.UpdateEmailPayload
	err := json.NewDecoder(r.Body).Decode(&changeEmail)
	if err != nil {
		log.Println(err.Error())
	}
	res := a.Mongo.UpdateUserEmail(changeEmail.MID, changeEmail.Email)
	if res {
		var s utils.SuccessStruct
		s.Disc = ""
		s.Status = "successful"
		json.NewEncoder(w).Encode(s)
		fmt.Println("Email is changed")
		return
	}
}

func (a *API_V1) LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		var s utils.SuccessStruct
		s.Disc = "Please send relevent data!"
		s.Status = "Unsuccessful"
		json.NewEncoder(w).Encode(s)
		return
	}

	var logoutReq utils.LogOutRequest
	err := json.NewDecoder(r.Body).Decode(&logoutReq)
	if err != nil {
		log.Println(err.Error())
	}

	res := a.Mongo.UpdateLogoutStatus(logoutReq.Mid, true)
	if res {
		var response gbp.Response
		response.Status = true
		response.Disc = "logout status changed to true"
		response.Data = ""
		responseBytes, err := proto.Marshal(&response)
		if err != nil {
			log.Println("[marshal error]", err.Error())
		}
		w.Write(responseBytes)
		return
	}

}
