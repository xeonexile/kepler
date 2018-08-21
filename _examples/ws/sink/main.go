package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lastexile/kepler"
	"github.com/lastexile/kepler/ws"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	spring := kepler.NewSpring("odd", func(ctx context.Context, ch chan<- kepler.Message) {

		i := 1
		for {
			ch <- kepler.NewMessage("odd", strconv.Itoa(i))
			i++
			time.Sleep(1 * time.Second)
			log.Println(i)
		}
	})

	broadcaster := kepler.NewBroadcastPipe("all", func(in kepler.Message) kepler.Message { return in })

	spring.LinkTo(broadcaster, kepler.Allways)

	url := "wss://abyss-unifeed-develop.marlin.onnisoft.com"
	url = "localhost:9090"

	go func() {
		err := http.ListenAndServe(url, nil)
		if err != nil {
			log.Fatal("Listen and serve error: ", err)
		} else {
			log.Println("ListenAndServe on: " + url)
		}
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// each ws connection is mapped to WsSink linked to broadcaster
		log.Println("opening connection")

		sink, _ := ws.NewSink(ws.ServeConnection(w, r), ws.JsonValue, func(c *websocket.Conn) {
			ws.SendTextMessage(c, []byte("hi"))
		}, func() {
			log.Println("close")
		})

		broadcaster.LinkTo(sink, kepler.Allways)

	})

	reader := bufio.NewReader(os.Stdin)
	log.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	log.Println(text)
}
