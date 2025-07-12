// lib/stores/auth.store.ts
import { writable, derived } from 'svelte/store';
import { goto } from '$app/navigation';
import { authService } from '$lib/services/auth.service';
import type { UserProfileResponse, LoginRequest, RegisterRequest } from '$lib/types/global.types';
import { toast } from 'svelte-sonner';

interface AuthState {
	isLoading: boolean;
	error: string | null;
	isAuthenticated: boolean;
	user: UserProfileResponse | null;
}

const initialState: AuthState = {
	error: null,
	user: null,
	isLoading: false,
	isAuthenticated: false
};

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	return {
		subscribe,

		/**
		 * Login user dengan email dan password
		 */
		async login(credentials: LoginRequest, redirectTo: string = '/dashboard') {
			this.setLoading(true);
			this.clearError();

			try {
				const response: any = await authService.login(credentials);

				if (response.success) {
					this.setUser(response.user);
					goto(redirectTo);
					return { success: true, user: response.user };
				} else {
					this.setError(response.message || 'Login failed. Please try again.');
					return { success: false, error: response.message };
				}
			} catch (error: any) {
				const errorMessage = error.response?.data?.message || 'Login failed. Please try again.';
				this.setError(errorMessage);
				return { success: false, error: errorMessage };
			} finally {
				this.setLoading(false);
			}
		},

		/**
		 * Register user baru
		 */
		async register(userData: RegisterRequest, redirectTo: string = '/dashboard') {
			this.setLoading(true);
			this.clearError();

			try {
				const response: any = await authService.register(userData);

				if (response.success) {
					this.setUser(response.user);
					goto(redirectTo);
					return { success: true, user: response.user };
				} else {
					this.setError(response.message || 'Registration failed. Please try again.');
					return { success: false, error: response.message };
				}
			} catch (error: any) {
				const errorMessage =
					error.response?.data?.message || 'Registration failed. Please try again.';
				this.setError(errorMessage);
				return { success: false, error: errorMessage };
			} finally {
				this.setLoading(false);
			}
		},

		/**
		 * Logout user
		 */
		async logout(redirectTo: string = '/signin') {
			this.setLoading(true);
			try {
				await authService.logout();
			} catch (error) {
				console.warn('Logout API failed, but continuing with local logout');
			} finally {
				this.reset();
				goto(redirectTo);
			}
		},

		/**
		 * Check authentication status
		 */
		async checkAuth() {
			this.setLoading(true);

			try {
				const response: any = await authService.refreshSession();

				if (response?.success && response?.user) {
					this.setUser(response.user);
					return { authenticated: true, user: response.user };
				} else {
					this.reset();
					return { authenticated: false, user: null };
				}
			} catch (error: any) {
				this.reset();
				return { authenticated: false, user: null };
			} finally {
				this.setLoading(false);
			}
		},

		/**
		 * Refresh session
		 */
		async refreshSession() {
			try {
				const response: any = await authService.refreshSession();

				if (response?.success && response?.user) {
					this.setUser(response.user);
					return { success: true, user: response.user };
				} else {
					this.reset();
					return { success: false, error: 'Session refresh failed' };
				}
			} catch (error: any) {
				this.reset();
				return {
					success: false,
					error: error.response?.data?.message || 'Session refresh failed'
				};
			}
		},

		// Helper methods
		setUser(user: UserProfileResponse | null) {
			update((state) => ({
				...state,
				user,
				isAuthenticated: !!user,
				error: null
			}));
		},

		setError(error: string) {
			update((state) => ({ ...state, error, isLoading: false }));
		},

		setLoading(isLoading: boolean) {
			update((state) => ({ ...state, isLoading }));
		},

		clearError() {
			update((state) => ({ ...state, error: null }));
		},

		reset() {
			set(initialState);
		}
	};
}

export const authStore = createAuthStore();

export const resetAuth = () => authStore.reset();
export const clearAuthError = () => authStore.clearError();
export const currentUser = derived(authStore, ($authStore) => $authStore.user);
export const authError = derived(authStore, ($authStore) => $authStore.error);
export const isLoading = derived(authStore, ($authStore) => $authStore.isLoading);
export const isAuthenticated = derived(authStore, ($authStore) => $authStore.isAuthenticated);
