<template>
  <div>
    <button type="button" class="inline-flex items-center p-2 mt-2 ml-3 text-sm text-gray-500 rounded-lg sm:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600">
      <span class="sr-only">Open sidebar</span>
      <svg class="w-6 h-6" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
        <path clip-rule="evenodd" fill-rule="evenodd" d="M2 4.75A.75.75 0 012.75 4h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 4.75zm0 10.5a.75.75 0 01.75-.75h7.5a.75.75 0 010 1.5h-7.5a.75.75 0 01-.75-.75zM2 10a.75.75 0 01.75-.75h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 10z"></path>
      </svg>
    </button>

    <aside id="logo-sidebar" class="fixed top-0 left-0 z-40 w-64 h-screen transition-transform -translate-x-full sm:translate-x-0" aria-label="Sidebar">
      <div class="h-full px-3 py-4 overflow-y-auto bg-gray-50 dark:bg-gray-800">
        <a href="#" class="flex items-center pl-2.5 mb-5">
<!--          <img src="https://flowbite.com/docs/images/logo.svg" class="h-6 mr-3 sm:h-7" alt="Flowbite Logo" />-->
          <span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">Crawl Manager</span>
        </a>
        <ul class="space-y-2 font-medium">
          <li>
            <nuxt-link to="/portal" class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group">
              <UIcon name="i-heroicons-squares-2x2" />
              <span class="ml-3">Dashboard</span>
            </nuxt-link>
          </li>
          <li>
            <nuxt-link to="/portal/sites" class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group">
              <UIcon name="i-heroicons-globe-alt" />
              <span class="ml-3">Sites</span>
            </nuxt-link>
          </li>
          <li>
            <nuxt-link to="/portal/crawlinghistory" class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group">
              <UIcon name="i-heroicons-cpu-chip" />
              <span class="ml-3">Crawling History</span>
            </nuxt-link>
          </li>
<!--          <li>-->
<!--            <nuxt-link to="/portal/logs" class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group">-->
<!--              <UIcon name="i-heroicons-archive-box" />-->
<!--              <span class="ml-3">Logs</span>-->
<!--            </nuxt-link>-->
<!--          </li>-->
          <span class="text-gray-400 text-[12px]">Settings</span>
          <li>
            <nuxt-link to="/portal/setting/general" class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group">
              <UIcon name="i-heroicons-paint-brush" />
              <span class="ml-3">General</span>
            </nuxt-link>
          </li>
        </ul>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
</script>
