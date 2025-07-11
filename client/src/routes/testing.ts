
<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import AssetUpdate from '$lib/components/assets/AssetUpdate.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Eye, Trash2 } from '@lucide/svelte';
	import { assets } from '$lib/stores/asset.store';
	
	// ... other imports and logic
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>Asset Name</Table.Head>
			<Table.Head>Location</Table.Head>
			<Table.Head>Price</Table.Head>
			<Table.Head>Condition</Table.Head>
			<Table.Head>Actions</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $assets as asset}
			<Table.Row>
				<Table.Cell>{asset.name}</Table.Cell>
				<Table.Cell>{asset.location?.name}</Table.Cell>
				<Table.Cell>{asset.price}</Table.Cell>
				<Table.Cell>{asset.condition}</Table.Cell>
				<Table.Cell>
					<div class="flex items-center space-x-2">
						<!-- View Button -->
						<Button size="icon" variant="outline">
							<Eye class="h-4 w-4" />
						</Button>
						
						<!-- Edit Button (using AssetUpdate component) -->
						<AssetUpdate {asset} />
						
						<!-- Delete Button -->
						<Button size="icon" variant="outline">
							<Trash2 class="h-4 w-4 text-red-500" />
						</Button>
					</div>
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>

// ========================================
// 5. ADVANCED PATTERN - WITH CUSTOM TRIGGER
// ========================================

// Custom trigger button
<script lang="ts">
	import AssetUpdate from '$lib/components/assets/AssetUpdate.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Settings } from '@lucide/svelte';
	
	export let asset: Asset;
	
	// Custom trigger component
	let CustomTrigger = Button;
</script>

<AssetUpdate 
	{asset} 
	buttonElement={CustomTrigger}
	let:handleTriggerClick
>
	<svelte:fragment slot="trigger">
		<Button 
			size="sm" 
			variant="secondary"
			on:click={handleTriggerClick}
		>
			<Settings class="h-4 w-4 mr-2" />
			Edit Asset
		</Button>
	</svelte:fragment>
</AssetUpdate>

// ========================================
// 6. KEY DIFFERENCES FROM REACT
// ========================================

/*
REACT vs SVELTE PATTERNS:

✅ REACT (your current):
- useForm hook for form state
- FormProvider for context
- Controlled components with value/onChange

✅ SVELTE (new pattern):
- bind:value for two-way binding
- Slot props for passing data down
- Custom events for communication
- Built-in reactivity

✅ DATA PASSING:
REACT: Props drilling + Context
SVELTE: Slot props + bind: directives

✅ FORM HANDLING:
REACT: react-hook-form + validation
SVELTE: Manual state + validation utility

✅ EVENT HANDLING:
REACT: Callback props
SVELTE: Custom events (createEventDispatcher)

✅ STATE MANAGEMENT:
REACT: useState + useEffect
SVELTE: Reactive statements ($:) + stores
*/

// ========================================
// 7. USAGE EXAMPLES SUMMARY
// ========================================

/*
✅ BASIC USAGE:
<AssetUpdate {asset} />

✅ WITH CUSTOM TRIGGER:
<AssetUpdate {asset} buttonElement={CustomButton} />

✅ WITH EVENT HANDLERS:
<AssetUpdate 
  {asset} 
  on:success={handleSuccess}
  on:error={handleError}
/>

✅ IN TABLE CELL:
<Table.Cell>
  <AssetUpdate {asset} />
</Table.Cell>

✅ MULTIPLE ACTIONS:
<div class="flex gap-2">
  <AssetView {asset} />
  <AssetUpdate {asset} />
  <AssetDelete {asset} />
</div>
*/