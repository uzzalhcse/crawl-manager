export function useAuthApi() {
  const login = (credentials: { username: string; password: string }) => useApi('/api/auth/login', { method: 'POST', body: JSON.stringify(credentials) });
  const register = (payload: any) => useApi('/api/auth/register', { method:'POST',body:JSON.stringify(payload) });

  return {
    login,
    register
  };
}
