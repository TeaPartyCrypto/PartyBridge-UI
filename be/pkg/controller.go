package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var once sync.Once

// Sell is a http route handler that accepts a sell order
// sell orders are stored in an on prem MongoDB database
func (c *Controller) RequestBridge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Log.Debug("received sell order")

	bridgeRequest := &BridgeRequest{}
	err := json.NewDecoder(r.Body).Decode(bridgeRequest)
	if err != nil {
		c.Log.Error("error decoding BridgeRequest order", zap.Error(err))
		sendErrorResponse(w, http.StatusBadRequest, "error decoding BridgeRequest order")
		return
	}

	fmt.Printf("BridgeRequest: %+v", bridgeRequest)

	// TODO: verify that sellOrder.Amount != nil
	// verify that all fields are not empty
	if bridgeRequest.Currency == "" || bridgeRequest.FromChain == "" || bridgeRequest.BridgeTo == "" || bridgeRequest.ShippingAddress == "" {
		sendErrorResponse(w, http.StatusBadRequest, "error decoding BridgeRequest order")
		return
	}

	bridgeRequestJSON, err := json.Marshal(bridgeRequest)
	if err != nil {
		c.Log.Error("error marshalling bridgeRequest", zap.Error(err))
		sendErrorResponse(w, http.StatusBadRequest, "error marshalling bridgeRequest")
		return
	}

	io := bytes.NewBuffer(bridgeRequestJSON)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is missing", http.StatusBadRequest)
		return
	}
	fmt.Println("id: " + id)

	saasURL, err := url.Parse(c.SAASAddress + "/requestBridge?id=" + id)
	if err != nil {
		c.Log.Error("error parsing URL", zap.Error(err))
		sendErrorResponse(w, http.StatusBadRequest, "error parsing URL")
		return
	}

	req, err := http.NewRequest("POST", saasURL.String(), io)
	if err != nil {
		c.Log.Error("error creating http request", zap.Error(err))
		sendErrorResponse(w, http.StatusBadRequest, "error creating http request")
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Log.Error("error sending sell order to Party", zap.Error(err))
		sendErrorResponse(w, http.StatusBadRequest, "error sending sell order to Party")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Log.Error("error reading response body from Party", zap.Error(err))
		sendErrorResponse(w, http.StatusBadRequest, "error reading response body from Party")
		return
	}

	// return accepted to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(body)
}

func (c *Controller) RootHandler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		kdp := path.Join(os.Getenv("KO_DATA_PATH"))
		if kdp == "" {
			kdp = "kodata"
		}

		if !strings.HasSuffix(kdp, "/") {
			kdp = kdp + "/"
		}
		c.rootHandler = http.FileServer(http.Dir(kdp))
	})

	c.rootHandler.ServeHTTP(w, r)
}

// Helper function to send error responses as JSON
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
