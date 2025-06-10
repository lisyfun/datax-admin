<template>
  <a-modal
    v-model:visible="modelVisible"
    title="生成 Cron 表达式"
    :width="680"
    @cancel="handleCancel"
    @ok="handleSubmit"
  >
    <div class="cron-container">
      <div class="cron-preview">
        <div class="preview-title">表达式预览</div>
        <div class="preview-content">{{ cronExpression }}</div>
        <div class="preview-desc">从左到右依次为：秒 分 时 日 月 周</div>
      </div>

      <a-tabs type="card">
        <a-tab-pane key="second" title="秒">
          <div class="tab-content">
            <a-radio-group v-model="cronForm.second.type">
              <a-space direction="vertical" size="large">
                <a-radio value="*">
                  <span class="radio-label">每秒</span>
                  <span class="radio-desc">每秒钟触发</span>
                </a-radio>
                <a-radio value="cycle">
                  <div class="radio-content">
                    <span class="radio-label">周期</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.second.start" :min="0" :max="59" size="small" style="width: 60px" />
                      秒开始，每
                      <a-input-number v-model="cronForm.second.interval" :min="1" :max="59" size="small" style="width: 60px" />
                      秒执行一次
                    </div>
                  </div>
                </a-radio>
                <a-radio value="-">
                  <div class="radio-content">
                    <span class="radio-label">区间</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.second.from" :min="0" :max="59" size="small" style="width: 60px" />
                      到
                      <a-input-number v-model="cronForm.second.to" :min="0" :max="59" size="small" style="width: 60px" />
                      秒之间的每秒钟触发
                    </div>
                  </div>
                </a-radio>
                <a-radio value="/">
                  <div class="radio-content">
                    <span class="radio-label">指定</span>
                    <div class="radio-desc">
                      <a-select
                        v-model="cronForm.second.specify"
                        multiple
                        size="small"
                        :max-tag-count="3"
                        style="width: 360px"
                      >
                        <a-option v-for="i in 60" :key="i-1" :value="i-1">{{ i-1 }}秒</a-option>
                      </a-select>
                    </div>
                  </div>
                </a-radio>
              </a-space>
            </a-radio-group>
          </div>
        </a-tab-pane>

        <a-tab-pane key="minute" title="分">
          <div class="tab-content">
            <a-radio-group v-model="cronForm.minute.type">
              <a-space direction="vertical" size="large">
                <a-radio value="*">
                  <span class="radio-label">每分</span>
                  <span class="radio-desc">每分钟触发</span>
                </a-radio>
                <a-radio value="cycle">
                  <div class="radio-content">
                    <span class="radio-label">周期</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.minute.start" :min="0" :max="59" size="small" style="width: 60px" />
                      分开始，每
                      <a-input-number v-model="cronForm.minute.interval" :min="1" :max="59" size="small" style="width: 60px" />
                      分执行一次
                    </div>
                  </div>
                </a-radio>
                <a-radio value="-">
                  <div class="radio-content">
                    <span class="radio-label">区间</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.minute.from" :min="0" :max="59" size="small" style="width: 60px" />
                      到
                      <a-input-number v-model="cronForm.minute.to" :min="0" :max="59" size="small" style="width: 60px" />
                      分之间的每分钟触发
                    </div>
                  </div>
                </a-radio>
                <a-radio value="/">
                  <div class="radio-content">
                    <span class="radio-label">指定</span>
                    <div class="radio-desc">
                      <a-select
                        v-model="cronForm.minute.specify"
                        multiple
                        size="small"
                        :max-tag-count="3"
                        style="width: 360px"
                      >
                        <a-option v-for="i in 60" :key="i-1" :value="i-1">{{ i-1 }}分</a-option>
                      </a-select>
                    </div>
                  </div>
                </a-radio>
              </a-space>
            </a-radio-group>
          </div>
        </a-tab-pane>

        <a-tab-pane key="hour" title="时">
          <div class="tab-content">
            <a-radio-group v-model="cronForm.hour.type">
              <a-space direction="vertical" size="large">
                <a-radio value="*">
                  <span class="radio-label">每小时</span>
                  <span class="radio-desc">每小时触发</span>
                </a-radio>
                <a-radio value="cycle">
                  <div class="radio-content">
                    <span class="radio-label">周期</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.hour.start" :min="0" :max="23" size="small" style="width: 60px" />
                      时开始，每
                      <a-input-number v-model="cronForm.hour.interval" :min="1" :max="23" size="small" style="width: 60px" />
                      小时执行一次
                    </div>
                  </div>
                </a-radio>
                <a-radio value="-">
                  <div class="radio-content">
                    <span class="radio-label">区间</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.hour.from" :min="0" :max="23" size="small" style="width: 60px" />
                      到
                      <a-input-number v-model="cronForm.hour.to" :min="0" :max="23" size="small" style="width: 60px" />
                      时之间的每小时触发
                    </div>
                  </div>
                </a-radio>
                <a-radio value="/">
                  <div class="radio-content">
                    <span class="radio-label">指定</span>
                    <div class="radio-desc">
                      <a-select
                        v-model="cronForm.hour.specify"
                        multiple
                        size="small"
                        :max-tag-count="3"
                        style="width: 360px"
                      >
                        <a-option v-for="i in 24" :key="i-1" :value="i-1">{{ i-1 }}时</a-option>
                      </a-select>
                    </div>
                  </div>
                </a-radio>
              </a-space>
            </a-radio-group>
          </div>
        </a-tab-pane>

        <a-tab-pane key="day" title="日">
          <div class="tab-content">
            <a-radio-group v-model="cronForm.day.type">
              <a-space direction="vertical" size="large">
                <a-radio value="*">
                  <span class="radio-label">每日</span>
                  <span class="radio-desc">每天触发</span>
                </a-radio>
                <a-radio value="?">
                  <span class="radio-label">不指定</span>
                  <span class="radio-desc">使用星期时必须选择此项</span>
                </a-radio>
                <a-radio value="L">
                  <span class="radio-label">最后一天</span>
                  <span class="radio-desc">每月最后一天触发</span>
                </a-radio>
                <a-radio value="W">
                  <span class="radio-label">工作日</span>
                  <span class="radio-desc">每月最近的工作日触发</span>
                </a-radio>
                <a-radio value="cycle">
                  <div class="radio-content">
                    <span class="radio-label">周期</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.day.start" :min="1" :max="31" size="small" style="width: 60px" />
                      日开始，每
                      <a-input-number v-model="cronForm.day.interval" :min="1" :max="31" size="small" style="width: 60px" />
                      天执行一次
                    </div>
                  </div>
                </a-radio>
                <a-radio value="-">
                  <div class="radio-content">
                    <span class="radio-label">区间</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.day.from" :min="1" :max="31" size="small" style="width: 60px" />
                      到
                      <a-input-number v-model="cronForm.day.to" :min="1" :max="31" size="small" style="width: 60px" />
                      日之间的每天触发
                    </div>
                  </div>
                </a-radio>
                <a-radio value="/">
                  <div class="radio-content">
                    <span class="radio-label">指定</span>
                    <div class="radio-desc">
                      <a-select
                        v-model="cronForm.day.specify"
                        multiple
                        size="small"
                        :max-tag-count="3"
                        style="width: 360px"
                      >
                        <a-option v-for="i in 31" :key="i" :value="i">{{ i }}日</a-option>
                      </a-select>
                    </div>
                  </div>
                </a-radio>
              </a-space>
            </a-radio-group>
          </div>
        </a-tab-pane>

        <a-tab-pane key="month" title="月">
          <div class="tab-content">
            <a-radio-group v-model="cronForm.month.type">
              <a-space direction="vertical" size="large">
                <a-radio value="*">
                  <span class="radio-label">每月</span>
                  <span class="radio-desc">每月触发</span>
                </a-radio>
                <a-radio value="cycle">
                  <div class="radio-content">
                    <span class="radio-label">周期</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.month.start" :min="1" :max="12" size="small" style="width: 60px" />
                      月开始，每
                      <a-input-number v-model="cronForm.month.interval" :min="1" :max="12" size="small" style="width: 60px" />
                      月执行一次
                    </div>
                  </div>
                </a-radio>
                <a-radio value="-">
                  <div class="radio-content">
                    <span class="radio-label">区间</span>
                    <div class="radio-desc">从
                      <a-input-number v-model="cronForm.month.from" :min="1" :max="12" size="small" style="width: 60px" />
                      到
                      <a-input-number v-model="cronForm.month.to" :min="1" :max="12" size="small" style="width: 60px" />
                      月之间的每月触发
                    </div>
                  </div>
                </a-radio>
                <a-radio value="/">
                  <div class="radio-content">
                    <span class="radio-label">指定</span>
                    <div class="radio-desc">
                      <a-select
                        v-model="cronForm.month.specify"
                        multiple
                        size="small"
                        :max-tag-count="3"
                        style="width: 360px"
                      >
                        <a-option v-for="i in 12" :key="i" :value="i">{{ i }}月</a-option>
                      </a-select>
                    </div>
                  </div>
                </a-radio>
              </a-space>
            </a-radio-group>
          </div>
        </a-tab-pane>

        <a-tab-pane key="week" title="周">
          <div class="tab-content">
            <a-radio-group v-model="cronForm.week.type">
              <a-space direction="vertical" size="large">
                <a-radio value="*">
                  <span class="radio-label">每周</span>
                  <span class="radio-desc">每周的每一天触发</span>
                </a-radio>
                <a-radio value="?">
                  <span class="radio-label">不指定</span>
                  <span class="radio-desc">使用日期时必须选择此项</span>
                </a-radio>
                <a-radio value="L">
                  <span class="radio-label">最后一周</span>
                  <span class="radio-desc">每月最后一个星期几触发</span>
                </a-radio>
                <a-radio value="cycle">
                  <div class="radio-content">
                    <span class="radio-label">周期</span>
                    <div class="radio-desc">从星期
                      <a-input-number v-model="cronForm.week.start" :min="1" :max="7" size="small" style="width: 60px" />
                      开始，每隔
                      <a-input-number v-model="cronForm.week.interval" :min="1" :max="7" size="small" style="width: 60px" />
                      天执行一次
                    </div>
                  </div>
                </a-radio>
                <a-radio value="-">
                  <div class="radio-content">
                    <span class="radio-label">区间</span>
                    <div class="radio-desc">从星期
                      <a-input-number v-model="cronForm.week.from" :min="1" :max="7" size="small" style="width: 60px" />
                      到星期
                      <a-input-number v-model="cronForm.week.to" :min="1" :max="7" size="small" style="width: 60px" />
                      之间的每天触发
                    </div>
                  </div>
                </a-radio>
                <a-radio value="/">
                  <div class="radio-content">
                    <span class="radio-label">指定</span>
                    <div class="radio-desc">
                      <a-select
                        v-model="cronForm.week.specify"
                        multiple
                        size="small"
                        :max-tag-count="3"
                        style="width: 360px"
                      >
                        <a-option v-for="i in 7" :key="i" :value="i">星期{{ ['日', '一', '二', '三', '四', '五', '六'][i % 7] }}</a-option>
                      </a-select>
                    </div>
                  </div>
                </a-radio>
              </a-space>
            </a-radio-group>
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>
  </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from 'vue';

