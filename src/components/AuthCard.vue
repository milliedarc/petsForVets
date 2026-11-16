<script setup lang="ts">
import {PawPrint} from "@lucide/vue";

defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<!-- Centred card with the logo, used by the logged-out pages (sign in, register). -->
<template>
  <main class="auth-page">
    <div class="auth-card">
      <div class="auth-logo">
        <PawPrint :size="40" aria-hidden="true"/>
      </div>
      <p class="auth-brand">Pets for Vets</p>
      <h1 class="auth-title">{{ title }}</h1>
      <p v-if="subtitle" class="auth-subtitle">{{ subtitle }}</p>
      <slot/>
    </div>
  </main>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px var(--pfv-gutter);
  background: radial-gradient(circle at top, var(--p-primary-100), var(--pfv-page-bg) 60%);
}

.auth-card {
  width: 100%;
  max-width: 420px;
  padding: 32px 28px;
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--pfv-radius);
  box-shadow: 0 12px 32px rgba(66, 31, 49, 0.08);
}

.auth-logo {
  width: 72px;
  height: 72px;
  margin: 0 auto 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--pfv-accent-soft);
  color: var(--p-primary-color);
}

.auth-brand {
  margin: 0;
  text-align: center;
  font-weight: 700;
  font-size: 0.8rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--p-primary-color);
}

.auth-title {
  text-align: center;
  margin: 4px 0 4px;
}

.auth-subtitle {
  text-align: center;
  color: var(--p-text-muted-color);
  margin-bottom: 24px;
}

/* Shared by the forms in the slot. */
.auth-card :deep(form) {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.auth-card :deep(.auth-field) {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.auth-card :deep(.auth-field label) {
  font-weight: 600;
  font-size: 0.9rem;
}

.auth-card :deep(.auth-footer) {
  text-align: center;
  margin: 8px 0 0;
  color: var(--p-text-muted-color);
}
</style>
