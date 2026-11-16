<script setup lang="ts">
import {useRouter} from 'vue-router'
import {useSession} from "@/stores/session";
import MyNavbar from "@/components/MyNavbar.vue";

// Who is signed in, and redirecting people who aren't, is handled by the
// router guard in src/router/router.ts before any page is shown.
const session = useSession();
const router = useRouter();

async function logOut(): Promise<void> {
  await session.logout();
  await router.push({name: 'UserLogin'});
}
</script>

<template>
  <header v-if="session.user.value" class="app-header">
    <MyNavbar :user="session.user.value" @logout="logOut"/>
  </header>

  <!-- Keyed on the path so going from one pet to another loads the new pet. -->
  <RouterView :key="$route.fullPath"/>

  <!--  PrimeVue components to be used in other components -->
  <Toast/>
  <ConfirmDialog/>
</template>

<style scoped>
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
}
</style>
