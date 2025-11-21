<script setup lang="ts">
import ColorPicker from '@/components/ColorPicker.vue'
import PixelGrid from '@/components/PixelGrid.vue'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref<any>(null)

onMounted(() => {
  const userInfo = localStorage.getItem('discord_user')
  if (userInfo) {
    user.value = JSON.parse(userInfo)
  }
})

const handleLogout = () => {
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

      <div class="main-layout">
        <aside class="sidebar">
          <ColorPicker />
        </aside>

        <main class="main-content">
          <PixelGrid />
        </main>
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

.main-layout {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 1.5rem;
  flex: 1;
  min-height: 0;
  animation: fade-in 0.6s ease-out 0.2s both;
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

.sidebar {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1.5rem;
  padding: 2rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow-y: auto;
  max-height: calc(100vh - 180px);
}

.main-content {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1.5rem;
  padding: 2rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  min-height: 0;
}

@media (max-width: 1024px) {
  .main-layout {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
  }

  .sidebar {
    max-height: 400px;
  }

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
}

.sidebar::-webkit-scrollbar {
  width: 8px;
}

.sidebar::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

.sidebar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

.sidebar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
