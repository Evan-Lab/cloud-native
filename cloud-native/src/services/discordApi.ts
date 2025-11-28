/**
 * Service pour les appels à l'API Discord
 * Gère automatiquement le proxy en production et les appels directs en développement
 */

const DISCORD_API_BASE_URL = import.meta.env.DEV
  ? '/api/discord' // Proxy Vite en développement
  : 'https://rplace-gateway4-5uir24en.ew.gateway.dev/web/api/discord' // API Gateway en production

const API_KEY = import.meta.env.VITE_API_KEY

export interface DiscordUser {
  id: string
  username: string
  discriminator: string
  avatar: string | null
  email?: string
  verified?: boolean
  global_name?: string
}

/**
 * Récupère les informations de l'utilisateur Discord
 * @param token Token d'accès Discord
 * @param tokenType Type de token (par défaut 'Bearer')
 * @returns Les données utilisateur Discord
 */
export async function fetchDiscordUser(
  token: string,
  tokenType: string = 'Bearer'
): Promise<DiscordUser> {
  const url = import.meta.env.DEV
    ? `${DISCORD_API_BASE_URL}/users/@me`
    : `${DISCORD_API_BASE_URL}/users/@me`

  const headers: HeadersInit = {
    Authorization: `${tokenType} ${token}`,
  }

  // En production, on doit passer l'API key et le token Discord
  if (!import.meta.env.DEV) {
    headers['X-API-KEY'] = API_KEY || ''
    headers['X-Discord-Token'] = `${tokenType} ${token}`
  }

  const response = await fetch(url, {
    headers,
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Token Discord invalide ou expiré')
    }
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.message || 'Erreur lors de la récupération du profil Discord')
  }

  return await response.json()
}

/**
 * Récupère les informations de l'utilisateur Discord via OAuth2
 * Utilisé par AuthCallback.vue qui attend un format { user: {...} }
 * @param token Token d'accès Discord
 * @param tokenType Type de token (par défaut 'Bearer')
 * @returns Les données utilisateur Discord dans le format OAuth2
 */
export async function fetchDiscordUserOAuth2(
  token: string,
  tokenType: string = 'Bearer'
): Promise<{ user: DiscordUser }> {
  const url = import.meta.env.DEV
    ? `${DISCORD_API_BASE_URL}/oauth2/@me`
    : `${DISCORD_API_BASE_URL}/oauth2/@me`

  const headers: HeadersInit = {
    Authorization: `${tokenType} ${token}`,
  }

  // En production, on doit passer l'API key et le token Discord
  if (!import.meta.env.DEV) {
    headers['X-API-KEY'] = API_KEY || ''
    headers['X-Discord-Token'] = `${tokenType} ${token}`
  }

  const response = await fetch(url, {
    headers,
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Token Discord invalide ou expiré')
    }
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.message || 'Erreur lors de la récupération du profil Discord')
  }

  return await response.json()
}

