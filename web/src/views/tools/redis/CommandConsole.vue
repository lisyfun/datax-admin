<template>
  <div>
    <a-card :bordered="true" title="命令行">
      <a-space direction="vertical" fill>
        <a-input v-model="line" placeholder="输入命令，例如：GET foo 或 LRANGE mylist 0 10" />
        <a-space>
          <a-button type="primary" @click="run" :loading="running">执行</a-button>
          <a-button @click="clear">清空</a-button>
        </a-space>
        <a-textarea v-model="output" :auto-size="{ minRows: 12 }" readonly placeholder="结果输出" />
      </a-space>
    </a-card>
  </div>
  </template>

<script setup lang="ts">
import { ref } from 'vue'
import { execCommand } from '@/api/redis'
import { Message } from '@arco-design/web-vue'

const props = defineProps<{ connId?: number; db?: number }>()
const line = ref('')
const output = ref('')
const running = ref(false)

function clear(){ line.value = ''; output.value = '' }
async function run(){
  if (!props.connId || !line.value.trim()) return
  running.value = true
  try{
    const { data } = await execCommand({ conn_id: props.connId!, db: props.db, line: line.value })
    const res = data.result
    if (typeof res === 'string') {
      output.value = res
    } else {
      output.value = JSON.stringify(res, null, 2)
    }
    Message.success('执行成功')
  } catch(e:any){
    output.value = e?.response?.data?.error || String(e)
    Message.error('执行失败')
  } finally {
    running.value = false
  }
}

</script>

<style scoped>
</style>
