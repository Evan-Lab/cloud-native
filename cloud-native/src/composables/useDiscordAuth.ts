import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

export interface DiscordUser {
  id: string
  username: string
  discriminator: string
  avatar: string | null
  email?: string
  verified?: boolean
  global_name?: string
}

const user = ref<DiscordUser | null>(null)
const token = ref<string | null>(null)
const tokenType = ref<string>('Bearer')
const isAuthenticated = computed(() => !!user.value && !!token.value)

export function useDiscordAuth() {
  const router = useRouter()

  const initAuth = () => {
    const storedToken = localStorage.getItem('discord_token')
    const storedTokenType = localStorage.getItem('discord_token_type')
    const storedUser = localStorage.getItem('discord_user')
    const tokenExpiry = localStorage.getItem('discord_token_expiry')

    if (tokenExpiry && Date.now() > parseInt(tokenExpiry)) {
      console.warn('⚠️ Token expiré, déconnexion...')
      logout()
      return false
    }

    if (storedToken && storedUser) {
      try {
        token.value = storedToken
        tokenType.value = storedTokenType || 'Bearer'
        user.value = JSON.parse(storedUser)
        return true
      } catch (error) {
        console.error('❌ Erreur lors de la restauration de la session:', error)
        logout()
        return false
      }
    }

    return false
  }

  const logout = () => {
    user.value = null
    token.value = null
    tokenType.value = 'Bearer'
    localStorage.removeItem('discord_token')
    localStorage.removeItem('discord_token_type')
    localStorage.removeItem('discord_token_expiry')
    localStorage.removeItem('discord_user')
    router.push({ name: 'login' })
  }

  const fetchUserProfile = async (): Promise<DiscordUser> => {
    if (!token.value) {
      throw new Error('Pas de token disponible')
    }

    try {
      const response = await fetch('https://discord.com/api/users/@me', {
        headers: {
          Authorization: `${tokenType.value} ${token.value}`,
        },
      })

      if (!response.ok) {
        if (response.status === 401) {
          logout()
          throw new Error('Token invalide ou expiré')
        }
        throw new Error('Erreur lors de la récupération du profil')
      }

      const userData = await response.json()
      user.value = userData
      localStorage.setItem('discord_user', JSON.stringify(userData))
      return userData
    } catch (error) {
      console.error('❌ Erreur fetchUserProfile:', error)
      throw error
    }
  }

  const checkTokenValidity = async (): Promise<boolean> => {
    if (!token.value) return false

    try {
      const response = await fetch('https://discord.com/api/users/@me', {
        headers: {
          Authorization: `${tokenType.value} ${token.value}`,
        },
      })

      if (response.ok) {
        return true
      } else {
        console.warn('⚠️ Token non valide')
        logout()
        return false
      }
    } catch (error) {
      console.error('❌ Erreur lors de la vérification du token:', error)
      return false
    }
  }

  const getTimeUntilExpiry = (): number => {
    const tokenExpiry = localStorage.getItem('discord_token_expiry')
    if (!tokenExpiry) return 0

    const expiryTime = parseInt(tokenExpiry)
    const now = Date.now()

    return Math.max(0, Math.floor((expiryTime - now) / 1000))
  }

  const getAvatarUrl = (size: number = 128): string | null => {
    if (!user.value || !user.value.avatar) return null
    return `https://cdn.discordapp.com/avatars/${user.value.id}/${user.value.avatar}.png?size=${size}`
  }

  return {
    user: computed(() => user.value),
    token: computed(() => token.value),
    tokenType: computed(() => tokenType.value),
    isAuthenticated,
    initAuth,
    logout,
    fetchUserProfile,
    checkTokenValidity,
    getTimeUntilExpiry,
    getAvatarUrl,
  }
}
