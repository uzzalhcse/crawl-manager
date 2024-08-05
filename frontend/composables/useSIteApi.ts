import type { SearchParameters } from '~/types';

export function useSiteApi() {
  const findAll = (query: SearchParameters) => useApi('/api/site', { query });
  const save = (payload: any) => useApi('/api/site', { method:'POST',body:JSON.stringify(payload) });
  const update = (payload: any, id:string) => useApi(`/api/site/${id}`, { method:'PUT',body:JSON.stringify(payload) });
  const remove = (id:number) => useApi(`/api/site/${id}`, { method:'DELETE' });

  const getSecrets = (site_id:string) => useApi(`/api/site-secret/${site_id}`);
  const addSecrets = (payload: any) => useApi(`/api/site-secret`,{ method:'POST',body:JSON.stringify(payload) })
  const startCrawler = (site_id:string) => useApi(`/api/start-crawler/${site_id}`);
  const stopCrawler = (instance_name:string) => useApi(`/api/stop-crawler/${instance_name}`);
  const crawlingHistory = () => useApi(`/api/crawling-history`);
  return {
    findAll,
    save,
    update,
    remove,
    getSecrets,
    addSecrets,
    startCrawler,
    crawlingHistory,
    stopCrawler
  };
}
