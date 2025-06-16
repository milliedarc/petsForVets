<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const microchipNumber = defineModel('microchipNumber', {
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
        <span>'s microchip number?</span>
      </p>
    </div>
    <div>
      <FloatLabel variant="on">
        <InputText id="microchip"
                   v-model="microchipNumber"
                   class="myInput"/>
        <label for="microchip">Microchip Number</label>
      </FloatLabel>
    </div>
  </section>
</template>

<style scoped>

</style>