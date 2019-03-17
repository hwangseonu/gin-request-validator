package security

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful-example/basic/server/models"
	"net/http"
	"strings"
	"time"
)

const (
	ACCESS = "access"
	REFRESH = "refresh"
)

var secret = "awesome-jwt"

type CustomClaims struct {
	jwt.StandardClaims
	Identity string `json:"identity"`
}

func (c CustomClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if !models.ExistsUserByUsername(c.Identity){
		return errors.New("cannot find user by username")
	}
	return nil
}

func GenerateToken(t, username string) string {
	var expire int64

	if t == ACCESS {
		expire = time.Now().Add(time.Hour).Unix()
	} else if t == REFRESH {
		expire = time.Now().AddDate(0, 1, 0).Unix()
	}

	claims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: expire,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "",
		Subject: t,
		NotBefore: time.Now().Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, CustomClaims{claims, username}).SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return token
}

func AuthRequired(sub string, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var token *jwt.Token
		var err error
		if token, err = jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		claims := token.Claims.(*CustomClaims)
		if claims.Subject != sub {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "jwt subject must " + sub})
			return
		}
		u := models.FindUserByUsername(claims.Identity)
		if u == nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "cannot find user by jwt identity"})
			return
		} else if !containsAll(u.Roles, roles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "cannot access this resource"})
			return
		}
		c.Set("user", u)

		if sub == REFRESH {
			c.Set("exp", claims.ExpiresAt)
		}
		c.Next()
	}
}

func containsAll(s []string, e[] string) bool {
	cnt := len(e)
	for _, a := range s {
		for _, b := range e {
			if a == b {
				cnt--
			}
		}
	}
	return cnt == 0
}