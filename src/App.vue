<script setup lang="ts">
import {ref, watchEffect} from "vue";
import {useRoute, useRouter} from 'vue-router'
import {useSession} from "@/stores/session";
import AppSidebar from "@/components/AppSidebar.vue";
import AppBreadcrumbs from "@/components/AppBreadcrumbs.vue";

// Who is signed in, and redirecting people who aren't, is handled by the
// router guard in src/router/router.ts before any page is shown.
const session = useSession();
const router = useRouter();
const route = useRoute();

// Name the browser tab after the current page (the last breadcrumb), e.g.
// "Prescriptions · Pets for Vets". Reactive, so it updates once a pet's name loads.
watchEffect(() => {
  const current = route.meta.breadcrumbs?.(route).at(-1)?.label;
  document.title = current ? `${current} · Pets for Vets` : 'Pets for Vets';
});

/** The sidebar drawer on phones and tablets. */
const isDrawerOpen = ref<boolean>(false);

async function logOut(): Promise<void> {
  isDrawerOpen.value = false;
  await session.logout();
  await router.push({name: 'UserLogin'});
}
</script>

<template>
  <div v-if="session.user.value" class="app-shell">
    <!-- Wide screens: sidebar always visible. -->
    <aside class="app-sidebar">
      <AppSidebar :user="session.user.value" @logout="logOut"/>
    </aside>

    <!-- Narrow screens: the sidebar opens as a drawer from the navbar's menu button. -->
    <Drawer v-model:visible="isDrawerOpen" class="app-drawer" aria-label="Menu">
      <template #container>
        <AppSidebar :user="session.user.value"
                    @logout="logOut"
                    @navigate="isDrawerOpen = false"/>
      </template>
    </Drawer>

    <div class="app-content">
      <header class="app-navbar">
        <Button icon="pi pi-bars"
                class="menu-button"
                variant="text"
                severity="secondary"
                aria-label="Open menu"
                :aria-expanded="isDrawerOpen"
                @click="isDrawerOpen = true"/>
        <AppBreadcrumbs/>
      </header>

      <!-- Keyed on the path so going from one pet to another loads the new pet.
           (Not fullPath: pages that keep state in ?query=, like the patient
           search, would be rebuilt on every keystroke.) -->
      <RouterView :key="$route.path"/>
    </div>
  </div>

  <!-- Signed out (sign in, register): no sidebar. -->
  <RouterView v-else :key="$route.path"/>

  <!--  PrimeVue components to be used in other components -->
  <Toast/>
  <ConfirmDialog/>
</template>

<style scoped>
.app-sidebar {
  display: none;
}

.app-navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 56px;
  padding: 8px var(--pfv-gutter) 8px 8px;
  background: var(--p-content-background);
  border-bottom: 1px solid var(--p-content-border-color);
}

@media (min-width: 992px) {
  .app-shell {
    --sidebar-width: 260px;
  }

  .app-sidebar {
    display: block;
    position: fixed;
    inset: 0 auto 0 0;
    width: var(--sidebar-width);
    overflow-y: auto;
    background: var(--p-content-background);
    border-right: 1px solid var(--p-content-border-color);
  }

  .app-content {
    margin-left: var(--sidebar-width);
  }

  .app-navbar {
    padding-left: 24px;
  }

  .menu-button {
    display: none;
  }
}
</style>

<style>
/* The drawer is teleported to <body>, outside this component's scoped styles. */
.app-drawer.p-drawer {
  width: min(300px, 85vw);
}
</style>
