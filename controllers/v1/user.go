package v1

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/handler"
	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Handler    *handler.Handler
	Middleware *middleware.Middleware
	Cloudinary *config.Cloudinary
	Logging    *logger.Logger
}

// func (v3 *UserController) requestEmailUpdateOTP(oldEmail string, newEmail string) (string, error) {
// 	id1, otp1 := v3.Handler.CacheHandler.RegisterOTP()
// 	id2, otp2 := v3.Handler.CacheHandler.RegisterOTP()

// 	var otpData map[string]string = map[string]string{
// 		"purpose":     "email",
// 		"oldemail":    oldEmail,
// 		"oldemailotp": otp1,
// 		"newemail":    newEmail,
// 		"newEmailotp": otp2,
// 	}
// 	_, err := json.Marshal(otpData)
// 	if err != nil {
// 		return "", err
// 	} else {
// 		fmt.Println("Old Email: ", oldEmail, "OTP: ", otp1)
// 		fmt.Println("New Email: ", newEmail, "OTP: ", otp2)

// 		claim := map[string]interface{}{
// 			"exp":      time.Now().Add(time.Minute * (60 * 5)).Unix(),
// 			"tokenid1": id1,
// 			"tokenid2": id2,
// 			"purpose":  "email",
// 		}
// 		token, err := middleware.GenerateJWT(claim, v3.Middleware.SecretKey)
// 		if err != nil {
// 			return "", err
// 		} else {
// 			v3.Handler.CacheHandler.Client.Set(id2+"_updateemail", newEmail, time.Minute*(60*5))
// 			v3.Handler.CacheHandler.Client.Set(id1+id2+"_purpose", "email", time.Minute*(60*5))
// 			return token, nil
// 		}
// 	}
// }

// func (v3 *UserController) requestOTP(email string, purpose string) (string, string, error) {
// 	// generating random OTP and UID
// 	id, otp := v3.Handler.CacheHandler.RegisterOTP()

// 	// sending OTP to OTP_SERVICE for sending it to the user
// 	var otpData map[string]string = map[string]string{
// 		"otp":     otp,
// 		"email":   email,
// 		"purpose": purpose,
// 	}
// 	_, err := json.Marshal(otpData)
// 	if err == nil {
// 		// produced to the queue that is listened by OTP_SERVICE
// 		// v3.Handler.QueueHandler.Produce("OTPd3hdzl8", b)
// 		fmt.Println("OTP: ", otp) // temp

// 		// generating authorization token
// 		claim := map[string]interface{}{
// 			"exp":     time.Now().Add(time.Minute * (60 * 5)).Unix(),
// 			"tokenid": id,
// 			"email":   email,
// 			"purpose": purpose,
// 		}
// 		token, err := middleware.GenerateJWT(claim, v3.Middleware.SecretKey)
// 		fmt.Println(len(token), " | Token Generated: ", token)
// 		if err != nil {
// 			return "", "", err
// 		}

// 		// storing id and number for future authorization.
// 		v3.Handler.CacheHandler.Client.Set(id+"_id", email, time.Minute*(60*5))
// 		v3.Handler.CacheHandler.Client.Set(id+"_purpose", purpose, time.Minute*(60*5))

// 		return token, id, nil
// 	} else {
// 		return "", "", err
// 	}
// }

// func (v3 UserController) RequestOtpForSignup(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)
// 	email := request["email"].(string)
// 	if email != "" {
// 		token, _, err := v3.requestOTP(email, "signup")
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
// 		} else {
// 			response := map[string]string{
// 				"token": token,
// 			}
// 			c.JSON(http.StatusCreated, response)
// 		}
// 	} else {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "email not found")
// 	}
// }

// func (v3 *UserController) VarifySignupOTP(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	id := c.Value("tokenid").(string)
// 	fmt.Println("id from VarifyOTP: ", id)
// 	otp := request["otp"].(string)

// 	if v3.Handler.CacheHandler.VarifyOTP(id, otp) {
// 		v3.Handler.CacheHandler.Client.Set(id+"_status", "varified", time.Minute*(60*5))
// 		response := map[string]string{
// 			"status": "successful",
// 		}
// 		c.JSON(http.StatusAccepted, response)
// 	} else {
// 		response := map[string]string{
// 			"status": "unsucessful",
// 		}
// 		c.JSON(http.StatusNotAcceptable, response)
// 	}
// }

