<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { deletePost, getPostById } from '../api/post'
import { getCurrentUsername, hasToken } from '../utils/auth'

const route = useRoute()
const router = useRouter()

const post = ref(null)
const loading = ref(false)
const deleting = ref(false)
const errorMessage = ref('')

const isOwner = computed(() => {
  return hasToken() && post.value?.author === getCurrentUsername()
})

function formatDate(value) {
  return new Date(value).toLocaleString('zh-CN')
}

async function fetchPost() {
  loading.value = true
  errorMessage.value = ''

  try {
    const result = await getPostById(route.params.id)
    post.value = result.data
  } catch (error) {
    errorMessage.value = error.message || '加载文章失败'
  } finally {
    loading.value = false
  }
}

async function handleDelete() {
  const confirmed = window.confirm('确定要删除这篇文章吗？')
  if (!confirmed) {
    return
  }

  deleting.value = true
  errorMessage.value = ''

  try {
    await deletePost(route.params.id)
    router.push('/')
  } catch (error) {
    errorMessage.value = error.message || '删除文章失败'
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchPost()
})
</script>

<template>
  <section class="page">
    <p v-if="loading" class="state-text">正在加载文章...</p>
    <p v-else-if="errorMessage" class="state-text error-text">{{ errorMessage }}</p>

    <article v-else-if="post" class="post-detail">
      <img
        v-if="post.cover_image"
        :src="post.cover_image"
        :alt="post.title"
        class="detail-cover"
        @error="post.cover_image = ''"
      />

      <div class="detail-heading">
        <div>
          <p class="post-meta">作者 {{ post.author }} · {{ formatDate(post.CreatedAt) }}</p>
          <h1>{{ post.title }}</h1>
        </div>

        <div v-if="isOwner" class="detail-actions">
          <RouterLink class="secondary-action" :to="`/posts/${post.ID}/edit`">编辑</RouterLink>
          <button class="danger-action" type="button" :disabled="deleting" @click="handleDelete">
            {{ deleting ? '删除中...' : '删除' }}
          </button>
        </div>
      </div>

      <div class="article-content">
        {{ post.content }}
      </div>
    </article>
  </section>
</template>
