<template>
  <div class="container">
    <a-card class="crypto-tool">
      <template #title>加解密工具</template>
      <template #extra></template>
      <a-tabs default-active-key="1">
        <a-tab-pane key="1" title="AES">
          <div class="section">
            <div class="section-header" @click="toggleAesKey">
              <div class="section-title">
                <icon-safe />
                AES密钥：
              </div>
              <icon-down v-if="!isAesKeyCollapsed" />
              <icon-right v-else />
            </div>
            <a-collapse :active-key="isAesKeyCollapsed ? [] : ['1']">
              <a-collapse-item key="1">
                <a-input v-model="aesKey" placeholder="请输入AES密钥" allow-clear />
              </a-collapse-item>
            </a-collapse>
          </div>

          <div class="section">
            <div class="section-header" @click="toggleAesEncrypt">
              <div class="section-title">
                <icon-lock />
                加密：
              </div>
              <icon-down v-if="!isAesEncryptCollapsed" />
              <icon-right v-else />
            </div>
            <a-collapse :active-key="isAesEncryptCollapsed ? [] : ['1']">
              <a-collapse-item key="1">
                <a-textarea
                  v-model="aesEncryptInput"
                  :auto-size="{ minRows: 4, maxRows: 8 }"
                  placeholder="请输入需要加密的字符串"
                  allow-clear
                />
                <div class="button-group">
                  <a-space>
                    <a-button type="primary" @click="handleAesEncrypt">确定</a-button>
                    <a-button @click="aesEncryptInput = ''">清空</a-button>
                  </a-space>
                </div>
                <div class="result-label">加密后数据：</div>
                <div class="result-content" @click="copyToClipboard(aesEncryptResult)">
                  {{ aesEncryptResult }}
                  <a-button
                    v-if="aesEncryptResult"
                    type="text"
                    size="mini"
                    class="copy-btn"
                  >
                    <template #icon>
                      <icon-copy />
                    </template>
                    点击复制
                  </a-button>
                </div>
              </a-collapse-item>
            </a-collapse>
          </div>

          <div class="section">
            <div class="section-header" @click="toggleAesDecrypt">
              <div class="section-title">
                <icon-unlock />
                解密：
              </div>
              <icon-down v-if="!isAesDecryptCollapsed" />
              <icon-right v-else />
            </div>
            <a-collapse :active-key="isAesDecryptCollapsed ? [] : ['1']">
              <a-collapse-item key="1">
                <a-textarea
                  v-model="aesDecryptInput"
                  :auto-size="{ minRows: 4, maxRows: 8 }"
                  placeholder="请输入需要解密的字符串"
                  allow-clear
                />
                <div class="button-group">
                  <a-space>
                    <a-button type="primary" @click="handleAesDecrypt">确定</a-button>
                    <a-button @click="aesDecryptInput = ''">清空</a-button>
                  </a-space>
                </div>
                <div class="result-label">解密后数据：</div>
                <div class="result-content" @click="copyToClipboard(aesDecryptResult)">
                  {{ aesDecryptResult }}
                  <a-button
                    v-if="aesDecryptResult"
                    type="text"
                    size="mini"
                    class="copy-btn"
                  >
                    <template #icon>
                      <icon-copy />
                    </template>
                    点击复制
                  </a-button>
                </div>
              </a-collapse-item>
            </a-collapse>
          </div>
        </a-tab-pane>

        <a-tab-pane key="2" title="SM4">
          <div class="section">
            <div class="section-header" @click="toggleSm4Key">
              <div class="section-title">
                <icon-safe />
                SM4加密key：
              </div>
              <icon-down v-if="!isSm4KeyCollapsed" />
              <icon-right v-else />
            </div>
            <a-collapse :active-key="isSm4KeyCollapsed ? [] : ['1']">
              <a-collapse-item key="1">
                <a-input v-model="sm4Key" placeholder="请输入SM4加密key" allow-clear />
              </a-collapse-item>
            </a-collapse>
          </div>

          <div class="section">
            <div class="section-header" @click="toggleSm4Encrypt">
              <div class="section-title">
                <icon-lock />
                加密：
              </div>
              <icon-down v-if="!isSm4EncryptCollapsed" />
              <icon-right v-else />
            </div>
            <a-collapse :active-key="isSm4EncryptCollapsed ? [] : ['1']">
              <a-collapse-item key="1">
                <a-textarea
                  v-model="sm4EncryptInput"
                  :auto-size="{ minRows: 4, maxRows: 8 }"
                  placeholder="请输入需要加密的字符串"
                  allow-clear
                />
                <div class="button-group">
                  <a-space>
                    <a-button type="primary" @click="handleSm4Encrypt">确定</a-button>
                    <a-button @click="sm4EncryptInput = ''">清空</a-button>
                  </a-space>
                </div>
                <div class="result-label">加密后数据：</div>
                <div class="result-content" @click="copyToClipboard(sm4EncryptResult)">
                  {{ sm4EncryptResult }}
                  <a-button
                    v-if="sm4EncryptResult"
                    type="text"
                    size="mini"
                    class="copy-btn"
                  >
                    <template #icon>
                      <icon-copy />
                    </template>
                    点击复制
                  </a-button>
                </div>
              </a-collapse-item>
            </a-collapse>
          </div>

          <div class="section">
            <div class="section-header" @click="toggleSm4Decrypt">
              <div class="section-title">
                <icon-unlock />
                解密：
              </div>
              <icon-down v-if="!isSm4DecryptCollapsed" />
              <icon-right v-else />
            </div>
            <a-collapse :active-key="isSm4DecryptCollapsed ? [] : ['1']">
              <a-collapse-item key="1">
                <a-textarea
                  v-model="sm4DecryptInput"
                  :auto-size="{ minRows: 4, maxRows: 8 }"
                  placeholder="请输入需要解密的字符串"
                  allow-clear
                />
                <div class="button-group">
                  <a-space>
                    <a-button type="primary" @click="handleSm4Decrypt">确定</a-button>
                    <a-button @click="sm4DecryptInput = ''">清空</a-button>
                  </a-space>
                </div>
                <div class="result-label">解密后数据：</div>
                <div class="result-content" @click="copyToClipboard(sm4DecryptResult)">
                  {{ sm4DecryptResult }}
                  <a-button
                    v-if="sm4DecryptResult"
                    type="text"
                    size="mini"
                    class="copy-btn"
                  >
                    <template #icon>
                      <icon-copy />
                    </template>
                    点击复制
                  </a-button>
                </div>
              </a-collapse-item>
            </a-collapse>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- JSON 格式化弹窗 -->
    <a-modal
      v-model:visible="showJsonModal"
      title="JSON 格式化预览"
      :width="1200"
      :footer="false"
      :mask-closable="true"
      @cancel="closeJsonModal"
      class="json-modal"
    >
      <div class="json-formatter">
        <a-row class="toolbar" :gutter="16">
          <a-col :span="24">
            <a-space>
              <a-button type="primary" @click="formatJson">格式化</a-button>
              <a-button @click="compressJson">压缩</a-button>
              <a-button @click="clearJsonContent">清空</a-button>
              <a-button @click="copyJsonContent">复制</a-button>
              <a-button type="text" class="expand-btn" @click="expandAll">
                <template #icon>
                  <icon-expand />
                </template>
                展开全部
              </a-button>
              <a-button type="text" class="collapse-btn" @click="collapseAll">
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
              @input="handleJsonInput"
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
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconCopy, IconLock, IconUnlock, IconSafe, IconDown, IconRight, IconFile, IconExpand, IconShrink } from '@arco-design/web-vue/es/icon'
import CryptoJS from 'crypto-js'
import { sm4 } from 'sm-crypto'

