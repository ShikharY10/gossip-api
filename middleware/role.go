package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ShikharY10/gbAPI/handler"
	"github.com/gin-gonic/gin"
)

func RoleBasedAccess(userHandler *handler.UserHandler, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Value("username").(string)
		userRole := c.Value("role").(string)
		fmt.Println("[MiddleWare] -> role: " + userRole)
		_, err := userHandler.GetUserRole(username)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("user not found"))
		}
		var allowed bool = false
		for _, role := range roles {
			if role == userRole {
				allowed = true
				break
			}
		}

		if allowed {
			c.Next()
		} else {
			c.AbortWithError(http.StatusUnauthorized, errors.New("role not allowed"))
		}
	}
}
