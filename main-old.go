package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/FelineJTD/secure-chat-kripto/server/handlers"
	"github.com/FelineJTD/secure-chat-kripto/server/logger"
	// "github.com/FelineJTD/secure-chat-kripto/server/middlewares"
)

type Message struct {
	Sender  int    `json:"sender"`
	Message string `json:"message"`
}

type PublicKey struct {
	PublicKey string `json:"public_key"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(key string, conn *websocket.Conn) {
	var err error = nil
	defer logger.HandleError(err)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		logger.Info("Message received from client: " + string(p))

		// DO SOMETHING HERE
		message := Message{}
		json.Unmarshal(p, &message)
		// this just returns the same message to client
		msgToSend := "You said: " + string(message.Message)
		// payload in json with structure
		// { sender: "server", message: msgToSend }
		payload := []byte(`{"sender":"server","message":"` + msgToSend + `"}`)

		logger.Debug("Shared Key: " + key) // This is just to silence the linter
		// TODO: Uncomment this to enable encryption, need testing
		// payload := handlers.Encrypt(key, payload) // some json marshalled version of this

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

	reader(r.RemoteAddr, ws)
}

// TODO: Test this endpoint
// Since the spec requested a handshake, It might be better to emulate it using a websocket, but this will do for now
// In essence the client makes a PUT request sending its public key, the server then generates a shared key and sends back its public key
// The client then calculates the shared key and can now send encrypted messages
func keyEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	defer logger.HandleError(err)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pubKey := PublicKey{}
	json.Unmarshal(body, &pubKey)

	SharedKey, err := handlers.GenerateKey(r.RemoteAddr, pubKey.PublicKey)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("Shared Key Generated: " + SharedKey) // We dont really need to send back the shared key because the client can calculate it for himself, I leave it up to your judgement

	key, err := handlers.GetPubKey()

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	payload := []byte(`{"public_key": "` + key + `"}`)

	w.Write(payload)

	logger.Info("Public Key Sent")
}

func setupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", homePage)
	r.Get("/chat", wsEndpoint)

	// TODO: Uncomment this to enable decryption middleware, need testing
	// r.Route("/chat", func(r chi.Router) {
	// 	r.Use(middlewares.DecryptMiddleware)
	// 	r.Get("/", wsEndpoint)
	// })

	r.Put("/key", keyEndpoint)

	return r
}

func main() {
	logger.Info("Initiating server...")
	r := setupRoutes()

	logger.Info("Server started at http://localhost:8080")
	logger.HandleFatal(http.ListenAndServe(":8080", r))
}
