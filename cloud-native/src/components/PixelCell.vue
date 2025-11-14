<script setup lang="ts">
import { computed } from 'vue'
import { usePixelStore } from '@/stores/pixelStore'

interface Props {
  x: number
  y: number
}

const props = defineProps<Props>()
const pixelStore = usePixelStore()

const pixelColor = computed(() => pixelStore.getPixelColor(props.x, props.y))

const handleClick = () => {
  pixelStore.placePixel(props.x, props.y)
}
</script>

<template>
  <div
    @click="handleClick"
    :style="{ backgroundColor: pixelColor }"
    class="w-2 h-2 cursor-pointer hover:opacity-80 transition-opacity border border-gray-200"
    :title="`(${x}, ${y})`"
  />
</template>
