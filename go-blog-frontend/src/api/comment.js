import { request } from './client'

export function getComments(postId, page = 1, pageSize = 10) {
  return request(`/api/v1/posts/${postId}/comments?page=${page}&page_size=${pageSize}`)
}

export function createComment(postId, content) {
  return request(`/api/v1/posts/${postId}/comments`, {
    method: 'POST',
    body: {
      content,
    },
  })
}

export function deleteComment(commentId) {
  return request(`/api/v1/comments/${commentId}`, {
    method: 'DELETE',
  })
}
