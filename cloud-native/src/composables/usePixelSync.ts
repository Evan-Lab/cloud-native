import { usePixelStore } from '@/stores/pixelStore'
import { onBeforeUnmount, readonly, ref } from 'vue'

const WS_URL = import.meta.env.DEV
  ? `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`
  : 'wss://rplace-gateway-5uir24en.ew.gateway.dev/ws'

const API_KEY = import.meta.env.VITE_API_KEY
const RECONNECT_DELAY = 3000

export function usePixelSync() {
  const pixelStore = usePixelStore()
  const isConnected = ref(false)
  const lastError = ref<string | null>(null)

  let socket: WebSocket | null = null
  let reconnectTimeout: number | null = null
  let shouldReconnect = true

  const connect = () => {
    try {
      const canvasId = 'zGYJpT1GTkWY95li4q0q'
      const url = `${WS_URL}?apiKey=${encodeURIComponent(API_KEY)}&canvasId=${encodeURIComponent(canvasId)}`
      socket = new WebSocket(url)

      socket.onopen = () => {
        console.log('✅ WebSocket connecté')
        isConnected.value = true
        lastError.value = null
      }

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)

          if (
            typeof data.x === 'number' &&
            typeof data.y === 'number' &&
            typeof data.color === 'string'
          ) {
            pixelStore.syncPixel(data.x, data.y, data.color)
          } else {
            console.warn('Format de message WebSocket invalide:', data)
          }
        } catch (error) {
          console.error('Erreur parsing message WebSocket:', error)
        }
      }

      socket.onerror = (error) => {
        console.error('❌ Erreur WebSocket:', error)
        lastError.value = 'Erreur de connexion WebSocket'
      }

      socket.onclose = (event) => {
        console.log('🔌 WebSocket fermé', event.code, event.reason)
        isConnected.value = false

        if (shouldReconnect && !event.wasClean) {
          console.log(`⏳ Reconnexion dans ${RECONNECT_DELAY / 1000}s...`)
          reconnectTimeout = window.setTimeout(() => {
            console.log('🔄 Tentative de reconnexion...')
            connect()
          }, RECONNECT_DELAY)
        }
      }
    } catch (error) {
      console.error('Erreur lors de la connexion WebSocket:', error)
      lastError.value = 'Impossible de se connecter au serveur'
    }
  }

  const disconnect = () => {
    shouldReconnect = false

    if (reconnectTimeout !== null) {
      clearTimeout(reconnectTimeout)
      reconnectTimeout = null
    }

    if (socket) {
      socket.close(1000, 'Déconnexion volontaire')
      socket = null
    }

    isConnected.value = false
  }

  onBeforeUnmount(() => {
    disconnect()
  })

  return {
    isConnected: readonly(isConnected),
    lastError: readonly(lastError),
    connect,
    disconnect,
  }
}
