package handlers

import (
	"crypto/ecdh"
	"crypto/rand"

	"github.com/FelineJTD/secure-chat-kripto/server/providers"
	"github.com/FelineJTD/secure-chat-kripto/server/util"
)

var (
	PrivKey *ecdh.PrivateKey
)

func init() {
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		util.HandleFatal(err) // Fatal, because always needed
	}

	PrivKey = priv
}

// Generate Using Hash function and PRNG
func GenerateKey(address string, pubkey string) (string, error) {
	conn := providers.Pool.Get()
	defer util.HandleError(conn.Err())
	defer conn.Close()

	// TODO: Implement And Test
	// pubkey := ecdh.PublicKey(pubkey)
	// key := PrivKey.ECDH(pubkey)
	key := "generatedkey"

	_, err := conn.Do("SET", address, key)

	if err != nil {
		return "", err
	}

	return key, nil
}

// Get Server's Public Key, for client to generate shared key
func GetPubKey() string {
	return string(PrivKey.PublicKey().Bytes())
}

// Get Shared Key from cache, untested!!
func GetSharedKey(address string) (string, error) {
	conn := providers.Pool.Get()
	defer util.HandleError(conn.Err())
	defer conn.Close()

	key, err := conn.Do("GET", address)
	if err != nil {
		return "", err
	}

	return string(key.([]byte)), nil
}
