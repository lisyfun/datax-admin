<template>
  <a-drawer :visible="visible" :width="720" @cancel="close" :footer="false" title="Key 详情">
    <a-descriptions :column="2" bordered size="small">
      <a-descriptions-item label="Key">{{ keyName }}</a-descriptions-item>
      <a-descriptions-item label="类型">{{ type }}</a-descriptions-item>
      <a-descriptions-item label="TTL">{{ ttlText }}</a-descriptions-item>
    </a-descriptions>

    <a-row :gutter="16" style="margin-top:12px">
      <a-col :span="16">
        <a-card title="值">
          <a-textarea v-model="textValue" auto-size />
          <a-space style="margin-top:12px">
            <a-button type="primary" @click="save">保存</a-button>
            <a-button status="danger" @click="del">删除</a-button>
          </a-space>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="TTL 设置">
          <a-space direction="vertical" fill>
            <a-input-number v-model="ttlInput" :min="-1" placeholder="秒，-1 表示永久" />
            <a-button @click="applyTTL">应用 TTL</a-button>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { getKey, deleteKey, setKey, getTTL, setTTL } from '@/api/redis'
import { Message } from '@arco-design/web-vue'

const props = defineProps<{ visible: boolean; connId?: number; db?: number; keyName: string }>()
const emit = defineEmits(['update:visible','refresh'])
const type = ref('')
const ttl = ref<number>(-1)
const textValue = ref('')
const ttlInput = ref<number | undefined>(undefined)

const ttlText = computed(() => ttl.value === -1 ? '永久' : `${ttl.value}s`)

watch(() => props.visible, async (v) => {
  if (v && props.connId && props.keyName) {
    const { data } = await getKey({ conn_id: props.connId!, db: props.db, key: props.keyName })
    type.value = data.type
    if (data.type === 'string') {
      textValue.value = data.value || ''
    } else {
      textValue.value = JSON.stringify(data.value, null, 2)
    }
    const { data: t } = await getTTL({ conn_id: props.connId!, db: props.db, key: props.keyName })
    ttl.value = t.ttl
  }
})

function close() { emit('update:visible', false) }
async function save() {
  if (!props.connId) return
  let v: any = textValue.value
  if (type.value !== 'string') {
    try { v = JSON.parse(textValue.value) } catch { /* 保持原值 */ }
  }
  await setKey({ conn_id: props.connId!, db: props.db, type: type.value, key: props.keyName, value: v })
  Message.success('保存成功')
  close()
}
async function del() {
  if (!props.connId) return
  await deleteKey({ conn_id: props.connId!, db: props.db, key: props.keyName })
  Message.success('删除成功')
  emit('refresh')
  close()
}

async function applyTTL(){
  if (!props.connId || ttlInput.value === undefined) return
  await setTTL({ conn_id: props.connId!, db: props.db, key: props.keyName, seconds: ttlInput.value })
  Message.success('TTL 已更新')
}
</script>
