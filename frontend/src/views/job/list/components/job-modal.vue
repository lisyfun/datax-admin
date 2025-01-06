<template>
  <a-modal
    v-model:visible="visible"
    :title="title"
    @before-ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-align="right"
      :label-col-props="{ span: 6 }"
      :wrapper-col-props="{ span: 18 }"
      @submit="handleSubmit"
    >
      <a-form-item field="name" label="任务名称" :rules="[{ required: true, message: '请输入任务名称' }, { maxLength: 100, message: '任务名称最大长度为100个字符' }]">
        <a-input v-model="form.name" placeholder="请输入任务名称" allow-clear />
      </a-form-item>
      <a-form-item field="type" label="任务类型" :rules="[{ required: true, message: '请选择任务类型' }]">
        <a-select
          v-model="form.type"
          :options="jobTypes"
          placeholder="请选择任务类型"
          @change="handleTypeChange"
        />
      </a-form-item>
      <a-form-item field="description" label="任务描述" :rules="[{ maxLength: 500, message: '任务描述最大长度为500个字符' }]">
        <a-textarea
          v-model="form.description"
          placeholder="请输入任务描述"
          allow-clear
        />
      </a-form-item>
      <a-form-item field="cron_expr" label="Cron表达式" :rules="[{ required: true, message: '请输入Cron表达式' }, { maxLength: 100, message: 'Cron表达式最大长度为100个字符' }]">
        <a-input-group>
          <a-input v-model="form.cron_expr" placeholder="请输入Cron表达式" allow-clear style="width: calc(100% - 100px)" />
          <a-button type="primary" @click="showCronHelper">Cron助手</a-button>
        </a-input-group>
      </a-form-item>
      <a-form-item field="timeout" label="超时时间(秒)" :rules="[{ type: 'number', min: 0, message: '超时时间必须大于等于0' }]">
        <a-input-number
          v-model="form.timeout"
          :min="0"
          placeholder="请输入超时时间"
          mode="button"
        />
      </a-form-item>
      <a-form-item field="retry_count" label="重试次数" :rules="[{ type: 'number', min: 0, message: '重试次数必须大于等于0' }]">
        <a-input-number
          v-model="form.retry_count"
          :min="0"
          placeholder="请输入重试次数"
          mode="button"
        />
      </a-form-item>
      <a-form-item field="retry_delay" label="重试间隔(秒)" :rules="[{ type: 'number', min: 0, message: '重试间隔必须大于等于0' }]">
        <a-input-number
          v-model="form.retry_delay"
          :min="0"
          placeholder="请输入重试间隔"
          mode="button"
        />
      </a-form-item>

      <!-- Shell任务参数 -->
      <template v-if="form.type === 'shell'">
        <a-form-item
          field="params.command"
          label="执行命令"
          :rules="[{ required: true, message: '请输入执行命令' }]"
        >
          <a-textarea
            v-model="(form.params as JobShellParams).command"
            placeholder="请输入执行命令"
            :auto-size="{ minRows: 3, maxRows: 5 }"
          />
        </a-form-item>
        <a-form-item field="params.work_dir" label="工作目录">
          <a-input
            v-model="(form.params as JobShellParams).work_dir"
            placeholder="请输入工作目录"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="params.environment" label="环境变量">
          <a-textarea
            v-model="environmentText"
            placeholder="请输入环境变量(JSON格式)"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @blur="validateEnvironment"
          />
        </a-form-item>
      </template>

      <!-- HTTP任务参数 -->
      <template v-if="form.type === 'http'">
        <a-form-item
          field="params.url"
          label="请求URL"
          :rules="[{ required: true, message: '请输入请求URL' }, { type: 'url', message: '请输入有效的URL' }]"
        >
          <a-input v-model="(form.params as JobHTTPParams).url" placeholder="请输入请求URL" allow-clear />
        </a-form-item>
        <a-form-item field="params.method" label="请求方法" :rules="[{ required: true, message: '请选择请求方法' }]">
          <a-select
            v-model="(form.params as JobHTTPParams).method"
            :options="httpMethods"
            placeholder="请选择请求方法"
          />
        </a-form-item>
        <a-form-item field="params.headers" label="请求头">
          <a-textarea
            v-model="headersText"
            placeholder="请输入请求头(JSON格式)"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @blur="validateHeaders"
          />
        </a-form-item>
        <a-form-item field="params.body" label="请求体">
          <a-textarea
            v-model="(form.params as JobHTTPParams).body"
            placeholder="请输入请求体"
            :auto-size="{ minRows: 3, maxRows: 5 }"
          />
        </a-form-item>
        <a-form-item field="params.success_code" label="成功状态码">
          <a-input
            v-model="successCodeText"
            placeholder="请输入成功状态码(多个用逗号分隔)"
            @blur="validateSuccessCode"
          />
        </a-form-item>
      </template>

      <!-- DataX任务参数 -->
      <template v-if="form.type === 'datax'">
        <a-form-item
          field="params.job_path"
          label="任务文件路径"
          :rules="[{ required: true, message: '请输入任务文件路径' }]"
        >
          <a-input
            v-model="(form.params as JobDataXParams).job_path"
            placeholder="请输入任务文件路径"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="params.parameters" label="任务参数">
          <a-textarea
            v-model="parametersText"
            placeholder="请输入任务参数(JSON格式)"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @blur="validateParameters"
          />
        </a-form-item>
        <a-form-item field="params.jvm_options" label="JVM参数">
          <a-textarea
            v-model="jvmOptionsText"
            placeholder="请输入JVM参数(每行一个)"
            :auto-size="{ minRows: 2, maxRows: 4 }"
            @blur="validateJvmOptions"
          />
        </a-form-item>
        <a-form-item field="params.speed" label="速率限制" :rules="[{ type: 'number', min: 0, message: '速率限制必须大于等于0' }]">
          <a-input-number
            v-model="(form.params as JobDataXParams).speed"
            :min="0"
            placeholder="请输入速率限制"
            mode="button"
          />
        </a-form-item>
      </template>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FormInstance } from '@arco-design/web-vue';
