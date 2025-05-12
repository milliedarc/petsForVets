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
    await pb.collection('users').authWithPassword(email.value, password.value);
    hasLoginError.value = false;
    await router.push('/')
  } catch (error) {
    hasLoginError.value = true;
    console.log(error);
  }
}

</script>

<template>
  <div class="vh-100 vw-100" style="background-color: whitesmoke">
    <div class="d-flex justify-content-center align-items-center vh-100 w-100">
      <div class="m-auto p-4 border rounded shadow-sm" style="width: 400px; background-color: white">
        <div class="d-flex justify-content-center align-items-center mb-2">
          <font-awesome-icon :icon="['fadr', 'paw']" class="fa-4x" style="--fa-secondary-color: #d6b6c2"/>
        </div>
        <h3 class="text-center">Welcome back</h3>
        <form>
          <div class="row mb-3 d-flex align-items-center">
            <label for="inputEmail" class="col-auto col-form-label">Email</label>
            <div class="col-sm-10 w-100">
              <InputText type="email" autofocus class="form-control" id="inputEmail" v-model="email"/>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <label for="inputPassword" class="col-auto col-form-label">Password</label>
            <div class="col-sm-10 w-100">
              <InputText type="password" class="form-control" id="inputPassword" v-model="password"/>
            </div>
            <div>
              <Message v-if="hasLoginError" severity="error" icon="pi pi-times-circle" class="mt-4">Invalid email or
                password
              </Message>
            </div>
          </div>
          <div class="d-flex justify-content-center">
            <Button class="btn btn-primary" @click.prevent="login" severity="contrast"><i class="pi pi-user"></i>Sign in
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>