<template>
  <div class="home-page">
    <div class="container">
      <h1 class="title">欢迎</h1>
      <p class="subtitle">博客与视频内容展示</p>
      <div class="links">
        <router-link to="/blog" class="link-button">📝 博客</router-link>
        <router-link to="/videos" class="link-button">🎬 视频</router-link>
        <router-link to="/login" class="link-button">🔐 后台</router-link>
      </div>
    </div>
    <IcpFooter />
  </div>
</template>

<script setup lang="ts">
import IcpFooter from '@/components/common/IcpFooter.vue'
</script>

<style scoped lang="scss">
.home-page {
  height: 100vh; // 固定高度，避免出现滚动条
  overflow: hidden; // 隐藏滚动条
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fff7e6 0%, #ffe9d9 45%, #ffd1a4 100%);
  position: relative;
  box-sizing: border-box;

  .container {
    text-align: center;
    padding: 10px 20px; // 减小上下内边距
    max-width: 500px; // 进一步限制最大宽度
    width: 100%;
    margin-bottom: 40px; // 为备案号留出空间

    .title {
      font-size: 32px; // 进一步减小标题字体
      font-weight: 700;
      color: #333;
      margin: 0 0 8px 0; // 减小边距
      line-height: 1.2;
    }

    .subtitle {
      font-size: 14px; // 进一步减小副标题字体
      color: #666;
      margin: 0 0 24px 0; // 减小边距
      line-height: 1.4;
    }

    .links {
      display: flex;
      gap: 12px; // 进一步减小间距
      justify-content: center;
      flex-wrap: wrap;

      .link-button {
        display: inline-block;
        padding: 10px 20px; // 进一步减小按钮内边距
        background: #fff;
        border-radius: 6px; // 减小圆角
        text-decoration: none;
        color: #333;
        font-size: 14px; // 进一步减小字体
        font-weight: 500;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1); // 减小阴影
        transition: all 0.3s;
        line-height: 1.4;

        &:hover {
          transform: translateY(-2px); // 减小悬停位移
          box-shadow: 0 4px 10px rgba(0, 0, 0, 0.15);
          background: #f5f5f5;
        }
      }
    }
  }
}

// 备案号组件样式
.home-page {
  :deep(.icp-footer) {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100%;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    padding: 8px 16px; // 减小备案号内边距
    font-size: 11px; // 减小备案号字体
    line-height: 1.2;
  }
}

@include mobile {
  .home-page {
    .container {
      padding: 8px 16px; // 减小移动端内边距
      margin-bottom: 35px; // 减小移动端底部间距

      .title {
        font-size: 24px; // 移动端更小的标题
        margin-bottom: 6px;
      }

      .subtitle {
        font-size: 12px; // 移动端更小的副标题
        margin-bottom: 20px;
      }

      .links {
        flex-direction: column;
        align-items: center;
        gap: 10px; // 减小移动端按钮间距

        .link-button {
          width: 160px; // 移动端按钮宽度
          padding: 8px 16px; // 减小移动端按钮内边距
          font-size: 13px;
        }
      }
    }

    :deep(.icp-footer) {
      padding: 6px 12px; // 移动端备案号内边距
      font-size: 10px; // 移动端备案号字体
    }
  }
}
</style>
