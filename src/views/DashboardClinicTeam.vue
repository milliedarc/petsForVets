<script setup lang="ts">
import {Stethoscope} from "@lucide/vue";
</script>

<template>
  <main class="page">
    <div class="page-header">
      <h1>Clinic dashboard</h1>
    </div>
    <div class="state-panel">
      <Stethoscope :size="48" class="icon" aria-hidden="true"/>
      <h2>Coming soon</h2>
      <p>
        This is where the clinic team will review repeat prescription requests and write prescriptions.
        Switch back to the pet owner view from your account menu in the meantime.
      </p>
    </div>
  </main>
</template>

<style scoped>
.icon {
  color: var(--p-primary-color);
}
</style>
