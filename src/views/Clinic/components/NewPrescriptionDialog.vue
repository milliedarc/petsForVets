<script setup lang="ts">
import type {ClinicPet, ClinicPrescription} from "@/types/Clinic";
import type {Medicine} from "@/types/Prescription";
import {computed, ref, watch} from "vue";
import {useToast} from "primevue/usetoast";
import {createPrescription, listMedicines} from "@/api/clinic";
import {errorMessage} from "@/api/client";

// ****************** PROPS / MODELS / EMITS *********************

const props = defineProps<{
  pet: ClinicPet
}>()

const visible = defineModel<boolean>('visible', {required: true})

const emit = defineEmits<{
  (e: 'created', prescription: ClinicPrescription): void
}>()

// ****************** CONSTS / REFS *********************

const toast = useToast();

const medicines = ref<Medicine[]>([]);
const medicinesError = ref<string>('');
const medicineId = ref<string>('');
const dosage = ref<string>('');
const quantity = ref<number | null>(1);
const isSaving = ref<boolean>(false);
const saveError = ref<string>('');

// ****************** COMPUTED *********************

const canSave = computed<boolean>(() => {
  return medicineId.value !== ''
      && dosage.value.trim() !== ''
      && quantity.value !== null && quantity.value >= 1
      && !isSaving.value;
});

// ****************** FUNCTIONS *********************

async function loadMedicines(): Promise<void> {
  if (medicines.value.length > 0) {
    return;
  }
  try {
    medicines.value = await listMedicines();
    medicinesError.value = '';
  } catch (error) {
    console.error(error);
    medicinesError.value = errorMessage(error, 'The medicine list could not be loaded.');
  }
}

function resetForm(): void {
  medicineId.value = '';
  dosage.value = '';
  quantity.value = 1;
  saveError.value = '';
}

async function save(): Promise<void> {
  if (!canSave.value) {
    return;
  }
  isSaving.value = true;
  saveError.value = '';
  try {
    const prescription = await createPrescription(props.pet.id, {
      medicine_id: medicineId.value,
      dosage: dosage.value.trim(),
      quantity: quantity.value ?? 1,
    });
    emit('created', prescription);
    toast.add({
      severity: 'success',
      summary: 'Prescription added',
      detail: `${props.pet.owner.name || 'The owner'} can now request repeats of ${prescription.medicine.name} for ${props.pet.name}.`,
      life: 5000,
    });
    visible.value = false;
  } catch (error) {
    console.error(error);
    saveError.value = errorMessage(error, 'The prescription could not be saved. Please try again.');
  } finally {
    isSaving.value = false;
  }
}

// ****************** WATCHERS *********************

// Start each time with a clean form and make sure the medicines are loaded.
watch(visible, (open) => {
  if (open) {
    resetForm();
    loadMedicines();
  }
});

</script>

<template>
  <Dialog v-model:visible="visible"
          modal
          :header="`New prescription for ${props.pet.name}`"
          :style="{width: 'min(520px, calc(100vw - 32px))'}">
    <form class="prescription-form" @submit.prevent="save">
      <Message v-if="medicinesError" severity="error" size="small">{{ medicinesError }}</Message>

      <div class="field">
        <label for="medicine">Medicine</label>
        <Select v-model="medicineId"
                inputId="medicine"
                :options="medicines"
                optionLabel="name"
                optionValue="id"
                placeholder="Choose a medicine"
                :loading="medicines.length === 0 && !medicinesError"
                filter
                filterPlaceholder="Search medicines"
                fluid>
          <template #option="{option}">
            <div>
              <div class="font-bold">{{ option.name }}</div>
              <div class="text-sm text-muted">{{ option.description }}</div>
            </div>
          </template>
        </Select>
      </div>

      <div class="field">
        <label for="dosage">Dosage instructions</label>
        <InputText id="dosage"
                   v-model="dosage"
                   maxlength="200"
                   placeholder="e.g. One tablet twice a day with food"
                   aria-describedby="dosageHint"
                   fluid/>
        <small id="dosageHint" class="field-hint m-0">The owner sees this on their prescriptions page.</small>
      </div>

      <div class="field quantity-field">
        <label for="quantity">Quantity per repeat</label>
        <InputNumber v-model="quantity"
                     inputId="quantity"
                     :min="1"
                     :max="1000"
                     showButtons
                     fluid/>
      </div>

      <Message v-if="saveError" severity="error" size="small" role="alert">{{ saveError }}</Message>

      <div class="actions">
        <Button label="Cancel" severity="secondary" variant="outlined" @click="visible = false"/>
        <Button type="submit" label="Add prescription" icon="pi pi-check" :disabled="!canSave" :loading="isSaving"/>
      </div>
    </form>
  </Dialog>
</template>

<style scoped>
.prescription-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field label {
  font-weight: 600;
}

.quantity-field {
  max-width: 200px;
}

.actions {
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
}
</style>
