package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ShikharY10/gbAPI/handler"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var OneYear int64 = 31543204048
var FiveMin int64 = 300000

type Middleware struct {
	UserHandler  *handler.UserHandler
	CacheHandler *handler.CacheHandler
	SecretKey    []byte
}

func GenerateJWT(claim map[string]interface{}, secretekey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for k, v := range claim {
		claims[k] = v
	}

	tokenString, err := token.SignedString(secretekey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VarifyJWT(token string, secretekey []byte) (jwt.MapClaims, error) {
	newToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("something went wrong")
		}
		return secretekey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := newToken.Claims.(jwt.MapClaims); ok && newToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("bad token")
	}
}

func isTokenExpired(exp string, duration int64) bool {
	oldTime, err := time.Parse(time.RFC1123, exp)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(oldTime)

	if elapsed.Milliseconds() < duration {
		return false
	} else {
		return true
	}
}

// func (mw *Middleware) APIV2Auth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var headerToken = r.Header.Get("x-access-token")

// 		token := strings.TrimSpace(headerToken)

// 		if token == "" {
// 			w.WriteHeader(http.StatusForbidden)
// 			json.NewEncoder(w).Encode("Missing auth token")
// 			return
// 		} else {
// 			var secreteKey string = ""

// 			secreteKey, err := mw.CacheHandler.GetSecretekey()
// 			if err != nil {
// 				// secreteKey, _ = j.Mongo.Secretekey()
// 				secreteKey, found := os.LookupEnv("SECRETE_KEY")
// 				if found {
// 					mw.CacheHandler.SetSecretekey(secreteKey)
// 				} else {
// 					return
// 				}

// 			}
// 			claim, err := VarifyJWT(token, []byte(secreteKey))
// 			if err != nil {
// 				w.WriteHeader(http.StatusForbidden)
// 				json.NewEncoder(w).Encode("bad token")
// 				return
// 			}
// 			if mw.UserHandler.CheckUserExistence(claim["username"].(string)) {
// 				next.ServeHTTP(w, r)
// 			} else {
// 				w.WriteHeader(http.StatusForbidden)
// 				json.NewEncoder(w).Encode("Access Denied. Token Authorization Failed. This is a secure route and cannot be access directly")
// 				return
// 			}
// 		}
// 	})
// }

func (mw *Middleware) APIV3Varification(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("New Register Request=== " + c.Request.URL.Path + " ===")
		var token string
		if key == "Authorization" {
			bearer := c.GetHeader(key)
			if bearer != "" {
				token = bearer[len("Bearer "):]
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid JWT Token")
				return
			}
		} else if key == "Auth-Token" {
			token = c.GetHeader(key)
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid JWT Token")
			return
		} else {
			claim, err := VarifyJWT(token, mw.SecretKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid JWT Token")
				return
			}
			result1 := mw.CacheHandler.Client.Get(claim["tokenid"].(string) + "_id")
			email := result1.Val()

			result2 := mw.CacheHandler.Client.Get(claim["tokenid"].(string) + "_purpose")
			purpose := result2.Val()

			if email == claim["email"].(string) && purpose == claim["purpose"].(string) {
				// data := map[string]interface{}{
				// 	"tokenid": claim["tokenid"].(string),
				// 	"email":   claim["email"].(string),
				// }
				c.Set("tokenid", claim["tokenid"].(string))
				c.Set("email", claim["email"].(string))
				// c.Keys = data
				fmt.Println("===Register Request Varified ===")
				c.Next()
				return
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "token data compromised")
				return
			}
		}

	}
}

func (mw *Middleware) APIV3EmailUpdateVarification() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("New Email Update Request=== " + c.Request.URL.Path + " ===")
		token := c.GetHeader("Auth-Token")
		if token != "" {
			claim, err := VarifyJWT(token, mw.SecretKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid JWT Token")
				return
			}
			tokenID1 := claim["tokenid1"].(string)
			tokenID2 := claim["tokenid2"].(string)

			purpose := mw.CacheHandler.Client.Get(tokenID1 + tokenID2 + "_purpose").Val()

			if purpose == claim["purpose"].(string) {
				c.Set("tokenid1", tokenID1)
				c.Set("tokenid2", tokenID2)
				fmt.Println("===Email Update Request Varified ===")
				c.Next()
				return
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "token data compromised")
				return
			}
		}
	}
}

func (mw *Middleware) APIV3Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("New Request=== " + c.Request.URL.Path + " ===")
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, "token not found")
			return
		} else {
			token := bearer[len("Bearer "):]
			if token == "" {
				c.AbortWithStatusJSON(http.StatusForbidden, "token not found")
				return
			} else {
				claim, err := VarifyJWT(token, mw.SecretKey)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid JWT Token, "+err.Error())
					return
				} else {
					data := map[string]interface{}{
						"uuid":     claim["uuid"].(string),
						"username": claim["username"].(string),
						"role":     claim["role"].(string),
					}
					c.Keys = data
					fmt.Println("=== Request Varified ===")
					c.Next()
				}
			}
		}
	}
}

// Adding one year duration
// time.Now().AddDate(1, 0, 0).Unix(),
