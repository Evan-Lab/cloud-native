<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { usePixelStore } from '@/stores/pixelStore'
import { DEFAULT_COLOR } from '@/types/pixel'

const pixelStore = usePixelStore()
const canvasRef = ref<HTMLCanvasElement | null>(null)

const PIXEL_SIZE = 1
const canvasWidth = pixelStore.gridWidth * PIXEL_SIZE
const canvasHeight = pixelStore.gridHeight * PIXEL_SIZE

const drawGrid = () => {
  if (!canvasRef.value) return

  const canvas = canvasRef.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ctx.fillStyle = DEFAULT_COLOR
  ctx.fillRect(0, 0, canvasWidth, canvasHeight)
  pixelStore.pixels.forEach((color, key) => {
    const [x, y] = key.split(',').map(Number)
    if (x === undefined || y === undefined) return
    ctx.fillStyle = color
    ctx.fillRect(x * PIXEL_SIZE, y * PIXEL_SIZE, PIXEL_SIZE, PIXEL_SIZE)
  })
}

const handleCanvasClick = (event: MouseEvent) => {
  if (!canvasRef.value) return

  const canvas = canvasRef.value
  const rect = canvas.getBoundingClientRect()
  const scaleX = canvas.width / rect.width
  const scaleY = canvas.height / rect.height
  const canvasX = (event.clientX - rect.left) * scaleX
  const canvasY = (event.clientY - rect.top) * scaleY
  const gridX = Math.floor(canvasX / PIXEL_SIZE)
  const gridY = Math.floor(canvasY / PIXEL_SIZE)
  pixelStore.placePixel(gridX, gridY)
  drawGrid()
}

onMounted(() => {
  drawGrid()
})

// Redessiner quand les pixels changent
watch(
  () => pixelStore.pixels.size,
  () => {
    drawGrid()
  },
)
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-bold text-gray-900">Grille de pixels</h2>
      <div class="text-sm text-gray-600">{{ pixelStore.totalPixelsPlaced }} pixels placés</div>
    </div>

    <div class="border-4 border-gray-900 overflow-auto max-w-full max-h-[600px] bg-white">
      <canvas
        ref="canvasRef"
        :width="canvasWidth"
        :height="canvasHeight"
        @click="handleCanvasClick"
        class="cursor-crosshair"
      />
    </div>

    <div class="flex gap-2">
      <button
        @click="pixelStore.clearGrid"
        class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition-colors"
      >
        Effacer la grille
      </button>
    </div>
  </div>
</template>
