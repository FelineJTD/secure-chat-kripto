<script lang="ts">
  export let onSend: (message: string) => void

  let input: HTMLInputElement

  const sendMessage = (e: Event) => {
    const formData = new FormData(e.target as HTMLFormElement)
    const message = formData.get("message") as string
    onSend(message)
    input.value = ""
  }
</script>

<form class="flex sticky bottom-0 px-4 md:px-8 py-4 bg-neutral-200" on:submit|preventDefault={sendMessage}>
  <input bind:this={input} name="message" type="text" class="w-full border border-neutral-400 rounded-xl p-2 mr-4" />
  <button class="bg-neutral-600 text-neutral-100 rounded-xl py-2 px-4 flex items-center hover:bg-neutral-700 duration-100">
    <p class="mr-2 mb-0.5">Send</p>
    <img src="../src/assets/icons/paper-plane-right.svg" alt="" class="w-4 h-4 mr-2" />
  </button>
</form>