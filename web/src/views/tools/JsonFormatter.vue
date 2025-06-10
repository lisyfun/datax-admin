<template>
  <div class="container">
    <a-card>
      <template #title>JSON 格式化</template>
      <template #extra></template>
      <a-row class="toolbar" :gutter="16">
        <a-col :span="24">
          <a-space>
            <a-button type="primary" @click="formatJson">格式化</a-button>
            <a-button @click="compressJson">压缩</a-button>
            <a-button @click="clearContent">清空</a-button>
            <a-button @click="copyContent">复制</a-button>
            <a-button
              type="text"
              class="expand-btn"
              @click="expandAll"
            >
              <template #icon>
                <icon-expand />
              </template>
              展开全部
            </a-button>
            <a-button
              type="text"
              class="collapse-btn"
              @click="collapseAll"
            >
              <template #icon>
                <icon-shrink />
              </template>
              收起全部
            </a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-row class="content-area" :gutter="16">
        <a-col :span="12" class="input-area">
          <a-textarea
            v-model="jsonContent"
            :auto-size="false"
            placeholder="请输入需要格式化的 JSON 内容"
            allow-clear
            @input="handleInput"
            @paste="handlePaste"
          ></a-textarea>
        </a-col>
        <a-col :span="12" class="preview-area">
          <div v-if="!jsonData" class="empty-state">
            <icon-file />
            <p>在左侧输入 JSON 内容后将在此处显示树形结构</p>
          </div>
          <a-tree
            v-else
            :data="treeData"
            v-model:expanded-keys="expandedKeys"
            :show-line="true"
            :show-icon="false"
            block-node
          >
            <template #title="nodeData">
              <span class="node-title">
                <template v-if="!isObject(nodeData.value)">
                  <span class="key">{{ nodeData.title }}: </span>
                  <span :class="['value', getValueType(nodeData.value)]">{{ formatValue(nodeData.value) }}</span>
                </template>
                <template v-else>
                  <span class="key">{{ nodeData.title }}</span>
                </template>
              </span>
            </template>
          </a-tree>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconFile, IconExpand, IconShrink } from '@arco-design/web-vue/es/icon'

const jsonContent = ref('')
const jsonData = ref<any>(null)
const expandedKeys = ref<string[]>([])

// 将 JSON 数据转换为树形结构
const convertToTree = (data: any, parentKey = 'root'): any[] => {
  if (typeof data !== 'object' || data === null) {
    return []
  }

  return Object.entries(data).map(([key, value], index) => {
    const nodeKey = `${parentKey}-${key}-${index}`
    const isObject = typeof value === 'object' && value !== null

    return {
      key: nodeKey,
      title: key,
      value: value,
      children: isObject ? convertToTree(value, nodeKey) : undefined
    }
  })
}

const treeData = computed(() => {
  if (!jsonData.value) return []
  return convertToTree(jsonData.value)
})

const handleInput = () => {
  try {
    if (!jsonContent.value.trim()) {
      jsonData.value = null
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonData.value = obj
    // 输入新内容时自动展开所有节点
    expandAll()
  } catch (error) {
    // 解析错误时不更新树形结构
  }
}

const handlePaste = (e: ClipboardEvent) => {
  e.preventDefault()
  const pastedText = e.clipboardData?.getData('text') || ''
  try {
    if (!pastedText.trim()) {
      return
    }
    const obj = JSON.parse(pastedText)
    jsonContent.value = JSON.stringify(obj, null, 2)
    jsonData.value = obj
    Message.success('格式化成功')
    expandAll()
  } catch (error) {
    jsonContent.value = pastedText
    Message.error('无效的 JSON 格式')
  }
}

const handleBlur = () => {
  try {
    if (!jsonContent.value.trim()) {
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonContent.value = JSON.stringify(obj, null, 2)
    jsonData.value = obj
    Message.success('格式化成功')
  } catch (error) {
    Message.error('无效的 JSON 格式')
  }
}

const formatJson = () => {
  try {
    if (!jsonContent.value) {
      Message.warning('请输入 JSON 内容')
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonContent.value = JSON.stringify(obj, null, 2)
    jsonData.value = obj
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
    jsonData.value = obj
    Message.success('压缩成功')
  } catch (error) {
    Message.error('无效的 JSON 格式')
  }
}

const clearContent = () => {
  jsonContent.value = ''
  jsonData.value = null
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

const expandAll = () => {
  expandedKeys.value = treeData.value
    .map(node => getAllKeys(node))
    .flat()
}

const collapseAll = () => {
  expandedKeys.value = []
}

const getAllKeys = (node: any): string[] => {
  const keys = [node.key]
  if (node.children) {
    node.children.forEach((child: any) => {
      keys.push(...getAllKeys(child))
    })
  }
  return keys
}

const getValueType = (value: any): string => {
  if (value === null) return 'null'
  if (Array.isArray(value)) return 'array'
  return typeof value
}

const formatValue = (value: any): string => {
  if (value === null) return 'null'
  if (typeof value === 'object') return ''
  if (typeof value === 'string') return `"${value}"`
  return String(value)
}

const isObject = (value: any): boolean => {
  return typeof value === 'object' && value !== null
}

// 监听 treeData 变化，自动展开所有节点
watch(treeData, () => {
  if (treeData.value.length > 0) {
    expandAll()
  }
}, { immediate: true })
</script>

<style scoped>
.container {
  padding: 16px;
}

.json-formatter {
  min-height: 600px;
}

.toolbar {
  margin-bottom: 16px;
}

.content-area {
  height: calc(100vh - 300px);
  min-height: 400px;
}

.input-area {
  height: 100%;
}

.input-area :deep(.arco-textarea-wrapper) {
  height: 100%;
}

.input-area :deep(.arco-textarea) {
  height: 100%;
  resize: none;
  font-family: monospace;
}

.preview-area {
  height: 100%;
  overflow: auto;
  border: 1px solid var(--color-neutral-3);
  border-radius: 4px;
  padding: 8px;
}

.empty-state {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-3);
}

.empty-state p {
  margin: 0;
}

.node-title {
  font-family: monospace;
  padding-left: 16px;
}

.key {
  color: var(--color-text-2);
}

.value {
  &.string {
    color: #22863a;
  }
  &.number {
    color: #005cc5;
  }
  &.boolean {
    color: #e36209;
  }
  &.null {
    color: #b31d28;
  }
}

:deep(.arco-tree-node-title) {
  width: 100%;
  overflow-x: auto;
}

:deep(.arco-textarea-wrapper) {
  height: 100% !important;
}

:deep(.arco-textarea) {
  height: 100% !important;
  font-family: monospace;
}

.separator {
  color: var(--color-text-3);
}
</style>
