<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login, register } from '../api/auth'
import { setToken } from '../utils/auth'

const route = useRoute()
const router = useRouter()

const mode = ref('login')
const username = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

async function submitForm() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    if (mode.value === 'register') {
      await register(username.value, password.value)
      successMessage.value = '注册成功，正在为你登录...'
    }

    const result = await login(username.value, password.value)
    setToken(result.token)
    router.push(route.query.redirect || '/')
  } catch (error) {
    errorMessage.value = error.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="auth-layout">
    <div class="auth-copy">
      <p class="eyebrow">Account</p>
      <h1>{{ mode === 'login' ? '欢迎回来' : '创建账号' }}</h1>
      <p>登录后可以发布文章、上传封面，并管理自己的内容。</p>
    </div>

    <form class="form-panel" @submit.prevent="submitForm">
      <div class="segmented-control">
        <button type="button" :class="{ active: mode === 'login' }" @click="mode = 'login'">
          登录
        </button>
        <button type="button" :class="{ active: mode === 'register' }" @click="mode = 'register'">
          注册
        </button>
      </div>

      <label class="field">
        <span>用户名</span>
        <input v-model.trim="username" type="text" autocomplete="username" placeholder="imai" />
      </label>

      <label class="field">
        <span>密码</span>
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="至少 6 位"
        />
      </label>

      <p v-if="errorMessage" class="form-message error-text">{{ errorMessage }}</p>
      <p v-if="successMessage" class="form-message success-text">{{ successMessage }}</p>

      <button class="submit-button" type="submit" :disabled="loading">
        {{ loading ? '处理中...' : mode === 'login' ? '登录' : '注册并登录' }}
      </button>
    </form>
  </section>
</template>
