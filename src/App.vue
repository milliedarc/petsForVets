<script setup lang="ts">
import type {User} from "@/types/User";
import {onBeforeMount, ref} from "vue";
import {useRouter} from 'vue-router'
import {getCurrentUser, logout} from "@/api/auth";

import MyNavbar from "@/components/MyNavbar.vue";

const showRouter = ref<boolean>(false);
const user = ref<User | undefined>(undefined);

const router = useRouter();

/** Pages a logged-out visitor may stay on. */
const publicRouteNames = ['UserLogin', 'UserRegistration'];

async function logOut(): Promise<void> {
  await logout()
  location.href = '/login'
}

onBeforeMount(async () => {
  try {
    const currentUser = await getCurrentUser()
    user.value = currentUser;
    if (currentUser.app_mode === 'clinicTeam') {
      await router.push("/clinic")
    } else {
      // Already signed in, so the sign-in and register pages have nothing to offer.
      await router.isReady()
      if (publicRouteNames.includes(router.currentRoute.value.name as string)) {
        await router.replace("/")
      }
    }
  } catch (error) {
    await router.isReady()
    if (!publicRouteNames.includes(router.currentRoute.value.name as string)) {
      await router.push("/login")
    }
  } finally {
    showRouter.value = true
  }
})

</script>

<template>
  <header v-if="user" class="app-header">
    <MyNavbar :user="user" @logout="logOut"/>
  </header>

  <RouterView v-if="showRouter"/>

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
