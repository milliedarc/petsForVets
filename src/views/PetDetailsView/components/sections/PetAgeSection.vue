<script setup lang="ts">
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";
import {ref} from "vue";

const props = defineProps<{
  name: string
}>()

const dateOfBirth = defineModel<Date | null>('dateOfBirth', {})

const years = ref<string>('')
const months = ref<string>('')

/**
 * Turns years and months into a date
 */
function calculateBirthDate(): void {
  const today = new Date();
  const birthDate = new Date(today);

  let yearsLocal = 0;
  if (!isNaN(Number(years.value))) {
    yearsLocal = Number(years.value);
  }

  let monthsLocal = 0;
  if (!isNaN(Number(months.value))) {
    monthsLocal = Number(months.value);
  }

  birthDate.setFullYear(birthDate.getFullYear() - yearsLocal);
  birthDate.setMonth(birthDate.getMonth() - monthsLocal);

  birthDate.setDate(1);

  birthDate.setHours(0, 0, 0, 0);

  dateOfBirth.value = birthDate;
}

</script>

<template>
  <section>
    <div>
      <p class="mb-1">How old is
        <PetNameFormatted :name="props.name"/>
        <span>?</span>
      </p>
    </div>
    <div>

      <Tabs value="dobDate" style="width: 350px">
        <TabList>
          <Tab value="dobDate">Date of birth</Tab>
          <Tab value="dobAge">Age in years</Tab>
        </TabList>

        <TabPanels>
          <TabPanel value="dobDate">
            <div class="d-flex">
              <label for="dob" class="font-bold block mb-2"></label>
              <DatePicker v-model="dateOfBirth"
                          showIcon fluid
                          iconDisplay="input"
                          inputId="dob"
                          class="myInput"
                          dateFormat="dd/mm/yy"/>
            </div>
          </TabPanel>

          <TabPanel value="dobAge">
            <div class="d-flex gap-3">
              <FloatLabel variant="on">
                <InputText type="number" min="0" id="years" style="width: 150px"
                           v-model="years"
                           @blur="calculateBirthDate"
                />
                <label for="years">Years *</label>
              </FloatLabel>
              <FloatLabel variant="on">
                <InputText type="number" min="0" id="months" style="width: 150px"
                           v-model="months"
                           @blur="calculateBirthDate"
                />
                <label for="months">Months *</label>
              </FloatLabel>
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
      <p class="fw-light" style="font-size: 0.8rem">
        Don't worry if you don't know the exact date. Choose
        <span class="fw-bold">Age in years</span>
        instead and give us an estimate!
      </p>
    </div>
  </section>
</template>

<style scoped>

</style>