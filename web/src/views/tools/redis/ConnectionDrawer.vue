<template>
  <a-drawer :visible="visible" title="连接管理" :width="520" @cancel="close" ok-text="关闭" cancel-text="关闭">
    <a-input v-model="keyword" placeholder="搜索连接..." allow-clear @input="onSearch" />
    <a-list :bordered="false" style="margin-top:8px">
      <a-list-item v-for="c in list" :key="c.id">
        <div style="flex:1">{{ c.name }}（{{ c.host }}:{{ c.port }} / db={{ c.db }}）</div>
        <a-space>
          <a-button size="small" @click="fav(c.id)" :type="c.favorite ? 'primary' : 'secondary'">★ 收藏</a-button>
          <a-button size="small" @click="test(c.id)">测试</a-button>
          <a-button size="small" status="danger" @click="remove(c.id)">删除</a-button>
        </a-space>
      </a-list-item>
    </a-list>
    <a-divider />
    <a-form :model="form" label-col="6">
      <a-form-item label="名称"><a-input v-model="form.name" /></a-form-item>
      <a-form-item label="地址"><a-input v-model="form.host" /></a-form-item>
      <a-form-item label="端口"><a-input-number v-model="form.port" :min="1" :max="65535" /></a-form-item>
      <a-form-item label="用户名"><a-input v-model="form.username" /></a-form-item>
      <a-form-item label="密码"><a-input-password v-model="form.password" /></a-form-item>
      <a-form-item label="DB"><a-input-number v-model="form.db" :min="0" /></a-form-item>
      <a-form-item label="TLS"><a-switch v-model="form.use_tls" /></a-form-item>
      <a-form-item>
        <a-space>
          <a-button type="primary" @click="create">新增连接</a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </a-drawer>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { listConnections, createConnection, deleteConnection, testConnection } from '@/api/redis'
import { useRedisStore } from '@/stores/redis'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits(['update:visible','refresh'])
const store = useRedisStore()
const connections = ref<any[]>([])
const keyword = ref('')
const list = computed(() => {
  store.setKeyword(keyword.value)
  return store.filtered()
})
const form = ref({ name:'', host:'', port:6379, username:'', password:'', db:0, use_tls:false })

function close(){ emit('update:visible', false) }

async function fetchList(){ const { data } = await listConnections({ page:1, page_size:100 }); connections.value = data.items || [] }
onMounted(fetchList)

async function test(id:number){ await testConnection(id) }
async function remove(id:number){ await deleteConnection(id); await fetchList(); emit('refresh') }
async function create(){ await createConnection(form.value); form.value = { name:'', host:'', port:6379, username:'', password:'', db:0, use_tls:false }; await fetchList(); emit('refresh') }
function fav(id:number){ store.toggleFavorite(id) }
function onSearch(){ /* 已通过 computed 自动过滤 */ }
</script>
