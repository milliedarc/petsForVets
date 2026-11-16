<script setup lang="ts">
import {PawPrint} from "@lucide/vue";
</script>

<template>
  <main class="page page--narrow">
    <div class="state-panel">
      <PawPrint :size="48" class="icon" aria-hidden="true"/>
      <h1>Page not found</h1>
      <p>We couldn't find the page you were looking for. It may have moved, or the link may be wrong.</p>
      <RouterLink to="/" custom v-slot="{navigate}">
        <Button label="Go to my pets" icon="pi pi-home" @click="navigate"/>
      </RouterLink>
    </div>
  </main>
</template>

<style scoped>
.icon {
  color: var(--p-primary-color);
}
</style>
