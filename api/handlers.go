package api

import (
	"net/http"
	"ws_protocol/protocol"
)

func CommonHandler(upgrader *protocol.Upgrader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WSConn, err := upgrader.Upgrade(w, r)
		if err != nil {
			return
		}
		if WSConn == nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}

		defer WSConn.Conn.Close()

		w.Write([]byte("Hello, World!"))
	}
}
