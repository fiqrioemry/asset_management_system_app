<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import { formatPrice } from '$lib/utils/formatter';
	import { Edit, Eye, Trash2 } from '@lucide/svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import type { Asset } from '$lib/types/global.types';
	import Button from '$lib/components/ui/button/button.svelte';
	import TableLoading from '$lib/components/loading/TableLoading.svelte';

	// state props imports
	export let loading: boolean = false;
	export let assets: Asset[] = [];

	// functions props imports
	export let onDeleteAsset: (id: string) => void = () => {};
	export let onOpenDetail: (asset: Asset) => void = () => {};
	export let onUpdateAsset: (asset: Asset) => void = () => {};
</script>

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
		{#if loading}
			<TableLoading rows={5} />
		{:else if assets && assets.length > 0}
			{#each assets as asset}
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
						<Badge variant="secondary">{asset.condition}</Badge>
					</Table.Cell>
					<Table.Cell>
						{asset?.category?.name || 'N/A'}
					</Table.Cell>

					<Table.Cell>
						<div class="flex items-center space-x-2">
							<Button
								size="icon"
								variant="outline"
								class="rounded-full"
								onclick={() => onOpenDetail(asset)}
							>
								<Eye class="h-3 w-3" />
							</Button>
							<Button
								size="icon"
								variant="outline"
								class="rounded-full"
								onclick={() => onUpdateAsset(asset)}
							>
								<Edit class="h-3 w-3" />
							</Button>

							<Button
								size="icon"
								variant="outline"
								class="rounded-full bg-red-100"
								onclick={() => onDeleteAsset(asset.id)}
							>
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
