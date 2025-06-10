<template>
  <div class="menu-management">
    <a-card>
      <template #title>菜单管理</template>
      <template #extra>
        <a-space>
          <a-input-search
            v-model="searchKeyword"
            placeholder="请输入菜单名称"
            style="width: 300px"
            @search="handleSearch"
          />
          <a-button type="primary" @click="() => handleAdd()">
            <template #icon><icon-plus /></template>
            新增菜单
          </a-button>
        </a-space>
      </template>

      <a-table
        :data="menuList"
        :loading="loading"
        :pagination="false"
        row-key="id"
        :tree-props="{
          children: 'children'
        }"
      >
        <template #columns>
          <a-table-column title="菜单名称" data-index="name" />
          <a-table-column title="路由路径" data-index="path" />
          <a-table-column title="组件路径" data-index="component" />
          <a-table-column title="图标" align="center">
            <template #cell="{ record }">
              <component v-if="record.icon" :is="iconMap[record.icon]" />
              <span v-else>-</span>
            </template>
          </a-table-column>
          <a-table-column title="类型" align="center">
            <template #cell="{ record }">
              <a-tag :color="record.type === 1 ? 'blue' : 'green'">
                {{ record.type === 1 ? '菜单' : '按钮' }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="排序" data-index="sort" align="center" />
          <a-table-column title="状态" align="center">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.status === 1"
                @update:model-value="(value) => handleStatusChange(record, Boolean(value))"
              />
            </template>
          </a-table-column>
          <a-table-column title="缓存" align="center">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.cache === 1"
                @update:model-value="(value) => handleCacheChange(record, Boolean(value))"
              />
            </template>
          </a-table-column>
          <a-table-column title="隐藏" align="center">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.hidden === 1"
                @update:model-value="(value) => handleHiddenChange(record, Boolean(value))"
              />
            </template>
          </a-table-column>
          <a-table-column title="操作" align="center">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="() => handleAdd(record)">
                  <template #icon><icon-plus /></template>
                  添加子菜单
                </a-button>
                <a-button type="text" size="small" @click="() => handleEdit(record)">
                  <template #icon><icon-edit /></template>
                  编辑
                </a-button>
                <a-popconfirm
                  content="确定要删除该菜单吗？"
                  @ok="() => handleDelete(record)"
                >
                  <a-button type="text" size="small" status="danger">
                    <template #icon><icon-delete /></template>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <!-- 菜单表单对话框 -->
    <a-modal
      v-model:visible="showMenuForm"
      :title="isEdit ? '编辑菜单' : '新增菜单'"
      @ok="handleSubmit"
      @cancel="() => showMenuForm = false"
    >
      <a-form ref="formRef" :model="formData" :rules="formRules">
        <a-form-item field="name" label="菜单名称">
          <a-input v-model="formData.name" placeholder="请输入菜单名称" />
        </a-form-item>
        <a-form-item field="type" label="类型">
          <a-radio-group v-model="formData.type">
            <a-radio :value="1">菜单</a-radio>
            <a-radio :value="2">按钮</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item field="parent_id" label="上级菜单">
          <a-tree-select
            v-model="formData.parent_id"
            :data="menuTreeData"
            placeholder="请选择上级菜单"
            allow-clear
          />
        </a-form-item>
        <template v-if="formData.type === 1">
          <a-form-item field="path" label="路由路径">
            <a-input v-model="formData.path" placeholder="请输入路由路径" />
          </a-form-item>
          <a-form-item field="component" label="组件路径">
            <a-input v-model="formData.component" placeholder="请输入组件路径" />
          </a-form-item>
          <a-form-item field="icon" label="图标">
            <a-input v-model="formData.icon" placeholder="请输入图标类名" />
          </a-form-item>
        </template>
        <a-form-item field="sort" label="排序">
          <a-input-number v-model="formData.sort" placeholder="请输入排序" :min="0" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-switch
            :model-value="formData.status === 1"
            @update:model-value="(value) => formData.status = value ? 1 : 0"
          />
        </a-form-item>
        <a-form-item field="cache" label="缓存">
          <a-switch
            :model-value="formData.cache === 1"
            @update:model-value="(value) => formData.cache = value ? 1 : 0"
          />
        </a-form-item>
        <a-form-item field="hidden" label="隐藏">
          <a-switch
            :model-value="formData.hidden === 1"
            @update:model-value="(value) => formData.hidden = value ? 1 : 0"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { FormInstance } from '@arco-design/web-vue';
import type { TreeNodeData } from '@arco-design/web-vue';
import { IconPlus, IconEdit, IconDelete, IconApps, IconUser, IconUserGroup, IconSafe, IconCalendar, IconUnorderedList, IconClockCircle, IconDashboard, IconDesktop, IconCloud, IconFile, IconBulb, IconCode, IconRobot, IconCommon, IconCommand, IconLock, IconList } from '@arco-design/web-vue/es/icon';
import { getMenuList, createMenu, updateMenu, deleteMenu } from '@/api/menu';
import type { MenuResponse, CreateMenuRequest, UpdateMenuRequest } from '@/api/menu';
import { useMenuStore } from '@/stores/modules/menu';

// 搜索关键字
const searchKeyword = ref('');
// 加载状态
const loading = ref(false);
// 菜单列表数据
const menuList = ref<MenuResponse[]>([]);
// 菜单树形数据（用于选择上级菜单）
const menuTreeData = ref<TreeNodeData[]>([]);
// 是否显示菜单表单
const showMenuForm = ref(false);
// 是否为编辑模式
const isEdit = ref(false);
// 表单引用
const formRef = ref<FormInstance>();

// 表单数据
const formData = reactive<Omit<CreateMenuRequest, 'parent_id'> & { id: number; parent_id?: number; status: number; cache: number; hidden: number }>({
  id: 0,
  name: '',
  type: 1,
  parent_id: undefined,
  path: '',
  component: '',
  icon: '',
  sort: 0,
  status: 1,
  cache: 1,
  hidden: 0,
});

// 表单校验规则
const formRules = {
  name: [{ required: true, message: '请输入菜单名称' }],
  type: [{ required: true, message: '请选择菜单类型' }],
  path: [{ required: true, message: '请输入路由路径', trigger: 'blur' }],
  component: [{ required: true, message: '请输入组件路径', trigger: 'blur' }],
};

// 图标映射
const iconMap: Record<string, any> = {
  'icon-dashboard': IconDashboard,
  'icon-apps': IconApps,
  'icon-user': IconUser,
  'icon-user-group': IconUserGroup,
  'icon-safe': IconSafe,
  'icon-calendar': IconCalendar,
  'icon-unordered-list': IconUnorderedList,
  'icon-clock-circle': IconClockCircle,
  'icon-desktop': IconDesktop,
  'icon-cloud': IconCloud,
  'icon-file': IconFile,
  'icon-bulb': IconBulb,
  'icon-code': IconCode,
  'icon-robot': IconRobot,
  'icon-common': IconCommon,
  'icon-command': IconCommand,
  'icon-lock': IconLock,
  'icon-menu': IconList,
};

// 添加 menuStore
const menuStore = useMenuStore();

// 转换为树形选择数据
const convertToTreeData = (menus: MenuResponse[]): TreeNodeData[] => {
  return menus.map(menu => ({
    key: menu.id,
    title: menu.name,
    children: menu.children ? convertToTreeData(menu.children) : undefined,
  }));
};

// 获取菜单列表
const fetchMenuList = async () => {
  loading.value = true;
  try {
    const { data } = await getMenuList();
    menuList.value = data.list;
    menuTreeData.value = convertToTreeData(data.list);
  } catch (error: any) {
    Message.error(error.message || '获取菜单列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  fetchMenuList();
};

// 新增菜单
const handleAdd = (parent?: any) => {
  isEdit.value = false;
  formData.id = 0;
  formData.name = '';
  formData.type = 1;
  formData.parent_id = parent?.id;
  formData.path = '';
  formData.component = '';
  formData.icon = '';
  formData.sort = 0;
  formData.status = 1;
  formData.cache = 1;
  formData.hidden = 0;
  showMenuForm.value = true;
};

// 编辑菜单
const handleEdit = (record: any) => {
  isEdit.value = true;
  formData.id = record.id;
  formData.name = record.name;
  formData.type = record.type;
  formData.parent_id = record.parent_id;
  formData.path = record.path;
  formData.component = record.component;
  formData.icon = record.icon;
  formData.sort = record.sort;
  formData.status = record.status;
  formData.cache = record.cache;
  formData.hidden = record.hidden;
  showMenuForm.value = true;
};

// 删除菜单
const handleDelete = async (record: any) => {
  try {
    await deleteMenu(record.id);
    Message.success('删除成功');
    fetchMenuList();
    // 重新获取用户菜单
    await menuStore.fetchUserMenus();
  } catch (error: any) {
    Message.error(error.message || '删除失败');
  }
};

// 更新菜单状态
const handleStatusChange = async (record: MenuResponse, value: boolean) => {
  try {
    await updateMenu(record.id, {
      name: record.name,
      type: record.type,
      parent_id: record.parent_id,
      path: record.path,
      component: record.component,
      icon: record.icon,
      sort: record.sort,
      status: value ? 1 : 0,
      cache: record.cache as 0 | 1,
      hidden: record.hidden as 0 | 1,
    });
    Message.success('更新成功');
    fetchMenuList();
    // 重新获取用户菜单
    await menuStore.fetchUserMenus();
  } catch (error: any) {
    Message.error(error.message || '更新失败');
  }
};

// 更新菜单缓存状态
const handleCacheChange = async (record: MenuResponse, value: boolean) => {
  try {
    await updateMenu(record.id, {
      name: record.name,
      type: record.type,
      parent_id: record.parent_id,
      path: record.path,
      component: record.component,
      icon: record.icon,
      sort: record.sort,
      status: record.status as 0 | 1,
      cache: value ? 1 : 0,
      hidden: record.hidden as 0 | 1,
    });
    Message.success('更新成功');
    fetchMenuList();
    // 重新获取用户菜单
    await menuStore.fetchUserMenus();
  } catch (error: any) {
    Message.error(error.message || '更新失败');
  }
};

// 更新菜单隐藏状态
const handleHiddenChange = async (record: MenuResponse, value: boolean) => {
  try {
    await updateMenu(record.id, {
      name: record.name,
      type: record.type,
      parent_id: record.parent_id,
      path: record.path,
      component: record.component,
      icon: record.icon,
      sort: record.sort,
      status: record.status as 0 | 1,
      cache: record.cache as 0 | 1,
      hidden: value ? 1 : 0,
    });
    Message.success('更新成功');
    fetchMenuList();
    // 重新获取用户菜单
    await menuStore.fetchUserMenus();
  } catch (error: any) {
    Message.error(error.message || '更新失败');
  }
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    const params: CreateMenuRequest | UpdateMenuRequest = {
      name: formData.name,
      type: formData.type,
      parent_id: formData.parent_id || 0,
      path: formData.path,
      component: formData.component,
      icon: formData.icon,
      sort: formData.sort,
      status: formData.status as 0 | 1,
      cache: formData.cache as 0 | 1,
      hidden: formData.hidden as 0 | 1,
    };

    if (isEdit.value) {
      await updateMenu(formData.id, params);
      Message.success('更新成功');
    } else {
      await createMenu(params as CreateMenuRequest);
      Message.success('创建成功');
    }

    showMenuForm.value = false;
    fetchMenuList();
    // 重新获取用户菜单
    await menuStore.fetchUserMenus();
  } catch (error: any) {
    Message.error(error.message || (isEdit.value ? '更新失败' : '创建失败'));
  }
};

// 页面加载时获取菜单列表
onMounted(() => {
  fetchMenuList();
});
</script>

<style scoped>
.menu-management {
  padding: 16px;
}
</style>
