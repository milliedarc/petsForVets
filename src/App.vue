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
  <nav v-if="route.name !== 'UserLogin'" class="navbar navbar-expand-md navbar-dark bg-dark mb-4">
    <div class="container-fluid">
      <a class="navbar-brand" href="#">Pets for Vets</a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarCollapse"
              aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarCollapse">
        <ul class="navbar-nav me-auto mb-2 mb-md-0">
          <li class="nav-item">
            <a class="nav-link active" aria-current="page" href="#">Home</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Link</a>
          </li>
          <li class="nav-item">
            <a class="nav-link disabled" aria-disabled="true">Disabled</a>
          </li>
        </ul>
        <div class="me-4 text-white">
          {{ pb.authStore.record?.name }}
        </div>
        <!--        <form class="d-flex" role="search">-->
        <!--          <input class="form-control me-2" type="search" placeholder="Search" aria-label="Search">-->
        <!--          <button class="btn btn-outline-success" type="submit">Search</button>-->
        <!--        </form>-->
        <button class="btn btn-outline-light" @click="logOut">Log out</button>
      </div>
    </div>
  </nav>


  <!--  <UserLogin v-if="showUserLogin" @login="doLogin"/>-->
  <!--  <Dashboard v-if="showDashboard" @viewPet="viewPet"/>-->
  <!--  <ViewPet v-if="showViewPet"/>-->
  <RouterView v-if="showRouter"/>

</template>

<style scoped>

</style>
