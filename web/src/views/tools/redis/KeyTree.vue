<template>
  <div class="key-tree">
    <a-input v-model="search" placeholder="Enter to search" @press-enter="emitSearch" allow-clear />
    <a-list style="margin-top:8px; max-height: 420px; overflow: auto">
      <a-list-item v-for="k in flatKeys" :key="k" @click="select(k)">{{ k }}</a-list-item>
    </a-list>
    <a-space style="margin-top:8px">
      <a-button size="mini" @click="loadMore">load more</a-button>
      <a-button size="mini" status="danger" @click="loadAll">load all</a-button>
    </a-space>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { listKeys, getKey } from '@/api/redis'

const props = defineProps<{ connId?: number; pattern?: string }>()
const emit = defineEmits(['select','search'])
const cursor = ref<number>(0)
const flatKeys = ref<string[]>([])
const search = ref('')

async function fetch(next = false) {
  if (!props.connId) return
  const { data } = await listKeys({ conn_id: props.connId, pattern: props.pattern || '*', cursor: next ? cursor.value : 0, count: 200 })
  cursor.value = data.cursor
  const ks: string[] = data.keys || []
  flatKeys.value = next ? [...flatKeys.value, ...ks] : ks
}

function select(k: string) { emit('select', k) }
function loadMore(){ fetch(true) }
async function loadAll(){ while (cursor.value !== 0) { await fetch(true) } }
function emitSearch(){ emit('search', search.value) }

watch(() => [props.connId, props.pattern], () => { cursor.value = 0; fetch(false) })
onMounted(() => fetch(false))
</script>

<style scoped>
.key-tree { display: flex; flex-direction: column; gap: 8px; }
</style>

