<template>
  <div class="json-formatter">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <h2 class="title">JSON 格式化工具</h2>
        <div class="status-info">
          <a-tag v-if="isValidJson" color="green">
            <template #icon><icon-check-circle /></template>
            有效 JSON
          </a-tag>
          <a-tag v-else-if="jsonContent && !isValidJson" color="red">
            <template #icon><icon-close-circle /></template>
            无效 JSON
          </a-tag>
          <span v-if="jsonStats.size" class="stats">
            {{ jsonStats.lines }} 行 · {{ jsonStats.size }} 字符
          </span>
        </div>
      </div>
      <div class="toolbar-right">
        <a-space>
          <a-button-group>
            <a-button
              :type="viewMode === 'tree' ? 'primary' : 'outline'"
              @click="viewMode = 'tree'"
            >
              <template #icon><icon-mind-mapping /></template>
              树形视图
            </a-button>
            <a-button
              :type="viewMode === 'text' ? 'primary' : 'outline'"
              @click="viewMode = 'text'"
            >
              <template #icon><icon-code /></template>
              文本视图
            </a-button>
          </a-button-group>
          <a-divider direction="vertical" />
          <a-button @click="formatJson" :disabled="!jsonContent">
            <template #icon><icon-code-square /></template>
            格式化
          </a-button>
          <a-button @click="compressJson" :disabled="!jsonContent">
            <template #icon><icon-safe /></template>
            压缩
          </a-button>
          <a-button @click="copyContent" :disabled="!jsonContent">
            <template #icon><icon-copy /></template>
            复制
          </a-button>
          <a-button @click="clearContent">
            <template #icon><icon-delete /></template>
            清空
          </a-button>
        </a-space>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 左侧编辑器 -->
      <div class="editor-panel" :style="{ width: leftWidth + '%' }">
        <div class="panel-header">
          <span class="panel-title">JSON 编辑器</span>
          <a-space size="small">
            <a-button size="mini" @click="uploadFile">
              <template #icon><icon-upload /></template>
              上传文件
            </a-button>
            <a-button size="mini" @click="downloadFile" :disabled="!isValidJson">
              <template #icon><icon-download /></template>
              下载
            </a-button>
          </a-space>
        </div>
        <div class="editor-container">
          <a-textarea
            v-model="jsonContent"
            :auto-size="false"
            placeholder="请输入或粘贴 JSON 内容..."
            allow-clear
            @input="handleInput"
            @paste="handlePaste"
            class="json-editor"
          />
          <div v-if="errorMessage" class="error-message">
            <icon-exclamation-circle />
            {{ errorMessage }}
          </div>
        </div>
      </div>

      <!-- 可拖动的分割线 -->
      <div
        class="resizer"
        @mousedown="startResize"
        @dblclick="resetSplit"
      >
        <div class="resizer-line"></div>
      </div>

      <!-- 右侧预览面板 -->
      <div class="preview-panel" :style="{ width: rightWidth + '%' }">
        <div class="panel-header">
          <span class="panel-title">
            {{ viewMode === 'tree' ? '树形结构' : '格式化预览' }}
          </span>
          <a-space size="small" v-if="viewMode === 'tree' && jsonData">
            <a-input-search
              v-model="searchText"
              placeholder="搜索键或值..."
              size="mini"
              style="width: 200px"
              @search="handleSearch"
              @clear="clearSearch"
              allow-clear
            />
            <a-button size="mini" @click="expandAll">
              <template #icon><icon-expand /></template>
              展开
            </a-button>
            <a-button size="mini" @click="collapseAll">
              <template #icon><icon-shrink /></template>
              收起
            </a-button>
          </a-space>
        </div>
        <div class="preview-container">
          <!-- 空状态 -->
          <div v-if="!jsonData" class="empty-state">
            <icon-file />
            <p>在左侧输入 JSON 内容</p>
            <p class="empty-tip">支持拖拽文件上传</p>
          </div>

          <!-- 树形视图 -->
          <div v-else-if="viewMode === 'tree'" class="tree-view">
            <a-tree
              :data="filteredTreeData"
              v-model:expanded-keys="expandedKeys"
              :show-line="true"
              :show-icon="false"
              block-node
            >
              <template #title="nodeData">
                <div class="tree-node" :class="{ 'highlight': nodeData.highlight }">
                  <span class="node-key">{{ nodeData.title }}</span>
                  <span v-if="!isObject(nodeData.value)" class="node-separator">:</span>
                  <span v-if="!isObject(nodeData.value)">
                    {{ formatValue(nodeData.value) }}
                  </span>
                  <span v-if="isObject(nodeData.value)" class="node-type">
                    {{ getObjectTypeLabel(nodeData.value) }}
                  </span>
                </div>
              </template>
            </a-tree>
          </div>

          <!-- 文本视图 -->
          <div v-else class="text-view">
            <pre class="formatted-json">{{ formattedJson }}</pre>
          </div>
        </div>
      </div>
    </div>

    <!-- 隐藏的文件输入 -->
    <input
      ref="fileInput"
      type="file"
      accept=".json,.txt"
      style="display: none"
      @change="handleFileUpload"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  IconFile, IconExpand, IconShrink, IconCheckCircle, IconCloseCircle,
  IconExclamationCircle, IconMindMapping, IconCode, IconCodeSquare,
  IconSafe, IconCopy, IconDelete, IconUpload, IconDownload
} from '@arco-design/web-vue/es/icon'

