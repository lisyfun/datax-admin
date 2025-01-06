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
      <a-form-item field="name" :label="$t('jobForm.name')" :rules="[{ required: true }]">
        <a-input v-model="form.name" :placeholder="$t('jobForm.name.placeholder')" />
      </a-form-item>
      <a-form-item field="type" :label="$t('jobForm.type')" :rules="[{ required: true }]">
        <a-select
          v-model="form.type"
          :options="jobTypes"
          :placeholder="$t('jobForm.type.placeholder')"
          @change="handleTypeChange"
        />
      </a-form-item>
      <a-form-item field="description" :label="$t('jobForm.description')">
        <a-textarea
          v-model="form.description"
          :placeholder="$t('jobForm.description.placeholder')"
        />
      </a-form-item>
      <a-form-item field="cron_expr" :label="$t('jobForm.cronExpr')" :rules="[{ required: true }]">
        <a-input v-model="form.cron_expr" :placeholder="$t('jobForm.cronExpr.placeholder')" />
      </a-form-item>
      <a-form-item field="timeout" :label="$t('jobForm.timeout')">
        <a-input-number
          v-model="form.timeout"
          :min="0"
          :placeholder="$t('jobForm.timeout.placeholder')"
        />
      </a-form-item>
      <a-form-item field="retry_count" :label="$t('jobForm.retryCount')">
        <a-input-number
          v-model="form.retry_count"
          :min="0"
          :placeholder="$t('jobForm.retryCount.placeholder')"
        />
      </a-form-item>
      <a-form-item field="retry_delay" :label="$t('jobForm.retryDelay')">
        <a-input-number
          v-model="form.retry_delay"
          :min="0"
          :placeholder="$t('jobForm.retryDelay.placeholder')"
        />
      </a-form-item>

      <!-- Shell任务参数 -->
      <template v-if="form.type === 'shell'">
        <a-form-item
          field="params.command"
          :label="$t('jobForm.shell.command')"
          :rules="[{ required: true }]"
        >
          <a-textarea
            v-model="form.params.command"
            :placeholder="$t('jobForm.shell.command.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.work_dir" :label="$t('jobForm.shell.workDir')">
          <a-input
            v-model="form.params.work_dir"
            :placeholder="$t('jobForm.shell.workDir.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.environment" :label="$t('jobForm.shell.environment')">
          <a-textarea
            v-model="environmentText"
            :placeholder="$t('jobForm.shell.environment.placeholder')"
          />
        </a-form-item>
      </template>

      <!-- HTTP任务参数 -->
      <template v-if="form.type === 'http'">
        <a-form-item
          field="params.url"
          :label="$t('jobForm.http.url')"
          :rules="[{ required: true }]"
        >
          <a-input v-model="form.params.url" :placeholder="$t('jobForm.http.url.placeholder')" />
        </a-form-item>
        <a-form-item field="params.method" :label="$t('jobForm.http.method')">
          <a-select
            v-model="form.params.method"
            :options="httpMethods"
            :placeholder="$t('jobForm.http.method.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.headers" :label="$t('jobForm.http.headers')">
          <a-textarea
            v-model="headersText"
            :placeholder="$t('jobForm.http.headers.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.body" :label="$t('jobForm.http.body')">
          <a-textarea
            v-model="form.params.body"
            :placeholder="$t('jobForm.http.body.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.success_code" :label="$t('jobForm.http.successCode')">
          <a-input
            v-model="successCodeText"
            :placeholder="$t('jobForm.http.successCode.placeholder')"
          />
        </a-form-item>
      </template>

      <!-- DataX任务参数 -->
      <template v-if="form.type === 'datax'">
        <a-form-item
          field="params.job_path"
          :label="$t('jobForm.datax.jobPath')"
          :rules="[{ required: true }]"
        >
          <a-input
            v-model="form.params.job_path"
            :placeholder="$t('jobForm.datax.jobPath.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.parameters" :label="$t('jobForm.datax.parameters')">
          <a-textarea
            v-model="parametersText"
            :placeholder="$t('jobForm.datax.parameters.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.jvm_options" :label="$t('jobForm.datax.jvmOptions')">
          <a-textarea
            v-model="jvmOptionsText"
            :placeholder="$t('jobForm.datax.jvmOptions.placeholder')"
          />
        </a-form-item>
        <a-form-item field="params.speed" :label="$t('jobForm.datax.speed')">
          <a-input-number
            v-model="form.params.speed"
            :min="0"
            :placeholder="$t('jobForm.datax.speed.placeholder')"
          />
        </a-form-item>
      </template>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Message } from '@arco-design/web-vue';
