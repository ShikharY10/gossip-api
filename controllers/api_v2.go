package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/ShikharY10/gbAPI/models"
	gbp "github.com/ShikharY10/gbAPI/protobuf"
	"github.com/ShikharY10/gbAPI/utils"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

type API_V2 struct {
	MsgModel   *models.MsgModel
	UserModel  *models.UserModel
	RedisModel *models.Redis
	RMQ        *config.RMQ
	AuthJwt    *middleware.JWT
}

func (a *API_V2) Apiv2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Version 2 APIs..."))
}

func (a *API_V2) SendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	params := mux.Vars(r)
	number := params["number"]

	id, otp := a.RedisModel.RegisterOTP()
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
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response

	if r.Body == nil {
		response.Status = false
		response.Disc = "empty body"
	} else {
		var __otp utils.VOTP
		_ = json.NewDecoder(r.Body).Decode(&__otp)
		res := a.RedisModel.VarifyOTP(__otp.Id, __otp.Otp)

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
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")

	var response gbp.Response
	response.Data = ""

	var newUserdata utils.NewUser

	err := json.NewDecoder(r.Body).Decode(&newUserdata)
	if err != nil {
		response.Status = false
		response.Disc = "bad json data"
	}
	id, _ := a.MsgModel.AddUserMsgField()
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

			fmt.Println("uD: ", uD)

			uid, err := a.UserModel.AddUser(uD)
			if err != nil {
				response.Status = false
				response.Disc = "unable to add user data to the database"
			} else {
				secretKey, found := os.LookupEnv("jwt-key")
				if !found {
					response.Status = false
					response.Disc = "unable to add user data to the database"
				} else {
					token, err := a.AuthJwt.GenerateJWT(newUserdata.Email, []byte(secretKey))
					if err != nil {
						fmt.Println("error in newuser: ", err.Error())
					}
					var responseMap map[string]string = map[string]string{
						"uid":     uid,
						"mid":     newUserdata.MsgId,
						"Eaeskey": utils.Encode(cipherText),
						"token":   token,
					}

					res, err := json.Marshal(responseMap)
					if err != nil {
						response.Status = false
						response.Disc = "error while preparing response"
					} else {
						response.Status = true
						response.Disc = "success"
						response.Data = utils.Encode(res) // utils.Encode(res)
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
	fmt.Println("ADDUSER RESPONSE SEND")
}

func (a *API_V2) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""

	if r.Body == nil {
		response.Disc = "empty body"

	} else {
		var loginData utils.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginData)

		if err != nil {
			response.Disc = "bad json data"
		} else {
			publicKey, err := utils.LoadKey([]byte(loginData.PublicKey))
			if err != nil {
				response.Disc = "bad pem formated public key"
			} else {
				seperated := strings.Split(loginData.Password, "+++")
				passwordHash := seperated[0]
				signature := []byte(seperated[1])
				digest := []byte(seperated[2])

				if utils.VarifySignature(publicKey, digest, signature) {
					myData, err := a.UserModel.GetUserDataByNumber(loginData.Number)
					if err != nil {
						response.Disc = "user not found"
					} else {
						if myData.Password == passwordHash {
							var loginResponsePayload gbp.LoginResponse
							var loginEnginePayload gbp.LoginEnginePayload

							encryptedMainKey, err := utils.RsaEncrypt(*publicKey, []byte(myData.MainKey))

							if err != nil {
								response.Disc = "unable to encrypt main key"
							} else {
								var myDataPayload gbp.UserData
								myDataPayload.Dob = myData.Age
								myDataPayload.Email = myData.Email
								myDataPayload.Gender = myData.Gender
								myDataPayload.MainKey = utils.Encode(encryptedMainKey)
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

									connData, err := a.UserModel.GetUserDataByMID(mid)
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

								secretkey, found := os.LookupEnv("jwt-key")
								if !found {

								} else {
									token, err := a.AuthJwt.GenerateJWT(myData.Email, []byte(secretkey))
									if err != nil {
										fmt.Println("error in newuser: ", err.Error())
									}

									loginResponsePayload.ConnData = connDatalist
									loginEnginePayload.AllConn = allConn
									loginResponsePayload.Token = token

									enginePayloadBytes, err := proto.Marshal(&loginEnginePayload)
									if err != nil {
										response.Disc = "internal server error"
									} else {
										var trans gbp.Transport
										trans.Id = ""
										trans.Msg = enginePayloadBytes
										trans.Tp = 10
										transBytes, err := proto.Marshal(&trans)

										if err != nil {
											response.Disc = "internal server error"
										} else {
											engineName := a.RMQ.GetEngineChannel()
											a.RMQ.Produce(engineName, transBytes)

											userPayloadBytes, err := proto.Marshal(&loginResponsePayload)

											if err != nil {
												response.Disc = "internal server error"
											} else {
												response.Status = true
												response.Disc = "ok"
												response.Data = utils.Encode(userPayloadBytes)
												a.UserModel.UpdateLogoutStatus(myData.MsgId, false)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	} else {
		w.Write(responseBytes)
		return
	}
}

func (a *API_V2) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout route hit")
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""
	if r.Body == nil {
		response.Disc = "empty body"
	} else {
		var logoutReq utils.LogOutRequest
		err := json.NewDecoder(r.Body).Decode(&logoutReq)
		if err != nil {
			response.Disc = "bad json data"
		} else if len(logoutReq.Mid) == 0 {
			response.Disc = "empty json data"
		} else {
			res := a.UserModel.UpdateLogoutStatus(logoutReq.Mid, true)
			if res {
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
	}
}

func (a *API_V2) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	fmt.Println("dashboard route hit")
	w.Write([]byte("Secured DashBoard of Gossip..."))
}

func (a *API_V2) ToggleBlock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Toggle Block")
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""
	if r.Body == nil {
		response.Disc = "empty body"
	} else {
		var blockRequest utils.ToggleBlocking
		err := json.NewDecoder(r.Body).Decode(&blockRequest)
		if err != nil {
			response.Disc = "bad json data" + err.Error()
		} else if len(blockRequest.SenderMID) == 0 || len(blockRequest.TargetNUM) == 0 {
			response.Disc = "bad json data: some fields are not found"
		} else {
			userData, err := a.UserModel.GetUserDataByMID(blockRequest.SenderMID)
			if err != nil {
				response.Disc = "sender not found"
			} else {
				tp := userData.Blocked[blockRequest.TargetNUM]

				if blockRequest.Type == -1 {
					if tp == 1 {
						res := a.UserModel.DeleteFromBlocking(blockRequest.SenderMID, blockRequest.TargetNUM)
						if res == 0 {
							response.Disc = "somthing went wrong"
						} else {
							response.Status = true
							response.Disc = "removed from blocking"
						}
					}
				} else if blockRequest.Type == 1 {
					if tp != 1 {
						res := a.UserModel.AddTOBlocking(blockRequest.SenderMID, blockRequest.TargetNUM)
						if res == 0 {
							response.Disc = "something went wrong"
						} else {
							response.Status = true
							response.Disc = "added to blocking"
						}
					}
				}
			}
		}
	}

	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	} else {
		w.Write(responseBytes)
		return
	}
}

func (a *API_V2) CheckAwailibity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""
	if r.Body == nil {
		response.Disc = "empty body"
	} else {
		var checkAccReq utils.CheckAccPayload
		err := json.NewDecoder(r.Body).Decode(&checkAccReq)
		if err != nil {
			response.Disc = "bad json data"
		} else if len(checkAccReq.SenderMID) < 1 {
			response.Disc = "empty json data"
		} else {
			found := a.UserModel.CheckAccountPresence(checkAccReq.TargetNUM)
			if !found {
				response.Disc = "user not found"
			} else {
				bFound := a.UserModel.CheckBlocking(checkAccReq.SenderMID, checkAccReq.TargetNUM)
				if bFound {
					response.Disc = "sender is blocked"
				} else {
					response.Status = true
					response.Disc = "success"
				}
			}
		}
	}

	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	} else {
		w.Write(responseBytes)
		return
	}
}

func (a *API_V2) RemoveFromHandshake(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""
	if r.Body == nil {
		response.Disc = "empty body"
	} else {
		var hsRemove utils.HandshakeDeletePaylaod
		err := json.NewDecoder(r.Body).Decode(&hsRemove)
		if err != nil {
			response.Disc = "bad json data"
		} else if len(hsRemove.TargetNUM) < 1 {
			response.Disc = "empty json data"
		} else {
			targetMId := a.UserModel.GetMsgIdByNumber(hsRemove.TargetNUM)
			senderNUM := a.UserModel.GetNumberByMsgId(hsRemove.UserMID)
			res1 := a.UserModel.RemoveFromConnection(hsRemove.UserMID, targetMId)
			res2 := a.UserModel.RemoveFromConnection(targetMId, hsRemove.UserMID)

			if res1 && res2 {
				var notify gbp.HandshakeDeleteNotify
				notify.Number = senderNUM
				notify.SenderMID = hsRemove.UserMID
				notify.TargetMID = targetMId
				notify.Mloc = "mloc"
				payload, err := proto.Marshal(&notify)
				if err != nil {
					response.Disc = "internal server error"
				} else {
					var trans gbp.Transport
					trans.Id = hsRemove.UserMID
					trans.Msg = payload
					trans.Tp = 7
					transBytes, err := proto.Marshal(&trans)
					if err != nil {
						response.Disc = "internal server error"
					} else {
						engineName := a.RMQ.GetEngineChannel()
						a.RMQ.Produce(engineName, transBytes)
						response.Status = true
						response.Disc = "success"
						response.Data = targetMId
					}
				}
			} else {
				response.Disc = "internal server error"
			}
		}
	}
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	} else {
		w.Write(responseBytes)
		return
	}
}

func (a *API_V2) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""
	if r.Body == nil {
		response.Disc = "empty body"
	} else {
		var changePic utils.ChangeProfilePayload
		err := json.NewDecoder(r.Body).Decode(&changePic)
		if err != nil {
			response.Disc = "bad json data"
		} else if len(changePic.PicData) < 1 {
			response.Disc = "empty json data"
		} else {
			userData, err := a.UserModel.GetUserDataByMID(changePic.SenderMID)
			if err != nil {
				response.Disc = "no user found"
			} else {
				if len(userData.Connections) > 0 {
					var notify gbp.ChangeProfilePayloads
					notify.PicData = changePic.PicData
					notify.SenderMID = changePic.SenderMID

					for key := range userData.Connections {
						notify.All = append(notify.All, key)
					}

					notifyBytes, err := proto.Marshal(&notify)
					if err != nil {
						response.Disc = "internal server error"
					} else {
						var trans gbp.Transport
						trans.Id = ""
						trans.Msg = notifyBytes
						trans.Tp = 8

						transBytes, err := proto.Marshal(&trans)
						if err != nil {
							response.Disc = "internal server error"
						} else {
							engineName := a.RMQ.GetEngineChannel()
							a.RMQ.Produce(engineName, transBytes)
							res := a.UserModel.UpdateUserProfilePic(changePic.SenderMID, changePic.PicData)
							if res {
								response.Status = true
								response.Disc = "success"
							} else {
								response.Disc = "internal server error"
							}
						}
					}
				} else {
					res := a.UserModel.UpdateUserProfilePic(changePic.SenderMID, changePic.PicData)
					if res {
						response.Status = true
						response.Disc = "success"
					} else {
						response.Disc = "internal server error"
					}
				}
			}
		}
	}
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	} else {
		w.Write(responseBytes)
		return
	}
}

func (a *API_V2) UpdateNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""
	if r.Body == nil {
		response.Disc = "empty body"
	} else {
		var __otp utils.UpdateNumber
		err := json.NewDecoder(r.Body).Decode(&__otp)

		if err != nil {
			response.Disc = "bad json data"
		} else if len(__otp.Otp) < 1 {
			response.Disc = "empty json data"
		} else {
			res := a.RedisModel.VarifyOTP(__otp.OtpId, __otp.Otp)

			if res {
				if __otp.Notify == 1 {
					userData, err := a.UserModel.GetUserDataByMID(__otp.MID)
					if err != nil {
						response.Disc = "no user found"
					} else {
						var notify gbp.NotifyChangeNumber
						notify.Number = __otp.Number
						notify.SenderMID = __otp.MID
						for mid := range userData.Connections {
							notify.All = append(notify.All, mid)
						}

						notifyBytes, err := proto.Marshal(&notify)

						if err != nil {
							response.Disc = "internal server error"
						} else {
							var trans gbp.Transport
							trans.Id = ""
							trans.Msg = notifyBytes
							trans.Tp = 9
							transBytes, err := proto.Marshal(&trans)
							if err != nil {
								response.Disc = "internal server error"
							} else {
								engineName := a.RMQ.GetEngineChannel()
								a.RMQ.Produce(engineName, transBytes)
								response.Status = true
								response.Disc = "success"
								response.Data = "1"
							}
						}
					}
				} else if __otp.Notify == 0 {
					res := a.UserModel.UpdateUserNumber(__otp.MID, __otp.Number)
					if res {
						response.Status = true
						response.Disc = "success"
						response.Data = "1"
					}
				}
			} else {
				response.Disc = "otp not changed"
				response.Data = "0"
			}
		}
	}
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	} else {
		w.Write(responseBytes)
		return
	}
}

func (a *API_V2) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/protobuf")
	w.Header().Set("service", "Gossip API")
	var response gbp.Response
	response.Status = false
	response.Data = ""
	if r.Body == nil {
		response.Disc = "empty body"
	} else {
		var changeEmail utils.UpdateEmailPayload
		err := json.NewDecoder(r.Body).Decode(&changeEmail)

		if err != nil {
			response.Disc = "bad json data"
		} else if len(changeEmail.Email) < 1 {
			response.Disc = "empty json data"
		} else {
			res := a.UserModel.UpdateUserEmail(changeEmail.MID, changeEmail.Email)
			if res {
				response.Status = true
				response.Disc = "success"
			} else {
				response.Disc = "no user found"
			}
		}
	}
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	} else {
		w.Write(responseBytes)
		return
	}
}
