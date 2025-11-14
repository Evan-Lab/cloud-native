<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { usePixelStore } from '@/stores/pixelStore'
import { DEFAULT_COLOR, PIXEL_SIZE } from '@/types/pixel'

const pixelStore = usePixelStore()
const canvasRef = ref<HTMLCanvasElement | null>(null)
const containerRef = ref<HTMLDivElement | null>(null)

const zoom = ref(1)
const minZoom = 0.5
const maxZoom = 4
const panX = ref(0)
const panY = ref(0)
const isPanning = ref(false)
const lastPanX = ref(0)
const lastPanY = ref(0)

const hoverX = ref<number | null>(null)
const hoverY = ref<number | null>(null)

const canvasWidth = pixelStore.gridWidth * PIXEL_SIZE
const canvasHeight = pixelStore.gridHeight * PIXEL_SIZE

const displayWidth = computed(() => canvasWidth * zoom.value)
const displayHeight = computed(() => canvasHeight * zoom.value)

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

  ctx.strokeStyle = 'rgba(0, 0, 0, 0.1)'
  ctx.lineWidth = 0.5

  for (let x = 0; x <= pixelStore.gridWidth; x++) {
    ctx.beginPath()
    ctx.moveTo(x * PIXEL_SIZE, 0)
    ctx.lineTo(x * PIXEL_SIZE, canvasHeight)
    ctx.stroke()
  }

  for (let y = 0; y <= pixelStore.gridHeight; y++) {
    ctx.beginPath()
    ctx.moveTo(0, y * PIXEL_SIZE)
    ctx.lineTo(canvasWidth, y * PIXEL_SIZE)
    ctx.stroke()
  }

  if (hoverX.value !== null && hoverY.value !== null) {
    ctx.fillStyle = pixelStore.selectedColor + '80' // 50% opacity
    ctx.fillRect(hoverX.value * PIXEL_SIZE, hoverY.value * PIXEL_SIZE, PIXEL_SIZE, PIXEL_SIZE)

    ctx.strokeStyle = pixelStore.selectedColor
    ctx.lineWidth = 2
    ctx.strokeRect(hoverX.value * PIXEL_SIZE, hoverY.value * PIXEL_SIZE, PIXEL_SIZE, PIXEL_SIZE)
  }
}

const getGridCoordinates = (event: MouseEvent): { x: number; y: number } | null => {
  if (!canvasRef.value) return null

  const canvas = canvasRef.value
  const rect = canvas.getBoundingClientRect()
  const scaleX = canvas.width / rect.width
  const scaleY = canvas.height / rect.height
  const canvasX = (event.clientX - rect.left) * scaleX
  const canvasY = (event.clientY - rect.top) * scaleY
  const gridX = Math.floor(canvasX / PIXEL_SIZE)
  const gridY = Math.floor(canvasY / PIXEL_SIZE)

  if (gridX < 0 || gridX >= pixelStore.gridWidth || gridY < 0 || gridY >= pixelStore.gridHeight) {
    return null
  }

  return { x: gridX, y: gridY }
}

const handleCanvasClick = (event: MouseEvent) => {
  if (isPanning.value) return

  const coords = getGridCoordinates(event)
  if (!coords) return

  pixelStore.placePixel(coords.x, coords.y)
  drawGrid()
}

const handleMouseMove = (event: MouseEvent) => {
  const coords = getGridCoordinates(event)
  if (coords) {
    hoverX.value = coords.x
    hoverY.value = coords.y
  } else {
    hoverX.value = null
    hoverY.value = null
  }
  drawGrid()
}

const handleMouseLeave = () => {
  hoverX.value = null
  hoverY.value = null
  drawGrid()
}

const handleWheel = (event: WheelEvent) => {
  event.preventDefault()
  const delta = event.deltaY > 0 ? -0.1 : 0.1
  const newZoom = Math.max(minZoom, Math.min(maxZoom, zoom.value + delta))
  zoom.value = newZoom
}

const zoomIn = () => {
  zoom.value = Math.min(maxZoom, zoom.value + 0.2)
}

const zoomOut = () => {
  zoom.value = Math.max(minZoom, zoom.value - 0.2)
}

const resetZoom = () => {
  zoom.value = 1
  panX.value = 0
  panY.value = 0
}

const handleMouseDown = (event: MouseEvent) => {
  if (event.button === 1 || event.shiftKey) {
    event.preventDefault()
    isPanning.value = true
    lastPanX.value = event.clientX
    lastPanY.value = event.clientY
  }
}

