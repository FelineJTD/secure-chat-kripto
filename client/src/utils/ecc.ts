// max: 115792089210356248762697446949407573530086143415290314195533631308867097853951
// curve: y² = x³ + ax + b
// a = 115792089210356248762697446949407573530086143415290314195533631308867097853948
// b = 41058363725152142129326129780047268409114441015993725554835256314039467401291

// Parameters for the curve
const p = 270292550838977n;
const a = 305656052155099n;
const b = 150865464689957n;
const Gx = 415909356139339n;
const Gy = 641599020309769n;
// const n =
//   115792089210356248762697446949407573529996955224135760342422259061068512044369n;

// Define a point class to represent points on the curve
export class Point {
  x: bigint;
  y: bigint;

  constructor(x: bigint, y: bigint) {
    this.x = x;
    this.y = y;
  }

  // comparison function for points
  equals(other: Point): boolean {
    return this.x === other.x && this.y === other.y;
  }

  sameXDiffY(other: Point): boolean {
    return this.x === other.x && this.y !== other.y;
  }
}

function modInverse(a: bigint, m: bigint): bigint {
  a = ((a % m) + m) % m;
  for (let x = 1n; x < m; x++) {
    if ((a * x) % m === 1n) {
      return x;
    }
  }
  return 1n;
}

// Point addition on the curve
function pointAddition(p1: Point, p2: Point): Point {
  if (p1.equals(p2)) {
    return pointDoubling(p1);
  } else if (p1.sameXDiffY(p2)) {
    return new Point(BigInt(0), BigInt(0));
  } else {
    // Calculate slope
    const slope = (p2.y - p1.y) * modInverse(p2.x - p1.x, p);
    // Calculate x-coordinate
    const x = (slope * slope - p1.x - p2.x) % p;
    // Calculate y-coordinate
    const y = (slope * (p1.x - x) - p1.y) % p;
    return new Point(x, y);
  }
}

// Point doubling on the curve
function pointDoubling(p1: Point): Point {
  // Calculate slope
  const slope = (3n * p1.x * p1.x + a) * modInverse(2n * p1.y, p);
  // Calculate x-coordinate
  const x = (slope * slope - 2n * p1.x) % p;
  // Calculate y-coordinate
  const y = (slope * (p1.x - x) - p1.y) % p;
  return new Point(x, y);
}

// Scalar multiplication on the curve
function scalarMultiply(k: bigint, p1: Point): Point {
  let result = new Point(BigInt(0), BigInt(0));
  let addend = p1;
  while (k > 0) {
    if (k % 2n === 1n) {
      result = pointAddition(result, addend);
    }
    addend = pointDoubling(addend);
    k = k / 2n;
  }
  return result;
}

// Generate a random private key
export function generatePrivateKey(): bigint {
  const hexString = Array(16)
    .fill(0)
    .map(() => Math.round(Math.random() * 0xf).toString(16))
    .join("");

  const randomBigInt = BigInt(`0x${hexString}`);
  return randomBigInt;
}

// Generate the corresponding public key
export function generatePublicKey(privateKey: bigint): Point {
  return scalarMultiply(privateKey, new Point(Gx, Gy));
}

// Encrypt a message using ECC
export function encryptMessage(publicKey: Point, message: string): string {
  const k = generatePrivateKey();
  // a = g^k mod p
  const a = scalarMultiply(k, new Point(Gx, Gy));
  // b = m*P^k mod p
  const m = BigInt(message);
  const b = pointAddition(scalarMultiply(k, publicKey), new Point(m, m));

  return `${a},${b},${message}`;
}

// Decrypt a message using ECC
export function decryptMessage(privateKey: bigint, ciphertext: string): string {
  const [a, b, message] = ciphertext.split(",");
  const aPoint = new Point(BigInt(a), BigInt(b));
  const m = pointAddition(aPoint, scalarMultiply(privateKey, aPoint));
  return m.x.toString();
}

// Example usage
const privateKey = generatePrivateKey();
const publicKey = generatePublicKey(privateKey);
console.log("Private Key:", privateKey);
console.log("Public Key:", publicKey);

const plaintext = "Hello, world!";
const ciphertext = encryptMessage(publicKey, plaintext);
console.log("Ciphertext:", ciphertext);

const decryptedMessage = decryptMessage(privateKey, ciphertext);
console.log("Decrypted Message:", decryptedMessage);
