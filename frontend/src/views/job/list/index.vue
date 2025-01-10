<template>
  <div class="jobs">
    <a-card>
      <template #title>任务列表</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchForm.name"
            placeholder="请输入任务名称"
            style="width: 300px"
            allow-clear
            @search="handleSearch"
            @press-enter="handleSearch"
          />
          <a-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            style="width: 120px"
            allow-clear
            @change="handleSearch"
          >
            <a-option :value="1">运行中</a-option>
            <a-option :value="0">已停止</a-option>
          </a-select>
          <a-button type="primary" @click="handleCreate">
            <template #icon><IconPlus /></template>
            新建任务
          </a-button>
          <a-button type="primary" status="success" @click="handleBatchExecute" :disabled="!selectedKeys.length">
            <template #icon><IconPlayCircle /></template>
            批量执行一次
          </a-button>
          <a-button @click="fetchData">
            <template #icon><IconRefresh /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        :columns="columns"
        :data="renderData"
        :row-selection="rowSelection"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '运行中' : '已停止' }}
          </a-tag>
        </template>
        <template #type="{ record }">
          <a-tag :color="getJobTypeColor(record.type)">
            {{ getJobTypeText(record.type) }}
          </a-tag>
        </template>
        <template #created_at="{ record }">
          {{ formatDateTime(record.created_at) }}
        </template>
        <template #updated_at="{ record }">
          {{ formatDateTime(record.updated_at) }}
        </template>
        <template #operations="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleEdit(record)">
              <template #icon><IconEdit /></template>
              编辑
            </a-button>
            <a-button
              type="text"
              size="small"
              status="success"
              @click="handleExecute(record)"
            >
              <template #icon><IconPlayCircle /></template>
              执行一次
            </a-button>
            <a-button
              v-if="record.status === 0"
              type="text"
              size="small"
              @click="handleStart(record)"
            >
              <template #icon><IconPlayCircle /></template>
              启动
            </a-button>
            <a-button
              v-else
              type="text"
              size="small"
              status="warning"
              @click="handleStop(record)"
            >
              <template #icon><IconPauseCircle /></template>
              停止
            </a-button>
            <a-button
              type="text"
              size="small"
              @click="handleHistory(record)"
            >
              <template #icon><IconHistory /></template>
              执行历史
            </a-button>
            <a-popconfirm
              content="确定要删除该任务吗？"
              @ok="handleDelete(record)"
            >
              <a-button type="text" size="small" status="danger">
                <template #icon><IconDelete /></template>
                删除
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 任务表单对话框 -->
    <a-modal
      v-model:visible="showForm"
      :title="isEdit ? '编辑任务' : '新建任务'"
      @ok="handleFormSubmit"
      @cancel="handleFormCancel"
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
                    <template #icon><IconCode /></template>
                    格式化JSON
                  </a-button>
                  <a-button type="text" @click="handleLoadTemplate">
                    <template #icon><IconFile /></template>
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
              <a-button @click="showCronModal = true">
                <template #icon><IconEdit /></template>
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
    </a-modal>

    <!-- Cron表达式生成器对话框 -->
    <a-modal
      v-model:visible="showCronModal"
      title="生成 Cron 表达式"
      :width="680"
      @cancel="handleCronModalCancel"
      @ok="handleCronModalSubmit"
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
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useRouter } from 'vue-router';
import {
  IconPlus,
  IconPlayCircle,
  IconPauseCircle,
  IconRefresh,
  IconEdit,
  IconDelete,
  IconCode,
  IconFile,
  IconHistory,
} from '@arco-design/web-vue/es/icon';
import type {
  Job,
  JobStatus,
  JobListRequest,
  JobShellParams,
  JobHTTPParams,
  JobDataXParams,
} from '@/api/types';
import {
  getJobList,
  startJob,
  stopJob,
  deleteJob,
  executeJob,
  executeJobs,
  createJob,
  updateJob,
} from '@/api/job';

const router = useRouter();

const loading = ref(false);
const renderData = ref<Job[]>([]);
const selectedKeys = ref<number[]>([]);

interface SearchFormState {
  name: string;
  status: JobStatus | '';
}

const searchForm = reactive<SearchFormState>({
  name: '',
  status: '',
});

const columns = [
  {
    title: '任务名称',
    dataIndex: 'name',
    align: 'center' as const,
  },
  {
    title: '任务类型',
    dataIndex: 'type',
    slotName: 'type',
    align: 'center' as const,
  },
  {
    title: '任务描述',
    dataIndex: 'description',
    align: 'center' as const,
  },
  {
    title: 'Cron 表达式',
    dataIndex: 'cron_expr',
    width: 200,
    align: 'center' as const,
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
    align: 'center' as const,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    slotName: 'created_at',
    align: 'center' as const,
  },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
    width: 320,
    align: 'center' as const,
  },
];

const rowSelection = {
  type: 'checkbox' as const,
  showCheckedAll: true,
  onlyCurrent: false,
  onChange: (selectedRowKeys: (string | number)[]) => {
    selectedKeys.value = selectedRowKeys as number[];
  },
};

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
});

// 表单相关
const showForm = ref(false);
const isEdit = ref(false);
const formRef = ref();

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

// 添加表单验证规则
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

