 <template>
  <div class="editor-panel">
    <a-card :bordered="true">
      <a-space align="center" wrap>
        <a-tag v-if="type">{{ type }}</a-tag>
        <a-input v-model="keyNameLocal" placeholder="Key 名称" style="width:280px" disabled />
        <a-input-number v-model="ttlInput" :min="-1" placeholder="TTL 秒(-1永久)" />
        <a-button @click="applyTTL">TTL</a-button>
        <a-button status="danger" @click="del">删除</a-button>
        <a-button @click="refresh">刷新</a-button>
        <a-button type="primary" @click="save">保存</a-button>
      </a-space>
    </a-card>
    <a-card style="margin-top:12px">
      <template v-if="type === 'set'">
        <a-space direction="vertical" fill>
          <a-typography-paragraph>成员 ({{ members.length }})</a-typography-paragraph>
          <a-table :data="members" :pagination="false" :bordered="false">
            <a-table-column title="Member" v-slot="{ record }">
              <span>{{ record }}</span>
            </a-table-column>
            <template #empty>
              <a-empty description="当前 Set 为空" />
            </template>
          </a-table>
        </a-space>
      </template>
      <template v-else>
        <a-textarea v-model="textValue" :auto-size="{ minRows: 16 }" placeholder="键值内容" />
      </template>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { getKey, deleteKey, setKey, getTTL, setTTL } from '@/api/redis'
import { Message } from '@arco-design/web-vue'

const props = defineProps<{ connId?: number; db?: number; keyName?: string }>()
const emit = defineEmits(['refresh'])
const type = ref('')
const ttl = ref<number>(-1)
const textValue = ref('')
const members = ref<string[]>([])
const ttlInput = ref<number | undefined>(undefined)
const keyNameLocal = computed(() => props.keyName || '')

async function load(){
  if (!props.connId || !props.keyName) return
  const { data } = await getKey({ conn_id: props.connId!, db: props.db, key: props.keyName! })
  type.value = data.type
  if (data.type === 'set') {
    members.value = Array.isArray(data.value) ? data.value : []
    textValue.value = ''
  } else {
    textValue.value = data.type === 'string' ? (data.value || '') : JSON.stringify(data.value, null, 2)
    members.value = []
  }
  const { data: t } = await getTTL({ conn_id: props.connId!, db: props.db, key: props.keyName! })
  ttl.value = t.ttl
}

watch(() => [props.connId, props.keyName], () => load())
load()

async function save(){
  if (!props.connId || !props.keyName) return
  let v: any = textValue.value
  if (type.value !== 'string') {
    try { v = JSON.parse(textValue.value) } catch { /* 保持原值 */ }
  }
  await setKey({ conn_id: props.connId!, db: props.db, type: type.value, key: props.keyName!, value: v })
  Message.success('保存成功')
}

async function del(){ if (!props.connId || !props.keyName) return; await deleteKey({ conn_id: props.connId!, db: props.db, key: props.keyName! }); Message.success('删除成功'); emit('refresh') }
async function applyTTL(){ if (!props.connId || !props.keyName || ttlInput.value === undefined) return; await setTTL({ conn_id: props.connId!, db: props.db, key: props.keyName!, seconds: ttlInput.value! }); Message.success('TTL 已更新') }
async function refresh(){ await load(); Message.success('已刷新') }
</script>

<style scoped>
.editor-panel { display: flex; flex-direction: column; }
</style>
