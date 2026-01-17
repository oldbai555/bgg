<template>
  <div class="home-page">
    <div class="container">
      <h1 class="title">欢迎</h1>
      <p class="subtitle">博客与视频内容展示</p>
      <div class="links">
        <NuxtLink to="/blog" class="link-button">📝 博客文章</NuxtLink>
        <NuxtLink to="/videos" class="link-button">🎬 视频内容</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 定义页面元数据
definePageMeta({
  layout: false
})

// 也可以选择自动重定向到博客列表页
// onMounted(() => {
//   navigateTo('/blog')
// })
</script>

<style scoped lang="scss">
.home-page {
  height: 100vh; // 固定高度，不使用 min-height，避免出现滚动条
  overflow: hidden; // 隐藏最外层滚动条
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fff7e6 0%, #ffe9d9 45%, #ffd1a4 100%);

  .container {
    text-align: center;
    padding: 40px 20px;

    .title {
      font-size: 48px;
      font-weight: 700;
      color: #333;
      margin-bottom: 16px;
    }

    .subtitle {
      font-size: 20px;
      color: #666;
      margin-bottom: 40px;
    }

    .links {
      display: flex;
      gap: 20px;
      justify-content: center;
      flex-wrap: wrap;

      .link-button {
        display: inline-block;
        padding: 16px 32px;
        background: #fff;
        border-radius: 12px;
        text-decoration: none;
        color: #333;
        font-size: 18px;
        font-weight: 500;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        transition: all 0.3s;

        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
          background: #f5f5f5;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .home-page {
    .container {
      .title {
        font-size: 36px;
      }

      .subtitle {
        font-size: 18px;
      }

      .links {
        flex-direction: column;
        align-items: center;

        .link-button {
          width: 200px;
        }
      }
    }
  }
}
</style>
