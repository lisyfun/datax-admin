<template>
  <a-card :bordered="true">
    <a-space direction="vertical" fill>
      <a-typography-title :heading="6">数据类型</a-typography-title>
      <a-radio-group v-model="current" direction="vertical" @change="emitType">
        <a-radio v-for="t in types" :key="t.value" :value="t.value">{{ t.label }}</a-radio>
      </a-radio-group>
    </a-space>
  </a-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const emit = defineEmits(['change'])
const types = [
  { value: 'string', label: 'String' },
  { value: 'hash', label: 'Hash' },
  { value: 'list', label: 'List' },
  { value: 'set', label: 'Set' },
  { value: 'zset', label: 'ZSet' },
  { value: 'stream', label: 'Stream' },
]
const current = ref('string')
function emitType(){ emit('change', current.value) }
</script>

