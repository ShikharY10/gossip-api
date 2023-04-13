package routes_v1

import (
	v1 "github.com/ShikharY10/gbAPI/controllers/v1"
	"github.com/gin-gonic/gin"
)

func PartnerRoute(router *gin.RouterGroup, api_v3 *v1.PartnerController) {

	authorizedRoutes := router.Group("/a")
	authorizedRoutes.Use(api_v3.Middleware.APIV1_Authorization())
	authorizedRoutes.GET("/searchusername", api_v3.SearchUsername)
	authorizedRoutes.POST("/partnerrequest", api_v3.PartnerRequest)
	authorizedRoutes.POST("partnerresponse", api_v3.PartnerResponse)
	authorizedRoutes.POST("/unfollow", api_v3.UnFollow)
	authorizedRoutes.GET("/followers/:user", api_v3.GetAllFollowers)
	authorizedRoutes.GET("/following/:user", api_v3.GetAllFollowing)
	authorizedRoutes.POST("/block", api_v3.BlockFriend)
	authorizedRoutes.POST("/unblock", api_v3.UnBlockFriend)
	authorizedRoutes.POST("/exist/:user", api_v3.CheckForFriendExistance)
}
