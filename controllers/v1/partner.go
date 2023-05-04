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
	username := c.Query("username")
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

func (v3 *PartnerController) MakePartnerRequest(c *gin.Context) { // Post Request
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
			c.AbortWithStatusJSON(400, gin.H{
				"statusstring": err.Error(),
			})
		} else {
			// add request detail to target part of db
			err = v3.Handler.UserHandler.SetNewFollowRequest(request.TargetId, "partnerrequests", request)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{
					"statusstring": err.Error(),
				})
			} else {
				jsonBytes, err := json.Marshal(request)
				if err != nil {
					c.AbortWithStatusJSON(500, gin.H{
						"statusstring": err.Error(),
					})
				} else {
					payload := schema.Payload{
						Data: jsonBytes,
						Type: "011",
					}
					bpBytes, err := proto.Marshal(&payload)
					if err != nil {
						c.AbortWithStatusJSON(500, gin.H{
							"statusstring": err.Error(),
						})
					} else {
						// send request detail to gossip_engine
						engineName, err := v3.Handler.CacheHandler.GetEngineChannel()
						if err != nil {
							c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
								"statusstring": err.Error(),
							})
						} else {
							v3.Handler.QueueHandler.Produce(engineName, bpBytes)
							c.JSON(200, gin.H{
								"statusstring": "successfull",
							})
						}
					}
				}
			}
		}
	}
}

func (v3 *PartnerController) MakePartnerResponse(c *gin.Context) {
	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	// collecting request body
	var request models.PartnerResponse
	err := c.BindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"statusstring": err.Error(),
		})
	} else {
		if request.IsAccepted {
			// Step1: insert responserId to target's partner section in db
			err := v3.Handler.UserHandler.SetNewPartner(request.TargetId, request.ResponserId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"statusstring": err.Error(),
				})
				return
			}
			// Step2: insert targetId to responser's partner section in db
			err = v3.Handler.UserHandler.SetNewPartner(request.ResponserId, request.TargetId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"statusstring": err.Error(),
				})
				return
			}
		}
		jsonBytes, err := json.Marshal(request)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"statusstring": err.Error(),
			})
		} else {
			payload := schema.Payload{
				Data: jsonBytes,
				Type: "021",
			}
			bpBytes, err := proto.Marshal(&payload)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"statusstring": err.Error(),
				})
			} else {
				// send request detail to gossip_engine
				engineName, err := v3.Handler.CacheHandler.GetEngineChannel()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"statusstring": err.Error(),
					})
				} else {
					v3.Handler.QueueHandler.Produce(engineName, bpBytes)
					c.JSON(200, gin.H{
						"statusstring": "successfull",
					})
				}
			}
		}
	}
}

func (v3 *PartnerController) RemovePartner(c *gin.Context) {
	// unfollow/:partnerid  <-- choosing

	// setting response headers
	c.Header("Content-Type", "application/json")
	c.Header("service", "Gossip API")

	partnerId := c.Param("id")

	if partnerId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		userId := c.Value("id").(string)
		err1 := v3.Handler.UserHandler.RemovePartner(userId, partnerId)
		err2 := v3.Handler.UserHandler.RemovePartner(partnerId, userId)
		if err1 != nil && err2 != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"statusstring": err1.Error() + " | " + err2.Error(),
			})
		} else {
			removePartnerInfo := gin.H{
				"initiaterId": userId,
				"targetId":    partnerId,
			}
			jsonBytes, err := json.Marshal(removePartnerInfo)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"statusstring": err.Error(),
				})
				return
			}
			payload := schema.Payload{
				Data: jsonBytes,
				Type: "031",
			}
			bpBytes, err := proto.Marshal(&payload)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"statusstring": err.Error(),
				})
			} else {
				// send request detail to gossip_engine
				engineName, err := v3.Handler.CacheHandler.GetEngineChannel()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"statusstring": err.Error(),
					})
				} else {
					v3.Handler.QueueHandler.Produce(engineName, bpBytes)
					c.JSON(200, gin.H{
						"statusstring": "successfull",
					})
				}
			}
		}
	}
}

func (v3 *PartnerController) GetAllPartners(c *gin.Context) {

}

func (v3 *PartnerController) BlockPartner(c *gin.Context) {}

func (v3 *PartnerController) UnBlockPartner(c *gin.Context) {}

func (v3 *PartnerController) CheckForFriendExistance(c *gin.Context) {}

func (v3 *PartnerController) GetOneFriend(c *gin.Context) {}
