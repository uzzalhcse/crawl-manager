<template>
  <div class="">

    <DashboardToolbar>
      <template #left>
        <h1 class="flex items-center gap-1.5 font-semibold text-gray-900 dark:text-white min-w-0">
          <span class="truncate">Proxy Management
        <UKbd class="ml-2">{{ filteredRows.length }}</UKbd></span>
        </h1>
      </template>
      <template #right>
        <UInput v-model="q" placeholder="Filter Proxy..." class="ml-auto" />
        <UButton color="gray" label="Sync Proxy" trailing-icon="i-heroicons-plus" @click="syncProxyList" />
<!--        <UButton color="gray" label="New Proxy" trailing-icon="i-heroicons-plus" @click="handleAdd" />-->
      </template>

    </DashboardToolbar>
    <UTable
      :columns="columns"
      :loading="itemsPending || loading"
      :progress="{ color: 'primary', animation: 'carousel' }"
      :rows="filteredRows"
      :ui="{ divide: 'divide-gray-200 dark:divide-gray-800' }"
      class="w-full"
      sort-mode="manual"
    >
      <template #server-data="{ row }">

        <span :class="!row.valid ? 'text-red-500' : 'text-green-500'">{{row.server}}</span>
        <UKbd v-if="row.site_proxies" class="ml-2">{{ row.site_proxies.length }}</UKbd>
      </template>
      <template #valid-data="{ row }">
        <UKbd class="ml-2">{{ row.valid ? 'Valid' : 'Invalid' }}</UKbd>
        <UTooltip v-if="!row.valid" text="Error Logs" :popper="{ arrow: true }">
          <UButton class="ml-2" size="xs" color="green" icon="i-heroicons-clipboard-document-list" @click="handleLog(row)"/>
        </UTooltip>
      </template>
      <template #site_proxies-data="{ row }">
        <UKbd v-for="site in row.site_proxies" :key="site.site_id" class="mr-1">{{ site.site_id }}</UKbd>
      </template>
      <template #action-data="{ row }">

<!--        <UPopover class="inline-flex mr-2" overlay>-->
<!--          <UTooltip text="Delete Proxy" :popper="{ arrow: true }">-->
<!--            <UButton color="red" icon="i-heroicons-trash"/>-->
<!--          </UTooltip>-->
<!--          <template #panel="{ close }">-->
<!--            <UCard class="max-w-xs mx-auto flex flex-col items-center">-->
<!--              &lt;!&ndash; Icon and Message &ndash;&gt;-->
<!--              <div class="flex items-center">-->
<!--                <i class="i-heroicons-exclamation-triangle text-yellow mr-2"></i>-->
<!--                <div class="font-semibold">-->
<!--                  Do you want to proceed with this action?-->
<!--                </div>-->
<!--              </div>-->

<!--              &lt;!&ndash; Buttons &ndash;&gt;-->
<!--              <div class="mt-4 flex justify-end space-x-4 w-full">-->
<!--                &lt;!&ndash; Cancel Button &ndash;&gt;-->
<!--                <UButton-->
<!--                    :disabled="loading"-->
<!--                    text="true"-->
<!--                    label="No, Thanks"-->
<!--                    size="2xs"-->
<!--                    @click="close"-->
<!--                />-->
<!--                &lt;!&ndash; Confirm Button &ndash;&gt;-->
<!--                <UButton-->
<!--                    :loading="loading"-->
<!--                    label="OK"-->
<!--                    color="yellow"-->
<!--                    size="2xs"-->
<!--                    @click="() => { deleteCrawler(row); }"-->
<!--                />-->
<!--              </div>-->
<!--            </UCard>-->
<!--          </template>-->
<!--        </UPopover>-->

<!--        <UTooltip text="Edit" :popper="{ arrow: true }">-->
<!--          <UButton color="orange" icon="i-heroicons-pencil-square" @click="handleEdit(row)"/>-->
<!--        </UTooltip>-->
        <UTooltip v-if="!row.valid" text="Error Logs" :popper="{ arrow: true }">
          <UButton class="ml-2" color="green" icon="i-heroicons-clipboard-document-list" @click="handleLog(row)"/>
        </UTooltip>
        <span v-else>N/A</span>
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
    <PortalModal v-model="isEditModalOpen" :ui="{ width: 'sm:max-w-md' }" description="Add a new Proxy" title="New Proxy" prevent-close>
      <UForm :state="proxy" :validate="validate" :validate-on="['submit']" class="space-y-4" @submit="updateItem">
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
          <UButton color="gray" label="Cancel" variant="ghost" @click="handleEdit" />
          <UButton :loading="loading" color="black" label="Update" type="submit" />
        </div>
      </UForm>
    </PortalModal>
    <PortalModal v-model="isLogModalOpen" :ui="{ width: 'sm:max-w-md' }" title="Error Logs" prevent-close>
      <pre>
        {{proxyLog}}
      </pre>
    </PortalModal>
  </div>