// 响应式数据
const jsonContent = ref('')
const jsonData = ref<any>(null)
const expandedKeys = ref<string[]>([])
const viewMode = ref<'tree' | 'text'>('tree')
const searchText = ref('')
const errorMessage = ref('')
const fileInput = ref<HTMLInputElement>()

// 分割线拖拽相关
const leftWidth = ref(50)
const rightWidth = computed(() => 100 - leftWidth.value)
const isResizing = ref(false)

// JSON 验证和统计
const isValidJson = computed(() => {
  if (!jsonContent.value.trim()) return false
  try {
    JSON.parse(jsonContent.value)
    return true
  } catch {
    return false
  }
})

const jsonStats = computed(() => {
  if (!jsonContent.value) return { lines: 0, size: 0 }
  return {
    lines: jsonContent.value.split('\n').length,
    size: jsonContent.value.length
  }
})

const formattedJson = computed(() => {
  if (!isValidJson.value) return jsonContent.value
  try {
    return JSON.stringify(JSON.parse(jsonContent.value), null, 2)
  } catch {
    return jsonContent.value
  }
})

// 将 JSON 数据转换为树形结构
const convertToTree = (data: any, parentKey = 'root', path = ''): any[] => {
  if (typeof data !== 'object' || data === null) {
    return []
  }

  return Object.entries(data).map(([key, value], index) => {
    const nodeKey = `${parentKey}-${key}-${index}`
    const currentPath = path ? `${path}.${key}` : key
    const isObject = typeof value === 'object' && value !== null

    return {
      key: nodeKey,
      title: key,
      value: value,
      path: currentPath,
      children: isObject ? convertToTree(value, nodeKey, currentPath) : undefined,
      highlight: false
    }
  })
}

const treeData = computed(() => {
  if (!jsonData.value) return []
  return convertToTree(jsonData.value)
})

// 搜索过滤后的树形数据
const filteredTreeData = computed(() => {
  if (!searchText.value) return treeData.value
  return filterTreeData(treeData.value, searchText.value.toLowerCase())
})

// 递归过滤树形数据
const filterTreeData = (nodes: any[], searchTerm: string): any[] => {
  return nodes.reduce((filtered: any[], node: any) => {
    const matchesKey = node.title.toLowerCase().includes(searchTerm)
    const matchesValue = !isObject(node.value) &&
      String(node.value).toLowerCase().includes(searchTerm)

    let filteredChildren: any[] = []
    if (node.children) {
      filteredChildren = filterTreeData(node.children, searchTerm)
    }

    if (matchesKey || matchesValue || filteredChildren.length > 0) {
      filtered.push({
        ...node,
        highlight: matchesKey || matchesValue,
        children: filteredChildren.length > 0 ? filteredChildren : node.children
      })
    }

    return filtered
  }, [])
}

