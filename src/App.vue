<script setup lang="ts">
import type {User} from "@/types/User";
import {onBeforeMount, ref} from "vue";
import {useRouter} from 'vue-router'
import {pb} from "@/components/Pocketbase"

import MyNavbar from "@/components/MyNavbar.vue";

const showRouter = ref<boolean>(false);
const user = ref<User | undefined>(undefined);

const router = useRouter();

function logOut(): void {
  pb.authStore.clear()
  location.reload()
}

onBeforeMount(async () => {
  try {
    const {record} = await pb.collection("users").authRefresh<User>()
    user.value = record;
    if (record.app_mode === 'clinicTeam') {
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
