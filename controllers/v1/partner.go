package v1

import (
	"encoding/json"
	"net/http"

	"github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/handler"
	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/ShikharY10/gbAPI/models"
	"github.com/ShikharY10/gbAPI/schema"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type PartnerController struct {
	Handler    *handler.Handler
	Middleware *middleware.Middleware
	Cloudinary *config.Cloudinary
	Logging    *logger.Logger
}

func (v3 *PartnerController) SearchUsername(c *gin.Context) {
	username := c.Query("q")
	if username == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		result, err := v3.Handler.UserHandler.SearchUsername(username)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func (v3 *PartnerController) PartnerRequest(c *gin.Context) { // Post Request
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request models.PartnerRequest
	err := c.BindJSON(&request)

	if err != nil {
		c.AbortWithStatus(400)
	} else {
		// add request detail to sender part of db
		err = v3.Handler.UserHandler.SetNewFollowRequest(request.RequesterId, "partnerrequested", request)
		if err != nil {
			c.AbortWithStatusJSON(400, "SENDER ERROR: "+err.Error())
		} else {
			// add request detail to target part of db
			err = v3.Handler.UserHandler.SetNewFollowRequest(request.TargetId, "partnerrequests", request)
			if err != nil {
				c.AbortWithStatusJSON(400, "Target ERROR: "+err.Error())
			} else {
				jsonBytes, err := json.Marshal(request)
				if err != nil {
					c.AbortWithStatus(500)
				} else {
					payload := schema.Payload{
						Data: jsonBytes,
						Type: "021",
					}
					bpBytes, err := proto.Marshal(&payload)
					if err != nil {
						c.AbortWithStatus(500)
					} else {
						// send request detail to gossip_engine
						engineName, err := v3.Handler.CacheHandler.GetEngineChannel()
						if err != nil {
							c.AbortWithStatusJSON(http.StatusInternalServerError, "intarnal server error")
						} else {
							v3.Handler.QueueHandler.Produce(engineName, bpBytes)
							c.JSON(200, "")
						}
					}
				}
			}
		}
	}
}

func (v3 *PartnerController) PartnerResponse(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request models.PartnerResponse
	err := c.BindJSON(&request)

	if err != nil {
		c.AbortWithStatus(400)
	} else {
		if request.IsAccepted {
			// Step1: insert responserId to target's partner section in db
			err := v3.Handler.UserHandler.SetNewPartner(request.TargetId, request.ResponserId)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Step2: insert targetId to responser's partner section in db
			err = v3.Handler.UserHandler.SetNewPartner(request.ResponserId, request.TargetId)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		jsonBytes, err := json.Marshal(request)
		if err != nil {
			c.AbortWithStatus(500)
		} else {
			payload := schema.Payload{
				Data: jsonBytes,
				Type: "021",
			}
			bpBytes, err := proto.Marshal(&payload)
			if err != nil {
				c.AbortWithStatus(500)
			} else {
				// send request detail to gossip_engine
				engineName, err := v3.Handler.CacheHandler.GetEngineChannel()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, "intarnal server error")
				} else {
					v3.Handler.QueueHandler.Produce(engineName, bpBytes)
					c.JSON(200, "")
				}
			}
		}
	}
}

func (v3 *PartnerController) UnFollow(c *gin.Context) {
	// unfollow/:partnerid  <-- choosing

	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	partnerId := c.Param("id")

	if partnerId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		userId := c.Value("uuid").(string)
		err := v3.Handler.UserHandler.RemovePartner(userId, partnerId)
		if err != nil {
			c.AbortWithStatus(500)
		} else {
			c.JSON(200, "succesfully removed")
		}
	}
}

func (v3 *PartnerController) GetAllFollowers(c *gin.Context) {}

func (v3 *PartnerController) GetAllFollowing(c *gin.Context) {}

func (v3 *PartnerController) BlockFriend(c *gin.Context) {}

func (v3 *PartnerController) UnBlockFriend(c *gin.Context) {}

func (v3 *PartnerController) CheckForFriendExistance(c *gin.Context) {}

func (v3 *PartnerController) GetOneFriend(c *gin.Context) {}
