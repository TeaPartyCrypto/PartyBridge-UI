package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	controller "github.com/TeaPartyCrypto/partybridge/be/pkg"
)

func main() {
	c := &controller.Controller{}
	// create a new sugard logger
	var err error
	c.Log, err = zap.NewProduction()
	if err != nil {
		log.Fatal(err)
		return
	}

	if c.SAASAddress == "" {
		c.SAASAddress = "http://143.42.111.52:8080"
		// c.SAASAddress = "http://192.168.50.90:8080"
	}

	http.HandleFunc("/", c.RootHandler)
	http.HandleFunc("/requestbridge", c.RequestBridge)
	log.Println("Listening on :8081")
	http.ListenAndServe(":8081", nil)
}
