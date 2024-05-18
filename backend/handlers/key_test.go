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
}

func TestEncryptDecrypt(t *testing.T) {
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

	shared, err := GenerateKey("127.0.0.1", string(RemotePublicKeyEncoded))

	if err != nil {
		log.Println(err)
		t.Error("Error generating key, from RemotePublicKeykey")
	}

	CachedKey, err := GetSharedKey("127.0.0.1")

	if err != nil {
		t.Error("Error getting key from cache")
	}

	if CachedKey != string(shared) {
		log.Println(CachedKey)
		t.Error("Keys do not match")
	}

	// Encrypt
	plaintext := []byte("Lorem Ipsum Dolor Sit Amet Consectetur Adipiscing Elit")
	ciphertext, err := Encrypt(shared, plaintext)

	log.Println(ciphertext)

	if err != nil {
		log.Println(err)
		t.Error("Error encrypting")
	}

	// Decrypt
	decrypted, err := Decrypt(shared, []byte(ciphertext))

	if err != nil {
		log.Println(err)
		t.Error("Error decrypting")
	}

	if string(plaintext) != decrypted {
		log.Println(string(plaintext))
		log.Println(decrypted)
		t.Error("Decrypted does not match plaintext")
	}
}