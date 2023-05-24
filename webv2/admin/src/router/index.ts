import {createRouter, createWebHistory} from 'vue-router'
import Login from '../views/login.vue'
import Admin from '../views/admin.vue'
import Index from '../components/admin/index.vue'
import NotFound from '../components/admin/404.vue'
import ArticleList from '../components/article/list.vue'
import ArticleEdit from '../components/article/edit.vue'
import BillList from '../components/bill/list.vue'
import BillCategoryList from '../components/bill_category/list.vue'
import CategoryList from '../components/category/list.vue'
import CommentList from '../components/comment/list.vue'
import FileList from '../components/file/list.vue'
import UserList from '../components/user/list.vue'
import Profile from '../components/user/profile.vue'
import {clearAllCaches, getToken} from "@/plugin/utils/cache";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/login',
            name: 'login',
            meta: {
                title: '请登录'
            },
            component: Login
        },
        {
            path: '/',
            name: 'admin',
            meta: {
                title: '管理后台'
            },
            redirect: '/index',
            component: Admin,
            children: [
                {
                    path: '/index',
                    name: 'index',
                    component: Index,
                    meta: {
                        title: '管理后台'
                    }
                },
                {
                    path: '/article_list',
                    name: 'article_list',
                    component: ArticleList,
                    meta: {
                        title: '文章管理'
                    }
                },
                {
                    path: '/article_edit',
                    name: 'article_edit',
                    component: ArticleEdit,
                    meta: {
                        title: '文章操作'
                    }
                },
                {
                    path: '/bill_list',
                    name: 'bill_list',
                    component: BillList,
                    meta: {
                        title: '清单管理'
                    }
                },
                {
                    path: '/bill_category_list',
                    name: 'bill_category_list',
                    component: BillCategoryList,
                    meta: {
                        title: '分类管理'
                    }
                },
                {
                    path: '/category_list',
                    name: 'category_list',
                    component: CategoryList,
                    meta: {
                        title: '分类管理'
                    }
                },
                {
                    path: '/comment_list',
                    name: 'comment_list',
                    component: CommentList,
                    meta: {
                        title: '评论管理'
                    }
                },
                {
                    path: '/file_list',
                    name: 'file_list',
                    component: FileList,
                    meta: {
                        title: '文件管理'
                    }
                },
                {
                    path: '/user_list',
                    name: 'user_list',
                    component: UserList,
                    meta: {
                        title: '用户管理'
                    }
                },
                {
                    path: '/profile',
                    name: 'profile',
                    component: Profile,
                    meta: {
                        title: '个人中心'
                    }
                },
            ]
        },
        {
            path: '/:error*', // /:error -> 匹配 /, /one, /one/two, /one/two/three, 等
            name: '404',
            meta: {
                title: '404'
            },
            component: NotFound
        },
    ]
})

router.beforeEach((to, from, next) => {

    if (to.meta.title) {
        document.title = to.meta.title as string
    }

    if (to.path === '/login') {
        next()
        return;
    }

    const userToken = getToken()
    if (userToken) {
        next()
        return;
    }
    clearAllCaches()
    next('/login')
})

export default router
