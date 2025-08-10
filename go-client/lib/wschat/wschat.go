package wschat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wschat struct {
	port         int
	staticDir    string
	interlocutor func(inputMsg string) (string, error)
}

func New(port int, interlocutor func(inputMsg string) (string, error)) *wschat {
	ch := &wschat{
		port:         port,
		staticDir:    "./static",
		interlocutor: interlocutor,
	}
	return ch
}

func (ch *wschat) Serve() {

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

		log.Println("New WebSocket connection")
		for {
			_, receivedMsg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Chat read error:", err)
				break
			}
			log.Printf("Received: %s", receivedMsg)

			responseMsg, err := ch.interlocutor(string(receivedMsg))
			if err != nil {
				log.Println("Interlocutor error:", err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(responseMsg)); err != nil {
				log.Println("Chat write error:", err)
				break
			}
			log.Printf("Response: %s", responseMsg)
		}
	})

	log.Printf("Serwer starting on :%d\n", ch.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ch.port), nil))
}
