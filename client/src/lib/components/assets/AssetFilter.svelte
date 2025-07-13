<script lang="ts">
	import { onMount } from 'svelte';
	import SearchInput from '$lib/components/ui/SearchInput.svelte';
	import FilterSelect from '$lib/components/ui/FilterSelect.svelte';
	import { locations, locationStore } from '$lib/stores/location.store';
	import { categories, categoryStore } from '$lib/stores/category.store';

	export let filters;

	onMount(async () => {
		await categoryStore.fetchCategories();
		await locationStore.fetchLocations();
	});

	export let onSearchChange: (searchTerm: string) => void = () => {};
	export let onLocationChange: (locationId: string) => void = () => {};
	export let onCategoryChange: (categoryId: string) => void = () => {};
</script>

<div
	class="flex flex-col items-center justify-between gap-2 border-b border-gray-200 px-4 py-4 md:flex-row"
>
	<!-- Search filter -->
	<SearchInput
		size="sm"
		className="md:w-96 w-full"
		onSearch={onSearchChange}
		placeholder="Search assets..."
		value={filters.search || ''}
	/>

	<!-- Filters & Sort -->
	<div class="flex w-full items-center justify-end gap-2 md:w-auto">
		<!-- Location Filter -->
		<FilterSelect
			size="sm"
			clearable={true}
			options={$locations}
			placeholder="All Locations"
			onChange={onLocationChange}
			value={filters.locationId || ''}
		/>
		<!-- Category Filter -->
		<FilterSelect
			size="sm"
			clearable={true}
			options={$categories}
			placeholder="All Categories"
			onChange={onCategoryChange}
			value={filters.categoryId || ''}
		/>
	</div>
</div>
