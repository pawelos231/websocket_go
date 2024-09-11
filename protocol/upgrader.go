package protocol

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"ws_protocol/helpers"

	"golang.org/x/crypto/blake2b"
)

type Upgrader struct {
	ReadBufferSize    uint16
	WriteBufferSize   uint16
	EnableCompression bool
}

// static for now
func NewUpgrader() *Upgrader {
	return &Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: false,
	}
}

// registerError is a helper function to handle errors during WebSocket upgrade.
func (upgrader *Upgrader) registerError(w http.ResponseWriter, status int, err error) (*WSConn, error) {
	// Check if the error is a HeaderError
	if herr, ok := err.(*HeaderError); ok {
		http.Error(w, herr.Error(), herr.Code)
	} else {
		http.Error(w, fmt.Sprintf("Status: %v. Error during WebSocket upgrade: %v", http.StatusText(status), err), status)
	}

	return nil, err
}

// optional to add: Sec-WebSocket-Protocol, Sec-WebSocket-Extensions
// https://datatracker.ietf.org/doc/html/rfc6455#section-4.2.2
/*
If the connection is happening on an HTTPS (HTTP-over-TLS) port,
perform a TLS handshake over the connection.  If this fails
*/
func (upgrader *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*WSConn, error) {
	headers := map[string]CustomHeader{
		"Connection":            {value: "Upgrade", toThrow: fmt.Errorf("connection header must be Upgrade"), code: http.StatusBadRequest},
		"Upgrade":               {value: "websocket", toThrow: fmt.Errorf("upgrade header must be websocket"), code: http.StatusBadRequest},
		"Sec-WebSocket-Version": {value: "13", toThrow: fmt.Errorf("Sec-WebSocket-Version must be 13"), code: http.StatusUpgradeRequired},
	}

	for headerName, headerValue := range headers {
		if err := upgrader.checkConnectionHeader(r, headerName, headerValue); err != nil {
			return upgrader.registerError(w, http.StatusInternalServerError, err)
		}
	}

	//hijack response writer to get the underlying connection
	hj, ok := w.(http.Hijacker)
	if !ok {
		err := fmt.Errorf("hijacking not supported")
		return upgrader.registerError(w, http.StatusInternalServerError, err)
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		return upgrader.registerError(w, http.StatusInternalServerError, err)
	}

	// Handle WebSocket handshake
	key, err := upgrader.createSecWebsocketKey(r)
	if err != nil {
		return upgrader.registerError(w, http.StatusInternalServerError, err)
	}

	//set response headers
	resHeaders := http.Header{}
	resHeaders.Set("Upgrade", "websocket")
	resHeaders.Set("Connection", "Upgrade")
	resHeaders.Set("Sec-WebSocket-Accept", key)

	res := http.Response{
		StatusCode: http.StatusSwitchingProtocols,
		Header:     resHeaders,
	}

	err = helpers.WriteToConnection(conn, res, upgrader.WriteBufferSize)
	if err != nil {
		return upgrader.registerError(w, http.StatusInternalServerError, err)

	}

	return &WSConn{Conn: conn}, nil

}

func (upgrader *Upgrader) createSecWebsocketKey(r *http.Request) (string, error) {
	//https://pkg.go.dev/crypto/sha1@go1.19.4?GOOS=windows
	//use blake2b to hash the key, beacuse 'SHA-1 is cryptographically broken and should not be used for secure applications.'
	//although rfc 6455 specifies SHA-1, will have to figure out a way how to use it
	secWebSocketKey := r.Header.Get("Sec-WebSocket-Key")
	hasher, err := blake2b.New256(nil)
	if err != nil {
		return "", fmt.Errorf("Error creating blake2b hasher: %v", err)
	}
	hasher.Write([]byte(secWebSocketKey + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	responseKey := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	return responseKey, nil
}

// checkConnectionHeader is a helper function to check if the request contains the required headers.
func (upgrader *Upgrader) checkConnectionHeader(r *http.Request, headerName string, headerContent CustomHeader) error {
	if r.Header.Get(headerName) != headerContent.value {
		return &HeaderError{
			HeaderName: headerName,
			Message:    headerContent.toThrow.Error(),
			Code:       headerContent.code,
		}
	}
	return nil
}
