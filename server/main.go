// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/FelineJTD/secure-chat-kripto/server/handlers"
	"github.com/FelineJTD/secure-chat-kripto/server/logger"
	// "github.com/FelineJTD/secure-chat-kripto/server/middlewares"
)

var addr = flag.String("addr", ":8080", "http service address")

type Message struct {
	Sender  int    `json:"sender"`
	Message string `json:"message"`
}

type PublicKey struct {
	PublicKey string `json:"public_key"`
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

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

func setupRoutes(hub *Hub) http.Handler {
	r := chi.NewRouter()

	r.Get("/", serveHome)
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// TODO: Uncomment this to enable decryption middleware, need testing
	// r.Route("/chat", func(r chi.Router) {
	// 	r.Use(middlewares.DecryptMiddleware)
	// 	r.Get("/", wsEndpoint)
	// })

	r.Put("/key", keyEndpoint)

	return r
}

func main() {
	fmt.Println("Initiating server...")
	flag.Parse()
	hub := newHub()
	go hub.run()

	setupRoutes(hub)
	fmt.Println("Server started at http://localhost:8080")

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 5 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
