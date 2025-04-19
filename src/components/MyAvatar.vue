<script setup>
import {ref} from "vue";

const props = defineProps(['user'])
const emit = defineEmits(['logout']);
const menu = ref();
const menuItems = ref([
  {
    label: 'Profile',
    items: [
      // {
      //   label: 'Refresh',
      //   icon: 'pi pi-refresh'
      // },
      // {
      //   label: 'Export',
      //   icon: 'pi pi-upload'
      // },
      {
        label: 'Settings',
        icon: 'pi pi-cog',
        shortcut: '⌘+O'
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
]);

const toggle = (event) => {
  menu.value.toggle(event);
};

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