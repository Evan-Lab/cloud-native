
<script setup lang="ts">
import { ref } from 'vue'

// Configuration Discord OAuth2 - Implicit Grant
const DISCORD_CLIENT_ID = import.meta.env.VITE_DISCORD_CLIENT_ID || 'YOUR_DISCORD_CLIENT_ID'
const REDIRECT_URI = import.meta.env.VITE_REDIRECT_URI || `${window.location.origin}/auth/callback`

// 🎯 CHANGEMENT IMPORTANT: response_type=token (au lieu de code)
// Cela active l'Implicit Grant Flow de Discord
const DISCORD_OAUTH_URL = `https://discord.com/oauth2/authorize?response_type=token&client_id=${DISCORD_CLIENT_ID}&redirect_uri=${encodeURIComponent(REDIRECT_URI)}&scope=identify%20email`

const isHovering = ref(false)

const handleDiscordLogin = () => {
  // Sauvegarder l'état pour la redirection post-auth
  sessionStorage.setItem('auth_redirect', window.location.pathname)
  window.location.href = DISCORD_OAUTH_URL
}
</script>

<template>
  <div class="login-container">
    <div class="bg-gradient"></div>
    <div class="bg-pattern"></div>

    <div class="login-content">
      <div class="login-card">
        <!-- Logo et titre -->
        <div class="logo-section">
          <div class="logo-container">
            <svg class="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1v-4zM14 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
            </svg>
          </div>
          <h1 class="login-title">Pixel Place</h1>
          <p class="login-subtitle">Connectez-vous pour commencer à créer</p>
        </div>

        <!-- Séparateur décoratif -->
        <div class="divider">
          <div class="divider-line"></div>
          <span class="divider-text">Bienvenue</span>
          <div class="divider-line"></div>
        </div>

        <!-- Bouton Discord -->
        <button
          @click="handleDiscordLogin"
          @mouseenter="isHovering = true"
          @mouseleave="isHovering = false"
          class="discord-button"
          :class="{ 'hovered': isHovering }"
        >
          <svg class="discord-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515a.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0a12.64 12.64 0 0 0-.617-1.25a.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057a19.9 19.9 0 0 0 5.993 3.03a.078.078 0 0 0 .084-.028a14.09 14.09 0 0 0 1.226-1.994a.076.076 0 0 0-.041-.106a13.107 13.107 0 0 1-1.872-.892a.077.077 0 0 1-.008-.128a10.2 10.2 0 0 0 .372-.292a.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127a12.299 12.299 0 0 1-1.873.892a.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028a19.839 19.839 0 0 0 6.002-3.03a.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419c0-1.333.956-2.419 2.157-2.419c1.21 0 2.176 1.096 2.157 2.42c0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419c0-1.333.955-2.419 2.157-2.419c1.21 0 2.176 1.096 2.157 2.42c0 1.333-.946 2.418-2.157 2.418z"/>
          </svg>
          <span class="button-text">Se connecter avec Discord</span>
          <svg class="arrow-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
        </button>

        <!-- Information -->
        <div class="info-box">
          <svg class="info-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="info-text">
            Authentification sécurisée avec Discord
          </p>
        </div>

        <!-- Footer badge -->
        <div class="footer-badge">
          <span class="badge-icon">✨</span>
          <span>Inspiré de r/place</span>
        </div>
      </div>

      <!-- Décoration -->
      <div class="decoration-circles">
        <div class="circle circle-1"></div>
        <div class="circle circle-2"></div>
        <div class="circle circle-3"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ... tous les styles existants restent identiques ... */
.login-container {
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
  background: radial-gradient(circle at 20% 50%, rgba(99, 102, 241, 0.1) 0%, transparent 50%),
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

.login-content {
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

.login-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1.5rem;
  padding: 3rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  position: relative;
}

.logo-section {
  text-align: center;
  margin-bottom: 2.5rem;
}

.logo-container {
  width: 5rem;
  height: 5rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #6366f1 0%, #ec4899 100%);
  border-radius: 1.5rem;
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.3);
  color: white;
  margin-bottom: 1.5rem;
  animation: pulse-glow 2s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% {
    box-shadow: 0 8px 24px rgba(99, 102, 241, 0.3);
    transform: scale(1);
  }
  50% {
    box-shadow: 0 8px 32px rgba(99, 102, 241, 0.5);
    transform: scale(1.05);
  }
}

.logo-container svg {
  width: 3rem;
  height: 3rem;
}

.login-title {
  font-size: 2.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, #ffffff 0%, #e2e8f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.02em;
  margin-bottom: 0.5rem;
}

.login-subtitle {
  font-size: 1rem;
  color: #94a3b8;
  font-weight: 500;
}

.divider {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin: 2rem 0;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
}

.divider-text {
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.discord-button {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  background: linear-gradient(135deg, #5865f2 0%, #4752c4 100%);
  border: none;
  border-radius: 0.75rem;
  color: white;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(88, 101, 242, 0.3);
  position: relative;
  overflow: hidden;
}

.discord-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.discord-button.hovered::before {
  left: 100%;
}

.discord-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(88, 101, 242, 0.5);
}

.discord-button:active {
  transform: translateY(0);
}

.discord-icon {
  width: 1.5rem;
  height: 1.5rem;
  transition: transform 0.3s ease;
}

.discord-button:hover .discord-icon {
  transform: rotate(-10deg) scale(1.1);
}

.button-text {
  flex: 1;
}

.arrow-icon {
  width: 1.25rem;
  height: 1.25rem;
  transition: transform 0.3s ease;
}

.discord-button:hover .arrow-icon {
  transform: translateX(4px);
}

.info-box {
  display: flex;
  gap: 0.75rem;
  margin-top: 2rem;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 0.75rem;
  animation: pulse-subtle 3s ease-in-out infinite;
}

@keyframes pulse-subtle {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}

.info-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #60a5fa;
  flex-shrink: 0;
  margin-top: 0.125rem;
}

.info-text {
  font-size: 0.875rem;
  color: #cbd5e1;
  line-height: 1.5;
}

.footer-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 2rem;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 9999px;
  color: #e2e8f0;
  font-size: 0.875rem;
  font-weight: 500;
}

.badge-icon {
  font-size: 1rem;
  animation: sparkle 2s ease-in-out infinite;
}

@keyframes sparkle {
  0%, 100% {
    transform: scale(1) rotate(0deg);
  }
  25% {
    transform: scale(1.2) rotate(-10deg);
  }
  75% {
    transform: scale(1.1) rotate(10deg);
  }
}

.decoration-circles {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  overflow: hidden;
}

.circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  animation: float 20s ease-in-out infinite;
}

.circle-1 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, #6366f1, #ec4899);
  top: -150px;
  right: -150px;
  animation-delay: 0s;
}

.circle-2 {
  width: 200px;
  height: 200px;
  background: linear-gradient(135deg, #ec4899, #8b5cf6);
  bottom: -100px;
  left: -100px;
  animation-delay: 7s;
}

.circle-3 {
  width: 150px;
  height: 150px;
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  top: 50%;
  left: 10%;
  animation-delay: 14s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -30px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

@media (max-width: 640px) {
  .login-content {
    padding: 1rem;
  }

  .login-card {
    padding: 2rem 1.5rem;
  }

  .login-title {
    font-size: 2rem;
  }

  .logo-container {
    width: 4rem;
    height: 4rem;
  }

  .logo-container svg {
    width: 2.5rem;
    height: 2.5rem;
  }
}
</style>
