<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const isNeutered = defineModel('isNeutered', {
  type: Boolean,
  required: true
})

</script>

<template>
  <section>
    <div>
      <p>Is
        <PetNameFormatted :name="props.name"/>
        <span> neutered?</span>
      </p>
    </div>

    <div class="flex items-center gap-2">
      <Checkbox v-model="isNeutered"
                binary
                inputId="neutered" name="neutered" value="true"/>
      <label class="ms-3" for="neutered">
        Please tick this box if your pet is neutered
      </label>
    </div>
  </section>
</template>

<style scoped>

</style>