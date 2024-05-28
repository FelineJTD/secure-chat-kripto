<script lang="ts">
  // props
  export let isSelf: boolean = false
  export let message: string
  export let isVerified: boolean = false
  export let isGroup: boolean = false
</script>

<div class={`${isSelf ? "flex-row" : "flex-row-reverse"} flex w-full items-end`}>
  <div class="flex-grow" />
  {#if isVerified}
    <img src="../src/assets/icons/seal-check.svg" alt="message verified" class="w-4 h-4 mr-2 opacity-50" />
  {/if}
  <div class={`${isSelf ? "bg-primary rounded-l-xl rounded-tr-xl self-end" : "bg-secondary rounded-r-xl rounded-tl-xl self-start"} ${isGroup ? "mt-0" : "mt-2"} text-neutral-800 py-2 px-4 whitespace-pre-wrap max-w-[36rem]`}>
    {message}
  </div>
</div>
