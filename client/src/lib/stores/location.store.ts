// lib/stores/location.store.ts
import { writable, derived } from 'svelte/store';
import { locationService } from '$lib/services/location.service';
import type { Location, LocationRequest, UpdateLocationRequest } from '$lib/types/global.types';

interface LocationState {
	locations: Location[];
	loading: boolean;
	error: string | null;
}

const initialState: LocationState = {
	locations: [],
	loading: false,
	error: null
};

function createLocationStore() {
	const { subscribe, set, update } = writable<LocationState>(initialState);

	return {
		subscribe,

		async fetchLocations() {
			this.setLoading(true);
			try {
				const response: any = await locationService.getLocations();
				if (response.success) {
					this.setLocations(response.data.locations || []);
				}
			} catch (error: any) {
				console.error('Error fetching locations:', error);
				this.setError(error.response?.data?.message || 'Failed to fetch locations');
			} finally {
				this.setLoading(false);
			}
		},

		async createLocation(locationRequest: LocationRequest) {
			this.setLoading(true);
			try {
				const response: any = await locationService.createLocation(locationRequest);

				if (response.success) {
					// Add new location to the list using server response data
					update((state) => ({
						...state,
						locations: [...state.locations, response.data].sort(
							(a, b) => a.name.localeCompare(b.name) // Sort alphabetically
						),
						loading: false,
						error: null
					}));

					return response.data; // Return created location for caller
				} else {
					this.setError(response.message || 'Failed to create location');
					throw new Error(response.message || 'Failed to create location');
				}
			} catch (error: any) {
				console.error('Error creating location:', error);
				this.setError(error.response?.data?.message || 'Failed to create location');
				throw error; // Re-throw for caller to handle
			} finally {
				this.setLoading(false);
			}
		},

		async updateLocation(updateLocationRequest: UpdateLocationRequest) {
			this.setLoading(true);
			try {
				const response: any = await locationService.updateLocation(updateLocationRequest);

				if (response.success) {
					// Update specific location with server response data
					update((state) => ({
						...state,
						locations: state.locations
							.map((location) =>
								location.id === updateLocationRequest.id
									? { ...location, ...response.data } // Use updated data from server
									: location
							)
							.sort((a, b) => a.name.localeCompare(b.name)), // Re-sort after update
						loading: false,
						error: null
					}));

					return response.data; // Return updated location for caller
				} else {
					this.setError(response.message || 'Failed to update location');
					throw new Error(response.message || 'Failed to update location');
				}
			} catch (error: any) {
				console.error('Error updating location:', error);
				this.setError(error.response?.data?.message || 'Failed to update location');
				throw error; // Re-throw for caller to handle
			} finally {
				this.setLoading(false);
			}
		},

		async deleteLocation(locationId: string) {
			if (!locationId) {
				this.setError('Location ID is required');
				return false;
			}

			this.setLoading(true);
			try {
				const response: any = await locationService.deleteLocation(locationId);

				if (response.success) {
					// Remove location from list
					update((state) => ({
						...state,
						locations: state.locations.filter((location) => location.id !== locationId),
						loading: false,
						error: null
					}));

					return true; // Return success indicator
				} else {
					this.setError(response.message || 'Failed to delete location');
					return false;
				}
			} catch (error: any) {
				console.error('Error deleting location:', error);
				this.setError(error.response?.data?.message || 'Failed to delete location');
				throw error; // Re-throw for caller to handle
			} finally {
				this.setLoading(false);
			}
		},

		// Optimistic delete for better UX
		async deleteLocationOptimistic(locationId: string) {
			if (!locationId) {
				this.setError('Location ID is required');
				return false;
			}

			// Store deleted location for potential rollback
			let deletedLocation: Location | null = null;

			// Optimistic update - remove from UI immediately
			update((state) => {
				deletedLocation = state.locations.find((location) => location.id === locationId) || null;
				return {
					...state,
					locations: state.locations.filter((location) => location.id !== locationId)
				};
			});

			try {
				const response: any = await locationService.deleteLocation(locationId);

				if (!response.success) {
					// Rollback if server says failed
					this.rollbackDelete(deletedLocation);
					this.setError(response.message || 'Failed to delete location');
					return false;
				}

				return true;
			} catch (error: any) {
				console.error('Error deleting location:', error);

				// Rollback on error
				this.rollbackDelete(deletedLocation);
				this.setError(error.response?.data?.message || 'Failed to delete location');
				return false;
			}
		},

		// Helper method for rollback
		rollbackDelete(deletedLocation: Location | null) {
			if (deletedLocation) {
				update((state) => ({
					...state,
					locations: [...state.locations, deletedLocation].sort(
						(a, b) => a.name.localeCompare(b.name) // Sort alphabetically
					)
				}));
			}
		},

		searchLocations(searchTerm: string): Location[] {
			let filteredLocations: Location[] = [];

			const { subscribe } = this;
			const unsubscribe = subscribe((state) => {
				filteredLocations = state.locations.filter((loc) =>
					loc.name.toLowerCase().includes(searchTerm.toLowerCase())
				);
			});
			unsubscribe();

			return filteredLocations;
		},

		setLocations(locations: Location[]) {
			update((state) => ({
				...state,
				locations: locations.sort((a, b) => a.name.localeCompare(b.name)), // Always sort alphabetically
				loading: false,
				error: null
			}));
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

export const locationStore = createLocationStore();

// Derived stores
export const locations = derived(locationStore, ($state) => $state.locations);
export const locationsError = derived(locationStore, ($state) => $state.error);
export const locationsLoading = derived(locationStore, ($state) => $state.loading);
export const locationOptions = derived(locations, ($locations) =>
	$locations.map((location) => ({
		value: location.id,
		label: location.name
	}))
);

// Derived store for location count
export const locationsCount = derived(locations, ($locations) => $locations.length);

// Helper functions
export const clearLocationsError = () => {
	locationStore.clearError();
};

export const searchLocationsByName = (searchTerm: string) => {
	return locationStore.searchLocations(searchTerm);
};
