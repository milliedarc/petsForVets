<script setup lang="ts">
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const props = defineProps<{
  name: string
}>()

const dobTab = defineModel('dobTab', {
  type: String,
  default: ''
})

const dob = defineModel('dob', {
  type: Date
})

const months = defineModel('months', {
  type: String
})

const years = defineModel('years', {
  type: String
})


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

      <Tabs v-model:value="dobTab" style="width: 350px">
        <TabList>
          <Tab value="dobDate">Date of birth</Tab>
          <Tab value="dobAge">Age in years</Tab>
        </TabList>

        <TabPanels>
          <TabPanel value="dobDate">
            <div class="d-flex">
              <label for="dob" class="font-bold block mb-2"></label>
              <DatePicker v-model="dob" showIcon fluid iconDisplay="input" inputId="dob" class="myInput"/>
            </div>
          </TabPanel>

          <TabPanel value="dobAge">
            <div class="d-flex gap-3">
              <FloatLabel variant="on">
                <InputText type="number" min="0" max="30" id="years" style="width: 150px"
                           v-model="years"/>
                <label for="years">Years *</label>
              </FloatLabel>
              <FloatLabel variant="on">
                <InputText type="number" min="0" max="11" id="months" style="width: 150px"
                           v-model="months"/>
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