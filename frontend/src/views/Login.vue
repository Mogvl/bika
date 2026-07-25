<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <h1 class="login-logo">PicACG</h1>
        <p class="login-subtitle">哔咔漫画 Web 版</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <input
            v-model="email"
            type="email"
            class="input-field"
            placeholder="邮箱"
            required
            autocomplete="email"
          />
        </div>
        <div class="form-group">
          <input
            v-model="password"
            type="password"
            class="input-field"
            placeholder="密码"
            required
            autocomplete="current-password"
          />
        </div>

        <p v-if="errorMsg" class="error-tip">{{ errorMsg }}</p>

        <button type="submit" class="btn btn-primary login-btn" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <div class="login-footer">
        <p class="disclaimer">
          使用哔咔漫画账号登录<br />
          本项目仅用于技术研究
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  if (!email.value || !password.value) {
    errorMsg.value = '请输入邮箱和密码'
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    await auth.login(email.value, password.value)
    router.push('/home')
  } catch (e: any) {
    errorMsg.value = e.message || '登录失败，请检查账号密码'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fce4ec 0%, #f8bbd0 100%);
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 360px;
  background: white;
  border-radius: 20px;
  padding: 40px 28px;
  box-shadow: 0 8px 32px rgba(231, 76, 139, 0.15);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  font-size: 36px;
  color: var(--primary);
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  width: 100%;
}

.error-tip {
  color: #e74c3c;
  font-size: 13px;
  text-align: center;
}

.login-btn {
  width: 100%;
  padding: 14px;
  font-size: 16px;
  margin-top: 8px;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
}

.disclaimer {
  font-size: 12px;
  color: #999;
  line-height: 1.6;
}
</style>
