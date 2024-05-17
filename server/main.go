package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FelineJTD/secure-chat-kripto/server/handlers"
	"github.com/FelineJTD/secure-chat-kripto/server/logger"
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
	var err error = nil
	defer logger.HandleError(err)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
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
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	defer logger.HandleError(err)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	logger.Info("Client Connected")

	reader(ws)
}

func keyEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	defer logger.HandleError(err)
	pubKey := r.URL.Query().Get("key")

	logger.Info("Generating Key for " + r.RemoteAddr)
	key, err := handlers.GenerateKey(r.RemoteAddr, pubKey)

	logger.Info("Key Generated: " + key)
	if err != nil {
		return
	}

	fmt.Fprintf(w, key)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/chat", wsEndpoint)
	http.HandleFunc("/key", keyEndpoint)
}

func main() {
	var err error = nil
	defer logger.HandleError(err)
	fmt.Println("Initiating server...")
	setupRoutes()
	fmt.Println("Server started at http://localhost:8080")
	logger.HandleFatal(http.ListenAndServe(":8080", nil))
}