// func (v3 *UserController) IsUsernameAwailable(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	username := c.Query("username")
// 	if username == "" {
// 		v3.Logging.LogError(errors.New("username not specified"))
// 		c.AbortWithStatus(400)
// 	} else {
// 		err := v3.Handler.UserHandler.IsUsernameAwailable(username)
// 		if err != nil {
// 			v3.Logging.LogError(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 		} else {
// 			c.JSON(200, "")
// 		}
// 	}
// }

// func (v3 *UserController) SignUp(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	var username string = request["username"].(string)
// 	var fullName string = request["fullname"].(string)

// 	id := c.Value("tokenid").(string)
// 	fmt.Println("id by key/value: ", id)

// 	objectId := primitive.NewObjectID()

// 	result := v3.Handler.CacheHandler.Client.Get(id + "_status")
// 	fmt.Println("result: ", result.Val())
// 	if result.Val() == "varified" {

// 		avatar := request["avatar"].(map[string]interface{})
// 		avatarData, err := v3.Cloudinary.UploadUserAvatar(id, avatar["imagedata"].(string), avatar["imageext"].(string))
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, "1. something went wrong, "+err.Error())
// 			return
// 		}

// 		imageId, err := v3.Handler.ImageHandler.SaveUserAvatar(objectId, avatar["imagedata"].(string), avatar["imageext"].(string))
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, "2. Something went wrong, "+err.Error())
// 		} else {
// 			avatarData.FileName = imageId

// 			messageID, err := v3.Handler.MsgHandler.AddUserMsgField()
// 			if err != nil {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, "3. Something went wrong, "+err.Error())
// 			}

// 			claims := map[string]interface{}{
// 				"uuid":     objectId.Hex(),
// 				"username": username,
// 				"role":     "user",
// 				"exp":      time.Now().AddDate(1, 0, 0).Unix(),
// 			}
// 			token, err := middleware.GenerateJWT(claims, v3.Middleware.SecretKey)

// 			if err != nil {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, "3. Something went wrong, "+err.Error())
// 			} else {
// 				var user models.User
// 				user.Avatar = *avatarData
// 				user.CreatedAt = time.Now().Format(time.RFC822)
// 				user.DeletedAt = ""
// 				user.Email = c.Value("email").(string)
// 				user.Partners = []primitive.ObjectID{}
// 				user.PartnerRequests = []models.PartnerRequest{}
// 				user.PartnerRequested = []models.PartnerRequest{}
// 				user.Posts = []primitive.ObjectID{}
// 				user.ID = objectId
// 				user.Logout = false
// 				user.MessageID = *messageID
// 				user.Name = fullName
// 				user.Role = "user"
// 				user.Token = token
// 				user.UpdatedAt = time.Now().Format(time.RFC822)
// 				user.Username = username

// 				err := v3.Handler.UserHandler.CreateNewUser(user)
// 				err1 := v3.Handler.UserHandler.InsetUserInFrequencyTable(user.ID, user.Username)
// 				if err != nil || err1 != nil {
// 					c.AbortWithStatusJSON(http.StatusInternalServerError, "4. something went wrong, "+err.Error())
// 				} else {
// 					c.JSON(http.StatusCreated, user)
// 					v3.Handler.CacheHandler.Client.Del(id)
// 					v3.Handler.CacheHandler.Client.Del(id + "_auth")
// 					v3.Handler.CacheHandler.Client.Del(id + "_status")
// 					v3.Handler.CacheHandler.Client.Del(id + "_purpose")
// 				}
// 			}
// 		}
// 	} else {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "5. Something went wrong")
// 	}
// }

// func (uc *UserController) RefreshAccessToken(c *gin.Context) {
// 	refreshToken, err := c.Cookie("refresh")
// 	if err != nil {
// 		c.AbortWithStatusJSON(500, err.Error())
// 		return
// 	}

// 	id := c.Param("id")
// 	if id == "" {
// 		c.AbortWithStatus(400)
// 		return
// 	}

