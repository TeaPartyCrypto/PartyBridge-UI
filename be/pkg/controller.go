package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/google/uuid"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("BridgeRequest: %+v", bridgeRequest)

	// TODO: verify that sellOrder.Amount != nil
	// verify that all fields are not empty
	if bridgeRequest.Currency == "" || bridgeRequest.FromChain == "" || bridgeRequest.BridgeTo == "" || bridgeRequest.ShippingAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid sell order")
		return
	}

	bridgeRequest.ClientID = uuid.New().String()

	bridgeRequestJSON, err := json.Marshal(bridgeRequest)
	if err != nil {
		c.Log.Error("error marshalling bridgeRequest", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	io := bytes.NewBuffer(bridgeRequestJSON)
	req, err := http.NewRequest("POST", c.SAASAddress+"/requestBridge", io)
	if err != nil {
		c.Log.Error("error creating http request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error creating http request")
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Log.Error("error sending sell order to Party", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error sending sell order to Party")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Log.Error("error reading response body from Party", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("error reading response body from Party")
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
