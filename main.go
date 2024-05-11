package main

import (
	"net/http"
	"os"

	"ws_protocol/api"
	"ws_protocol/protocol"

	"fmt"

	"github.com/gorilla/mux"
)

//https://datatracker.ietf.org/doc/html/rfc6455#section-4.1
//The request MUST contain an |Upgrade| header field whose value MUST include the "websocket" keyword.
//The request MUST contain a |Connection| header field whose value MUST include the "Upgrade" token.

/*
The request MUST include a header field with the name
|Sec-WebSocket-Key|.  The value of this header field MUST be a
nonce consisting of a randomly selected 16-byte value that has
been base64-encoded (see Section 4 of [RFC4648]).  The nonce
MUST be selected randomly for each connection.
*/

/*
 The request MUST include a header field with the name
|Sec-WebSocket-Version|.  The value of this header field MUST be 13.
*/

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		panic("No arguments provided, must provide port")
	}
	port := args[0]

	upgrader := protocol.NewUpgrader()

	r := mux.NewRouter()
	r.HandleFunc("/", api.CommonHandler(upgrader)).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Server running on port: ", port)
	http.ListenAndServe(port, nil)

}