// 	isTokenValid := uc.Cache.IsTokenValid(id, refreshToken, "refresh")
// 	if isTokenValid {
// 		refreshClaims, err := uc.Jwt.VarifyRefreshToken(refreshToken)
// 		if err != nil {
// 			c.AbortWithStatusJSON(401, "logged out")
// 			return
// 		}
// 		_id, err := primitive.ObjectIDFromHex(refreshClaims["id"].(string))
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, err.Error())
// 			return
// 		}

// 		user, err := uc.Database.GetUserData(bson.M{"_id": _id}, nil)
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, err.Error())
// 			return
// 		}

// 		newAccessTokenClaim := map[string]interface{}{
// 			"id":       user.Id.Hex(),
// 			"username": user.Username,
// 			"role":     user.Role,
// 			"exp":      time.Now().Add(time.Hour * 1).Unix(),
// 		}
// 		accessToken, err := uc.Jwt.GenerateJWT(newAccessTokenClaim, "access")
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, err.Error())
// 			return
// 		}

// 		uc.Cache.SetAccessTokenExpiry(user.Id.Hex(), accessToken, time.Hour*1)

// 		c.JSON(200, map[string]string{
// 			"accessToken": accessToken,
// 		})
// 	} else {
// 		c.AbortWithStatusJSON(401, "logged out")
// 		return
// 	}
// }

// func (v3 UserController) LogOut(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	username := c.Value("username").(string)

// 	result := v3.Handler.UserHandler.UpdateLogoutStatus(username, true)
// 	if result {
// 		c.JSON(http.StatusCreated, "sucessfully logged out")
// 	} else {
// 		c.AbortWithStatus(http.StatusPreconditionFailed)
// 	}
// }

func (v3 *UserController) GetUser(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(400, "id not found")
	} else {
		fmt.Println("id: ", id)
		uuid := c.Value("id").(string)
		isPartners := v3.Handler.UserHandler.IsPartners(id, uuid)
		var userData map[string]interface{}
		var err error

		if isPartners {
			userData, err = v3.Handler.UserHandler.GetUserDetails(id, "L2")
			if err != nil {
				// v3.Logging.LogError(err)
				c.AbortWithStatusJSON(500, err.Error())
				return
			}
			v3.Handler.UserHandler.IncrementUserInFrequencyTable(userData["username"].(string))

		} else {
			userData, err = v3.Handler.UserHandler.GetUserDetails(id, "L3")
			if err != nil {
				// v3.Logging.LogError(err)
				c.AbortWithStatusJSON(500, err.Error())
				return
			}
			v3.Handler.UserHandler.IncrementUserInFrequencyTable(userData["username"].(string))
		}
		c.JSON(200, userData)
	}
}

// func (v3 *UserController) RequestOtpForLogin(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// {
// 	// 	  "type": "email"
// 	//    "email": "yshikharfzd10@gmail.com"
// 	// }

// 	// {
// 	//    "type": "username",
// 	// 	  "username": "shikhary10"
// 	// }

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)
// 	userIdType := request["type"].(string)
// 	var email string
// 	var err error
// 	if userIdType == "username" {
// 		username := request["username"].(string)
// 		email, err = v3.Handler.UserHandler.GetUserEmail(username)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, "wrong username")
// 		}
// 	} else if userIdType == "email" {
// 		email = request["email"].(string)
// 		if email == "" {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, "wrong email")
// 		}
// 	}
// 	token, _, err := v3.requestOTP(email, "login")
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
// 	} else {
// 		response := map[string]string{
// 			"token": token,
// 		}
// 		c.JSON(http.StatusCreated, response)
// 	}
// }

// func (v3 *UserController) LogIn(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	bearer := c.GetHeader("Authorization")
// 	if bearer == "" {
// 		c.AbortWithStatusJSON(http.StatusForbidden, "token not found")
// 	} else {
// 		token := bearer[len("Bearer "):]
// 		if token == "" {
// 			c.AbortWithStatusJSON(http.StatusForbidden, "token not found")
// 		} else {
// 			claim, err := middleware.VarifyJWT(token, v3.Middleware.SecretKey)
// 			if err != nil {
// 				c.AbortWithStatusJSON(http.StatusForbidden, "token varification failed")
// 			} else {
// 				tokenID := claim["tokenid"].(string)
// 				result := v3.Handler.CacheHandler.Client.Get(tokenID + "_purpose")
// 				purpose := result.Val()

