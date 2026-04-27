import { request } from './client'

export function getPostList(page = 1, pageSize = 10) {
  return request(`/api/v1/posts?page=${page}&page_size=${pageSize}`)
}

export function getPostById(id) {
  return request(`/api/v1/posts/${id}`)
}

export function createPost(post) {
  return request('/api/v1/posts', {
    method: 'POST',
    body: post,
  })
}

export function updatePost(id, post) {
  return request(`/api/v1/posts/${id}`, {
    method: 'PUT',
    body: post,
  })
}

export function deletePost(id) {
  return request(`/api/v1/posts/${id}`, {
    method: 'DELETE',
  })
}