const props = defineProps<{
  visible: boolean;
  expression: string;
}>();

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void;
  (e: 'update:expression', expression: string): void;
}>();

interface CronField {
  type: '*' | '?' | '-' | '/' | 'L' | 'W' | 'cycle';
  start: number;
  interval: number;
  from: number;
  to: number;
  specify: number[];
}

interface CronForm {
  second: CronField;
  minute: CronField;
  hour: CronField;
  day: CronField;
  month: CronField;
  week: CronField;
}

const cronForm = reactive<CronForm>({
  second: { type: '*', start: 0, interval: 1, from: 0, to: 59, specify: [] },
  minute: { type: '*', start: 0, interval: 1, from: 0, to: 59, specify: [] },
  hour: { type: '*', start: 0, interval: 1, from: 0, to: 23, specify: [] },
  day: { type: '*', start: 1, interval: 1, from: 1, to: 31, specify: [] },
  month: { type: '*', start: 1, interval: 1, from: 1, to: 12, specify: [] },
  week: { type: '?', start: 1, interval: 1, from: 1, to: 7, specify: [] }
});

// 生成Cron表达式
const generateCronExpression = (field: CronField): string => {
  switch (field.type) {
    case '*':
      return '*';
    case '?':
      return '?';
    case '-':
      return `${field.from}-${field.to}`;
    case '/':
      return field.specify.sort((a, b) => a - b).join(',');
    case 'L':
      return 'L';
    case 'W':
      return 'W';
    case 'cycle':
      return `${field.start}/${field.interval}`;
    default:
      return '*';
  }
};

