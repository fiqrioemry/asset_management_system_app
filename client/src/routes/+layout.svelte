<script lang="ts">
	import '../app.css';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import ToastContainer from '$lib/components/ui/toast/ToastContainer.svelte';
	import ConfirmDialog from '$lib/components/common/ConfirmDialog.svelte';
</script>

<Toaster position="top-right" richColors={true} closeButton={true} duration={2000} />
<ToastContainer />
<ConfirmDialog />
<slot />
