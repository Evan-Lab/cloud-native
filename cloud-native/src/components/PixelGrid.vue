<script setup lang="ts">
import { usePixelStore } from '@/stores/pixelStore'
import { DEFAULT_COLOR, PIXEL_SIZE, PixelColor, type Tool } from '@/types/pixel'
import { computed, onMounted, ref, watch } from 'vue'

const pixelStore = usePixelStore()
const canvasRef = ref<HTMLCanvasElement | null>(null)
const containerRef = ref<HTMLDivElement | null>(null)

const hoverX = ref<number | null>(null)
const hoverY = ref<number | null>(null)
const isDrawing = ref(false)

const canvasWidth = pixelStore.gridWidth * PIXEL_SIZE
const canvasHeight = pixelStore.gridHeight * PIXEL_SIZE

const primaryColors = [
  PixelColor.RED,
  PixelColor.BLUE,
  PixelColor.YELLOW,
  PixelColor.GREEN,
  PixelColor.PURPLE,
  PixelColor.PINK,
  PixelColor.WHITE,
  PixelColor.BLACK,
]

const colorGroups = [
  {
    name: 'Primaires',
    colors: [
      PixelColor.RED,
      PixelColor.BLUE,
      PixelColor.YELLOW,
      PixelColor.GREEN,
      PixelColor.PURPLE,
      PixelColor.PINK,
    ],
  },
  {
    name: 'Secondaires',
    colors: [
      PixelColor.ORANGE,
      PixelColor.CYAN,
      PixelColor.LIME,
      PixelColor.INDIGO,
      PixelColor.ROSE,
      PixelColor.TEAL,
    ],
  },
  {
    name: 'Neutres',
    colors: [
      PixelColor.WHITE,
      PixelColor.LIGHT_GRAY,
      PixelColor.GRAY,
      PixelColor.DARK_GRAY,
      PixelColor.BLACK,
    ],
  },
  {
    name: 'Tons Terre',
    colors: [PixelColor.BROWN, PixelColor.BEIGE, PixelColor.CORAL],
  },
]

const showColorPicker = ref(false)

const toggleColorPicker = () => {
  showColorPicker.value = !showColorPicker.value
}

const canvasCursor = computed(() => {
  if (pixelStore.isOnCooldown && pixelStore.currentTool !== 'eyedropper') {
    return 'cursor-not-allowed'
  }

  switch (pixelStore.currentTool) {
    case 'brush':
      return 'cursor-crosshair'
    case 'eraser':
      return 'cursor-not-allowed'
    case 'eyedropper':
      return 'cursor-copy'
    default:
      return 'cursor-crosshair'
  }
})

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
    ctx.fillStyle = pixelStore.selectedColor + '80'
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
  const coords = getGridCoordinates(event)
  if (!coords) return

  console.log('Clic sur:', coords)

  if (pixelStore.currentTool === 'eyedropper') {
    const color = pixelStore.getPixelColor(coords.x, coords.y)
    if (color !== DEFAULT_COLOR) {
      pixelStore.setSelectedColor(color)
      console.log('Couleur:', color)
    }
  } else {
    if (pixelStore.isOnCooldown) {
      return
    }
    pixelStore.placePixel(coords.x, coords.y)
  }
  drawGrid()
}

