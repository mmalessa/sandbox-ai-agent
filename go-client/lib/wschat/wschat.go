package wschat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wschat struct {
	port      int
	staticDir string
}

func New(port int) *wschat {
	ch := &wschat{
		port:      port,
		staticDir: "./static",
	}
	return ch
}

func (ch *wschat) Serve(handler func(c *websocket.Conn, w http.ResponseWriter, r *http.Request)) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // we allow to connect from any source
		},
	}

	http.Handle("/", http.FileServer(http.Dir(ch.staticDir)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()
		handler(conn, w, r)
	})

	log.Printf("Serwer starting on :%d\n", ch.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ch.port), nil))
}
