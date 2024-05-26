package ecdh

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Define the Point structure
type Point struct {
	X *big.Int
	Y *big.Int
}

// Define the Curve structure
type Curve struct {
	a  *big.Int
	b  *big.Int
	p  *big.Int
	g  Point
}

// Initialize the curve parameters
func NewCurve() *Curve {
	p := new(big.Int)
	p.SetString("501093B9D01F6A995131F06A99E971704FBED42E0260730567E466958EAF21A5ED4B8458626AFFC325C1172017291D86E8E6818AED1BC6D59AD1E773DD9F2401B577263802C21D82870008D37DB9C72122F4E190EBEA973382E5C18EB3DBFB565EAAB90B98EE880104767", 16)

	a := new(big.Int)
	a.SetString("4EB5DF9B1A22357A539C1F8EE0DBD02F8F5373C1F3DE064A2547D8DF185A0B9411BE82EB93FB61BA51C1EA46A35141B5BBC8083CD642B1F6419BD0263C61C6BA128EE64B224BCCE25A1794C30E20DBEEF8B163DC9662EC0739455849AD2AEAE67CCCDBC84968674D299BF", 16)
	b := big.NewInt(0)

	gx := new(big.Int)
	gx.SetString("1471CFA725EB7FB877EC8F8DE8B3DD9E6F3B880BDD984289BB180E372968D6CDA1A667AF0B859FFC1A2700B8853541FF0416AC5D3C7B00EFDD550614B1956618413453C90E06ED4EC0060204CDC287F140B3153E161D078452734A03510532E3EABAE6E105FFFE2656D59", 16)

	gy := new(big.Int)
	gy.SetString("190AE1E86693B1EE024FE0811766A67437F43AE62FE6E431A83EC7ACE5BB5A58AB8C9D296A08622C8ED6572BC9F38B5FCF8FB883C8A8E673C483A978ED3C4A4D9D7CE8F0B7A732107EB4AF7EA85D99B9BDF7C465831A3FB76C8451AF68756E9E6B8FA95EF5A317EF4011F", 16)

	g := Point{gx, gy}

	return &Curve{a, b, p, g}
}

// Modular inverse: returns x such that (x * k) % p == 1
func modInverse(k, p *big.Int) *big.Int {
	return new(big.Int).ModInverse(k, p)
}

// Point addition: R = P + Q
func (curve *Curve) Add(p, q *Point) *Point {
	if p.X == nil && p.Y == nil {
		return q
	}
	if q.X == nil && q.Y == nil {
		return p
	}

	if p.X.Cmp(q.X) == 0 && p.Y.Cmp(new(big.Int).Neg(q.Y)) == 0 {
		return &Point{nil, nil}
	}

	if p.X.Cmp(q.X) == 0 && p.Y.Cmp(q.Y) == 0 {
		// Point doubling
		dX := new(big.Int).Add(new(big.Int).Mul(big.NewInt(3), new(big.Int).Mul(p.X, p.X)), curve.a)
		dY := new(big.Int).Mul(big.NewInt(2), p.Y)
		inv := modInverse(dY, curve.p)
		slope := new(big.Int).Mul(dX, inv)
		slope_mod := new(big.Int).Mod(slope, curve.p)
		rx := new(big.Int).Sub(new(big.Int).Mul(slope_mod, slope_mod), new(big.Int).Mul(big.NewInt(2), p.X))
		rx.Mod(rx, curve.p)
		ry := new(big.Int).Sub(new(big.Int).Mul(slope_mod, new(big.Int).Sub(p.X, rx)), p.Y)
		ry.Mod(ry, curve.p)
		return &Point{rx, ry}
	} else {
		dX := new(big.Int).Sub(q.X, p.X)
		dY := new(big.Int).Sub(q.Y, p.Y)
		if dX.Cmp(big.NewInt(0)) == -1 {
			dX.Neg(dX)
			dY.Neg(dY)
		}
		inv := modInverse(dX, curve.p)
		slope := new(big.Int).Mod(new(big.Int).Mul(dY, inv), curve.p)
		rx := new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(slope, slope), new(big.Int).Add(p.X, q.X)), curve.p)
		ry := new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(slope, new(big.Int).Sub(p.X, rx)), p.Y), curve.p)
		return &Point{rx, ry}
	}
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
	p := curve.p
	d, err := rand.Int(rand.Reader, p)
	if err != nil {
		fmt.Println(err)
	}
	return d
}

// Generate a public key
func GeneratePublicKey(curve *Curve, d *big.Int) *Point {
	return curve.ScalarMult(d, &curve.g)
}

// ECDH: generate a shared key
func GenerateKeyPair() (*big.Int, *Point) {
	curve := NewCurve()
	privKey := GeneratePrivateKey(curve)
	pubKey := GeneratePublicKey(curve, privKey)
	return privKey, pubKey
}

// ECDH: generate a shared key
func GenerateSharedKey(privKey *big.Int, pubKey *Point) *big.Int {
	fmt.Println("Private Key: ", privKey)
	fmt.Println("Pub Key: ", pubKey.X, pubKey.Y)
	curve := NewCurve()
	sharedKey := curve.ScalarMult(privKey, pubKey)
	fmt.Println("Shared Key: ", sharedKey)
	return sharedKey.X
}