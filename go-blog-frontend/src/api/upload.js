import { request } from './client'

export function uploadImage(file) {
  const formData = new FormData()
  formData.append('file', file)

  return request('/api/v1/upload', {
    method: 'POST',
    body: formData,
  })
}
