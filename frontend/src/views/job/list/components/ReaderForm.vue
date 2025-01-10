<template>
  <a-form :model="modelValue" layout="vertical">
    <a-form-item label="数据源类型" required>
      <a-select
        v-model="modelValue.name"
        placeholder="请选择数据源类型"
      >
        <a-option value="mysqlreader">MySQL</a-option>
        <a-option value="postgresqlreader">PostgreSQL</a-option>
        <a-option value="sqlserverreader">SQLServer</a-option>
        <a-option value="oraclereader">Oracle</a-option>
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
    <a-form-item label="查询SQL">
      <a-textarea
        v-model="modelValue.parameter.selectSql"
        placeholder="请输入查询SQL，如不填写将根据columns和where生成"
        :auto-size="{ minRows: 2, maxRows: 5 }"
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
    <a-form-item label="Where条件">
      <a-input
        v-model="modelValue.parameter.where"
        placeholder="请输入Where条件，如：status = 1"
        allow-clear
      />
    </a-form-item>
    <a-form-item label="批量大小">
      <a-input-number
        v-model="modelValue.parameter.batchSize"
        placeholder="请输入批量大小"
        :min="1"
        :max="100000"
        :step="1000"
        :default-value="20000"
      />
    </a-form-item>
  </a-form>
</template>

<script lang="ts" setup>
interface ReaderConfig {
  name: string;
  parameter: {
    username: string;
    password: string;
    host: string;
    port: number;
    database: string;
    table: string;
    columns: string[];
    selectSql?: string;
    where?: string;
    batchSize?: number;
  };
}

const props = defineProps<{
  modelValue: ReaderConfig;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: ReaderConfig): void;
}>();

const handleColumnChange = (val: string) => {
  const newValue = { ...props.modelValue };
  newValue.parameter.columns = val.split(',').map(item => item.trim()).filter(Boolean);
  emit('update:modelValue', newValue);
};
</script>