// 添加Cron表达式生成器对话框控制变量
const showCronModal = ref(false);

// Cron表达式生成器相关
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

// 处理Cron表达式生成器对话框
const handleCronModalCancel = () => {
  // 重置表单状态
  cronForm.second.type = '*';
  cronForm.minute.type = '*';
  cronForm.hour.type = '*';
  cronForm.day.type = '*';
  cronForm.month.type = '*';
  cronForm.week.type = '?';
  showCronModal.value = false;
};

const handleCronModalSubmit = () => {
  form.cron_expr = cronExpression.value;
  showCronModal.value = false;
};

// 获取任务列表数据
const fetchData = async () => {
  try {
    loading.value = true;
    const { data } = await getJobList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: searchForm.name,
      status: searchForm.status ? searchForm.status as JobStatus : undefined,
    } as JobListRequest);
    renderData.value = data.items;
    pagination.total = data.total;
  } catch (err) {
    Message.error('获取任务列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  pagination.current = 1;
  fetchData();
};

// 分页变化
const onPageChange = (current: number) => {
  pagination.current = current;
  fetchData();
};

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  fetchData();
};

// 创建任务
const handleCreate = () => {
  isEdit.value = false;
  form.id = 0;
  form.name = '';
  form.description = '';
  form.type = 'shell';
  form.command = '';
  form.working_dir = '';
  form.url = '';
  form.method = 'GET';
  form.headers = '{}';
  form.body = '';
  form.success_codes = '200';
  form.cron_expr = '';
  form.timeout = 0;
  form.retry_count = 0;
  form.retry_delay = 0;
  form.params = {};
  showForm.value = true;
};

// 编辑任务
const handleEdit = (record: Job) => {
  isEdit.value = true;
  form.id = record.id;
  form.name = record.name;
  form.description = record.description || '';
  form.type = record.type as 'shell' | 'http' | 'datax';

  // 从params中解析出相应的参数
  try {
    const params = typeof record.params === 'string' ? JSON.parse(record.params) : record.params;
    if (record.type === 'shell') {
      form.command = params.command || '';
      form.working_dir = params.work_dir || '';
    } else if (record.type === 'http') {
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

  form.cron_expr = record.cron_expr;
  form.timeout = record.timeout || 0;
  form.retry_count = record.retry_count || 0;
  form.retry_delay = record.retry_delay || 0;
  form.params = {};
  showForm.value = true;
};

// 执行一次任务
const handleExecute = async (record: Job) => {
  try {
    await executeJob(record.id);
    Message.success('执行任务成功');
  } catch (err) {
    Message.error('执行任务失败');
  }
};

// 批量执行任务
const handleBatchExecute = async () => {
  if (!selectedKeys.value.length) {
    Message.warning('请选择要执行的任务');
    return;
  }

  try {
    await executeJobs(selectedKeys.value);
    Message.success('批量执行任务成功');
  } catch (err) {
    Message.error('批量执行任务失败');
  }
};

// 启动任务
const handleStart = async (record: Job) => {
  try {
    await startJob(record.id);
    Message.success('启动任务成功');
    fetchData();
  } catch (err) {
    Message.error('启动任务失败');
  }
};

// 停止任务
const handleStop = async (record: Job) => {
  try {
    await stopJob(record.id);
    Message.success('停止任务成功');
    fetchData();
  } catch (err) {
    Message.error('停止任务失败');
  }
};

// 删除任务
const handleDelete = async (record: Job) => {
  try {
    await deleteJob(record.id);
    Message.success('删除任务成功');
    fetchData();
  } catch (err) {
    Message.error('删除任务失败');
  }
};

// 提交表单
const handleFormSubmit = async () => {
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

    if (isEdit.value) {
      await updateJob(form.id, data);
      Message.success('编辑任务成功');
    } else {
      await createJob(data);
      Message.success('创建任务成功');
    }

    showForm.value = false;
    fetchData();
  } catch (err: any) {
    if (err.response?.data?.error) {
      Message.error(err.response.data.error);
    } else if (err.errors) {
      Message.error('表单验证失败，请检查输入');
    } else {
      Message.error(isEdit.value ? '编辑任务失败' : '创建任务失败');
    }
  }
};

// 取消表单
const handleFormCancel = () => {
  showForm.value = false;
  if (formRef.value) {
    formRef.value.resetFields();
  }
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

// 添加任务类型相关的工具函数
const getJobTypeText = (type: string) => {
  switch (type) {
    case 'shell':
      return 'Shell脚本';
    case 'http':
      return 'HTTP请求';
    case 'datax':
      return 'DataX任务';
    default:
      return '未知类型';
  }
};

const getJobTypeColor = (type: string) => {
  switch (type) {
    case 'shell':
      return 'blue';
    case 'http':
      return 'green';
    case 'datax':
      return 'purple';
    default:
      return 'gray';
  }
};

// 添加时间格式化函数
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  });
};

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

// 跳转到执行历史
const handleHistory = (record: Job) => {
  router.push({
    path: '/job/history',
    query: {
      job_id: record.id.toString(),
      job_name: record.name
    }
  });
};

// 初始化加载数据
fetchData();
</script>

<style scoped>
.jobs {
  padding: 16px;
}


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
