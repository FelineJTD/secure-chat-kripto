package ecdh

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Define the Point structure
type Point struct {
	x *big.Int
	y *big.Int
}

// Define the Curve structure
type Curve struct {
	a  *big.Int
	b  *big.Int
	p  *big.Int
	g  Point
	n  *big.Int
	h  *big.Int
}

// Initialize the curve parameters (example with secp256k1)
func NewSecp256k1Curve() *Curve {
	p := new(big.Int)
	p.SetString("B051", 16)

	a := big.NewInt(38699)
	b := big.NewInt(7)

	gx := new(big.Int)
	gx.SetString("4613", 16)

	gy := new(big.Int)
	gy.SetString("746B", 16)

	g := Point{gx, gy}

	n := new(big.Int)
	n.SetString("B051", 16)

	h := big.NewInt(1)

	return &Curve{a, b, p, g, n, h}
}

// Modular inverse: returns x such that (x * k) % p == 1
func modInverse(k, p *big.Int) *big.Int {
	return new(big.Int).ModInverse(k, p)
}

// Point addition: R = P + Q
func (curve *Curve) Add(p, q *Point) *Point {
	if p.x == nil && p.y == nil {
		return q
	}
	if q.x == nil && q.y == nil {
		return p
	}

	if p.x.Cmp(q.x) == 0 && p.y.Cmp(new(big.Int).Neg(q.y)) == 0 {
		return &Point{nil, nil}
	}

	lambda := new(big.Int)
	if p.x.Cmp(q.x) == 0 && p.y.Cmp(q.y) == 0 {
		// Point doubling
		dX := new(big.Int).Add(new(big.Int).Mul(big.NewInt(3), new(big.Int).Mul(p.x, p.x)), curve.a)
		dY := new(big.Int).Mul(big.NewInt(2), p.y)
		inv := modInverse(dY, curve.p)
		slope := new(big.Int).Mul(dX, inv)
		slope_mod := new(big.Int).Mod(slope, curve.p)
		rx := new(big.Int).Sub(new(big.Int).Mul(slope_mod, slope_mod), new(big.Int).Mul(big.NewInt(2), p.x))
		rx.Mod(rx, curve.p)
		ry := new(big.Int).Sub(new(big.Int).Mul(slope_mod, new(big.Int).Sub(p.x, rx)), p.y)
		ry.Mod(ry, curve.p)
		return &Point{rx, ry}
	} else {
		// Point addition
		num := new(big.Int).Sub(q.y, p.y)
		den := new(big.Int).Sub(q.x, p.x)
		lambda.Mul(num, modInverse(den, curve.p))
		lambda.Mod(lambda, curve.p)
	}

	rx := new(big.Int).Sub(new(big.Int).Mul(lambda, lambda), p.x)
	rx.Sub(rx, q.x)
	rx.Mod(rx, curve.p)

	ry := new(big.Int).Sub(new(big.Int).Mul(lambda, new(big.Int).Sub(p.x, rx)), p.y)
	ry.Mod(ry, curve.p)

	return &Point{rx, ry}
}

// Point doubling: R = 2P
func (curve *Curve) Double(p *Point) *Point {
	return curve.Add(p, p)
}

// Scalar multiplication: R = kP
func (curve *Curve) ScalarMult(k *big.Int, p *Point) *Point {
	res := &Point{nil, nil}
	addend := p

	for k.Cmp(big.NewInt(0)) != 0 {
		if new(big.Int).And(k, big.NewInt(1)).Cmp(big.NewInt(0)) != 0 {
			res = curve.Add(res, addend)
		}
		addend = curve.Double(addend)
		k.Rsh(k, 1)
	}

	return res
}

// Generate a private key
func GeneratePrivateKey(curve *Curve) *big.Int {
	n := curve.n
	d, err := rand.Int(rand.Reader, n)
	if err != nil {
		fmt.Println(err)
	}
	return d
}

// func GenerateRandomBigIntBelow(limit *big.Int) (*big.Int, error) {
// 	n, err := rand.Int(rand.Reader, limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return n, nil
// }


// Generate a public key
func GeneratePublicKey(curve *Curve, d *big.Int) *Point {
	return curve.ScalarMult(d, &curve.g)
}

// ECDH: generate a shared key
func GenerateKeyPair() (*big.Int, *Point) {
	curve := NewSecp256k1Curve()
	privKey := GeneratePrivateKey(curve)
	pubKey := GeneratePublicKey(curve, privKey)
	return privKey, pubKey
}

// ECDH: generate a shared key
func GenerateSharedKey(privKey *big.Int, pubKey *Point) *big.Int {
	curve := NewSecp256k1Curve()
	sharedKey := curve.ScalarMult(privKey, pubKey)
	fmt.Println("Shared Key: ", sharedKey)
	return sharedKey.x
}

func TestDoubling() {
	gx := new(big.Int)
	gx.SetString("4613", 16)

	gy := new(big.Int)
	gy.SetString("746B", 16)

	curve := NewSecp256k1Curve()
	p := Point{gx, gy}
	fmt.Println("Generator: ", p.x, p.y)
	add := curve.Add(&p, &p)
	fmt.Println("Addition: ", add.x, add.y)
}