<script setup lang="ts">
import {ref} from "vue";
import PetTypeIds from "@/types/PetTypeIds";

const emit = defineEmits(["petTypeClicked"]);

const petTypes = [
  {
    icon: ['fadl', 'cat'],
    title: 'Cat',
    id: PetTypeIds.catId,
  },
  {
    icon: ['fadl', 'dog'],
    title: 'Dog',
    id: PetTypeIds.dogId,
  },
  {
    icon: ['fadl', 'rabbit'],
    title: 'Small Animal',
    id: PetTypeIds.smallAnimalId,
  },
  {
    icon: ['fadl', 'turtle'],
    title: 'Reptile',
    id: PetTypeIds.reptileId,
  }
]

const isExpanded = ref<boolean>(false)

const petTypeId = defineModel<string>()

function isSelected(inPetType: string): boolean {
  return (inPetType === petTypeId.value)
}

function clickPetType(petT): void {
  petTypeId.value = petT.id;
  isExpanded.value = false;
  emit("petTypeClicked", petT);
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
  width: 150px;
  height: 150px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0;
  transition: transform 0.3s, background-color 0.3s;
  cursor: pointer;
}

.icon-container:hover {
  transform: scale(1.05);
  background-color: #e0e0e0;
}

.icon-container.selected {
  border-color: seagreen;
  border-width: 4px;
  border-style: solid;
}
</style>