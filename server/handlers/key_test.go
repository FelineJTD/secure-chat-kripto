package handlers

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"log"
	"testing"

	"github.com/FelineJTD/secure-chat-kripto/server/logger"
)

func TestGenerateKeyAndCache(t *testing.T) {
	logger.SetVerbosity(3)

	RemotePrivateKey, err := ecdh.X25519().GenerateKey(rand.Reader)

	if err != nil {
		log.Println(err)
		t.Error("Library error generating key")
	}

	RemotePublicKey, err := x509.MarshalPKIXPublicKey(RemotePrivateKey.PublicKey())
	RemotePublicKeyEncoded := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: RemotePublicKey,
		})

	if err != nil {
		log.Println(err)
		t.Error("Library error marshalling key")
	}

	keyOne, err := GenerateKey("127.0.0.1", string(RemotePublicKeyEncoded))

	if err != nil {
		log.Println(err)
		t.Error("Error generating key, from RemotePublicKeykey")
	}

	ServerPublicKey, err := GetPubKey()

	if err != nil {
		t.Error("Error getting server public key")
	}

	decoded, _ := pem.Decode([]byte(ServerPublicKey))
	ServerPublicKeyParsed, err := x509.ParsePKIXPublicKey(decoded.Bytes)

	if err != nil {
		t.Error("Error parsing server public key")
	}

	keyTwoRaw, err := RemotePrivateKey.ECDH(ServerPublicKeyParsed.(*ecdh.PublicKey))

	keyTwo := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "SHARED KEY",
			Bytes: keyTwoRaw,
		}))

	if err != nil {
		log.Println(err)
		t.Error("Error generating key, from RemotePrivateKey")
	}

	if keyOne != string(keyTwo) {
		log.Println(keyOne)
		log.Println(keyTwo)
		t.Error("Keys do not match")
	}

	CachedKey, err := GetSharedKey("127.0.0.1")

	if err != nil {
		t.Error("Error getting key from cache")
	}

	if CachedKey != string(keyOne) {
		log.Println(CachedKey)
		t.Error("Keys do not match")
	}
}
