<script lang="ts">
  export let doSign
  export let setSignKey
  export let setVerifyKey
  export let onGenerateSign
  export let onGenerate
  export let setPrivKeyECC
  export let setPubKeyECC
  export let status
</script>

<div class="w-0 lg:w-1/3 pr-24 flex flex-col gap-2 py-4 pl-4">
  <h1 class="text-3xl font-semibold mb-4 text-neutral-600">Keys</h1>
  <label for="private-key" class="text-neutral-500">Enter your private key</label>
  <input id="private-key" type="file" on:change={(e) => setPrivKeyECC(e)} placeholder="Enter your private key" />
  <label for="public-key" class="text-neutral-500">Enter your partner's public key</label>
  <input id="public-key" type="file" on:change={(e) => setPubKeyECC(e)} placeholder="Enter your partner's public key" />

  <h2 class="text-xl font-semibold mb-1 text-neutral-500">No keys yet?</h2>
  <button class="bg-primary hover:bg-secondary duration-100 w-fit text-neutral-800 py-2 px-4 rounded-md" on:click={onGenerate}>Generate keys</button>
  <!-- TODO: Change to Toast Instead -->
  {#if status}
    <p class="text-green
    ">{status}</p>
  {/if}
  
  <label for="public-key" class="text-neutral-500">Enter your signing key</label>
  <input id="public-key" type="file" on:change={(e) => setSignKey(e)} placeholder="Enter your signing private key" />
  <label for="public-key" class="text-neutral-500">Enter your verifying key</label>
  <input id="public-key" type="file" on:change={(e) => setVerifyKey(e)} placeholder="Enter your partner's public key" />
  <div class="flex flex-row items-center">
    <h2 class="text-xl font-semibold mb-1 text-neutral-500">Sign Messages?</h2>
    <input id="digisign" type="checkbox" on:change={(e) => doSign(e)} />
  </div>
  
  <h2 class="text-xl font-semibold mb-1 text-neutral-500">No Signing keys yet?</h2>
  <button class="bg-primary hover:bg-secondary duration-100 w-fit text-neutral-800 py-2 px-4 rounded-md" on:click={onGenerateSign}>Generate keys</button>
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

  input[type="checkbox"] {
    width: 1rem;
    height: 1rem;
    margin: 0;
    margin-bottom: 0;
    margin-left: 1rem;
  }
</style>