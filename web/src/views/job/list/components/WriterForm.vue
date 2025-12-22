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
    <a-form-item
      v-if="modelValue.name === 'postgresqlwriter'"
      label="Schema"
    >
      <a-input
        v-model="modelValue.parameter.schema"
        placeholder="请输入Schema名称，默认为public"
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
      <div class="sql-list">
        <div
          v-for="(sql, index) in (modelValue.parameter.preSql || [])"
          :key="index"
          class="sql-item"
        >
          <a-space :size="8" fill align="start">
            <a-textarea
              :model-value="sql"
              @update:model-value="(val: string) => handlePreSqlItemChange(index, val)"
              placeholder="请输入前置SQL"
              :auto-size="{ minRows: 1, maxRows: 5 }"
              allow-clear
              style="width: 560px"
            />
            <a-button
              type="text"
              status="danger"
              @click="removePreSql(index)"
              style="margin-top: 2px"
            >
              <template #icon><icon-delete /></template>
            </a-button>
          </a-space>
        </div>
        <div class="sql-add">
          <a-button type="dashed" long @click="addPreSql">
            <template #icon><icon-plus /></template>
            添加前置SQL
          </a-button>
        </div>
      </div>
    </a-form-item>
    <a-form-item label="后置SQL">
      <div class="sql-list">
        <div
          v-for="(sql, index) in (modelValue.parameter.postSql || [])"
          :key="index"
          class="sql-item"
        >
          <a-space :size="8" fill align="start">
            <a-textarea
              :model-value="sql"
              @update:model-value="(val: string) => handlePostSqlItemChange(index, val)"
              placeholder="请输入后置SQL"
              :auto-size="{ minRows: 1, maxRows: 5 }"
              allow-clear
              style="width: 560px"
            />
            <a-button
              type="text"
              status="danger"
              @click="removePostSql(index)"
              style="margin-top: 2px"
            >
              <template #icon><icon-delete /></template>
            </a-button>
          </a-space>
        </div>
        <div class="sql-add">
          <a-button type="dashed" long @click="addPostSql">
            <template #icon><icon-plus /></template>
            添加后置SQL
          </a-button>
        </div>
      </div>
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
    schema?: string;
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

// PreSQL handlers
const addPreSql = () => {
  const newValue = { ...props.modelValue };
  if (!newValue.parameter.preSql) {
    newValue.parameter.preSql = [];
  }
  newValue.parameter.preSql.push('');
  emit('update:modelValue', newValue);
};

const removePreSql = (index: number) => {
  const newValue = { ...props.modelValue };
  if (newValue.parameter.preSql) {
    newValue.parameter.preSql.splice(index, 1);
    emit('update:modelValue', newValue);
  }
};

const handlePreSqlItemChange = (index: number, val: string) => {
  const newValue = { ...props.modelValue };
  if (newValue.parameter.preSql) {
    newValue.parameter.preSql[index] = val;
    emit('update:modelValue', newValue);
  }
};

// PostSQL handlers
const addPostSql = () => {
  const newValue = { ...props.modelValue };
  if (!newValue.parameter.postSql) {
    newValue.parameter.postSql = [];
  }
  newValue.parameter.postSql.push('');
  emit('update:modelValue', newValue);
};

const removePostSql = (index: number) => {
  const newValue = { ...props.modelValue };
  if (newValue.parameter.postSql) {
    newValue.parameter.postSql.splice(index, 1);
    emit('update:modelValue', newValue);
  }
};

const handlePostSqlItemChange = (index: number, val: string) => {
  const newValue = { ...props.modelValue };
  if (newValue.parameter.postSql) {
    newValue.parameter.postSql[index] = val;
    emit('update:modelValue', newValue);
  }
};
</script>

<style scoped>
.sql-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.sql-item {
  display: flex;
  align-items: center;
}

.sql-item :deep(.arco-space-fill) {
  width: 100%;
}

.sql-item :deep(.arco-input-wrapper),
.sql-item :deep(.arco-textarea-wrapper) {
  flex: 1;
}

.sql-add {
  margin-top: 4px;
}

.sql-add :deep(.arco-btn) {
  border-style: dashed;
}
</style>
