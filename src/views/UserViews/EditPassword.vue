<script setup lang="ts">
import {computed, ref} from "vue";
import {useToast} from "primevue/usetoast";
import {errorMessage} from "@/api/client";
import {changePassword} from "@/api/auth";
import {useRouter} from "vue-router";
import {useSession} from "@/stores/session";

// ****************** CONSTS *********************

const toast = useToast()
const router = useRouter()
const session = useSession()

// ****************** REFS *********************

const currentPassword = ref<string>("");
const newPassword = ref<string>("");
const confirmPassword = ref<string>("");
const isSaving = ref<boolean>(false);

// ****************** COMPUTED *********************


const passwordsMatch = computed<boolean>(() => {
  return newPassword.value === confirmPassword.value;
})

const isValidPassword = computed<boolean>(() => {
  return currentPassword.value !== ''
      && newPassword.value.length >= 8
      && passwordsMatch.value;
})

// ****************** FUNCTIONS *********************

async function savePassword(): Promise<void> {
  if (!isValidPassword.value) {
    return;
  }
  isSaving.value = true;
  try {
    await changePassword(currentPassword.value, newPassword.value)
    // The server has ended every session, including this one.
    session.clear();
    toast.add({severity: 'success', summary: 'Password changed', detail: 'Please sign in with your new password.', life: 6000});
    await router.push({name: 'UserLogin'});
  } catch (error) {
    toast.add({severity: 'error', summary: 'Password change failed', detail: errorMessage(error, 'Please try again'), life: 10_000});
    console.error(error);
    isSaving.value = false;
  }
}

</script>

<template>
  <main class="page page--narrow">
    <RouterLink :to="{name: 'Settings'}" class="back-link">
      <i class="pi pi-arrow-left" aria-hidden="true"></i>Settings
    </RouterLink>
    <div class="page-header">
      <h1>Change password</h1>
    </div>

    <Card>
      <template #content>
        <form class="password-form" @submit.prevent="savePassword">
          <div class="field">
            <label for="currentPassword">Current password</label>
            <Password inputId="currentPassword"
                      v-model="currentPassword"
                      :feedback="false"
                      :inputProps="{autocomplete: 'current-password'}"
                      toggleMask fluid/>
          </div>

          <div class="field">
            <label for="newPassword">New password</label>
            <Password inputId="newPassword"
                      v-model="newPassword"
                      :inputProps="{autocomplete: 'new-password'}"
                      :invalid="newPassword !== '' && newPassword.length < 8"
                      toggleMask fluid/>
            <small v-if="newPassword !== '' && newPassword.length < 8" class="field-error m-0">
              The new password must be at least 8 characters.
            </small>
            <small v-else class="field-hint m-0">At least 8 characters.</small>
          </div>

          <div class="field">
            <label for="confirmNewPassword">Confirm new password</label>
            <Password inputId="confirmNewPassword"
                      v-model="confirmPassword"
                      :feedback="false"
                      :inputProps="{autocomplete: 'new-password'}"
                      :invalid="confirmPassword !== '' && !passwordsMatch"
                      toggleMask fluid/>
            <small v-if="confirmPassword !== '' && !passwordsMatch" class="field-error m-0">
              The passwords don't match.
            </small>
          </div>

          <Message severity="info" icon="pi pi-info-circle" size="small">
            For your security, you'll be signed out everywhere once your password has changed.
          </Message>

          <div class="flex flex-wrap gap-3">
            <Button type="submit" label="Change password" :disabled="!isValidPassword" :loading="isSaving"/>
            <RouterLink :to="{name: 'Settings'}" custom v-slot="{navigate}">
              <Button label="Cancel" severity="secondary" variant="outlined" @click="navigate"/>
            </RouterLink>
          </div>
        </form>
      </template>
    </Card>
  </main>
</template>

<style scoped>
.password-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 420px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field label {
  font-weight: 600;
}
</style>