import type {
  JobType,
  JobShellParams,
  JobHTTPParams,
  JobDataXParams,
  CreateJobRequest,
  UpdateJobRequest,
} from '@/api/types';
import { createJob, updateJob } from '@/api/job';

interface FormState {
  id?: number;
  name: string;
  type: JobType | undefined;
  description: string;
  cron_expr: string;
  timeout: number;
  retry_count: number;
  retry_delay: number;
  params: JobShellParams | JobHTTPParams | JobDataXParams;
}

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  type: {
    type: String as () => 'create' | 'edit',
    default: 'create',
  },
  data: {
    type: Object as () => Partial<FormState>,
    default: () => ({}),
  },
});

const emit = defineEmits(['update:visible', 'success']);
const formRef = ref<FormInstance>();

const title = computed(() =>
  props.type === 'create' ? '创建任务' : '编辑任务'
);

const jobTypes = computed(() => [
  { label: 'Shell脚本', value: 'shell' as const },
  { label: 'HTTP请求', value: 'http' as const },
  { label: 'DataX同步', value: 'datax' as const },
]);

const httpMethods = computed(() => [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
  { label: 'HEAD', value: 'HEAD' },
]);

const generateFormModel = (): FormState => ({
  name: '',
  type: undefined,
  description: '',
  cron_expr: '',
  timeout: 3600, // 默认1小时超时
  retry_count: 3, // 默认重试3次
  retry_delay: 60, // 默认重试间隔1分钟
  params: {} as JobShellParams | JobHTTPParams | JobDataXParams,
});

const form = ref<FormState>(generateFormModel());
const environmentText = ref('');
const headersText = ref('');
const successCodeText = ref('');
const parametersText = ref('');
const jvmOptionsText = ref('');

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入任务名称' },
    { maxLength: 100, message: '任务名称最大长度为100个字符' },
  ],
  type: [{ required: true, message: '请选择任务类型' }],
  description: [
    { maxLength: 500, message: '任务描述最大长度为500个字符' },
  ],
  cron_expr: [
    { required: true, message: '请输入Cron表达式' },
    { maxLength: 100, message: 'Cron表达式最大长度为100个字符' },
  ],
  timeout: [
    { type: 'number' as const, min: 0, message: '超时时间必须大于等于0' },
  ],
  retry_count: [
    { type: 'number' as const, min: 0, message: '重试次数必须大于等于0' },
  ],
  retry_delay: [
    { type: 'number' as const, min: 0, message: '重试间隔必须大于等于0' },
  ],
};

// 处理任务类型变更
const handleTypeChange = (value: unknown) => {
  const jobType = value as JobType;
  switch (jobType) {
    case 'shell':
      form.value.params = {
        command: '',
        work_dir: '',
        environment: {},
      } as JobShellParams;
      break;
    case 'http':
      form.value.params = {
        url: '',
        method: 'GET',
        headers: {},
        body: '',
        success_code: [200],
      } as JobHTTPParams;
      break;
    case 'datax':
      form.value.params = {
        job_path: '',
        parameters: {},
        jvm_options: [],
        speed: 0,
      } as JobDataXParams;
      break;
  }
};

// 验证JSON格式的输入
const validateJSON = (text: string, fieldName: string): Record<string, string> | null => {
  try {
    if (!text) return {};
    const obj = JSON.parse(text);
    if (typeof obj !== 'object' || Array.isArray(obj)) {
      Message.error(`${fieldName}必须是一个JSON对象`);
      return null;
    }
    return obj;
  } catch (e) {
    Message.error(`${fieldName}格式不正确，请输入有效的JSON格式`);
    return null;
  }
};

