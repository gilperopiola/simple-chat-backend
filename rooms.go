package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Room struct {
	ID uint `json:"id" gorm:"auto_increment;unique;not null"`
}

type Message struct {
	ID     uint   `json:"id" gorm:"auto_increment;unique;not null"`
	Room   *Room  `json:"room,omitempty" database:"-"`
	RoomID int    `json:"room_id,omitempty"`
	User   *User  `json:"user,omitempty" database:"-"`
	UserID int    `json:"user_id,omitempty"`
	Data   string `json:"data" gorm:"not_null"`
}

func GetRoomMessages(c *gin.Context) {
	roomID := c.Param("id")
	since := c.Query("since")

	var messages []Message
	database.Where("id > " + since + " AND room_id = " + roomID).Find(&messages)
	c.JSON(http.StatusOK, messages)
}

func GetRooms(c *gin.Context) {
	var rooms []Room
	database.Find(&rooms)
	c.JSON(http.StatusOK, rooms)
}

func PostMessage(c *gin.Context) {
	roomID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("ID")

	var message Message
	c.BindJSON(&message)

	//checks if room exists
	var room Room
	if err := database.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, database.BeautifyError(err))
		return
	}
	message.RoomID = int(room.ID)

	//checks if user exists
	var user User
	if err := database.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, database.BeautifyError(err))
		return
	}
	message.UserID = int(user.ID)

	if err := database.Create(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, database.BeautifyError(err))
		return
	}

	c.JSON(http.StatusOK, message)
}
