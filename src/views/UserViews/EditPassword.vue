<script setup lang="ts">
import {computed, ref} from "vue";
import {useRouter} from "vue-router";
import {useToast} from "primevue/usetoast";
import {pb} from "@/components/Pocketbase"

// ****************** CONSTS *********************

const router = useRouter();
const toast = useToast()

// ****************** REFS *********************

const currentPassword = ref<string>("");
const newPassword = ref<string>("");
const confirmPassword = ref<string>("");

// ****************** COMPUTED *********************


const isValidPassword = computed<boolean>(() => {
  return currentPassword.value !== '' && newPassword.value !== '' && confirmPassword.value !== '';
})

// ****************** FUNCTIONS *********************

function goToSettings(): void {
  router.push({name: 'Settings'})
}

async function savePassword(): Promise<void> {

  try {
    const updatedPassword = {
      oldPassword: currentPassword.value,
      password: newPassword.value,
      passwordConfirm: confirmPassword.value
    }
    await pb.collection('users').update(pb.authStore.record.id, updatedPassword)
    await router.push({name: 'UserLogin'})
    toast.add({severity: 'success', summary: 'Success!', detail: 'Password changed successfully', life: 3000});

  } catch (error) {
    let errorMessage = "Password change failed! "
    const responsePasswordMessage = error?.response?.data?.password?.message;
    const responseOldPasswordMessage = error?.response?.data?.oldPassword?.message;

    if (responsePasswordMessage) {
      errorMessage = errorMessage + responsePasswordMessage;

    } else if (responseOldPasswordMessage) {
      errorMessage = errorMessage + responseOldPasswordMessage;
    }

    toast.add({severity: 'error', summary: 'Error', detail: errorMessage, life: 10_000});
    console.log(error.response.data.password.message);
  }
}

</script>

<template>
  <main class="container">
    <div class="bg-body-tertiary p-5 rounded">

      <div class="d-flex justify-content-center mt-4 w-100">

        <Card style="width: 40rem" class="pt-4">
          <template #header>
            <div class="text-center">
              <h3>Change password</h3>
            </div>
          </template>
          <template #content>
            <div class="text-center w-100">
              <p>Keep your account secure by changing your password regularly.</p>
            </div>
            <div>
              <div class="inputDiv">
                <div class="d-flex flex-column gap-2 w-100">
                  <label for="currentPassword" class="fw-bold">Current password</label>
                  <Password fluid
                            id="currentPassword"
                            v-model="currentPassword"
                            :feedback="false"
                            toggleMask/>
                </div>
              </div>

              <div class="inputDiv">
                <div class="d-flex flex-column gap-2 w-100">
                  <label for="newPassword" class="fw-bold">New password</label>
                  <Password fluid
                            id="newPassword"
                            v-model="newPassword"
                            toggleMask/>
                </div>
              </div>

              <div class="inputDiv">
                <div class="d-flex flex-column gap-2 w-100">
                  <label for="confirmNewPassword" class="fw-bold">Confirm new password</label>
                  <Password fluid
                            id="confirmNewPassword"
                            v-model="confirmPassword"
                            toggleMask
                            :feedback="false"
                  />
                </div>
              </div>
              <div class="inputDiv">
                <p>For security reasons, after changing your password successfully, you will be logged out.</p>
              </div>
            </div>

          </template>
          <template #footer>
            <div class="pt-3 px-3">
              <div class="mb-3 text-center">
                <Button @click="savePassword"
                        label="Save" aria-label="Save password"
                        class="w-100"
                        :disabled="!isValidPassword"/>
              </div>
              <div class="text-center">
                <Button
                    @click="goToSettings"
                    label="Cancel" aria-label="Cancel password change"
                    severity="secondary" outlined class="w-100"/>
              </div>
            </div>
          </template>
        </Card>

      </div>

    </div>
  </main>
</template>

<style scoped>

.inputDiv {
  margin-top: 40px;
  margin-left: 20px;
  margin-right: 20px;
  display: flex;
}

</style>