<template>
  <div class="redis-page">
    <a-page-header title="Redis 管理" subtitle="更直观的键浏览与编辑" />
    <div class="content">
      <Browser />
    </div>
  </div>
</template>

<script setup lang="ts">
import Browser from './Browser.vue'
</script>

<style scoped>
.content { padding: 16px; }
.redis-page { background: var(--color-bg-2); min-height: 100%; }
</style>