</template>


<script lang="ts" setup>

import type { FormError } from '#ui/types';
import type { Proxy } from '~/types';
import {useSiteApi} from "~/composables/useSIteApi";
import {base64urlEncode} from "iron-webcrypto";

const route = useRoute();
const router = useRouter();
const q = ref<string>('');
const loading = ref<boolean>(false);
const isNewModalOpen = ref<boolean>(false);
const isEditModalOpen = ref<boolean>(false);
const isLogModalOpen = ref<boolean>(false);
const toast = useToast()
const columns = [
  { key: 'server', label: 'Server', sortable: true },
  { key: 'country_code', label: 'country_code' ,sortable: true},
  { key: 'city_name', label: 'city_name' ,sortable: true},
  { key: 'site_proxies', label: 'Sites' ,sortable: true},
  { key: 'last_verification', label: 'last_verification' ,sortable: true},
  { key: 'valid', label: 'Status' ,sortable: true},
  // { key: 'action', label: 'Action' },
];
const validate = (state: Proxy): FormError[] => {
  const errors = []
  if (!state.server) errors.push({ path: 'server', message: 'Please enter valid server.' })
  if (!state.username) errors.push({ path: 'username', message: 'Please enter valid username.' })
  if (!state.password) errors.push({ path: 'password', message: 'Please enter valid password.' })
  return errors
}
const proxyLog = ref("")
const proxy = ref({
  id: "",
  server: "http://",
  username: "lnvmpyru",
  password: "5un1tb1azapa",
  status: 'active',
})
const filteredRows = computed(() => {
  if (!q.value) {
    return items.value; // Return all items if search query is empty
  }

  return items.value.filter((site:any) => {
    // Check main site properties
    const siteMatches = Object.values(site).some((value) =>
        String(value).toLowerCase().includes(q.value.toLowerCase())
    );

    // Check site_proxies properties
    const siteProxiesMatch = site.site_proxies?.some((proxy:any) =>
        Object.values(proxy).some((value) =>
            String(value).toLowerCase().includes(q.value.toLowerCase())
        )
    );

    return siteMatches || siteProxiesMatch;
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
    id: "",
    server: "http://",
    username: "lnvmpyru",
    password: "5un1tb1azapa",
    status: "active",
  }
}
function handleAdd(){
  isNewModalOpen.value = !isNewModalOpen.value
  resetItem()
}
async function syncProxyList() {
  loading.value = true
  useSiteApi().syncProxy().then(res=>{
    if(res.status.value!="error"){
      refresh()
      toast.add({ title: "Proxy Sync Done" })
    }
    loading.value = false
  })
}
async function saveItem() {
  loading.value = true
  useSiteApi().saveProxy(proxy.value).then(res=>{
    if(res.status.value!="error"){
      isNewModalOpen.value = false;
      refresh()
      toast.add({ title: "Proxy Saved" })
    }
    loading.value = false
  })
}
async function updateItem() {
  loading.value = true
  useSiteApi().updateProxy(proxy.value,proxy.value.id).then(res=>{
    if(res.status.value!="error"){
      isEditModalOpen.value = false;
      refresh()
      toast.add({ title: "Proxy Update" })
    }
    loading.value = false
  })
}
async function deleteCrawler(proxy:Proxy) {
  loading.value = true
  await useSiteApi().removeProxy(proxy.id).then(res=>{
    if(res.status.value!="error"){
      loading.value = false
      toast.add({ title: "Proxy Deleted" })
      refresh()
    }
  })
}
function handleEdit(row:any) {
  if (row.id){
    proxy.value = row
  }
  isEditModalOpen.value = !isEditModalOpen.value;
}
function handleLog(row:any) {
  proxyLog.value = row.error_log
  isLogModalOpen.value = !isLogModalOpen.value;
}
</script>
