export function useAuthApi() {
  const login = (credentials: { email: string; password: string }) => useApi('/auth/login', { method: 'POST', body: JSON.stringify(credentials) });

  return {
    login
  };
}
