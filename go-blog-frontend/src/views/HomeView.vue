<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { getPostList } from '../api/post'

const posts = ref([])
const page = ref(1)
const pageSize = 6
const total = ref(0)
const loading = ref(false)
const errorMessage = ref('')

function formatDate(value) {
  return new Date(value).toLocaleString('zh-CN')
}

function excerpt(content) {
  if (!content) {
    return ''
  }
  return content.length > 120 ? `${content.slice(0, 120)}...` : content
}

async function fetchPosts(targetPage = page.value) {
  loading.value = true
  errorMessage.value = ''

  try {
    const result = await getPostList(targetPage, pageSize)
    posts.value = result.data || []
    total.value = result.total || 0
    page.value = result.page || targetPage
  } catch (error) {
    errorMessage.value = error.message || '加载文章失败'
  } finally {
    loading.value = false
  }
}

function goPrevious() {
  if (page.value > 1) {
    fetchPosts(page.value - 1)
  }
}

function goNext() {
  if (page.value * pageSize < total.value) {
    fetchPosts(page.value + 1)
  }
}

onMounted(() => {
  fetchPosts()
})
</script>

<template>
  <section class="page stack-page">
    <div class="page-heading">
      <div>
        <p class="eyebrow">Go Blog</p>
        <h1>最新文章</h1>
      </div>

      <RouterLink class="primary-action" to="/create">写新文章</RouterLink>
    </div>

    <p v-if="loading" class="state-text">正在加载文章...</p>
    <p v-else-if="errorMessage" class="state-text error-text">{{ errorMessage }}</p>

    <div v-else-if="posts.length === 0" class="empty-state">
      <h2>还没有文章</h2>
      <p>发布第一篇博客后，它会出现在这里。</p>
      <RouterLink class="primary-action" to="/create">开始写作</RouterLink>
    </div>

    <div v-else class="post-grid">
      <article v-for="post in posts" :key="post.ID" class="post-card">
        <RouterLink :to="`/posts/${post.ID}`" class="post-cover-link">
          <img
            v-if="post.cover_image"
            :src="post.cover_image"
            :alt="post.title"
            class="post-cover"
            @error="post.cover_image = ''"
          />
          <div v-else class="post-cover placeholder-cover">
            {{ post.title.slice(0, 1) }}
          </div>
        </RouterLink>

        <div class="post-card-body">
          <p class="post-meta">作者 {{ post.author }} · {{ formatDate(post.CreatedAt) }}</p>
          <h2>
            <RouterLink :to="`/posts/${post.ID}`">{{ post.title }}</RouterLink>
          </h2>
          <p class="post-preview">{{ excerpt(post.content) }}</p>
        </div>
      </article>
    </div>

    <div v-if="posts.length > 0" class="pager">
      <button type="button" :disabled="page <= 1 || loading" @click="goPrevious">上一页</button>
      <span>第 {{ page }} 页 / 共 {{ total }} 篇</span>
      <button type="button" :disabled="page * pageSize >= total || loading" @click="goNext">
        下一页
      </button>
    </div>
  </section>
</template>
