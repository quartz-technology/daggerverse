package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		fmt.Println("Receive ping request")

		url := os.Getenv("BACKEND_URL")

		fmt.Println("Sending request to backend at", url)
		res, err := http.Get(fmt.Sprintf("%s/ping", url))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		defer res.Body.Close()
		
		// Read the response body
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Convert the body to a string
		bodyString := string(bodyBytes)

		c.JSON(http.StatusOK, gin.H{
			"message": os.Getenv("MESSAGE"),
			"message_from_backend": bodyString,
		})
	})

	fmt.Printf("os.Getenv(MESSAGE): %s\n", os.Getenv("MESSAGE"))
	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
