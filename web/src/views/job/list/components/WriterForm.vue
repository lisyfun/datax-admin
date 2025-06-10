<template>
  <a-form :model="modelValue" layout="vertical">
    <a-form-item label="数据源类型" required>
      <a-select
        v-model="modelValue.name"
        placeholder="请选择数据源类型"
      >
        <a-option value="mysqlwriter">MySQL</a-option>
        <a-option value="postgresqlwriter">PostgreSQL</a-option>
        <a-option value="sqlserverwriter">SQLServer</a-option>
        <a-option value="oraclewriter">Oracle</a-option>
      </a-select>
    </a-form-item>
    <a-form-item label="用户名" required>
      <a-input
        v-model="modelValue.parameter.username"
        placeholder="请输入用户名"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="密码" required>
      <a-input-password
        v-model="modelValue.parameter.password"
        placeholder="请输入密码"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="主机地址" required>
      <a-input
        v-model="modelValue.parameter.host"
        placeholder="请输入主机地址"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="端口" required>
      <a-input-number
        v-model="modelValue.parameter.port"
        placeholder="请输入端口号"
        :min="1"
        :max="65535"
      />
    </a-form-item>
    <a-form-item label="数据库" required>
      <a-input
        v-model="modelValue.parameter.database"
        placeholder="请输入数据库名"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="表名" required>
      <a-input
        v-model="modelValue.parameter.table"
        placeholder="请输入表名"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="字段">
      <a-input
        :model-value="modelValue.parameter.columns?.join(',')"
        @update:model-value="handleColumnChange"
        placeholder="请输入字段名，多个字段用逗号分隔"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="写入模式">
      <a-select
        v-model="modelValue.parameter.writeMode"
        placeholder="请选择写入模式"
      >
        <a-option value="insert">Insert</a-option>
        <a-option value="replace">Replace</a-option>
        <a-option value="update">Update</a-option>
      </a-select>
    </a-form-item>
    <a-form-item label="批量大小">
      <a-input-number
        v-model="modelValue.parameter.batchSize"
        placeholder="请输入批量大小"
        :min="1"
        :max="100000"
        :step="1000"
        :default-value="10000"
      />
    </a-form-item>
    <a-form-item label="前置SQL">
      <a-textarea
        :model-value="modelValue.parameter.preSql?.join('\n')"
        @update:model-value="handlePreSqlChange"
        placeholder="请输入前置SQL，每行一条"
        :auto-size="{ minRows: 2, maxRows: 5 }"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="后置SQL">
      <a-textarea
        :model-value="modelValue.parameter.postSql?.join('\n')"
        @update:model-value="handlePostSqlChange"
        placeholder="请输入后置SQL，每行一条"
        :auto-size="{ minRows: 2, maxRows: 5 }"
        allow-clear
      />
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
interface WriterConfig {
  name: string;
  parameter: {
    username: string;
    password: string;
    host: string;
    port: number;
    database: string;
    table: string;
    columns: string[];
    writeMode: 'insert' | 'replace' | 'update';
    batchSize?: number;
    preSql?: string[];
    postSql?: string[];
  };
}

const props = defineProps<{
  modelValue: WriterConfig;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: WriterConfig): void;
}>();

const handleColumnChange = (val: string) => {
  const newValue = { ...props.modelValue };
  newValue.parameter.columns = val.split(',').map(item => item.trim()).filter(Boolean);
  emit('update:modelValue', newValue);
};

const handlePreSqlChange = (val: string) => {
  const newValue = { ...props.modelValue };
  newValue.parameter.preSql = val.split('\n').map(item => item.trim()).filter(Boolean);
  emit('update:modelValue', newValue);
};

const handlePostSqlChange = (val: string) => {
  const newValue = { ...props.modelValue };
  newValue.parameter.postSql = val.split('\n').map(item => item.trim()).filter(Boolean);
  emit('update:modelValue', newValue);
};
</script>
