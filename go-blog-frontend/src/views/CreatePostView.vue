<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createPost, getPostById, updatePost } from '../api/post'
import { uploadImage } from '../api/upload'
import { getCurrentUsername } from '../utils/auth'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => Boolean(route.params.id))
const pageTitle = computed(() => (isEdit.value ? '编辑文章' : '发布文章'))

const title = ref('')
const content = ref('')
const coverImage = ref('')
const coverFile = ref(null)
const coverPreview = ref('')
const loading = ref(false)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

function handleFileChange(event) {
  const file = event.target.files?.[0]
  coverFile.value = file || null
  coverPreview.value = file ? URL.createObjectURL(file) : ''
}

async function loadPostForEdit() {
  if (!isEdit.value) {
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const result = await getPostById(route.params.id)
    const post = result.data

    if (post.author !== getCurrentUsername()) {
      errorMessage.value = '只能编辑自己发布的文章'
      return
    }

    title.value = post.title
    content.value = post.content
    coverImage.value = post.cover_image || ''
  } catch (error) {
    errorMessage.value = error.message || '加载文章失败'
  } finally {
    loading.value = false
  }
}

async function submitPost() {
  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    let finalCoverImage = coverImage.value
    if (coverFile.value) {
      const uploadResult = await uploadImage(coverFile.value)
      finalCoverImage = uploadResult.url
    }

    const payload = {
      title: title.value,
      content: content.value,
      cover_image: finalCoverImage,
    }

    const savedPost = isEdit.value
      ? await updatePost(route.params.id, payload)
      : await createPost(payload)

    successMessage.value = isEdit.value ? '文章已更新' : '文章已发布'
    router.push(`/posts/${savedPost.ID}`)
  } catch (error) {
    errorMessage.value = error.message || '保存文章失败'
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadPostForEdit()
})
</script>

<template>
  <section class="editor-layout">
    <div class="editor-heading">
      <p class="eyebrow">Compose</p>
      <h1>{{ pageTitle }}</h1>
    </div>

    <p v-if="loading" class="state-text">正在加载文章...</p>

    <form v-else class="editor-form" @submit.prevent="submitPost">
      <label class="field">
        <span>标题</span>
        <input v-model.trim="title" type="text" maxlength="100" placeholder="写一个清楚的标题" />
      </label>

      <label class="field">
        <span>封面图片</span>
        <input type="file" accept="image/png,image/jpeg" @change="handleFileChange" />
      </label>

      <div v-if="coverPreview || coverImage" class="cover-preview">
        <img :src="coverPreview || coverImage" alt="封面预览" />
      </div>

      <label class="field">
        <span>正文</span>
        <textarea v-model.trim="content" rows="14" placeholder="写下你的文章内容"></textarea>
      </label>

      <p v-if="errorMessage" class="form-message error-text">{{ errorMessage }}</p>
      <p v-if="successMessage" class="form-message success-text">{{ successMessage }}</p>

      <button class="submit-button" type="submit" :disabled="submitting">
        {{ submitting ? '保存中...' : isEdit ? '更新文章' : '发布文章' }}
      </button>
    </form>
  </section>
</template>
