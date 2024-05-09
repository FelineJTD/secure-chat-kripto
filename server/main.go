package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender  int    `json:"sender"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Message received from client: " + string(p))

		// DO SOMETHING HERE
		message := Message{}
		json.Unmarshal(p, &message)
		// this just returns the same message to client
		msgToSend := "You said: " + string(message.Message)
		// payload in json with structure
		// { sender: "server", message: msgToSend }
		payload := []byte(`{"sender":"server","message":"` + msgToSend + `"}`)
		if err := conn.WriteMessage(messageType, payload); err != nil {
			log.Println(err)
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected")

  reader(ws)
}


func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/chat", wsEndpoint)
}

func main() {
	fmt.Println("Initiating server...")
	setupRoutes()
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}