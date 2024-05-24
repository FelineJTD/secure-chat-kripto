<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import ChatHeader from "./lib/ChatHeader.svelte";
  import ChatInput from "./lib/ChatInput.svelte";
  import Container from "./lib/Container.svelte";
  import { onMount } from "svelte";
  // import { Point, generatePrivateKey, generatePublicKey } from "./utils/ecc";

  type Message = {
    sender: string
    message: string
  }

  let messages: Message[] = []
  let id: string
  // let privKey: bigint | null
  // let pubKey: Point | null

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

  onMount(() => {
    // // Try to get private key from local storage
    // privKey = localStorage.getItem("privKey") ? BigInt(localStorage.getItem("privKey") as string) : null
    // pubKey = localStorage.getItem("pubKey") ? JSON.parse(localStorage.getItem("pubKey") as string) : null
    // if (!privKey || !pubKey) {
    //   // If private key is not found, generate a new one
    //   privKey = generatePrivateKey()
    //   pubKey = generatePublicKey(privKey)
    //   localStorage.setItem("privKey", privKey.toString())
    //   localStorage.setItem("pubKey", pubKey.toString())
    // } 

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
