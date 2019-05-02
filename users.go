package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       uint   `json:"id" gorm:"auto_increment;unique;not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password,omitempty" gorm:"not null"`

	Token string `json:"token,omitempty" gorm:"-"`
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := database.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
