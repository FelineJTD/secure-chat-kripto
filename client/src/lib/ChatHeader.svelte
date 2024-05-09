<script lang="ts">
  // props
  export let sender: string
  export let isConnected: boolean
</script>

<div class="flex px-4 md:px-8 py-2 bg-neutral-600 z-10 text-neutral-100 h-16 items-center justify-between">
  <!-- name and profile picture -->
  <div class="flex h-full">
    <div class="rounded-full min-h-full aspect-square bg-neutral-400 mr-4 overflow-clip">
      <img src="https://picsum.photos/200" alt="Profile" class="rounded-full" />
    </div>
    <div class="-mt-0.5">
      <h1 class="font-semibold text-2xl">{sender}</h1>
      <p class="text-neutral-300 text-sm">Chat is secure.</p>
    </div>
  </div>
  <!-- connection status -->
  <div class={`${isConnected ? 'bg-green' : 'bg-red'} size-4 rounded-full`} />
</div>