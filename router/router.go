package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"price-tracker/handler"
)

func SetupRouter(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/all", func(c *gin.Context) {
		data, err := handler.GetAll()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"data": data})
	})

	router.GET("/price", func(c *gin.Context) {
		// get raw byte data from request body
		data, err := c.GetRawData()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		url := string(data)

		fmt.Println(url)

		price, _ := handler.GetPrice(url)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"price": price})
	})

	router.POST("/track", func(c *gin.Context) {
		// get raw byte data from request body
		data, err := c.GetRawData()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		url := string(data)

		product, err := handler.TrackPrice(url)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"tracking": product})
	})

	return router
}
