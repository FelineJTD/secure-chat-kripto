<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
</script>

<main class="bg-neutral-100">
  <div class="max-w-[48rem] min-h-screen m-auto flex flex-col p-12 gap-4">
    <ChatBubble message={`Hello, \nWorld!`} />
    <ChatBubble isSelf message="Hello, World!" />
    <ChatBubble isSelf message="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." />
  </div>
</main>

<style>

</style>
