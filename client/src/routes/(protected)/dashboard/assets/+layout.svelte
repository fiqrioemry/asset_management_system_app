<script lang="ts">
	// asset store imports
	import {
		assets,
		error,
		filters,
		isLoading,
		assetStore,
		pagination
	} from '$lib/stores/asset.store';
	import type { Asset } from '$lib/types/global.types';

	// other functionality
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { confirmStore } from '$lib/stores/confirm.store';
	import { browser } from '$app/environment';

	// asset component imports
	import AlertCard from '$lib/components/ui/AlertCard.svelte';
	import Pagination from '$lib/components/ui/Pagination.svelte';
	import AssetsTable from '$lib/components/assets/AssetsTable.svelte';
	import AssetFilter from '$lib/components/assets/AssetFilter.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import AssetDetail from '$lib/components/assets/AssetDetail.svelte';

	// Dialog state management
	let showDetailDialog = false;
	let selectedAssetId = '';

	console.log('Assets page loaded'); // Debug log

	// Fetch initial data
	onMount(async () => {
		await assetStore.getAll($filters);
		if (browser) {
			setTimeout(() => {
				checkForAssetId();
			}, 100);
		}
	});

	function checkForAssetId() {
		const assetId = $page.url.searchParams.get('assetId');
		console.log('Checking for assetId:', assetId); // Debug log
		console.log('Current URL:', $page.url.href); // Debug log

		if (assetId && assetId !== selectedAssetId) {
			console.log('Opening dialog for asset:', assetId); // Debug log
			openDetailDialog(assetId);

			// Clean URL after dialog opens
			setTimeout(() => {
				if (browser) {
					console.log('Cleaning URL'); // Debug log
					window.history.replaceState({}, '', '/dashboard/assets');
				}
			}, 200);
		}
	}

	// Watch for URL changes - this is crucial
	$: if (browser && $page.url) {
		const assetId = $page.url.searchParams.get('assetId');
		console.log('URL changed, assetId:', assetId, 'current selectedAssetId:', selectedAssetId); // Debug log

		if (assetId && assetId !== selectedAssetId && !showDetailDialog) {
			console.log('Reactive opening dialog for:', assetId); // Debug log
			openDetailDialog(assetId);
		}
	}

	// set filter variables
	let initialFilters = {
		...$filters
	};

	// handle search input with debounce
	async function handleSearch(searchTerm: string) {
		assetStore.setFilters({ search: searchTerm, page: 1 });
		await assetStore.getAll($filters);
	}

	// handle location filter
	async function handleLocationFilter(locationId: string) {
		assetStore.setFilters({ locationId, page: 1 });
		await assetStore.getAll($filters);
	}

	// handle category filter
	async function handleCategoryFilter(categoryId: string) {
		assetStore.setFilters({ categoryId, page: 1 });
		await assetStore.getAll($filters);
	}

	// handle pagination
	async function handlePageChange(page: number) {
		assetStore.setFilters({ page });
		await assetStore.getAll($filters);
	}

	// Handle error retry
	async function handleErrorRetry() {
		await assetStore.getAll($filters);
	}

	// open modal for asset detail
	function handleOpenDetail(asset: Asset) {
		console.log('handleOpenDetail called for asset:', asset.id); // Debug log
		goto(`/dashboard/assets/${asset.id}`, { replaceState: true });
	}

	// open detail dialog
	function openDetailDialog(assetId: string) {
		console.log('openDetailDialog called with:', assetId); // Debug log
		selectedAssetId = assetId;
		showDetailDialog = true;
		console.log(
			'Dialog state set - showDetailDialog:',
			showDetailDialog,
			'selectedAssetId:',
			selectedAssetId
		); // Debug log
	}

	// close detail dialog
	function closeDetailDialog() {
		console.log('closeDetailDialog called'); // Debug log
		showDetailDialog = false;
		selectedAssetId = '';
		// Ensure we're back to the base assets page
		goto('/dashboard/assets', { replaceState: true });
	}

	// open modal for asset update
	function handleUpdateAsset(asset: Asset) {
		goto(`/dashboard/assets/${asset.id}/edit`, { replaceState: true });
	}

	// open confirmation modal for asset deletion
	function handleDeleteAsset(id: string) {
		confirmStore
			.delete('Delete Asset', `Are you sure you want to delete this asset?`)
			.then(async (confirmed) => {
				if (confirmed) {
					await assetStore.delete(id);
				}
			});
	}

	// Debug reactive statements
	$: console.log(
		'Reactive - showDetailDialog:',
		showDetailDialog,
		'selectedAssetId:',
		selectedAssetId
	);
	$: console.log('Current URL:', $page.url.href);
</script>

<div class="assets-page">
	<div class="mx-auto w-full max-w-7xl p-6">
		<!-- Debug info (remove in production) -->
		<div class="mb-4 rounded bg-gray-100 p-2 text-sm">
			Debug: Dialog Open: {showDetailDialog}, Asset ID: {selectedAssetId}
		</div>

		<!-- Error Message -->
		{#if $error}
			<AlertCard onAction={handleErrorRetry} message={$error.message} />
		{/if}

		<div class="rounded-lg border border-gray-200 bg-white shadow-sm">
			<!-- Page Heading -->
			<PageHeading title="Assets" description="Manage your company assets and categories" />

			<!-- Asset Filter section -->
			<AssetFilter
				filters={initialFilters}
				onSearchChange={handleSearch}
				onLocationChange={handleLocationFilter}
				onCategoryChange={handleCategoryFilter}
			/>

			<div class="overflow-x-auto">
				<!-- asset table -->
				<AssetsTable
					assets={$assets}
					loading={$isLoading}
					onOpenDetail={handleOpenDetail}
					onUpdateAsset={handleUpdateAsset}
					onDeleteAsset={handleDeleteAsset}
				/>
				<!-- asset table pagination -->
				<Pagination
					page={$pagination.page}
					onPageChange={handlePageChange}
					perPage={$pagination?.limit || 10}
					total={$pagination?.totalItems || 0}
					totalPages={$pagination?.totalPages || 1}
				/>
			</div>
		</div>
	</div>
</div>

<!-- Asset Detail Dialog -->
{#if selectedAssetId}
	<AssetDetail assetId={selectedAssetId} open={showDetailDialog} onClose={closeDetailDialog} />
{/if}
