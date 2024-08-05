import { defineStore } from 'pinia';

interface User {
  name: string | null,
  email: string | null,
  roles: string[];
}

interface AuthStoreState {
  isAuthenticated: boolean;
  accessToken: string | null| undefined;
  user: User | null;
  isAdmin: boolean;
}

export const useAuthStore = defineStore('AuthStore', {
  state: (): AuthStoreState => ({
    isAuthenticated: !!useCookie('access_token').value, // Check if token exists in cookies
    accessToken: useCookie('access_token').value,
    user: {
      name: null,
      email: null,
      roles: [],
    },
    isAdmin: false,
  }),
  actions: {
    setToken(token: string): Promise<boolean> {
      return new Promise((resolve, reject) => {
        this.isAuthenticated = true;
        this.accessToken = token;
        useCookie('access_token').value = token;
        resolve(true);
      });
    },
    removeToken(): Promise<void> {
      return new Promise((resolve, reject) => {
        this.resetToken();
        resolve();
      });
    },
    resetToken(): void {
      this.isAuthenticated = false;
      this.accessToken = null;
      this.user = null;
      useCookie('access_token').value = null;
    },
  },
});
