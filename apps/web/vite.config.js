import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    svelte(),
  ],
  resolve: {
    alias: {
      "@exorcist-dto": path.resolve(__dirname, "./src/dto"),
      "@exorcist-type": path.resolve(__dirname, "./src/lib/types")
    }
  }
})
