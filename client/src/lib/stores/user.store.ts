// src/lib/stores/user.store.ts
import { toast } from 'svelte-sonner';
import { writable, derived } from 'svelte/store';
import { authStore } from '$lib/stores/auth.store';
import { userService } from '$lib/services/user.service';
import type {
	UserProfileResponse,
	UpdateUserRequest,
	ChangePasswordRequest
} from '$lib/types/global.types';

interface UserState {
	isLoading: boolean;
	error: string | null;
	user: UserProfileResponse | null;
}

const initialState: UserState = {
	isLoading: false,
	error: null,
	user: null
};

function createUserStore() {
	const { subscribe, set, update } = writable<UserState>(initialState);

	return {
		subscribe,

		/**
		 * Get current user profile
		 */
		async getUser() {
			this.setLoading(true);
			this.clearError();

			try {
				const response: any = await userService.getMe();

				if (response.success) {
					this.setUser(response.user);
					// Sync dengan auth store
					authStore.setUser(response.user);
					return { success: true, user: response.user };
				} else {
					this.setError(response.message || 'Failed to get user profile');
					return { success: false, error: response.message };
				}
			} catch (error: any) {
				const errorMessage = error.response?.data?.message || 'Failed to get user profile';
				this.setError(errorMessage);
				return { success: false, error: errorMessage };
			} finally {
				this.setLoading(false);
			}
		},

		/**
		 * Update user profile
		 */
		async updateUser(userData: UpdateUserRequest) {
			this.setLoading(true);
			this.clearError();

			try {
				const response: any = await userService.updateUser(userData);

				if (response.success) {
					this.setUser(response.user);
					// Sync dengan auth store
					authStore.setUser(response.user);

					// Success toast
					toast.success('Profile updated successfully', {
						description: response.message || 'Your profile has been updated.'
					});

					return { success: true, user: response.user };
				} else {
					this.setError(response.message || 'Failed to update profile');
					return { success: false, error: response.message };
				}
			} catch (error: any) {
				const errorMessage = error.response?.data?.message || 'Failed to update profile';
				this.setError(errorMessage);
				return { success: false, error: errorMessage };
			} finally {
				this.setLoading(false);
			}
		},

		/**
		 * Change user password
		 */
		async changePassword(passwordData: ChangePasswordRequest) {
			this.setLoading(true);
			this.clearError();

			try {
				const response: any = await userService.changePassword(passwordData);

				toast.success(response.message || 'Password changed successfully');

				if (response.success) {
					return { success: true };
				} else {
					this.setError(response.message || 'Failed to change password');
					return { success: false, error: response.message };
				}
			} catch (error: any) {
				const errorMessage = error.response?.data?.message || 'Failed to change password';
				this.setError(errorMessage);
				return { success: false, error: errorMessage };
			} finally {
				this.setLoading(false);
			}
		},

		/**
		 * Refresh user profile
		 */
		async refreshUser() {
			return this.getUser();
		},

		// Helper methods
		setUser(user: UserProfileResponse | null) {
			update((state) => ({
				...state,
				user,
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

export const userStore = createUserStore();

// Derived stores
export const userError = derived(userStore, ($userStore) => $userStore.error);
export const currentUserProfile = derived(userStore, ($userStore) => $userStore.user);
export const isUserLoading = derived(userStore, ($userStore) => $userStore.isLoading);

// Convenience exports
export const resetUser = () => userStore.reset();
export const clearUserError = () => userStore.clearError();
