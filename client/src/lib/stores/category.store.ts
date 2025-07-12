// lib/stores/category.store.ts
import { writable, derived } from 'svelte/store';
import { categoryService } from '$lib/services/category.service';
import type { CategoryRequest, ParentCategory } from '$lib/types/global.types';

interface CategoryState {
	categories: ParentCategory[];
	loading: boolean;
	error: string | null;
}

const initialState: CategoryState = {
	categories: [],
	loading: false,
	error: null
};

function createCategoryStore() {
	const { subscribe, set, update } = writable<CategoryState>(initialState);

	return {
		subscribe,

		async fetchCategories() {
			this.setLoading(true);
			try {
				const response: any = await categoryService.getCategories();
				if (response.success) {
					this.setCategories(response.data.categories || []);
				}
			} catch (error: any) {
				console.error('Error fetching categories:', error);
				this.setError(error.response?.data?.message || 'Failed to fetch categories');
			} finally {
				this.setLoading(false);
			}
		},

		async createCategory(categoryRequest: CategoryRequest) {
			this.setLoading(true);
			try {
				const response: any = await categoryService.createCategory(categoryRequest);

				if (response.success) {
					// Add new category to the list using server response data
					update((state) => ({
						...state,
						categories: [...state.categories, response.data],
						loading: false,
						error: null
					}));

					return response.data; // Return created category for caller
				} else {
					this.setError(response.message || 'Failed to create category');
					throw new Error(response.message || 'Failed to create category');
				}
			} catch (error: any) {
				console.error('Error creating category:', error);
				this.setError(error.response?.data?.message || 'Failed to create category');
				throw error; // Re-throw for caller to handle
			} finally {
				this.setLoading(false);
			}
		},

		async updateCategory(CategoryRequest: CategoryRequest) {
			this.setLoading(true);
			try {
				const response: any = await categoryService.updateCategory(CategoryRequest);

				if (response.success) {
					// Update specific category with server response data
					update((state) => ({
						...state,
						categories: state.categories.map((category) =>
							category.id === CategoryRequest.id
								? { ...category, ...response.data } // Use updated data from server
								: category
						),
						loading: false,
						error: null
					}));

					return response.data; // Return updated category for caller
				} else {
					this.setError(response.message || 'Failed to update category');
					throw new Error(response.message || 'Failed to update category');
				}
			} catch (error: any) {
				console.error('Error updating category:', error);
				this.setError(error.response?.data?.message || 'Failed to update category');
				throw error; // Re-throw for caller to handle
			} finally {
				this.setLoading(false);
			}
		},

		async deleteCategory(categoryId: string) {
			if (!categoryId) {
				this.setError('Category ID is required');
				return false;
			}

			this.setLoading(true);
			try {
				const response: any = await categoryService.deleteCategory(categoryId);

				if (response.success) {
					// Remove category from list
					update((state) => ({
						...state,
						categories: state.categories.filter((category) => category.id !== categoryId),
						loading: false,
						error: null
					}));

					return true;
				} else {
					this.setError(response.message || 'Failed to delete category');
					return false;
				}
			} catch (error: any) {
				console.error('Error deleting category:', error);
				this.setError(error.response?.data?.message || 'Failed to delete category');
				throw error; // Re-throw for caller to handle
			} finally {
				this.setLoading(false);
			}
		},

		findCategoryById(categoryId: string): ParentCategory | undefined {
			let foundCategory: ParentCategory | undefined;

			const { subscribe } = this;
			const unsubscribe = subscribe((state) => {
				foundCategory = state.categories.find((cat) => cat.id === categoryId);
			});
			unsubscribe();

			return foundCategory;
		},

		// Helper method to get all subcategories of a parent
		getSubcategoriesByParent(parentId: string) {
			let subcategories: any[] = [];

			const { subscribe } = this;
			const unsubscribe = subscribe((state) => {
				const parent = state.categories.find((cat) => cat.id === parentId);
				subcategories = parent?.children || [];
			});
			unsubscribe();

			return subcategories;
		},

		setCategories(categories: ParentCategory[]) {
			update((state) => ({ ...state, categories, loading: false, error: null }));
		},

		setError(error: string) {
			update((state) => ({ ...state, error, loading: false }));
		},

		setLoading(loading: boolean) {
			update((state) => ({ ...state, loading, error: null }));
		},

		clearError() {
			update((state) => ({ ...state, error: null }));
		},

		reset() {
			set(initialState);
		}
	};
}

export const categoryStore = createCategoryStore();

// Derived stores
export const categories = derived(categoryStore, ($state) => $state.categories);
export const categoriesError = derived(categoryStore, ($state) => $state.error);
export const categoriesLoading = derived(categoryStore, ($state) => $state.loading);

// Derived store for flattened categories (useful for dropdowns)
export const flatCategories = derived(categories, ($categories) => {
	const flat: Array<{ id: string; name: string; parentId?: string }> = [];

	$categories.forEach((parent) => {
		flat.push({ id: parent.id, name: parent.name });

		if (parent.children) {
			parent.children.forEach((sub) => {
				flat.push({
					id: sub.id,
					name: `${parent.name} > ${sub.name}`,
					parentId: parent.id
				});
			});
		}
	});

	return flat;
});

// Helper functions
export const clearCategoriesError = () => {
	categoryStore.clearError();
};

export const getCategoryById = (categoryId: string) => {
	return categoryStore.findCategoryById(categoryId);
};

export const getSubcategoriesByParent = (parentId: string) => {
	return categoryStore.getSubcategoriesByParent(parentId);
};
