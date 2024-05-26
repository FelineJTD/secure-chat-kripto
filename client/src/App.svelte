<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import ChatHeader from "./lib/ChatHeader.svelte";
  import ChatInput from "./lib/ChatInput.svelte";
  import Container from "./lib/Container.svelte";
  import KeyInputs from "./lib/KeyInputs.svelte";
  import { onMount } from "svelte";
  import { decryptMessage, encryptMessage, generateKeyPair, Point } from "./utils/ecc";

  type Message = {
    sender: string
    message: string
  }

  let messages: Message[] = []
  let id: string

  let privKeyECC: bigint | null
  let pubKeyECC: Point | null

  let socket: WebSocket
  let isConnected = false

  // Connect to WebSocket server
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
        message: decryptMessage(privKeyECC as bigint, JSONToPoints(payload.message))
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

  // Convert points to JSON string
  function pointsToJSON(point: [Point, Point]): string {
    return JSON.stringify([{x: point[0].x.toString(), y: point[0].y.toString()}, {x: point[1].x.toString(), y: point[1].y.toString()}])
  }

  // Convert JSON string to points
  function JSONToPoints(json: string): [Point, Point] {
    const points = JSON.parse(json)
    return [new Point(BigInt(points[0].x), BigInt(points[0].y)), new Point(BigInt(points[1].x), BigInt(points[1].y))]
  }

  // Generate key pairs
  function generate() {
    const [priv, pub] = generateKeyPair()
    privKeyECC = priv
    pubKeyECC = pub
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

    const url = window.location.href
    id = url.split(":")[2].split("/")[0]
    connectWS()

    return () => {
      socket.close()
    }
  })

  const onSend = (message: string) => {
    if (!pubKeyECC) return

    const payload = {
      sender: id,
      message: pointsToJSON(encryptMessage(pubKeyECC, message))
    }
    const payloadString = JSON.stringify(payload)
    console.log("Sending ", payloadString)
    socket.send(payloadString)
  }
</script>

<main class="bg-neutral-100 h-screen">
  <div class="flex flex-col lg:flex-row w-full h-screen">
    <KeyInputs onGenerate={generate} />
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
