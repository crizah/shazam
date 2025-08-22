package server

// import (
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// type WSHandler struct {
// 	Upgrader websocket.Upgrader
// }

// func MakeHandler() WSHandler {
// 	return WSHandler{Upgrader: websocket.Upgrader{}}

// }

// func (wsh WSHandler) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) error {

// 	conn, err := wsh.Upgrader.Upgrade(responseWriter, req, nil) // upgrades to websocket from http
// 	if err != nil {
// 		return err
// 	}

// 	defer conn.Close()
// 	return nil

// }
