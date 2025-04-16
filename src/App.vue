<script setup lang="ts">
import {onBeforeMount, onMounted, ref} from "vue";
import {useRoute, useRouter} from 'vue-router'

import PocketBase from 'pocketbase';
import MyNavbar from "@/components/MyNavbar.vue";

const pb = new PocketBase('http://127.0.0.1:8090');

const showRouter = ref(false);
const router = useRouter()
const route = useRoute();

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

function logOut() {
  pb.authStore.clear()
  location.reload()
}

</script>

<template>
  <div v-if="route.name !== 'UserLogin'" class="card">
    <MyNavbar :user="pb.authStore.record" @logout="logOut"/>
  </div>

  <RouterView v-if="showRouter"/>

</template>

<style scoped>

</style>