const cronExpression = computed(() => {
  const second = generateCronExpression(cronForm.second);
  const minute = generateCronExpression(cronForm.minute);
  const hour = generateCronExpression(cronForm.hour);
  let day = generateCronExpression(cronForm.day);
  const month = generateCronExpression(cronForm.month);
  let week = generateCronExpression(cronForm.week);

  if (week !== '?' && week !== '*') {
    day = '?';
  }
  else if (day !== '?' && day !== '*') {
    week = '?';
  }
  else if (day === '*' && week === '*') {
    week = '?';
  }

  return `${second} ${minute} ${hour} ${day} ${month} ${week}`;
});

// 监听日期和星期的变化，自动处理互斥关系
watch(() => cronForm.day.type, (newType) => {
  if (newType !== '?' && newType !== '*') {
    // 如果日期被指定了具体值，则星期必须设置为不指定
    cronForm.week.type = '?';
  }
});

watch(() => cronForm.week.type, (newType) => {
  if (newType !== '?' && newType !== '*') {
    // 如果星期被指定了具体值，则日期必须设置为不指定
    cronForm.day.type = '?';
  }
});

// 处理对话框
const handleCancel = () => {
  // 重置表单状态
  cronForm.second.type = '*';
  cronForm.minute.type = '*';
  cronForm.hour.type = '*';
  cronForm.day.type = '*';
  cronForm.month.type = '*';
  cronForm.week.type = '?';
  emit('update:visible', false);
};

