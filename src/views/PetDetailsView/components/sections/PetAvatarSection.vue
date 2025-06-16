<script setup lang="ts">
import {computed, ref} from "vue";
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const tempAvatarDisplay = ref(null)
const deleteAvatarFlag = ref<boolean>(false)
const fileUploadRef = ref();

const props = defineProps<{
  name: string,
  avatarUrl: string
}>()

const emit = defineEmits(['deleteClicked', 'fileSelected'])

const hasAvatarUrl = computed<boolean>(() => {
  return props.avatarUrl !== null && props.avatarUrl !== '' && !tempAvatarDisplay.value
})

function onFileSelect(event): void {
  const file = event.files[0];
  const reader = new FileReader();

  reader.onload = async (e) => {
    tempAvatarDisplay.value = e.target.result;
  };
  reader.readAsDataURL(file);
  deleteAvatarFlag.value = false;
  emit('fileSelected', file);
}

function deletePetAvatar(): void {
  deleteAvatarFlag.value = true;
  tempAvatarDisplay.value = null;
  fileUploadRef.value?.clear(); // This clears the filename and resets input
  emit('deleteClicked')
}

</script>
<template>
  <section>
    <div>
      <p class="mb-0">Upload a picture of
        <PetNameFormatted :name="props.name"/>
        <span>:</span>
      </p>
      <p class="fw-light" style="font-size: 0.8rem">
        You can skip this step and add it later.
      </p>
    </div>

    <div class="card">
      <div class="card-body">
        <div class="d-flex flex-column align-items-center">

          <div v-if="!deleteAvatarFlag"
               style="position: relative">
            <img v-if="hasAvatarUrl"
                 :src="avatarUrl"
                 :alt="`${props.name}'s uploaded picture`"
                 class="shadow-lg rounded mb-3"
                 style="max-width: 240px"/>
            <img v-else-if="tempAvatarDisplay"
                 :src="tempAvatarDisplay"
                 :alt="`Selected avatar for ${props.name}`"
                 class="shadow-lg rounded mb-3"
                 style="max-width: 240px"/>

            <Button v-if="hasAvatarUrl || tempAvatarDisplay" icon="pi pi-times"
                    @click="deletePetAvatar"
                    severity="danger" size="small" rounded
                    :aria-label="`Remove uploaded photo of ${props.name}`"
                    style="position: absolute; top: -15px; right: -15px;"/>
          </div>
          <FileUpload ref="fileUploadRef"
                      mode="basic" @select="onFileSelect" customUpload
                      severity="secondary"
                      class="p-button-outlined">
          </FileUpload>
        </div>
      </div>
    </div>
  </section>

</template>

<style scoped>

</style>