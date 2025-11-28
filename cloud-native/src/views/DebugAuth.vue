<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

const isAuthenticated = ref(false)
const user = ref<any>(null)
const token = ref<string | null>(null)
const tokenExpiry = ref<string | null>(null)
const tokenType = ref<string | null>(null)
const tokenScope = ref<string | null>(null)

const checkAuth = () => {
  token.value = localStorage.getItem('discord_token')
  tokenType.value = localStorage.getItem('discord_token_type')
  tokenScope.value = localStorage.getItem('discord_token_scope')
  const userStr = localStorage.getItem('discord_user')
  const expiry = localStorage.getItem('discord_token_expiry')

  if (userStr) {
    user.value = JSON.parse(userStr)
    isAuthenticated.value = true
  }

  if (expiry) {
    const expiryDate = new Date(parseInt(expiry))
    tokenExpiry.value = expiryDate.toLocaleString('fr-FR')
  }
}

const localStorageData = computed(() => ({
  discord_token: token.value?.substring(0, 30) + '...',
  discord_token_type: tokenType.value,
  discord_token_expiry: tokenExpiry.value,
  discord_token_scope: tokenScope.value,
  discord_user: user.value
}))

const logout = () => {
  localStorage.removeItem('discord_token')
  localStorage.removeItem('discord_token_type')
  localStorage.removeItem('discord_token_expiry')
  localStorage.removeItem('discord_token_scope')
  localStorage.removeItem('discord_user')
  checkAuth()
  window.location.href = '/login'
}

onMounted(() => {
  checkAuth()
  // Écouter les changements d'authentification
  window.addEventListener('discord-auth-success', checkAuth)
})
</script>

<template>
  <div class="debug-container">
    <h1 class="title">🔍 Debug Authentification Discord</h1>
    
    <div class="card">
      <h2 class="section-title">Statut</h2>
      <div class="status-badge" :class="{ authenticated: isAuthenticated }">
        {{ isAuthenticated ? '✅ Authentifié' : '❌ Non authentifié' }}
      </div>
    </div>

    <div v-if="isAuthenticated" class="card">
      <h2 class="section-title">👤 Informations Utilisateur</h2>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">Username:</span>
          <span class="value">{{ user?.username }}</span>
        </div>
        <div class="info-item">
          <span class="label">ID:</span>
          <span class="value">{{ user?.id }}</span>
        </div>
        <div class="info-item">
          <span class="label">Email:</span>
          <span class="value">{{ user?.email || 'N/A' }}</span>
        </div>
        <div class="info-item">
          <span class="label">Global Name:</span>
          <span class="value">{{ user?.global_name || user?.username }}</span>
        </div>
        <div v-if="user?.avatar" class="info-item">
          <span class="label">Avatar:</span>
          <img 
            :src="`https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png`"
            alt="Avatar"
            class="avatar"
          />
        </div>
      </div>
    </div>

    <div v-if="isAuthenticated" class="card">
      <h2 class="section-title">🔑 Token</h2>
      <div class="token-info">
        <div class="info-item">
          <span class="label">Token:</span>
          <code class="token">{{ token?.substring(0, 30) }}...</code>
        </div>
        <div class="info-item">
          <span class="label">Expire le:</span>
          <span class="value">{{ tokenExpiry }}</span>
        </div>
      </div>
    </div>

    <div class="card">
      <h2 class="section-title">🔗 Actions</h2>
      <div class="actions">
        <button v-if="!isAuthenticated" @click="$router.push('/login')" class="btn btn-primary">
          Se connecter
        </button>
        <button v-if="isAuthenticated" @click="logout" class="btn btn-danger">
          Se déconnecter
        </button>
        <button @click="checkAuth" class="btn btn-secondary">
          ♻️ Rafraîchir
        </button>
        <button @click="$router.push('/')" class="btn btn-secondary">
          🏠 Accueil
        </button>
      </div>
    </div>

    <div v-if="isAuthenticated" class="card">
      <h2 class="section-title">📦 localStorage</h2>
      <pre class="code-block">{{ localStorageData }}</pre>
    </div>
  </div>
</template>

<style scoped>
.debug-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  padding: 2rem;
}

.title {
  text-align: center;
  font-size: 2rem;
  font-weight: 800;
  color: white;
  margin-bottom: 2rem;
}

.card {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  backdrop-filter: blur(10px);
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #e2e8f0;
  margin-bottom: 1rem;
}

.status-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-weight: 600;
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.status-badge.authenticated {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.info-grid {
  display: grid;
  gap: 1rem;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.label {
  font-weight: 600;
  color: #94a3b8;
  min-width: 120px;
}

.value {
  color: #e2e8f0;
  font-family: monospace;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: 2px solid #6366f1;
}

.token {
  background: rgba(0, 0, 0, 0.3);
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  color: #10b981;
  font-family: monospace;
  font-size: 0.875rem;
}

.token-info {
  display: grid;
  gap: 1rem;
}

.actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
}

.btn-primary {
  background: linear-gradient(135deg, #5865f2 0%, #4752c4 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(88, 101, 242, 0.5);
}

.btn-danger {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
}

.btn-danger:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(239, 68, 68, 0.5);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

.code-block {
  background: rgba(0, 0, 0, 0.3);
  padding: 1rem;
  border-radius: 0.5rem;
  color: #10b981;
  font-family: monospace;
  font-size: 0.875rem;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
}

@media (max-width: 640px) {
  .debug-container {
    padding: 1rem;
  }

  .title {
    font-size: 1.5rem;
  }

  .actions {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>


