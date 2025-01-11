<template>
  <div class="json-form-wrapper">
    <template v-if="!value || typeof value !== 'object'">
      <a-input
        :model-value="String(value)"
        @update:model-value="(val) => $emit('update', { path: [], value: val })"
        placeholder="请输入值"
      />
    </template>

    <template v-else-if="Array.isArray(value)">
      <div class="json-array">
        <div v-for="(item, index) in value" :key="index" class="json-array-item">
          <JsonForm
            :value="item"
            :path="[...currentPath, String(index)]"
            @update="({ value: val }) => handleArrayUpdate(index, val)"
          />
          <a-button type="text" status="danger" @click="handleArrayDelete(index)">
            <template #icon><icon-delete /></template>
          </a-button>
        </div>
        <a-button type="dashed" long @click="handleArrayAdd">
          <template #icon><icon-plus /></template>
          添加项
        </a-button>
      </div>
    </template>

    <template v-else>
      <div class="json-object">
        <div v-for="(val, key) in value" :key="key" class="json-field">
          <div class="json-field-label">{{ key }}</div>
          <div class="json-field-value">
            <JsonForm
              :value="val"
              :path="[...currentPath, String(key)]"
              @update="({ value: newVal }) => handleObjectUpdate(String(key), newVal)"
            />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { IconDelete, IconPlus } from '@arco-design/web-vue/es/icon';

interface JsonFormProps {
  value: any;
  path?: string[];
}

const props = defineProps<JsonFormProps>();

const emit = defineEmits<{
  (e: 'update', payload: { path: string[], value: any }): void;
}>();

const currentPath = props.path || [];

const handleArrayUpdate = (index: number, value: any) => {
  const newArray = [...(props.value as any[])];
  newArray[index] = value;
  emit('update', { path: currentPath, value: newArray });
};

const handleArrayDelete = (index: number) => {
  const newArray = [...(props.value as any[])];
  newArray.splice(index, 1);
  emit('update', { path: currentPath, value: newArray });
};

const handleArrayAdd = () => {
  const array = props.value as any[];
  const newArray = [...array];
  newArray.push(typeof array[0] === 'object' ? {} : '');
  emit('update', { path: currentPath, value: newArray });
};

const handleObjectUpdate = (key: string, value: any) => {
  const newObject = { ...props.value };
  newObject[key] = value;
  emit('update', { path: currentPath, value: newObject });
};
</script>

<style scoped>
.json-form-wrapper {
  width: 100%;
}

.json-object {
  padding: 8px 0;
}

.json-field {
  margin-bottom: 16px;
  padding-left: 16px;
  border-left: 2px solid var(--color-neutral-3);
}

.json-field:last-child {
  margin-bottom: 0;
}

.json-field-label {
  font-size: 14px;
  color: var(--color-text-2);
  margin-bottom: 4px;
}

.json-field-value {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.json-array {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px 0;
}

.json-array-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.json-array-item :deep(.arco-input-wrapper) {
  flex: 1;
}
</style>
