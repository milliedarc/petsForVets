<script setup>
import {ref} from "vue";
import PocketBase from 'pocketbase'
import {useRouter} from "vue-router";

const pb = new PocketBase('http://127.0.0.1:8090');

const emit = defineEmits(["login"]);

const router = useRouter()

const email = ref("");
const password = ref("");
const hasLoginError = ref(false);

async function login() {
  try {
    const userData = await pb.collection('users').authWithPassword(email.value, password.value);
    hasLoginError.value = false;
    await router.push('/')
    console.log(userData);
  } catch (error) {
    hasLoginError.value = true;
    console.log(error);
  }
}

</script>

<template>
  <div class="d-flex justify-content-center align-items-center vh-100 w-100">
    <div class="m-auto p-4 border rounded shadow" style="width: 400px;">
      <h1 class="text-center">User Login</h1>
      <form>
        <div class="row mb-3 d-flex align-items-center">
          <label for="inputEmail3" class="col-auto col-form-label">Email</label>
          <div class="col-sm-10 w-100">
            <input type="email" autofocus class="form-control" id="inputEmail3" v-model="email"/>
          </div>
        </div>
        <div class="row mb-3 d-flex align-items-center">
          <label for="inputPassword3" class="col-auto col-form-label">Password</label>
          <div class="col-sm-10 w-100">
            <input type="password" class="form-control" id="inputPassword3" v-model="password"/>
          </div>
          <div class="alert alert-danger text-center mt-4" role="alert" v-if="hasLoginError">
            Invalid email or password
          </div>
        </div>
        <div class="d-flex justify-content-center">
          <button class="btn btn-primary px-4" @click.prevent="login">Sign in</button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>

</style>