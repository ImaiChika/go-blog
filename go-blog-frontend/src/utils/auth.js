const TOKEN_KEY = 'go_blog_token'

export const AUTH_CHANGED_EVENT = 'go-blog-auth-changed'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
  window.dispatchEvent(new Event(AUTH_CHANGED_EVENT))
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
  window.dispatchEvent(new Event(AUTH_CHANGED_EVENT))
}

export function hasToken() {
  return Boolean(getToken())
}

export function getCurrentUsername() {
  const token = getToken()
  if (!token) {
    return ''
  }

  try {
    const payload = JSON.parse(window.atob(token.split('.')[1]))
    return payload.username || ''
  } catch {
    return ''
  }
}
