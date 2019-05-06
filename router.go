package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterActions interface {
	Setup()
}

type MyRouter struct {
	*gin.Engine
}

func (router *MyRouter) Setup() {
	gin.SetMode(gin.DebugMode)
	router.Engine = gin.New()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Authentication", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Authentication", "Authorization", "Content-Type"},
	}))

	router.POST("/login", LogIn)
	authorized := router.Group("/", validateToken())
	{
		authorized.GET("user/:id", GetUser, validateToken())
		authorized.GET("rooms", GetRooms, validateToken())
		authorized.GET("room/:id", GetRoomMessages, validateToken())
		authorized.POST("room/:id", PostMessage, validateToken())
		authorized.POST("command", ExecuteCommand, validateToken())
	}
}
