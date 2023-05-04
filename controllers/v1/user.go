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
