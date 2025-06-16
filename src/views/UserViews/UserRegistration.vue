<script setup lang="ts">
import {computed, ref} from "vue";

const email = ref<string>("");
const password = ref<string>("");
const passwordConfirm = ref<string>("");

const isValidPassword = computed(() => {
  return password.value !== '' && passwordConfirm.value !== '';
})

</script>

<template>
  <div style="background-color: whitesmoke">
    <div class="d-flex justify-content-center align-items-center vh-100 w-100">
      <div class="m-auto p-4 border rounded shadow-sm" style="width: 400px; background-color: white">
        <div class="d-flex justify-content-center align-items-center mb-2">
          <font-awesome-icon :icon="['fadr', 'paw']" class="fa-4x" style="--fa-secondary-color: #d6b6c2"/>
        </div>
        <h3 class="text-center">Welcome</h3>
        <div class="m-4">
          <p>To register, please provide us with your email address and a password.</p>
        </div>
        <form>
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <InputText type="email" autofocus class="form-control" id="inputEmail" v-model="email"/>
                <label for="inputEmail">Email</label>
              </FloatLabel>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <Password fluid
                          id="password"
                          v-model="password"
                          toggleMask/>
                <label for="password">Password</label>
              </FloatLabel>
            </div>
          </div>
          <div class="row mb-4 d-flex align-items-center">
            <div class="col-sm-10 w-100">
              <FloatLabel variant="on">
                <Password fluid
                          id="passwordConfirm"
                          v-model="passwordConfirm"
                          toggleMask/>
                <label for="passwordConfirm">Confirm password</label>
              </FloatLabel>
            </div>
          </div>
          <div class="d-flex justify-content-center">
            <Button severity="contrast" aria-label="Next step">Next
              <i class="pi pi-angle-double-right"></i>
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>