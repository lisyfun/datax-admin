<template>
  <a-modal
    v-model:visible="modelVisible"
    :title="isEdit ? '编辑任务' : '新建任务'"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form ref="formRef" :model="form" :rules="rules">
      <a-form-item field="name" label="任务名称" :rules="[{ required: true, message: '请输入任务名称' }]">
        <a-input v-model="form.name" placeholder="请输入任务名称" />
      </a-form-item>
      <a-form-item field="type" label="任务类型" :rules="[{ required: true, message: '请选择任务类型' }]">
        <a-radio-group v-model="form.type">
          <a-radio value="shell">Shell脚本</a-radio>
          <a-radio value="http">HTTP请求</a-radio>
          <a-radio value="datax">DataX任务</a-radio>
        </a-radio-group>
      </a-form-item>
      <template v-if="form.type === 'http'">
        <a-form-item field="url" label="请求URL" :rules="[{ required: true, message: '请输入请求URL' }, { type: 'url', message: '请输入有效的URL' }]">
          <a-input v-model="form.url" placeholder="请输入请求URL" allow-clear />
        </a-form-item>
        <a-form-item field="method" label="请求方法" :rules="[{ required: true, message: '请选择请求方法' }]">
          <a-select v-model="form.method" placeholder="请选择请求方法">
            <a-option value="GET">GET</a-option>
            <a-option value="POST">POST</a-option>
            <a-option value="PUT">PUT</a-option>
            <a-option value="DELETE">DELETE</a-option>
            <a-option value="PATCH">PATCH</a-option>
            <a-option value="HEAD">HEAD</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="headers" label="请求头">
          <a-textarea
            v-model="form.headers"
            placeholder="请输入请求头(JSON格式)"
            :auto-size="{ minRows: 2, maxRows: 4 }"
          />
        </a-form-item>
        <a-form-item field="body" label="请求体">
          <a-textarea
            v-model="form.body"
            placeholder="请输入请求体"
            :auto-size="{ minRows: 3, maxRows: 5 }"
          />
        </a-form-item>
        <a-form-item field="success_codes" label="成功状态码">
          <a-input
            v-model="form.success_codes"
            placeholder="请输入成功状态码，多个用逗号分隔，如：200,201,204"
          />
        </a-form-item>
      </template>

      <template v-if="form.type === 'shell'">
        <a-form-item
          field="command"
          label="执行命令"
          :rules="[{ required: true, message: '请输入执行命令' }]"
        >
          <a-textarea
            v-model="form.command"
            placeholder="请输入要执行的命令"
          />
        </a-form-item>
        <a-form-item field="working_dir" label="工作目录">
          <a-input v-model="form.working_dir" placeholder="请输入工作目录，默认为当前目录" />
        </a-form-item>
      </template>

      <template v-if="form.type === 'datax'">
        <a-form-item
          field="command"
          label="DataX配置"
          :rules="[{ required: true, message: '请输入DataX配置' }]"
        >
          <a-textarea
            v-model="form.command"
            placeholder="请输入DataX任务配置JSON"
            :auto-size="{ minRows: 10, maxRows: 20 }"
          />
          <template #extra>
            <div class="datax-tools">
              <a-space>
                <a-button type="text" @click="handleFormatJson">
                  <template #icon><icon-code /></template>
                  格式化JSON
                </a-button>
                <a-button type="text" @click="handleLoadTemplate">
                  <template #icon><icon-file /></template>
                  加载模板
                </a-button>
              </a-space>
            </div>
          </template>
        </a-form-item>
      </template>

      <a-form-item field="description" label="任务描述">
        <a-textarea v-model="form.description" placeholder="请输入任务描述" />
      </a-form-item>
      <a-form-item field="cron_expr" label="Cron 表达式" :rules="[{ required: true, message: '请输入 Cron 表达式' }]">
        <a-input v-model="form.cron_expr" placeholder="请输入 Cron 表达式">
          <template #append>
            <a-button @click="showCronGenerator">
              <template #icon><icon-edit /></template>
              生成表达式
            </a-button>
          </template>
        </a-input>
      </a-form-item>
      <a-form-item field="timeout" label="超时时间(秒)">
        <a-input-number
          v-model="form.timeout"
          placeholder="请输入超时时间"
          :min="0"
          :max="86400"
          :step="1"
        />
      </a-form-item>
      <a-form-item field="retry_count" label="重试次数">
        <a-input-number
          v-model="form.retry_count"
          placeholder="请输入重试次数"
          :min="0"
          :max="10"
          :step="1"
        />
      </a-form-item>
      <a-form-item field="retry_delay" label="重试间隔(秒)">
        <a-input-number
          v-model="form.retry_delay"
          placeholder="请输入重试间隔"
          :min="1"
          :max="3600"
          :step="1"
        />
      </a-form-item>
    </a-form>

    <!-- Cron表达式生成器组件 -->
    <cron-generator
      v-model:visible="showCronModal"
      v-model:expression="form.cron_expr"
    />
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from 'vue';
import { Message } from '@arco-design/web-vue';
import { IconCode, IconFile, IconEdit } from '@arco-design/web-vue/es/icon';
import type { Job, JobShellParams, JobHTTPParams, JobDataXParams } from '@/api/types';
import { createJob, updateJob } from '@/api/job';
import CronGenerator from './CronGenerator.vue';

const props = defineProps<{
  visible: boolean;
  isEdit: boolean;
  jobData?: Job;
}>();

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void;
  (e: 'success'): void;
}>();

