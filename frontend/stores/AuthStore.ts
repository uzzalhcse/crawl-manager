import { defineStore } from 'pinia';
import { useAuthApi } from '~/composables/useAuthApi'; // Adjust the import path as needed

interface User {
  name: string | null,
  email: string | null,
  roles: string[];
  id?: string;
}

interface AuthStoreState {
  isAuthenticated: boolean;
  accessToken: string | null | undefined;
  user: User | null;
  isAdmin: boolean;
  loginError: string | null;
}

export const useAuthStore = defineStore('AuthStore', {
  state: (): AuthStoreState => ({
    isAuthenticated: !!useCookie('access_token').value,
    accessToken: useCookie('access_token').value,
    user: {
      name: null,
      email: null,
      roles: [],
    },
    isAdmin: false,
    loginError: null,
  }),

  actions: {
    setToken(token: string): Promise<boolean> {
      return new Promise((resolve) => {
        this.isAuthenticated = true;
        this.accessToken = token;
        const tokenCookie = useCookie('access_token', {
          maxAge: 60 * 60 * 24 * 7, // 7 days
          httpOnly: false,
          secure: true,
        });
        tokenCookie.value = token;
        resolve(true);
      });
    },

    removeToken(): Promise<void> {
      return new Promise((resolve) => {
        this.resetToken();
        resolve();
      });
    },

    resetToken(): void {
      this.isAuthenticated = false;
      this.accessToken = null;
      this.user = null;
      const tokenCookie = useCookie('access_token');
      tokenCookie.value = null;
    },

    async login(credentials: { username: string, password: string }): Promise<boolean> {
      this.loginError = null;
      try {
        const { data, error } = await useAuthApi().login({
          username: credentials.username,
          password: credentials.password
        });

        if (error.value) {
          this.loginError = error.value?.message || 'Login failed';
          return false;
        }

        if (data.value && data.value.token) {
          await this.setToken(data.value.token);
          await this.fetchUserProfile();

          // Redirect to dashboard or home
          await navigateTo('/portal');
          return true;
        }

        return false;
      } catch (err) {
        this.loginError = 'An unexpected error occurred';
        console.error('Login error:', err);
        return false;
      }
    },

    async register(userData: {
      name: string,
      email: string,
      username: string,
      password: string
    }): Promise<boolean> {
      this.loginError = null;
      try {
        const { data, error } = await useAuthApi().register({
          name: userData.name,
          email: userData.email,
          username: userData.username,
          password: userData.password
        });

        if (error.value) {
          this.loginError = error.value?.message || 'Registration failed';
          return false;
        }

        await navigateTo('/portal');
        return !!data.value;

      } catch (err) {
        this.loginError = 'An unexpected error occurred during registration';
        console.error('Registration error:', err);
        return false;
      }
    },


    async fetchUserProfile(): Promise<boolean> {
      if (!this.accessToken) return false;

      try {
        const { data, error } = await useFetch<{ user?: { id: string; name: string; email: string; roles?: string[] } }>('/api/auth/profile', {
          headers: {
            'Authorization': `Bearer ${this.accessToken}`
          }
        });

        if (error.value) {
          this.resetToken();
          return false;
        }

        if (data.value?.user) {
          this.user = {
            id: data.value.user.id,
            name: data.value.user.name,
            email: data.value.user.email,
            roles: data.value.user.roles || []
          };

          // Check if user is admin
          this.isAdmin = this.user.roles.includes('ADMIN');

          return true;
        }

        return false;
      } catch (err) {
        console.error('Fetch profile error:', err);
        return false;
      }
    },

    async logout(): Promise<void> {
      try {
        // Optional: Call logout endpoint to invalidate token on server
        await useFetch('/api/auth/logout', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${this.accessToken}`
          }
        });
      } catch (err) {
        console.error('Logout error:', err);
      } finally {
        this.resetToken();
        await navigateTo('/login');
      }
    },

  },
});

// Optional: Add type guards and helpers
export function isAdmin(store: ReturnType<typeof useAuthStore>): boolean {
  return store.isAdmin;
}

export function hasRole(store: ReturnType<typeof useAuthStore>, role: string): boolean {
  return store.user?.roles.includes(role) || false;
}