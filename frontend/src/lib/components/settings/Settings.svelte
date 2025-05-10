<script>
	import GitConfiguration from '$lib/components/settings/tabs/git/GitConfiguration.svelte';
	import AdminConfiguration from '$lib/components/settings/tabs/git/AdminConfiguration.svelte';
	import AdvancedConfiguration from '$lib/components/settings/tabs/git/AdvancedConfiguration.svelte';

	function toggleTab(tabId) {
		const tabs = document.querySelectorAll('.tab-content');
		tabs.forEach((tab) => tab.classList.add('hidden'));

		const activeTab = document.getElementById(tabId);
		activeTab.classList.remove('hidden');
	}

	function closeSettingsModal() {
		const modal = document.getElementById('settingsModal');
		modal.classList.add('hidden');
	}

	function toggleMenu() {
		const menu = document.getElementById('menuNavigation');
		menu.classList.toggle('hidden');
	}
</script>

<!-- Configuration Modal -->
<div
	id="settingsModal"
	class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center hidden"
>
	<div
		class="bg-gray-800 text-white rounded-2xl shadow-lg w-full max-w-4xl h-full max-h-[90%] flex flex-col md:flex-row relative">
		<!-- Close Button -->
		<button
			class="absolute top-4 right-4 bg-gray-700 text-gray-400 hover:text-white w-10 h-10 rounded-full flex items-center justify-center border border-gray-500"
			on:click={()=>closeSettingsModal()}
		>
			<i class="fas fa-times"></i>
		</button>

		<!-- Left Navigation -->
		<div class="w-full md:w-1/4 bg-gray-700 p-6 flex flex-col md:space-y-4 md:pr-6">
			<!-- Hamburger Menu for Mobile -->
			<button
				class="md:hidden flex items-center bg-gray-600 text-white py-2 px-4 rounded mb-4"
				on:click={()=>toggleMenu()}
			>
				<i class="fas fa-bars text-lg mr-2"></i>
				<span>Menu</span>
			</button>
			<div id="menuNavigation" class="space-y-4 hidden md:block">
				<h2 class="text-xl font-bold mb-6">Configuration</h2>
				<button
					class="flex items-center space-x-2 text-white hover:bg-gray-600 py-2 px-4 rounded"
					on:click={()=>toggleTab('adminSettingsTab')}
				>
					<i class="fas fa-user-cog text-lg"></i>
					<span>Admin</span>
				</button>
				<button
					class="flex items-center space-x-2 text-white hover:bg-gray-600 py-2 px-4 rounded"
					on:click={()=>toggleTab('gitConfigTab')}
				>
					<i class="fas fa-code-branch text-lg"></i>
					<span>Git</span>
				</button>
				<button
					class="flex items-center space-x-2 text-white hover:bg-gray-600 py-2 px-4 rounded"
					on:click={()=>toggleTab('advancedSettingsTab')}
				>
					<i class="fas fa-cogs text-lg"></i>
					<span>Advanced</span>
				</button>
			</div>
		</div>

		<!-- Right Content -->
		<div class="w-full md:w-3/4 p-6 md:pt-12 overflow-y-auto">
			<AdminConfiguration/>
			<GitConfiguration />
			<AdvancedConfiguration/>
		</div>
	</div>
</div>
