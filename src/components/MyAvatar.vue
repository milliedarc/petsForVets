<script setup lang="ts">
import {ref, computed} from "vue";
import {useRouter} from "vue-router";
import {pb} from "@/components/Pocketbase";

const props = defineProps<{
  user: User
}>()

const emit = defineEmits<{
  (e: 'logout'): void
  (e: 'appModeSwitched', newMode: string): void
}>()

const router = useRouter();

const menu = ref();

const switchLabel = computed(() => {
  return `Switch to ${targetAppMode.value === 'clinicTeam' ? 'Clinic Team' : 'Pet Owner'} view`;
});

const menuItems = computed(() => {
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
      {
        label: switchLabel.value,
        icon: 'pi pi-arrow-right-arrow-left',
        shortcut: '⌘+1',
        command: () => {
          switchAppMode()
        }
      },
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

const toggle = (event) => {
  menu.value.toggle(event);
};

function goToSettings() {
  router.push({name: 'Settings'});
}

function goToUserDetails() {
  router.push({name: 'UserDetails'});
}

async function switchAppMode() {
  const newMode = targetAppMode.value;
  try {
    await pb.collection('users').update(pb.authStore.record.id, {
      app_mode: newMode
    });

    await pb.collection('users').authRefresh();

    window.location.href = '/'

  } catch (error) {
    console.error("Failed to switch app mode:", error);
  }
}

</script>

<template>
  <div class="d-flex justify-content-center align-items-center">
    <div class="me-3">
      Hi, {{ props.user?.name }}
    </div>
    <Avatar @click="toggle" icon="pi pi-user" class="mr-2" size="large"
            style="background-color: #d6b6c2; color: #2a1261; cursor: pointer"
            shape="circle"/>
    <Menu ref="menu" id="overlay_menu" :model="menuItems" :popup="true"/>
  </div>
</template>

<style scoped>

</style>