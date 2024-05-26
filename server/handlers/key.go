package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/nart4hire/goblockc"
	"github.com/nart4hire/goschnorr"

	"github.com/gomodule/redigo/redis"

	"github.com/FelineJTD/secure-chat-kripto/server/ecdh"
	"github.com/FelineJTD/secure-chat-kripto/server/logger"
	"github.com/FelineJTD/secure-chat-kripto/server/providers"
)

var (
	PrivKey *big.Int
	PubKey *ecdh.Point
	Schnorr schnorr.Schnorr
)

func init() {
	priv, pub := ecdh.GenerateKeyPair()
	PrivKey = priv
	PubKey = pub

	if schnorr, err := schnorr.NewSchnorr(rand.Reader, sha256.New()); err != nil {
		logger.HandleFatal(err) // Fatal, because always needed
	} else {
		Schnorr = schnorr
	}
}

// func StringToPubKey(pubkey string) (*ecdh.PublicKey, error) {
// 	decoded, _ := pem.Decode([]byte(pubkey))

// 	parsed, err := x509.ParsePKIXPublicKey(decoded.Bytes)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return parsed.(*ecdh.PublicKey), nil
// }

// func GenerateSharedKey(pub *ecdh.PublicKey) (string, error) {
// 	key, err := PrivKey.ECDH(pub)

// 	if err != nil {
// 		return "", err
// 	}

// 	return string(pem.EncodeToMemory(
// 		&pem.Block{
// 			Type:  "SHARED KEY",
// 			Bytes: key,
// 		})), nil
// }

// Generate Using Hash function and PRNG
func GenerateKey(address string, pubkey *ecdh.Point) (string, error) {
	conn := providers.Pool.Get()
	defer logger.HandleError(conn.Err())
	defer conn.Close()

	key := ecdh.GenerateSharedKey(PrivKey, pubkey)
	keyHash, err := Hash(key.Text(16))
	if err != nil {
		return "", err
	}

	// Store Key in Cache
	if _, err := conn.Do("SET", address, keyHash); err != nil {
		return "", err
	}

	logger.Debug("Key Generated for: " + address + " => " + keyHash)

	return keyHash, nil
}

// Get Server's Public Key, for client to generate shared key
func GetPubKey() (*ecdh.Point, error) {
	return PubKey, nil
}

// Get Shared Key from cache, key is Hex encoded
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

// Key is always hex encoded, String is UTF-8, converted plainly
func Encrypt(key, plaintext string) (string, error) {
	k, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	gbc, err := goblockc.NewBlock(k[:16])
	if err != nil {
		return "", err
	}

	iv := make([]byte, gbc.BlockSize())

    ctr, err := goblockc.NewCTR(gbc, iv)

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plaintext))
	copy(ciphertext, []byte(plaintext))
    ctr.XORKeyStream(ciphertext, ciphertext)

	return hex.EncodeToString(ciphertext), nil
}

// Key is always hex encoded, ciphertext is also hex encoded to preserve data
func Decrypt(key, ciphertext string) (string, error) {
	k, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	gbc, err := goblockc.NewBlock(k[:16])
	if err != nil {
		return "", err
	}

	iv := make([]byte, gbc.BlockSize())

    ctr, err := goblockc.NewCTR(gbc, iv)

	if err != nil {
		return "", err
	}

	c, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(c))
	copy(plaintext, c)
    ctr.XORKeyStream(plaintext, plaintext)

	return string(plaintext), nil
}

func GetSchnorr() ([]byte, []byte, []byte) {
	p, q, gen := Schnorr.GetParams()

	return p.Bytes(), q.Bytes(), gen.Bytes()
}

func Hash(hexString string) (string, error) {
	if len(hexString) % 2 != 0 {
		hexString = "0" + hexString
	}

	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(bytes)

	return hex.EncodeToString(hash[:]), nil
}