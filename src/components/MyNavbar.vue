<script setup>
import {ref} from "vue";
import MyAvatar from "@/components/MyAvatar.vue";
import {useRouter} from "vue-router";

const router = useRouter()
const emit = defineEmits(['logout']);
const props = defineProps(['user'])

const items = ref([
  {
    label: 'Home',
    icon: 'pi pi-home',
    command: () => {
      router.push('/')
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


</script>

<template>
  <Menubar :model="items">
    <template #start>
      <h5 class="mb-0">Pets for Vets</h5>
    </template>
    <template #end>
      <MyAvatar @logout="emit('logout')" :user="props.user"/>
    </template>
  </Menubar>
</template>

<style scoped>

</style>