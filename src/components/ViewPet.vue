<script setup lang="ts">
// imports, consts, models, props, emits, refs, computed, watchers, functions, hooks
import {useRoute} from "vue-router";
import PocketBase from "pocketbase";
import {onMounted, ref} from "vue";

const pb = new PocketBase('http://127.0.0.1:8090');

const pet = ref<Pet | undefined>(undefined);
const petNotFound = ref(false)
const isLoading = ref(true)
const route = useRoute()

async function fetchPet() {
  const result = await pb.collection('pets').getOne(route.params.id, {
    expand: 'species,breed,breed_secondary'
  })
  // @ts-ignore
  pet.value = result; // assigns content to 'pets'
  console.log(result)
}

onMounted(async () => {
  try {
    await fetchPet()
  } catch (error) {
    petNotFound.value = true
  } finally {
    isLoading.value = false
  }
})

</script>

<template>
  <div v-if="!isLoading">
    <div v-if="petNotFound">
      <h1>Pet not found</h1>
    </div>
    <div v-else>
      <h1>View & edit your pet {{ route.params.id }}</h1>
      <h2> {{ pet?.name }}</h2>
    </div>
  </div>
</template>

<style scoped>

</style>