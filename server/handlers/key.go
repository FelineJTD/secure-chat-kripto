package handlers

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"crypto/ecdh" // TODO: implement own ECDH
	// "crypto/aes" // TODO: replace with GoBlockC
	"github.com/nart4hire/goblockc"

	"github.com/gomodule/redigo/redis"

	"github.com/FelineJTD/secure-chat-kripto/server/logger"
	"github.com/FelineJTD/secure-chat-kripto/server/providers"
)

var (
	PrivKey *ecdh.PrivateKey
)

func init() {
	// TODO: Use BBS for randomness
	if priv, err := ecdh.X25519().GenerateKey(rand.Reader); err != nil {
		logger.HandleFatal(err) // Fatal, because always needed
	} else {
		PrivKey = priv
	}
}

func StringToPubKey(pubkey string) (*ecdh.PublicKey, error) {
	decoded, _ := pem.Decode([]byte(pubkey))

	parsed, err := x509.ParsePKIXPublicKey(decoded.Bytes)

	if err != nil {
		return nil, err
	}

	return parsed.(*ecdh.PublicKey), nil
}

func GenerateSharedKey(pub *ecdh.PublicKey) (string, error) {
	key, err := PrivKey.ECDH(pub)

	if err != nil {
		return "", err
	}

	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "SHARED KEY",
			Bytes: key,
		})), nil
}

// Generate Using Hash function and PRNG
func GenerateKey(address string, pubkey string) (string, error) {
	conn := providers.Pool.Get()
	defer logger.HandleError(conn.Err())
	defer conn.Close()

	pub, err := StringToPubKey(pubkey)

	if err != nil {
		return "", err
	}

	key, err := GenerateSharedKey(pub)

	if err != nil {
		return "", err
	}

	// Store Key in Cache
	if _, err = conn.Do("SET", address, key); err != nil {
		return "", err
	}

	logger.Debug("Key Generated for: " + address + " => " + key)

	return key, nil
}

// Get Server's Public Key, for client to generate shared key
func GetPubKey() (string, error) {
	key, err := x509.MarshalPKIXPublicKey(PrivKey.PublicKey())
	if err != nil {
		return "", err
	}

	keyEnc := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: key,
		}))

	return keyEnc, nil
}

// Get Shared Key from cache
func GetSharedKey(address string) (string, error) {
	conn := providers.Pool.Get()
	defer logger.HandleError(conn.Err())
	defer conn.Close()

	key, err := redis.String(conn.Do("GET", address))
	if err != nil {
		return "", err
	}

	return key, nil
}

func Encrypt(key string, plaintext []byte) (string, error) {
	gbc, err := goblockc.NewBlock([]byte(key)[:16])

	if err != nil {
		return "", err
	}

	iv := make([]byte, gbc.BlockSize())

    ctr, err := goblockc.NewCTR(gbc, iv)

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plaintext))
	copy(ciphertext, plaintext)
    ctr.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func Decrypt(key string, ciphertext []byte) (string, error) {
	gbc, err := goblockc.NewBlock([]byte(key)[:16])

	if err != nil {
		return "", err
	}

	iv := make([]byte, gbc.BlockSize())

    ctr, err := goblockc.NewCTR(gbc, iv)

	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertext))
	copy(plaintext, ciphertext)
    ctr.XORKeyStream(plaintext, plaintext)

	return string(plaintext), nil
}