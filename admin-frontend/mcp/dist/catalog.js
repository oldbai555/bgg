/** 项目 UI 组件元数据 SSOT */
export const COMPONENT_CATALOG = [
    {
        name: 'D2Table',
        category: 'common',
        path: 'src/components/common/D2Table.vue',
        summary: '基于 el-table 的列表+分页+详情/编辑/新增抽屉封装，标准 CRUD 页面首选',
        readme: 'src/components/common/README.md',
        relatedTypes: ['src/types/table.ts'],
    },
    {
        name: 'ImageUpload',
        category: 'common',
        path: 'src/components/common/ImageUpload.vue',
        summary: '图片上传组件，配合 baseUrl 与后端 /upload 接口',
    },
    {
        name: 'IcpFooter',
        category: 'common',
        path: 'src/components/common/IcpFooter.vue',
        summary: 'ICP 备案页脚，公共展示页底部必挂',
    },
    {
        name: 'MetricReporter',
        category: 'common',
        path: 'src/components/common/MetricReporter.vue',
        summary: '公共页埋点上报封装，列表/详情页统一接入',
    },
    {
        name: 'NoticeReader',
        category: 'common',
        path: 'src/components/common/NoticeReader.vue',
        summary: '公告阅读器组件',
    },
    {
        name: 'TaskFloatBall',
        category: 'common',
        path: 'src/components/common/TaskFloatBall.vue',
        summary: '异步任务浮球，导出/下载类操作走 admin_task 时使用',
    },
    {
        name: 'AppHeader',
        category: 'layout',
        path: 'src/components/layout/AppHeader.vue',
        summary: '后台顶栏：Logo、折叠、消息、用户菜单',
    },
    {
        name: 'AppSidebar',
        category: 'layout',
        path: 'src/components/layout/AppSidebar.vue',
        summary: '后台侧边栏，渲染动态菜单路由',
    },
    {
        name: 'Breadcrumb',
        category: 'layout',
        path: 'src/components/layout/Breadcrumb.vue',
        summary: '面包屑导航',
    },
    {
        name: 'MessageNotification',
        category: 'layout',
        path: 'src/components/layout/MessageNotification.vue',
        summary: '顶栏消息通知下拉',
    },
    {
        name: 'PageHeader',
        category: 'layout',
        path: 'src/components/layout/PageHeader.vue',
        summary: '页面标题区（标题 + 操作按钮槽位）',
    },
    {
        name: 'UserMenu',
        category: 'layout',
        path: 'src/components/layout/UserMenu.vue',
        summary: '用户下拉菜单（个人中心、退出登录）',
    },
    {
        name: 'BlogAuthorCard',
        category: 'blog',
        path: 'src/components/blog/BlogAuthorCard.vue',
        summary: '博客作者信息卡片',
    },
    {
        name: 'BlogCategoryNav',
        category: 'blog',
        path: 'src/components/blog/BlogCategoryNav.vue',
        summary: '博客分类导航',
    },
    {
        name: 'BlogHeader',
        category: 'blog',
        path: 'src/components/blog/BlogHeader.vue',
        summary: '博客页头区域',
    },
    {
        name: 'BlogSocialLinks',
        category: 'blog',
        path: 'src/components/blog/BlogSocialLinks.vue',
        summary: '博客社交链接展示',
    },
    {
        name: 'BlogTOC',
        category: 'blog',
        path: 'src/components/blog/BlogTOC.vue',
        summary: '博客文章目录（TOC）',
    },
];
export function findComponent(name) {
    const normalized = name.trim();
    return COMPONENT_CATALOG.find((c) => c.name.toLowerCase() === normalized.toLowerCase());
}
export function listByCategory(category) {
    if (!category) {
        return COMPONENT_CATALOG;
    }
    const cat = category.trim().toLowerCase();
    return COMPONENT_CATALOG.filter((c) => c.category === cat);
}
