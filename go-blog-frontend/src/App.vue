<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink, RouterView } from 'vue-router'
import { AUTH_CHANGED_EVENT, clearToken, getCurrentUsername, hasToken } from './utils/auth'

const router = useRouter()
const isLoggedIn = ref(hasToken())
const username = ref(getCurrentUsername())

function syncAuthState() {
  isLoggedIn.value = hasToken()
  username.value = getCurrentUsername()
}

function logout() {
  clearToken()
  router.push('/')
}

onMounted(() => {
  window.addEventListener(AUTH_CHANGED_EVENT, syncAuthState)
})

onUnmounted(() => {
  window.removeEventListener(AUTH_CHANGED_EVENT, syncAuthState)
})
</script>

<template>
  <div class="app-shell">
    <header class="site-header">
      <RouterLink class="brand" to="/">
        <span class="brand-mark">G</span>
        <span>
          <strong>Go Blog</strong>
          <small>Vue + Gin</small>
        </span>
      </RouterLink>

      <nav class="site-nav">
        <RouterLink to="/">首页</RouterLink>
        <RouterLink to="/create">写文章</RouterLink>
        <RouterLink v-if="!isLoggedIn" to="/login">登录</RouterLink>
        <button v-else type="button" @click="logout">退出</button>
      </nav>

      <div v-if="isLoggedIn && username" class="account-box">
        <span>{{ username }}</span>
      </div>
    </header>

    <main class="site-main">
      <RouterView />
    </main>
  </div>
</template>
