<script setup lang="ts">
import { PixelColor } from '@/types/pixel'
import { usePixelStore } from '@/stores/pixelStore'

const pixelStore = usePixelStore()

const colors = Object.values(PixelColor)

const colorGroups = [
  {
    name: 'Primaires',
    colors: [PixelColor.RED, PixelColor.BLUE, PixelColor.YELLOW, PixelColor.GREEN, PixelColor.PURPLE, PixelColor.PINK]
  },
  {
    name: 'Secondaires',
    colors: [PixelColor.ORANGE, PixelColor.CYAN, PixelColor.LIME, PixelColor.INDIGO, PixelColor.ROSE, PixelColor.TEAL]
  },
  {
    name: 'Neutres',
    colors: [PixelColor.WHITE, PixelColor.LIGHT_GRAY, PixelColor.GRAY, PixelColor.DARK_GRAY, PixelColor.BLACK]
  },
  {
    name: 'Tons Terre',
    colors: [PixelColor.BROWN, PixelColor.BEIGE, PixelColor.CORAL]
  }
]

const selectColor = (color: PixelColor) => {
  pixelStore.setSelectedColor(color)
}
</script>

<template>
  <div class="color-picker-container">
    <div class="glass-card p-6 mb-6">
      <div class="flex flex-col items-center gap-4 text-center">
        <div
          class="w-20 h-20 rounded-xl shadow-lg border-4 border-white/20"
          :style="{ backgroundColor: pixelStore.selectedColor }"
        ></div>
        <div class="w-full">
          <div class="text-sm text-gray-400 mb-2">Couleur active</div>
          <div class="text-lg font-mono text-white font-bold">{{ pixelStore.selectedColor }}</div>
        </div>
      </div>
    </div>

    <div class="space-y-5">
      <div v-for="group in colorGroups" :key="group.name" class="glass-card p-5">
        <h3 class="text-sm font-semibold text-gray-400 mb-4 text-center">{{ group.name }}</h3>
        <div class="grid grid-cols-6 gap-3">
          <button
            v-for="color in group.colors"
            :key="color"
            @click="selectColor(color)"
            :class="[
              'color-swatch',
              pixelStore.selectedColor === color && 'selected'
            ]"
            :style="{ backgroundColor: color }"
            :title="color"
          >
            <div v-if="pixelStore.selectedColor === color" class="checkmark">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </div>
          </button>
        </div>
      </div>
    </div>

    <div class="glass-card p-4 mt-5 text-center">
      <div class="text-sm text-gray-400">
        {{ colors.length }} couleurs disponibles
      </div>
    </div>
  </div>
</template>

<style scoped>
.color-picker-container {
  height: 100%;
  overflow-y: auto;
}

.color-swatch {
  aspect-ratio: 1;
  border-radius: 0.75rem;
  border: 3px solid transparent;
  transition: all 0.2s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.color-swatch:hover {
  transform: scale(1.1);
  border-color: rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}

.color-swatch.selected {
  border-color: rgba(255, 255, 255, 0.8);
  box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.2), 0 8px 20px rgba(0, 0, 0, 0.4);
  transform: scale(1.05);
}

.checkmark {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.5));
  animation: checkmark-appear 0.2s ease;
}

@keyframes checkmark-appear {
  from {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.5);
  }
  to {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
}
</style>
