package ws

import (
	"time"

	"github.com/lastexile/kepler"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 3 * 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// New creates new WS outgoing sink
func New(connFactory func() (*websocket.Conn, error), formattter kepler.MarshallerFunc, onConnect func(conn *websocket.Conn), onClose func()) (sink kepler.Sink, err error) {
	conn, err := connFactory()
	go readPump(conn, onClose)
	onConnect(conn)

	sink = kepler.NewSink("ws", func(m kepler.Message) {

		value, err := formattter(m)
		if err != nil {
			log.Error("Unable to marshall message value")
			return
		}

		SendTextMessage(conn, value)
	})

	return
}

// SendTextMessage data over opened conn
func SendTextMessage(conn *websocket.Conn, data []byte) {
	writer, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Errorf("Unable to create response TextMessage: %v\n", err)
		return
	}

	if err == nil {
		log.Info("Sending current state for new session")
		n, err := writer.Write(data)
		if err != nil {
			log.Warnf("Unable to write TextMessage: %v", err)
		} else {
			log.Infof("%v bytes sent", n)
		}
	}

	if err = writer.Close(); err != nil {
		log.Errorf("writer close error: %v\n", err)
		return
	}
}

func readPump(conn *websocket.Conn, onClose func()) {
	defer func() {
		conn.Close()
		onClose()
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
				if onClose != nil {
					onClose()
					log.Info("Connction Closed with onClose ")
				} else {
					log.Warn("No OnClose handler specified")
				}
			}
			break
		}
	}
}
