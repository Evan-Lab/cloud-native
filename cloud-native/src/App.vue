<script setup lang="ts">
import ColorPicker from '@/components/ColorPicker.vue'
import PixelGrid from '@/components/PixelGrid.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-100 p-8">
    <div class="max-w-7xl mx-auto">
      <header class="mb-8">
        <h1 class="text-5xl font-bold text-gray-900 mb-2">Pixel Place</h1>
        <p class="text-lg text-gray-600">Créez votre art pixel par pixel - Inspiré de r/place</p>
      </header>

      <div class="grid grid-cols-1 lg:grid-cols-[300px_1fr] gap-8">
        <aside class="bg-white p-6 rounded-lg shadow-lg h-fit">
          <ColorPicker />
        </aside>

        <main class="bg-white p-6 rounded-lg shadow-lg">
          <PixelGrid />
        </main>
      </div>
    </div>
  </div>
</template>