const handleMouseMoveGlobal = (event: MouseEvent) => {
  if (isPanning.value) {
    const dx = event.clientX - lastPanX.value
    const dy = event.clientY - lastPanY.value
    panX.value += dx
    panY.value += dy
    lastPanX.value = event.clientX
    lastPanY.value = event.clientY
  }
}

const handleMouseUp = () => {
  isPanning.value = false
}

onMounted(() => {
  drawGrid()
  window.addEventListener('mousemove', handleMouseMoveGlobal)
  window.addEventListener('mouseup', handleMouseUp)
})

watch([() => pixelStore.pixels.size, zoom, hoverX, hoverY], () => {
  drawGrid()
})
</script>

<template>
  <div class="flex flex-col gap-8 h-full">
    <div class="flex items-center justify-between flex-wrap gap-6">
      <div class="flex items-center gap-6">
        <div class="glass-card px-6 py-4 text-center">
          <div class="text-sm text-gray-400 mb-2">Pixels placés</div>
          <div class="text-2xl font-bold pixel-gradient">{{ pixelStore.totalPixelsPlaced }}</div>
        </div>
        <div v-if="hoverX !== null && hoverY !== null" class="glass-card px-6 py-4 text-center">
          <div class="text-sm text-gray-400 mb-2">Position</div>
          <div class="text-lg font-mono text-white">{{ hoverX }}, {{ hoverY }}</div>
        </div>
      </div>

      <div class="flex items-center gap-3 glass-card px-5 py-3">
        <button @click="zoomOut" class="zoom-btn" :disabled="zoom <= minZoom" title="Zoom arrière">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM13 10H7"
            />
          </svg>
        </button>
        <div
          class="px-4 py-2 bg-gray-800 rounded text-sm font-mono text-white min-w-[70px] text-center"
        >
          {{ Math.round(zoom * 100) }}%
        </div>
        <button @click="zoomIn" class="zoom-btn" :disabled="zoom >= maxZoom" title="Zoom avant">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v6m3-3H7"
            />
          </svg>
        </button>
        <button @click="resetZoom" class="zoom-btn" title="Réinitialiser">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
        </button>
      </div>
    </div>

    <div ref="containerRef" class="flex-1 canvas-wrapper">
      <div class="canvas-inner">
        <canvas
          ref="canvasRef"
          :width="canvasWidth"
          :height="canvasHeight"
          :style="{
            width: displayWidth + 'px',
            height: displayHeight + 'px',
            transform: `translate(${panX}px, ${panY}px)`,
          }"
          @click="handleCanvasClick"
          @mousemove="handleMouseMove"
          @mouseleave="handleMouseLeave"
          @mousedown="handleMouseDown"
          @wheel="handleWheel"
          class="canvas-pixel"
          :class="{ 'cursor-grab': isPanning, 'cursor-crosshair': !isPanning }"
        />
      </div>
    </div>

    <div class="flex gap-4 flex-wrap items-center">
      <button
        @click="pixelStore.clearGrid"
        class="action-btn bg-gradient-to-r from-red-600 to-pink-600 hover:from-red-700 hover:to-pink-700"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
          />
        </svg>
        Effacer tout
      </button>
      <div class="glass-card px-6 py-4 text-sm text-gray-400 text-center leading-relaxed">
        💡 Astuce: Shift+Clic ou molette pour zoomer/déplacer
      </div>
    </div>
  </div>
</template>

<style scoped>
.canvas-wrapper {
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  position: relative;
}

.canvas-inner {
  width: 100%;
  height: 100%;
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.canvas-pixel {
  border-radius: 0.5rem;
  image-rendering: pixelated;
  image-rendering: -moz-crisp-edges;
  image-rendering: crisp-edges;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  transition: transform 0.1s ease;
}

/* Custom scrollbar for canvas area */
.canvas-inner::-webkit-scrollbar {
  width: 12px;
  height: 12px;
}

.canvas-inner::-webkit-scrollbar-track {
  background: rgba(15, 23, 42, 0.5);
  border-radius: 6px;
}

.canvas-inner::-webkit-scrollbar-thumb {
  background: rgba(99, 102, 241, 0.4);
  border-radius: 6px;
  border: 2px solid rgba(15, 23, 42, 0.5);
}

.canvas-inner::-webkit-scrollbar-thumb:hover {
  background: rgba(99, 102, 241, 0.6);
}

.canvas-inner::-webkit-scrollbar-corner {
  background: rgba(15, 23, 42, 0.5);
}
</style>
