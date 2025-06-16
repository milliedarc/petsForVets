<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const gender = defineModel('gender', {
  type: String,
  default: '',
  required: true
})

</script>

<template>
  <section>
    <div>
      <p>What's
        <PetNameFormatted :name="props.name"/>
        <span>'s gender?</span>
      </p>
    </div>

    <div class="d-flex gap-4">
      <div class="d-flex align-items-center">
        <RadioButton v-model="gender"
                     inputId="male"
                     name="gender"
                     value="male"/>
        <label class="ms-3" for="male">Male</label>
      </div>
      <div class="d-flex align-items-center">
        <RadioButton v-model="gender"
                     inputId="female"
                     name="gender"
                     value="female"/>
        <label class="ms-3" for="female">Female</label>
      </div>
    </div>
  </section>
</template>

<style scoped>

</style>