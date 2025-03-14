<script setup>
import {ref} from "vue";
import UserLogin from "@/components/UserLogin.vue";
import Dashboard from "@/components/Dashboard.vue";

const showDashboard = ref(true)
const showUserLogin = ref(false)

// change values later
function doLogin() {
  showUserLogin.value = false;
  showDashboard.value = true;
}

</script>

<template>

  <UserLogin v-if="showUserLogin" @login="doLogin"/>
  <Dashboard v-if="showDashboard"/>

</template>

<style scoped>

</style>
