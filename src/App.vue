<script setup lang="ts">
import {computed, onBeforeMount, ref} from "vue";
import {useRouter} from 'vue-router'
import {pb} from "@/components/Pocketbase"

import MyNavbar from "@/components/MyNavbar.vue";

const showRouter = ref<boolean>(false);
const user = ref<User | undefined>(undefined);

const router = useRouter();

const showNavbar = computed<boolean>(() => {
  return user.value !== undefined
})

function logOut(): void {
  pb.authStore.clear()
  location.reload()
}

onBeforeMount(async () => {
  try {
    await pb.collection("users").authRefresh()
    user.value = pb.authStore.record as any;
    if (user.value.app_mode === 'clinicTeam') {
      await router.push("/clinic")
    }
  } catch (error) {
    await router.push("/login")
  } finally {
    showRouter.value = true
  }
})

</script>

<template>
  <div v-if="showNavbar" class="card">
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
