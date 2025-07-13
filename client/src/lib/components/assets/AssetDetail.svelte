<script lang="ts">
	import { onMount } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { formatDate, formatPrice } from '$lib/utils/formatter';
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import DialogError from '$lib/components/shared/DialogError.svelte';
	import { assetStore, asset, isLoading } from '$lib/stores/asset.store';
	import DialogLoading from '$lib/components/shared/DialogLoading.svelte';
	import { MapPin, Tag, Calendar, Package, ImageIcon } from '@lucide/svelte';

	export let assetId: string;
	export let open: boolean = false;
	export let onClose: () => void = () => {};

	let errorState = false;
	let errorMessage = '';

	onMount(async () => {
		// Load data when component mounts and has required props
		if (assetId && open) {
			await loadAssetData();
		}
	});

	async function loadAssetData() {
		try {
			errorState = false;
			errorMessage = '';

			const result = await assetStore.getById(assetId);

			if (!result) {
				setError('Asset not found or failed to load');
			}
		} catch (error) {
			console.error('Error loading asset:', error);
			setError('Failed to load asset data');
		}
	}

	function setError(message: string) {
		errorState = true;
		errorMessage = message;
	}

	function handleClose() {
		// Reset states when closing
		errorState = false;
		errorMessage = '';
		onClose();
	}

	// Reactive statement to handle when dialog opens with new assetId
	$: if (open && assetId) {
		loadAssetData();
	}
</script>

<Dialog.Root bind:open onOpenChange={handleClose}>
	<Dialog.Content class="max-h-[90vh] max-w-3xl overflow-hidden p-0">
		{#if $isLoading}
			<!-- Loading State -->
			<DialogLoading loadingText="Loading asset details..." />
		{:else if errorState}
			<!-- Error State -->
			<DialogError onClose={handleClose} title="Error Loading Asset" description={errorMessage} />
		{:else if $asset}
			<!-- Asset Content -->
			<div class="flex flex-col">
				<!-- Header with Image Background -->
				<div class="relative">
					{#if $asset.image}
						<!-- Asset Image -->
						<div class="h-48 overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200">
							<img
								src={$asset.image}
								alt={$asset.name}
								class="h-full w-full object-cover"
								loading="lazy"
							/>
							<div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent"></div>
						</div>
					{:else}
						<!-- Placeholder -->
						<div
							class="flex h-48 items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200"
						>
							<ImageIcon class="h-16 w-16 text-gray-400" />
						</div>
					{/if}

					<!-- Asset Header Info Overlay -->
					<div
						class="absolute right-0 bottom-0 left-0 bg-gradient-to-t from-black/60 to-transparent p-6"
					>
						<div class="text-white">
							<h1 class="mb-2 text-2xl font-bold">{$asset.name}</h1>
							<div class="flex items-center gap-4">
								<span class="text-3xl font-bold">
									{formatPrice($asset.price)}
								</span>
								<StatusBadge value={$asset.condition || 'Unknown'} />
							</div>
						</div>
					</div>
				</div>

				<!-- Content Body -->
				<div class="space-y-6 p-6">
					<!-- Description -->
					{#if $asset.description}
						<div>
							<h3 class="mb-2 text-sm font-semibold tracking-wide text-gray-900 uppercase">
								Description
							</h3>
							<p class="leading-relaxed text-gray-700">
								{$asset.description}
							</p>
						</div>
					{/if}

					<!-- Quick Info Cards -->
					<div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
						<!-- Location -->
						<div class="rounded-lg border border-blue-100 bg-blue-50 p-4">
							<div class="mb-1 flex items-center gap-2">
								<MapPin class="h-4 w-4 text-blue-600" />
								<span class="text-xs font-medium tracking-wide text-blue-900 uppercase"
									>Location</span
								>
							</div>
							<p class="text-sm font-semibold text-blue-800">
								{$asset.location?.name || 'Not specified'}
							</p>
						</div>

						<!-- Category -->
						<div class="rounded-lg border border-purple-100 bg-purple-50 p-4">
							<div class="mb-1 flex items-center gap-2">
								<Tag class="h-4 w-4 text-purple-600" />
								<span class="text-xs font-medium tracking-wide text-purple-900 uppercase"
									>Category</span
								>
							</div>
							<p class="text-sm font-semibold text-purple-800">
								{$asset.category?.name || 'Uncategorized'}
							</p>
						</div>

						<!-- Serial Number -->
						<div class="rounded-lg border border-green-100 bg-green-50 p-4">
							<div class="mb-1 flex items-center gap-2">
								<Package class="h-4 w-4 text-green-600" />
								<span class="text-xs font-medium tracking-wide text-green-900 uppercase"
									>Serial</span
								>
							</div>
							<p class="font-mono text-sm font-semibold text-green-800">
								{$asset.serialNumber || 'N/A'}
							</p>
						</div>

						<!-- Purchase Date -->
						<div class="rounded-lg border border-orange-100 bg-orange-50 p-4">
							<div class="mb-1 flex items-center gap-2">
								<Calendar class="h-4 w-4 text-orange-600" />
								<span class="text-xs font-medium tracking-wide text-orange-900 uppercase"
									>Purchased</span
								>
							</div>
							<p class="text-sm font-semibold text-orange-800">
								{#if $asset.purchaseDate}
									{formatDate($asset.purchaseDate)}
								{:else}
									Unknown
								{/if}
							</p>
						</div>
					</div>

					<!-- Additional Details -->
					{#if $asset.warranty}
						<div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
							<h3 class="mb-3 text-sm font-semibold tracking-wide text-gray-900 uppercase">
								Additional Information
							</h3>
							<div class="space-y-2">
								{#if $asset.warranty}
									<div class="flex items-center justify-between text-sm">
										<span class="text-gray-600">Warranty Until:</span>
										<span class="font-medium">{formatDate($asset.warranty)}</span>
									</div>
								{/if}
							</div>
						</div>
					{/if}
				</div>
			</div>
		{:else}
			<!-- No Asset State -->
			<DialogError
				onClose={handleClose}
				title="Asset not found"
				description="The asset you're looking for doesn't exist or has been removed."
			/>
		{/if}
	</Dialog.Content>
</Dialog.Root>
