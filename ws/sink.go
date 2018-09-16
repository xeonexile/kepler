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

	go func() {
		ticker := time.NewTicker(pingPeriod)
		for range ticker.C {
			if write(conn, websocket.PingMessage, []byte{}); err != nil {
				log.Errorln("PING error: %v\n", err)
			}
		}
	}()
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
	defer func() {
		if conn != nil {
			conn.Close()
		}

		if onClose != nil {
			onClose()
		}
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(string) error {
		log.Debug("Ping")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	conn.SetPongHandler(func(string) error {
		log.Debug("Pong")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warnf("WS Connection UnexpectedCloseError: %v\n", err)
			}
			log.Info("WS Connection closed: %v\n", err)
			break
		}
	}
}
