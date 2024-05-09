<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import ChatHeader from "./lib/ChatHeader.svelte";
  import ChatInput from "./lib/ChatInput.svelte";
  import Container from "./lib/Container.svelte";
  import { onMount } from "svelte";

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
    socket.send(message)
  }
  
</script>

<main class="bg-neutral-100">
  <div class="flex flex-col">
    <Container>
      <ChatHeader sender="Conan" isConnected={isConnected} />
      <ChatContainer>
        <ChatBubble message={`Hello, \nWorld!`} />
        <ChatBubble isSelf message="Hello, World!" />
        <ChatBubble message={`Hello, \nWorld!`} />
        <ChatBubble isSelf message="Hello, World!" />
        <ChatBubble message={`Hello, \nWorld!`} />
        <ChatBubble isSelf message="Hello, World!" />
        <ChatBubble isSelf message="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." />
      </ChatContainer>
      <ChatInput onSend={onSend} />
    </Container>
  </div>
</main>

<style>

</style>
