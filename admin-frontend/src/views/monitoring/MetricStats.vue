<template>
  <div class="page">
    <!-- 搜索表单 -->
    <el-card class="mb-12">
      <el-form :inline="true" :model="query" class="search-form">
        <el-form-item label="业务模块">
          <el-select
            v-model="query.module"
            placeholder="请选择模块"
            clearable
            style="width: 200px"
          >
            <el-option label="博客文章列表" value="blog_article_list" />
            <el-option label="博客文章详情" value="blog_article_detail" />
            <el-option label="视频列表" value="video_list" />
            <el-option label="视频详情" value="video_detail" />
          </el-select>
        </el-form-item>
        <el-form-item label="业务ID">
          <el-input-number
            v-model="query.bizId"
            :min="0"
            placeholder="业务ID（0表示列表页）"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="统计日期">
          <el-date-picker
            v-model="query.day"
            type="date"
            format="YYYYMMDD"
            value-format="YYYYMMDD"
            placeholder="选择日期（默认今天）"
            style="width: 180px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadData">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计结果 -->
    <el-card>
      <div v-if="stats" class="stats-result">
        <div class="stats-header">
          <h3>统计数据</h3>
          <div class="stats-meta">
            <span>模块：{{ stats.module }}</span>
            <span>业务ID：{{ stats.bizId }}</span>
            <span>日期：{{ formatDay(stats.day) }}</span>
          </div>
        </div>
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-label">PV（页面访问量）</div>
            <div class="stat-value">{{ stats.pv }}</div>
            <div class="stat-desc">累计访问次数</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">UV（独立访客）</div>
            <div class="stat-value">{{ stats.uv }}</div>
            <div class="stat-desc">基于 IP + User-Agent 去重</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">VV（访问次数）</div>
            <div class="stat-value">{{ stats.vv }}</div>
            <div class="stat-desc">访问次数统计</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">IP（独立IP数）</div>
            <div class="stat-value">{{ stats.ip }}</div>
            <div class="stat-desc">基于 IP 地址去重</div>
          </div>
        </div>

        <!-- 简单图表展示：当前 PV/UV/VV/IP 的对比柱状图 -->
        <div class="chart-wrapper">
          <h4 class="chart-title">访问数据对比图</h4>
          <div class="chart-bars">
            <div
              v-for="item in chartData"
              :key="item.key"
              class="chart-bar-item"
            >
              <div class="chart-bar">
                <div
                  class="chart-bar-inner"
                  :style="{ height: item.percent + '%', background: item.color }"
                ></div>
              </div>
              <div class="chart-bar-label">{{ item.label }}</div>
              <div class="chart-bar-value">{{ item.value }}</div>
            </div>
          </div>
        </div>
      </div>
      <div v-else-if="!loading" class="empty-tip">
        请选择业务模块和日期后点击查询
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {computed, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {monitoringApi} from '@/api/monitoring'
import type {MetricStatsReq, MetricStatsResp} from '@/api/generated/admin'

const loading = ref(false)
const stats = ref<MetricStatsResp | null>(null)

const query = reactive<MetricStatsReq>({
  module: '',
  bizId: 0,
  day: ''
})

const formatDay = (day: string): string => {
  if (!day || day.length !== 8) {
return day
}
  return `${day.substring(0, 4)}-${day.substring(4, 6)}-${day.substring(6, 8)}`
}

// 图表数据：基于当前返回的 PV/UV/VV/IP 生成简单对比柱状图
const chartData = computed(() => {
  if (!stats.value) {
return []
}
  const items = [
    {key: 'pv', label: 'PV', value: stats.value.pv, color: '#67c23a'},
    {key: 'uv', label: 'UV', value: stats.value.uv, color: '#409eff'},
    {key: 'vv', label: 'VV', value: stats.value.vv, color: '#e6a23c'},
    {key: 'ip', label: 'IP', value: stats.value.ip, color: '#f56c6c'}
  ]
  const max = Math.max(...items.map((i) => i.value), 1) // 避免全 0 时除以 0
  return items.map((i) => ({
    ...i,
    percent: (i.value / max) * 100
  }))
})

const loadData = async () => {
  if (!query.module) {
    ElMessage.warning('请选择业务模块')
    return
  }

  loading.value = true
  try {
    const resp = await monitoringApi.metricStats(query)
    stats.value = resp
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '查询失败'
    ElMessage.error(message)
    stats.value = null
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  query.module = ''
  query.bizId = 0
  query.day = ''
  stats.value = null
}
</script>

<style scoped lang="scss">
.page {
  padding: 16px 24px;
}

.search-form {
  .el-form-item {
    margin-bottom: 0;
  }
}

.stats-result {
  .stats-header {
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid #e4e7ed;

    h3 {
      margin: 0 0 12px;
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }

    .stats-meta {
      display: flex;
      gap: 24px;
      font-size: 14px;
      color: #606266;

      span {
        display: inline-flex;
        align-items: center;
      }
    }
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
  }

  .stat-card {
    padding: 20px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    color: #fff;
    text-align: center;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
    transition: transform 0.2s, box-shadow 0.2s;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
    }

    .stat-label {
      font-size: 14px;
      opacity: 0.9;
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 32px;
      font-weight: 700;
      margin-bottom: 8px;
      line-height: 1.2;
    }

    .stat-desc {
      font-size: 12px;
      opacity: 0.8;
    }
  }
}

.chart-wrapper {
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px dashed #e4e7ed;

  .chart-title {
    margin: 0 0 16px;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }

  .chart-bars {
    display: flex;
    align-items: flex-end;
    gap: 24px;
    padding: 8px 4px 0;
    height: 220px;
  }

  .chart-bar-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: 13px;
    color: #606266;
  }

  .chart-bar {
    width: 40px;
    flex: 1;
    display: flex;
    align-items: flex-end;
    justify-content: center;
    padding: 4px 0;
    border-radius: 999px;
    background: #f5f7fa;
    box-shadow: inset 0 0 0 1px #ebeef5;
    margin-bottom: 8px;
  }

  .chart-bar-inner {
    width: 24px;
    border-radius: 999px;
    transition: height 0.3s ease;
  }

  .chart-bar-label {
    font-weight: 600;
    margin-bottom: 4px;
  }

  .chart-bar-value {
    font-family: Menlo, Monaco, Consolas, 'Courier New', monospace;
  }
}

.empty-tip {
  text-align: center;
  color: #909399;
  padding: 60px 0;
  font-size: 14px;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stats-header .stats-meta {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
