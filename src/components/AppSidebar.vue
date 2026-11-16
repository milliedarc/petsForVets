<script setup lang="ts">
import type {User} from "@/types/User";
import {computed} from "vue";
import {useRoute, useRouter, type RouteLocationRaw} from "vue-router";
import {useToast} from "primevue/usetoast";
import {PawPrint} from "@lucide/vue";
import {useSession} from "@/stores/session";
import {errorMessage} from "@/api/client";

// ****************** TYPES *********************

interface NavItem {
  label: string
  icon: string
  to: RouteLocationRaw
  /** Route names that highlight this item (a pet's pages count as "My pets"). */
  activeFor: string[]
}

// ****************** PROPS / EMITS *********************

const props = defineProps<{
  user: User
}>()

const emit = defineEmits<{
  (e: 'logout'): void
  /** A link was followed, so the mobile drawer can close. */
  (e: 'navigate'): void
}>()

// ****************** CONSTS *********************

const route = useRoute();
const router = useRouter();
const toast = useToast();
const session = useSession();

const accountItems: NavItem[] = [
  {label: 'Personal details', icon: 'pi pi-user', to: {name: 'UserDetails'}, activeFor: ['UserDetails']},
  {label: 'Settings', icon: 'pi pi-cog', to: {name: 'Settings'}, activeFor: ['Settings', 'EditPassword']},
];

// ****************** COMPUTED *********************

const isClinicView = computed<boolean>(() => props.user.app_mode === 'clinicTeam');

const mainItems = computed<NavItem[]>(() => {
  if (isClinicView.value) {
    return [
      {label: 'Clinic dashboard', icon: 'pi pi-briefcase', to: '/clinic', activeFor: ['DashboardClinicTeam']},
    ];
  }
  return [
    {
      label: 'My pets',
      icon: 'pi pi-home',
      to: '/',
      activeFor: ['DashboardPetOwner', 'PetDetails', 'Prescriptions'],
    },
    {label: 'Add a pet', icon: 'pi pi-plus-circle', to: {name: 'PetAdd'}, activeFor: ['PetAdd']},
  ];
});

const fullName = computed<string>(() => {
  return [props.user.name, props.user.surname].filter(Boolean).join(' ');
});

const initials = computed<string>(() => {
  const first = props.user.name.trim().charAt(0);
  const last = props.user.surname.trim().charAt(0);
  return (first + last).toUpperCase() || '?';
});

// ****************** FUNCTIONS *********************

function isActive(item: NavItem): boolean {
  return item.activeFor.includes(route.name as string);
}

async function switchAppMode(): Promise<void> {
  const newMode = isClinicView.value ? 'petOwner' : 'clinicTeam';
  try {
    await session.switchAppMode(newMode);
    emit('navigate');
    await router.push(newMode === 'clinicTeam' ? '/clinic' : '/');
  } catch (error) {
    console.error("Failed to switch app mode:", error);
    toast.add({severity: 'error', summary: 'Error', detail: errorMessage(error, 'Could not switch view'), life: 5000});
  }
}

</script>

<template>
  <nav class="sidebar" aria-label="Main">
    <RouterLink :to="isClinicView ? '/clinic' : '/'" class="brand" @click="emit('navigate')">
      <span class="brand-logo"><PawPrint :size="22" aria-hidden="true"/></span>
      <span class="brand-name">Pets for Vets</span>
    </RouterLink>

    <div class="nav-section">
      <ul class="nav-list">
        <li v-for="item in mainItems" :key="item.label">
          <RouterLink :to="item.to"
                      class="nav-link"
                      :class="{active: isActive(item)}"
                      :aria-current="isActive(item) ? 'page' : undefined"
                      @click="emit('navigate')">
            <i :class="item.icon" aria-hidden="true"></i>
            {{ item.label }}
          </RouterLink>
        </li>
      </ul>
    </div>

    <div class="nav-section">
      <h2 class="nav-heading">Account</h2>
      <ul class="nav-list">
        <li v-for="item in accountItems" :key="item.label">
          <RouterLink :to="item.to"
                      class="nav-link"
                      :class="{active: isActive(item)}"
                      :aria-current="isActive(item) ? 'page' : undefined"
                      @click="emit('navigate')">
            <i :class="item.icon" aria-hidden="true"></i>
            {{ item.label }}
          </RouterLink>
        </li>
        <!-- Only clinic staff can switch views (the API enforces this too). -->
        <li v-if="props.user.veterinary">
          <button type="button" class="nav-link" @click="switchAppMode">
            <i class="pi pi-arrow-right-arrow-left" aria-hidden="true"></i>
            {{ isClinicView ? 'Switch to pet owner view' : 'Switch to clinic view' }}
          </button>
        </li>
      </ul>
    </div>

    <div class="sidebar-footer">
      <div class="user">
        <Avatar :label="initials" shape="circle" class="user-avatar"/>
        <div class="user-text">
          <div class="user-name">{{ fullName }}</div>
          <div class="user-email">{{ props.user.email }}</div>
        </div>
      </div>
      <button type="button" class="nav-link" @click="emit('logout')">
        <i class="pi pi-sign-out" aria-hidden="true"></i>
        Log out
      </button>
    </div>
  </nav>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  gap: 24px;
  height: 100%;
  padding: 20px 14px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 8px;
  color: var(--p-text-color);
}

.brand-logo {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: var(--p-primary-color);
  color: var(--p-primary-contrast-color);
}

.brand-name {
  font-weight: 800;
  font-size: 1.125rem;
}

.nav-heading {
  margin: 0 0 6px;
  padding: 0 12px;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--p-text-muted-color);
}

.nav-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

/* Shared by links and buttons so they look the same. */
.nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 9px 12px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  color: var(--p-text-color);
  font: inherit;
  font-weight: 600;
  font-size: 0.95rem;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.15s, color 0.15s;
}

.nav-link i {
  font-size: 1rem;
  color: var(--p-text-muted-color);
  transition: color 0.15s;
}

.nav-link:hover {
  background: var(--p-content-hover-background);
}

.nav-link:focus-visible {
  outline: 2px solid var(--p-primary-color);
  outline-offset: -2px;
}

.nav-link.active {
  background: var(--pfv-accent-soft);
  color: var(--p-primary-700);
}

.nav-link.active i {
  color: var(--p-primary-color);
}

.sidebar-footer {
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid var(--p-content-border-color);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 8px;
}

.user-avatar {
  flex-shrink: 0;
  background: var(--pfv-accent-soft);
  color: var(--p-primary-700);
  font-weight: 700;
}

.user-text {
  min-width: 0;
}

.user-name {
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-email {
  font-size: 0.8rem;
  color: var(--p-text-muted-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
