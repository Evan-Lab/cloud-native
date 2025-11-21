import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import { defineConfig } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueJsx(), vueDevTools(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      // Proxy pour contourner CORS de Discord (optionnel)
      '/api/discord': {
        target: 'https://discord.com/api',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/discord/, ''),
        secure: true,
        headers: {
          'User-Agent': 'PixelPlace/1.0',
        },
      },
    },
  },
})
