// src/lib/stores/addDialog.store.ts
import { writable } from 'svelte/store';

type AddDialogState = {
	isOpen: boolean;
	isLoading: boolean;
	title: string;
	submitHandler?: (data: any) => Promise<void>;
};

const initialState: AddDialogState = {
	isOpen: false,
	isLoading: false,
	title: ''
};

const { subscribe, set, update } = writable<AddDialogState>(initialState);

export const addDialogStore = {
	subscribe,

	/**
	 * Open add dialog
	 */
	open(title: string, submitHandler: (data: any) => Promise<void>) {
		set({
			isOpen: true,
			isLoading: false,
			title,
			submitHandler
		});
	},

	/**
	 * Set loading state
	 */
	setLoading(loading: boolean) {
		update((state) => ({ ...state, isLoading: loading }));
	},

	/**
	 * Close dialog
	 */
	close() {
		set(initialState);
	},

	/**
	 * Submit form data
	 */
	async submit(data: any) {
		update((state) => ({ ...state, isLoading: true }));

		try {
			const currentState = get(addDialogStore);
			if (currentState.submitHandler) {
				await currentState.submitHandler(data);
			}
			// Dialog will be closed by the component after success
		} catch (error) {
			update((state) => ({ ...state, isLoading: false }));
			throw error;
		}
	}
};

// Helper to get current state
function get(store: typeof addDialogStore) {
	let state: AddDialogState;
	store.subscribe((s) => (state = s))();
	return state!;
}
