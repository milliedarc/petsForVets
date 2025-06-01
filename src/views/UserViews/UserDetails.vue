<script setup lang="ts">
import ButtonPencilEdit from "@/components/ButtonPencilEdit.vue";
import {useRouter} from "vue-router";

const router = useRouter();

function goToContactDetails() {
  router.push({name: "EditContactDetails"});
}

</script>

<template>
  <main class="container">
    <div class="bg-body-tertiary p-5 rounded">
      <h1 class="mb-4">Your details</h1>

      <div class="card mb-3">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-center">
            <div class="mt-2">
              <div>
                <h6 class="fw-bold">Contact details</h6>
              </div>
              <div class="mt-3">
                <p></p>
              </div>
            </div>
            <div class="d-flex justify-content-center">
              <ButtonPencilEdit @click="goToContactDetails"/>
            </div>
          </div>
        </div>
      </div>

      <div class="card mb-3">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-center">
            <div class="mt-2">
              <div>
                <h6 class="fw-bold">My addresses</h6>
              </div>
              <div class="mt-3">
                <p></p>
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