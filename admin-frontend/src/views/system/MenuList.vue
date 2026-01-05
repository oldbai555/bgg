<template>
  <div class="page">
    <el-card>
      <div class="toolbar">
        <el-button type="success" @click="openCreate()" v-permission="'menu:create'">
          {{ t('common.create') }}
        </el-button>
      </div>
      <el-tree
        :data="treeData"
        node-key="id"
        :props="{label: 'name', children: 'children'}"
        default-expand-all
      >
        <template #default="{data}">
          <span class="menu-item">
            <el-icon v-if="getMenuIcon(data.icon)" class="menu-icon">
              <component :is="getMenuIcon(data.icon)" />
            </el-icon>
            <span class="menu-name">{{ data.name }}</span>
            <el-tag v-if="data.type === 1" size="small" type="info">目录</el-tag>
            <el-tag v-else-if="data.type === 2" size="small" type="success">菜单</el-tag>
            <el-tag v-else-if="data.type === 3" size="small" type="warning">按钮</el-tag>
          </span>
          <span class="ops">
            <el-button link type="primary" @click.stop="openCreate(data)" v-permission="'menu:create'">
              {{ t('common.create') }}
            </el-button>
            <el-button link type="primary" @click.stop="openEdit(data)" v-permission="'menu:update'">
              {{ t('common.edit') }}
            </el-button>
            <el-button link type="danger" @click.stop="handleDelete(data)" v-permission="'menu:delete'">
              {{ t('common.delete') }}
            </el-button>
          </span>
        </template>
      </el-tree>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? t('common.edit') : t('common.create')" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item :label="t('common.parent')" prop="parentId">
          <el-tree-select
            v-model="form.parentId"
            :data="parentOptions"
            :props="{label: 'name', children: 'children', value: 'id'}"
            :render-after-expand="false"
            check-strictly
            :filter-node-method="filterParentNode"
            placeholder="请选择父级（根节点选择0）"
            style="width: 100%"
            :disabled="form.type === 1"
          >
            <template #default="{data}">
              <span class="tree-select-label">
                <el-icon v-if="getMenuIcon(data.icon)" class="menu-icon">
                  <component :is="getMenuIcon(data.icon)" />
                </el-icon>
                <span>{{ data.name }}</span>
                <el-tag v-if="data.type === 1" size="small" type="info" style="margin-left: 8px">目录</el-tag>
                <el-tag v-else-if="data.type === 2" size="small" type="success" style="margin-left: 8px">菜单</el-tag>
              </span>
            </template>
          </el-tree-select>
          <div class="form-tip" v-if="form.type === 1">目录只能存在于根节点下</div>
          <div class="form-tip" v-else-if="form.type === 2">菜单可以存在于根节点或目录下</div>
          <div class="form-tip" v-else-if="form.type === 3">按钮只能存在于菜单下</div>
        </el-form-item>
        <el-form-item :label="t('common.name')" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('common.type')" prop="type">
          <el-select v-model="form.type" style="width: 100%" @change="handleTypeChange">
            <el-option :label="t('menu.type.directory')" :value="1" />
            <el-option :label="t('menu.type.menu')" :value="2" />
            <el-option :label="t('menu.type.button')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('menu.path')" prop="path" v-if="form.type !== 3">
          <el-input v-model="form.path" :placeholder="t('menu.pathPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('menu.component')" prop="component" v-if="form.type === 2">
          <el-input v-model="form.component" :placeholder="t('menu.componentPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('menu.icon')" prop="icon">
          <el-input v-model="form.icon" :placeholder="t('menu.iconPlaceholder')" />
        </el-form-item>
        <!-- 权限编码字段已移除，改用权限-菜单关联表 -->
        <el-form-item :label="t('common.order')">
          <el-input-number v-model="form.orderNum" :min="0" />
        </el-form-item>
        <el-form-item :label="t('menu.visible')">
          <el-switch v-model="form.visible" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          {{ t('common.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {ref, reactive, onMounted, computed} from 'vue';
import {ElMessage, ElMessageBox, ElForm} from 'element-plus';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import {menuTree, menuCreate, menuUpdate, menuDelete} from '@/api/generated/admin';
import type {MenuItem, MenuCreateReq, MenuUpdateReq} from '@/api/generated/admin';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();

const treeData = ref<MenuItem[]>([]);
const loading = ref(false);

const dialogVisible = ref(false);
const isEdit = ref(false);
const formRef = ref<InstanceType<typeof ElForm>>();
const form = reactive({
  id: 0,
  parentId: 0,
  name: '',
  path: '',
  component: '',
  icon: '',
  type: 1, // 1 目录 2 菜单 3 按钮
  orderNum: 0,
  visible: 1,
  status: 1
});

// 父级选项（根据菜单类型过滤）
const parentOptions = computed(() => {
  // 添加根节点选项
  const options: MenuItem[] = [
    {
      id: 0,
      parentId: 0,
      name: '根节点',
      path: '',
      component: '',
      icon: '',
      type: 0,
      orderNum: 0,
      visible: 1,
      status: 1,
      children: []
    }
  ];

  // 根据菜单类型过滤可选的父级
  const filterTree = (items: MenuItem[], excludeId?: number): MenuItem[] => {
    return items
      .filter(item => item.id !== excludeId) // 排除自己
      .map(item => {
        let include = false;
        
        if (form.type === 1) {
          // 目录：只能选择根节点（已经在options中添加了）
          include = false;
        } else if (form.type === 2) {
          // 菜单：可以选择根节点或目录（type=1）
          include = item.type === 1;
        } else if (form.type === 3) {
          // 按钮：只能选择菜单（type=2）
          include = item.type === 2;
        }

        if (!include) {
          return null;
        }

        const children = item.children ? filterTree(item.children, excludeId) : [];
        return {
          ...item,
          children: children.length > 0 ? children : undefined
        };
      })
      .filter((item): item is MenuItem => item !== null);
  };

  const filtered = filterTree(treeData.value, isEdit.value ? form.id : undefined);
  return [...options, ...filtered];
});

// 过滤父级节点（用于搜索）
const filterParentNode = (value: string, data: MenuItem) => {
  if (!value) return true;
  return data.name.toLowerCase().includes(value.toLowerCase());
};

const rules = {
  name: [{required: true, message: t('common.nameRequired'), trigger: 'blur'}],
  type: [{required: true, message: t('common.typeRequired'), trigger: 'change'}],
  parentId: [
    {
      validator: (rule: any, value: number, callback: any) => {
        if (form.type === 1 && value !== 0) {
          callback(new Error('目录只能存在于根节点下'));
        } else if (form.type === 3 && value === 0) {
          callback(new Error('按钮只能存在于菜单下'));
        } else {
          callback();
        }
      },
      trigger: 'change'
    }
  ]
};

const getMenuIcon = (iconName?: string) => {
  if (!iconName) return null;
  const iconMap: Record<string, any> = ElementPlusIconsVue;
  // 处理 icon 名称，可能是 "ele-DataBoard" 格式，需要转换为 "DataBoard"
  const iconKey = iconName.startsWith('ele-') ? iconName.substring(4) : iconName;
  return iconMap[iconKey] || null;
};

const loadData = async () => {
  loading.value = true;
  try {
    const resp = await menuTree();
    treeData.value = resp.list || [];
  } catch (err: any) {
    ElMessage.error(err.message || t('common.loadFailed'));
  } finally {
    loading.value = false;
  }
};

const openCreate = (parent?: MenuItem) => {
  isEdit.value = false;
  Object.assign(form, {
    id: 0,
    parentId: parent?.id || 0,
    name: '',
    path: '',
    component: '',
    icon: '',
    type: 1,
    orderNum: 0,
    visible: 1,
    status: 1
  });
  dialogVisible.value = true;
};

const openEdit = (data: MenuItem) => {
  isEdit.value = true;
  Object.assign(form, {
    id: data.id,
    parentId: data.parentId,
    name: data.name,
    path: data.path || '',
    component: data.component || '',
    icon: data.icon || '',
    type: data.type,
    orderNum: data.orderNum,
    visible: data.visible,
    status: data.status
  });
  dialogVisible.value = true;
};

const submitLoading = ref(false);

const handleSubmit = () => {
  formRef.value?.validate(async (valid) => {
    if (!valid) return;
    submitLoading.value = true;
    try {
      if (isEdit.value) {
        await menuUpdate(form as MenuUpdateReq);
        ElMessage.success(t('common.updateSuccess'));
      } else {
        await menuCreate(form as MenuCreateReq);
        ElMessage.success(t('common.createSuccess'));
      }
      dialogVisible.value = false;
      loadData();
    } catch (err: any) {
      ElMessage.error(err.message || t('common.submitFailed'));
    } finally {
      submitLoading.value = false;
    }
  });
};

const handleTypeChange = () => {
  // 当菜单类型改变时，自动调整parentId
  if (form.type === 1) {
    // 目录只能存在于根节点
    form.parentId = 0;
  } else if (form.type === 3 && form.parentId === 0) {
    // 按钮不能存在于根节点，如果没有有效的父级，清空parentId让用户选择
    // 这里不清空，让用户手动选择
  }
  // 重新验证
  formRef.value?.validateField('parentId');
};

const handleDelete = (data: MenuItem) => {
  ElMessageBox.confirm(t('common.confirmDelete'), t('common.confirm'), {type: 'warning'})
    .then(async () => {
      await menuDelete({id: data.id});
      ElMessage.success(t('common.deleteSuccess'));
      loadData();
    })
    .catch(() => {});
};

onMounted(loadData);
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.toolbar {
  margin-bottom: 8px;
}
.menu-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.menu-icon {
  font-size: 16px;
}
.menu-name {
  margin-right: 8px;
}
.ops {
  margin-left: 12px;
  display: inline-flex;
  gap: 6px;
}
.tree-select-label {
  display: flex;
  align-items: center;
  gap: 4px;
}
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>

