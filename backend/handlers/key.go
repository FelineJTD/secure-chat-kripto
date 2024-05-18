package handlers

import (
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"crypto/aes" // replace with GoBlockC

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

// TODO: Replace with GoBlockC
func Encrypt(key string, plaintext []byte) (string, error) {
	// TODO: Hash this key instead to get some bytes if needed
	aes, err := aes.NewCipher([]byte(key)[:32])

	if err != nil {
		return "", err
	}

    gcm, err := cipher.NewGCM(aes)

    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    _, err = rand.Read(nonce)

    if err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

func Decrypt(key string, ciphertext []byte) (string, error) {
	aes, err := aes.NewCipher([]byte(key)[:32])

	if err != nil {
		return "", err
	}

    gcm, err := cipher.NewGCM(aes)

    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
    if err != nil {
        return "", err
    }

	return string(plaintext), nil
}