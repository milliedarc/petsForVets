<script setup>
import {onMounted, ref} from "vue";
import UserLogin from "@/components/UserLogin.vue";
import Dashboard from "@/components/Dashboard.vue";

import PocketBase from 'pocketbase';

const pb = new PocketBase('http://127.0.0.1:8090');

const showDashboard = ref(false)
const showUserLogin = ref(true)

// change values later
function doLogin() {
  showUserLogin.value = false;
  showDashboard.value = true;
}

onMounted(async () => {
  try {
    const pepe = await pb.collection("users").authRefresh()
    doLogin()
    console.log('result:', pepe)
  } catch (error) {
    showUserLogin.value = true;
    showDashboard.value = false;
    console.log(error)
  }
})

</script>

<template>

  <UserLogin v-if="showUserLogin" @login="doLogin"/>
  <Dashboard v-if="showDashboard"/>

</template>

<style scoped>

</style>
