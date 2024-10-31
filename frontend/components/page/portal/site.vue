<template>
  <div class="">

    <DashboardToolbar>
      <template #left>
        <h1 class="flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0">
          <span class="truncate">Site</span>
        </h1>
      </template>
      <template #right>
        <UInput v-model="q" placeholder="Filter site..." class="ml-auto" />
        <UButton color="gray" label="New Site" trailing-icon="i-heroicons-plus" @click="handleAdd" />
      </template>

    </DashboardToolbar>

    <PortalModal v-model="isNewModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Add a new site" title="New site" prevent-close>
      <UForm :state="site" :validate="validate" :validate-on="['submit']" class="space-y-4" @submit="saveItem">
        <UFormGroup label="Site ID" name="site_id">
          <UInput v-model="site.site_id" autofocus type="text" />
        </UFormGroup>

        <UFormGroup label="Initial Url" name="url">
          <UInput v-model="site.url" type="text" />
        </UFormGroup>
        <UFormGroup label="Git Branch (only for staging)" name="git_branch">
          <UInput v-model="site.git_branch" type="text" placeholder="LZ-XX" />
        </UFormGroup>
        <UCheckbox v-model="site.use_proxy" name="use_proxy" label="Use Webshare Proxy" />

        <UFormGroup v-if="site.use_proxy" label="Number of Proxy" name="num_of_proxy">
          <UInput v-model="site.number_of_proxies" type="number" placeholder="5" />
        </UFormGroup>

        <UFormGroup label="Frequency" name="frequency">
          <UInput v-model="site.frequency" type="text" />
          <template #description>
            Schedules are specified using unix-cron format. E.g. every minute: "* * * * *", every 3 hours: "0 */3 * * *", every Monday at 9:00: "0 9 * * 1"
            <UButton size="2xs" to="https://cloud.google.com/scheduler/docs/configuring/cron-job-schedules" icon="i-heroicons-arrow-top-right-on-square"  target="_blank">Learn More</UButton>
          </template>
        </UFormGroup>

        <UFormGroup label="Status" name="status">
          <USelectMenu v-model="site.status" :options="['active', 'inactive']" :ui-menu="{ select: 'capitalize', option: { base: 'capitalize' } }" />
        </UFormGroup>
