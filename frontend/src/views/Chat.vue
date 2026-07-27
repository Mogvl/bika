<template>
  <div class="chat-page">
    <!-- 房间列表 -->
    <div v-if="!currentRoom" class="room-list">
      <div class="section-title">💬 聊天室</div>

      <div v-if="loading" class="loading">加载中</div>

      <div v-else-if="!chatToken" class="chat-login">
        <p>需要登录聊天服务器</p>
        <button class="btn btn-primary" @click="loginChat">自动登录</button>
        <p v-if="loginError" class="error-tip">{{ loginError }}</p>
      </div>

      <div v-else class="rooms">
        <div
          v-for="room in rooms"
          :key="room.id"
          class="room-item"
          @click="joinRoom(room)"
        >
          <div class="room-icon">
            <img v-if="room.icon" :src="room.icon" alt="" />
            <span v-else class="room-icon-placeholder">💬</span>
          </div>
          <div class="room-info">
            <div class="room-title">{{ room.title }}</div>
            <div class="room-desc">{{ room.description }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 聊天界面 -->
    <div v-else class="chat-room">
      <div class="chat-header">
        <button class="btn-back" @click="leaveRoom">←</button>
        <span class="chat-title">{{ currentRoom.title }}</span>
        <button class="btn-refresh" @click="loadMessages">刷新</button>
      </div>

      <div class="messages-container" ref="messagesContainer">
        <div v-if="loadingMessages" class="loading">加载中</div>
        <div
          v-for="msg in messages"
          :key="msg._id"
          class="message-item"
          :class="{ 'message-mine': msg.sender?._id === myUserId }"
        >
          <div class="message-avatar">
            <img v-if="msg.sender?.avatar" :src="msg.sender.avatar" alt="" />
            <span v-else>{{ msg.sender?.name?.[0] || '?' }}</span>
          </div>
          <div class="message-content">
            <div class="message-sender">{{ msg.sender?.name || '匿名' }}</div>
            <div class="message-text">{{ msg.message }}</div>
          </div>
        </div>
      </div>

      <div class="chat-input">
        <input
          v-model="newMessage"
          type="text"
          placeholder="输入消息..."
          @keydown.enter="sendMessage"
        />
        <button class="btn-send" @click="sendMessage" :disabled="!newMessage.trim()">发送</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { chatLogin, getChatRooms, getChatMessages, sendChatMessage, getChatProfile } from '@/api'
import { useAuthStore } from '@/stores/auth'

interface Room {
  id: string
  title: string
  description: string
  icon: string
  isAvailable: boolean
  minLevel: number
}

interface ChatMessage {
  _id: string
  roomId: string
  message: string
  sender: {
    _id: string
    name: string
    avatar: string
    level: number
  }
  createdAt: string
}

const auth = useAuthStore()
const rooms = ref<Room[]>([])
const currentRoom = ref<Room | null>(null)
const messages = ref<ChatMessage[]>([])
const newMessage = ref('')
const loading = ref(true)
const loadingMessages = ref(false)
const chatToken = ref('')
const myUserId = ref('')
const loginError = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

onMounted(() => {
  // 检查本地存储的聊天token
  const savedToken = localStorage.getItem('chat_token')
  if (savedToken) {
    chatToken.value = savedToken
    loadRooms()
  } else {
    loading.value = false
  }
})

async function loginChat() {
  loginError.value = ''
  try {
    const res = await chatLogin(auth.user?.email || '', '')
    chatToken.value = res.data?.token || ''
    if (chatToken.value) {
      localStorage.setItem('chat_token', chatToken.value)
      await loadRooms()
      await loadProfile()
    }
  } catch (e: any) {
    loginError.value = e.message || '聊天登录失败'
  }
}

async function loadRooms() {
  loading.value = true
  try {
    const res = await getChatRooms()
    rooms.value = res.data?.rooms || []
  } catch {} finally { loading.value = false }
}

async function loadProfile() {
  try {
    const res = await getChatProfile()
    myUserId.value = res.data?._id || ''
  } catch {}
}

async function joinRoom(room: Room) {
  currentRoom.value = room
  messages.value = []
  await loadMessages()
}

function leaveRoom() {
  currentRoom.value = null
  messages.value = []
}

async function loadMessages() {
  if (!currentRoom.value) return
  loadingMessages.value = true
  try {
    const res = await getChatMessages(currentRoom.value.id)
    messages.value = (res.data?.messages || []).reverse()
    await nextTick()
    scrollToBottom()
  } catch {} finally { loadingMessages.value = false }
}

async function sendMessage() {
  if (!newMessage.value.trim() || !currentRoom.value) return
  try {
    await sendChatMessage(currentRoom.value.id, newMessage.value)
    newMessage.value = ''
    await loadMessages()
  } catch (e: any) {
    alert(e.message || '发送失败')
  }
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
</script>

<style scoped>
.chat-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 56px);
}

.room-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.chat-login {
  text-align: center;
  padding: 40px 20px;
}

.chat-login .btn {
  margin-top: 16px;
}

.error-tip {
  color: #e74c3c;
  font-size: 13px;
  margin-top: 8px;
}

.rooms {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.room-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: var(--bg-card);
  border-radius: var(--radius);
  cursor: pointer;
  box-shadow: var(--shadow);
}

.room-item:hover {
  background: var(--bg);
}

.room-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.room-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.room-icon-placeholder {
  font-size: 20px;
  color: white;
}

.room-info {
  flex: 1;
  min-width: 0;
}

.room-title {
  font-size: 15px;
  font-weight: 500;
}

.room-desc {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 4px;
}

/* 聊天界面 */
.chat-room {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 56px);
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
}

.btn-back, .btn-refresh {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: var(--text);
  padding: 4px 8px;
}

.chat-title {
  flex: 1;
  font-size: 16px;
  font-weight: 500;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-item {
  display: flex;
  gap: 8px;
  max-width: 80%;
}

.message-mine {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: white;
}

.message-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.message-content {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 8px 12px;
  box-shadow: var(--shadow);
}

.message-mine .message-content {
  background: var(--primary);
  color: white;
}

.message-sender {
  font-size: 11px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.message-mine .message-sender {
  color: rgba(255,255,255,0.7);
}

.message-text {
  font-size: 14px;
  line-height: 1.4;
  word-break: break-word;
}

.chat-input {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  background: var(--bg-card);
  border-top: 1px solid var(--border);
}

.chat-input input {
  flex: 1;
  padding: 10px 16px;
  border: 1px solid var(--border);
  border-radius: 20px;
  font-size: 14px;
  outline: none;
}

.chat-input input:focus {
  border-color: var(--primary);
}

.btn-send {
  padding: 10px 20px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
}

.btn-send:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
