package server

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"

	// "log"
	"net/http"
	"strings"
)

// https for inital handshake and then upgrades to websocket

// client makes a GET request
// server responds with https 101

// Handler responds to http request
// Handler.ServeHTTP(ResponseWriter, *Request), should write reply headers and data to the ResponseWriter and then return

type Conn interface {
	Close() error
}

type WebSocket struct {
	conn   Conn
	bufrw  *bufio.ReadWriter
	header http.Header
	status uint16
}

var Errors []error

func Handler(responseWriter http.ResponseWriter, req *http.Request) {
	// cant return error, needs to fit into a certain type

	webSocket, err := New(responseWriter, req)
	if err != nil {
		Errors = append(Errors, err)

	}

	err = webSocket.Handshake()
	if err != nil {
		Errors = append(Errors, err)
	}

	defer webSocket.Close()

}

func (ws *WebSocket) Close() error {

	return nil

}

func (ws *WebSocket) Handshake() error {

	key := ws.header.Get("Sec-WebSocket-Key")
	accept := calculateAccept(key)

	response := []string{
		"HTTP/1.1 101 Web Switching Protocols",
		"Server: go/echoserver",
		"Upgrade: WebSocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: " + accept,
		"", // required for extra CRLF
		"", // required for extra CRLF
	}

	joined := []byte(strings.Join(response, "\r\n"))

	_, err := ws.bufrw.Write(joined)

	if err != nil {
		return err

	}
	return ws.bufrw.Flush()

}

func calculateAccept(key string) string {
	// sha1 hash of concated key +guid
	// 16 byte, base 64

	h := sha1.New()
	guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11" // magic guid

	h.Write([]byte(key))
	h.Write([]byte(guid))

	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return encoded

}

func New(rw http.ResponseWriter, req *http.Request) (*WebSocket, error) {

	hj, ok := rw.(http.Hijacker)
	if !ok {
		return nil, errors.New("cant hijack")
	}

	conn, bufrw, err := hj.Hijack()

	if err != nil {
		return nil, err
	}

	return &WebSocket{conn: conn, bufrw: bufrw, header: req.Header, status: 101}, nil

}