const formRef = ref();
const showCronModal = ref(false);

interface FormState {
  id: number;
  name: string;
  description: string;
  type: 'shell' | 'http' | 'datax';
  command: string;
  working_dir: string;
  url: string;
  method: string;
  headers: string;
  body: string;
  success_codes: string;
  cron_expr: string;
  timeout: number;
  retry_count: number;
  retry_delay: number;
  params: Record<string, any>;
}

const form = reactive<FormState>({
  id: 0,
  name: '',
  description: '',
  type: 'shell',
  command: '',
  working_dir: '',
  url: '',
  method: 'GET',
  headers: '{}',
  body: '',
  success_codes: '200',
  cron_expr: '',
  timeout: 0,
  retry_count: 0,
  retry_delay: 0,
  params: {},
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入任务名称' }
  ],
  type: [
    { required: true, message: '请选择任务类型' }
  ],
  cron_expr: [
    { required: true, message: '请输入 Cron 表达式' }
  ]
};

// 格式化 JSON
const handleFormatJson = () => {
  try {
    const jsonObj = JSON.parse(form.command);
    form.command = JSON.stringify(jsonObj, null, 2);
  } catch (err) {
    Message.error('JSON 格式错误，无法格式化');
  }
};

// 加载模板
const handleLoadTemplate = () => {
  const template = {
    "job": {
      "setting": {
        "speed": {
          "channel": 3
        }
      },
      "content": [
        {
          "reader": {
            "name": "mysqlreader",
            "parameter": {
              "username": "root",
              "password": "password",
              "column": ["*"],
              "connection": [
                {
                  "table": ["table"],
                  "jdbcUrl": ["jdbc:mysql://localhost:3306/database"]
                }
              ]
            }
          },
          "writer": {
            "name": "mysqlwriter",
            "parameter": {
              "username": "root",
              "password": "password",
              "column": ["*"],
              "connection": [
                {
                  "table": ["table"],
                  "jdbcUrl": "jdbc:mysql://localhost:3306/database"
                }
              ]
            }
          }
        }
      ]
    }
  };
  form.command = JSON.stringify(template, null, 2);
};

// 显示Cron表达式生成器
const showCronGenerator = () => {
  showCronModal.value = true;
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();

    let params: JobShellParams | JobHTTPParams | JobDataXParams;
    if (form.type === 'shell') {
      params = {
        command: form.command,
        work_dir: form.working_dir,
        environment: {}
      } as JobShellParams;
    } else if (form.type === 'http') {
      params = {
        url: form.url,
        method: form.method,
        headers: JSON.parse(form.headers),
        body: form.body,
        success_code: form.success_codes.split(',').map(code => parseInt(code.trim())).filter(code => !isNaN(code))
      } as JobHTTPParams;
    } else {
      params = JSON.parse(form.command) as JobDataXParams;
    }

    const data = {
      name: form.name,
      description: form.description,
      type: form.type,
      cron_expr: form.cron_expr,
      timeout: form.timeout,
      retry_count: form.retry_count,
      retry_delay: form.retry_delay,
      params
    };

    if (props.isEdit) {
      await updateJob(form.id, data);
      Message.success('编辑任务成功');
    } else {
      await createJob(data);
      Message.success('创建任务成功');
    }

    emit('update:visible', false);
    emit('success');
  } catch (err: any) {
    if (err.response?.data?.error) {
      Message.error(err.response.data.error);
    } else if (err.errors) {
      Message.error('表单验证失败，请检查输入');
    } else {
      Message.error(props.isEdit ? '编辑任务失败' : '创建任务失败');
    }
  }
};

// 取消表单
const handleCancel = () => {
  emit('update:visible', false);
  if (formRef.value) {
    formRef.value.resetFields();
  }
};

// 监听jobData变化，更新表单数据
watch(() => props.jobData, (newVal: Job | undefined) => {
  if (newVal) {
    form.id = newVal.id;
    form.name = newVal.name;
    form.description = newVal.description || '';
    form.type = newVal.type as 'shell' | 'http' | 'datax';

    try {
      const params = typeof newVal.params === 'string' ? JSON.parse(newVal.params) : newVal.params;
      if (newVal.type === 'shell') {
        form.command = params.command || '';
        form.working_dir = params.work_dir || '';
      } else if (newVal.type === 'http') {
        form.url = params.url || '';
        form.method = params.method || 'GET';
        form.headers = JSON.stringify(params.headers || {}, null, 2);
        form.body = params.body || '';
        form.success_codes = (params.success_code || [200]).join(',');
      } else {
        form.command = typeof params === 'string' ? params : JSON.stringify(params, null, 2);
      }
    } catch (e) {
      console.error('解析参数失败:', e);
      form.command = '';
      form.working_dir = '';
      form.url = '';
      form.method = 'GET';
      form.headers = '{}';
      form.body = '';
      form.success_codes = '200';
    }

    form.cron_expr = newVal.cron_expr;
    form.timeout = newVal.timeout || 0;
    form.retry_count = newVal.retry_count || 0;
    form.retry_delay = newVal.retry_delay || 0;
    form.params = {};
  }
}, { immediate: true });

// 处理visible的双向绑定
const modelVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
});
</script>

<style scoped>
.datax-tools {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed var(--color-neutral-3);
}

.datax-tools :deep(.arco-btn) {
  padding: 0 4px;
  font-size: 14px;
  color: var(--color-text-3);
}

.datax-tools :deep(.arco-btn:hover) {
  color: rgb(var(--primary-6));
}
</style>
