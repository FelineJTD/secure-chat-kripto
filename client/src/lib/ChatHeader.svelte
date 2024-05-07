<script lang="ts">
  // props
  export let sender: string
</script>

<div class="flex sticky top-0 px-4 md:px-8 py-2 bg-neutral-600 z-10 text-neutral-100 h-16 items-center">
  <div class="rounded-full min-h-full aspect-square bg-neutral-400 mr-4 overflow-clip">
    <!-- INSERT IMG HERE -->
  </div>
  <div class="-mt-0.5">
    <h1 class="font-semibold text-2xl">{sender}</h1>
    <p class="text-neutral-300 text-sm">Chat is secure.</p>
  </div>
</div>