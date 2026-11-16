<script setup lang="ts">
import {computed, ref} from "vue";
import {PawPrint} from "@lucide/vue";
import {register} from "@/api/auth";
import {ApiError} from "@/api/client";

const name = ref<string>("");
const surname = ref<string>("");
const email = ref<string>("");
const password = ref<string>("");
const passwordConfirm = ref<string>("");
const errorMessage = ref<string>("");
const isSubmitting = ref<boolean>(false);

// Same rules as the API (checkPasswordRules in api/internal/httpapi/auth.go),
// checked here too so people get feedback before submitting.
const passwordProblem = computed<string>(() => {
  if (password.value !== '' && password.value.length < 8) {
    return 'Your password must be at least 8 characters.';
  }
  if (passwordConfirm.value !== '' && password.value !== passwordConfirm.value) {
    return "The passwords don't match.";
  }
  return '';
})

const canSubmit = computed<boolean>(() => {
  return name.value.trim() !== ''
      && email.value.includes('@')
      && password.value.length >= 8
      && password.value === passwordConfirm.value
      && !isSubmitting.value;
})

async function registerUser(): Promise<void> {
  errorMessage.value = '';
  isSubmitting.value = true;
  try {
    await register(email.value, password.value, name.value, surname.value);
    // Registering also logs you in; a full page load lets App.vue pick up the session.
    window.location.href = '/';
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : 'Registration failed, please try again.';
    console.error(error);
  } finally {
    isSubmitting.value = false;
  }
}

</script>

<template>
  <div style="background-color: whitesmoke">
    <div class="d-flex justify-content-center align-items-center vh-100 w-100">
      <div class="m-auto p-4 border rounded shadow-sm" style="width: 400px; background-color: white">
        <div class="d-flex justify-content-center align-items-center mb-2">
          <PawPrint :size="64" color="#d6b6c2"/>
        </div>
        <h3 class="text-center">Welcome</h3>
        <div class="m-4">
          <p>To register, please tell us your name, email address and a password.</p>
        </div>
        <form @submit.prevent="registerUser">
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <InputText autofocus class="form-control" id="inputName" v-model="name" autocomplete="given-name"/>
                <label for="inputName">First name</label>
              </FloatLabel>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <InputText class="form-control" id="inputSurname" v-model="surname" autocomplete="family-name"/>
                <label for="inputSurname">Surname</label>
              </FloatLabel>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <InputText type="email" class="form-control" id="inputEmail" v-model="email" autocomplete="email"/>
                <label for="inputEmail">Email</label>
              </FloatLabel>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <Password fluid
                          inputId="password"
                          v-model="password"
                          :inputProps="{autocomplete: 'new-password'}"
                          toggleMask/>
                <label for="password">Password</label>
              </FloatLabel>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <Password fluid
                          inputId="passwordConfirm"
                          v-model="passwordConfirm"
                          :inputProps="{autocomplete: 'new-password'}"
                          :feedback="false"
                          toggleMask/>
                <label for="passwordConfirm">Confirm password</label>
              </FloatLabel>
              <small v-if="passwordProblem" class="text-danger">{{ passwordProblem }}</small>
            </div>
          </div>

          <Message v-if="errorMessage" severity="error" size="small" class="mb-3">{{ errorMessage }}</Message>

          <div class="d-flex justify-content-center">
            <Button type="submit" severity="contrast" :disabled="!canSubmit" aria-label="Create account">
              Create account
              <i class="pi pi-angle-double-right"></i>
            </Button>
          </div>
          <div class="text-center mt-3">
            Already have an account?
            <RouterLink :to="{name: 'UserLogin'}">Sign in</RouterLink>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
