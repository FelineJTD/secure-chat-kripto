import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import goWasm from 'vite-plugin-golang-wasm'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    svelte(),
    goWasm({
      goBinaryPath: 'C:\\Program Files\\Go\\bin\\go.exe',
      wasmExecPath: 'C:\\Program Files\\Go\\misc\\wasm\\wasm_exec.js',
      goBuildExtraArgs: ['-C', './src/wasm'],
    }),
  ],
})