import type { FormInstance } from '@arco-design/web-vue';
import { createJob, updateJob } from '@/api/job';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  type: {
    type: String,
    default: 'create',
  },
  data: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(['update:visible', 'success']);

const { t } = useI18n();
const formRef = ref<FormInstance>();

const title = computed(() =>
  props.type === 'create' ? t('jobModal.create') : t('jobModal.edit')
);

const jobTypes = computed(() => [
  { label: 'Shell', value: 'shell' },
  { label: 'HTTP', value: 'http' },
  { label: 'DataX', value: 'datax' },
]);

const httpMethods = computed(() => [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]);

const generateFormModel = () => ({
  name: '',
  type: undefined,
  description: '',
  cron_expr: '',
  timeout: 0,
  retry_count: 0,
  retry_delay: 0,
  params: {},
});

const form = ref(generateFormModel());
const environmentText = ref('');
const headersText = ref('');
const successCodeText = ref('');
const parametersText = ref('');
const jvmOptionsText = ref('');

watch(
  () => props.visible,
  (val) => {
    if (val && props.type === 'edit' && props.data) {
      form.value = { ...props.data };
      if (form.value.type === 'shell' && form.value.params) {
        environmentText.value = JSON.stringify(form.value.params.environment || {}, null, 2);
      } else if (form.value.type === 'http' && form.value.params) {
        headersText.value = JSON.stringify(form.value.params.headers || {}, null, 2);
        successCodeText.value = (form.value.params.success_code || []).join(',');
      } else if (form.value.type === 'datax' && form.value.params) {
        parametersText.value = JSON.stringify(form.value.params.parameters || {}, null, 2);
        jvmOptionsText.value = (form.value.params.jvm_options || []).join('\n');
      }
    } else {
      form.value = generateFormModel();
      environmentText.value = '';
      headersText.value = '';
      successCodeText.value = '';
      parametersText.value = '';
      jvmOptionsText.value = '';
    }
  }
);

const handleTypeChange = () => {
  form.value.params = {};
  environmentText.value = '';
  headersText.value = '';
  successCodeText.value = '';
  parametersText.value = '';
  jvmOptionsText.value = '';
};

const parseParams = () => {
  if (form.value.type === 'shell') {
    try {
      form.value.params.environment = JSON.parse(environmentText.value || '{}');
    } catch (err) {
      throw new Error(t('jobForm.shell.environment.error'));
    }
  } else if (form.value.type === 'http') {
    try {
      form.value.params.headers = JSON.parse(headersText.value || '{}');
    } catch (err) {
      throw new Error(t('jobForm.http.headers.error'));
    }
    form.value.params.success_code = successCodeText.value
      .split(',')
      .map((code) => parseInt(code.trim()))
      .filter((code) => !isNaN(code));
  } else if (form.value.type === 'datax') {
    try {
      form.value.params.parameters = JSON.parse(parametersText.value || '{}');
    } catch (err) {
      throw new Error(t('jobForm.datax.parameters.error'));
    }
    form.value.params.jvm_options = jvmOptionsText.value
      .split('\n')
      .map((opt) => opt.trim())
      .filter((opt) => opt);
  }
};

const handleSubmit = async () => {
  if (!formRef.value) return false;
  try {
    await formRef.value.validate();
    parseParams();
    if (props.type === 'create') {
      await createJob(form.value);
      Message.success(t('jobModal.create.success'));
    } else {
      await updateJob(props.data.id, form.value);
      Message.success(t('jobModal.edit.success'));
    }
    emit('success');
    return true;
  } catch (err) {
    Message.error(err.message || t('jobModal.error'));
    return false;
  }
};

const handleCancel = () => {
  form.value = generateFormModel();
  emit('update:visible', false);
};
</script>

<style scoped>
.arco-form-item {
  margin-bottom: 20px;
}
</style>
