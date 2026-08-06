<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-header">
        <h1 class="register-logo">注册账号</h1>
        <p class="register-subtitle">哔咔漫画</p>
      </div>

      <form @submit.prevent="handleRegister" class="register-form">
        <input v-model="form.email" type="text" class="input-field" placeholder="邮箱" required />
        <input v-model="form.password" type="password" class="input-field" placeholder="密码" required />
        <input v-model="form.name" type="text" class="input-field" placeholder="昵称" required />
        <input v-model="form.birthday" type="date" class="input-field" placeholder="生日" />
        <select v-model="form.gender" class="input-field">
          <option value="m">男</option>
          <option value="f">女</option>
          <option value="bot">其他</option>
        </select>

        <div class="question-group">
          <select v-model="form.question1" class="input-field">
            <option :value="1">你最喜欢的漫画角色？</option>
            <option :value="2">你的第一个宠物叫什么？</option>
            <option :value="3">你的家乡在哪里？</option>
          </select>
          <input v-model="form.answer1" type="text" class="input-field" placeholder="答案1" />
        </div>
        <div class="question-group">
          <select v-model="form.question2" class="input-field">
            <option :value="1">你最喜欢的漫画角色？</option>
            <option :value="2">你的第一个宠物叫什么？</option>
            <option :value="3">你的家乡在哪里？</option>
          </select>
          <input v-model="form.answer2" type="text" class="input-field" placeholder="答案2" />
        </div>
        <div class="question-group">
          <select v-model="form.question3" class="input-field">
            <option :value="1">你最喜欢的漫画角色？</option>
            <option :value="2">你的第一个宠物叫什么？</option>
            <option :value="3">你的家乡在哪里？</option>
          </select>
          <input v-model="form.answer3" type="text" class="input-field" placeholder="答案3" />
        </div>

        <p v-if="errorMsg" class="error-tip">{{ errorMsg }}</p>
        <p v-if="successMsg" class="success-tip">{{ successMsg }}</p>

        <button type="submit" class="btn btn-primary register-btn" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>

        <router-link to="/login" class="login-link">已有账号？去登录</router-link>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '@/api'

const router = useRouter()
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const form = reactive({
  email: '',
  password: '',
  name: '',
  birthday: '2000-01-01',
  gender: 'm',
  question1: 1,
  answer1: '',
  question2: 2,
  answer2: '',
  question3: 3,
  answer3: '',
})

async function handleRegister() {
  if (!form.email || !form.password || !form.name) {
    errorMsg.value = '邮箱、密码和昵称不能为空'
    return
  }
  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''
  try {
    await register(form)
    successMsg.value = '注册成功！请前往邮箱激活账号'
    setTimeout(() => router.push('/login'), 2000)
  } catch (e: any) {
    errorMsg.value = e.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(ellipse 80% 60% at 20% -10%, rgba(255, 95, 168, 0.22), transparent 60%),
    radial-gradient(ellipse 60% 50% at 90% 110%, rgba(255, 143, 194, 0.16), transparent 60%),
    var(--bg);
  padding: 20px;
}
.register-container {
  width: 100%;
  max-width: 400px;
  background: rgba(44, 28, 41, 0.75);
  backdrop-filter: blur(24px) saturate(140%);
  -webkit-backdrop-filter: blur(24px) saturate(140%);
  border: 1px solid var(--border);
  border-radius: 24px;
  padding: 32px 24px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 95, 168, 0.12);
}
.register-header { text-align: center; margin-bottom: 24px; }
.register-logo {
  font-family: var(--font-display);
  font-size: 30px; font-weight: 800;
  background: linear-gradient(135deg, #ff9acb 0%, #ff5fa8 55%, #ff3d94 100%);
  -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent;
}
.register-subtitle { font-size: 13px; color: var(--text-muted); }
.register-form { display: flex; flex-direction: column; gap: 12px; }
.question-group { display: flex; flex-direction: column; gap: 8px; padding: 8px 0; border-top: 1px solid var(--border); }
.error-tip { color: var(--danger); font-size: 13px; text-align: center; }
.success-tip { color: var(--success); font-size: 13px; text-align: center; }
.register-btn { width: 100%; padding: 12px; font-size: 15px; }
.login-link { text-align: center; font-size: 13px; color: var(--primary); margin-top: 8px; }
</style>