// 				if purpose != claim["purpose"].(string) || purpose != "login" {
// 					c.AbortWithStatusJSON(http.StatusPreconditionFailed, "token is compromised")
// 				} else {
// 					otp := request["otp"].(string)
// 					if !v3.Handler.CacheHandler.VarifyOTP(tokenID, otp) {
// 						c.AbortWithStatusJSON(http.StatusUnauthorized, "wrong OTP")
// 					} else {
// 						//
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func (v3 *UserController) GetUserAvatar(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "image/jpeg")
	c.Header("service", "Gossip API")

	id := c.Param("id")
	width := c.Query("width")
	height := c.Query("height")
	scale := c.Query("scale")

	var imgData []byte
	var err error
	if id == c.Value("id").(string) || v3.Handler.UserHandler.IsPartners(id, c.Value("id").(string)) {
		fmt.Println("Same || Partner")
		imgData, err = v3.Handler.UserHandler.GetUserAvatar(id, width, height, scale)
	} else {
		imgData, err = v3.Handler.UserHandler.GetUserAvatar(id, "100", "100", "0.70")
	}

	fmt.Println("imageData: ", imgData)

	if err != nil {
		fmt.Println("err: ", err.Error())
		c.AbortWithStatusJSON(500, err.Error())
	} else {
		f, err := os.Create("temp/" + id + ".jpg")
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
		}

		_, err = f.Write(imgData)
		if err != nil {
			c.AbortWithStatusJSON(500, err.Error())
		}
		c.File("temp/" + id + ".jpg")

		defer os.Remove("temp/" + id + ".jpg")
		defer f.Close()

		c.File("temp/" + id + ".jpg")
		// c.JSON(200, string(imgData))
	}
	// authenticated by token: size >= 50x50 pixels
	// normal: size <= 50x50 pixels
}

// user token based update
// func (v3 *UserController) UpdateUserName(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	uuid := c.Value("uuid").(string)

// 	// name
// 	fullName := request["fullname"].(string)
// 	if fullName == "" {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "name not found")
// 	} else {
// 		result := v3.Handler.UserHandler.UpdateUserName(uuid, fullName)
// 		if result {
// 			c.JSON(http.StatusCreated, map[string]string{"name": fullName})
// 		} else {
// 			c.AbortWithStatus(http.StatusPreconditionFailed)
// 		}
// 	}
// }

// user token based update
// func (v3 *UserController) UpdateAvatar(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	uuid := c.Value("uuid").(string)

// 	// name
// 	imageData := request["imagedata"].(string)
// 	imageExt := request["imageext"].(string)

// 	if imageData != "" && imageExt != "" {
// 		avatar, err := v3.Cloudinary.UploadUserAvatar(uuid+"temp", imageData, imageExt)
// 		if err != nil {
// 			c.AbortWithStatus(http.StatusPreconditionFailed)
// 		} else {
// 			objectId, err := primitive.ObjectIDFromHex(uuid)
// 			if err != nil {
// 				c.AbortWithStatus(http.StatusPreconditionFailed)
// 			} else {
// 				imageId, err := v3.Handler.ImageHandler.SaveUserAvatar(objectId, imageData, imageExt)
// 				avatar.FileName = imageId
// 				if err != nil {
// 					c.AbortWithStatus(http.StatusPreconditionFailed)
// 				} else {
// 					result := v3.Handler.UserHandler.UpdateUserAvatar(uuid, *avatar)
// 					if result {
// 						c.JSON(http.StatusCreated, avatar)
// 					} else {
// 						c.AbortWithStatus(http.StatusPreconditionFailed)
// 					}
// 				}
// 			}
// 		}
// 	} else {
// 		c.AbortWithStatus(http.StatusPreconditionFailed)
// 	}
// }

