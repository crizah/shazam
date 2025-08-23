package server

// for getting raw audio data

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

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

// var PayloadChan = make(chan []byte, 10)

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

	for {
		frame, err := webSocket.Recv()
		if err != nil {
			Errors = append(Errors, err)
			break
		}

		//   *  %x0 denotes a continuation frame

		//   *  %x1 denotes a text frame

		//   *  %x2 denotes a binary frame

		//   *  %x3-7 are reserved for further non-control frames

		//   *  %x8 denotes a connection close

		//   *  %x9 denotes a ping

		//   *  %xA denotes a pong

		//   *  %xB-F are reserved for further control frames

		switch frame.Opcode {
		case 0x1: // text frame
			// PayloadChan <- frame.Payload

			fmt.Println("Received text:", string(frame.Payload))
		case 0x2: // binary frame
			fmt.Println("Received binary:", frame.Payload)
			// PayloadChan <- frame.Payload

		case 0x8: // close
			fmt.Println("Client closed connection")
			return
		case 0x9: // ping
			fmt.Println("Received ping")

		case 0xA: // pong
			fmt.Println("Received pong")
		default:
			fmt.Println("Unknown opcode", frame.Opcode)
		}
	}

}

func (ws *WebSocket) Close() error {

	return nil

}

type DataFrame struct {
	IsFragment bool
	Opcode     byte
	Reserved   byte
	IsMasked   bool
	Length     uint64
	Payload    []byte
}

func (ws *WebSocket) Recv() (DataFrame, error) {

	// 1= Fin(last frame), 3 is bs, last 4 is opcode
	// 1 = mask, rest = payload length

	// 1. Read first two-bytes a. find if the frame is a fragment b. find opcode c. find if the payload is masked d. find the payload length
	// 2. if length is less than 126, goto step#5
	// 3. if length equals to 126, read next two bytes in network byte order. This is the new payload length value
	// 4. if length equals to 127, read next eight bytes in network byte order. This is the new payload length value
	// 5. Read next 4 bytes as masking key
	// 6. Read next length bytes as masked payload data
	// 7. Decode the masked payload with masking key

	// opcodes

	// The length of the "Payload data", in bytes: if 0-125, that is the
	//   payload length.  If 126, the following 2 bytes interpreted as a
	//   16-bit unsigned integer are the payload length.  If 127, the
	//   following 8 bytes interpreted as a 64-bit unsigned integer (the
	//   most significant bit MUST be 0) are the payload length.  Multibyte
	//   length quantities are expressed in network byte order.  Note that
	//   in all cases, the minimal number of bytes MUST be used to encode
	//   the length, for example, the length of a 124-byte-long string
	//   can't be encoded as the sequence 126, 0, 124.  The payload length
	//   is the length of the "Extension data" + the length of the
	//   "Application data".  The length of the "Extension data" may be
	//   zero, in which case the payload length is the length of the
	//   "Application data".

	df := DataFrame{}
	head, err := ws.ReadFrame(2)
	if err != nil {
		return df, err

	}

	df.IsFragment = (head[0] & 0x80) == 0x00
	df.Opcode = head[0] & 0x0F
	df.Reserved = (head[0] & 0x70)

	df.IsMasked = (head[1] & 0x80) == 0x80

	var length uint64
	length = uint64(head[1] & 0x7F)

	if length == 126 {
		data, err := ws.ReadFrame(2)
		if err != nil {
			return df, err
		}
		length = uint64(binary.BigEndian.Uint16(data))
	} else if length == 127 {
		data, err := ws.ReadFrame(8)
		if err != nil {
			return df, err
		}
		length = uint64(binary.BigEndian.Uint64(data))
	}
	mask, err := ws.ReadFrame(4)
	if err != nil {
		return df, err
	}
	df.Length = length

	payload, err := ws.ReadFrame(int(length))
	if err != nil {
		return df, err
	}

	pl := decode(payload, length, mask)

	df.Payload = pl

	return df, nil

}

func decode(payload []byte, length uint64, mask []byte) []byte {
	// 	Octet i of the transformed data ("transformed-octet-i") is the XOR of
	//    octet i of the original data ("original-octet-i") with octet at index
	//    i modulo 4 of the masking key ("masking-key-octet-j"):

	//      j                   = i MOD 4
	//      transformed-octet-i = original-octet-i XOR masking-key-octet-j

	for i := uint64(0); i < length; i++ {
		payload[i] ^= mask[i%4]
	}

	return payload

}

func (ws *WebSocket) ReadFrame(n int) ([]byte, error) {
	// read the n bytes

	data := make([]byte, n)

	r, err := ws.bufrw.Read(data)
	if err != nil {
		return data, err
	}
	if r != n {
		return data, errors.New("dindt read the requested number of bytes")

	}

	return data, nil

}

func (ws *WebSocket) Handshake() error {

	key := ws.header.Get("Sec-WebSocket-Key")
	accept := calculateAccept(key)

	response := []string{
		"HTTP/1.1 101 Switching Protocols",
		"Server: go/server",
		"Upgrade: websocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: " + accept,
		"",
		"",
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
