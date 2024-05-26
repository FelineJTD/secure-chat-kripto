<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import ChatHeader from "./lib/ChatHeader.svelte";
  import ChatInput from "./lib/ChatInput.svelte";
  import Container from "./lib/Container.svelte";
  import KeyInputs from "./lib/KeyInputs.svelte";
  import { onMount } from "svelte";

  type Message = {
    sender: string
    message: string
  }

  let messages: Message[] = []
  let id: string
  let privKey: bigint | null
  let pubKey: Point | null

  let socket: WebSocket
  let isConnected = false

  const connectWS = () => {
    socket = new WebSocket("ws://localhost:8080/chat")
    socket.addEventListener("open", ()=> {
      console.log("Opened")
      isConnected = true
    })
    socket.addEventListener("message", (event) => {
      console.log("Message from server ", event.data)
      const payload = JSON.parse(event.data)
      const message = {
        sender: payload.sender,
        message: decryptMessage(privKey as bigint, JSONToPoints(payload.message))
      }
      messages = [message, ...messages]
    })
    socket.addEventListener("close", () => {
      console.log("Closed")
      isConnected = false
      // While connection is closed, try to reconnect every 3 seconds
      setTimeout(() => {
        connectWS()
      }, 3000)
    })
    socket.addEventListener("error", (error) => {
      console.log("Error ", error)
    })
  }

  const p = 49393n;
  const a = 71479n;
  const b = 52157n;
  const Gx = 14703n;
  const Gy = 38268n;

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

  isInfinite(): boolean {
    return this.x === 0n && this.y === 0n;
  }

  convertToMinus(): Point {
    return new Point(this.x, -this.y);
  }
}

function stringToHex(s: string): string {
  return s.split("").map((c) => c.charCodeAt(0).toString(16)).join("");
}

function hexToBigInt(h: string): bigint {
  return BigInt(`0x${h}`);
}

function bigIntToHex(n: bigint): string {
  return n.toString(16);
}

function hexToString(h: string): string {
  return h.match(/.{1,2}/g)!.map((byte) => String.fromCharCode(parseInt(byte, 16))).join("");
}

function calculateY(x: bigint): bigint {
  return (((x * x * x + a * x + b) % p) + p) % p;
}

// Input: a, b
// Output: multiplicative inverse of a
// 1 x = 0 y=1 lastx = 1 lasty = 0;
// 2 temp = 0 temp2 = a temp1 = b;
// 3 while b != 0 do
// 4 q = a/b;
// 5 r = a%b;
// 6 a = b;
// 7 b = r;
// 8 temp = x;
// 9 x = lastx - q * x;
// 10 lastx = temp;
// 11 temp = y;
// 12 y = lasty - q*y;
// 13 lasty = temp;
// end
// return lastx

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
// Algorithm 10: The pseudocode of adding two points on a curve
// Input: p1, p2
// Output: result
// 1 if p1 = p2 then
// 2 return double p1;
// end
// 3 dY = p2.y - p1.y;
// 4 dX = p2.x -p1.x;
// 5 if dX is negative then
// 6 flip signs of dX and xY;
// end
// 7 dX = gcdExtended(dX, p);
// 8 slope = dY*dX%p;
// 9 result.x = (slope.pow(2, p) - p1.x - p2.x) % p;
// 10 result.y = slope * (p1.x - result.x) - p1.y;
// 11 return result;
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
    const slope = ((((p2.y - p1.y) * modInverse(p2.x - p1.x, p)) % p) + p) % p;
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
  const slope = ((((3n * p1.x * p1.x + a) * modInverse(2n * p1.y, p)) % p) + p) %p;
  // Calculate x-coordinate
  const x = (((slope * slope - 2n * p1.x) % p) + p) % p;
  // Calculate y-coordinate
  const y = (((slope * (p1.x - x) - p1.y) % p) + p) % p;
  return new Point(x, y);
}

function isOdd(n: bigint) {
  // n^1 is n+1, then even, else odd
  if ((n ^ 1n) == (n + 1n))
    return false;
  else
    return true;
}

// Scalar multiplication on the curve
function scalarMultiply(k: bigint, p1: Point): Point {
  let result = new Point(0n, 0n);
  let addend = p1;
  while (k > 0n) {
    if (isOdd(k)) {
      console.log("point addition")
      result = pointAddition(result, addend);
    }
    console.log("point doubling")
    addend = pointDoubling(addend);
    k = k >> 1n;
  }
  return result;
}

// function scalarMultiply(k: BigInteger, P: {x: BigInteger, y: BigInteger}): {x: BigInteger, y: BigInteger} {
//     let result = { x: BigInteger.zero, y: BigInteger.zero }; // Initialize the result to the point at infinity
//     let addend = P;

//     while (!k.isZero()) {
//         if (k.isOdd()) {
//             result = pointAdd(result, addend);
//         }
//         addend = pointAdd(addend, addend);
//         k = k.shiftRight(1); // Equivalent to dividing k by 2
//     }

//     return result;
// }

