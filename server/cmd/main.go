package main

import (
	"net/http"

	"server/internal/chat"
	"server/internal/room"
	sfuSignal "server/internal/sfu-signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// http -> gin handler
func Handler(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
}

func main() {
	go chat.Run()
	go sfuSignal.Run()

	room.Rooms.Init()

	router := gin.New()

	router.Use(cors.New(cors.Config{AllowOrigins: []string{"http://localhost:3000"}}))

	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		chat.ServeWs(c.Writer, c.Request, roomId)
	})

	roomGroup := router.Group("")
	{
		roomGroup.GET("create", Handler(room.CreateRoom))
		roomGroup.GET("participants", Handler(room.GetRoom))
		roomGroup.GET("join", func(c *gin.Context) {
			room.JoinRoom(c.Writer, c.Request)
			roomID, _ := c.Request.URL.Query()["id"]
			chat.ServeWs(c.Writer, c.Request, roomID[0])

		})
	}

	router.Run("0.0.0.0:8080")
}