<!--        <UFormGroup label="Number of Crawling Per month" name="no_of_crawling_per_month">-->
<!--          <USelectMenu v-model="site.no_of_crawling_per_month" :options="[-->
<!--     { value: 1, label: '1 time' },-->
<!--      { value: 2, label: '2 times' },-->
<!--      { value: 3, label: '3 times' }]" :ui-menu="{ select: 'capitalize', option: { base: 'capitalize' } }" />-->
<!--        </UFormGroup>-->
        <h2 class="font-black">VM Config</h2>
        <UFormGroup label="Cores" name="cores">
          <UInput v-model.number="site.vm_config.cores" type="number" />
        </UFormGroup>

        <UFormGroup label="Memory (MB)" name="memory">
          <UInput v-model.number="site.vm_config.memory" type="number" />
        </UFormGroup>

        <UFormGroup label="Disk Sze (GB)" name="disk">
          <UInput v-model.number="site.vm_config.disk" type="number" />
        </UFormGroup>

        <UFormGroup label="Zone" name="zone">
          <USelectMenu v-model="site.vm_config.zone" :options="['asia-northeast1-a', 'asia-northeast1-b', 'asia-northeast1-c']" :ui-menu="{ select: 'capitalize', option: { base: 'capitalize' } }" />
        </UFormGroup>

        <div class="flex justify-end gap-3">
          <UButton color="gray" label="Cancel" variant="ghost" @click="handleAdd" />
          <UButton :loading="loading" color="black" label="Save" type="submit" />
        </div>
      </UForm>
    </PortalModal>
    <PortalModal v-model="isEditModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Edit a site" title="Edit site" prevent-close>
      <UForm :state="site" :validate="validate" :validate-on="['submit']" class="space-y-4" @submit="updateItem">
        <UFormGroup label="Site ID" name="site_id">
          <UInput v-model="site.site_id" autofocus type="text"/>
        </UFormGroup>

        <UFormGroup label="Initial Url" name="url">
          <UInput v-model="site.url" type="text" />
        </UFormGroup>

        <UFormGroup label="Git Branch (only for staging)" name="git_branch">
          <UInput v-model="site.git_branch" type="text" placeholder="LZ-XX" />
        </UFormGroup>
        <UCheckbox v-model="site.use_proxy" name="use_proxy" label="Use Webshare Proxy" />

        <UFormGroup v-if="site.use_proxy" label="Number of Proxy" name="num_of_proxy">
          <UInput v-model="site.number_of_proxies" type="number" placeholder="5" />
        </UFormGroup>
        <UFormGroup label="Frequency" name="frequency">
          <UInput v-model="site.frequency" type="text" :disabled="true" />
          <template #description>
            Schedules are specified using unix-cron format. E.g. every minute: "* * * * *", every 3 hours: "0 */3 * * *", every Monday at 9:00: "0 9 * * 1"
            <UButton size="2xs" to="https://cloud.google.com/scheduler/docs/configuring/cron-job-schedules" icon="i-heroicons-arrow-top-right-on-square"  target="_blank">Learn More</UButton>
          </template>
        </UFormGroup>

        <UFormGroup label="Status" name="status">
          <USelectMenu v-model="site.status" :options="['active', 'inactive']" :ui-menu="{ select: 'capitalize', option: { base: 'capitalize' } }" />
        </UFormGroup>
        <h2 class="font-black">VM Config</h2>
        <UFormGroup label="Cores" name="cores">
          <UInput v-model.number="site.vm_config.cores" type="number" />
        </UFormGroup>

        <UFormGroup label="Memory (MB)" name="memory">
          <UInput v-model.number="site.vm_config.memory" type="number" />
        </UFormGroup>

        <UFormGroup label="Disk Sze (GB)" name="disk">
          <UInput v-model.number="site.vm_config.disk" type="number" />
        </UFormGroup>

        <UFormGroup label="Zone" name="zone">
          <USelectMenu v-model="site.vm_config.zone" :options="['asia-northeast1-a', 'asia-northeast1-b', 'asia-northeast1-c']" :ui-menu="{ select: 'capitalize', option: { base: 'capitalize' } }" />
        </UFormGroup>

        <div class="flex justify-end gap-3">
          <UButton color="gray" label="Cancel" variant="ghost" @click="handleEdit" />
          <UButton :loading="loading" color="black" label="Save" type="submit" />
        </div>
      </UForm>
    </PortalModal>

    <PortalModal v-model="isSecretModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Add a new site" title="Site Secret" prevent-close>
      <UForm :state="secret" class="space-y-4" @submit="saveSecret">
        <UFormGroup label="Secrets (Json Key:Value)" name="secrets">
          <UTextarea v-model="secret.secrets" resize type="text" />
        </UFormGroup>

        <div class="flex justify-end gap-3">
          <UButton color="gray" label="Cancel" variant="ghost" @click="handleSecret" />
          <UButton :loading="loading" color="black" label="Save" type="submit" />
        </div>
      </UForm>
    </PortalModal>
    <UTable
      :columns="columns"
      :loading="itemsPending"
      :progress="{ color: 'primary', animation: 'carousel' }"
      :rows="filteredRows"
      :ui="{ divide: 'divide-gray-200 dark:divide-gray-800' }"
      class="w-full"
      sort-mode="manual"
    >
      <template #site_id-data="{ row }">
        <span :class="row.status === 'inactive' ? 'text-red-500' : ''">{{row.site_id}}</span>
        <UTooltip v-if="row.use_proxy" text="Number of used proxy" :popper="{ arrow: true }">
          <UKbd class="ml-2">{{ row.number_of_proxies }}</UKbd>
        </UTooltip>
      </template>
      <template #action-data="{ row }">
        <UPopover class="inline-flex mr-2" overlay>
          <UTooltip text="Run Crawler" :popper="{ arrow: true }">
            <UButton color="green" icon="i-heroicons-play" />
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
                    @click="() => { startCrawler(row).finally(() => close()); }"
                />
              </div>
            </UCard>
          </template>
        </UPopover>

        <UTooltip class="mr-2" text="Build Binary" :popper="{ arrow: true }">
          <UButton color="blue" icon="i-heroicons-wrench-screwdriver" :loading="row.loading" @click="buildCrawler(row)"/>
        </UTooltip>
        <UTooltip class="mr-2" text="Env Secrets" :popper="{ arrow: true }">
          <UButton color="yellow" icon="i-heroicons-key" @click="handleSecret(row)"/>
        </UTooltip>
        <UTooltip text="Edit" :popper="{ arrow: true }">
          <UButton color="orange" icon="i-heroicons-pencil-square" @click="handleEdit(row)"/>
        </UTooltip>
      </template>

      <template #empty-state>
        <div class="flex flex-col items-center justify-center py-6 gap-3">
          <span class="italic text-sm">No Sites Found!</span>
          <UButton label="Add New Site" @click="handleAdd" />
        </div>
      </template>
    </UTable>
  </div>
</template>


<script lang="ts" setup>

import type { FormError } from '#ui/types';
import type { Site } from '~/types';
import {useSiteApi} from "~/composables/useSIteApi";

const route = useRoute();
const router = useRouter();
const page = ref<number>(parseInt(route.query.page as string) || 1);
const limit = 10;
const q = ref<string>('');
const loading = ref<boolean>(false);
const isNewModalOpen = ref<boolean>(false);
const isEditModalOpen = ref<boolean>(false);
const isSecretModalOpen = ref<boolean>(false);
const toast = useToast()
const columns = [
  { key: 'site_id', label: 'Name', sortable: true },
  { key: 'url', label: 'Url' },
  { key: 'status', label: 'Status', sortable: true },
  { key: 'git_branch', label: 'Git Branch', sortable: true },
  { key: 'frequency', label: 'Crawling Frequency' },
  { key: 'vm_config', label: 'VM Config' },
  { key: 'action', label: 'Action' }
];

