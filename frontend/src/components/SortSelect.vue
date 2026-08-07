<template>
  <div class="sort-select-wrap" ref="wrap" @click="toggle">
    <div class="sort-select-trigger" :class="{ open: open }">
      <span class="sort-select-label">{{ modelValueLabel }}</span>
      <span class="sort-select-arrow" :class="{ open: open }">▾</span>
    </div>
    <div v-if="open" class="sort-select-menu">
      <div
        v-for="opt in options"
        :key="opt.value"
        class="sort-select-option"
        :class="{ active: opt.value === modelValue }"
        @click="pick(opt.value)"
      >
        <span class="opt-check">{{ opt.value === modelValue ? '✓' : '' }}</span>
        <span class="opt-label">{{ opt.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  modelValue: string
  options: { value: string; label: string }[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: string): void
  (e: 'change'): void
}>()

const open = ref(false)
const wrap = ref<HTMLElement | null>(null)

const modelValueLabel = computed(() => {
  return props.options.find(o => o.value === props.modelValue)?.label || props.modelValue
})

function toggle() {
  open.value = !open.value
}

function pick(v: string) {
  if (v !== props.modelValue) {
    emit('update:modelValue', v)
    emit('change')
  }
  open.value = false
}

function onDocClick(e: MouseEvent) {
  if (wrap.value && !wrap.value.contains(e.target as Node)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', onDocClick))
onUnmounted(() => document.removeEventListener('click', onDocClick))
</script>

<style scoped>
.sort-select-wrap {
  position: relative;
  display: inline-block;
  font-family: var(--font-body);
}

.sort-select-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 16px;
  min-width: 120px;
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: var(--radius-pill);
  cursor: pointer;
  font-size: 13px;
  color: var(--text);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  transition: border-color 0.2s, box-shadow 0.2s;
  user-select: none;
}

.sort-select-trigger:hover {
  border-color: var(--primary);
  box-shadow: var(--shadow);
}

.sort-select-trigger.open {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(164, 117, 208, 0.12);
}

.sort-select-label {
  flex: 1;
  font-weight: 500;
  white-space: nowrap;
}

.sort-select-arrow {
  font-size: 12px;
  color: var(--text-muted);
  transition: transform 0.2s;
}

.sort-select-arrow.open {
  transform: rotate(180deg);
}

.sort-select-menu {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  min-width: 100%;
  background: rgba(255, 255, 255, 0.76);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: 16px;
  box-shadow: var(--glass-highlight), var(--shadow);
  padding: 6px;
  z-index: 50;
  animation: sortPop 0.18s cubic-bezier(0.22, 1, 0.36, 1);
}

@keyframes sortPop {
  from { opacity: 0; transform: translateY(-4px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.sort-select-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 9px 10px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  color: var(--text);
  transition: background 0.15s, color 0.15s;
}

.sort-select-option:hover {
  background: var(--bg-soft);
}

.sort-select-option.active {
  color: var(--primary);
  font-weight: 600;
}

.opt-check {
  width: 16px;
  color: var(--primary);
  font-size: 13px;
}

.opt-label {
  white-space: nowrap;
}
</style>
