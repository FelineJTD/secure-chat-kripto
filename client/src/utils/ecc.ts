// ecc.ts

// Parameters for the curve
const p =
  9391992014528224411648387829887779452968700080555255747361254716652003361688800380330385633056245683444877536960503368327139028742945338465232435742415490790285260907481871880469373580741241260948876883488819453317735629504379161346245998683504143286421351n;
const a =
  9233124640038346995100978392313313511369492123558996708844011838096113732285757053789463384591671694743340670303624555512347783309238918503776319156266971458062018710184003825351802199040926442114671941913985501450211259391979794403797469337351719270652351n;
const Gx =
  2398249831791816336028631263148105895130194480594232706681452236075099255914833333176431307231368483589360217800747308174513483544427244190177932600264268134572314832192520846016006636941990751131120801543171652309984054140576933467418364998622752106966361n;
const Gy =
  2937610345392163240898320965449658011204813209985492263862609624172037174893838813802778590568695212500828225410673923612008661914350991091209447994751752919108072103550418611957669073010110009093661014607538825852066098351822205667186845598987024537485599n;

// Point class
export class Point {
  x: bigint;
  y: bigint;

  constructor(x: bigint, y: bigint) {
    this.x = x;
    this.y = y;
  }

  equals(other: Point): boolean {
    return this.x === other.x && this.y === other.y;
  }

  sameXDiffY(other: Point): boolean {
    return this.x === other.x && this.y !== other.y;
  }

  isInfinite(): boolean {
    return this.x === 0n && this.y === 0n;
  }

  convertToMinus(): Point {
    return new Point(this.x, -this.y);
  }
}

function stringToHex(s: string): string {
  return s
    .split("")
    .map((c) => c.charCodeAt(0).toString(16))
    .join("");
}

function hexToBigInt(h: string): bigint {
  return BigInt(`0x${h}`);
}

function bigIntToHex(n: bigint): string {
  return n.toString(16);
}

function hexToString(h: string): string {
  return h
    .match(/.{1,2}/g)!
    .map((byte) => String.fromCharCode(parseInt(byte, 16)))
    .join("");
}

function modInverse(aPoint: bigint, bPoint: bigint): bigint {
  let a = aPoint;
  let b = bPoint;
  let x = 0n;
  let y = 1n;
  let lastx = 1n;
  let lasty = 0n;
  let temp = 0n;

  while (b !== 0n) {
    const q = a / b;
    const r = a % b;
    a = b;
    b = r;
    temp = x;
    x = lastx - q * x;
    lastx = temp;
    temp = y;
    y = lasty - q * y;
    lasty = temp;
  }

  return lastx;
}

function pointAddition(p1: Point, p2: Point): Point {
  if (p1.equals(p2)) {
    return pointDoubling(p1);
  } else if (p1.sameXDiffY(p2)) {
    return new Point(0n, 0n);
  } else if (p1.isInfinite()) {
    return p2;
  } else if (p2.isInfinite()) {
    return p1;
  } else {
    let dX = p2.x - p1.x;
    let dY = p2.y - p1.y;
    if (dX < 0n) {
      dX = -dX;
      dY = -dY;
    }
    const slope = (((dY * modInverse(dX, p)) % p) + p) % p;
    // Calculate x-coordinate
    const x = (((slope * slope - p1.x - p2.x) % p) + p) % p;
    // Calculate y-coordinate
    const y = (((slope * (p1.x - x) - p1.y) % p) + p) % p;
    return new Point(x, y);
  }
}

// Point doubling on the curve
function pointDoubling(p1: Point): Point {
  // Calculate slope
  const slope =
    ((((3n * p1.x * p1.x + a) * modInverse(2n * p1.y, p)) % p) + p) % p;
  // Calculate x-coordinate
  const x = (((slope * slope - 2n * p1.x) % p) + p) % p;
  // Calculate y-coordinate
  const y = (((slope * (p1.x - x) - p1.y) % p) + p) % p;
  return new Point(x, y);
}

function isOdd(n: bigint) {
  // n^1 is n+1, then even, else odd
  if ((n ^ 1n) == n + 1n) return false;
  else return true;
}

// Scalar multiplication on the curve
export function scalarMultiply(k: bigint, p1: Point): Point {
  let result = new Point(0n, 0n);
  let addend = p1;
  while (k > 0n) {
    if (isOdd(k)) {
      result = pointAddition(result, addend);
    }
    addend = pointDoubling(addend);
    k = k >> 1n;
  }
  return result;
}

function generatePrivateKey(): bigint {
  const hexString = Array(250)
    .fill(0)
    .map(() => Math.round(Math.random() * 0xf).toString(16))
    .join("");

  const randomInt = BigInt(`0x${hexString}`);
  return randomInt;
}

export function generatePublicKey(privateKey: bigint): Point {
  return scalarMultiply(privateKey, new Point(Gx, Gy));
}

export function generateKeyPair(): [bigint, Point] {
  const privateKey = generatePrivateKey();
  const publicKey = generatePublicKey(privateKey);
  return [privateKey, publicKey];
}

export function deriveSharedSecret(
  privateKey: bigint,
  publicKey: Point
): Point {
  return scalarMultiply(privateKey, publicKey);
}

function encryptECC(publicKey: Point, message: Point): [Point, Point] {
  const k = generatePrivateKey();
  const a = scalarMultiply(k, new Point(Gx, Gy));
  const b = pointAddition(scalarMultiply(k, publicKey), message);
  return [a, b];
}

function decryptECC(privateKey: bigint, ciphertext: [Point, Point]): Point {
  const [a, b] = ciphertext;
  const m = pointAddition(b, scalarMultiply(privateKey, a).convertToMinus());
  return m;
}

export function encryptMessage(
  publicKey: Point,
  message: string
): [Point, Point] {
  const msgInt = hexToBigInt(stringToHex(message));
  const msgRemainder = msgInt % p;
  const msgMultiple = msgInt / p;
  const msgPoint = new Point(msgRemainder, msgMultiple);
  console.log("Message", msgPoint);
  return encryptECC(publicKey, msgPoint);
}

export function decryptMessage(
  privateKey: bigint,
  ciphertext: [Point, Point]
): string {
  const decryptedMessage = decryptECC(privateKey, ciphertext);
  console.log("Decrypted message", decryptedMessage);
  return hexToString(bigIntToHex(decryptedMessage.x + decryptedMessage.y * p));
}
