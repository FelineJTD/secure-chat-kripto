package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/FelineJTD/secure-chat-kripto/server/ecdh"
	"github.com/FelineJTD/secure-chat-kripto/server/handlers"
	"github.com/FelineJTD/secure-chat-kripto/server/logger"
	// "github.com/FelineJTD/secure-chat-kripto/server/middlewares"
)

type Message struct {
	Sender  int    `json:"sender"`
	Message string `json:"message"`
}

type PublicKey struct {
	Port string `json:"port"`
	PublicKey string `json:"public_key"`
}

type Handshake struct {
	Port string `json:"port"`
	PublicKey string `json:"public_key"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.RemoteAddr)
	w.Write([]byte("Home Page"))
}

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

	msgJSON := Handshake{}
	if err = json.Unmarshal(body, &msgJSON); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pubKeyClientTemp := PubKeyClient{}

	if err = json.Unmarshal([]byte(msgJSON.PublicKey), &pubKeyClientTemp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	X, ok := new(big.Int).SetString(pubKeyClientTemp.X, 10)
	if !ok {
		err = errors.New("error parsing X from client public key")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Y, ok := new(big.Int).SetString(pubKeyClientTemp.Y, 10)
	if !ok {
		err = errors.New("error parsing Y from client public key")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pubKeyClient := ecdh.Point{X: X, Y: Y}

	address := strings.Split(r.RemoteAddr, ":")[0] + ":" + msgJSON.Port
	logger.Info("Shaking Hands With: " + address)

	sk, err := handlers.GenerateKey(address, &pubKeyClient)
	if err != nil {
		logger.HandleError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("Generated Shared Key: " + sk)

	// // Send the public key to the client as string
	// pubKeyX := pubKey.X.String()
	// pubKeyY := pubKey.Y.String()
	pubKeyJSON, err := json.Marshal(PubKeyClient{X: handlers.PubKey.X.String(), Y: handlers.PubKey.Y.String()})
	if err != nil {
		logger.HandleError(err)
		return
	}

	w.Write(pubKeyJSON)

	logger.Info("Handshake Complete")
}

func getParams(w http.ResponseWriter, r *http.Request) {
	p, q, gen := handlers.GetSchnorr()

	payload := []byte(`{"p": "` + hex.EncodeToString(p) + `", "q": "` + hex.EncodeToString(q) + `", "gen": "` + hex.EncodeToString(gen) + `"}`)
	w.Write(payload)
}

func setupRoutes(hub *Hub) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RealIP)

	r.Put("/key", keyEndpoint)

	r.Get("/schnorr", getParams)

	r.Get("/", homePage)

	r.Route("/chat", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			serveWs(hub, w, r)
		})
	})

	return r
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	r:= setupRoutes(hub)

	logger.Info("Server started at http://localhost:8080")
	logger.HandleFatal(http.ListenAndServe(":8080", r))
}