/**
 * 字典选项管理 Composable
 * 统一管理所有 el-select 选项，从全局字典store获取，避免硬编码
 */
import {computed} from 'vue';
import {useDictStore} from '@/stores/dict';

/**
 * 从全局字典store获取选项列表
 * @param code 字典类型编码
 * @param defaultValue 默认值（如果字典中没有数据时使用）
 * @returns 选项列表，格式：{label: string, value: string | number}[]
 */
export function getDictOptions(
  code: string,
  defaultValue: Array<{label: string; value: string | number}> = []
): Array<{label: string; value: string | number}> {
  const dictStore = useDictStore();
  const options = dictStore.getDictOptions(code);
  return options.length > 0 ? options : defaultValue;
}

/**
 * 使用字典选项的 Composable
 * @param code 字典类型编码
 * @param defaultValue 默认值（如果字典中没有数据时使用）
 */
export function useDictOptions(
  code: string,
  defaultValue: Array<{label: string; value: string | number}> = []
) {
  const dictStore = useDictStore();

  // 从store获取选项（响应式）
  const options = computed(() => {
    const storeOptions = dictStore.getDictOptions(code);
    return storeOptions.length > 0 ? storeOptions : defaultValue;
  });

  // 根据值获取标签
  const getLabel = (value: string | number): string => {
    return dictStore.getDictLabel(code, value);
  };

  // 刷新字典（如果需要）
  const refresh = async () => {
    await dictStore.refreshDicts([code]);
  };

  return {
    options,
    loading: computed(() => !dictStore.loaded),
    loadOptions: refresh, // 兼容旧代码
    getLabel,
    refresh
  };
}

