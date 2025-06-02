<script setup lang="ts">
import {ref} from "vue";
import MyAvatar from "@/components/MyAvatar.vue";
import {useRouter} from "vue-router";

const router = useRouter()

const emit = defineEmits<{
  (e: 'logout'): void
  (e: 'appModeSwitched', newMode: string): void
}>()

const props = defineProps<{
  user: User
}>()

const items = ref([
  {
    label: 'Home',
    icon: 'pi pi-home',
    command: () => {
      goToDashboard()
    }
  },
  {
    label: 'Features',
    icon: 'pi pi-star'
  },
  {
    label: 'Projects',
    icon: 'pi pi-search',
    items: [
      {
        label: 'Components',
        icon: 'pi pi-bolt'
      },
      {
        label: 'Blocks',
        icon: 'pi pi-server'
      },
      {
        label: 'UI Kit',
        icon: 'pi pi-pencil'
      },
      {
        label: 'Templates',
        icon: 'pi pi-palette',
        items: [
          {
            label: 'Apollo',
            icon: 'pi pi-palette'
          },
          {
            label: 'Ultima',
            icon: 'pi pi-palette'
          }
        ]
      }
    ]
  },
  {
    label: 'Contact',
    icon: 'pi pi-envelope'
  }
]);

function goToDashboard() {
  if (props.user.app_mode === 'petOwner') {
    router.push({name: 'DashboardPetOwner'})
  } else {
    router.push({name: 'DashboardClinicTeam'})
  }
}

</script>

<template>
  <Menubar :model="items">
    <template #start>
      <font-awesome-icon :icon="['fadr', 'paw']" class="fa-xl" style="--fa-secondary-color: #d6b6c2"/>
      <h5 class="mb-0 ms-2">Pets for Vets</h5>
    </template>
    <template #end>
      <MyAvatar @logout="emit('logout')"
                @app-mode-switched="emit('appModeSwitched', $event)"
                :user="props.user"/>
    </template>
  </Menubar>
</template>

<style scoped>

</style>