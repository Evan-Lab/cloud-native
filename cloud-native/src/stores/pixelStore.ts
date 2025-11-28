import { drawPixel as apiDrawPixel, ApiError } from '@/services/gatewayApi'
import { DEFAULT_COLOR, GRID_HEIGHT, GRID_WIDTH, PixelColor, type Tool } from '@/types/pixel'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

const CANVAS_ID = 'zGYJpT1GTkWY95li4q0q'

export const usePixelStore = defineStore('pixel', () => {
  const pixels = ref<Map<string, PixelColor>>(new Map())
  const selectedColor = ref<PixelColor>(PixelColor.RED)
  const currentTool = ref<Tool>('brush')
  const getPixelKey = (x: number, y: number): string => `${x},${y}`

  const gridWidth = computed(() => GRID_WIDTH)
  const gridHeight = computed(() => GRID_HEIGHT)

  const getPixelColor = (x: number, y: number): PixelColor => {
    const key = getPixelKey(x, y)
    return pixels.value.get(key) || DEFAULT_COLOR
  }

  const totalPixelsPlaced = computed(() => pixels.value.size)

  const setPixel = (x: number, y: number, color: PixelColor) => {
    if (x < 0 || x >= GRID_WIDTH || y < 0 || y >= GRID_HEIGHT) {
      console.warn(`Invalid pixel coordinates: (${x}, ${y})`)
      return
    }
    const key = getPixelKey(x, y)
    pixels.value.set(key, color)
  }

  const setSelectedColor = (color: PixelColor) => {
    selectedColor.value = color
  }

  const setTool = (tool: Tool) => {
    currentTool.value = tool
  }

  const placePixel = async (x: number, y: number) => {
    const color = currentTool.value === 'eraser' ? DEFAULT_COLOR : selectedColor.value

    setPixel(x, y, color)

    try {
      await sendPixelToServer(x, y)
    } catch (error) {
      console.error('Erreur placement pixel:', error)
    }
  }

  const syncPixel = (x: number, y: number, color: string) => {
    setPixel(x, y, color as PixelColor)
  }

  const sendPixelToServer = async (x: number, y: number) => {
    const token = localStorage.getItem('discord_token')

    if (!token) {
      console.error('Pas de token Discord')
      throw new Error('Pas de token Discord disponible')
    }

    const color = currentTool.value === 'eraser' ? DEFAULT_COLOR : selectedColor.value

    try {
      await apiDrawPixel(x, y, color, token, CANVAS_ID)
    } catch (error) {
      if (error instanceof ApiError) {
        if (error.statusCode === 401) {
          console.error('Token invalide, déconnexion nécessaire')
        }
        throw error
      }
      throw error
    }
  }

  const clearGrid = () => {
    pixels.value.clear()
  }

  const exportGrid = () => {
    return {
      width: GRID_WIDTH,
      height: GRID_HEIGHT,
      pixels: Array.from(pixels.value.entries()).map(([key, color]) => {
        const [x, y] = key.split(',').map(Number)
        return { x, y, color }
      }),
    }
  }

  const importGrid = (data: { pixels: Array<{ x: number; y: number; color: PixelColor }> }) => {
    pixels.value.clear()
    data.pixels.forEach(({ x, y, color }) => {
      setPixel(x, y, color)
    })
  }

  return {
    pixels,
    selectedColor,
    currentTool,
    gridWidth,
    gridHeight,
    totalPixelsPlaced,
    getPixelColor,
    setPixel,
    setSelectedColor,
    setTool,
    placePixel,
    syncPixel,
    sendPixelToServer,
    clearGrid,
    exportGrid,
    importGrid,
  }
})
