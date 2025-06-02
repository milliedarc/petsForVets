<script setup lang="ts">
import {computed, onBeforeMount, ref} from "vue";
import {useRoute, useRouter} from 'vue-router'
import {pb} from "@/components/Pocketbase"

import MyNavbar from "@/components/MyNavbar.vue";

const showRouter = ref(false);
const user = ref(pb.authStore.record);

const router = useRouter();
const route = useRoute();

const showNavbar = computed(() => {
  return !(route.name === 'UserLogin' || route.name === 'UserRegistration')
})

function logOut() {
  pb.authStore.clear()
  location.reload()
}

function handleAppModeSwitched(newMode: string) {
  user.value = {...user.value, app_mode: newMode}
  pb.collection("users").authRefresh()
}

// if successful -> nav to dashboard; if unsuccessful -> stay in login
onBeforeMount(async () => {
  try {
    await pb.collection("users").authRefresh()
    // await router.push("/")
  } catch (error) {
    await router.push("/login")
  } finally {
    showRouter.value = true
  }
})

</script>

<template>
  <div v-if="showNavbar" class="card">
    <MyNavbar :user="user" @logout="logOut" @app-mode-switched="handleAppModeSwitched"/>
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
