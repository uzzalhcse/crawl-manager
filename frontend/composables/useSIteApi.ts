import type { SearchParameters } from '~/types';

export function useSiteApi() {
  const findAll = (query: SearchParameters) => useApi('/api/site', { query });
  const save = (payload: any) => useApi('/api/site', { method:'POST',body:JSON.stringify(payload) });
  const update = (payload: any, id:string) => useApi(`/api/site/${id}`, { method:'PUT',body:JSON.stringify(payload) });
  const remove = (id:number) => useApi(`/api/site/${id}`, { method:'DELETE' });

  const getSecrets = (site_id:string) => useApi(`/api/site-secret/${site_id}`);
  const addSecrets = (payload: any) => useApi(`/api/site-secret`,{ method:'POST',body:JSON.stringify(payload) })
  const updateSecrets = (payload: any,site_id:string) => useApi(`/api/site-secret/${site_id}`,{ method:'PUT',body:JSON.stringify(payload) })
  const startCrawler = (site_id:string) => useApi(`/api/start-crawler/${site_id}`);
  const stopCrawler = (instance_name:string) => useApi(`/api/stop-crawler/${instance_name}`);
  const buildCrawler = (site_id:string) => useApi(`/api/build-crawler/${site_id}`);
  const crawlingHistory = () => useApi(`/api/crawling-history`);
  const proxyList = () => useApi(`/api/proxy`);
  const saveProxy = (payload: any) => useApi(`/api/proxy`,{ method:'POST',body:JSON.stringify(payload) });
  const removeProxy = (server:string) => useApi(`/api/proxy/${server}`, { method:'DELETE' });
  const updateProxy = (payload: any, id:string) => useApi(`/api/proxy/${id}`, { method:'PUT',body:JSON.stringify(payload) });

  return {
    findAll,
    save,
    update,
    remove,
    getSecrets,
    addSecrets,
    updateSecrets,
    startCrawler,
    crawlingHistory,
    stopCrawler,
    buildCrawler,
    proxyList,
    saveProxy,
    updateProxy,
    removeProxy,
  };
}