// AES
const aesKey = ref('upbest@2019_best')
const isAesKeyCollapsed = ref(true)
const aesEncryptInput = ref('')
const aesEncryptResult = ref('')
const aesDecryptInput = ref('')
const aesDecryptResult = ref('')
const isAesEncryptCollapsed = ref(true)
const isAesDecryptCollapsed = ref(false)

// SM4
const sm4Key = ref('BC13D6BD076F612567EE1B145A85CCBC')
const isSm4KeyCollapsed = ref(true)
const sm4EncryptInput = ref('')
const sm4EncryptResult = ref('')
const sm4DecryptInput = ref('')
const sm4DecryptResult = ref('')
const isSm4EncryptCollapsed = ref(true)
const isSm4DecryptCollapsed = ref(false)

// AES 密钥和向量
const AES_KEY = 'upbest@2019_best'
const AES_IV = 'upbest@2019_best'

// JSON 格式化相关
const showJsonModal = ref(false)
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

const handleJsonInput = () => {
  try {
    if (!jsonContent.value.trim()) {
      jsonData.value = null
      return
    }
    const obj = JSON.parse(jsonContent.value)
    jsonData.value = obj
    expandAll()
  } catch (error) {
    // 解析错误时不更新树形结构
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

const clearJsonContent = () => {
  jsonContent.value = ''
  jsonData.value = null
  Message.success('已清空')
}

const copyJsonContent = async () => {
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

const closeJsonModal = () => {
  showJsonModal.value = false
  jsonContent.value = ''
  jsonData.value = null
}

// 切换展开/收起状态
const toggleAesKey = () => {
  isAesKeyCollapsed.value = !isAesKeyCollapsed.value
  if (!isAesKeyCollapsed.value) {
    isAesEncryptCollapsed.value = true
    isAesDecryptCollapsed.value = true
  }
}

const toggleAesEncrypt = () => {
  isAesEncryptCollapsed.value = !isAesEncryptCollapsed.value
  if (!isAesEncryptCollapsed.value) {
    isAesKeyCollapsed.value = true
    isAesDecryptCollapsed.value = true
  }
}

const toggleAesDecrypt = () => {
  isAesDecryptCollapsed.value = !isAesDecryptCollapsed.value
  if (!isAesDecryptCollapsed.value) {
    isAesKeyCollapsed.value = true
    isAesEncryptCollapsed.value = true
  }
}

const toggleSm4Key = () => {
  isSm4KeyCollapsed.value = !isSm4KeyCollapsed.value
  if (!isSm4KeyCollapsed.value) {
    isSm4EncryptCollapsed.value = true
    isSm4DecryptCollapsed.value = true
  }
}

const toggleSm4Encrypt = () => {
  isSm4EncryptCollapsed.value = !isSm4EncryptCollapsed.value
  if (!isSm4EncryptCollapsed.value) {
    isSm4KeyCollapsed.value = true
    isSm4DecryptCollapsed.value = true
  }
}

const toggleSm4Decrypt = () => {
  isSm4DecryptCollapsed.value = !isSm4DecryptCollapsed.value
  if (!isSm4DecryptCollapsed.value) {
    isSm4KeyCollapsed.value = true
    isSm4EncryptCollapsed.value = true
  }
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  if (!text) return

  try {
    // 首先尝试使用 navigator.clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      Message.success('复制成功')
      return
    }

    // 备选方案：创建临时文本区域
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()

    try {
      document.execCommand('copy')
      textArea.remove()
      Message.success('复制成功')
    } catch (error) {
      textArea.remove()
      Message.error('复制失败')
    }
  } catch (error) {
    Message.error('复制失败')
  }
}

// AES 加密
const handleAesEncrypt = () => {
  try {
    if (!aesEncryptInput.value) {
      Message.warning('请输入需要加密的内容')
      return
    }
    const key = CryptoJS.enc.Utf8.parse(aesKey.value)
    const iv = CryptoJS.enc.Utf8.parse(aesKey.value)

    const encrypted = CryptoJS.AES.encrypt(aesEncryptInput.value, key, {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7
    })

    aesEncryptResult.value = encrypted.toString()
    Message.success('加密成功')
  } catch (error) {
    Message.error('加密失败')
  }
}

// AES 解密
const handleAesDecrypt = () => {
  try {
    if (!aesDecryptInput.value) {
      Message.warning('请输入需要解密的内容')
      return
    }
    const key = CryptoJS.enc.Utf8.parse(aesKey.value)
    const iv = CryptoJS.enc.Utf8.parse(aesKey.value)

    const decrypted = CryptoJS.AES.decrypt(aesDecryptInput.value, key, {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7
    })

    aesDecryptResult.value = decrypted.toString(CryptoJS.enc.Utf8)
    Message.success('解密成功')

    // 尝试解析 JSON
    try {
      JSON.parse(aesDecryptResult.value)
      jsonContent.value = aesDecryptResult.value
      handleJsonInput()
      showJsonModal.value = true
    } catch (error) {
      // 不是有效的 JSON，不显示格式化弹窗
    }
  } catch (error) {
    console.error('解密失败:', error)
    Message.error('解密失败，请检查密文格式是否正确')
  }
}

// SM4 加密
const handleSm4Encrypt = () => {
  try {
    if (!sm4EncryptInput.value) {
      Message.warning('请输入需要加密的内容')
      return
    }
    if (!sm4Key.value) {
      Message.warning('请输入SM4加密key')
      return
    }

    const encrypted = sm4.encrypt(sm4EncryptInput.value, sm4Key.value)
    sm4EncryptResult.value = encrypted
    Message.success('加密成功')
  } catch (error) {
    Message.error('加密失败')
  }
}

// SM4 解密
const handleSm4Decrypt = () => {
  try {
    if (!sm4DecryptInput.value) {
      Message.warning('请输入需要解密的内容')
      return
    }
    if (!sm4Key.value) {
      Message.warning('请输入SM4加密key')
      return
    }

    const decrypted = sm4.decrypt(sm4DecryptInput.value, sm4Key.value)
    sm4DecryptResult.value = decrypted
    Message.success('解密成功')

    // 尝试解析 JSON
    try {
      JSON.parse(sm4DecryptResult.value)
      jsonContent.value = sm4DecryptResult.value
      handleJsonInput()
      showJsonModal.value = true
    } catch (error) {
      // 不是有效的 JSON，不显示格式化弹窗
    }
  } catch (error) {
    Message.error('解密失败')
  }
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
  padding: 12px;
}

.crypto-tool {
  min-height: 600px;
}

.section {
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  padding: 2px 0;
  margin-bottom: 4px;
}

.section-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
  display: flex;
  align-items: center;
  gap: 4px;
}

.button-group {
  margin: 12px 0;
}

.result-label {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 6px;
  color: var(--color-text-1);
}

.result-content {
  padding: 10px;
  background-color: var(--color-fill-2);
  border-radius: 4px;
  font-family: monospace;
  word-break: break-all;
  min-height: 120px;
  max-height: 200px;
  overflow-y: auto;
  position: relative;
  cursor: pointer;
  transition: background-color 0.2s;
}

.result-content:hover {
  background-color: var(--color-fill-3);
}

.copy-btn {
  position: absolute;
  right: 8px;
  top: 8px;
  opacity: 0.6;
}

.result-content:hover .copy-btn {
  opacity: 1;
}

:deep(.arco-collapse-item-content) {
  padding: 12px;
}

:deep(.arco-collapse-item-header) {
  display: none;
}

:deep(.arco-textarea-wrapper) {
  margin-bottom: 0;
}

:deep(.arco-input-wrapper) {
  margin-bottom: 0;
}

:deep(.arco-textarea) {
  min-height: 120px !important;
}

.arco-modal {
  width: 80%;
}

.json-modal {
  min-height: 600px;
}

.json-formatter {
  height: calc(90vh - 110px);
}

.toolbar {
  margin-bottom: 16px;
}

.content-area {
  height: calc(90vh - 180px);
  min-height: 500px;
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
  font-size: 14px;
  line-height: 1.6;
  padding: 12px;
}

.preview-area {
  height: 100%;
  overflow: auto;
  border: 1px solid var(--color-neutral-3);
  border-radius: 4px;
  padding: 16px;
}

.empty-state {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-3);
}

.empty-state :deep(svg) {
  font-size: 48px;
  margin-bottom: 16px;
}

.node-title {
  font-family: monospace;
  font-size: 14px;
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

</style>
