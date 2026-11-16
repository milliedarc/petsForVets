<script setup lang="ts">
</script>

<template>
  <main class="page page--narrow">
    <div class="page-header">
      <h1>Settings</h1>
    </div>

    <Card>
      <template #content>
        <div class="setting">
          <div>
            <h2 class="setting-title">Password</h2>
            <p class="setting-text">Keep your account secure by changing your password regularly.</p>
          </div>
          <RouterLink :to="{name: 'EditPassword'}" custom v-slot="{navigate}">
            <Button label="Change password" icon="pi pi-lock" severity="secondary" variant="outlined"
                    @click="navigate"/>
          </RouterLink>
        </div>

        <div class="setting">
          <div>
            <h2 class="setting-title">Contact preferences</h2>
            <p class="setting-text">
              Choose how the staff at your veterinary surgery contact you, and about what.
            </p>
          </div>
          <Tag value="Coming soon" severity="secondary"/>
        </div>
      </template>
    </Card>
  </main>
</template>

<style scoped>
.setting {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px 24px;
  padding: 16px 0;
  border-bottom: 1px solid var(--p-content-border-color);
}

.setting:first-child {
  padding-top: 0;
}

.setting:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.setting > div {
  flex: 1 1 280px;
}

.setting-title {
  font-size: 1rem;
  margin-bottom: 2px;
}

.setting-text {
  margin: 0;
  font-size: 0.9rem;
  color: var(--p-text-muted-color);
}
</style>
