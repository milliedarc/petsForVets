<script setup lang="ts">
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import {ref} from "vue";

const props = defineProps<{
  name: string
}>()

const dateOfBirth = defineModel<Date | null>('dateOfBirth', {})

const years = ref<number | null>(null)
const months = ref<number | null>(null)
const today = new Date()

/**
 * Turns an estimated age in years and months into a date of birth
 * (the 1st of the month, since we only know the month).
 */
function calculateBirthDate(): void {
  if (years.value === null && months.value === null) {
    return;
  }
  const birthDate = new Date();
  birthDate.setDate(1);
  birthDate.setHours(0, 0, 0, 0);
  birthDate.setFullYear(birthDate.getFullYear() - (years.value ?? 0));
  birthDate.setMonth(birthDate.getMonth() - (months.value ?? 0));

  dateOfBirth.value = birthDate;
}

</script>

<template>
  <section class="question">
    <span class="question-label">
      How old is
      <PetNameFormatted :name="props.name"/>?
    </span>

    <Tabs value="dobDate" class="age-tabs">
      <TabList>
        <Tab value="dobDate">Date of birth</Tab>
        <Tab value="dobAge">Estimated age</Tab>
      </TabList>

      <TabPanels>
        <TabPanel value="dobDate">
          <label for="dob" class="sr-only">Date of birth</label>
          <DatePicker v-model="dateOfBirth"
                      showIcon fluid
                      iconDisplay="input"
                      inputId="dob"
                      :maxDate="today"
                      placeholder="dd/mm/yyyy"
                      dateFormat="dd/mm/yy"/>
        </TabPanel>

        <TabPanel value="dobAge">
          <div class="age-inputs">
            <FloatLabel variant="on">
              <InputNumber v-model="years" inputId="years" :min="0" :max="50" fluid
                           @update:modelValue="calculateBirthDate"/>
              <label for="years">Years</label>
            </FloatLabel>
            <FloatLabel variant="on">
              <InputNumber v-model="months" inputId="months" :min="0" :max="11" fluid
                           @update:modelValue="calculateBirthDate"/>
              <label for="months">Months</label>
            </FloatLabel>
          </div>
        </TabPanel>
      </TabPanels>
    </Tabs>
    <small class="field-hint">
      Don't know the exact date? Choose <span class="font-bold">Estimated age</span> and give us your best guess.
    </small>
  </section>
</template>

<style scoped>
.age-tabs {
  max-width: var(--pfv-field-width);
}

.age-tabs :deep(.p-tabpanels) {
  padding-left: 0;
  padding-right: 0;
}

.age-inputs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
</style>
