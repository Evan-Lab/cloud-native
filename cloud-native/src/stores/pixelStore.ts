import { DEFAULT_COLOR, GRID_HEIGHT, GRID_WIDTH, PixelColor } from '@/types/pixel'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const usePixelStore = defineStore('pixel', () => {
  const pixels = ref<Map<string, PixelColor>>(new Map())
  const selectedColor = ref<PixelColor>(PixelColor.RED)
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

  const placePixel = (x: number, y: number) => {
    setPixel(x, y, selectedColor.value)
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
    gridWidth,
    gridHeight,
    totalPixelsPlaced,
    getPixelColor,
    setPixel,
    setSelectedColor,
    placePixel,
    clearGrid,
    exportGrid,
    importGrid,
  }
})
