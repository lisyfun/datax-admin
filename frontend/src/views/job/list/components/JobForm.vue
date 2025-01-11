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
        <a-form-item label="任务参数">
          <div class="datax-params">
            <div v-for="(param, index) in form.datax_params.parameters" :key="index" class="parameter-item">
              <a-space :size="8" fill>
                <a-input
                  v-model="param.key"
                  placeholder="参数名"
                  allow-clear
                />
                <a-input
                  v-model="param.value"
                  placeholder="参数值"
                  allow-clear
                />
                <a-button
                  type="text"
                  status="danger"
                  @click="removeParameter(index)"
                >
                  <template #icon><icon-delete /></template>
                </a-button>
              </a-space>
            </div>
            <div class="parameter-add">
              <a-button type="dashed" long @click="addParameter">
                <template #icon><icon-plus /></template>
                添加参数
              </a-button>
            </div>
          </div>
        </a-form-item>

        <a-form-item
          label="任务内容"
          field="datax_params.job_content"
          :rules="[{
            required: true,
            validator: (value: string) => {
              try {
                const content = JSON.parse(value || '{}');
                const reader = content.job?.content?.[0]?.reader;
                const writer = content.job?.content?.[0]?.writer;
                if (!reader || !writer) {
                  return false;
                }
                return true;
              } catch {
                return false;
              }
            },
            message: '请配置DataX任务的Reader和Writer'
          }]"
        >
          <div class="datax-buttons">
            <a-space>
              <a-button type="primary" @click="handleConfigReader">
                <template #icon><icon-edit /></template>
                Reader配置
              </a-button>
              <a-button type="primary" @click="handleConfigWriter">
                <template #icon><icon-edit /></template>
                Writer配置
              </a-button>
            </a-space>
          </div>
        </a-form-item>

        <!-- Reader配置弹窗 -->
        <a-modal
          v-model:visible="showReaderModal"
          title="Reader配置"
          @ok="handleReaderModalOk"
          @cancel="handleReaderModalCancel"
          :width="600"
        >
          <ReaderForm
            v-if="currentReader"
            v-model="currentReader"
            @update:model-value="handleReaderUpdate"
          />
        </a-modal>

        <!-- Writer配置弹窗 -->
        <a-modal
          v-model:visible="showWriterModal"
          title="Writer配置"
          @ok="handleWriterModalOk"
          @cancel="handleWriterModalCancel"
          :width="600"
        >
          <WriterForm
            v-if="currentWriter"
            v-model="currentWriter"
            @update:model-value="handleWriterUpdate"
          />
        </a-modal>
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
import { Message, Modal } from '@arco-design/web-vue';
import {
  IconDelete,
  IconPlus,
  IconEdit,
  IconCode,
  IconFile
} from '@arco-design/web-vue/es/icon';
import type { Job, JobShellParams, JobHTTPParams, JobDataXParams } from '@/api/types';
import { createJob, updateJob } from '@/api/job';
import CronGenerator from './CronGenerator.vue';
import ReaderForm from './ReaderForm.vue';
import WriterForm from './WriterForm.vue';

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

// 状态变量
const showReaderModal = ref(false);
const showWriterModal = ref(false);
const currentReader = ref<any>(null);
const currentWriter = ref<any>(null);

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
  datax_params: {
    parameters: Array<{ key: string; value: string }>;
    job_content: string;
  };
}

// 初始化表单数据
const initFormState = (): FormState => ({
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
  datax_params: {
    parameters: [],
    job_content: '',
  },
});

const form = reactive<FormState>(initFormState());

// 重置表单
const resetForm = () => {
  // 只在新建时重置表单
  if (!props.isEdit) {
    Object.assign(form, initFormState());
    if (formRef.value) {
      formRef.value.resetFields();
    }
  }
};

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
  ],
  'datax_params.job_content': [{
    required: true,
    validator: (value: string) => {
      try {
        const content = JSON.parse(value || '{}');
        const reader = content.job?.content?.[0]?.reader;
        const writer = content.job?.content?.[0]?.writer;
        if (!reader || !writer) {
          return false;
        }
        return true;
      } catch {
        return false;
      }
    },
    message: '请配置DataX任务的Reader和Writer'
  }]
};

