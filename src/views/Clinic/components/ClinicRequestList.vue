<script setup lang="ts">
import type {ClinicRequest, RequestAction} from "@/types/Clinic";
import {computed, ref} from "vue";
import {useToast} from "primevue/usetoast";
import {updateClinicRequest} from "@/api/clinic";
import {errorMessage} from "@/api/client";
import {formatTimeAgo} from "@/utils/dates";
import {requestStatusInfo} from "@/views/Clinic/requestStatus";

// ****************** PROPS / EMITS *********************

const props = withDefaults(defineProps<{
  requests: ClinicRequest[]
  /** Show each request's status tag (for lists mixing statuses). */
  showStatus?: boolean
  /** Show the pet and owner (hidden on a patient's own page). */
  showPatient?: boolean
}>(), {showStatus: false, showPatient: true})

const emit = defineEmits<{
  /** A request changed status; the parent decides whether to keep showing it. */
  (e: 'updated', request: ClinicRequest): void
}>()

// ****************** CONSTS *********************

const toast = useToast();

const successMessages: Record<RequestAction, {summary: string, detail: (request: ClinicRequest) => string}> = {
  Approved: {
    summary: 'Request approved',
    detail: (request) => `${request.prescription.medicine_name} for ${request.pet.name} is now being prepared.`,
  },
  Ready: {
    summary: 'Ready to collect',
    detail: (request) => `${request.pet.name}'s ${request.prescription.medicine_name} is ready for collection.`,
  },
  Rejected: {
    summary: 'Request rejected',
    detail: (request) => `${request.owner.name || 'The owner'} will see your reason.`,
  },
};

// ****************** REFS *********************

/** Id of the request whose action is in progress, so only its buttons spin. */
const busyId = ref<string>('');
const busyAction = ref<RequestAction | ''>('');

const requestToReject = ref<ClinicRequest | null>(null);
const rejectionReason = ref<string>('');

// ****************** COMPUTED *********************

const isRejectDialogOpen = computed<boolean>({
  get: () => requestToReject.value !== null,
  set: (open) => {
    if (!open) {
      requestToReject.value = null;
    }
  },
});

// ****************** FUNCTIONS *********************

function ownerName(request: ClinicRequest): string {
  return [request.owner.name, request.owner.surname].filter(Boolean).join(' ') || request.owner.email;
}

async function changeStatus(request: ClinicRequest, status: RequestAction, reason = ''): Promise<boolean> {
  busyId.value = request.id;
  busyAction.value = status;
  try {
    const updated = await updateClinicRequest(request.id, status, reason);
    const message = successMessages[status];
    toast.add({severity: 'success', summary: message.summary, detail: message.detail(updated), life: 4000});
    emit('updated', updated);
    return true;
  } catch (error) {
    console.error(error);
    toast.add({severity: 'error', summary: 'Could not update the request', detail: errorMessage(error, 'Please try again.'), life: 6000});
    return false;
  } finally {
    busyId.value = '';
    busyAction.value = '';
  }
}

function openRejectDialog(request: ClinicRequest): void {
  requestToReject.value = request;
  rejectionReason.value = '';
}

async function confirmReject(): Promise<void> {
  const request = requestToReject.value;
  if (!request || rejectionReason.value.trim() === '') {
    return;
  }
  if (await changeStatus(request, 'Rejected', rejectionReason.value.trim())) {
    requestToReject.value = null;
  }
}

function isBusy(request: ClinicRequest, action: RequestAction): boolean {
  return busyId.value === request.id && busyAction.value === action;
}

</script>

