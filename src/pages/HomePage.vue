<script setup lang="ts">
import { useCounterStore } from '@/stores/counter'

const counter = useCounterStore()
</script>

<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center">
    <div class="bg-white p-8 rounded-lg shadow-md">
      <h1 class="text-3xl font-bold text-gray-900 mb-4">Home Page</h1>
      <p class="text-gray-600 mb-4">Welcome to Vue Boilerplate with Vite + TypeScript</p>
      <div class="flex items-center gap-4">
        <button
          class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
          @click="counter.decrement"
        >
          -
        </button>
        <span class="text-xl font-semibold">Count: {{ counter.count }}</span>
        <button
          class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
          @click="counter.increment"
        >
          +
        </button>
      </div>
    </div>
  </div>
</template>
