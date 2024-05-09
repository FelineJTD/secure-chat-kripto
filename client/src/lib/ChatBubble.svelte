<script lang="ts">
  // props
  export let isSelf: boolean = false
  export let message: string
  export let isGroup: boolean = false
</script>

<div class={`${isSelf ? "bg-primary rounded-l-xl rounded-tr-xl self-end" : "bg-secondary rounded-r-xl rounded-tl-xl self-start"} ${isGroup ? "mt-0" : "mt-2"} text-neutral-800 py-2 px-4 whitespace-pre-wrap max-w-[36rem]`}>
  {message}
</div>
