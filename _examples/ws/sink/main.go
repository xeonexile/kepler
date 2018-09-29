package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lastexile/kepler"
	"github.com/lastexile/kepler/ws"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("starting...")

	spring := kepler.NewSpring(func(ctx context.Context, ch chan<- kepler.Message) {

		i := 1
		for {
			ch <- kepler.NewMessage("odd", strconv.Itoa(i))
			i++
			time.Sleep(1 * time.Second)
			log.Println(i)
		}
	})

	broadcaster := kepler.NewBroadcastPipe(func(in kepler.Message) kepler.Message { return in })

	spring.LinkTo(".", broadcaster, kepler.Allways)

	url := "localhost:9090"

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// each ws connection is mapped to WsSink linked to broadcaster
		log.Println("opening connection")

		var closer func()

		onClose := func() {
			log.Println("closing")
			closer()
			log.Println("close")
		}
		sink, _ := ws.NewSink(ws.ServeConnection(w, r), ws.JSONValue, func(c *websocket.Conn) {
			ws.SendTextMessage(c, []byte("hi"))
		}, onClose)

		closer = broadcaster.LinkTo(".", sink, kepler.Allways)

	})

	go func() {
		err := http.ListenAndServe(url, nil)
		if err != nil {
			log.Fatal("Listen and serve error: ", err)
		} else {
			log.Println("ListenAndServe on: " + url)
		}
	}()

	kepler.Await()
}
