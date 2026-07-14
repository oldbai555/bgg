/**
 * 字典管理 Store
 * 统一管理全局字典数据，在登录后一次性加载所有需要的字典
 */
import {defineStore} from 'pinia';
import {systemApi} from '@/api/system';
import type {DictItemItem} from '@/api/generated/admin';

interface DictState {
  // 字典数据，key为字典类型编码，value为字典项列表
  dicts: Record<string, DictItemItem[]>;
  // 是否已加载
  loaded: boolean;
  // 加载时间戳
  loadAt: number;
}

// 需要批量获取的字典类型编码列表（在登录后一次性获取）
export const REQUIRED_DICT_CODES = [
  'performance_log_slow_status', // 性能日志慢查询状态
  'menu_type', // 菜单类型
  'read_status', // 已读状态
  'operation_type', // 操作类型
  'http_method', // HTTP请求方法
  'notice_type', // 公告类型
  'notice_status', // 公告状态
  'login_status', // 登录状态
  'audit_type', // 审计类型
  'chat_message_type', // 消息类型
  'daily_short_sentence_type', // 短句类型
  'notification_source_type', // 消息来源类型
  'storage_base_url', // 存储baseURL（配置）
  'websocket_base_url', // WebSocket baseURL（配置）
  'sdk_http_method', // SDK 接口管理-HTTP方法
  'sdk_status', // SDK Key/接口状态
  'task_type', // 任务类型
  'task_execution_type', // 任务执行类型
  'task_status', // 任务状态
  'task_config', // 任务配置（如最近任务数量）
  // 博客相关字典
  'blog_article_status', // 博客文章业务状态
  'blog_article_audit_status', // 博客文章审核状态
  'blog_tag_status', // 博客标签状态
  'blog_friend_link_status', // 友情链接状态
  'blog_social_info_status' // 社交信息状态
] as const;

export const useDictStore = defineStore('dict', {
  state: (): DictState => ({
    dicts: {},
    loaded: false,
    loadAt: 0
  }),

  getters: {
    /**
     * 根据字典类型编码获取字典项列表
     */
    getDictItems: (state) => (code: string): DictItemItem[] => {
      return state.dicts[code] || [];
    },

    /**
     * 根据字典类型编码和值获取标签
     */
    getDictLabel: (state) => (code: string, value: string | number): string => {
      const items = state.dicts[code] || [];
      const item = items.find((item) => item.value === String(value));
      return item ? item.label : String(value);
    },

    /**
     * 根据字典类型编码获取选项列表（用于el-select等组件）
     */
    getDictOptions: (state) => (code: string): Array<{label: string; value: string | number}> => {
      const items = state.dicts[code] || [];
      return items.map((item) => ({
        label: item.label,
        value: item.value
      }));
    }
  },

  actions: {
    /**
     * 批量加载字典数据
     * @param codes 需要加载的字典类型编码列表，如果不传则加载所有必需的字典
     */
    async loadDicts(codes?: string[]) {
      try {
        const codesToLoad = codes || [...REQUIRED_DICT_CODES];
        const resp = await systemApi.dictBatchGet({codes: codesToLoad});
        
        if (resp && resp.dicts) {
          // 更新字典数据
          Object.keys(resp.dicts).forEach((code) => {
            this.dicts[code] = resp.dicts[code].items || [];
          });
          
          this.loaded = true;
          this.loadAt = Date.now();
        }
      } catch (err) {
        console.error('批量加载字典失败:', err);
        throw err;
      }
    },

    /**
     * 刷新字典数据
     */
    async refreshDicts(codes?: string[]) {
      this.loaded = false;
      await this.loadDicts(codes);
    },

    /**
     * 清除字典数据
     */
    clearDicts() {
      this.dicts = {};
      this.loaded = false;
      this.loadAt = 0;
    }
  }
});

