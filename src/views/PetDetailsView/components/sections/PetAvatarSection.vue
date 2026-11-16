<script setup lang="ts">
import {computed, ref} from "vue";
import type {FileUploadSelectEvent} from "primevue/fileupload";
import PetNameFormatted from "@/views/PetDetailsView/components/PetNameFormatted.vue";

const tempAvatarDisplay = ref<string | null>(null)
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

function onFileSelect(event: FileUploadSelectEvent): void {
  const file: File = event.files[0];
  const reader = new FileReader();

  reader.onload = () => {
    tempAvatarDisplay.value = reader.result as string;
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
  <section class="question">
    <span class="question-label">
      Upload a picture of
      <PetNameFormatted :name="props.name"/>
    </span>
    <small class="field-hint mb-4">You can skip this step and add it later.</small>

    <div>
      <div>
        <div class="flex flex-col items-center">

          <div v-if="!deleteAvatarFlag"
               style="position: relative">
            <img v-if="hasAvatarUrl"
                 :src="avatarUrl"
                 :alt="`${props.name}'s uploaded picture`"
                 class="avatar-preview"
                 style="max-width: 240px"/>
            <img v-else-if="tempAvatarDisplay"
                 :src="tempAvatarDisplay"
                 :alt="`Selected avatar for ${props.name}`"
                 class="avatar-preview"
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
.avatar-preview {
  border-radius: var(--pfv-radius);
  margin-bottom: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}
</style>