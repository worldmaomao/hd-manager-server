package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"log"
	"net/http"
	"strings"
	"time"
	"worldmaomao/harddisk/internal/config"
	"worldmaomao/harddisk/internal/constant"
	"worldmaomao/harddisk/internal/dao/model"
)

type Claims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.StandardClaims
}

func GenerateJWTToken(user model.User, jwtKey string, platformAudience string) string {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(48 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: user.Username,
		Roles:    user.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Audience:  platformAudience,
			Issuer:    "worldmaomao",
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return tokenString
}

func RequireAuthenticated(c *gin.Context, container di.Container) {
	config := container.Get(constant.Configuration).(*config.Configuration)
	token := c.GetHeader("Authorization")
	if token == "" {
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	jwtToken := strings.Split(token, "Bearer ")
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(jwtToken[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJwtKey()), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
		}
		c.Abort()
		c.Writer.WriteHeader(http.StatusBadRequest)
	}
	if !tkn.Valid {
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
	}
	return

}
