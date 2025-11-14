<script setup lang="ts">
import { PixelColor } from '@/types/pixel'
import { usePixelStore } from '@/stores/pixelStore'

const pixelStore = usePixelStore()

const colors = Object.values(PixelColor)

const selectColor = (color: PixelColor) => {
  pixelStore.setSelectedColor(color)
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <h2 class="text-xl font-bold text-gray-900">Sélectionnez une couleur</h2>
    <div class="flex gap-2 flex-wrap">
      <button
        v-for="color in colors"
        :key="color"
        @click="selectColor(color)"
        :class="[
          'w-12 h-12 rounded-lg border-4 transition-all hover:scale-110',
          pixelStore.selectedColor === color
            ? 'border-gray-900 shadow-lg scale-110'
            : 'border-gray-300',
        ]"
        :style="{ backgroundColor: color }"
        :title="color"
      />
    </div>
    <div class="text-sm text-gray-600">
      Couleur sélectionnée:
      <span class="font-bold" :style="{ color: pixelStore.selectedColor }">
        {{ pixelStore.selectedColor }}
      </span>
    </div>
  </div>
</template>
