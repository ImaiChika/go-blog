import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import PostDetailView from '../views/PostDetailView.vue'
import CreatePostView from '../views/CreatePostView.vue'
import { hasToken } from '../utils/auth'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
  },
  {
    path: '/posts/:id',
    name: 'post-detail',
    component: PostDetailView,
  },
  {
    path: '/create',
    name: 'create-post',
    component: CreatePostView,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: '/posts/:id/edit',
    name: 'edit-post',
    component: CreatePostView,
    meta: {
      requiresAuth: true,
    },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !hasToken()) {
    return {
      path: '/login',
      query: {
        redirect: to.fullPath,
      },
    }
  }
})

export default router
