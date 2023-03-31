package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	controller "github.com/TeaPartyCrypto/partybridge/be/pkg"
)

type Account struct {
	PrivateKey  string `json:"PrivateKey"`
	PublicKey   string `json:"PublicKey"`
	ProgramHash string `json:"ProgramHash"`
}

func main() {
	c := &controller.Controller{}
	// look for a local NKN wallet
	// if one is not found, create one
	// and save it to the local file system
	// for future use.

	if c.SAASAddress == "" {
		c.SAASAddress = "http://38.109.255.242:8080"
	}

	// create a new sugard logger
	var err error
	c.Log, err = zap.NewProduction()
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", c.RootHandler)
	http.HandleFunc("/requestbridge", c.RequestBridge)
	log.Println("Listening on :8081")
	http.ListenAndServe(":8081", nil)
}
