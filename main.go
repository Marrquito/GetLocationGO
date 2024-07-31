package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	server := gin.Default()

	server.Use(cors.Default())

	server.POST("/", func(c *gin.Context) {
		var loc Location
		if err := c.ShouldBindJSON(&loc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Recebido: Latitude %f, Longitude %f\n", loc.Latitude, loc.Longitude)

		c.JSON(http.StatusOK, gin.H{"status": "sucesso"})
	})

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
