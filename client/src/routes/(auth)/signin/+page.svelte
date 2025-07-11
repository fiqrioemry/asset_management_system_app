<script lang="ts">
	import { authStore } from '$lib/stores/auth.store';
	import LoginForm from '$lib/components/business/LoginForm.svelte';
</script>

<svelte:head>
	<title>Sign In</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gray-100">
	<div class="bg-background w-full max-w-md rounded-lg p-4 shadow-md">
		<div class="mb-6 text-center">
			<h1 class="text-2xl font-bold">Welcome Back</h1>
			<p class="text-muted-foreground text-sm">Please enter your credentials to continue</p>
		</div>
		<LoginForm />

		<div class="text-muted-foreground mt-4 text-center text-sm">
			<p class="switch-mode">
				Don't have an account? <a href="/signup" on:click={() => authStore.clearError()}>
					Sign up here
				</a>
			</p>
		</div>
	</div>
</div>
