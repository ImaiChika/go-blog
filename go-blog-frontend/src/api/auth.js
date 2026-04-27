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
