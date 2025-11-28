<script setup lang="ts">
import { usePixelStore } from '@/stores/pixelStore'
import { PixelColor } from '@/types/pixel'
import { computed, ref } from 'vue'

const pixelStore = usePixelStore()

const colors = Object.values(PixelColor)

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

const selectColor = (color: PixelColor) => {
  pixelStore.setSelectedColor(color)
}

const showCustomPicker = ref(false)
const customColorInput = ref('#ff0000')
const hue = ref(0)
const saturation = ref(100)
const lightness = ref(50)

const hslToHex = (h: number, s: number, l: number): string => {
  l /= 100
  const a = (s * Math.min(l, 1 - l)) / 100
  const f = (n: number) => {
    const k = (n + h / 30) % 12
    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1)
    return Math.round(255 * color)
      .toString(16)
      .padStart(2, '0')
  }
  return `#${f(0)}${f(8)}${f(4)}`
}

const customColorFromHSL = computed(() => {
  return hslToHex(hue.value, saturation.value, lightness.value)
})

const applyCustomColor = () => {
  const color = showCustomPicker.value ? customColorFromHSL.value : customColorInput.value
  pixelStore.setSelectedColor(color as PixelColor)
}

const togglePickerMode = () => {
  showCustomPicker.value = !showCustomPicker.value
}
</script>

<template>
  <div class="h-full overflow-y-auto">
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

    <div class="glass-card p-5 mb-5">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-sm font-semibold text-gray-400">Couleur personnalisée</h3>
        <button
          @click="togglePickerMode"
          class="text-xs px-3 py-1.5 rounded-lg bg-primary/20 text-indigo-300 hover:bg-primary/30 transition-colors duration-200"
        >
          {{ showCustomPicker ? 'Input natif' : 'Palette HSL' }}
        </button>
      </div>

      <div v-if="!showCustomPicker" class="space-y-3">
        <div class="flex items-center gap-3">
          <input
            type="color"
            v-model="customColorInput"
            class="w-16 h-16 rounded-lg cursor-pointer border-2 border-white/20 bg-transparent"
          />
          <div class="flex-1">
            <input
              type="text"
              v-model="customColorInput"
              placeholder="#000000"
              class="w-full px-4 py-2 rounded-lg bg-dark/50 border border-white/10 text-white font-mono text-sm focus:outline-none focus:border-primary/50 transition-colors"
            />
          </div>
        </div>
        <button
          @click="applyCustomColor"
          class="w-full px-4 py-2.5 rounded-lg bg-primary/80 hover:bg-primary text-white font-semibold transition-all duration-200 hover:scale-[1.02]"
        >
          Appliquer
        </button>
      </div>

      <div v-else class="space-y-4">
        <div
          class="w-full h-20 rounded-lg border-2 border-white/20 mb-3"
          :style="{ backgroundColor: customColorFromHSL }"
        ></div>

        <div class="space-y-2">
          <div class="flex justify-between text-xs text-gray-400">
            <span>Teinte (Hue)</span>
            <span>{{ hue }}°</span>
          </div>
          <input
            type="range"
            v-model="hue"
            min="0"
            max="360"
            class="w-full h-2 rounded-lg appearance-none cursor-pointer"
            style="
              background: linear-gradient(
                to right,
                #ff0000,
                #ffff00,
                #00ff00,
                #00ffff,
                #0000ff,
                #ff00ff,
                #ff0000
              );
            "
          />
        </div>

        <div class="space-y-2">
          <div class="flex justify-between text-xs text-gray-400">
            <span>Saturation</span>
            <span>{{ saturation }}%</span>
          </div>
          <input
            type="range"
            v-model="saturation"
            min="0"
            max="100"
            class="w-full h-2 rounded-lg appearance-none cursor-pointer bg-linear-to-r from-gray-400 to-blue-500"
          />
        </div>

        <div class="space-y-2">
          <div class="flex justify-between text-xs text-gray-400">
            <span>Luminosité (Lightness)</span>
            <span>{{ lightness }}%</span>
          </div>
          <input
            type="range"
            v-model="lightness"
            min="0"
            max="100"
            class="w-full h-2 rounded-lg appearance-none cursor-pointer bg-linear-to-r from-black via-gray-500 to-white"
          />
        </div>

        <button
          @click="applyCustomColor"
          class="w-full px-4 py-2.5 rounded-lg bg-primary/80 hover:bg-primary text-white font-semibold transition-all duration-200 hover:scale-[1.02]"
        >
          Appliquer
        </button>
      </div>
    </div>

    <div class="space-y-4">
      <div v-for="group in colorGroups" :key="group.name" class="glass-card p-4">
        <h3 class="text-xs font-semibold text-gray-400 mb-3 text-center">{{ group.name }}</h3>
        <div class="grid grid-cols-8 gap-2">
          <button
            v-for="color in group.colors"
            :key="color"
            @click="selectColor(color)"
            :class="[
              'aspect-square rounded-xl border-2 transition-all duration-200 cursor-pointer relative overflow-hidden',
              pixelStore.selectedColor === color
                ? 'border-white/80 shadow-[0_0_0_3px_rgba(255,255,255,0.2),0_6px_16px_rgba(0,0,0,0.4)] scale-105'
                : 'border-transparent hover:border-white/30 hover:scale-110 hover:shadow-[0_6px_12px_rgba(0,0,0,0.3)]',
            ]"
            :style="{ backgroundColor: color }"
            :title="color"
          >
            <div
              v-if="pixelStore.selectedColor === color"
              class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-white drop-shadow-[0_2px_4px_rgba(0,0,0,0.5)] animate-checkmark-appear"
            >
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>
          </button>
        </div>
      </div>
    </div>

    <div class="glass-card p-4 mt-4 text-center">
      <div class="text-sm text-gray-400">{{ colors.length }} couleurs disponibles</div>
    </div>
  </div>
</template>

<style scoped>
input[type='range'] {
  -webkit-appearance: none;
  appearance: none;
}

input[type='range']::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  cursor: pointer;
  border: 2px solid rgba(99, 102, 241, 0.5);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

input[type='range']::-moz-range-thumb {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  cursor: pointer;
  border: 2px solid rgba(99, 102, 241, 0.5);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

input[type='color'] {
  -webkit-appearance: none;
  appearance: none;
}

input[type='color']::-webkit-color-swatch-wrapper {
  padding: 0;
}

input[type='color']::-webkit-color-swatch {
  border: none;
  border-radius: 0.5rem;
}
</style>
