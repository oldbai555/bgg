<template>
  <div class="dashboard">
    <el-card class="dashboard__card">
      <template #header>
        <div class="card-header">
          <h3>{{ t('common.welcome') }}</h3>
        </div>
      </template>
      <el-calendar v-model="selectedDate">
        <template #date-cell="{ data }">
          <el-tooltip
            class="item"
            effect="dark"
            :content="getDailyQuote(data.day)"
            placement="top-start"
          >
            <div :class="data.isSelected ? 'is-selected' : ''">
              <p class="date-text">
                {{ getDateText(data.day) }}
                <span v-if="data.isSelected" class="selected-icon">✔️</span>
              </p>
              <p class="daily-quote">
                {{ truncateText(getDailyQuote(data.day), 30) }}
              </p>
            </div>
          </el-tooltip>
        </template>
      </el-calendar>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { dailyShortSentenceList } from '@/api/generated/admin';
import type { DailyShortSentenceItem } from '@/api/generated/adminComponents';
import { ElMessage } from 'element-plus';

const { t } = useI18n();

const selectedDate = ref<Date>(new Date());
const quotes = ref<string[]>([]);

// 加载短句列表
const loadQuotes = async () => {
  try {
    const resp = await dailyShortSentenceList({
      page: 1,
      pageSize: 2000, // 获取所有短句
    });

    if (resp && resp.list) {
      quotes.value = resp.list.map((item: DailyShortSentenceItem) => item.content);
    }
  } catch (error) {
    console.error('加载每日短句失败:', error);
    ElMessage.error('加载每日短句失败');
  }
};

// 获取日期文本（只显示日期部分）
const getDateText = (dateStr: string): string => {
  const parts = dateStr.split('-');
  if (parts.length >= 3) {
    return `${parts[1]}-${parts[2]}`;
  }
  return dateStr;
};

// 根据日期获取对应的短句
const getDailyQuote = (date: string): string => {
  if (quotes.value.length === 0) {
    return '';
  }

  try {
    // 计算当前日期是今年的第几天
    const dateObj = new Date(date);
    const startOfYear = new Date(dateObj.getFullYear(), 0, 1);
    const diffTime = dateObj.getTime() - startOfYear.getTime();
    const dayOfYear = Math.ceil(diffTime / (24 * 60 * 60 * 1000));

    // 调整为从0开始的索引
    const selectedDay = dayOfYear - 1;
    const quotesLength = quotes.value.length;

    // 使用模运算循环获取短句
    return quotes.value[selectedDay % quotesLength] || '';
  } catch (error) {
    console.error('计算日期短句失败:', error);
    return '';
  }
};

// 截断文本
const truncateText = (text: string, maxLength: number): string => {
  if (typeof text !== 'string') {
    return '';
  }

  if (text.length > maxLength) {
    return text.substring(0, maxLength) + '...';
  }
  return text;
};

onMounted(() => {
  loadQuotes();
});
</script>

<style scoped lang="scss">
.dashboard {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.dashboard__card {
  flex: 1;

  .card-header {
    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 500;
    }
  }
}

.is-selected {
  color: #409eff;
}

.date-text {
  margin: 0;
  font-size: 14px;
  font-weight: 500;

  .selected-icon {
    margin-left: 4px;
  }
}

.daily-quote {
  margin: 4px 0 0 0;
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  line-height: 1.4;
}
</style>

