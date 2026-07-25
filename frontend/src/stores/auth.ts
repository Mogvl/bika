import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, getProfile as apiGetProfile } from '@/api'
import type { UserProfile } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<UserProfile | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  async function login(email: string, password: string) {
    loading.value = true
    try {
      const res = await apiLogin(email, password)
      const newToken = res.data.token
      if (!newToken) {
        throw new Error('登录失败，请检查账号密码')
      }
      token.value = newToken
      localStorage.setItem('token', newToken)
      await fetchProfile()
      return true
    } catch (e: any) {
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchProfile() {
    try {
      const res = await apiGetProfile()
      user.value = res.data.user || res.data
      localStorage.setItem('user', JSON.stringify(user.value))
    } catch {
      // ignore
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function init() {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser)
      } catch {}
    }
  }

  return { token, user, loading, isLoggedIn, login, fetchProfile, logout, init }
})
