<script setup lang="ts">
import type {Pet} from "@/types/Pet";
import {computed} from "vue";
import PetIcon from "@/components/PetIcon.vue";
import {formatAge} from "@/utils/dates";

// ****************** PROPS *********************

const props = defineProps<{
  pet: Pet
}>();

// ****************** COMPUTED *********************

/** e.g. "Dog · Labrador Retriever" */
const speciesAndBreed = computed<string>(() => {
  return [props.pet.species.name, props.pet.breed?.name]
      .filter(Boolean)
      .join(' · ');
})

/** e.g. "3 years, 2 months · Female · Neutered" */
const facts = computed<string[]>(() => {
  const list: string[] = [];
  list.push(formatAge(props.pet.date_of_birth) || 'Age unknown');
  if (props.pet.gender) {
    list.push(props.pet.gender === 'male' ? 'Male' : 'Female');
  }
  if (props.pet.neutered) {
    list.push('Neutered');
  }
  return list;
})

const editPath = computed<string>(() => `/pets/${props.pet.id}`);
const prescriptionsPath = computed<string>(() => `/pets/${props.pet.id}/prescriptions`);

</script>

<template>
  <Card class="pet-card">
    <template #content>
      <div class="pet-summary">
        <PetIcon :pet="pet"/>
        <div class="pet-text">
          <h2 class="pet-name">
            <RouterLink :to="editPath">{{ pet.name }}</RouterLink>
          </h2>
          <div class="pet-species">{{ speciesAndBreed }}</div>
          <div class="pet-facts">{{ facts.join(' · ') }}</div>
        </div>
      </div>
    </template>
    <template #footer>
      <div class="pet-actions">
        <RouterLink :to="prescriptionsPath" custom v-slot="{navigate}">
          <Button label="Prescriptions"
                  icon="pi pi-receipt"
                  size="small"
                  @click="navigate"
                  :aria-label="`${pet.name}'s repeat prescriptions`"/>
        </RouterLink>
        <RouterLink :to="editPath" custom v-slot="{navigate}">
          <Button label="Edit details"
                  icon="pi pi-pencil"
                  size="small"
                  severity="secondary"
                  variant="outlined"
                  @click="navigate"
                  :aria-label="`Edit ${pet.name}'s details`"/>
        </RouterLink>
      </div>
    </template>
  </Card>
</template>

<style scoped>
.pet-card {
  height: 100%;
}

.pet-summary {
  display: flex;
  align-items: center;
  gap: 16px;
}

.pet-text {
  min-width: 0; /* lets long names wrap instead of overflowing the card */
}

.pet-name {
  font-size: 1.25rem;
  margin: 0 0 2px;
  overflow-wrap: anywhere;
}

.pet-name a {
  color: var(--p-text-color);
}

.pet-species {
  font-weight: 600;
  color: var(--p-primary-color);
}

.pet-facts {
  font-size: 0.9rem;
  color: var(--p-text-muted-color);
}

.pet-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>
