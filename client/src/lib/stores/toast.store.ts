// src/lib/stores/toast.store.ts
import { writable } from 'svelte/store';

export interface Toast {
	id: string;
	type: 'success' | 'error' | 'warning' | 'info';
	title?: string;
	message: string;
	duration?: number;
	dismissible?: boolean;
	position?:
		| 'top-right'
		| 'top-left'
		| 'bottom-right'
		| 'bottom-left'
		| 'top-center'
		| 'bottom-center';
}

interface ToastState {
	toasts: Toast[];
}

const initialState: ToastState = {
	toasts: []
};

function createToastStore() {
	const { subscribe, update } = writable<ToastState>(initialState);

	return {
		subscribe,

		// Add a new toast
		addToast: (toast: Omit<Toast, 'id'>) => {
			const id = Math.random().toString(36).substr(2, 9);
			const newToast: Toast = {
				id,
				duration: 4000,
				dismissible: true,
				position: 'top-right',
				...toast
			};

			update((state) => ({
				...state,
				toasts: [...state.toasts, newToast]
			}));

			return id;
		},

		// Remove a toast by id
		removeToast: (id: string) => {
			update((state) => ({
				...state,
				toasts: state.toasts.filter((toast) => toast.id !== id)
			}));
		},

		// Clear all toasts
		clearToasts: () => {
			update((state) => ({
				...state,
				toasts: []
			}));
		},

		// Convenience methods for different toast types
		success: (message: string, options?: Partial<Omit<Toast, 'id' | 'type' | 'message'>>) => {
			const id = Math.random().toString(36).substr(2, 9);
			const toast: Toast = {
				id,
				type: 'success',
				message,
				duration: 4000,
				dismissible: true,
				position: 'top-right',
				...options
			};

			update((state) => ({
				...state,
				toasts: [...state.toasts, toast]
			}));

			return id;
		},

		error: (message: string, options?: Partial<Omit<Toast, 'id' | 'type' | 'message'>>) => {
			const id = Math.random().toString(36).substr(2, 9);
			const toast: Toast = {
				id,
				type: 'error',
				message,
				duration: 6000, // Error toasts stay longer
				dismissible: true,
				position: 'top-right',
				...options
			};

			update((state) => ({
				...state,
				toasts: [...state.toasts, toast]
			}));

			return id;
		},

		warning: (message: string, options?: Partial<Omit<Toast, 'id' | 'type' | 'message'>>) => {
			const id = Math.random().toString(36).substr(2, 9);
			const toast: Toast = {
				id,
				type: 'warning',
				message,
				duration: 5000,
				dismissible: true,
				position: 'top-right',
				...options
			};

			update((state) => ({
				...state,
				toasts: [...state.toasts, toast]
			}));

			return id;
		},

		info: (message: string, options?: Partial<Omit<Toast, 'id' | 'type' | 'message'>>) => {
			const id = Math.random().toString(36).substr(2, 9);
			const toast: Toast = {
				id,
				type: 'info',
				message,
				duration: 4000,
				dismissible: true,
				position: 'top-right',
				...options
			};

			update((state) => ({
				...state,
				toasts: [...state.toasts, toast]
			}));

			return id;
		}
	};
}

export const toastStore = createToastStore();
