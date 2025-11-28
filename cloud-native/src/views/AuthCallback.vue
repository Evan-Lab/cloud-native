<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const status = ref<'loading' | 'success' | 'error'>('loading')
const errorMessage = ref('')

onMounted(async () => {
  try {
    const hash = window.location.hash.substring(1) // Enlever le '#'
    const params = new URLSearchParams(hash)

    const accessToken = params.get('access_token')
    const tokenType = params.get('token_type')
    const expiresIn = params.get('expires_in')
    const scope = params.get('scope')

    if (!accessToken) {
      throw new Error("Token d'accès manquant dans la réponse Discord")
    }

    console.log('Token reçu:', {
      tokenType,
      accessToken,
      expiresIn: `${expiresIn}s (${parseInt(expiresIn || '0') / 3600}h)`,
      scope,
    })

    localStorage.setItem('discord_token', accessToken)
    localStorage.setItem('discord_token_type', tokenType || 'Bearer')
    localStorage.setItem(
      'discord_token_expiry',
      String(Date.now() + parseInt(expiresIn || '0') * 1000),
    )
    localStorage.setItem('discord_token_scope', scope || '')

    const userResponse = await fetch('/api/discord/oauth2/@me', {
      headers: {
        Authorization: `${tokenType} ${accessToken}`,
      },
    })

    if (!userResponse.ok) {
      const errorData = await userResponse.json().catch(() => ({}))
      throw new Error(errorData.message || 'Erreur lors de la récupération du profil')
    }

    const oauthData = await userResponse.json()
    const userData = oauthData.user

    console.log('Profil utilisateur récupéré:', userData.username)

    localStorage.setItem('discord_user', JSON.stringify(userData))

    window.dispatchEvent(
      new CustomEvent('discord-auth-success', {
        detail: { user: userData },
      }),
    )

    status.value = 'success'

    const redirectPath = sessionStorage.getItem('auth_redirect') || '/'
    sessionStorage.removeItem('auth_redirect')

    window.history.replaceState({}, document.title, window.location.pathname)

    setTimeout(() => {
      router.push(redirectPath)
    }, 1000)
  } catch (error) {
    console.error("Erreur lors de l'authentification:", error)
    status.value = 'error'
    errorMessage.value = error instanceof Error ? error.message : 'Une erreur est survenue'

    setTimeout(() => {
      router.push({ name: 'login' })
    }, 3000)
  }
})
</script>

<template>
  <div class="callback-container">
    <div class="bg-gradient"></div>
    <div class="bg-pattern"></div>

    <div class="callback-content">
      <div class="callback-card">
        <div v-if="status === 'loading'" class="status-section">
          <div class="spinner-container">
            <div class="spinner"></div>
          </div>
          <h2 class="status-title">Authentification en cours...</h2>
          <p class="status-message">Récupération de votre profil Discord...</p>
        </div>

        <div v-else-if="status === 'success'" class="status-section">
          <div class="success-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M5 13l4 4L19 7"
              />
            </svg>
          </div>
          <h2 class="status-title success-text">Authentification réussie !</h2>
          <p class="status-message">Connexion établie avec succès</p>
        </div>

        <div v-else-if="status === 'error'" class="status-section">
          <div class="error-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </div>
          <h2 class="status-title error-text">Échec de l'authentification</h2>
          <p class="status-message">{{ errorMessage }}</p>
          <p class="status-redirect">Redirection vers la page de connexion...</p>
        </div>

        <div class="progress-bar">
          <div
            class="progress-fill"
            :class="{
              'progress-loading': status === 'loading',
              'progress-success': status === 'success',
              'progress-error': status === 'error',
            }"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.callback-container {
  min-height: 100vh;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
  display: flex;
  align-items: center;
  justify-content: center;
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

.callback-content {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 480px;
  padding: 2rem;
  animation: fade-in-up 0.6s ease-out;
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.callback-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1.5rem;
  padding: 3rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.status-section {
  text-align: center;
  margin-bottom: 2rem;
}

.spinner-container {
  display: flex;
  justify-content: center;
  margin-bottom: 2rem;
}

.spinner {
  width: 4rem;
  height: 4rem;
  border: 4px solid rgba(99, 102, 241, 0.2);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.success-icon {
  width: 4rem;
  height: 4rem;
  margin: 0 auto 2rem;
  background: linear-gradient(135deg, #10b981, #059669);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  animation: scale-in 0.5s ease-out;
}

.success-icon svg {
  width: 2.5rem;
  height: 2.5rem;
  stroke-width: 3;
}

@keyframes scale-in {
  from {
    transform: scale(0);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.error-icon {
  width: 4rem;
  height: 4rem;
  margin: 0 auto 2rem;
  background: linear-gradient(135deg, #ef4444, #dc2626);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  animation: shake 0.5s ease-out;
}

.error-icon svg {
  width: 2.5rem;
  height: 2.5rem;
  stroke-width: 3;
}

@keyframes shake {
  0%,
  100% {
    transform: translateX(0);
  }
  25% {
    transform: translateX(-10px);
  }
  75% {
    transform: translateX(10px);
  }
}

.status-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #e2e8f0;
  margin-bottom: 0.75rem;
}

.success-text {
  color: #10b981;
}

.error-text {
  color: #ef4444;
}

.status-message {
  font-size: 1rem;
  color: #94a3b8;
  line-height: 1.5;
}

.status-info {
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 0.5rem;
  font-style: italic;
}

.status-redirect {
  font-size: 0.875rem;
  color: #64748b;
  margin-top: 0.5rem;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 9999px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.3s ease;
}

.progress-loading {
  width: 60%;
  background: linear-gradient(90deg, #6366f1, #8b5cf6);
  animation: progress-loading 1.5s ease-in-out infinite;
}

@keyframes progress-loading {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(250%);
  }
}

.progress-success {
  width: 100%;
  background: linear-gradient(90deg, #10b981, #059669);
  animation: progress-complete 1s ease-out;
}

.progress-error {
  width: 100%;
  background: linear-gradient(90deg, #ef4444, #dc2626);
  animation: progress-complete 0.5s ease-out;
}

@keyframes progress-complete {
  from {
    width: 0%;
  }
  to {
    width: 100%;
  }
}

@media (max-width: 640px) {
  .callback-content {
    padding: 1rem;
  }

  .callback-card {
    padding: 2rem 1.5rem;
  }

  .status-title {
    font-size: 1.5rem;
  }
}
</style>
