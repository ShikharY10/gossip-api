package v1

import (
	"github.com/ShikharY10/gbAPI/config"
	"github.com/ShikharY10/gbAPI/handler"
	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/middleware"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	Handler    *handler.Handler
	Middleware *middleware.Middleware
	Cloudinary *config.Cloudinary
	Logging    *logger.Logger
}

// // Post Related APIs========================================================
func (v3 *PostController) CreateNewPost(c *gin.Context) {}

// func (v3 **PostController) CreateNewPost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) DeletePost(c *gin.Context) {}

// func (v3 **PostController) DeletePost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) LikeAPost(c *gin.Context) {}

// func (v3 **PostController) LikeAPost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) CommentToPost(c *gin.Context) {}

// func (v3 **PostController) CommentToPost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) GetOnePost(c *gin.Context) {}

// func (v3 **PostController) GetOnePost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) GetAllPost(c *gin.Context) {}

// func (v3 **PostController) GetAllPost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) AdminGetAllPosts(c *gin.Context) {}

// func (v3 **PostController) AdminGetAllPosts(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) AdminGetOnePost(c *gin.Context) {}

// func (v3 **PostController) AdminGetOnePost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) AdminUpdateOnePost(c *gin.Context) {}

// func (v3 **PostController) AdminUpdateOnePost(w http.ResponseWriter, r *http.Request) {}

func (v3 *PostController) AdminDeleteOnePost(c *gin.Context) {}

// func (v3 **PostController) AdminDeleteOnePost(w http.ResponseWriter, r *http.Request) {}

// // upload APIs============================================================
// func (v3 **PostController) UploadChatImage(c *gin.Context) {}

// func (v3 **PostController) UploadChatImage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("[V3] : Chat Image Upload")
// 	var response Response
// 	response.Data = ""
// 	// var body map[string]string
// 	var body imageUpload
// 	err := json.NewDecoder(r.Body).Decode(&body)
// 	if err != nil {
// 		response.Status = 1
// 		response.Disc = err.Error()
// 		fmt.Println("Error")
// 	} else {
// 		response.Data = "working"
// 		secureUrl, publicId, err := v3.Cloudinary.UploadChatImage(body.Image, body.Ext)
// 		if err != nil {
// 			fmt.Println("err: ", err)
// 			response.Status = 1
// 			response.Disc = err.Error()
// 			fmt.Println("Upload Error")
// 		} else {
// 			var resMap map[string]string = make(map[string]string)
// 			resMap["secureUrl"] = secureUrl
// 			resMap["publicId"] = publicId
// 			resBytes, err := json.Marshal(&resMap)
// 			if err != nil {
// 				response.Status = 1
// 				response.Disc = err.Error()
// 				fmt.Println("Response Parse Error")
// 			} else {
// 				response.Data = string(resBytes)
// 				fmt.Println("Image Upload Response Written!")
// 			}
// 		}
// 	}
// 	json.NewEncoder(w).Encode(response)
// }

func (v3 *PostController) UploadChatVideo(c *gin.Context) {}

// func (v3 **PostController) UploadChatVideo(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Version 3 UploadChatVideo APIs..."))
// }

func (v3 *PostController) UploadProfilePicture(c *gin.Context) {}

// func (v3 **PostController) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Version 3 UploadProfilePicture APIs..."))
// }

func (v3 *PostController) UpdateProfilePicture(c *gin.Context) {}
