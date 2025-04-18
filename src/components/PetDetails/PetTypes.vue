<script setup>
import {computed, ref} from "vue";

const petTypes = [
  {
    icon: ['fadl', 'cat'],
    title: 'Cat',
    speciesName: 'Cat'
  },
  {
    icon: ['fadl', 'dog'],
    title: 'Dog',
    speciesName: 'Dog'
  },
  {
    icon: ['fadl', 'rabbit'],
    title: 'Rabbit',
    speciesName: 'Rabbit'
  },
  {
    icon: ['fadl', 'turtle'],
    title: 'Exotic',
    speciesName: 'Exotic'
  }
]

const isExpanded = ref(false)

const petType = defineModel()

function isSelected(inPetType) {
  if (inPetType === petType.value) {
    return true
  }
  return false
}

function clickPetType(petT) {
  petType.value = petT.speciesName;
  isExpanded.value = false;
}

</script>

<template>
  <div>
    <div class="d-flex gap-4">

      <template

          v-for="petT in petTypes"
          :key="petT.title">

        <div v-if="isSelected(petT.speciesName) || isExpanded"
             @click="clickPetType(petT)"
             :class="{'icon-container':true, 'selected': isSelected(petT.speciesName)}"
             style="display: flex; flex-direction: column">

          <div v-if="isSelected(petT.speciesName)"
               style="position: relative; top: 20px; left: 63px;">
            <font-awesome-layers class="fa-xl">
              <font-awesome-icon :icon="['fas', 'circle']" style="color: white"/>
              <font-awesome-icon :icon="['fas', 'circle-check']"
                                 style="color: seagreen"/>
            </font-awesome-layers>
          </div>

          <font-awesome-icon :icon="petT.icon" size="5x"/>
          <div class="text-center mt-2"> {{ petT.title }}</div>
        </div>
      </template>
    </div>

    <div v-if="isExpanded === false" class="d-flex align-items-center mt-4">
      <i class="pi pi-angle-down"></i>
      <button @click="isExpanded = true"
              class="btn btn-link px-0" style="color: black">
        <span class="fw-bold">Show all pet types</span>
      </button>
    </div>

  </div>
</template>

<style scoped>

.icon-container {
  width: 150px; /* Or any fixed size */
  height: 150px;
  border-radius: 50%; /* Makes it circular */
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0; /* Light gray background */
  transition: transform 0.3s, background-color 0.3s;
  cursor: pointer;
}

.icon-container:hover {
  transform: scale(1.05); /* Slight zoom on hover */
  background-color: #e0e0e0; /* Darker background on hover */
}

.icon-container.selected {
  border-color: seagreen;
  border-width: 4px;
  border-style: solid;
}
</style>