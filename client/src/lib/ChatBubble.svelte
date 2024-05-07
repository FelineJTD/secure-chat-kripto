<script lang="ts">
  // props
  export let isSelf: boolean = false
  export let message: string
</script>

<div class={`${isSelf ? "bg-primary rounded-l-xl rounded-tr-xl self-end" : "bg-secondary rounded-r-xl rounded-tl-xl self-start"} text-neutral-800 py-2 px-4 whitespace-pre-wrap max-w-[36rem]`}>
  {message}
</div>
