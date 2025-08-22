package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// https for inital handshake and then upgrades to websocket

// client makes a GET request
// server responds with https 101

// Handler responds to http request
// Handler.ServeHTTP(ResponseWriter, *Request), should write reply headers and data to the ResponseWriter and then return

// type Upg struct {
// 	upgrader websocket.Upgrader
// }

type WSHandler struct {
	Upgrader websocket.Upgrader
}

func MakeHandler() WSHandler {
	return WSHandler{Upgrader: websocket.Upgrader{}}

}

func (wsh *WSHandler) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) error {

	conn, err := wsh.Upgrader.Upgrade(responseWriter, req, nil) // upgrades to websocket from http
	if err != nil {
		return err
	}

	defer conn.Close()
	return nil

}
