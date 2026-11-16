<script setup lang="ts">
import {computed} from "vue";
import {useRoute} from "vue-router";
import type {Crumb} from "@/router/router";

const route = useRoute();

/** The trail for the current page, from its route's meta (see src/router/router.ts). */
const crumbs = computed<Crumb[]>(() => route.meta.breadcrumbs?.(route) ?? []);
</script>

<template>
  <nav v-if="crumbs.length > 0" aria-label="Breadcrumb" class="breadcrumbs">
    <ol>
      <li v-for="(crumb, index) in crumbs" :key="index">
        <i v-if="index > 0" class="pi pi-angle-right separator" aria-hidden="true"></i>
        <span v-if="index === crumbs.length - 1" class="current" aria-current="page">{{ crumb.label }}</span>
        <RouterLink v-else-if="crumb.to" :to="crumb.to">{{ crumb.label }}</RouterLink>
        <span v-else class="plain">{{ crumb.label }}</span>
      </li>
    </ol>
  </nav>
</template>

<style scoped>
.breadcrumbs {
  min-width: 0; /* allow shrinking inside the navbar */
}

ol {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
  padding: 0;
  list-style: none;
  font-size: 0.9rem;
  white-space: nowrap;
}

li {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.separator {
  font-size: 0.75rem;
  color: var(--p-text-muted-color);
}

a, a:visited {
  color: var(--p-text-muted-color);
  font-weight: 600;
}

a:hover {
  color: var(--p-primary-color);
}

.plain {
  color: var(--p-text-muted-color);
}

.current {
  font-weight: 700;
  color: var(--p-text-color);
  overflow: hidden;
  text-overflow: ellipsis;
}

/* On phones, only the last two steps fit: hide the earlier ones. */
@media (max-width: 575px) {
  li:not(:nth-last-child(-n + 2)) {
    display: none;
  }

  li:nth-last-child(2) .separator {
    display: none;
  }
}
</style>
