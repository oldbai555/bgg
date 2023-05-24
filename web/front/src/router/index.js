import Vue from 'vue'
import VueRouter from 'vue-router'

const ArticleList = () =>
  import(/* webpackChunkName: "group-index" */ '../components/articleList.vue')
const Detail = () =>
  import(/* webpackChunkName: "group-detail" */ '../components/details.vue')
const Category = () =>
  import(/* webpackChunkName: "group-category" */ '../components/cateList.vue')
const Search = () =>
  import(/* webpackChunkName: "group-search" */ '../components/search.vue')

Vue.use(VueRouter)

//获取原型对象上的push函数
const originalPush = VueRouter.prototype.push
//修改原型对象中的push方法
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}

const routes = [
  { path: '/', component: ArticleList, meta: { title: 'LB小破站' } },
  {
    path: '/article/detail/:id',
    component: Detail,
    meta: { title: window.sessionStorage.getItem('title') },
    props: true
  },
  {
    path: '/category/:cid',
    component: Category,
    meta: { title: '分类信息' },
    props: true
  },
  {
    path: '/search/:title',
    component: Search,
    meta: { title: '搜索结果' },
    props: true
  }
]

const router = new VueRouter({
  mode: 'hash',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title ? to.meta.title : '加载中'
  }
  next()
})

export default router
