import type { UseFetchOptions } from '#app';

export const useApi = (endpoint: string, options?: UseFetchOptions<object>) => {
  const config = useRuntimeConfig();
  const route = useRoute();
  const baseUrl = config.public.apiBase;
  const url = `${baseUrl}${endpoint}`;
  const token = useCookie('access_token');


  const toast = useToast()
  return useFetch(url, {
    ...options,
    transform: (data: any) => data.data,
    async onRequest({ options }) {
      options.headers = {
        ...options.headers,
        "Authorization": "Bearer " + token.value
      };
    },
    async onResponse({ response }) {
      if (response.ok) return;

      // Redirect to login page if we get an 401(unauthorized) error
      if (route.name !== 'login' && response.status === 401) {
        navigateTo('/login');
      }
      if (!response._data.success){
        toast.add({ title: response._data.message })
        return
      }
    }
  });
};
