<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { createComment, deleteComment, getComments } from '../api/comment'
import { deletePost, getPostById } from '../api/post'
import { getCurrentUsername, hasToken } from '../utils/auth'

const route = useRoute()
const router = useRouter()

const post = ref(null)
const loading = ref(false)
const deleting = ref(false)
const errorMessage = ref('')
const comments = ref([])
const commentPage = ref(1)
const commentPageSize = 10
const commentTotal = ref(0)
const commentsLoading = ref(false)
const commentContent = ref('')
const commentError = ref('')
const submittingComment = ref(false)
const deletingCommentId = ref(null)

const isOwner = computed(() => {
  return hasToken() && post.value?.author === getCurrentUsername()
})

const isLoggedIn = computed(() => hasToken())

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

async function fetchComments(targetPage = commentPage.value) {
  commentsLoading.value = true
  commentError.value = ''

  try {
    const result = await getComments(route.params.id, targetPage, commentPageSize)
    comments.value = result.data || []
    commentTotal.value = result.total || 0
    commentPage.value = result.page || targetPage
  } catch (error) {
    commentError.value = error.message || '加载评论失败'
  } finally {
    commentsLoading.value = false
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

function canDeleteComment(comment) {
  const username = getCurrentUsername()
  return hasToken() && (comment.author === username || post.value?.author === username)
}

async function submitComment() {
  const content = commentContent.value.trim()
  if (!content) {
    commentError.value = '评论内容不能为空'
    return
  }

  submittingComment.value = true
  commentError.value = ''

  try {
    await createComment(route.params.id, content)
    commentContent.value = ''
    await fetchComments(1)
  } catch (error) {
    commentError.value = error.message || '发表评论失败'
  } finally {
    submittingComment.value = false
  }
}

async function handleDeleteComment(comment) {
  const confirmed = window.confirm('确定要删除这条评论吗？')
  if (!confirmed) {
    return
  }

  deletingCommentId.value = comment.ID
  commentError.value = ''

  try {
    await deleteComment(comment.ID)
    await fetchComments(commentPage.value)
  } catch (error) {
    commentError.value = error.message || '删除评论失败'
  } finally {
    deletingCommentId.value = null
  }
}

function goPreviousComments() {
  if (commentPage.value > 1) {
    fetchComments(commentPage.value - 1)
  }
}

function goNextComments() {
  if (commentPage.value * commentPageSize < commentTotal.value) {
    fetchComments(commentPage.value + 1)
  }
}

onMounted(() => {
  fetchPost()
  fetchComments()
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

    <section v-if="post" class="comments-section">
      <div class="comments-heading">
        <div>
          <p class="eyebrow">Comments</p>
          <h2>评论</h2>
        </div>
        <span>{{ commentTotal }} 条</span>
      </div>

      <form v-if="isLoggedIn" class="comment-form" @submit.prevent="submitComment">
        <label class="field">
          <span>写评论</span>
          <textarea
            v-model="commentContent"
            rows="4"
            maxlength="500"
            placeholder="写下你的想法"
          ></textarea>
        </label>

        <button class="submit-button" type="submit" :disabled="submittingComment">
          {{ submittingComment ? '发送中...' : '发表评论' }}
        </button>
      </form>

      <div v-else class="login-tip">
        <RouterLink to="/login">登录后发表评论</RouterLink>
      </div>

      <p v-if="commentError" class="form-message error-text">{{ commentError }}</p>
      <p v-if="commentsLoading" class="state-text">正在加载评论...</p>

      <div v-else-if="comments.length === 0" class="empty-comments">
        还没有评论。
      </div>

      <div v-else class="comment-list">
        <article v-for="comment in comments" :key="comment.ID" class="comment-item">
          <div>
            <p class="comment-meta">{{ comment.author }} · {{ formatDate(comment.CreatedAt) }}</p>
            <p class="comment-content">{{ comment.content }}</p>
          </div>

          <button
            v-if="canDeleteComment(comment)"
            class="text-danger-button"
            type="button"
            :disabled="deletingCommentId === comment.ID"
            @click="handleDeleteComment(comment)"
          >
            {{ deletingCommentId === comment.ID ? '删除中...' : '删除' }}
          </button>
        </article>
      </div>

      <div v-if="comments.length > 0" class="pager">
        <button
          type="button"
          :disabled="commentPage <= 1 || commentsLoading"
          @click="goPreviousComments"
        >
          上一页
        </button>
        <span>第 {{ commentPage }} 页 / 共 {{ commentTotal }} 条</span>
        <button
          type="button"
          :disabled="commentPage * commentPageSize >= commentTotal || commentsLoading"
          @click="goNextComments"
        >
          下一页
        </button>
      </div>
    </section>
  </section>
</template>
