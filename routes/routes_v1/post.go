package routes_v1

import (
	v1 "github.com/ShikharY10/gbAPI/controllers/v1"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/gin-gonic/gin"
)

func PostRoute(router *gin.RouterGroup, api_v3 *v1.PostController) {
	post := router.Group("/api/v3/post/s")
	post.Use(api_v3.Middleware.APIV1_Authorization())

	post.POST("/post", api_v3.CreateNewPost)
	post.DELETE("/post/:id", api_v3.DeletePost)
	post.POST("/like", api_v3.LikeAPost)
	post.POST("/comment", api_v3.CommentToPost)
	post.GET("/post/:id", api_v3.GetOnePost)
	post.GET("/posts", api_v3.GetAllPost)

	postAdmin := post.Group("/admin")
	postAdmin.Use(middleware.RoleBasedAccess(api_v3.Handler.UserHandler, "admin"))
	postAdmin.GET("/posts", api_v3.AdminGetAllPosts)
	postAdmin.GET("/post/:id", api_v3.AdminGetOnePost)
	postAdmin.PUT("/post/:id", api_v3.AdminUpdateOnePost)
	postAdmin.DELETE("/post/:id", api_v3.AdminDeleteOnePost)
}
