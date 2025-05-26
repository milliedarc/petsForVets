<script setup lang="ts">
import {useRouter} from "vue-router";
import ButtonPencilEdit from "@/components/ButtonPencilEdit.vue";

const router = useRouter();

function goToEditPassword() {
  router.push({name: 'EditPassword'});
}

</script>

<template>
  <main class="container">
    <div class="bg-body-tertiary p-5 rounded">
      <h1 class="mb-4">Settings</h1>

      <div class="card mb-3">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-center">

            <div class="mt-2">
              <div>
                <h6 class="fw-bold">Manage password</h6>
              </div>
              <div class="mt-3">
                <p>Keep your account secure by changing your password regularly.</p>
              </div>
            </div>

            <div class="d-flex justify-content-center">
              <ButtonPencilEdit @click="goToEditPassword"/>
            </div>
          </div>
        </div>
      </div>

      <div class="card mb-3">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-center">

            <div class="mt-2">
              <div>
                <h6 class="fw-bold">Contact preferences</h6>
              </div>
              <div class="mt-3">
                <p>Choose how you would like the staff at your veterinary surgery to contact you, and regarding
                  what.</p>
              </div>
            </div>

            <div class="d-flex justify-content-center">
              <ButtonPencilEdit/>
            </div>
          </div>
        </div>
      </div>

    </div>
  </main>
</template>

<style scoped>

</style>