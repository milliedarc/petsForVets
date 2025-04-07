<script setup>
import {onBeforeMount, onMounted, ref} from "vue";
import {useRoute, useRouter} from 'vue-router'

import UserLogin from "@/components/UserLogin.vue";
import Dashboard from "@/components/Dashboard.vue";
import ViewPet from "@/components/ViewPet.vue";

import PocketBase from 'pocketbase';

const pb = new PocketBase('http://127.0.0.1:8090');

// const showDashboard = ref(false)
// const showUserLogin = ref(true)
// const showViewPet = ref(false)
//
// function doLogin() {
//   showUserLogin.value = false;
//   showDashboard.value = true;
// }

// function viewPet() {
//   showDashboard.value = false;
//   showViewPet.value = true;
// }
const showRouter = ref(false);
const router = useRouter()

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

  <!--  <UserLogin v-if="showUserLogin" @login="doLogin"/>-->
  <!--  <Dashboard v-if="showDashboard" @viewPet="viewPet"/>-->
  <!--  <ViewPet v-if="showViewPet"/>-->
  <RouterView v-if="showRouter"/>

</template>

<style scoped>

</style>