const handleMouseDown = (event: MouseEvent) => {
  handleCanvasClick(event)
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

const handleMouseUp = () => {
  isDrawing.value = false
}

const handleMouseLeave = () => {
  hoverX.value = null
  hoverY.value = null
  isDrawing.value = false
  drawGrid()
}

const selectTool = (tool: Tool) => {
  pixelStore.setTool(tool)
}

const selectColor = (color: PixelColor) => {
  pixelStore.setSelectedColor(color)
}

onMounted(() => {
  drawGrid()
})

watch([() => pixelStore.pixels.size, hoverX, hoverY], () => {
  drawGrid()
})
</script>

<template>
  <div class="relative w-full h-full overflow-hidden">
    <div class="absolute top-4 left-1/2 -translate-x-1/2 z-100 pointer-events-auto">
      <div
        class="flex items-center gap-3 bg-gray-300 my-4 backdrop-blur-xl border border-white/15 rounded-xl p-4"
      >
        <button
          @click="selectTool('brush')"
          :class="[
            'p-2.5 rounded-lg border-2 transition-all duration-200 flex items-center justify-center bg-black',
            pixelStore.currentTool === 'brush'
              ? 'border-white text-white shadow-[0_0_20px_rgba(255,255,255,0.4),0_4px_12px_rgba(0,0,0,0.3)]'
              : 'border-transparent text-gray-400 hover:border-white/30 hover:text-white hover:scale-105',
          ]"
          title="Pinceau (dessiner)"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
            />
          </svg>
        </button>

        <button
          @click="selectTool('eraser')"
          :class="[
            'p-2.5 rounded-lg border-2 transition-all duration-200 flex items-center justify-center bg-black',
            pixelStore.currentTool === 'eraser'
              ? 'border-white text-white shadow-[0_0_20px_rgba(255,255,255,0.4),0_4px_12px_rgba(0,0,0,0.3)]'
              : 'border-transparent text-gray-400 hover:border-white/30 hover:text-white hover:scale-105',
          ]"
          title="Gomme (effacer)"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>

        <button
          @click="selectTool('eyedropper')"
          :class="[
            'p-2.5 rounded-lg border-2 transition-all duration-200 flex items-center justify-center bg-black',
            pixelStore.currentTool === 'eyedropper'
              ? 'border-white text-white shadow-[0_0_20px_rgba(255,255,255,0.4),0_4px_12px_rgba(0,0,0,0.3)]'
              : 'border-transparent text-gray-400 hover:border-white/30 hover:text-white hover:scale-105',
          ]"
          title="Pipette (copier couleur)"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
            />
          </svg>
        </button>

        <div class="w-px h-7 bg-white/15 mx-1"></div>

        <div class="flex gap-2">
          <button
            v-for="color in primaryColors"
            :key="color"
            @click="selectColor(color)"
            :class="[
              'w-8 h-8 rounded-lg border-2 transition-all duration-200 relative overflow-hidden',
              pixelStore.selectedColor === color
                ? 'border-white/80 scale-110 shadow-[0_0_0_3px_rgba(255,255,255,0.2),0_4px_12px_rgba(0,0,0,0.4)]'
                : 'border-transparent hover:border-white/30 hover:scale-105 hover:shadow-lg',
            ]"
            :style="{ backgroundColor: color }"
            :title="color"
          >
            <svg
              v-if="pixelStore.selectedColor === color"
              class="w-5 h-5 absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-white drop-shadow-[0_2px_4px_rgba(0,0,0,0.5)]"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fill-rule="evenodd"
                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                clip-rule="evenodd"
              />
            </svg>
          </button>
        </div>

        <div class="w-px h-7 bg-white/15 mx-1"></div>

        <button
          @click="toggleColorPicker"
          :class="[
            'p-2.5 rounded-lg border-2 transition-all duration-200 flex items-center justify-center bg-black',
            showColorPicker
              ? 'border-white text-white shadow-[0_0_20px_rgba(255,255,255,0.4)]'
              : 'border-transparent text-gray-400 hover:border-white/30 hover:text-white hover:scale-105',
          ]"
          title="Plus de couleurs"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
            />
          </svg>
        </button>

        <div class="w-px h-7 bg-white/15 mx-1"></div>

        <div
          v-if="pixelStore.isOnCooldown"
          class="flex items-center gap-2 px-4 py-2 bg-bgray-500/20 border-2 border-bgray-500/50 rounded-lg"
        >
          <svg
            class="w-5 h-5 text-gray-400 animate-pulse"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span class="text-orange-300 font-bold text-lg font-mono">
            {{ pixelStore.cooldownRemaining }}s
          </span>
        </div>

        <div class="w-px h-7 bg-white/15 mx-1"></div>

        <button
          @click="pixelStore.clearGrid"
          class="p-2.5 rounded-lg border-2 border-transparent bg-black text-gray-400 transition-all duration-200 flex items-center justify-center hover:border-red-500/60 hover:text-white hover:scale-105"
          title="Effacer tout"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
            />
          </svg>
        </button>
      </div>
    </div>

    <div
      v-if="showColorPicker"
      class="absolute top-24 left-1/2 -translate-x-1/2 z-100 pointer-events-auto"
    >
      <div
        class="bg-gray-800/95 backdrop-blur-xl border border-white/15 rounded-xl p-6 shadow-[0_20px_60px_rgba(0,0,0,0.5)] max-w-md"
      >
        <div class="space-y-4">
          <div v-for="group in colorGroups" :key="group.name" class="space-y-2">
            <h3 class="text-xs font-semibold text-gray-400 text-center">{{ group.name }}</h3>
            <div class="grid grid-cols-6 gap-2">
              <button
                v-for="color in group.colors"
                :key="color"
                @click="selectColor(color)"
                :class="[
                  'w-10 h-10 rounded-lg border-2 transition-all duration-200 relative overflow-hidden',
                  pixelStore.selectedColor === color
                    ? 'border-white/80 scale-110 shadow-[0_0_0_3px_rgba(255,255,255,0.2),0_4px_12px_rgba(0,0,0,0.4)]'
                    : 'border-transparent hover:border-white/30 hover:scale-105 hover:shadow-lg',
                ]"
                :style="{ backgroundColor: color }"
                :title="color"
              >
                <svg
                  v-if="pixelStore.selectedColor === color"
                  class="w-5 h-5 absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-white drop-shadow-[0_2px_4px_rgba(0,0,0,0.5)]"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fill-rule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  />
                </svg>
              </button>
            </div>
          </div>

          <div class="pt-4 border-t border-white/10">
            <div class="flex items-center justify-center gap-4">
              <div
                class="w-16 h-16 rounded-xl border-4 border-white/20 shadow-lg"
                :style="{ backgroundColor: pixelStore.selectedColor }"
              ></div>
              <div>
                <div class="text-xs text-gray-400 mb-1">Couleur active</div>
                <div class="text-sm font-mono text-white font-bold">
                  {{ pixelStore.selectedColor }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div ref="containerRef" class="absolute top-0 left-0 w-full h-full overflow-hidden">
      <canvas
        ref="canvasRef"
        :width="canvasWidth"
        :height="canvasHeight"
        @mousedown="handleMouseDown"
        @mouseup="handleMouseUp"
        @mousemove="handleMouseMove"
        @mouseleave="handleMouseLeave"
        :class="['canvas-pixel w-full h-full', canvasCursor]"
        style="image-rendering: pixelated; image-rendering: crisp-edges"
      />
    </div>
  </div>
</template>
