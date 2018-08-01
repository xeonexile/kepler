package ws

import (
	"context"
	"time"

	"github.com/lastexile/kepler"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Connection returns based on Default dialer connectionFactory
func Connection(path string) func() (*websocket.Conn, error) {
	return func() (*websocket.Conn, error) {
		conn, _, err := websocket.DefaultDialer.Dial(path, nil)
		return conn, err
	}
}

// New creates instance of web socket spring
func New(ctx context.Context, connFactory func() (*websocket.Conn, error), formatter kepler.UnmarshallFunction, onConnect func(*websocket.Conn)) (kepler.Spring, error) {

	return kepler.NewSpring("foo", func(ctx context.Context, ch chan<- kepler.Message) {
		var writeCtx context.Context
		var writeCancel context.CancelFunc
		var conn *websocket.Conn
		var err error
		for {
			if conn == nil {
				log.Info(conn)
				for {
					conn, err = connFactory()
					if err != nil {
						log.Errorf("Failed to open the connection: %v \n", err)
						time.Sleep(10 * time.Second)
						continue
					}

					log.Infoln("Connected")
					writeCtx, writeCancel = context.WithCancel(context.Background())
					go writePump(writeCtx, conn, pingPeriod)
					onConnect(conn)
					break
				}
			}

			bytesRead, data, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("Read Websocket Error: %v. %v bytes \n", err, bytesRead)
				writeCancel()
				conn.Close()
				conn = nil

			}

			if msg, err := formatter(data); err == nil {
				ch <- msg
			} else {
				log.Errorf("formatter error %v\n", err)
			}
		}
	}), nil
}

func writePump(ctx context.Context, conn *websocket.Conn, pingPeriod time.Duration) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Infoln("Closing")
		ticker.Stop()
		conn.Close()
		log.Infoln("Closed")
	}()

	for {
		select {
		case <-ticker.C:
			if err := write(conn, websocket.PingMessage, []byte{}); err != nil {
				log.Errorln("PING error: %v\n", err)
				return
			}
			log.Infoln("PING")
		case <-ctx.Done():
			log.Infoln("ctx Done")
			return
		}
	}
}

// write writes a message with the given message type and payload.
func write(conn *websocket.Conn, mt int, payload []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(mt, payload)
}
