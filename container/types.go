package container

import (
	"github.com/gorilla/websocket"
)

type WsWriter struct {
	conn *websocket.Conn
}

func (w WsWriter) Write(p []byte) (n int, err error) {
	return len(p), w.conn.WriteMessage(websocket.TextMessage, p)
}
