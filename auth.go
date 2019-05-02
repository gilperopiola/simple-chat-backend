package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//LogIn takes {username, givenPassword}, checks if the user exists and returns it with a token
func LogIn(c *gin.Context) {
	var user User
	var givenPassword string
	c.BindJSON(&user)
	givenPassword = user.Password

	if len(user.Username) == 0 || len(givenPassword) == 0 {
		c.JSON(http.StatusBadRequest, "both fields required")
		return
	}

	database.Where("username = ?", user.Username).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, "wrong username")
		return
	}

	if user.Password != hash(user.Username, givenPassword) {
		c.JSON(http.StatusBadRequest, "wrong password")
		return
	}

	user.Token = generateToken(user)
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func validateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if len(tokenString) < 40 {
			c.JSON(http.StatusUnauthorized, "authentication error")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWT.SECRET), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, "authentication error")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			c.Set("ID", claims.Id)
			c.Set("Username", claims.Audience)
		} else {
			c.JSON(http.StatusUnauthorized, "authentication error")
			c.Abort()
		}
	}
}

func generateToken(user User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        fmt.Sprint(user.ID),
		Audience:  user.Username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(config.JWT.SESSION_DURATION)).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.JWT.SECRET))
	return tokenString
}

func hash(salt string, data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(salt + data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