// 格式化 JSON
const handleFormatJson = () => {
  try {
    const content = JSON.parse(form.datax_params.job_content);
    form.datax_params.job_content = JSON.stringify(content, null, 2);
  } catch (err) {
    Message.error('JSON 格式错误，无法格式化');
  }
};

// 加载模板
const handleLoadTemplate = () => {
  const template = {
    job: {
      content: [{
        reader: {
          name: 'mysqlreader',
          parameter: {
            username: 'root',
            password: 'password',
            column: ['*'],
            connection: [
              {
                table: ['table'],
                jdbcUrl: ['jdbc:mysql://localhost:3306/database']
              }
            ]
          }
        },
        writer: {
          name: 'mysqlwriter',
          parameter: {
            username: 'root',
            password: 'password',
            column: ['*'],
            connection: [
              {
                table: ['table'],
                jdbcUrl: 'jdbc:mysql://localhost:3306/database'
              }
            ]
          }
        }
      }],
      setting: {
        speed: {
          channel: 24,
          bytes: 52428800
        },
        errorLimit: {
          record: 0,
          percentage: 0.02
        }
      }
    }
  };
  form.datax_params.job_content = JSON.stringify(template, null, 2);
};

// 显示Cron表达式生成器
const showCronGenerator = () => {
  showCronModal.value = true;
};

// 添加参数
const addParameter = () => {
  form.datax_params.parameters.push({ key: '', value: '' });
};

// 删除参数
const removeParameter = (index: number) => {
  form.datax_params.parameters.splice(index, 1);
};

// 处理JSON更新
const handleJsonUpdate = ({ value }: { path: string[], value: any }) => {
  try {
    form.datax_params.job_content = JSON.stringify(value, null, 2);
  } catch (err) {
    Message.error('JSON格式错误');
  }
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    console.log('提交表单时的表单数据:', form);

    let params: JobShellParams | JobHTTPParams | JobDataXParams;
    if (form.type === 'shell') {
      params = {
        command: form.command.trim(),
        work_dir: form.working_dir.trim(),
        environment: {}
      } as JobShellParams;
    } else if (form.type === 'http') {
      params = {
        url: form.url.trim(),
        method: form.method,
        headers: JSON.parse(form.headers),
        body: form.body,
        success_code: form.success_codes.split(',').map(code => parseInt(code.trim())).filter(code => !isNaN(code))
      } as JobHTTPParams;
    } else {
      // 处理DataX参数
      const parameters: Record<string, string> = {};
      form.datax_params.parameters.forEach(({ key, value }) => {
        if (key && value) {
          parameters[key] = value;
        }
      });

      try {
        // 解析当前的 job_content
        const currentContent = JSON.parse(form.datax_params.job_content || '{}');

        // 构造完整的 DataX 参数
        const jobConfig = {
          job: {
            content: [{
              reader: currentContent.job?.content?.[0]?.reader || {},
              writer: currentContent.job?.content?.[0]?.writer || {}
            }],
            setting: {
              speed: {
                channel: 24,
                bytes: 52428800
              },
              errorLimit: {
                record: 0,
                percentage: 0.02
              }
            }
          }
        };

        params = {
          job_config: JSON.stringify(jobConfig),
          parameters
        } as unknown as JobDataXParams;
      } catch (err) {
        console.error('构造DataX参数失败:', err);
        Message.error('构造DataX参数失败');
        return;
      }
    }

    const data = {
      name: form.name.trim(),
      description: form.description.trim(),
      type: form.type,
      cron_expr: form.cron_expr.trim(),
      timeout: form.timeout,
      retry_count: form.retry_count,
      retry_delay: form.retry_delay,
      params
    };

    console.log('提交的数据:', {
      id: props.jobData?.id,
      data
    });

    if (props.isEdit) {
      if (!props.jobData?.id) {
        Message.error('任务ID无效');
        return;
      }
      console.log('正在更新任务，ID:', props.jobData.id);
      await updateJob(props.jobData.id, data);
      Message.success('编辑任务成功');
    } else {
      console.log('正在创建新任务');
      await createJob(data);
      Message.success('创建任务成功');
    }

    emit('update:visible', false);
    emit('success');
  } catch (err: any) {
    console.error('表单提交失败:', err);
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
};

