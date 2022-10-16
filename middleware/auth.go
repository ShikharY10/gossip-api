package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ShikharY10/gbAPI/models"
	"github.com/golang-jwt/jwt"
)

type JWT struct {
	I         int
	MsgModel  *models.MsgModel
	userModel *models.UserModel
	Redis     *models.Redis
}

type VarifiedClaim struct {
	authorized bool
	email      string
	exp        float64
}

func (j *JWT) GenerateJWT(email string, secretekey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(secretekey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) VarifyJWT(token string, secretekey []byte) (*VarifiedClaim, error) {
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
		var newClaim VarifiedClaim
		newClaim.authorized = claims["authorized"].(bool)
		newClaim.email = claims["email"].(string)
		newClaim.exp = claims["exp"].(float64)
		return &newClaim, nil
	} else {
		return nil, errors.New("bad token")
	}
}

func (j *JWT) APIV2Auth(next http.Handler) http.Handler {
	fmt.Println("secure route")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headerToken = r.Header.Get("x-access-token")

		token := strings.TrimSpace(headerToken)

		if token == "" {
			fmt.Println("Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Missing auth token")
			return
		} else {
			var secreteKey string = ""

			secreteKey, err := j.Redis.GetSecretekey()
			if err != nil {
				// secreteKey, _ = j.Mongo.Secretekey()
				secreteKey, found := os.LookupEnv("SECRETE_KEY")
				if found {
					j.Redis.SetSecretekey(secreteKey)
				} else {
					return
				}

			}
			claim, err := j.VarifyJWT(token, []byte(secreteKey))
			if err != nil {
				fmt.Println("Bad token")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("bad token")
				return
			}
			if j.userModel.CheckUserExistence(claim.email) {
				fmt.Println("Authenticated | Bypassing...")
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Access Denied. Token Authorization Failed. This is a secure route and cannot be access directly")
				return
			}
		}
	})
}
