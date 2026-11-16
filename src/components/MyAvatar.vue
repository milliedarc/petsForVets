<script setup lang="ts">
import type {User} from "@/types/User";
import {ref, computed} from "vue";
import {useRouter} from "vue-router";
import {updateAppMode} from "@/api/auth";
import {useToast} from "primevue/usetoast";
import {errorMessage} from "@/api/client";
import type {MenuItem} from "primevue/menuitem";

// ****************** REFS *********************

const menu = ref();

// ****************** PROPS *********************
const props = defineProps<{
  user: User
}>()

// ****************** EMITS *********************

const emit = defineEmits<{
  (e: 'logout'): void
  (e: 'appModeSwitched', newMode: string): void
}>()

// ****************** CONST *********************

const router = useRouter();
const toast = useToast();

// ****************** COMPUTED *********************

const switchLabel = computed(() => {
  return `Switch to ${targetAppMode.value === 'clinicTeam' ? 'Clinic Team' : 'Pet Owner'} view`;
});

const menuItems = computed<MenuItem[]>(() => {
  return [{
    label: 'Profile',
    items: [
      {
        label: 'Personal details',
        icon: 'pi pi-user',
        command: () => {
          goToUserDetails()
        }
      },
      {
        label: 'Settings',
        icon: 'pi pi-cog',
        command: () => {
          goToSettings()
        }
      },
      // Only clinic staff can switch views (the API enforces this too).
      ...(props.user.veterinary ? [{
        label: switchLabel.value,
        icon: 'pi pi-arrow-right-arrow-left',
        command: () => {
          switchAppMode()
        }
      }] : []),
      {
        separator: true
      },
      {
        label: 'Logout',
        icon: 'pi pi-sign-out',
        command: () => {
          emit('logout')
        }
      }
    ]
  }
  ]
})

const targetAppMode = computed(() => {
  return props.user?.app_mode === 'petOwner' ? 'clinicTeam' : 'petOwner'
})

const initials = computed<string>(() => {
  const first = props.user.name.trim().charAt(0);
  const last = props.user.surname.trim().charAt(0);
  return (first + last).toUpperCase() || '?';
})

const toggle = (event: Event) => {
  menu.value.toggle(event);
};

// ****************** FUNCTIONS *********************

function goToSettings(): void {
  router.push({name: 'Settings'});
}

function goToUserDetails(): void {
  router.push({name: 'UserDetails'});
}

async function switchAppMode(): Promise<void> {
  const newMode = targetAppMode.value;
  try {
    await updateAppMode(newMode);
    window.location.href = newMode === 'clinicTeam' ? '/clinic' : '/'
  } catch (error) {
    console.error("Failed to switch app mode:", error);
    toast.add({severity: 'error', summary: 'Error', detail: errorMessage(error, 'Could not switch view'), life: 5000});
  }
}

</script>

<template>
  <button type="button"
          class="account-button"
          @click="toggle"
          aria-haspopup="true"
          aria-controls="account_menu">
    <span class="greeting">Hi, {{ props.user.name }}</span>
    <Avatar :label="initials" shape="circle" class="account-avatar"/>
    <i class="pi pi-angle-down text-muted" aria-hidden="true"></i>
    <span class="sr-only">Open account menu</span>
  </button>
  <Menu ref="menu" id="account_menu" :model="menuItems" :popup="true"/>
</template>

<style scoped>
.account-button {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 4px 8px 4px 12px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--p-text-color);
  font: inherit;
  cursor: pointer;
}

.account-button:hover {
  background: var(--p-content-hover-background);
}

.account-button:focus-visible {
  outline: 2px solid var(--p-primary-color);
  outline-offset: 2px;
}

.account-avatar {
  background: var(--pfv-accent-soft);
  color: var(--p-primary-700);
  font-weight: 700;
}

@media (max-width: 575px) {
  .greeting {
    display: none;
  }
}
</style>
