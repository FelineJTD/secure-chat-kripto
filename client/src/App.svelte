<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import ChatHeader from "./lib/ChatHeader.svelte";
  import ChatInput from "./lib/ChatInput.svelte";
  import Container from "./lib/Container.svelte";
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
      messages = [payload, ...messages]
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

  const p = 76511n;
  const a = 71479n;
  const b = 52157n;
  const Gx = 14703n;
  const Gy = 38268n;
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

  isInfinite(): boolean {
    return this.x === 0n && this.y === 0n;
  }

  convertToMinus(): Point {
    return new Point(this.x, -this.y);
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
    return new Point(0n, 0n);
  } else if (p1.isInfinite()) {
    return p2;
  } else if (p2.isInfinite()) {
    return p1;
  } else {
    // Calculate slope
    const slope = ((((p2.y - p1.y) * modInverse(p2.x - p1.x, p)) % p) + p) % p;
    console.log("slope", slope)
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

function isOdd(n: bigint)
{
 
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
      result = pointAddition(result, addend);
    }
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
export function generatePrivateKey(): bigint {
  const hexString = Array(4)
    .fill(0)
    .map(() => Math.round(Math.random() * 0xf).toString(16))
    .join("");

  const randomInt = BigInt(`0x${hexString}`);
  return randomInt;
}

// Generate the corresponding public key
export function generatePublicKey(privateKey: bigint): Point {
  return scalarMultiply(privateKey, new Point(Gx, Gy));
}

// Encrypt a message using ECC
export function encryptMessage(publicKey: Point, message: string): [Point, Point] {
  const dummy = new Point(24601n, 33894n);
  const k = generatePrivateKey();
  // a = g^k mod p
  const a = scalarMultiply(k, new Point(Gx, Gy));
  // b = m*P^k mod p
  // const m = BigInt(message);
  // console.log("m", dummy)
  const b = pointAddition(scalarMultiply(k, publicKey), dummy);
  console.log("a", a)
  console.log("b", b)

  return [a, b]
}

// Decrypt a message using ECC
export function decryptMessage(privateKey: bigint, ciphertext: [Point, Point]): Point {
  // const [a, b, message] = ciphertext.split(",");
  // const aPoint = new Point(BigInt(a), BigInt(b));
  const [a, b] = ciphertext;
  const m = pointAddition(b, scalarMultiply(privateKey, a).convertToMinus());
  console.log("m", m.x)
  console.log("m", m.y)
  return m;

  // const m = pointAddition(aPoint, scalarMultiply(privateKey, aPoint));
  // return m.x.toString();
}

  onMount(() => {
    // // Try to get private key from local storage
    privKey = localStorage.getItem("privKey") ? BigInt(localStorage.getItem("privKey") as string) : null
    pubKey = localStorage.getItem("pubKey") ? JSON.parse(localStorage.getItem("pubKey") as string) : null
    if (!privKey || !pubKey) {
      // If private key is not found, generate a new one
      privKey = generatePrivateKey()
      pubKey = generatePublicKey(privKey)
      localStorage.setItem("privKey", privKey.toString())
      // localStorage.setItem("pubKey", JSON.stringify(pubKey))
    } 
    console.log("privKey", privKey)
    console.log("pubKey", pubKey)

    const message = "Hello World!"
    const encryptedMessage = encryptMessage(pubKey, message)
    console.log("Encrypted message", encryptedMessage)
    const decryptedMessage = decryptMessage(privKey, encryptedMessage)
    console.log("Decrypted message", decryptedMessage)

    const url = window.location.href
    id = url.split(":")[2].split("/")[0]
    connectWS()

    return () => {
      socket.close()
    }
  })

  const onSend = (message: string) => {
    const payload = {
      sender: id,
      message
    }
    const payloadString = JSON.stringify(payload)
    socket.send(payloadString)
  }
  
</script>

<main class="bg-neutral-100 h-screen">
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
</main>

<style>

</style>
