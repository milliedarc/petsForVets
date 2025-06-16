<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import {computed} from "vue";
import countries from '@/types/countries.json'
import Country from "@/types/Country";

const props = defineProps<{
  name: string
}>()

const isImported = defineModel('isImported', {
  type: Boolean,
  required: true
})

const importCountryCode = defineModel('importCountryCode', {
  type: String,
  required: false
})


const sortedCountries = computed<Country[]>(() => {
  return countries.countries.toSorted(function (countryA, countryB) {
    if (countryA.name < countryB.name) {
      return -1;
    }
    if (countryA.name > countryB.name) {
      return 1;
    }
    return 0;
  });
})

</script>

<template>
  <section>
    <div>
      <p>Is
        <PetNameFormatted :name="props.name"/>
        <span> imported?</span>
      </p>
    </div>

    <div class="flex items-center gap-2">
      <Checkbox v-model="isImported"
                binary
                inputId="imported" name="imported" value="true"/>
      <label class="ms-3" for="imported">
        Please tick this box if your pet is imported
      </label>
    </div>

    <div class="mt-3" v-if="isImported">

      <FloatLabel variant="on">
        <Select v-model="importCountryCode"
                :options="sortedCountries"
                editable
                inputId="country"
                optionLabel="name"
                optionValue="country_code"
                class="myInput"/>
        <label for="country">Select import country</label>
      </FloatLabel>
    </div>
  </section>
</template>

<style scoped>

</style>