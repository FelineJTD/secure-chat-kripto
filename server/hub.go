// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/big"

	"github.com/FelineJTD/secure-chat-kripto/server/ecdh"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients with private key, big int
	clients map[*Client]*big.Int

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]*big.Int),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// Generate a random key for the client
			privKey, pubKey := ecdh.GenerateKeyPair()
			fmt.Println("Generated key pair: ", privKey, pubKey)
			// Send the public key to the client
			// client.send <- []byte(pubKey)
			// Register the client
			h.clients[client] = privKey
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
