<script setup lang="ts">
import type {User} from "@/types/User";
import {computed} from "vue";
import MyAvatar from "@/components/MyAvatar.vue";
import {useRouter} from "vue-router";
import {PawPrint} from "@lucide/vue";

const router = useRouter()

const emit = defineEmits<{
  (e: 'logout'): void
}>()

const props = defineProps<{
  user: User
}>()

/** Clinic staff working in the clinic view have a different home page. */
const homePath = computed<string>(() => {
  return props.user.app_mode === 'clinicTeam' ? '/clinic' : '/'
})

const items = computed(() => [
  {
    label: props.user.app_mode === 'clinicTeam' ? 'Clinic' : 'My pets',
    icon: 'pi pi-home',
    command: () => {
      router.push(homePath.value)
    }
  }
]);

</script>

<template>
  <Menubar :model="items" class="navbar">
    <template #start>
      <RouterLink :to="homePath" class="brand" aria-label="Pets for Vets home">
        <PawPrint :size="26" aria-hidden="true"/>
        <span class="brand-name">Pets for Vets</span>
      </RouterLink>
    </template>
    <template #end>
      <MyAvatar
          @logout="emit('logout')"
          :user="props.user"/>
    </template>
  </Menubar>
</template>

<style scoped>
.navbar {
  border-radius: 0;
  border-width: 0 0 1px 0;
  padding: 8px var(--pfv-gutter);
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-right: 16px;
  color: var(--p-primary-color);
}

.brand:hover {
  text-decoration: none;
}

.brand-name {
  font-weight: 800;
  font-size: 1.125rem;
  color: var(--p-text-color);
}
</style>