const site = ref({
  id: null,
  site_id: "",
  name: "",
  url: "",
  use_proxy: false,
  number_of_proxies: 0,
  frequency: "",
  status: "active",
  git_branch: "dev",
  vm_config: {
    cores:2,
    memory:4096,
    disk:10,
    zone:"asia-northeast1-a"
  }
})
const secret = ref({
  site_id: "",
  secrets: ""
})

// https://ui.nuxt.com/components/form
const validate = (state: Site): FormError[] => {
  const errors = []
  if (!state.site_id) errors.push({ path: 'site_id', message: 'Please enter a site_id.' })
  if (!state.url) errors.push({ path: 'url', message: 'Please enter a url.' })
  if (!state.frequency) errors.push({ path: 'frequency', message: 'Please enter a frequency.' })
  if (!state.vm_config.cores || state.vm_config.cores < 2) errors.push({ path: 'cores', message: 'Please enter valid cores.' })
  if (!state.vm_config.memory) errors.push({ path: 'memory', message: 'Please enter memory.' })
  if (!state.vm_config.disk) errors.push({ path: 'disk', message: 'Please enter disk.' })
  return errors
}
const filteredRows = computed(() => {
  if (!q.value) {
    return items.value; // Return all items if search query is empty
  }

  return items.value.filter((site: Site) => {
    return Object.values(site).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });
});
const { data: items, pending: itemsPending, refresh } = await useSiteApi().findAll({ page, limit });
watch(page, (newValue) => {
  updatePage(newValue);
});

function updatePage(newPage: number) {
  router.push({ query: { ...route.query, page: newPage } });
}


function resetItem(){
  site.value = {
    id: null,
    site_id: "",
    name: "",
    url: "",
    use_proxy: false,
    number_of_proxies: 0,
    status: "active",
    git_branch: "dev",
    frequency:"0 0 1 * *",
    vm_config: {
      cores:2,
      memory:4096,
      disk:10,
      zone:"asia-northeast1-a"
    }
  }
}
function handleAdd(){
  isNewModalOpen.value = !isNewModalOpen.value
  resetItem()
}
function handleEdit(row:any) {
  if (row.site_id){
    site.value = row
  }
  isEditModalOpen.value = !isEditModalOpen.value;
}
async function buildCrawler(site: any) {
  if (site.site_id){
    site.loading = true
    useSiteApi().buildCrawler(site.site_id).then(res=>{
      console.log('res', res.data.value)
      toast.add({ title: "Binary Build Success" })
      site.loading = false
    })
  }

}
async function handleSecret(site: any) {
  if (site.site_id){
    useSiteApi().getSecrets(site.site_id).then(res=>{
      console.log('res', res.data.value)
      secret.value =res.data.value
      secret.value.secrets =JSON.stringify(res.data.value.secrets)
    })
  }

  isSecretModalOpen.value = !isSecretModalOpen.value;

}

async function handleRemove(id: number) {
  // Add prompt confirm dialog
  if (confirm("This will remove permanently")) {
    await useSiteApi().remove(id)
    refresh()
    toast.add({ title: "Site Deleted" })
  }
}

async function saveItem() {
  loading.value = true
  site.value.number_of_proxies = Number(site.value.number_of_proxies)
  useSiteApi().save(site.value).then(res=>{
    if(res.status.value!="error"){
      isNewModalOpen.value = false;
      refresh()
      toast.add({ title: "Site Saved" })
    }
    loading.value = false
  })
}
async function startCrawler(site:any) {
  loading.value = true
  await useSiteApi().startCrawler(site.site_id).then(res=>{
    if(res.status.value!="error"){
      toast.add({ title: "Crawler Started" })
    }
    loading.value = false
  })
}

async function saveSecret() {
  loading.value = true

  if(secret.value.secrets.length<3) {
    secret.value.secrets = "{}"
  }
  secret.value.secrets = JSON.parse(secret.value.secrets);
  console.log("secret.value", secret.value.secrets);
  useSiteApi().addSecrets(secret.value).then(res=>{
    if(res.status.value!="error"){
      isSecretModalOpen.value = false;
      toast.add({ title: "Secret Saved" })
    }
    loading.value = false
  })
}

async function updateItem() {
  loading.value = true
  site.value.number_of_proxies = Number(site.value.number_of_proxies)
  useSiteApi().update(site.value, site.value.site_id).then(res=>{
    if(res.status.value!="error"){
      isEditModalOpen.value = false;
      resetItem()
      refresh()
      toast.add({ title: "Site Updated" })
    }
    loading.value = false
  })

}
</script>
