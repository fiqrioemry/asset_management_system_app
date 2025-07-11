<!-- src/routes/dashboard/+layout.svelte -->
<script lang="ts">
	import Navbar from '$lib/components/layout/Navbar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppSidebar from '$lib/components/layout/AppSidebar.svelte';
	import ProtectedRoute from '$lib/components/middleware/ProtectedRoute.svelte';
</script>

<!-- <ProtectedRoute> -->
<Sidebar.Provider>
	<div class="flex min-h-screen w-full">
		<!-- Sidebar -->
		<div class="hidden md:block">
			<AppSidebar />
		</div>

		<!-- Main Content Area -->
		<div class="b flex w-full flex-1 flex-col">
			<!-- Top Navbar -->
			<Navbar />

			<!-- Page Content -->
			<main class="flex-1 overflow-auto">
				<slot />
			</main>
		</div>
	</div>

	<!-- Mobile Sidebar Overlay (optional) -->
	<Sidebar.Trigger
		class="bg-background border-border fixed top-4 left-4 z-50 rounded-lg border p-2 shadow-lg md:hidden"
	>
		<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M4 6h16M4 12h16M4 18h16"
			></path>
		</svg>
	</Sidebar.Trigger>
</Sidebar.Provider>
<!-- </ProtectedRoute> -->
