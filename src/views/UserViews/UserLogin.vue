<script setup lang="ts">
import {ref} from "vue";
import {useRouter} from "vue-router";
import {pb} from "@/components/Pocketbase"

const router = useRouter()

const email = ref<string>("");
const password = ref<string>("");
const hasLoginError = ref<string>(false);

async function login(): Promise<void> {
  try {
    await pb.collection('users').authWithPassword(email.value, password.value);
    hasLoginError.value = false;
    window.location.href = '/'
  } catch (error) {
    hasLoginError.value = true;
    console.log(error);
  }
}

</script>

<template>
  <div style="background-color: whitesmoke">
    <div class="d-flex justify-content-center align-items-center vh-100 w-100">
      <div class="m-auto p-4 border rounded shadow-sm" style="width: 400px; background-color: white">
        <div class="d-flex justify-content-center align-items-center mb-2">
          <font-awesome-icon :icon="['fadr', 'paw']" class="fa-4x" style="--fa-secondary-color: #d6b6c2"/>
        </div>
        <h3 class="text-center">Welcome back</h3>
        <form>
          <div class="row mb-3 d-flex align-items-center">
            <label for="inputEmail">Email</label>
            <div class="col-sm-10 w-100">
              <InputText type="email" autofocus class="form-control" id="inputEmail" v-model="email"/>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <label for="inputPassword">Password</label>
            <div class="col-sm-10 w-100">
              <Password fluid
                        id="password"
                        v-model="password"
                        :feedback="false"
                        toggleMask/>
            </div>
            <div>
              <Message v-if="hasLoginError"
                       severity="error"
                       icon="pi pi-times-circle"
                       class="mt-4"
                       aria-describedby="loginError">
                Invalid email or password
              </Message>
            </div>
          </div>
          <div class="text-center">
            <p>Don't have an account?
              <RouterLink to="/register">Register with us</RouterLink>
            </p>
          </div>
          <div class="d-flex justify-content-center">
            <Button
                @click.prevent="login"
                severity="contrast"
                aria-label="Login">
              <i class="pi pi-user"></i>Sign in
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>