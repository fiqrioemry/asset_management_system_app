<script lang="ts">
	import {
		assets,
		assetStore,
		assetsError,
		assetsLoading,
		assetsFilters,
		assetsPagination
	} from '$lib/stores/asset.store';
	import { onMount } from 'svelte';
	import * as Table from '$lib/components/ui/table';
	import { formatPrice } from '$lib/utils/formatter';
	import { Edit, Eye, Trash2 } from '@lucide/svelte';
	import AlertCard from '$lib/components/ui/AlertCard.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Pagination from '$lib/components/ui/Pagination.svelte';
	import SearchInput from '$lib/components/ui/SearchInput.svelte';
	import FilterSelect from '$lib/components/ui/FilterSelect.svelte';
	import { locations, locationStore } from '$lib/stores/location.store';
	import { categories, categoryStore } from '$lib/stores/category.store';
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import TableLoading from '$lib/components/loading/TableLoading.svelte';
	import FormDialog from '$lib/components/common/FormDialog.svelte';
	import UpdateAssetForm from '$lib/components/business/UpdateAssetForm.svelte';
	import UpdateAsset from '$lib/components/business/UpdateAsset.svelte';

	// Local state for inputs
	let searchValue = '';
	let selectedLocationId = '';
	let selectedCategoryId = '';

	onMount(async () => {
		await categoryStore.fetchCategories();
		await locationStore.fetchLocations();
		await assetStore.fetchAssets($assetsFilters);

		// Sync initial values dari store
		searchValue = $assetsFilters.search || '';
		selectedLocationId = $assetsFilters.locationId || '';
		selectedCategoryId = $assetsFilters.categoryId || '';
	});

	// handle search input with debounce
	async function handleSearch(searchTerm: string) {
		assetStore.setFilters({ search: searchTerm, page: 1 });
		await assetStore.fetchAssets($assetsFilters);
	}

	async function handleLocationFilter(locationId: string) {
		assetStore.setFilters({ locationId, page: 1 });
		await assetStore.fetchAssets($assetsFilters);
	}

	async function handleCategoryFilter(categoryId: string) {
		assetStore.setFilters({ categoryId, page: 1 });
		await assetStore.fetchAssets($assetsFilters);
	}

	// handle pagination
	async function handlePageChange(page: number) {
		assetStore.setFilters({ page });
		await assetStore.fetchAssets($assetsFilters);
	}

	$: currentPage = $assetsFilters.page || 1;
	$: perPage = $assetsPagination?.limit || 10;
	$: total = $assetsPagination?.totalItems || 0;
	$: totalPages = $assetsPagination?.totalPages || 1;

	// Handle error retry
	async function handleErrorRetry() {
		await assetStore.fetchAssets($assetsFilters);
	}

	export const assetConditionConfig = {
		new: { bg: 'bg-red-100', text: 'text-red-800' },
		Poor: { bg: 'bg-red-100', text: 'text-red-800' },
		Good: { bg: 'bg-blue-100', text: 'text-blue-800' },
		Fair: { bg: 'bg-yellow-100', text: 'text-yellow-800' },
		Excellent: { bg: 'bg-green-100', text: 'text-green-800' }
	};
</script>

<div class="assets-page">
	<div class="mx-auto w-full max-w-7xl p-6">
		<!-- Error Message -->
		{#if $assetsError}
			<AlertCard onAction={handleErrorRetry} message={$assetsError} />
		{/if}

		<div class="rounded-lg border border-gray-200 bg-white shadow-sm">
			<div class="border-b px-6 py-4">
				<h2 class="text-xl font-semibold text-gray-900">Asset Management</h2>
				<p class="mt-1 text-sm text-gray-500">Manage your company assets and categories</p>
			</div>
			<div
				class="flex flex-col items-center justify-between gap-2 border-b border-gray-200 px-4 py-4 md:flex-row"
			>
				<!-- Search -->
				<SearchInput
					size="sm"
					className="md:w-96 w-full"
					onSearch={handleSearch}
					placeholder="Search assets..."
					value={$assetsFilters.search || ''}
				/>

				<!-- Filters & Sort -->
				<div class="flex w-full items-center justify-end gap-2 md:w-auto">
					<!-- Location Filter -->
					<FilterSelect
						size="sm"
						clearable={true}
						options={$locations}
						placeholder="All Locations"
						onChange={handleLocationFilter}
						value={$assetsFilters.locationId || ''}
					/>

					<!-- Category Filter -->
					<FilterSelect
						size="sm"
						clearable={true}
						options={$categories}
						placeholder="All Categories"
						onChange={handleCategoryFilter}
						value={$assetsFilters.categoryId || ''}
					/>
				</div>
			</div>

			<div class="overflow-x-auto">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Name</Table.Head>
							<Table.Head>Price</Table.Head>
							<Table.Head>Location</Table.Head>
							<Table.Head>Condition</Table.Head>
							<Table.Head>Category</Table.Head>
							<Table.Head>Action</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#if $assetsLoading}
							<TableLoading rows={5} />
						{:else if $assets && $assets.length > 0}
							{#each $assets as asset}
								<Table.Row>
									<Table.Cell>
										<div class="flex gap-2">
											<img src={asset.image} alt={asset.name} class="h-10 w-10 rounded-md" />
											<div>
												<div class="text-sm font-medium">{asset.name}</div>
												<div class="text-muted-foreground truncate text-sm">
													{asset.description}
												</div>
											</div>
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="primary">{asset.location?.name}</Badge>
									</Table.Cell>
									<Table.Cell class="text-sm">
										{formatPrice(asset.price)}
									</Table.Cell>
									<Table.Cell>
										<StatusBadge
											showIcon={true}
											value={asset.condition}
											customConfig={assetConditionConfig}
										/>
									</Table.Cell>
									<Table.Cell>
										{asset?.category?.name || 'N/A'}
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center space-x-2">
											<Button size="icon" variant="outline" class="rounded-full">
												<Eye class="h-3 w-3" />
											</Button>
											<UpdateAsset {asset} />
											<Button size="icon" variant="outline" class="rounded-full bg-red-100">
												<Trash2 class="h-4 w-4 text-red-500" />
											</Button>
										</div>
									</Table.Cell>
								</Table.Row>
							{/each}
						{:else}
							<Table.Row>
								<Table.Cell colspan={6} class="text-muted-foreground text-center text-sm">
									No assets found.
								</Table.Cell>
							</Table.Row>
						{/if}
					</Table.Body>
				</Table.Root>
				<Pagination {currentPage} {totalPages} {total} {perPage} onPageChange={handlePageChange} />
			</div>
		</div>
	</div>
</div>
