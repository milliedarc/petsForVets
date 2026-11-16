<script setup lang="ts">

import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import countries from '@/types/countries.json'
import type {Country} from "@/types/Country";

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

// The list never changes, so sort it once rather than in a computed.
const sortedCountries: Country[] = countries.countries.toSorted(
    (countryA: Country, countryB: Country) => countryA.name.localeCompare(countryB.name)
);

</script>

<template>
  <section class="question">
    <span class="question-label">
      Was
      <PetNameFormatted :name="props.name"/>
      imported from another country?
    </span>
    <div class="flex items-center gap-3">
      <ToggleSwitch v-model="isImported" inputId="imported"/>
      <label for="imported">{{ isImported ? 'Yes, imported' : 'No' }}</label>
    </div>

    <div v-if="isImported" class="mt-4">
      <label for="country" class="sr-only">Country imported from</label>
      <Select v-model="importCountryCode"
              :options="sortedCountries"
              inputId="country"
              optionLabel="name"
              optionValue="country_code"
              placeholder="Country imported from"
              filter
              filterPlaceholder="Search countries"
              resetFilterOnHide
              class="field-input"/>
    </div>
  </section>
</template>
