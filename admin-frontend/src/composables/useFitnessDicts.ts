import {onMounted} from 'vue'
import {useDictStore} from '@/stores/dict'
import {FITNESS_DICT_CODES} from '@/constants/fitness'

/**
 * 身材管理手机端字典：手机端用户不走后台登录收尾（不会批量加载 REQUIRED_DICT_CODES），
 * 进页面时发现身材管理字典还没加载就补拉一次（/dict/batch 只要求登录，不要求权限）
 */
export function useFitnessDicts() {
  const dictStore = useDictStore()

  const ensure = async () => {
    const missing = FITNESS_DICT_CODES.some((code) => dictStore.getDictItems(code).length === 0)
    if (!missing) {
      return
    }
    try {
      await dictStore.loadDicts([...FITNESS_DICT_CODES])
    } catch {
      // 字典拉取失败时页面回退展示原始值，不阻塞主流程
    }
  }

  onMounted(ensure)

  return {
    label: (code: string, value: string | number) => dictStore.getDictLabel(code, value),
    options: (code: string) => dictStore.getDictOptions(code)
  }
}