// special token based update
// func (v3 *UserController) UpdateUsername(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	if request == nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "body not found")
// 	} else {
// 		uuid := c.Value("uuid").(string)
// 		fmt.Println("uuid: ", uuid)
// 		username := request["username"].(string)
// 		email, err := v3.Handler.UserHandler.GetUserEmail(uuid)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong: "+err.Error())
// 		} else {
// 			token, id, err := v3.requestOTP(email, "username")
// 			if err != nil {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong: "+err.Error())
// 			} else {
// 				v3.Handler.CacheHandler.Client.Set(id+"_updateusername", username, time.Duration(time.Minute*(60*5)))
// 				response := map[string]string{
// 					"token": token,
// 				}
// 				c.JSON(http.StatusCreated, response)
// 			}
// 		}
// 	}
// }

// func (v3 *UserController) VarifyUsernameUpdateOTP(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	if request == nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "body not found")
// 	} else {
// 		tokenID := c.Value("tokenid").(string)
// 		uuid := c.Value("uuid").(string)
// 		otp := request["otp"].(string)

// 		if !v3.Handler.CacheHandler.VarifyOTP(tokenID, otp) {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, "wrong OTP")
// 		} else {
// 			newUsername := v3.Handler.CacheHandler.Client.Get(tokenID + "_updateusername").Val()
// 			updated := v3.Handler.UserHandler.UpdateUsername(uuid, newUsername)
// 			if updated {
// 				response := map[string]string{
// 					"username": newUsername,
// 				}
// 				c.JSON(http.StatusOK, response)
// 			} else {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, "error while updating username")
// 			}
// 		}
// 	}
// }

// special token based update
// func (v3 *UserController) UpdateEmail(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	if request == nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "body not found")
// 	} else {
// 		uuid := c.Value("uuid").(string)
// 		newEmail := request["newemail"].(string)
// 		oldEmail, err := v3.Handler.UserHandler.GetUserEmail(uuid)

// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
// 		} else {
// 			token, err := v3.requestEmailUpdateOTP(oldEmail, newEmail)
// 			if err != nil {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
// 			} else {
// 				response := map[string]string{
// 					"token": token,
// 				}
// 				c.JSON(http.StatusCreated, response)
// 			}
// 		}
// 	}
// }

// func (v3 *UserController) VarifyEmailUpdateOTP(c *gin.Context) {
// 	// setting response headers
// 	c.Header("Content-Type", "application/json")
// 	c.Header("service", "Gossip API")

// 	// collecting request body
// 	var request map[string]any
// 	c.BindJSON(&request)

// 	if request == nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "body not found")
// 	} else {
// 		// oldemailotp
// 		// newemailotp

// 		oldEmailOTP := request["oldemailotp"].(string)
// 		newEmailOTP := request["newemailotp"].(string)

// 		tokenID1 := c.Value("tokenid1").(string)
// 		tokenID2 := c.Value("tokenid2").(string)

// 		if v3.Handler.CacheHandler.VarifyOTP(tokenID1, oldEmailOTP) && v3.Handler.CacheHandler.VarifyOTP(tokenID2, newEmailOTP) {
// 			uuid := c.Value("uuid").(string)
// 			newEmail := v3.Handler.CacheHandler.Client.Get(tokenID2 + "_updateemail").Val()
// 			if v3.Handler.UserHandler.UpdateUserEmail(uuid, newEmail) {
// 				response := map[string]string{
// 					"email": newEmail,
// 				}
// 				c.JSON(http.StatusOK, response)
// 			} else {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, "error while updating email")
// 			}
// 		} else {
// 			c.AbortWithStatusJSON(http.StatusNotAcceptable, "Invalid OTP")
// 		}
// 	}
// }

func (v3 UserController) AdminGetAllUsers(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	users, err := v3.Handler.UserHandler.AdminGetAllUser()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	} else {
		response := map[string]interface{}{
			"users": users,
		}
		c.JSON(200, response)
	}
}

func (v3 UserController) AdminGetOneUser(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	username, found := c.Params.Get("username")
	if found {
		user, found := v3.Handler.UserHandler.AdminGetOneUser(username)
		if found {
			response := map[string]interface{}{
				"user": user,
			}
			c.JSON(200, response)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "something went wrong")
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, "user id not found")
	}
}

func (v3 *UserController) AdminCreateAdmin(c *gin.Context) {}

func (v3 UserController) AdminUpdateOneUser(c *gin.Context) {}

func (v3 UserController) AdminDeleteOneUser(c *gin.Context) {}
