package ws

import (
	"context"
	"time"

	"github.com/lastexile/kepler"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// NewSpring creates instance of web socket spring. Acts as ws client
func NewSpring(ctx context.Context, connFactory ConnectionFactoryFunc, formatter kepler.UnmarshalFunction, onConnect func(*websocket.Conn)) (kepler.Spring, error) {

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
						time.Sleep(connectionRetryInterval)
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