// 验证环境变量
const validateEnvironment = () => {
  const env = validateJSON(environmentText.value, '环境变量');
  if (env) {
    (form.value.params as JobShellParams).environment = env;
  }
};

// 验证请求头
const validateHeaders = () => {
  const headers = validateJSON(headersText.value, '请求头');
  if (headers) {
    (form.value.params as JobHTTPParams).headers = headers;
  }
};

// 验证成功状态码
const validateSuccessCode = () => {
  try {
    const codes = successCodeText.value
      .split(',')
      .map(code => parseInt(code.trim()))
      .filter(code => !isNaN(code) && code >= 100 && code < 600);

    if (codes.length === 0) {
      codes.push(200); // 默认使用200
    }
    (form.value.params as JobHTTPParams).success_code = codes;
  } catch (e) {
    Message.error('成功状态码格式不正确');
  }
};

// 验证任务参数
const validateParameters = () => {
  const params = validateJSON(parametersText.value, '任务参数');
  if (params) {
    (form.value.params as JobDataXParams).parameters = params;
  }
};

// 验证JVM参数
const validateJvmOptions = () => {
  const options = jvmOptionsText.value
    .split('\n')
    .map(opt => opt.trim())
    .filter(opt => opt);
  (form.value.params as JobDataXParams).jvm_options = options;
};

// 显示Cron表达式助手
const showCronHelper = () => {
  // TODO: 实现Cron表达式助手
  Message.info('Cron表达式助手功能开发中...');
};

// 处理表单提交
const handleSubmit = async () => {
  try {
    await formRef.value?.validate();

    // 构造基础参数
    const params: CreateJobRequest = {
      name: form.value.name,
      type: form.value.type!,
      description: form.value.description || '',
      cron_expr: form.value.cron_expr,
      timeout: form.value.timeout || 0,
      retry_count: form.value.retry_count || 0,
      retry_delay: form.value.retry_delay || 0,
      params: {} as any, // 先初始化一个空对象
    };

    // 根据任务类型设置params
    if (form.value.type === 'shell') {
      params.params = {
        command: (form.value.params as JobShellParams).command,
        work_dir: (form.value.params as JobShellParams).work_dir || '',
        environment: (form.value.params as JobShellParams).environment || {},
      } as JobShellParams;
    } else if (form.value.type === 'http') {
      params.params = {
        url: (form.value.params as JobHTTPParams).url,
        method: (form.value.params as JobHTTPParams).method,
        headers: (form.value.params as JobHTTPParams).headers || {},
        body: (form.value.params as JobHTTPParams).body || '',
        success_code: (form.value.params as JobHTTPParams).success_code || [200],
      } as JobHTTPParams;
    } else if (form.value.type === 'datax') {
      params.params = {
        job_path: (form.value.params as JobDataXParams).job_path,
        parameters: (form.value.params as JobDataXParams).parameters || {},
        jvm_options: (form.value.params as JobDataXParams).jvm_options || [],
        speed: (form.value.params as JobDataXParams).speed || 0,
      } as JobDataXParams;
    }

    console.log(params);

    if (props.type === 'create') {
      await createJob(params);
      Message.success('任务创建成功');
    } else {
      await updateJob(form.value.id!, params);
      Message.success('任务更新成功');
    }

    emit('success');
    emit('update:visible', false);
  } catch (err) {
    // 表单验证失败或API调用失败
    console.error('提交失败:', err);
    if (err instanceof Error) {
      Message.error(err.message);
    } else {
      Message.error('提交失败，请检查表单内容');
    }
  }
};

// 处理取消
const handleCancel = () => {
  formRef.value?.resetFields();
  emit('update:visible', false);
};

// 监听visible变化，重置表单
watch(
  () => props.visible,
  (val) => {
    if (val && props.type === 'edit' && props.data) {
      form.value = { ...generateFormModel(), ...props.data };
      // 初始化文本框的值
      if (form.value.type === 'shell') {
        environmentText.value = JSON.stringify(
          (form.value.params as JobShellParams).environment,
          null,
          2
        );
      } else if (form.value.type === 'http') {
        headersText.value = JSON.stringify(
          (form.value.params as JobHTTPParams).headers,
          null,
          2
        );
        successCodeText.value = (form.value.params as JobHTTPParams).success_code.join(',');
      } else if (form.value.type === 'datax') {
        parametersText.value = JSON.stringify(
          (form.value.params as JobDataXParams).parameters,
          null,
          2
        );
        jvmOptionsText.value = (form.value.params as JobDataXParams).jvm_options.join('\n');
      }
    } else if (!val) {
      form.value = generateFormModel();
      environmentText.value = '';
      headersText.value = '';
      successCodeText.value = '';
      parametersText.value = '';
      jvmOptionsText.value = '';
    }
  }
);
</script>

<style scoped>
.arco-form-item {
  margin-bottom: 20px;
}
</style>
