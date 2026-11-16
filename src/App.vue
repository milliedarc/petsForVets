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
  <div v-if="user" class="card">
    <MyNavbar :user="user" @logout="logOut"/>
  </div>

  <RouterView v-if="showRouter"/>

  <!--  PrimeVue components to be used in other components -->
  <Toast/>
  <ConfirmDialog/>

</template>

<style>
a {
  text-decoration: none !important;
}

a:hover {
  text-decoration: underline !important;
}

a:visited {
  text-decoration: none !important;
}

a, a:visited {
  color: rgb(132 63 97);
}
</style>
