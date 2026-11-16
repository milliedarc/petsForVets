<script setup lang="ts">
import {ref} from "vue";
import {login as apiLogin} from "@/api/auth";
import {ApiError} from "@/api/client";
import AuthCard from "@/components/AuthCard.vue";

const email = ref<string>("");
const password = ref<string>("");
const loginError = ref<string>("");
const isSubmitting = ref<boolean>(false);

async function login(): Promise<void> {
  loginError.value = '';
  isSubmitting.value = true;
  try {
    await apiLogin(email.value, password.value);
    // A full page load lets App.vue pick up the new session.
    window.location.href = '/'
  } catch (error) {
    console.error(error);
    // A 401 means wrong credentials; anything else (server down, ...) gets its own message.
    loginError.value = error instanceof ApiError && error.status !== 401
        ? error.message
        : 'Invalid email or password';
    isSubmitting.value = false;
  }
}

</script>

<template>
  <AuthCard title="Welcome back" subtitle="Sign in to manage your pets">
    <form @submit.prevent="login">
      <div class="auth-field">
        <label for="email">Email</label>
        <InputText id="email" type="email" v-model="email" autocomplete="email" autofocus required fluid/>
      </div>
      <div class="auth-field">
        <label for="password">Password</label>
        <Password inputId="password"
                  v-model="password"
                  :feedback="false"
                  :inputProps="{autocomplete: 'current-password', required: true}"
                  toggleMask fluid/>
      </div>

      <Message v-if="loginError" severity="error" icon="pi pi-times-circle" role="alert">
        {{ loginError }}
      </Message>

      <Button type="submit" label="Sign in" icon="pi pi-sign-in" :loading="isSubmitting" fluid/>

      <p class="auth-footer">
        Don't have an account?
        <RouterLink :to="{name: 'UserRegistration'}">Register with us</RouterLink>
      </p>
    </form>
  </AuthCard>
</template>
