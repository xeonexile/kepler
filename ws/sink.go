package ws

import (
	"time"

	"github.com/lastexile/kepler"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// NewSink creates new WS outgoing sink. Acts as the ws server endpoint
func NewSink(connFactory ConnectionFactoryFunc, formatter kepler.MarshallerFunc, onConnect func(conn *websocket.Conn), onClose func()) (sink kepler.Sink, err error) {
	conn, err := connFactory()
	go readPump(conn, onClose)
	onConnect(conn)

	sink = kepler.NewSink(func(m kepler.Message) {

		value, err := formatter(m)
		if err != nil {
			log.Error("Unable to marshall message value")
			return
		}

		SendTextMessage(conn, value)
	})

	return
}

func readPump(conn *websocket.Conn, onClose func()) {
	closeHandler := onClose
	defer func() {
		if conn != nil {
			conn.Close()
		}

		if closeHandler != nil {
			closeHandler()
		}
		log.Info("defer close ")
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(string) error {
		log.Info("Ping")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	conn.SetPongHandler(func(string) error {
		log.Info("Pong")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if closeHandler != nil {
					closeHandler()
					closeHandler = nil
					log.Info("Connection Closed with onClose ")
				} else {
					log.Warn("No OnClose handler specified")
				}
			}
			break
		}
	}
}