// 监听visible变化，当弹窗关闭时重置表单
watch(() => props.visible, (newVal) => {
  if (!newVal) {
    // 关闭弹窗时，只在非编辑模式下重置表单
    if (!props.isEdit) {
      Object.assign(form, initFormState());
      if (formRef.value) {
        formRef.value.resetFields();
      }
    }
  }
});

// 监听jobData变化，更新表单数据
watch(() => props.jobData, (newVal: Job | undefined) => {
  if (newVal && props.visible) {
    console.log('编辑任务数据:', newVal);

    if (!newVal.id) {
      console.error('任务数据中缺少ID');
      Message.error('任务数据无效');
      return;
    }

    try {
      // 先解析参数
      let params: any;
      if (typeof newVal.params === 'string') {
        params = JSON.parse(newVal.params);
        console.log('解析后的参数:', params);
      } else {
        params = newVal.params;
      }

      // 根据任务类型设置表单数据
      if (newVal.type === 'shell') {
        Object.assign(form, {
          id: newVal.id,
          name: newVal.name || '',
          description: newVal.description || '',
          type: newVal.type,
          cron_expr: newVal.cron_expr || '',
          timeout: newVal.timeout || 0,
          retry_count: newVal.retry_count || 0,
          retry_delay: newVal.retry_delay || 0,
          command: params.command || '',
          working_dir: params.work_dir || '',
        });
      } else if (newVal.type === 'http') {
        Object.assign(form, {
          id: newVal.id,
          name: newVal.name || '',
          description: newVal.description || '',
          type: newVal.type,
          cron_expr: newVal.cron_expr || '',
          timeout: newVal.timeout || 0,
          retry_count: newVal.retry_count || 0,
          retry_delay: newVal.retry_delay || 0,
          url: params.url || '',
          method: params.method || 'GET',
          headers: typeof params.headers === 'string' ? params.headers : JSON.stringify(params.headers || {}, null, 2),
          body: params.body || '',
          success_codes: Array.isArray(params.success_code) ? params.success_code.join(',') : '200'
        });
      } else if (newVal.type === 'datax') {
        try {
          // 解析DataX参数
          const dataxParams = typeof params === 'string' ? JSON.parse(params) : params;
          console.log('DataX原始参数:', dataxParams);

          // 处理job_config
          let jobConfig = dataxParams.job_config;
          if (typeof jobConfig === 'string') {
            jobConfig = JSON.parse(jobConfig);
          }
          console.log('解析后的job_config:', jobConfig);

          Object.assign(form, {
            id: newVal.id,
            name: newVal.name || '',
            description: newVal.description || '',
            type: newVal.type,
            cron_expr: newVal.cron_expr || '',
            timeout: newVal.timeout || 0,
            retry_count: newVal.retry_count || 0,
            retry_delay: newVal.retry_delay || 0,
            datax_params: {
              parameters: Object.entries(dataxParams.parameters || {}).map(([key, value]) => ({
                key,
                value: String(value)
              })),
              job_content: JSON.stringify(jobConfig, null, 2)
            }
          });
        } catch (e) {
          console.error('解析DataX参数失败:', e);
          Message.error(`解析DataX参数失败: ${(e as Error).message}`);
        }
      }

      console.log('更新后的表单数据:', form);
    } catch (e) {
      console.error('解析任务参数失败:', e);
      Message.error(`解析任务参数失败: ${(e as Error).message}`);
    }
  }
}, { immediate: true });

// 处理visible的双向绑定
const modelVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
});

