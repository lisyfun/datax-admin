<template>
  <a-card :bordered="false">
    <a-space direction="vertical" fill>
      <a-card :bordered="true">
        <a-space wrap>
          <a-select v-model="connId" placeholder="选择连接" style="width:240px" @change="onConnChange">
            <a-option v-for="c in connections" :key="c.id" :value="c.id">{{ c.name }}</a-option>
          </a-select>
          <a-button type="outline" @click="openConnDrawer">
            <template #icon><icon-settings /></template>
            管理连接
          </a-button>
          <a-select v-model="db" placeholder="DB" style="width:100px" @change="onDbChange">
            <a-option v-for="d in dbs" :key="d" :value="d">DB{{ d }}</a-option>
          </a-select>
          <a-input v-model="pattern" placeholder="筛选模式（如 * 或 user:*）" style="width:280px" />
          <a-button type="primary" @click="() => loadKeys(false)">
            <template #icon><icon-search /></template>
            搜索
          </a-button>
          <a-button @click="doExport">
            <template #icon><icon-download /></template>
            导出 JSON
          </a-button>
          <a-button type="outline" @click="toggleCli">
            命令行模式
          </a-button>
        </a-space>
        <a-alert type="info" style="margin-top:12px" closable>
          使用筛选模式加载键；点击行或“查看”在侧栏打开详情，支持保存与删除。重命名/复制/移动在操作列。
        </a-alert>
      </a-card>

      <a-row :gutter="16">
        <a-col :span="24">
          <a-card :bordered="true" :title="cliMode ? '命令行' : '键列表'" style="margin-bottom:12px">
            <template v-if="cliMode">
              <CommandConsole :conn-id="connId" :db="db" />
            </template>
            <template v-else>
              <a-table :data="keys" :pagination="pagination" row-key="key" :bordered="false" size="small" :loading="loadingKeys" @row-click="onRowClick" @page-change="onPageChange" @page-size-change="onPageSizeChange">
                <template #columns>
                  <a-table-column title="Key" data-index="key" :ellipsis="true" />
                  <a-table-column title="类型" data-index="type" :width="120">
                    <template #cell="{ record }">
                      <a-tag>{{ record.type || '-' }}</a-tag>
                    </template>
                  </a-table-column>
                  <a-table-column title="操作">
                    <template #cell="{ record }">
                      <a-space>
                        <a-button size="small" type="primary" @click="openDetail(record.key)">查看</a-button>
                      </a-space>
                    </template>
                  </a-table-column>
                </template>
                <template #empty>
                  <a-empty description="未匹配到键" />
                </template>
              </a-table>
            </template>
          </a-card>

        </a-col>
      </a-row>


    </a-space>

    <KeyDetail v-model:visible="detailVisible" :conn-id="connId" :db="db" :key-name="currentKey" @refresh="loadKeys(false)" />
    <a-modal v-model:visible="renameVisible" title="重命名" @ok="doRename">
      <a-input v-model="renameTo" placeholder="新 Key 名" />
    </a-modal>
    <a-modal v-model:visible="copyVisible" title="复制" @ok="doCopy">
      <a-space direction="vertical" fill>
        <a-input v-model="copyTo" placeholder="目标 Key" />
        <a-input-number v-model="copyDb" :min="0" placeholder="目标DB(可选)" />
        <a-switch v-model="copyReplace">覆盖</a-switch>
      </a-space>
    </a-modal>
    <a-modal v-model:visible="moveVisible" title="移动" @ok="doMove">
      <a-input-number v-model="moveDb" :min="0" placeholder="目标DB" />
    </a-modal>
    <ConnectionDrawer v-model:visible="connDrawerVisible" @refresh="fetchConnections" />
  </a-card>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconSearch, IconSettings, IconDownload, IconEye, IconEdit, IconCopy, IconSwap } from '@arco-design/web-vue/es/icon'
import KeyDetail from './KeyDetail.vue'
import ConnectionDrawer from './ConnectionDrawer.vue'
import CommandConsole from './CommandConsole.vue'
// 类型侧栏已移除，展示全部键
import { useRedisStore } from '@/stores/redis'
import { listKeys, getKey, renameKey, copyKey, moveKey, exportKeys, selectDb, countKeys } from '@/api/redis'
import { usePermission } from '@/composables/usePermission'

const store = useRedisStore()
const connections = store.$state.connections
const connId = ref<number | undefined>(store.currentConnId)
const pattern = ref<string>('')
// 已不按类型筛选
  const keys = ref<{key:string; type:string}[]>([])
  const cursor = ref<number>(0)
  const nextCursor = ref<number>(0)
  const cursorStack = ref<number[]>([])
  const loadingKeys = ref(false)
  const pagination = reactive({ total: 0, current: 1, pageSize: 10, showTotal: true, showJumper: true, showPageSize: true, pageSizeOptions: [10, 15, 20, 50, 100] })
