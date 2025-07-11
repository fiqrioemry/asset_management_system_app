<script lang="ts">
	import { goto } from '$app/navigation';
	import { authService } from '$lib/services/auth.service';
	import { authStore, currentUser } from '$lib/stores/auth.store';

	async function handleLogout() {
		try {
			await authService.logout();
			authStore.logout();
			goto('/signin');
		} catch (error) {
			console.error('Logout failed:', error);
			authStore.logout();
			goto('/signin');
		}
	}
</script>

<svelte:head>
	<title>Dashboard</title>
</svelte:head>

<div class="dashboard">
	<header>
		<h1>Dashboard</h1>
		<button on:click={handleLogout}>Logout</button>
	</header>

	{#if $currentUser}
		<div class="user-info">
			<h2>Welcome, {$currentUser.fullname}!</h2>
			<div class="user-details">
				<p><strong>Email:</strong> {$currentUser.email}</p>
				<p><strong>User ID:</strong> {$currentUser.id}</p>
				<p><strong>Joined:</strong> {new Date($currentUser.joinedAt).toLocaleDateString()}</p>
				{#if $currentUser.avatar}
					<img src={$currentUser.avatar} alt="Avatar" class="avatar" />
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.dashboard {
		max-width: 800px;
		margin: 0 auto;
		padding: 20px;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 30px;
		padding-bottom: 20px;
		border-bottom: 1px solid #eee;
	}

	h1 {
		margin: 0;
		color: #333;
	}

	button {
		padding: 10px 20px;
		background-color: #dc3545;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	button:hover {
		background-color: #c82333;
	}

	.user-info {
		background-color: #f8f9fa;
		padding: 20px;
		border-radius: 8px;
	}

	.user-details {
		margin-top: 15px;
	}

	.user-details p {
		margin: 10px 0;
	}

	.avatar {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		margin-top: 10px;
	}
</style>
