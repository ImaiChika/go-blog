<script setup>
import { ref } from 'vue'
import { changePassword } from '../api/auth'

const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

async function submitForm() {
  errorMessage.value = ''
  successMessage.value = ''

  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = '两次输入的新密码不一致'
    return
  }

  if (oldPassword.value === newPassword.value) {
    errorMessage.value = '新密码不能和旧密码相同'
    return
  }

  loading.value = true
  try {
    const result = await changePassword(oldPassword.value, newPassword.value)
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
    successMessage.value = result.message || '密码修改成功'
  } catch (error) {
    errorMessage.value = error.message || '密码修改失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="auth-layout">
    <div class="auth-copy">
      <p class="eyebrow">Security</p>
      <h1>修改密码</h1>
      <p>为了保护账号安全，修改密码前需要先输入当前密码。修改成功后，下次登录请使用新密码。</p>
    </div>

    <form class="form-panel" @submit.prevent="submitForm">
      <label class="field">
        <span>当前密码</span>
        <input
          v-model="oldPassword"
          type="password"
          autocomplete="current-password"
          placeholder="请输入当前密码"
        />
      </label>

      <label class="field">
        <span>新密码</span>
        <input
          v-model="newPassword"
          type="password"
          autocomplete="new-password"
          placeholder="至少 6 位"
        />
      </label>

      <label class="field">
        <span>确认新密码</span>
        <input
          v-model="confirmPassword"
          type="password"
          autocomplete="new-password"
          placeholder="再输入一次新密码"
        />
      </label>

      <p v-if="errorMessage" class="form-message error-text">{{ errorMessage }}</p>
      <p v-if="successMessage" class="form-message success-text">{{ successMessage }}</p>

      <button class="submit-button" type="submit" :disabled="loading">
        {{ loading ? '修改中...' : '确认修改' }}
      </button>
    </form>
  </section>
</template>
