<script setup lang="ts">
import type {User} from "@/types/User";
import {ref, computed} from "vue";
import {useRouter} from "vue-router";
import {updateAppMode} from "@/api/auth";
import {useToast} from "primevue/usetoast";
import {ApiError} from "@/api/client";

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

const menuItems = computed<[{}]>(() => {
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
        shortcut: '⌘+O',
        command: () => {
          goToSettings()
        }
      },
      // Only clinic staff can switch views (the API enforces this too).
      ...(props.user.veterinary ? [{
        label: switchLabel.value,
        icon: 'pi pi-arrow-right-arrow-left',
        shortcut: '⌘+1',
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
        shortcut: '⌘+Q',
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
    const detail = error instanceof ApiError ? error.message : 'Could not switch view';
    toast.add({severity: 'error', summary: 'Error', detail, life: 5000});
  }
}

</script>

<template>
  <div class="d-flex justify-content-center align-items-center">
    <div class="me-3">
      Hi, {{ props.user?.name }}
    </div>
    <Avatar @click="toggle"
            icon="pi pi-user"
            class="mr-2" size="large"
            style="background-color: #d6b6c2; color: #2a1261; cursor: pointer"
            shape="circle"/>
    <Menu ref="menu" id="overlay_menu" :model="menuItems" :popup="true"/>
  </div>
</template>

<style scoped>

</style>