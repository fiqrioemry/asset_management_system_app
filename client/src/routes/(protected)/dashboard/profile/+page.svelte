<script lang="ts">
	import {
		Lock,
		User,
		Clock,
		Camera,
		Loader2,
		Calendar,
		RefreshCw,
		CheckCircle
	} from '@lucide/svelte';
	import { cn } from '$lib/utils';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import StatCard from '$lib/components/ui/StatCard.svelte';
	import AlertCard from '$lib/components/ui/AlertCard.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import FormDialog from '$lib/components/common/FormDialog.svelte';
	import { formatDate, getAvatarInitials } from '$lib/utils/formatter';
	import SubmitButton from '$lib/components/forms/SubmitButton.svelte';
	import InlineTextEdit from '$lib/components/forms/InlineTextEdit.svelte';
	import ProfileLoading from '$lib/components/loading/ProfileLoading.svelte';
	import ChangePasswordForm from '$lib/components/business/ChangePasswordForm.svelte';
	import { userStore, currentUserProfile, isUserLoading, userError } from '$lib/stores/user.store';

	let avatarFile: File | null = null;
	let fileInput: HTMLInputElement | null = null;
	let showChangePassword = false;

	onMount(async () => {
		await userStore.getUser();
	});

	async function handleAvatarChange(event: any) {
		const file = event.target.files[0];
		if (!file) return;

		if (!file.type.startsWith('image/')) {
			toast.error('Please select an image file');
			return;
		}

		// Validate file size (max 1MB)
		if (file.size > 1 * 1024 * 1024) {
			toast.error('File size must be less than 1MB');
			return;
		}

		avatarFile = file;
		await uploadAvatar(file);
	}

	async function uploadAvatar(file: any) {
		if (!file || !$currentUserProfile?.fullname) return;

		const updateData = {
			fullname: $currentUserProfile.fullname,
			avatar: file
		};

		const result = await userStore.updateUser(updateData);

		if (result.success) {
			avatarFile = null; // Clear file reference on success
		}
	}

	function triggerFileInput() {
		if (fileInput) {
			fileInput.click();
		}
	}

	async function handleRefresh() {
		const result = await userStore.refreshUser();

		if (result.success) {
			toast.success('Profile refreshed!');
		} else {
			toast.error('Failed to refresh profile');
		}
	}

	async function handleUpdateFullname(newValue: string) {
		if (!newValue || newValue.trim() === '' || newValue.trim().length < 3) {
			toast.error('Fullname required and must be at least 3 characters');
			return false;
		}
		const result = await userStore.updateUser({ fullname: newValue });
		return result.success;
	}

	function openChangePassword() {
		showChangePassword = true;
	}

	function handlePasswordChangeSuccess() {
		showChangePassword = false;
	}

	function handlePasswordChangeCancel() {
		showChangePassword = false;
	}

	async function handleErrorRetry() {
		await userStore.getUser();
	}
</script>

<FormDialog
	size="md"
	title="Change Password"
	bind:open={showChangePassword}
	description="Update your account password for better security"
>
	<ChangePasswordForm
		onSuccess={handlePasswordChangeSuccess}
		onCancel={handlePasswordChangeCancel}
	/>
</FormDialog>

<div class="bg-muted min-h-screen py-8">
	<div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
		<!-- Loading State -->
		{#if $isUserLoading}
			<ProfileLoading />
		{/if}

		<!-- Error State -->
		{#if $userError}
			<AlertCard onAction={handleErrorRetry} message={$userError} />
		{/if}

		<!-- Profile Content -->
		{#if $currentUserProfile && !$isUserLoading}
			<div class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm">
				<!-- Profile Header -->
				<div class="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-8">
					<div class="flex items-center space-x-6">
						<!-- Avatar with Edit -->
						<div class="relative flex-shrink-0">
							<!-- Loading overlay for avatar updates -->
							{#if $isUserLoading && avatarFile}
								<div
									class="absolute inset-0 z-10 flex items-center justify-center rounded-full bg-black/50"
								>
									<Loader2 class="h-8 w-8 animate-spin text-white" />
								</div>
							{/if}

							{#if $currentUserProfile.avatar}
								<img
									src={$currentUserProfile.avatar}
									alt="Profile Avatar"
									class={cn('h-20 w-20 rounded-full border-4 border-white object-cover shadow-lg', {
										'opacity-50': $isUserLoading && avatarFile
									})}
								/>
							{:else}
								<div
									class="flex h-20 w-20 items-center justify-center rounded-full border-4 border-white bg-gray-100 shadow-lg"
								>
									<span class="text-muted-foreground text-xl font-semibold">
										{getAvatarInitials($currentUserProfile.fullname)}
									</span>
								</div>
							{/if}

							<!-- Avatar Edit Button -->
							<button
								class="btn-float"
								aria-label="Change avatar"
								on:click={triggerFileInput}
								disabled={$isUserLoading}
							>
								<Camera class="h-4 w-4" />
							</button>

							<!-- Hidden File Input -->
							<input
								type="file"
								accept="image/*"
								class="hidden"
								bind:this={fileInput}
								on:change={handleAvatarChange}
							/>
						</div>

						<!-- User Headline Info -->
						<div class="flex-1">
							<div class="flex items-center space-x-2">
								<h2 class="text-2xl font-bold text-white">{$currentUserProfile.fullname}</h2>
							</div>
							<p class="text-blue-100">{$currentUserProfile.email}</p>
							<div class="mt-2 flex items-center text-blue-100">
								<Calendar class="mr-2 h-4 w-4" />
								<span class="text-sm">Member since {formatDate($currentUserProfile.joinedAt)}</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Profile Details -->
				<div class="px-6 py-6">
					<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
						<!-- Account Information -->
						<div>
							<h3 class="text-foreground mb-4 text-lg font-medium">Account Information</h3>
							<div class="space-y-4">
								<InlineTextEdit
									name="fullname"
									label="Fullname"
									loading={$isUserLoading}
									value={$currentUserProfile.fullname}
									onSave={handleUpdateFullname}
								/>

								<!-- Email Address (Read Only) -->
								<div>
									<span class="block text-sm font-medium">Email Address</span>
									<div
										class="mt-1 flex items-center justify-between rounded-md bg-gray-100 px-3 py-2"
									>
										<span class="text-foreground text-sm">{$currentUserProfile.email}</span>
										<span class="text-xs text-gray-500">Cannot be changed</span>
									</div>
								</div>
							</div>
						</div>

						<!-- Profile Statistics -->
						<div>
							<h3 class="text-foreground mb-4 text-lg font-medium">Account Statistics</h3>
							<div class="space-y-4">
								<StatCard icon={User} title="Account Status" value="Active" variant="blue" />
								<StatCard icon={CheckCircle} title="Email Verified" value="Yes" variant="purple" />
								<StatCard
									icon={Clock}
									variant="green"
									title="Member since"
									value={formatDate($currentUserProfile.joinedAt)}
								/>
							</div>
						</div>
					</div>
				</div>

				<!-- Action Buttons -->
				<div class="bg-muted border-muted border-t p-6">
					<div class="flex flex-wrap gap-3">
						<Button onclick={openChangePassword} variant="outline" class="w-full md:w-60">
							<Lock class="mr-2 h-4 w-4" />
							Change Password
						</Button>

						<SubmitButton
							buttonText="Refresh"
							onsubmit={handleRefresh}
							className="md:w-60 w-full"
							isLoading={$isUserLoading}
							buttonLoadingText="Refreshing..."
						>
							<RefreshCw class="mr-2 h-4 w-4" />
						</SubmitButton>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
