import { request } from './client'

export function login(username, password) {
  return request('/login', {
    method: 'POST',
    body: {
      username,
      password,
    },
  })
}

export function register(username, password) {
  return request('/register', {
    method: 'POST',
    body: {
      username,
      password,
    },
  })
}

export function changePassword(oldPassword, newPassword) {
  return request('/api/v1/me/password', {
    method: 'PUT',
    body: {
      old_password: oldPassword,
      new_password: newPassword,
    },
  })
}