// 配置Reader
const handleConfigReader = () => {
  try {
    const content = JSON.parse(form.datax_params.job_content || '{}');
    console.log('Reader配置 - 当前内容:', content);

    // 从 job.content[0] 中获取 reader 配置
    const reader = content.job?.content?.[0]?.reader || {
      name: 'mysqlreader',
      parameter: {
        username: '',
        password: '',
        host: 'localhost',
        port: 3306,
        database: '',
        table: '',
        columns: [],
        where: '1=1',
        batchSize: 20000
      }
    };

    currentReader.value = reader;
    console.log('Reader配置 - 解析结果:', currentReader.value);
    showReaderModal.value = true;
  } catch (err) {
    console.error('Reader配置解析失败:', err);
    Message.error('JSON格式错误，请先输入正确的JSON');
  }
};

// 配置Writer
const handleConfigWriter = () => {
  try {
    const content = JSON.parse(form.datax_params.job_content || '{}');
    console.log('Writer配置 - 当前内容:', content);

    // 从 job.content[0] 中获取 writer 配置
    const writer = content.job?.content?.[0]?.writer || {
      name: 'mysqlwriter',
      parameter: {
        username: '',
        password: '',
        host: 'localhost',
        port: 3306,
        database: '',
        table: '',
        columns: [],
        writeMode: 'insert',
        batchSize: 10000,
        preSql: [],
        postSql: []
      }
    };

    currentWriter.value = writer;
    console.log('Writer配置 - 解析结果:', currentWriter.value);
    showWriterModal.value = true;
  } catch (err) {
    console.error('Writer配置解析失败:', err);
    Message.error('JSON格式错误，请先输入正确的JSON');
  }
};

// 处理Reader配置更新
const handleReaderUpdate = (val: any) => {
  currentReader.value = val;
};

// 处理Writer配置更新
const handleWriterUpdate = (val: any) => {
  currentWriter.value = val;
};

// 处理Reader配置对话框确认
const handleReaderModalOk = () => {
  try {
    const content = JSON.parse(form.datax_params.job_content || '{}');
    if (!content.job) {
      content.job = {
        content: [{}],
        setting: {
          speed: {
            channel: 24,
            bytes: 52428800
          },
          errorLimit: {
            record: 0,
            percentage: 0.02
          }
        }
      };
    }
    if (!content.job.content) {
      content.job.content = [{}];
    }
    if (!content.job.content[0]) {
      content.job.content[0] = {};
    }
    content.job.content[0].reader = currentReader.value;
    form.datax_params.job_content = JSON.stringify(content, null, 2);
    console.log('Reader配置保存后:', content);
    showReaderModal.value = false;
  } catch (err) {
    console.error('Reader配置保存失败:', err);
    Message.error('保存Reader配置失败');
  }
};

// 处理Writer配置对话框确认
const handleWriterModalOk = () => {
  try {
    const content = JSON.parse(form.datax_params.job_content || '{}');
    if (!content.job) {
      content.job = {
        content: [{}],
        setting: {
          speed: {
            channel: 24,
            bytes: 52428800
          },
          errorLimit: {
            record: 0,
            percentage: 0.02
          }
        }
      };
    }
    if (!content.job.content) {
      content.job.content = [{}];
    }
    if (!content.job.content[0]) {
      content.job.content[0] = {};
    }
    content.job.content[0].writer = currentWriter.value;
    form.datax_params.job_content = JSON.stringify(content, null, 2);
    console.log('Writer配置保存后:', content);
    showWriterModal.value = false;
  } catch (err) {
    console.error('Writer配置保存失败:', err);
    Message.error('保存Writer配置失败');
  }
};

// 处理Reader配置对话框取消
const handleReaderModalCancel = () => {
  showReaderModal.value = false;
};

// 处理Writer配置对话框取消
const handleWriterModalCancel = () => {
  showWriterModal.value = false;
};
</script>

<style scoped>
.datax-params {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.parameter-item {
  display: flex;
  align-items: center;
}

.parameter-item :deep(.arco-space-fill) {
  width: 100%;
}

.parameter-item :deep(.arco-input-wrapper) {
  flex: 1;
}

.parameter-add {
  margin-top: 4px;
}

.parameter-add :deep(.arco-btn) {
  border-style: dashed;
}

.datax-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
}
</style>
