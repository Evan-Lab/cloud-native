import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import { defineConfig } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig({
  plugins: [vue(), vueJsx(), vueDevTools(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      '/api/discord': {
        target: 'https://discord.com/api',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/discord/, ''),
        secure: true,
        headers: {
          'User-Agent': 'PixelPlace/1.0',
        },
      },
      '/api/gateway': {
        target: 'https://rplace-gateway-5uir24en.ew.gateway.dev',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/gateway/, '/web/api'),
        secure: true,
      },
      '/ws': {
        target: 'wss://rplace-gateway-5uir24en.ew.gateway.dev',
        ws: true,
        changeOrigin: true,
      },
    },
  },
})
