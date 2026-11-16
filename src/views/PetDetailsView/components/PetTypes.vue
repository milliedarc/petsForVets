<script setup lang="ts">
import {ref} from "vue";
import PetTypeIds from "@/types/PetTypeIds";
import {Cat, Dog, Rabbit, Turtle, CircleCheck} from "@lucide/vue";

const emit = defineEmits(["petTypeClicked"]);

const petTypes = [
  {
    icon: Cat,
    title: 'Cat',
    id: PetTypeIds.catId,
  },
  {
    icon: Dog,
    title: 'Dog',
    id: PetTypeIds.dogId,
  },
  {
    icon: Rabbit,
    title: 'Small Animal',
    id: PetTypeIds.smallAnimalId,
  },
  {
    icon: Turtle,
    title: 'Reptile',
    id: PetTypeIds.reptileId,
  }
]

const isExpanded = ref<boolean>(false)

const petTypeId = defineModel<string>()

function isSelected(inPetType: string): boolean {
  return (inPetType === petTypeId.value)
}

function clickPetType(petT: typeof petTypes[number]): void {
  petTypeId.value = petT.id;
  isExpanded.value = false;
  emit("petTypeClicked", petT);
}

</script>

<template>
  <div>
    <div class="pet-type-list">
      <template
          v-for="petT in petTypes"
          :key="petT.title">
        <button v-if="isSelected(petT.id) || isExpanded || petTypeId === ''"
                type="button"
                class="pet-type"
                :class="{'selected': isSelected(petT.id)}"
                :aria-pressed="isSelected(petT.id)"
                @click="clickPetType(petT)">
          <span class="icon-circle">
            <span v-if="isSelected(petT.id)" class="selected-indicator" aria-hidden="true">
              <CircleCheck :size="22"/>
            </span>
            <component :is="petT.icon" :size="56" aria-hidden="true"/>
          </span>
          <span class="pet-type-title">{{ petT.title }}</span>
        </button>
      </template>
    </div>

    <Button v-if="isExpanded === false && petTypeId !== ''"
            @click="isExpanded = true"
            label="Change pet type"
            icon="pi pi-angle-down"
            variant="link"
            class="change-type"/>
  </div>
</template>

<style scoped>
.pet-type-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.pet-type {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 0;
  border: 0;
  background: none;
  font: inherit;
  font-weight: 600;
  color: var(--p-text-color);
  cursor: pointer;
}

.icon-circle {
  position: relative; /* anchor for .selected-indicator */
  width: 110px;
  height: 110px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--p-content-hover-background);
  color: var(--p-text-muted-color);
  border: 3px solid transparent;
  transition: transform 0.2s, background-color 0.2s, border-color 0.2s;
}

.pet-type:hover .icon-circle {
  transform: scale(1.05);
  background-color: var(--pfv-accent-soft);
}

.pet-type:focus-visible .icon-circle {
  outline: 2px solid var(--p-primary-color);
  outline-offset: 3px;
}

.pet-type.selected .icon-circle {
  border-color: var(--p-primary-color);
  background-color: var(--pfv-accent-soft);
  color: var(--p-primary-700);
}

.selected-indicator {
  position: absolute;
  top: 2px;
  right: 2px;
  display: flex;
  color: var(--p-primary-color);
  background-color: white;
  border-radius: 50%;
}

.change-type {
  padding-left: 0;
  margin-top: 8px;
}
</style>
