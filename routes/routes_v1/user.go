package routes_v1

import (
	v1 "github.com/ShikharY10/gbAPI/controllers/v1"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, api_v3 *v1.UserController) {

	// signup routes
	router.POST("/requestsignupotp", api_v3.RequestOtpForSignup)                                                  // tested: OK
	router.POST("/varifysignupotp", api_v3.Middleware.APIV3Varification("Authorization"), api_v3.VarifySignupOTP) // tested: OK
	router.GET("/isusernameawailable", api_v3.IsUsernameAwailable)
	router.POST("/signup", api_v3.Middleware.APIV3Varification("Authorization"), api_v3.SignUp) // tested: OK

	// login routes
	router.POST("/requestloginotp", api_v3.RequestOtpForLogin)
	router.POST("/login", api_v3.Middleware.APIV3Varification("Authorization"), api_v3.LogIn)
	router.GET("/getuseravatar/:id", api_v3.GetUserAvatar)

	// secured routes, only be accessed by authorized users.
	authorizedRoutes := router.Group("/a")
	authorizedRoutes.Use(api_v3.Middleware.APIV3Authorization())
	authorizedRoutes.POST("/logout", api_v3.LogOut)
	authorizedRoutes.GET("/getuserdetails", api_v3.GetUser)
	authorizedRoutes.GET("/getuseravatar/:id", api_v3.GetUserAvatar)

	authorizedRoutes.PUT("/updateavatar", api_v3.UpdateAvatar) // tested: Ok
	authorizedRoutes.PUT("/updatename", api_v3.UpdateUserName) // tested: OK

	authorizedRoutes.POST("/updateusername", api_v3.UpdateUsername)                                                                // tested: OK
	authorizedRoutes.POST("/varifyusernameotp", api_v3.Middleware.APIV3Varification("Auth-Token"), api_v3.VarifyUsernameUpdateOTP) // tested: OK

	authorizedRoutes.POST("/updateemail", api_v3.UpdateEmail)                                                               // tested: OK
	authorizedRoutes.POST("/varifyemailotp", api_v3.Middleware.APIV3EmailUpdateVarification(), api_v3.VarifyEmailUpdateOTP) // tested: OK

	// Not implemented yet and that is why it is commented.
	// securedUserAdmin := authorizedUser.Group("/admin")
	// securedUserAdmin.Use(middleware.RoleBasedAccess(&api_v3.Handler.UserHandler, "admin"))
	// securedUserAdmin.GET("/createadmin", api_v3.AdminCreateAdmin)
	// securedUserAdmin.GET("/users", api_v3.AdminGetAllUsers)
	// securedUserAdmin.GET("/user/:username", api_v3.AdminGetOneUser)
	// securedUserAdmin.PUT("/user/:username", api_v3.AdminUpdateOneUser)
	// securedUserAdmin.DELETE("/user/:username", api_v3.AdminDeleteOneUser)
}
