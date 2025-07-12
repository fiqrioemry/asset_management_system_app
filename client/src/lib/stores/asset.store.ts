// lib/stores/asset.store.ts

import type {
	Asset,
	GetAssetsRequest,
	Pagination,
	UpdateAssetRequest,
	CreateAssetRequest,
	AssetsListResponse
} from '$lib/types/global.types';
import { toast } from 'svelte-sonner';
import { writable, derived } from 'svelte/store';
import { assetService } from '$lib/services/asset.service';

interface AssetState {
	asset: Asset | null;
	assets: Asset[];
	isLoading: boolean;
	error: string | null;
	filters: GetAssetsRequest;
	pagination: Pagination;
}

const initialState: AssetState = {
	asset: null,
	assets: [],
	isLoading: false,
	error: null,
	filters: {
		page: 1,
		limit: 10,
		search: '',
		categoryId: undefined,
		locationId: undefined,
		condition: undefined,
		sortBy: 'createdAt',
		sortOrder: 'desc',
		minPrice: undefined,
		maxPrice: undefined
	},
	pagination: { currentPage: 1, limit: 10, totalItems: 0, totalPages: 0 }
};

function createAssetStore() {
	const { subscribe, set, update } = writable<AssetState>(initialState);

	return {
		subscribe,

		async fetchAssets(params: GetAssetsRequest = {}) {
			this.setLoading(true);

			try {
				const response: any = await assetService.getAssets(params);
				if (response.success) {
					this.setAssets(response.data);
				}
			} catch (error: any) {
				console.error('Error fetching assets:', error);
				this.setError(error.response?.data?.message || 'Failed to fetch assets');
			} finally {
				this.setLoading(false);
			}
		},

		async getAssetById(assetId: string) {
			this.setLoading(true);

			try {
				const response: any = await assetService.getAssetById(assetId);
				if (response.success) {
					this.setAsset(response.data);
					return true;
				}
			} catch (error: any) {
				console.error('Error fetching asset:', error);
				this.setError(error.response?.data?.message || 'Failed to fetch asset');
				return false;
			} finally {
				this.setLoading(false);
			}
		},

		async createAsset(createAssetRequest: CreateAssetRequest) {
			this.setLoading(true);
			try {
				const response: any = await assetService.createNewAsset(createAssetRequest);

				if (response.success) {
					update((state) => ({
						...state,
						assets: [response.data, ...state.assets],
						pagination: {
							...state.pagination,
							totalItems: state.pagination.totalItems + 1
						},
						isLoading: false,
						error: null
					}));

					return response.data; // Return created asset for caller
				}
			} catch (error: any) {
				console.error('Error creating asset:', error);
				this.setError(error.response?.data?.message || 'Failed to create asset');
				throw error; // Re-throw for caller to handle
			} finally {
				this.setLoading(false);
			}
		},

		async updateAsset(updateAssetRequest: UpdateAssetRequest) {
			this.setLoading(true);
			try {
				const response: any = await assetService.updateAsset(updateAssetRequest);

				if (response.success) {
					// Update specific asset with server response data
					update((state) => ({
						...state,
						assets: state.assets.map((asset) =>
							asset.id === updateAssetRequest.id ? { ...asset, ...response.data } : asset
						),
						isLoading: false,
						error: null
					}));

					return response.data;
				}
			} catch (error: any) {
				console.error('Error updating asset:', error);
				this.setError(error.response?.data?.message || 'Failed to update asset');
				throw error;
			} finally {
				this.setLoading(false);
			}
		},

		async deleteAsset(assetId: string) {
			if (!assetId) {
				this.setError('Asset ID is required');
				return;
			}

			this.setLoading(true);
			try {
				const response: any = await assetService.deleteAsset(assetId);

				if (response.success) {
					toast.success(response.message);
					update((state) => ({
						...state,
						assets: state.assets.filter((asset) => asset.id !== assetId),
						pagination: {
							...state.pagination,
							totalItems: state.pagination.totalItems - 1
						},
						isLoading: false,
						error: null
					}));

					return true;
				}
			} catch (error: any) {
				console.error('Error deleting asset:', error);
				this.setError(error.response?.data?.message || 'Failed to delete asset');
				throw error;
			} finally {
				this.setLoading(false);
			}
		},

		setFilters(filters: Partial<GetAssetsRequest>) {
			update((state) => {
				const newFilters = { ...state.filters, ...filters };
				return {
					...state,
					filters: newFilters,
					pagination: { ...state.pagination, page: 1 }
				};
			});
		},

		async findAsset(assetId: string): Promise<Asset | undefined> {
			let foundAsset: Asset | undefined;
			this.setLoading(true);
			update((state) => {
				foundAsset = state.assets.find((asset) => asset.id === assetId);
				return state;
			});
			this.setAsset(foundAsset || null);
			this.setLoading(false);
			return foundAsset;
		},

		setAsset(asset: Asset | null) {
			update((state) => ({
				...state,
				asset
			}));
		},

		setAssets(data: AssetsListResponse) {
			update((state) => ({
				...state,
				assets: data.assets,
				pagination: data.pagination,
				isLoading: false,
				error: null
			}));
		},

		setError(error: string) {
			update((state) => ({ ...state, error, isLoading: false }));
		},

		setLoading(isLoading: boolean) {
			update((state) => ({ ...state, isLoading, error: null }));
		},

		clearError() {
			update((state) => ({ ...state, error: null }));
		},

		reset() {
			set(initialState);
		}
	};
}

export const assetStore = createAssetStore();
export const asset = derived(assetStore, ($state) => $state.asset);
export const assets = derived(assetStore, ($state) => $state.assets);
export const assetsError = derived(assetStore, ($state) => $state.error);
export const assetsFilters = derived(assetStore, ($state) => $state.filters);
export const assetLoading = derived(assetStore, ($state) => $state.isLoading);
export const assetsLoading = derived(assetStore, ($state) => $state.isLoading);
export const assetsPagination = derived(assetStore, ($state) => $state.pagination);

export const createAsset = async (form: CreateAssetRequest) => {
	return assetStore.createAsset(form);
};

export const updateAsset = async (asset: UpdateAssetRequest) => {
	return assetStore.updateAsset(asset);
};

export const deleteAsset = async (assetId: string) => {
	return assetStore.deleteAsset(assetId);
};

export const fetchAssets = async (params: GetAssetsRequest = {}) => {
	return assetStore.fetchAssets(params);
};

export const clearAssetsError = () => {
	assetStore.clearError();
};
