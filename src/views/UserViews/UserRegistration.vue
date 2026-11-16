<script setup lang="ts">
import {computed, ref} from "vue";
import AuthCard from "@/components/AuthCard.vue";
import {register} from "@/api/auth";
import {errorMessage as apiErrorMessage} from "@/api/client";

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
    errorMessage.value = apiErrorMessage(error, 'Registration failed, please try again.');
    console.error(error);
    isSubmitting.value = false;
  }
}

</script>

<template>
  <AuthCard title="Create your account" subtitle="Keep your pets' details in one place and request repeat prescriptions online.">
    <form @submit.prevent="registerUser">
      <div class="name-row">
        <div class="auth-field">
          <label for="inputName">First name</label>
          <InputText id="inputName" v-model="name" autocomplete="given-name" autofocus required fluid/>
        </div>
        <div class="auth-field">
          <label for="inputSurname">Surname</label>
          <InputText id="inputSurname" v-model="surname" autocomplete="family-name" fluid/>
        </div>
      </div>
      <div class="auth-field">
        <label for="inputEmail">Email</label>
        <InputText type="email" id="inputEmail" v-model="email" autocomplete="email" required fluid/>
      </div>
      <div class="auth-field">
        <label for="password">Password</label>
        <Password inputId="password"
                  v-model="password"
                  :inputProps="{autocomplete: 'new-password', 'aria-describedby': 'passwordHint'}"
                  toggleMask fluid/>
        <small id="passwordHint" class="field-hint m-0">At least 8 characters.</small>
      </div>
      <div class="auth-field">
        <label for="passwordConfirm">Confirm password</label>
        <Password inputId="passwordConfirm"
                  v-model="passwordConfirm"
                  :inputProps="{autocomplete: 'new-password'}"
                  :feedback="false"
                  :invalid="passwordConfirm !== '' && password !== passwordConfirm"
                  toggleMask fluid/>
        <small v-if="passwordProblem" class="field-error m-0" role="alert">{{ passwordProblem }}</small>
      </div>

      <Message v-if="errorMessage" severity="error" role="alert">{{ errorMessage }}</Message>

      <Button type="submit" label="Create account" icon="pi pi-arrow-right" iconPos="right"
              :disabled="!canSubmit" :loading="isSubmitting" fluid/>

      <p class="auth-footer">
        Already have an account?
        <RouterLink :to="{name: 'UserLogin'}">Sign in</RouterLink>
      </p>
    </form>
  </AuthCard>
</template>

<style scoped>
.name-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

@media (max-width: 400px) {
  .name-row {
    grid-template-columns: 1fr;
  }
}
</style>
