package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// create a new sugard logger
	var err error
	c.Log, err = zap.NewProduction()
	if err != nil {
		log.Fatal(err)
		return
	}

	// fetch the IP address of the local machine
	ip, err := getPublicIP()
	if err != nil {
		c.Log.Error("error fetching public IP", zap.Error(err))
	} else {
		fmt.Println("Public IP:", ip)
	}
	c.IP = ip

	if c.SAASAddress == "" {
		c.SAASAddress = "http://143.42.111.52:8080"
		// c.SAASAddress = "http://192.168.50.90:8080"
	}

	http.HandleFunc("/", c.RootHandler)
	http.HandleFunc("/requestbridge", c.RequestBridge)
	log.Println("Listening on :8081")
	http.ListenAndServe(":8081", nil)
}

type IPResponse struct {
	IP string `json:"ip"`
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", fmt.Errorf("error fetching public IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var ipResponse IPResponse
	err = json.Unmarshal(body, &ipResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return ipResponse.IP, nil
}
