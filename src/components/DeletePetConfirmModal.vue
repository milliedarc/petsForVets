<script setup lang="ts">
import Pet from "./Pet.js"
import PocketBase from 'pocketbase';
import {inject} from "vue";

const bootstrap = inject('bootstrap');

const pb = new PocketBase('http://127.0.0.1:8090');

const petToEdit = defineModel<Pet>()

const emit = defineEmits(['deleted']);

async function deletePet() {
  await pb.collection('pets').delete(petToEdit.value.id);
// Get the modal instance using Bootstrap's library and hide it
  // @ts-ignore
  const myModal = bootstrap.Modal.getInstance(document.getElementById('deletePetModal'));
  if (myModal) {
    myModal.hide();
    emit("deleted")
  }
}

</script>

<template>
  <div class="modal fade" id="deletePetModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Are you sure?</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <p>Deleting your pet cannot be undone.</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
          <button @click="deletePet"
                  type="button" class="btn btn-danger">Yes, delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>