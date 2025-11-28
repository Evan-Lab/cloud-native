const API_BASE_URL = import.meta.env.DEV
  ? '/api/gateway'
  : 'https://rplace-gateway-5uir24en.ew.gateway.dev/web/api'

const API_KEY = import.meta.env.VITE_API_KEY

export interface PixelData {
  x: number
  y: number
  color: string
}

export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export async function drawPixel(
  x: number,
  y: number,
  color: string,
  discordToken: string,
  canvasId: string = 'zGYJpT1GTkWY95li4q0q',
): Promise<void> {
  try {
    const requestBody = { x, y, color, canvasId }
    console.log('Envoi pixel à API:', requestBody)

    const response = await fetch(`${API_BASE_URL}/draw-pixel`, {
      method: 'POST',
      headers: {
        'X-API-KEY': API_KEY,
        'X-Discord-Token': `Bearer ${discordToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestBody),
    })

    if (!response.ok) {
      if (response.status === 401) {
        console.error('Token invalide (401)')
        throw new ApiError('Token Discord invalide ou expiré', 401)
      }
      if (response.status === 429) {
        console.error('Rate limit (429)')
        throw new ApiError('Rate limit dépassé. Attendez 5 secondes.', 429)
      }

      const errorData = await response.json().catch(() => ({}))
      console.error('Erreur:', response.status, errorData)
      throw new ApiError(errorData.message || 'Erreur lors du placement du pixel', response.status)
    }
  } catch (error) {
    if (error instanceof ApiError) throw error
    console.error('Exception:', error)
    throw new ApiError('Erreur de connexion au serveur')
  }
}

export async function loadCanvas(canvasId: string = 'zGYJpT1GTkWY95li4q0q'): Promise<PixelData[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/canvas?canvasId=${canvasId}`, {
      method: 'GET',
      headers: {
        'X-API-KEY': API_KEY,
      },
    })

    if (!response.ok) {
      throw new ApiError('Erreur lors du chargement du canvas', response.status)
    }

    const data = await response.json()
    return data as PixelData[]
  } catch (error) {
    if (error instanceof ApiError) throw error
    throw new ApiError('Erreur de connexion au serveur')
  }
}
