<template>
  <div class="">

    <DashboardToolbar>
      <template #left>
        <h1 class="flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0">
          <span class="truncate">Proxy Management</span>
        </h1>
      </template>
      <template #right>
        <UInput v-model="q" placeholder="Filter Server..." class="ml-auto" />
        <UButton color="gray" label="New Proxy" trailing-icon="i-heroicons-plus" @click="handleAdd" />
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
          <span class="italic text-sm">No Proxy Server Found!</span>
          <UButton label="Add New Proxy" @click="handleAdd" />
        </div>
      </template>
    </UTable>

    <PortalModal v-model="isNewModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Add a new Proxy" title="New Proxy" prevent-close>
      <UForm :state="proxy" :validate="validate" :validate-on="['submit']" class="space-y-4" @submit="saveItem">
        <UFormGroup label="Server" name="server">
          <small>With a scheme (http://, https://, etc.)</small>
          <UInput v-model="proxy.server" autofocus type="text" placeholder="http://1.2.3.4:1010"/>
        </UFormGroup>
        <UFormGroup label="Username" name="username">
          <UInput v-model="proxy.username" autofocus type="text" />
        </UFormGroup>
        <UFormGroup label="Password" name="password">
          <UInput v-model="proxy.password" autofocus type="text" />
        </UFormGroup>


        <UFormGroup label="Status" name="status">
          <USelectMenu v-model="proxy.status" :options="['active', 'inactive']" :ui-menu="{ select: 'capitalize', option: { base: 'capitalize' } }" />
        </UFormGroup>

        <div class="flex justify-end gap-3">
          <UButton color="gray" label="Cancel" variant="ghost" @click="handleAdd" />
          <UButton :loading="loading" color="black" label="Save" type="submit" />
        </div>
      </UForm>
    </PortalModal>
  </div>
</template>


<script lang="ts" setup>

import type { FormError } from '#ui/types';
import type { Server } from '~/types';
import {useSiteApi} from "~/composables/useSIteApi";

const route = useRoute();
const router = useRouter();
const q = ref<string>('');
const loading = ref<boolean>(false);
const isNewModalOpen = ref<boolean>(false);
const toast = useToast()
const columns = [
  { key: 'server', label: 'Server', sortable: true },
  { key: 'status', label: 'Status' ,sortable: true},
  { key: 'uses', label: 'Uses' ,sortable: true},
  { key: 'action', label: 'Action' },
];
const validate = (state: Server): FormError[] => {
  const errors = []
  if (!state.server) errors.push({ path: 'server', message: 'Please enter valid server.' })
  if (!state.username) errors.push({ path: 'username', message: 'Please enter valid username.' })
  if (!state.password) errors.push({ path: 'password', message: 'Please enter valid password.' })
  return errors
}

const proxy = ref({
  server: "",
  username: "lnvmpyru",
  password: "5un1tb1azapa",
  status: 'active',
})
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
const { data: items, pending: itemsPending, refresh } = await useSiteApi().proxyList();
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

function resetItem(){
  proxy.value = {
    server: "",
    username: "lnvmpyru",
    password: "5un1tb1azapa",
    status: "active",
  }
}
function handleAdd(){
  isNewModalOpen.value = !isNewModalOpen.value
  resetItem()
}
async function saveItem() {
  loading.value = true
  useSiteApi().save(proxy.value).then(res=>{
    if(res.status.value!="error"){
      isNewModalOpen.value = false;
      refresh()
      toast.add({ title: "Site Saved" })
    }
    loading.value = false
  })
}
</script>