const handleSubmit = () => {
  emit('update:expression', cronExpression.value);
  emit('update:visible', false);
};

// 处理visible的双向绑定
const modelVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
});
</script>

<style scoped>
.cron-container {
  padding: 0 16px;

  .cron-preview {
    margin-bottom: 24px;
    padding: 16px;
    background: rgba(var(--primary-6), 0.08);
    border: 1px solid rgba(var(--primary-6), 0.15);
    border-radius: 4px;
    text-align: center;

    .preview-title {
      color: var(--color-text-2);
      font-size: 14px;
      margin-bottom: 8px;
    }

    .preview-content {
      color: rgb(var(--primary-6));
      font-size: 20px;
      font-family: monospace;
      margin: 8px 0;
      font-weight: bold;
      padding: 8px;
      background: var(--color-bg-2);
      border-radius: 4px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }

    .preview-desc {
      color: var(--color-text-3);
      font-size: 12px;
    }
  }

  :deep(.arco-tabs-card) {
    .arco-tabs-header {
      margin-bottom: 16px;
    }

    .arco-tabs-nav {
      background: var(--color-bg-2);
      border: 1px solid var(--color-border);
      border-radius: 4px;
      padding: 4px;
    }

    .arco-tabs-tab {
      border-radius: 4px;
      padding: 8px 16px;
      font-size: 14px;

      &:hover {
        color: rgb(var(--primary-6));
        background: rgba(var(--primary-6), 0.06);
      }

      &.arco-tabs-tab-active {
        background: rgba(var(--primary-6), 0.1);
      }
    }
  }

  .tab-content {
    padding: 8px;
    background: var(--color-bg-2);
    border-radius: 4px;

    .radio-content {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }

    .radio-label {
      font-size: 14px;
      font-weight: 500;
      color: var(--color-text-1);
    }

    .radio-desc {
      color: var(--color-text-3);
      font-size: 13px;
      line-height: 32px;
    }
  }

  :deep(.arco-radio) {
    align-items: flex-start;
    padding: 8px 12px;
    width: 100%;
    border-radius: 4px;
    transition: all 0.2s ease;
    border: 1px solid transparent;

    &:hover {
      background: rgba(var(--primary-6), 0.04);
      border-color: rgba(var(--primary-6), 0.1);
    }
  }

  :deep(.arco-radio-checked) {
    background: rgba(var(--primary-6), 0.08);
    border-color: rgba(var(--primary-6), 0.2);
  }

  :deep(.arco-input-number) {
    margin: 0 4px;
    border-color: var(--color-border);
    background: var(--color-bg-1);

    &:hover, &:focus {
      background: var(--color-bg-2);
      border-color: rgb(var(--primary-6));
    }
  }

  :deep(.arco-select) {
    width: 100%;

    .arco-select-view {
      background: var(--color-bg-1);
      border-color: var(--color-border);

      &:hover {
        background: var(--color-bg-2);
        border-color: rgb(var(--primary-6));
      }
    }
  }
}
</style>
