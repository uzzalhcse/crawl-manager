<template>
  <div class="">

    <DashboardToolbar>
      <template #left>
        <h1 class="flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0">
          <span class="truncate">Crawling History</span>
        </h1>
        <div class="">
          <small>Recent crawling history from the last 15 days, including running crawlers</small>
        </div>
      </template>
      <template #right>
        <div class="flex px-3 py-3.5 border-gray-200 dark:border-gray-700">
          <UInput v-model="q" placeholder="Filter History..." />
        </div>
      </template>

    </DashboardToolbar>
    <UTable
      :columns="columns"
      :loading="itemsPending"
      :progress="{ color: 'primary', animation: 'carousel' }"
      :rows="filteredRows"
      :ui="{ divide: 'divide-gray-200 dark:divide-gray-800' }"
      class="w-full"
      sort-mode="manual"
    >
      <template #logs-data="{ row }">

        <UTooltip text="Summary Logs" :popper="{ arrow: true }">
          <UButton class="mr-2" color="red" icon="i-heroicons-clipboard-document-list" @click="handleSummaryLog(row)" />
        </UTooltip>
        <UTooltip text="GCP Log explorer" :popper="{ arrow: true }">
          <UButton :to="row.log_url" icon="i-heroicons-arrow-top-right-on-square"  target="_blank"></UButton>
        </UTooltip>
      </template>
      <template #action-data="{ row }">

        <UPopover v-if="row.status == 'running'" class="inline-flex mr-2" overlay>
          <UTooltip text="Stop Crawler" :popper="{ arrow: true }">
            <UButton color="red" icon="i-heroicons-stop"/>
          </UTooltip>
          <template #panel="{ close }">
            <UCard class="max-w-xs mx-auto flex flex-col items-center">
              <!-- Icon and Message -->
              <div class="flex items-center">
                <i class="i-heroicons-exclamation-triangle text-yellow mr-2"></i>
                <div class="font-semibold">
                  Do you want to proceed with this action?
                </div>
              </div>

              <!-- Buttons -->
              <div class="mt-4 flex justify-end space-x-4 w-full">
                <!-- Cancel Button -->
                <UButton
                    :disabled="loading"
                    text="true"
                    label="No, Thanks"
                    size="2xs"
                    @click="close"
                />
                <!-- Confirm Button -->
                <UButton
                    :loading="loading"
                    label="OK"
                    color="yellow"
                    size="2xs"
                    @click="() => { stopCrawler(row); }"
                />
              </div>
            </UCard>
          </template>
        </UPopover>
        <UTooltip text="Restart Crawler (upcoming)" :popper="{ arrow: true }">
          <UButton :disabled="true" class="mr-2" color="yellow" icon="i-heroicons-arrow-path" @click="stopCrawler(row)"/>
        </UTooltip>
      </template>
      <template #empty-state>
        <div class="flex flex-col items-center justify-center py-6 gap-3">
          <span class="italic text-sm">No Crawling History Found!</span>
        </div>
      </template>
    </UTable>

    <PortalModal v-model="isLogModalOpen" :ui="{ width: 'max-w-xl' }" title="Summary Logs" prevent-close>
      <UTable
          :columns="summaryColumns"
          :loading="itemsPending"
          :progress="{ color: 'primary', animation: 'carousel' }"
          :rows="crawlingSummary"
          :ui="{ divide: 'divide-gray-200 dark:divide-gray-800' }"
          class="w-full"
          sort-mode="manual"
      >
        <template #errors-data="{ row }">
<!--          <span v-for="error in row.errors" class="text-red-500"> {{ error.url }}, </span>-->
          <h4>upcoming</h4>
        </template>
        <template #empty-state>
          <div class="flex flex-col items-center justify-center py-6 gap-3">
            <span class="italic text-sm">No Crawling Summary Found!</span>
          </div>
        </template>
      </UTable>
    </PortalModal>
  </div>
</template>


<script lang="ts" setup>

import type { FormError } from '#ui/types';
import type { Site } from '~/types';
import {useSiteApi} from "~/composables/useSIteApi";

const route = useRoute();
const router = useRouter();
const q = ref<string>('');
const loading = ref<boolean>(false);
const isLogModalOpen = ref<boolean>(false);
const toast = useToast()
const columns = [
  { key: 'site_id', label: 'Site', sortable: true },
  { key: 'instance_name', label: 'Instance Name', sortable: true },
  { key: 'status', label: 'Status' ,sortable: true},
  { key: 'start_date', label: 'Start Date',sortable: true },
  { key: 'end_date', label: 'End Date',sortable: true },
  { key: 'initiate_by', label: 'Initiate by',sortable: true },
  { key: 'action', label: 'Action' },
  { key: 'logs', label: 'Logs'}
];
const summaryColumns = [
  { key: 'collection_name', label: 'Collection',sortable: true },
  { key: 'data_count', label: 'Data' },
  { key: 'error_count', label: 'Error',sortable: true},
  { key: 'created_at', label: 'Date',sortable: true },
  { key: 'errors', label: 'Errors'}
];

const crawlingSummary = ref<any>([])
const filteredRows = computed(() => {
  if (!q.value) {
    return items.value; // Return all items if search query is empty
  }

  return items.value.filter((site: any) => {
    return Object.values(site).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });
});
const { data: items, pending: itemsPending, refresh } = await useSiteApi().crawlingHistory();
async function stopCrawler(history:any) {
  loading.value = true
  await useSiteApi().stopCrawler(history.instance_name).then(res=>{
    if(res.status.value!="error"){
      loading.value = false
      toast.add({ title: "Crawler Stopped" })
      refresh()
    }
  })
}
async function handleSummaryLog(row:any) {
  isLogModalOpen.value = true
  await useSiteApi().getCrawlerSummary(row.instance_name).then(res=>{
    crawlingSummary.value = res.data.value
  })
}
</script>