// Generate a random private key
function generatePrivateKey(): bigint {
  const hexString = Array(8)
    .fill(0)
    .map(() => Math.round(Math.random() * 0xf).toString(16))
    .join("");

  const randomInt = BigInt(`0x${hexString}`);
  return randomInt;
}

// Generate the corresponding public key
function generatePublicKey(privateKey: bigint): Point {
  return scalarMultiply(privateKey, new Point(Gx, Gy));
}

// Encrypt a message using ECC
function encryptECC(publicKey: Point, message: Point): [Point, Point] {
  // const dummy = new Point(24601n, 33894n);
  const k = generatePrivateKey();
  // a = g^k mod p
  const a = scalarMultiply(k, new Point(Gx, Gy));
  // b = m*P^k mod p
  // const m = BigInt(message);
  // console.log("m", dummy)
  const b = pointAddition(scalarMultiply(k, publicKey), message);
  console.log("a", a)
  console.log("b", b)

  return [a, b]
}

// Decrypt a message using ECC
function decryptECC(privateKey: bigint, ciphertext: [Point, Point]): Point {
  // const [a, b, message] = ciphertext.split(",");
  // const aPoint = new Point(BigInt(a), BigInt(b));
  const [a, b] = ciphertext;
  const m = pointAddition(b, scalarMultiply(privateKey, a).convertToMinus());
  return m;

  // const m = pointAddition(aPoint, scalarMultiply(privateKey, aPoint));
  // return m.x.toString();
}

function encryptMessage(publicKey: Point, message: string): [Point, Point] {
  const msgInt = hexToBigInt(stringToHex(message))
  const msgRemainder = msgInt % p
  const msgMultiple = msgInt /p
  // console.log(msgInt === msgRemainder + (msgMultiple * p))
  const msgPoint = new Point(msgRemainder, msgMultiple)
  console.log("Message", msgPoint)
  return encryptECC(publicKey, msgPoint)
}

function decryptMessage(privateKey: bigint, ciphertext: [Point, Point]): string {
  const decryptedMessage = decryptECC(privateKey, ciphertext)
  console.log("Decrypted message", decryptedMessage)
  return hexToString(bigIntToHex((decryptedMessage.x) + decryptedMessage.y * p))
}

function pointsToJSON(point: [Point, Point]): string {
  return JSON.stringify([{x: point[0].x.toString(), y: point[0].y.toString()}, {x: point[1].x.toString(), y: point[1].y.toString()}])
}

function JSONToPoints(json: string): [Point, Point] {
  const points = JSON.parse(json)
  return [new Point(BigInt(points[0].x), BigInt(points[0].y)), new Point(BigInt(points[1].x), BigInt(points[1].y))]
}

function generateKeyPairs() {
  privKey = generatePrivateKey()
  pubKey = generatePublicKey(privKey)
  console.log("privKey", privKey)
  console.log("pubKey", pubKey)
}


  onMount(() => {
    // // Try to get private key from local storage
    // privKey = localStorage.getItem("privKey") ? BigInt(localStorage.getItem("privKey") as string) : null
    // pubKey = localStorage.getItem("pubKey") ? JSON.parse(localStorage.getItem("pubKey") as string) : null
    // if (!privKey || !pubKey) {
    //   // If private key is not found, generate a new one
    //   privKey = generatePrivateKey()
    //   pubKey = generatePublicKey(privKey)
    //   localStorage.setItem("privKey", privKey.toString())
    //   // localStorage.setItem("pubKey", JSON.stringify(pubKey))
    // } 
    // console.log("privKey", privKey)
    // console.log("pubKey", pubKey)

    // const message = "   a"
    // const msgInt = hexToBigInt(stringToHex(message))
    // console.log("Message", msgInt)
    // const back = hexToString(bigIntToHex(msgInt))
    // console.log("Back", back)
    // const encryptedMessage = encryptMessage(pubKey, message)
    // console.log("Encrypted message", encryptedMessage)
    // const decryptedMessage = decryptMessage(privKey, encryptedMessage)
    // console.log("Decrypted message", decryptedMessage)

    const url = window.location.href
    id = url.split(":")[2].split("/")[0]
    connectWS()

    return () => {
      socket.close()
    }
  })

  const onSend = (message: string) => {
    if (!pubKey) return

    const payload = {
      sender: id,
      message: pointsToJSON(encryptMessage(pubKey, message))
    }
    const payloadString = JSON.stringify(payload)
    socket.send(payloadString)
  }
  
</script>

<main class="bg-neutral-100 h-screen">
  <div class="flex flex-col lg:flex-row w-full h-screen">
    <KeyInputs onGenerate={generateKeyPairs} />
    <Container>
      <ChatHeader sender={id} isConnected={isConnected} />
      <ChatContainer>
        <!-- loop with index -->
        {#each messages as message, i}
          <ChatBubble isSelf={message.sender === id} message={message.message} />
        {/each}
      </ChatContainer>
      <ChatInput onSend={onSend} />
    </Container>
    <div class="w-0 lg:w-1/3" />
  </div>
</main>

<style>

</style>
