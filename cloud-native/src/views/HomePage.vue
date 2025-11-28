<script setup lang="ts">
import PixelGrid from '@/components/PixelGrid.vue'
import { usePixelSync } from '@/composables/usePixelSync'
import { loadCanvas } from '@/services/gatewayApi'
import { usePixelStore } from '@/stores/pixelStore'
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref<any>(null)
const pixelStore = usePixelStore()
const { isConnected, connect, disconnect } = usePixelSync()
const isLoading = ref(true)

onMounted(async () => {
  const userInfo = localStorage.getItem('discord_user')
  if (userInfo) {
    user.value = JSON.parse(userInfo)
  }

  try {
    const pixels = await loadCanvas()
    pixels.forEach(({ x, y, color }) => {
      pixelStore.syncPixel(x, y, color)
    })

    connect()
  } catch (error) {
    console.error("Erreur lors de l'initialisation:", error)
  } finally {
    isLoading.value = false
  }
})

onBeforeUnmount(() => {
  disconnect()
})

const handleLogout = () => {
  disconnect()
  localStorage.removeItem('discord_token')
  localStorage.removeItem('discord_user')
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="app-container">
    <div class="bg-gradient"></div>
    <div class="bg-pattern"></div>

    <div class="content-wrapper">
      <header class="app-header">
        <div class="header-content">
          <div class="flex items-center gap-4">
            <div class="logo-container">
              <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M4 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1v-4zM14 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z"
                />
              </svg>
            </div>
            <div>
              <h1 class="app-title">Pixel Place</h1>
            </div>

            <div
              v-if="!isLoading"
              :class="['connection-indicator', isConnected ? 'connected' : 'disconnected']"
            >
              <span class="status-dot"></span>
              <span class="status-text">{{ isConnected ? 'En ligne' : 'Hors ligne' }}</span>
            </div>
          </div>
          <div class="header-right">
            <div v-if="user" class="user-info">
              <img
                v-if="user.avatar"
                :src="`https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png`"
                :alt="user.username"
                class="user-avatar"
              />
              <div v-else class="user-avatar-default">
                {{ user.username?.charAt(0).toUpperCase() }}
              </div>
              <span class="user-name">{{ user.username }}</span>
            </div>
            <button @click="handleLogout" class="logout-button">
              <svg class="logout-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                />
              </svg>
              <span>Déconnexion</span>
            </button>
          </div>
        </div>
      </header>

      <div class="grid-container">
        <PixelGrid />
      </div>
    </div>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
}

.bg-gradient {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background:
    radial-gradient(circle at 20% 50%, rgba(99, 102, 241, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(236, 72, 153, 0.1) 0%, transparent 50%);
  pointer-events: none;
}

.bg-pattern {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.02) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.02) 1px, transparent 1px);
  background-size: 50px 50px;
  pointer-events: none;
}

.content-wrapper {
  position: relative;
  z-index: 1;
  max-width: 1600px;
  margin: 0 auto;
  padding: 2rem;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  margin-bottom: 2rem;
  animation: slide-down 0.5s ease-out;
}

@keyframes slide-down {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.flex {
  display: flex;
}

.items-center {
  align-items: center;
}

.gap-4 {
  gap: 1rem;
}

.logo-container {
  width: 3rem;
  height: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #6366f1 0%, #ec4899 100%);
  border-radius: 1rem;
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.3);
  color: white;
  animation: pulse-glow 2s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%,
  100% {
    box-shadow: 0 8px 24px rgba(99, 102, 241, 0.3);
  }
  50% {
    box-shadow: 0 8px 32px rgba(99, 102, 241, 0.5);
  }
}

.w-12 {
  width: 3rem;
}

.h-12 {
  height: 3rem;
}

.app-title {
  font-size: 2.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, #ffffff 0%, #e2e8f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.02em;
}

.app-subtitle {
  font-size: 1rem;
  color: #94a3b8;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 9999px;
  backdrop-filter: blur(10px);
}

.user-avatar {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  border: 2px solid rgba(99, 102, 241, 0.5);
}

.user-avatar-default {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1 0%, #ec4899 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
}

.user-name {
  color: #e2e8f0;
  font-weight: 500;
  font-size: 0.875rem;
}

.logout-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 0.5rem;
  color: #fca5a5;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
}

.logout-button:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.5);
  transform: translateY(-1px);
}

.logout-icon {
  width: 1rem;
  height: 1rem;
}

.grid-container {
  flex: 1;
  min-height: 0;
  animation: fade-in 0.6s ease-out 0.2s both;
  display: flex;
  flex-direction: column;
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.connection-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.3s ease;
}

.connection-indicator.connected {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: #86efac;
}

.connection-indicator.disconnected {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #fca5a5;
}

.status-dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 50%;
  animation: pulse-dot 2s ease-in-out infinite;
}

.connected .status-dot {
  background: #22c55e;
}

.disconnected .status-dot {
  background: #ef4444;
}

@keyframes pulse-dot {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(1.2);
  }
}

.status-text {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

@media (max-width: 1024px) {
  .app-title {
    font-size: 2rem;
  }

  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-right {
    width: 100%;
    justify-content: space-between;
  }

  .connection-indicator {
    font-size: 0.75rem;
    padding: 0.375rem 0.75rem;
  }

  .status-text {
    font-size: 0.625rem;
  }
}
</style>
