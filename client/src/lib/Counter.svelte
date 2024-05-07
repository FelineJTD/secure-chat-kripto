<script lang="ts">
  let count: number = 0
  const increment = () => {
    count += 1
  }
</script>

<button on:click={increment} class="text-3xl text-red-700">
  count is {count}
</button>
