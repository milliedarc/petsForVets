<script setup lang="ts">
import {ref} from "vue";
import PetIcon from "@/components/PetIcon.vue";
import {useRouter} from "vue-router";
import ButtonPencilEdit from "@/components/ButtonPencilEdit.vue";

// ****************** CONSTS *********************

const router = useRouter()

// ****************** REFS *********************

const menu = ref();

const items = ref([
  {
    label: 'Manage your pet',
    items: [
      {
        label: 'Repeat prescriptions',
        icon: 'pi pi-cart-plus',
        command: () => {
          goToPrescriptionsView()
        }
      },
      {
        label: 'Appointments',
        icon: 'pi pi-calendar-clock'
      }
    ]
  }
]);

// ****************** PROPS *********************

const props = defineProps<{
  pet: Pet
}>();

// ****************** FUNCTIONS *********************

const toggle = (event) => {
  menu.value.toggle(event);
};

function calculateAge(dob: string): string {
  if (dob === '') {
    return ''
  }

  const birthDate = new Date(dob);
  const today = new Date();

  let years = today.getFullYear() - birthDate.getFullYear();
  let months = today.getMonth() - birthDate.getMonth();
  const dayDiff = today.getDate() - birthDate.getDate();

  // Adjust if birthday hasn't occurred yet this month
  if (dayDiff < 0) {
    months--;
  }

  // Adjust if month difference is negative
  if (months < 0) {
    years--;
    months += 12;
  }

  return `${years} years ${months} months`;
}

function goToPetView(): void {
  router.push(`/pets/${props.pet.id}`);
}

function goToPrescriptionsView(): void {
  router.push(`/pets/${props.pet.id}/prescriptions`)
}

</script>

<template>
  <div class="card mb-3">
    <div class="card-body">
      <div class="d-flex justify-content-between align-items-center">
        <div class="d-flex">
          <div>
            <PetIcon :pet="pet"/>
          </div>
          <div class="ms-4">
            <h5>{{ pet.name }}</h5>
            <div>
              <div class="mb-2">{{ calculateAge(pet.date_of_birth) }}</div>
            </div>
          </div>
        </div>
        <div class="d-flex justify-content-center">
          <ButtonPencilEdit @click="goToPetView"/>
          <div>
            <Button type="button"
                    icon="pi pi-ellipsis-v"
                    @click="toggle"
                    size="small"
                    severity="secondary"
                    aria-haspopup="true"
                    aria-controls="overlay_menu"/>
            <Menu ref="menu" id="overlay_menu" :model="items" :popup="true"/>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>