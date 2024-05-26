// + build js,wasm
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"syscall/js"

	"github.com/nart4hire/goblockc"
	"github.com/nart4hire/goschnorr"

	"github.com/teamortix/golang-wasm/wasm"
)

func GenSchnorrKeyPair(p, q, g string) (js.Value, error) {
	pi, succ := new(big.Int).SetString(p, 16)
	if !succ {
		return js.ValueOf(nil), errors.New("invalid p")
	}
	qi, succ := new(big.Int).SetString(q, 16)
	if !succ {
		return js.ValueOf(nil), errors.New("invalid q")
	}
	gi, succ := new(big.Int).SetString(g, 16)
	if !succ {
		return js.ValueOf(nil), errors.New("invalid gen")
	}

	s := schnorr.NewSchnorrFromParam(pi, qi, gi, rand.Reader, sha256.New())

	priv, pub, err := s.GenKeyPair()
	if err != nil {
		return js.ValueOf(nil), err
	}

	return js.ValueOf(map[string]interface{}{
		"private": hex.EncodeToString(priv),
		"public":  hex.EncodeToString(pub),
	}), nil
}

func Sign(p, q, g, privateKey, message string) (js.Value, error) {
	pi, succ := new(big.Int).SetString(p, 16)
	if !succ {
		return js.ValueOf(nil), errors.New("invalid p")
	}
	qi, succ := new(big.Int).SetString(q, 16)
	if !succ {
		return js.ValueOf(nil), errors.New("invalid q")
	}
	gi, succ := new(big.Int).SetString(g, 16)
	if !succ {
		return js.ValueOf(nil), errors.New("invalid gen")
	}

	s := schnorr.NewSchnorrFromParam(pi, qi, gi, rand.Reader, sha256.New())

	priv, err := hex.DecodeString(privateKey)
	if err != nil {
		return js.ValueOf(nil), err
	}

	sign, hash, err := s.Sign(priv, message)
	if err != nil {
		return js.ValueOf(nil), err
	}

	return js.ValueOf(map[string]interface{}{
		"sign": hex.EncodeToString(sign),
		"hash": hex.EncodeToString(hash),
	}),  nil
}

func Verify(p, q, g, publicKey, signature, hash, message string) (bool, error) {
	pi, succ := new(big.Int).SetString(p, 16)
	if !succ {
		return false, errors.New("invalid p")
	}
	qi, succ := new(big.Int).SetString(q, 16)
	if !succ {
		return false, errors.New("invalid q")
	}
	gi, succ := new(big.Int).SetString(g, 16)
	if !succ {
		return false, errors.New("invalid gen")
	}

	s := schnorr.NewSchnorrFromParam(pi, qi, gi, rand.Reader, sha256.New())

	pub, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	sign, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	hashb, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}

	return s.Verify(pub, sign, hashb, message), nil
}


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


func main() {
	wasm.Expose("keys", GenSchnorrKeyPair)
	wasm.Expose("sign", Sign)
	wasm.Expose("verify", Verify)
	wasm.Expose("encrypt", Encrypt)
	wasm.Expose("decrypt", Decrypt)
	wasm.Ready()

	select {}
}