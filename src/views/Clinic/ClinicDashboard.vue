<script setup lang="ts">
import type {ClinicRequest, ClinicSummary} from "@/types/Clinic";
import type {PrescriptionRequestStatus} from "@/types/PrescriptionRequest";
import {computed, onMounted, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
import {getClinicSummary, listClinicRequests} from "@/api/clinic";
import {errorMessage} from "@/api/client";
import {useSession} from "@/stores/session";
import {requestStatusInfo} from "@/views/Clinic/requestStatus";
import ClinicRequestList from "@/views/Clinic/components/ClinicRequestList.vue";

// ****************** CONSTS *********************

const route = useRoute();
const router = useRouter();
const {user} = useSession();

/** The queue's tabs, in the order a request moves through them. */
const tabs: PrescriptionRequestStatus[] = ['Pending', 'Approved', 'Ready', 'Rejected'];

const emptyMessages: Record<PrescriptionRequestStatus, string> = {
  Pending: 'No requests waiting for review. Nice work!',
  Approved: 'Nothing being prepared at the moment.',
  Ready: 'No prescriptions waiting to be collected.',
  Rejected: 'No rejected requests.',
};

// ****************** REFS *********************

const summary = ref<ClinicSummary>();
const summaryError = ref<string>('');

const requests = ref<ClinicRequest[]>([]);
const isLoadingRequests = ref<boolean>(true);
const requestsError = ref<string>('');

// ****************** COMPUTED *********************

/** The open tab lives in the URL (?tab=Approved), so it survives back/forward and reloads. */
const activeTab = computed<PrescriptionRequestStatus>({
  get: () => {
    const tab = route.query.tab as PrescriptionRequestStatus;
    return tabs.includes(tab) ? tab : 'Pending';
  },
  set: (tab) => {
    router.replace({query: tab === 'Pending' ? {} : {tab}});
  },
});

const greeting = computed<string>(() => {
  const hour = new Date().getHours();
  const timeOfDay = hour < 12 ? 'Good morning' : hour < 18 ? 'Good afternoon' : 'Good evening';
  return user.value?.name ? `${timeOfDay}, ${user.value.name}` : timeOfDay;
});

/** Finished requests (ready, rejected) are history: newest first. The queue is oldest first. */
const sortedRequests = computed<ClinicRequest[]>(() => {
  return activeTab.value === 'Ready' || activeTab.value === 'Rejected'
      ? requests.value.toReversed()
      : requests.value;
});

// ****************** FUNCTIONS *********************

function countFor(status: PrescriptionRequestStatus): number | undefined {
  if (!summary.value) {
    return undefined;
  }
  return {
    Pending: summary.value.pending_requests,
    Approved: summary.value.approved_requests,
    Ready: summary.value.ready_requests,
    Rejected: undefined,
  }[status];
}

async function fetchSummary(): Promise<void> {
  try {
    summary.value = await getClinicSummary();
    summaryError.value = '';
  } catch (error) {
    console.error(error);
    summaryError.value = errorMessage(error, 'The overview could not be loaded.');
  }
}

async function fetchRequests(): Promise<void> {
  isLoadingRequests.value = true;
  requestsError.value = '';
  try {
    requests.value = await listClinicRequests({status: activeTab.value});
  } catch (error) {
    console.error(error);
    requestsError.value = errorMessage(error, 'The requests could not be loaded.');
  } finally {
    isLoadingRequests.value = false;
  }
}

/** A request changed status, so it no longer belongs in this tab. */
function onRequestUpdated(updated: ClinicRequest): void {
  requests.value = requests.value.filter((request) => request.id !== updated.id);
  fetchSummary();
}

// ****************** WATCHERS / HOOKS *********************

watch(activeTab, fetchRequests);

onMounted(() => {
  fetchSummary();
  fetchRequests();
});

</script>

<template>
  <main class="page">
    <div class="page-header">
      <div>
        <h1>{{ greeting }}</h1>
        <p class="m-0 text-muted">Here's what's waiting for the clinic today.</p>
      </div>
      <RouterLink :to="{name: 'ClinicPatients'}" custom v-slot="{navigate}">
        <Button label="Find a patient" icon="pi pi-search" severity="secondary" variant="outlined" @click="navigate"/>
      </RouterLink>
    </div>

    <Message v-if="summaryError" severity="error" class="mb-4">{{ summaryError }}</Message>

    <!-- Overview tiles. The request tiles double as shortcuts to their tab. -->
    <section class="stats" aria-label="Overview">
      <button v-for="status in (['Pending', 'Approved', 'Ready'] as const)"
              :key="status"
              type="button"
              class="stat"
              :class="[`stat--${requestStatusInfo[status].severity}`, {active: activeTab === status}]"
              :aria-pressed="activeTab === status"
              @click="activeTab = status">
        <span class="stat-icon"><i :class="requestStatusInfo[status].icon" aria-hidden="true"></i></span>
        <span class="stat-value">
          <Skeleton v-if="countFor(status) === undefined" width="2rem" height="1.75rem"/>
          <template v-else>{{ countFor(status) }}</template>
        </span>
        <span class="stat-label">{{ requestStatusInfo[status].label }}</span>
      </button>
      <RouterLink :to="{name: 'ClinicPatients'}" class="stat stat--neutral">
        <span class="stat-icon"><i class="pi pi-heart" aria-hidden="true"></i></span>
        <span class="stat-value">
          <Skeleton v-if="!summary" width="2rem" height="1.75rem"/>
          <template v-else>{{ summary.patients }}</template>
        </span>
        <span class="stat-label">
          Patients<template v-if="summary"> · {{ summary.clients }} {{ summary.clients === 1 ? 'client' : 'clients' }}</template>
        </span>
      </RouterLink>
    </section>

    <section aria-labelledby="queueHeading">
      <h2 id="queueHeading">Repeat prescription requests</h2>

      <Tabs :value="activeTab" @update:value="(tab) => activeTab = tab as PrescriptionRequestStatus" class="queue-tabs">
        <TabList>
          <Tab v-for="status in tabs" :key="status" :value="status">
            {{ requestStatusInfo[status].label }}
            <Badge v-if="countFor(status)"
                   :value="countFor(status)"
                   :severity="status === 'Pending' ? 'warn' : 'secondary'"
                   class="tab-badge"/>
          </Tab>
        </TabList>
      </Tabs>

      <div v-if="isLoadingRequests" class="state-panel" aria-live="polite">
        <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="5"/>
        <span>Loading requests…</span>
      </div>

      <div v-else-if="requestsError" class="state-panel" role="alert">
        <p>{{ requestsError }}</p>
        <Button label="Try again" icon="pi pi-refresh" severity="secondary" @click="fetchRequests"/>
      </div>

      <div v-else-if="requests.length === 0" class="state-panel">
        <i :class="requestStatusInfo[activeTab].icon" class="empty-icon" aria-hidden="true"></i>
        <p>{{ emptyMessages[activeTab] }}</p>
      </div>

      <ClinicRequestList v-else
                         :requests="sortedRequests"
                         @updated="onRequestUpdated"/>
    </section>
  </main>
</template>

<style scoped>
.stats {
  display: grid;
  /* Four across on wide screens, two across on phones. */
  grid-template-columns: repeat(auto-fit, minmax(min(100%, 150px), 1fr));
  gap: 12px;
  margin-bottom: 32px;
}

.stat {
  --stat-color: var(--p-text-muted-color);
  --stat-bg: var(--p-content-hover-background);
  display: grid;
  grid-template-columns: auto 1fr;
  grid-template-rows: auto auto;
  column-gap: 12px;
  align-items: center;
  padding: 16px;
  text-align: left;
  font: inherit;
  color: var(--p-text-color);
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--pfv-radius);
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.stat:hover {
  border-color: var(--stat-color);
}

.stat:focus-visible {
  outline: 2px solid var(--p-primary-color);
  outline-offset: 2px;
}

.stat.active {
  border-color: var(--stat-color);
  box-shadow: inset 0 0 0 1px var(--stat-color);
}

.stat--warn { --stat-color: var(--p-amber-600); --stat-bg: var(--p-amber-50); }
.stat--info { --stat-color: var(--p-sky-600); --stat-bg: var(--p-sky-50); }
.stat--success { --stat-color: var(--p-green-600); --stat-bg: var(--p-green-50); }
.stat--neutral { --stat-color: var(--p-primary-color); --stat-bg: var(--pfv-accent-soft); }

.stat-icon {
  grid-row: span 2;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: var(--stat-bg);
  color: var(--stat-color);
  font-size: 1.15rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 800;
  line-height: 1.1;
}

.stat-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--p-text-muted-color);
}

.queue-tabs {
  margin-bottom: 16px;
}

.queue-tabs :deep(.p-tablist-tab-list) {
  background: transparent;
}

.tab-badge {
  margin-left: 6px;
}

.empty-icon {
  font-size: 1.75rem;
  color: var(--p-primary-color);
}
</style>
