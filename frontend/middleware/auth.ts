export default defineNuxtRouteMiddleware(async (to) => {
    const authStore = useAuthStore();
    const token = useCookie('access_token');

    // Define public routes
    const publicRoutes = ['/login', '/register'];

    // Redirect to login if trying to access any route other than login/register without token
    if (!publicRoutes.includes(to.path) && !token.value) {
        return navigateTo('/login');
    }

    // Redirect to dashboard if logged in and trying to access login/register
    if (publicRoutes.includes(to.path) && token.value) {
        return navigateTo('/portal');
    }

    // Fetch user profile if token exists but user not loaded
    if (token.value && !authStore.user) {
        await authStore.fetchUserProfile();
    }
});