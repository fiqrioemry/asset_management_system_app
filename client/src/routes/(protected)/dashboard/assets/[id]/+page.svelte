<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	onMount(() => {
		// Redirect ke products page karena modal akan handle tampilan
		goto('/dashboard/product', { replaceState: true });
	});
</script>

<!-- This page only handles the redirect, no content needed -->
