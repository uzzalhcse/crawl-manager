<template>
  <div class="">

    <DashboardToolbar>
      <template #left>
        <h1 class="flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0">
          <span class="truncate">Crawling History</span>
        </h1>
      </template>

    </DashboardToolbar>

    <UTable
      :columns="columns"
      :loading="itemsPending"
      :progress="{ color: 'primary', animation: 'carousel' }"
      :rows="item"
      :ui="{ divide: 'divide-gray-200 dark:divide-gray-800' }"
      class="w-full"
      sort-mode="manual"
    >
      <template #logs-data="{ row }">
        <UButton class="mr-2" color="red" icon="i-heroicons-clipboard-document-list" @click="isLogModalOpen=true" />
        <UButton v-if="row.status != 'running'" :to="`https://console.cloud.google.com/storage/browser/gen_crawled_data_venturas_asia-northeast1/maker/${row.site_id}/logs`" icon="i-heroicons-arrow-top-right-on-square"  target="_blank"></UButton>
        <UBadge v-else color="red" label="N/A" />
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

    <PortalModal v-model="isLogModalOpen" :ui="{ width: 'sm:max-w-md' }" title="Summary Logs" prevent-close>
      <pre>
        {{logs}}
      </pre>
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
  { key: 'action', label: 'Action' },
  { key: 'logs', label: 'Logs'}
];

const logs = [
  {
    msg:"Crawling Category",
    time:"2024-09-01 1:30:00"
  },
  {
    msg:"Total 15 Category Found",
    time:"2024-09-01 3:15:00"
  },
  {
    msg:"Crawling SubCategory",
    time:"2024-09-01 3:15:00"
  },
  {
    msg:"Total 256 SubCategory Found",
    time:"2024-09-01 7:20:00"
  },
  {
    msg:"Crawling Products",
    time:"2024-09-01 7:20:00"
  },
  {
    msg:"Total 5025 Products Found",
    time:"2024-09-01 16:30:00"
  },
]

const { data: item, pending: itemsPending, refresh } = await useSiteApi().crawlingHistory();
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
</script>
