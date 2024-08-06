<template>
  <div class="">

    <DashboardToolbar>
      <template #left>
        <h1 class="flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0">
          <span class="truncate">Site</span>
        </h1>
      </template>
      <template #right>

        <UButton color="gray" label="New Site" trailing-icon="i-heroicons-plus" @click="handleAdd" />
      </template>

    </DashboardToolbar>

    <PortalModal v-model="isNewModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Add a new site" title="New site">
      <UForm :state="site" :validate="validate" :validate-on="['submit']" class="space-y-4" @submit="saveItem">
        <UFormGroup label="Site ID" name="site_id">
          <UInput v-model="site.site_id" autofocus type="text" />
        </UFormGroup>
        <UFormGroup label="Name" name="name">
          <UInput v-model="site.name" autofocus type="text" />
        </UFormGroup>

        <UFormGroup label="Initial Url" name="url">
          <UInput v-model="site.url" type="text" />
        </UFormGroup>

        <UFormGroup label="Number of Crawling Per month" name="no_of_crawling_per_month">
          <USelectMenu v-model="site.no_of_crawling_per_month" :options="[
     { value: 1, label: '1 time' },
      { value: 2, label: '2 times' },
      { value: 3, label: '3 times' }]" :ui-menu="{ select: 'capitalize', option: { base: 'capitalize' } }" />
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
          <UButton color="gray" label="Cancel" variant="ghost" @click="handleAdd" />
          <UButton :loading="loading" color="black" label="Save" type="submit" />
        </div>
      </UForm>
    </PortalModal>
    <PortalModal v-model="isEditModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Edit a site" title="Edit site">
      <UForm :state="site" :validate="validate" :validate-on="['submit']" class="space-y-4" @submit="updateItem">
        <UFormGroup label="Site ID" name="site_id">
          <UInput v-model="site.site_id" autofocus type="text" />
        </UFormGroup>
        <UFormGroup label="Name" name="name">
          <UInput v-model="site.name" autofocus type="text" />
        </UFormGroup>

        <UFormGroup label="Initial Url" name="url">
          <UInput v-model="site.url" type="text" />
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

    <PortalModal v-model="isSecretModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Add a new site" title="Site Secret">
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
        <UButton class="mr-2" color="green" icon="i-heroicons-play" @click="startCrawler(row)"/>
        <UButton color="yellow" icon="i-heroicons-key" @click="handleSecret(row)"/>
        <UDropdown :items="handleAction(row)" position="bottom-end">
          <UButton color="gray" icon="i-heroicons-ellipsis-vertical" variant="ghost" />
        </UDropdown>
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
const page = ref<number>(parseInt(route.query.page as string) || 1);
const limit = 10;
const q = ref<string>('');
const loading = ref<boolean>(false);
const isNewModalOpen = ref<boolean>(false);
const isEditModalOpen = ref<boolean>(false);
const isSecretModalOpen = ref<boolean>(false);
const toast = useToast()
const columns = [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'url', label: 'Url' },
  { key: 'no_of_crawling_per_month', label: 'Monthly Crawling limit' },
  { key: 'vm_config', label: 'VM Config' },
  { key: 'action', label: 'Action' }
];

const site = ref({
  id: null,
  site_id: "",
  name: "",
  url: "",
  no_of_crawling_per_month: 1,
  status: "",
  vm_config: {
    cores:2,
    memory:4096,
    disk:10,
    zone:"asia-northeast1"
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
  if (!state.name) errors.push({ path: 'name', message: 'Please enter a name.' })
  if (!state.url) errors.push({ path: 'url', message: 'Please enter a url.' })
  if (!state.vm_config.cores) errors.push({ path: 'cores', message: 'Please enter cores.' })
  if (!state.vm_config.memory) errors.push({ path: 'memory', message: 'Please enter memory.' })
  if (!state.vm_config.disk) errors.push({ path: 'disk', message: 'Please enter disk.' })
  return errors
}
const { data: item, pending: itemsPending, refresh } = await useSiteApi().findAll({ page, limit });
watch(page, (newValue) => {
  updatePage(newValue);
});

function updatePage(newPage: number) {
  router.push({ query: { ...route.query, page: newPage } });
}

function handleAction (site:any) {
  return [
    [
      {
        label: 'Edit',
        click: () => handleEdit(site)
      },
      // {
      //   label: 'Remove',
      //   labelClass: 'text-red-500 dark:text-red-400',
      //   click: () => handleRemove(site.id)
      // }
    ]
  ]
}

function resetItem(){
  site.value = {
    id: null,
    site_id: "",
    name: "",
    url: "",
    status: "",
    no_of_crawling_per_month:1,
    vm_config: {
      cores:2,
      memory:4096,
      disk:10,
      zone:"asia-northeast1"
    }
  }
}
function handleAdd(){
  isNewModalOpen.value = !isNewModalOpen.value
  resetItem()
}
function handleEdit(site:any) {
  console.log('site',toRaw(site))
  site.value = site
  isEditModalOpen.value = !isEditModalOpen.value;
}
async function handleSecret(site: any) {
  if (site){
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
  site.value.no_of_crawling_per_month = site.value.no_of_crawling_per_month.value
  await useSiteApi().save(site.value);
  loading.value = false
  isNewModalOpen.value = false;
  refresh()
  toast.add({ title: "Site Saved" })
}
async function startCrawler(site:any) {
  loading.value = true
  await useSiteApi().startCrawler(site.site_id).then(res=>{
    if(res.status.value!="error"){
      loading.value = false
      toast.add({ title: "Crawler Started" })
    }
  })
}

async function saveSecret() {
  loading.value = true
  secret.value.secrets = JSON.parse(secret.value.secrets);
  console.log("secret.value", secret.value.secrets);
  await useSiteApi().addSecrets(secret.value);
  loading.value = false
  isSecretModalOpen.value = false;
  toast.add({ title: "Secret Saved" })
}

async function updateItem() {
  loading.value = true
  await useSiteApi().update(site.value, site.value.site_id);
  isEditModalOpen.value = false;
  resetItem()
  loading.value = false
  refresh()
  toast.add({ title: "Site Saved" })

}
</script>
