<script setup lang="ts">
import { ref } from 'vue'
import { computed } from 'vue'

const { id, isHover = true } = defineProps<{
  id: string
  isHover?: boolean
}>()

const imageUrl = computed(() => `https://q.trap.jp/api/v3/public/icon/${id}`)
</script>

<template>
  <div>
    <img :src="imageUrl" class="usericon" />
    <v-tooltip v-if="isHover" activator="parent" location="bottom">{{ id }}</v-tooltip>
  </div>
</template>

<style scoped>
.usericon {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
}
</style>
