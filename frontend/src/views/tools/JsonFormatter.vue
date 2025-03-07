<template>
  <div class="container">
    <a-card class="json-formatter" :bordered="false">
      <a-row :gutter="16">
        <a-col :span="24">
          <a-space>
            <a-button type="primary" @click="formatJson">格式化</a-button>
            <a-button @click="compressJson">压缩</a-button>
            <a-button @click="clearContent">清空</a-button>
            <a-button @click="copyContent">复制</a-button>
          </a-space>
        </a-col>
        <a-col :span="24" style="margin-top: 16px;">
          <a-textarea
            v-model="jsonContent"
            :auto-size="{ minRows: 10, maxRows: 20 }"
            placeholder="请输入需要格式化的 JSON 内容"
            allow-clear
          ></a-textarea>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'

const jsonContent = ref('')

const formatJson = () => {
  try {
    if (!jsonContent.value) {
      Message.warning('请输入 JSON 内容')
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonContent.value = JSON.stringify(obj, null, 2)
    Message.success('格式化成功')
  } catch (error) {
    Message.error('无效的 JSON 格式')
  }
}

const compressJson = () => {
  try {
    if (!jsonContent.value) {
      Message.warning('请输入 JSON 内容')
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonContent.value = JSON.stringify(obj)
    Message.success('压缩成功')
  } catch (error) {
    Message.error('无效的 JSON 格式')
  }
}

const clearContent = () => {
  jsonContent.value = ''
  Message.success('已清空')
}

const copyContent = async () => {
  if (!jsonContent.value) {
    Message.warning('没有可复制的内容')
    return
  }
  try {
    await navigator.clipboard.writeText(jsonContent.value)
    Message.success('复制成功')
  } catch (error) {
    Message.error('复制失败')
  }
}
</script>

<style scoped>
.container {
  padding: 16px;
}

.json-formatter {
  min-height: 500px;
}
</style>
