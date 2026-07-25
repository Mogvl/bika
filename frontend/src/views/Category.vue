<template>
  <div class="page-container">
    <div class="section-title">分类浏览</div>

    <div class="category-grid">
      <div
        v-for="cat in categories"
        :key="cat.title"
        class="category-card"
        @click="goCategory(cat.title)"
      >
        <div class="category-name">{{ cat.title }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories } from '@/api'
import type { Category } from '@/types'

const router = useRouter()
const categories = ref<Category[]>([])

onMounted(async () => {
  try {
    const res = await getCategories()
    categories.value = res.data?.categories?.filter((c: Category) => c.title) || []
  } catch {}
})

function goCategory(name: string) {
  router.push('/home')
  // 延迟一下等路由切换
  setTimeout(() => {
    window.dispatchEvent(new CustomEvent('select-category', { detail: name }))
  }, 100)
}
</script>

<style scoped>
.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
  padding: 4px;
}

.category-card {
  background: var(--bg-card);
  border-radius: var(--radius);
  padding: 24px 12px;
  text-align: center;
  box-shadow: var(--shadow);
  cursor: pointer;
  transition: transform 0.2s;
}

.category-card:hover { transform: translateY(-2px); }

.category-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text);
}
</style>
