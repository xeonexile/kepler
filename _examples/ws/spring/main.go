package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/lastexile/kepler"
	"github.com/lastexile/kepler/ws"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	url := "wss://abyss-unifeed-develop.marlin.onnisoft.com/ws"
	url = "ws://localhost:9090/ws"
	s, err := ws.NewSpring(context.Background(), ws.DialConnection(url), func(d []byte) (kepler.Message, error) {
		return kepler.NewMessage("foo", string(d)), nil
	}, func(conn *websocket.Conn) {
		log.Println("connected")
	})
	if err != nil {
		log.Fatalf("Unable to create wsspring: %v\n", err)
	}

	logSink := kepler.NewSink("odd", func(m kepler.Message) {

		log.Println(m.String())
	})

	s.LinkTo(logSink, kepler.Allways)

	reader := bufio.NewReader(os.Stdin)
	log.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	log.Println(text)
}
