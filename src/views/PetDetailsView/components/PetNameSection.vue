<script setup lang="ts">

const name = defineModel('name', {
  type: String,
  required: true
})

</script>

<template>
  <section>
    <p>My pet's name is:</p>
    <FloatLabel variant="on">
      <InputText v-model="name" id="petName" class="myInput"/>
      <label for="petName">Pet name *</label>
    </FloatLabel>
    <Message v-if="name === ''" severity="error" size="small" variant="simple" aria-describedby="nameError">
      Please enter a pet's name to continue
    </Message>
  </section>
</template>

<style scoped>

</style>