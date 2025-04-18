<script setup>
import {computed, ref} from "vue";

const petTypes = [
  {
    icon: ['fadl', 'cat'],
    title: 'Cat',
    id: 'kxiat7hgd9d7t04'
  },
  {
    icon: ['fadl', 'dog'],
    title: 'Dog',
    id: 'ds68m4ltt54b83b'
  },
  {
    icon: ['fadl', 'rabbit'],
    title: 'Small furry',
    id: '9sc981tg86690t0'
  },
  {
    icon: ['fadl', 'turtle'],
    title: 'Reptile',
    id: '43a455h8l2x0ch7'
  }
]

const isExpanded = ref(false)

const petTypeId = defineModel()

function isSelected(inPetType) {
  if (inPetType === petTypeId.value) {
    return true
  }
  return false
}

function clickPetType(petT) {
  petTypeId.value = petT.id;
  isExpanded.value = false;
}

</script>

<template>
  <div>
    <div class="d-flex gap-4">

      <template

          v-for="petT in petTypes"
          :key="petT.title">
        <div v-if="isSelected(petT.id) || isExpanded || petTypeId === ''">
          <div
              @click="clickPetType(petT)"
              :class="{'icon-container':true, 'selected': isSelected(petT.id)}"
              style="display: flex; flex-direction: column">

            <div v-if="isSelected(petT.id)"
                 style="position: relative; top: -11px; left: 51px;">
              <font-awesome-layers class="fa-xl">
                <font-awesome-icon :icon="['fas', 'circle']" style="color: white"/>
                <font-awesome-icon :icon="['fas', 'circle-check']"
                                   style="color: seagreen"/>
              </font-awesome-layers>
            </div>

            <font-awesome-icon :icon="petT.icon" size="5x"/>
          </div>
          <div class="text-center mt-2"> {{ petT.title }}</div>
        </div>
      </template>
    </div>

    <div v-if="isExpanded === false && petTypeId !== ''" class="d-flex align-items-center mt-4">
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