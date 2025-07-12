// src/lib/stores/confirm.store.ts
import { writable } from 'svelte/store';

type ConfirmType = 'simple' | 'typed';

type ConfirmState = {
	isOpen: boolean;
	type: ConfirmType;
	title: string;
	message: string;
	confirmWord?: string; // For typed confirmation
	resolve?: (value: boolean) => void;
};

const initialState: ConfirmState = {
	isOpen: false,
	type: 'simple',
	title: '',
	message: ''
};

const { subscribe, set, update } = writable<ConfirmState>(initialState);

export const confirmStore = {
	subscribe,

	/**
	 * Simple delete confirmation (just Yes/No buttons)
	 */
	delete(title: string, message: string): Promise<boolean> {
		return new Promise<boolean>((resolve) => {
			set({
				isOpen: true,
				type: 'simple',
				title,
				message,
				resolve
			});
		});
	},

	/**
	 * Typed confirmation (user must type a word to confirm)
	 */
	deleteTyped(title: string, message: string, confirmWord: string): Promise<boolean> {
		return new Promise<boolean>((resolve) => {
			set({
				isOpen: true,
				type: 'typed',
				title,
				message,
				confirmWord,
				resolve
			});
		});
	},

	/**
	 * User confirmed action
	 */
	confirm() {
		update((state) => {
			state.resolve?.(true);
			return { ...initialState };
		});
	},

	/**
	 * User cancelled action
	 */
	cancel() {
		update((state) => {
			state.resolve?.(false);
			return { ...initialState };
		});
	}
};
