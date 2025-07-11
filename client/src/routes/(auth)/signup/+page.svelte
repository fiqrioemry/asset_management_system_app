<script lang="ts">
	import { authStore } from '$lib/stores/auth.store';
	import RegisterForm from '$lib/components/business/RegisterForm.svelte';
</script>

<svelte:head>
	<title>Sign Up | Asset Management System</title>
	<meta
		name="description"
		content="Create an account to manage your assets efficiently. Sign up now!"
	/>
</svelte:head>

<div class="flex h-screen items-center justify-center bg-gray-100">
	<div class="bg-background w-full max-w-md rounded-lg p-4 shadow-md">
		<div class="mb-6 text-center">
			<h1 class="text-2xl font-bold">Create Your Account</h1>
			<p class="text-muted-foreground text-sm">
				Please fill in the details below to create your account.
			</p>
		</div>
		<RegisterForm />

		<div class="text-muted-foreground mt-4 text-center text-sm">
			<p class="switch-mode">
				Already have an account? <a href="/signin" on:click={() => authStore.clearError()}
					>Sign in here</a
				>
			</p>
		</div>
	</div>
</div>
