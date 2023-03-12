package api

import (
	"github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/controllers"
	"github.com/ShikharY10/gbAPI/handler"
	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/middleware"
	routes_v1 "github.com/ShikharY10/gbAPI/routes/routes_v1"
	"github.com/gin-gonic/gin"
)

func testRoutes(c *gin.Context) {
	c.String(200, "Home Route, Working!")
}

func StartVersionThreeAPIs(handle *handler.Handler, env *config.ENV, logging *logger.Logger) {
	var cloudinary config.Cloudinary = config.InitCloudinary()

	var authorization middleware.Middleware
	authorization.UserHandler = handle.UserHandler
	authorization.CacheHandler = handle.CacheHandler
	authorization.SecretKey = []byte(env.JWTSecret)

	var api_v3 controllers.API_V3
	api_v3.Handler = handle
	api_v3.Middleware = &authorization
	api_v3.Cloudinary = &cloudinary

	userController, partnerController, postController := controllers.GetController(
		handle,
		&authorization,
		&cloudinary,
		logging,
	)

	router := gin.New()
	v3Router := router.Group("/api/v3")

	router.GET("/", testRoutes)

	routes_v1.UserRoute(v3Router, userController)
	routes_v1.PartnerRoute(v3Router, partnerController)
	routes_v1.PostRoute(v3Router, postController)

	router.Run(":" + env.APIPort)
}
