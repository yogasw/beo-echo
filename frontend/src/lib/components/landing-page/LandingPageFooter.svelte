<script lang="ts">
	import { base } from '$app/paths';

	const year = new Date().getFullYear();

	const explore = [
		{ label: 'Concepts', href: `${base}/#concepts` },
		{ label: 'Operating modes', href: `${base}/#modes` },
		{ label: 'Connect to Claude (MCP)', href: `${base}/#install` },
		{ label: 'Deploy with Docker', href: `${base}/#deploy` }
	];

	const resources = [
		{ label: 'GitHub repository', href: 'https://github.com/yogasw/beo-echo' },
		{ label: 'Report an issue', href: 'https://github.com/yogasw/beo-echo/issues' },
		{ label: 'Product Hunt', href: 'https://www.producthunt.com/products/beo-echo' }
	];
</script>

<!-- Footer -->
<footer class="bg-gray-50 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800">
	<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
		<div class="grid grid-cols-1 md:grid-cols-12 gap-10">
			<!-- Brand -->
			<div class="md:col-span-6">
				<div class="flex items-center gap-2.5 mb-3">
					<span class="w-9 h-9 rounded-lg bg-blue-600 text-white flex items-center justify-center">
						<i class="fas fa-bolt"></i>
					</span>
					<span class="text-lg font-bold text-gray-900 dark:text-white">Beo Echo</span>
				</div>
				<p class="text-sm text-gray-600 dark:text-gray-400 max-w-sm leading-relaxed">
					Mock, proxy &amp; replay any HTTP API — self-hosted, and drivable from Claude over MCP.
				</p>
				<div class="flex items-center gap-3 mt-5">
					<a
						href="https://github.com/yogasw/beo-echo"
						target="_blank"
						rel="noopener noreferrer"
						class="w-9 h-9 rounded-lg border border-gray-200 dark:border-gray-700 flex items-center justify-center text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:border-blue-400 transition-colors"
						title="GitHub"
						aria-label="Beo Echo on GitHub"
					>
						<i class="fab fa-github"></i>
					</a>
					<a
						href="https://www.producthunt.com/products/beo-echo"
						target="_blank"
						rel="noopener noreferrer"
						class="w-9 h-9 rounded-lg border border-gray-200 dark:border-gray-700 flex items-center justify-center text-gray-600 dark:text-gray-300 hover:text-[#ff6154] hover:border-[#ff6154]/50 transition-colors"
						title="Product Hunt"
						aria-label="Beo Echo on Product Hunt"
					>
						<i class="fab fa-product-hunt"></i>
					</a>
				</div>
			</div>

			<!-- Explore -->
			<div class="md:col-span-3">
				<h4 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400 mb-4">
					Explore
				</h4>
				<ul class="space-y-2.5 text-sm">
					{#each explore as link}
						<li>
							<a
								href={link.href}
								class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
							>
								{link.label}
							</a>
						</li>
					{/each}
				</ul>
			</div>

			<!-- Resources -->
			<div class="md:col-span-3">
				<h4 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400 mb-4">
					Resources
				</h4>
				<ul class="space-y-2.5 text-sm">
					{#each resources as link}
						<li>
							<a
								href={link.href}
								target="_blank"
								rel="noopener noreferrer"
								class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
							>
								{link.label}
							</a>
						</li>
					{/each}
				</ul>
			</div>
		</div>

		<div
			class="mt-12 pt-6 border-t border-gray-200 dark:border-gray-800 flex flex-col sm:flex-row items-center justify-between gap-3"
		>
			<p class="text-xs text-gray-500 dark:text-gray-400">
				© {year} Beo Echo · MIT Licensed
			</p>
			<p class="text-xs text-gray-400 dark:text-gray-500">Built with ❤️ for developers</p>
		</div>
	</div>
</footer>
