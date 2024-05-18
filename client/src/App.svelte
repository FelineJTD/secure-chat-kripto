<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import ChatHeader from "./lib/ChatHeader.svelte";
  import ChatInput from "./lib/ChatInput.svelte";
  import Container from "./lib/Container.svelte";
  import { onMount } from "svelte";

  type Message = {
    sender: number
    message: string
  }

  let messages: Message[] = []

  const id = Math.floor(Math.random() * 1000000)

  let socket: WebSocket
  let isConnected = false

  const connectWS = () => {
    console.log("Connecting to WS")
    socket = new WebSocket("ws://localhost:8080/ws")
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
    <ChatHeader sender="Conan" isConnected={isConnected} />
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
