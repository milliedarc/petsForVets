<script setup lang="ts">
import {ref} from "vue";

const name = defineModel('name', {
  type: String,
  required: true
})

/** Only complain about an empty name once the user has left the field. */
const isTouched = ref<boolean>(false)

</script>

<template>
  <section class="question">
    <label for="petName" class="question-label">My pet's name is:</label>
    <InputText v-model="name"
               id="petName"
               class="field-input"
               autocomplete="off"
               placeholder="e.g. Biscuit"
               :invalid="isTouched && name.trim() === ''"
               aria-describedby="petNameError"
               @blur="isTouched = true"/>
    <small v-if="isTouched && name.trim() === ''" id="petNameError" class="field-error">
      Please enter your pet's name to continue
    </small>
  </section>
</template>
