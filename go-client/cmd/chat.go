package cmd

import (
	"fmt"
	"go-client/lib/wschat"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var chat = &cobra.Command{
	Use:   "chat",
	Short: "Run AI chat",
	Run:   cmd_chat,
}

func init() {
	rootCmd.AddCommand(chat)
}

func cmd_chat(cmd *cobra.Command, args []string) {
	ch := wschat.New(8000)
	ch.Serve(cmd_chat_handler)
}

func cmd_chat_handler(c *websocket.Conn, w http.ResponseWriter, r *http.Request) {

	for {
		_, receivedMsg, err := c.ReadMessage()
		if err != nil {
			log.Println("Chat read error:", err)
			break
		}
		log.Printf("Received: %s", receivedMsg)

		// responseMsg, err := ai.Send(string(receivedMsg))
		// if err != nil {
		// 	log.Println("AI error:", err)
		// 	break
		// }
		responseMsg := fmt.Sprintf("Repsonse for %s", receivedMsg)

		if err := c.WriteMessage(websocket.TextMessage, []byte(responseMsg)); err != nil {
			log.Println("Chat write error:", err)
			break
		}
		log.Printf("Response: %s", responseMsg)
	}
}
