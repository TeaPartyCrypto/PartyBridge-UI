package controller

import (
	"net/http"

	"go.uber.org/zap"
)

type Controller struct {
	rootHandler http.Handler
	SAASAddress string
	IP          string

	Log *zap.Logger
}

type BridgeRequest struct {
	Currency        string `json:"currency"`
	FromChain       string `json:"fromChain"`
	Amount          int    `json:"amount"`
	BridgeTo        string `json:"bridgeTo"`
	ShippingAddress string `json:"shippingAddress"`
	ClientID        string `json:"clientId"`
}
