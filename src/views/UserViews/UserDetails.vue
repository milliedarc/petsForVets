<script setup lang="ts">
import {useSession} from "@/stores/session";

// The router guard only shows this page to a signed-in user.
const {user} = useSession();

</script>

<template>
  <main class="page page--narrow">
    <div class="page-header">
      <h1>Personal details</h1>
    </div>

    <template v-if="user">
      <Card class="mb-4">
        <template #title>Contact details</template>
        <template #content>
          <dl class="details">
            <dt>Name</dt>
            <dd>{{ [user.name, user.surname].filter(Boolean).join(' ') }}</dd>
            <dt>Email</dt>
            <dd>{{ user.email }}</dd>
          </dl>
          <p class="text-muted text-sm m-0 mt-4">
            Changing your name or email isn't available online yet. Please contact the clinic if they need updating.
          </p>
        </template>
      </Card>

      <Card>
        <template #title>
          <div class="flex items-center justify-between gap-2">
            My addresses
            <Tag value="Coming soon" severity="secondary"/>
          </div>
        </template>
        <template #content>
          <p class="text-muted m-0">You'll soon be able to add your home address here.</p>
        </template>
      </Card>
    </template>
  </main>
</template>

<style scoped>
.details {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 8px 24px;
  margin: 0;
}

.details dt {
  font-weight: 600;
  color: var(--p-text-muted-color);
}

.details dd {
  margin: 0;
  overflow-wrap: anywhere;
}
</style>
