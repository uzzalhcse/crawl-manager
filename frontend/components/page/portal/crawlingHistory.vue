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
      v-if="item"
      :columns="columns"
      :loading="itemsPending"
      :progress="{ color: 'primary', animation: 'carousel' }"
      :rows="item"
      :ui="{ divide: 'divide-gray-200 dark:divide-gray-800' }"
      class="w-full"
      sort-mode="manual"
    >
      <template #action-data="{ row }">
        <UButton class="mr-2" color="green" icon="i-heroicons-pause-circle" @click="stopCrawler(row)"/>
      </template>
    </UTable>
    <h3 v-else class="text-center mt-5">No Data Found</h3>
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
const toast = useToast()
const columns = [
  { key: 'site_id', label: 'Site', sortable: true },
  { key: 'instance_name', label: 'Instance Name', sortable: true },
  { key: 'logs', label: 'Logs'},
  { key: 'status', label: 'Status' ,sortable: true},
  { key: 'start_date', label: 'Start Date',sortable: true },
  { key: 'end_date', label: 'End Date',sortable: true },
  { key: 'action', label: 'Action' }
];


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
