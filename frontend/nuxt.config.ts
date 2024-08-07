// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  // Configure environment variables
  runtimeConfig: {
    apiSecret: '', // can be overridden by NUXT_API_SECRET environment variable
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
    }
  },
  ssr:false,
  app: {
    head: {
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
      link: [
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,400;0,600;0,700;1,400;1,600;1,700&display=swap',
        },
      ],
    },
  },
  ui: {
    icons: ['heroicons'],
    safelistColors: ['primary', 'red', 'orange', 'green'],
  },
  colorMode: {
    preference: 'dark'
  },
  plugins: [
  ],
  modules: [
    '@pinia/nuxt',
    '@nuxt/ui',
    '@nuxt/image',
    '@samk-dev/nuxt-vcalendar',
    '@pinia/nuxt',
  ],
})
