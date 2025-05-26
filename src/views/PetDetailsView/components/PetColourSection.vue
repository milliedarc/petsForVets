<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const colour = defineModel('colour', {
  type: String,
  default: '',
  required: true
})

</script>

<template>
  <section>
    <div>
      <p>What colour is
        <PetNameFormatted :name="props.name"/>
        <span>?</span>
      </p>
    </div>

    <div>
      <FloatLabel class="myInput" variant="on">
        <InputText id="colour" v-model="colour" class="myInput"/>
        <label for="colour">Colour</label>
      </FloatLabel>
    </div>
  </section>
</template>

<style scoped>

</style>