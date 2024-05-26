<script lang="ts">
  export let onGenerate
</script>

<div class="w-0 lg:w-1/3 pr-24 flex flex-col gap-2 py-4 pl-4">
  <h1 class="text-3xl font-semibold mb-4 text-neutral-600">Keys</h1>
  <label for="private-key" class="text-neutral-500">Private Key</label>
  <input id="private-key" type="text" placeholder="Enter your private key" />
  <label for="public-key" class="text-neutral-500">Public Key</label>
  <input id="public-key" type="text" placeholder="Enter your partner's public key" />

  <h2 class="text-xl font-semibold mb-1 text-neutral-500">No keys yet?</h2>
  <button class="bg-primary hover:bg-secondary duration-100 w-fit text-neutral-800 py-2 px-4 rounded-md" on:click={onGenerate}>Generate keys</button>
</div>

<style>
  input {
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 0.25rem;
    display: block;
    width: 100%;
    margin-bottom: 1rem;
  }
</style>