const detailVisible = ref(false)
const currentKey = ref('')
const connDrawerVisible = ref(false)
const dbs = Array.from({ length: 16 }, (_, i) => i)
const db = ref<number>(0)
const cliMode = ref(false)
const renameVisible = ref(false)
const copyVisible = ref(false)
const moveVisible = ref(false)
const renameTo = ref('')
const copyTo = ref('')
const copyDb = ref<number | undefined>(undefined)
const copyReplace = ref(false)
const moveDb = ref<number>(0)

function openDetail(k: string) { currentKey.value = k; detailVisible.value = true }
const onRowClick = (record: { key: string; type: string }) => { openDetail(record.key) }
function openConnDrawer() { connDrawerVisible.value = true }
function toggleCli(){ cliMode.value = !cliMode.value }

async function fetchConnections() {
  await store.fetchConnections();
  connId.value = store.currentConnId
  if (connId.value) { pattern.value = '*'; await loadKeys(false) }
}

async function loadKeys(next = false) {
  if (!connId.value) return
  loadingKeys.value = true
  const pat = (pattern.value && pattern.value.trim()) ? pattern.value : '*'
  const { data } = await listKeys({ conn_id: connId.value, db: db.value, pattern: pat, cursor: next ? nextCursor.value : cursor.value, count: pagination.pageSize })
  nextCursor.value = data.cursor
  const arr: string[] = data.keys || []
  const typed = await Promise.all(arr.map(async (k) => {
    try { const { data: d } = await getKey({ conn_id: connId.value!, db: db.value, key: k }); return { key: k, type: d.type || '-' } }
    catch { return { key: k, type: '-' } }
  }))
  if (next) {
    keys.value = keys.value.concat(typed)
  } else {
    keys.value = typed
    pagination.current = 1
    try {
      const { data: cnt } = await countKeys({ conn_id: connId.value, db: db.value, pattern: pat, batch: 200 })
      pagination.total = cnt.total || keys.value.length
    } catch {
      pagination.total = keys.value.length
    }
  }
  Message.success('已加载键列表')
  loadingKeys.value = false
}

async function ensurePageData(page: number, size: number) {
  const needed = page * size
  let guard = 10
  while (keys.value.length < needed && nextCursor.value !== 0 && guard-- > 0) {
    cursorStack.value.push(cursor.value)
    cursor.value = nextCursor.value
    await loadKeys(true)
  }
}

const onPageChange = async (current: number) => {
  pagination.current = current
  await ensurePageData(current, pagination.pageSize)
}

const onPageSizeChange = async (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
  // 重新从当前位置加载以匹配新的 pageSize
  await ensurePageData(1, size)
}
function prevPage() {
  if (cursorStack.value.length === 0) return
  cursor.value = cursorStack.value.pop() as number
  loadKeys(false)
}
function nextPage() {
  if (nextCursor.value === 0) return
  cursorStack.value.push(cursor.value)
  cursor.value = nextCursor.value
  loadKeys(false)
}

function openRename(k:string){ currentKey.value = k; renameTo.value = `${k}-copy`; renameVisible.value = true }
function openCopy(k:string){ currentKey.value = k; copyTo.value = `${k}-copy`; copyVisible.value = true }
function openMove(k:string){ currentKey.value = k; moveDb.value = 1; moveVisible.value = true }

async function doRename(){ if (!connId.value) return; await renameKey({ conn_id: connId.value, db: db.value, key: currentKey.value, new_key: renameTo.value }); renameVisible.value = false; await loadKeys(false) }
async function doCopy(){ if (!connId.value) return; await copyKey({ conn_id: connId.value, db: db.value, key: currentKey.value, dest_key: copyTo.value, replace: copyReplace.value, dest_db: copyDb.value }); copyVisible.value = false; }
async function doMove(){ if (!connId.value) return; await moveKey({ conn_id: connId.value, db: db.value, key: currentKey.value, dest_db: moveDb.value }); moveVisible.value = false; await loadKeys(false) }

async function doExport(){ if (!connId.value) return; const { data } = await exportKeys({ conn_id: connId.value, db: db.value, pattern: pattern.value, count: 200 }); const blob = new Blob([JSON.stringify(data.items || [], null, 2)], { type: 'application/json' }); const a = document.createElement('a'); a.href = URL.createObjectURL(blob); a.download = 'redis-export.json'; a.click(); URL.revokeObjectURL(a.href); }

const onConnChange = () => { db.value = 0 }
const onDbChange = async (v: number) => { if (!connId.value) return; await selectDb(connId.value, v); await loadKeys(false) }

watch(() => store.currentConnId, async (id) => { connId.value = id; if (connId.value) { pattern.value = '*'; await loadKeys(false) } })
fetchConnections()
const onKeySearch = (kw: string) => { pattern.value = kw; loadKeys(false) }
// 已移除类型切换
const { hasPermission } = usePermission()
const can = (code: string) => hasPermission(code)
</script>
