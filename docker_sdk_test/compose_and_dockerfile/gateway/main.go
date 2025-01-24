package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
		message := os.Getenv("MESSAGE")
		if message == "" {
			message = "pong"
		}

    c.JSON(http.StatusOK, gin.H{
      "message": message,
    })
  })
  r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}