// 事件处理方法
const handleInput = () => {
  errorMessage.value = ''
  try {
    if (!jsonContent.value.trim()) {
      jsonData.value = null
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonData.value = obj
    nextTick(() => {
      if (viewMode.value === 'tree') {
        expandAll()
      }
    })
  } catch (error: any) {
    jsonData.value = null
    errorMessage.value = `JSON 解析错误: ${error.message}`
  }
}

// 搜索相关方法
const handleSearch = (value: string) => {
  searchText.value = value
  if (value && filteredTreeData.value.length > 0) {
    expandAll()
  }
}

const clearSearch = () => {
  searchText.value = ''
}

// 文件操作方法
const uploadFile = () => {
  fileInput.value?.click()
}

const handleFileUpload = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target?.result as string
    jsonContent.value = content
    handleInput()
    Message.success('文件上传成功')
  }
  reader.onerror = () => {
    Message.error('文件读取失败')
  }
  reader.readAsText(file)
}

const downloadFile = () => {
  if (!isValidJson.value) {
    Message.warning('请确保 JSON 格式正确')
    return
  }

  const blob = new Blob([formattedJson.value], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'formatted.json'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  Message.success('文件下载成功')
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

// JSON 操作方法
const formatJson = () => {
  try {
    if (!jsonContent.value.trim()) {
      Message.warning('请输入 JSON 内容')
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonContent.value = JSON.stringify(obj, null, 2)
    jsonData.value = obj
    errorMessage.value = ''
    Message.success('格式化成功')
  } catch (error: any) {
    errorMessage.value = `JSON 解析错误: ${error.message}`
    Message.error('无效的 JSON 格式')
  }
}

const compressJson = () => {
  try {
    if (!jsonContent.value.trim()) {
      Message.warning('请输入 JSON 内容')
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonContent.value = JSON.stringify(obj)
    jsonData.value = obj
    errorMessage.value = ''
    Message.success('压缩成功')
  } catch (error: any) {
    errorMessage.value = `JSON 解析错误: ${error.message}`
    Message.error('无效的 JSON 格式')
  }
}

const clearContent = () => {
  jsonContent.value = ''
  jsonData.value = null
  errorMessage.value = ''
  searchText.value = ''
  Message.success('已清空')
}

const copyContent = async () => {
  if (!jsonContent.value) {
    Message.warning('没有可复制的内容')
    return
  }
  try {
    const textToCopy = viewMode.value === 'text' ? formattedJson.value : jsonContent.value
    await navigator.clipboard.writeText(textToCopy)
    Message.success('复制成功')
  } catch (error) {
    Message.error('复制失败')
  }
}

// 树形操作方法
const expandAll = () => {
  const allKeys = treeData.value
    .map(node => getAllKeys(node))
    .flat()
  expandedKeys.value = allKeys
}

const collapseAll = () => {
  expandedKeys.value = []
}

// 分割线拖拽方法
const startResize = (e: MouseEvent) => {
  isResizing.value = true
  const startX = e.clientX
  const startLeftWidth = leftWidth.value

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizing.value) return

    const container = document.querySelector('.main-content') as HTMLElement
    if (!container) return

    const containerRect = container.getBoundingClientRect()
    const deltaX = e.clientX - startX
    const deltaPercent = (deltaX / containerRect.width) * 100

    let newLeftWidth = startLeftWidth + deltaPercent

    // 限制最小和最大宽度
    newLeftWidth = Math.max(20, Math.min(80, newLeftWidth))

    leftWidth.value = newLeftWidth
  }

  const handleMouseUp = () => {
    isResizing.value = false
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  }

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

const resetSplit = () => {
  leftWidth.value = 50
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

// 工具方法
const getValueType = (value: any): string => {
  if (value === null) return 'null'
  if (Array.isArray(value)) return 'array'
  return typeof value
}

const getValueLengthClass = (value: any): string => {
  if (typeof value === 'string' && value.length > 100) {
    return 'long-text'
  }
  return ''
}

const formatValue = (value: any): string => {
  if (value === null) return 'null'
  if (typeof value === 'object') return ''
  if (typeof value === 'string') return `"${value}"`
  return String(value)
}

const getObjectTypeLabel = (value: any): string => {
  if (Array.isArray(value)) {
    return `Array(${value.length})`
  }
  if (value === null) {
    return 'null'
  }
  if (typeof value === 'object') {
    const keys = Object.keys(value)
    return `Object(${keys.length})`
  }
  return ''
}

const isObject = (value: any): boolean => {
  return typeof value === 'object' && value !== null
}

// 监听器
watch(treeData, () => {
  if (treeData.value.length > 0 && viewMode.value === 'tree') {
    nextTick(() => {
      expandAll()
    })
  }
}, { immediate: true })

// 监听搜索文本变化
watch(searchText, (newValue) => {
  if (newValue && filteredTreeData.value.length > 0) {
    nextTick(() => {
      expandAll()
    })
  }
})
</script>

<style scoped>
.json-formatter {
  height: calc(100vh - 60px); /* 减去顶部导航栏高度 */
  display: flex;
  flex-direction: column;
  background: var(--color-bg-1);
}

/* 工具栏样式 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border-2);
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-1);
}

.status-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stats {
  font-size: 12px;
  color: var(--color-text-3);
}

.toolbar-right {
  display: flex;
  align-items: center;
}

/* 主要内容区域 */
.main-content {
  flex: 1;
  display: flex;
  min-height: 0;
  position: relative;
}

/* 面板样式 */
.editor-panel,
.preview-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.editor-panel {
  border-right: none;
}

/* 可拖动分割线 */
.resizer {
  width: 8px;
  background: var(--color-border-2);
  cursor: col-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  flex-shrink: 0;
  transition: background-color 0.2s ease;
}

.resizer:hover {
  background: var(--color-primary-light-4);
}

.resizer-line {
  width: 2px;
  height: 40px;
  background: var(--color-border-3);
  border-radius: 1px;
  transition: all 0.2s ease;
}

.resizer:hover .resizer-line {
  background: var(--color-primary-6);
  height: 60px;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border-2);
  flex-shrink: 0;
}

.panel-title {
  font-weight: 500;
  color: var(--color-text-1);
}

/* 编辑器容器 */
.editor-container {
  flex: 1;
  position: relative;
  display: flex;
  flex-direction: column;
}

.json-editor {
  flex: 1;
  border: none !important;
  border-radius: 0 !important;
}

.json-editor :deep(.arco-textarea-wrapper) {
  height: 100%;
  border: none;
  border-radius: 0;
}

.json-editor :deep(.arco-textarea) {
  height: 100%;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  border: none;
  resize: none;
  padding: 16px;
}

.error-message {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--color-danger-light-1);
  color: var(--color-danger-6);
  padding: 8px 16px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  border-top: 1px solid var(--color-danger-3);
}

/* 预览容器 */
.preview-container {
  flex: 1;
  overflow: auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-3);
  gap: 8px;
}

.empty-state svg {
  font-size: 48px;
  opacity: 0.5;
}

.empty-tip {
  font-size: 12px;
  opacity: 0.7;
}

/* 树形视图 */
.tree-view {
  flex: 1;
  overflow: auto;
  padding: 8px;
  height: 100%;
  min-height: 0;
  width: 100%;
}

.tree-node {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  display: flex;
  align-items: flex-start;
  gap: 4px;
  width: 100%;
  line-height: 1.4;

  /* 默认单行显示 */
  white-space: nowrap;
  overflow: hidden;

  /* 包含长文本时允许换行 */
  &:has(.long-text) {
    white-space: normal;
    overflow: visible;
    align-items: flex-start;
  }
}

.tree-node.highlight {
  background: var(--color-primary-light-1);
  border-radius: 2px;
  padding: 0 4px;
}

.node-key {
  color: var(--color-primary-6);
  font-weight: 500;
  flex-shrink: 0;
  white-space: nowrap;
}

.node-separator {
  color: var(--color-text-3);
}

.node-value {
  &.string {
    color: #22863a;
    /* 短文本：单行显示，超出显示省略号 */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 400px;

    /* 长文本：多行显示，允许换行，充分利用可用空间 */
    &.long-text {
      white-space: pre-wrap;
      word-break: break-word;
      max-width: calc(100% - 100px); /* 使用calc计算，减去缩进和图标空间 */
      min-width: 300px; /* 降低最小宽度 */
      width: auto;
      overflow: visible;
      text-overflow: unset;
    }
  }
  &.number {
    color: #005cc5;
    white-space: nowrap;
  }
  &.boolean {
    color: #e36209;
    white-space: nowrap;
  }
  &.null {
    color: #b31d28;
    font-style: italic;
    white-space: nowrap;
  }
}

.node-type {
  color: var(--color-text-3);
  font-size: 11px;
  background: var(--color-fill-2);
  padding: 1px 4px;
  border-radius: 2px;
  margin-left: 4px;
}

/* 文本视图 */
.text-view {
  flex: 1;
  overflow: auto;
  padding: 16px;
  height: 100%;
  min-height: 0;
}

.formatted-json {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  margin: 0;
  white-space: pre-wrap;
  color: var(--color-text-1);
  min-height: 100%;
}

/* 树形组件样式覆盖 */
:deep(.arco-tree) {
  height: 100%;
  overflow: auto;
}

:deep(.arco-tree-node) {
  width: 100%;
}

:deep(.arco-tree-node-title) {
  width: 100%;
  line-height: 1.4;
  min-height: auto;

  /* 默认单行显示 */
  overflow: hidden;
  white-space: nowrap;

  /* 包含长文本时允许换行 */
  &:has(.long-text) {
    overflow: visible;
    white-space: normal;
  }
}

:deep(.arco-tree-node-content) {
  padding: 2px 0;
  width: 100%;
  align-items: flex-start;
}

:deep(.arco-tree-node-switcher) {
  color: var(--color-text-3);
  flex-shrink: 0;
  margin-top: 2px;
}

:deep(.arco-tree-node-indent) {
  flex-shrink: 0;
}

/* 预览容器样式优化 */
.preview-container {
  flex: 1;
  overflow: auto;
  width: 100%;
}

.tree-view {
  padding: 16px;
  width: 100%;
  box-sizing: border-box;
}

/* 确保树形节点能够充分利用宽度 */
:deep(.arco-tree-node) {
  width: 100%;
}

:deep(.arco-tree-node-title) {
  width: 100%;
  max-width: none;
}

/* 处理短文本 - 使用省略号 */
:deep(.arco-tree-node-title .node-value.string:not(.long-text)) {
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 处理长文本 - 允许换行，充分利用可用空间 */
:deep(.arco-tree-node-title .node-value.string.long-text) {
  max-width: calc(100% - 100px); /* 使用calc计算，减去缩进和图标空间 */
  min-width: 300px; /* 降低最小宽度 */
  width: auto;
  white-space: pre-wrap;
  word-break: break-word;
  overflow: visible;
  text-overflow: unset;
  line-height: 1.4;
}

/* 树节点标题容器 */
:deep(.arco-tree-node-title) {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 4px 0;
  min-height: 24px;
  width: 100%; /* 确保容器占满可用宽度 */
  flex: 1; /* 允许容器伸缩 */
}

/* 树节点容器 */
:deep(.arco-tree-node) {
  width: 100%; /* 确保节点占满可用宽度 */
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
  }

  .editor-panel {
    border-right: none;
    border-bottom: 1px solid var(--color-border-2);
  }

  .toolbar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .toolbar-left,
  .toolbar-right {
    justify-content: center;
  }

  /* 移动端调整文本宽度 */
  :deep(.arco-tree-node-title .node-value.string.long-text) {
    max-width: calc(100vw - 100px);
    min-width: 200px;
  }

  :deep(.arco-tree-node-title .node-value.string:not(.long-text)) {
    max-width: calc(100vw - 100px);
  }
}
</style>