<template>
  <ul class="request-list">
    <li v-for="request in props.requests" :key="request.id" class="request">
      <div class="request-main">
        <div class="request-title">
          <span class="medicine">{{ request.prescription.medicine_name }}</span>
          <Tag v-if="props.showStatus"
               :value="requestStatusInfo[request.status].label"
               :severity="requestStatusInfo[request.status].severity"/>
        </div>
        <div class="dosage">
          {{ request.prescription.dosage }} · Quantity {{ request.prescription.quantity }}
        </div>

        <div v-if="props.showPatient" class="patient">
          <RouterLink :to="{name: 'ClinicPatient', params: {id: request.pet.id}}" class="pet-link">
            <i class="pi pi-heart" aria-hidden="true"></i>{{ request.pet.name }}
          </RouterLink>
          <span class="text-muted">({{ request.pet.species_name }})</span>
          <span class="owner">
            <i class="pi pi-user" aria-hidden="true"></i>{{ ownerName(request) }}
            <a :href="`mailto:${request.owner.email}`" class="owner-email">{{ request.owner.email }}</a>
          </span>
        </div>

        <div class="requested" :title="new Date(request.created_at).toLocaleString()">
          Requested {{ formatTimeAgo(request.created_at) }}
        </div>

        <Message v-if="request.rejection_reason" severity="error" size="small" class="mt-2">
          <span class="font-bold">Reason given:</span> {{ request.rejection_reason }}
        </Message>
      </div>

      <div v-if="request.status === 'Pending' || request.status === 'Approved'" class="request-actions">
        <Button v-if="request.status === 'Pending'"
                label="Approve"
                icon="pi pi-check"
                size="small"
                :loading="isBusy(request, 'Approved')"
                :disabled="busyId !== ''"
                @click="changeStatus(request, 'Approved')"/>
        <Button v-else
                label="Ready to collect"
                icon="pi pi-box"
                size="small"
                :loading="isBusy(request, 'Ready')"
                :disabled="busyId !== ''"
                @click="changeStatus(request, 'Ready')"/>
        <Button label="Reject"
                size="small"
                severity="secondary"
                variant="outlined"
                :disabled="busyId !== ''"
                :aria-label="`Reject ${request.pet.name}'s request for ${request.prescription.medicine_name}`"
                @click="openRejectDialog(request)"/>
      </div>
    </li>
  </ul>

  <Dialog v-model:visible="isRejectDialogOpen"
          modal
          header="Reject request"
          :style="{width: 'min(480px, calc(100vw - 32px))'}">
    <form v-if="requestToReject" class="reject-form" @submit.prevent="confirmReject">
      <p class="m-0">
        Rejecting <span class="font-bold">{{ requestToReject.prescription.medicine_name }}</span>
        for <span class="font-bold">{{ requestToReject.pet.name }}</span>.
        {{ requestToReject.owner.name || 'The owner' }} will see your reason on their prescriptions page.
      </p>
      <label for="rejectionReason" class="font-bold">Reason</label>
      <Textarea id="rejectionReason"
                v-model="rejectionReason"
                rows="3"
                maxlength="500"
                autoResize
                autofocus
                placeholder="e.g. Bella is due a check-up first. Please book an appointment."
                fluid/>
      <div class="flex justify-between items-center gap-2 flex-wrap">
        <small class="text-muted">{{ rejectionReason.length }}/500</small>
        <div class="flex gap-2">
          <Button label="Cancel" severity="secondary" variant="outlined" @click="requestToReject = null"/>
          <Button type="submit"
                  label="Reject request"
                  severity="danger"
                  :disabled="rejectionReason.trim() === ''"
                  :loading="isBusy(requestToReject, 'Rejected')"/>
        </div>
      </div>
    </form>
  </Dialog>
</template>

<style scoped>
.request-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.request {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px 24px;
  padding: 16px 20px;
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--pfv-radius);
}

.request-main {
  flex: 1 1 320px;
  min-width: 0;
}

.request-title {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.medicine {
  font-weight: 700;
  font-size: 1.05rem;
}

.dosage {
  color: var(--p-text-muted-color);
  font-size: 0.9rem;
}

.patient {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px 10px;
  margin-top: 8px;
  font-size: 0.9rem;
}

.patient i {
  font-size: 0.8rem;
  margin-right: 6px;
}

.pet-link {
  font-weight: 700;
}

.owner {
  display: inline-flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 2px 8px;
}

.owner-email {
  color: var(--p-text-muted-color);
  overflow-wrap: anywhere;
}

.requested {
  margin-top: 6px;
  font-size: 0.8rem;
  color: var(--p-text-muted-color);
}

.request-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.reject-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
