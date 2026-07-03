<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { toast } from '$lib/stores/toast';
	import LandingPageHeader from '$lib/components/landing-page/LandingPageHeader.svelte';
	import LandingPageFooter from '$lib/components/landing-page/LandingPageFooter.svelte';
	import MockScene from './MockScene.svelte';
	import { STATIC_MODE } from '$lib/config/appMode';

	const dockerCmd =
		'docker run -d --platform linux/amd64 -p 8080:80 -v $(pwd)/beo-echo-config:/app/configs/ ghcr.io/yogasw/beo-echo:latest';
	import {
		concepts,
		modes,
		installMethods,
		useCases,
		migrationStages
	} from './guideData';

	let endpoint = 'https://<your-beo-echo-host>/mcp';

	onMount(async () => {
		// In a pure static (no-backend) build the origin is a docs host, not a
		// Beo Echo instance — keep the placeholder so the command stays correct.
		if (STATIC_MODE) return;
		await auth.initialize();
		if (browser) {
			endpoint = `${window.location.origin}/mcp`;
		}
	});

	// Substitute the live endpoint into any code snippet.
	function render(line: string): string {
		return line.replaceAll('{ENDPOINT}', endpoint);
	}

	async function copyCode(lines: string[]) {
		try {
			await navigator.clipboard.writeText(lines.map(render).join('\n'));
			toast.success('Copied');
		} catch {
			toast.error('Could not copy — select and copy manually');
		}
	}

	// diagram helpers per mode
	const forwards = (id: string) => id === 'proxy' || id === 'forwarder';
	const isOff = (id: string) => id === 'disabled';
	const isCut = (id: string) => id === 'mock';
	// proxy forwards conditionally (only when no mock matches)
	const isCond = (id: string) => id === 'proxy';
	// a matched mock is answered at Beo Echo and stops there (mock + proxy-match)
	const stopsAtBeo = (id: string) => id === 'mock' || id === 'proxy';
	function wire2Class(id: string): string {
		if (id === 'mock' || id === 'disabled') return 'cut';
		if (id === 'proxy') return 'cond';
		return '';
	}

	async function handleLogout() {
		await auth.logout();
		window.location.reload();
	}
</script>

<div class="flex-1 w-full bg-white dark:bg-gray-900 transition-colors duration-200 flex flex-col">
	<LandingPageHeader showUserMenu={true} on:logout={handleLogout} on:back={() => goto('/')} />

	<main class="flex-1">
		<!-- ───────────────── Hero ───────────────── -->
		<section
			class="relative overflow-hidden bg-gradient-to-b from-blue-50 to-white dark:from-gray-800 dark:to-gray-900 py-12 pt-24"
		>
			<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
				<span
					class="inline-flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-blue-600 dark:text-blue-400 bg-blue-100 dark:bg-blue-500/15 px-3 py-1 rounded-full mb-5"
				>
					<i class="fas fa-bolt"></i> Mock · Proxy · Replay
				</span>
				<h1 class="text-3xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4">
					Mock, proxy &amp; replay any HTTP API
				</h1>
				<p class="text-lg text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
					Stand up mock endpoints in seconds, forward to the real service when you're ready, and
					drive it all from the dashboard — or from Claude over MCP.
				</p>

				<div class="flex flex-wrap gap-3 justify-center mt-8">
					<a
						href="#deploy"
						class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
					>
						<i class="fab fa-docker"></i> Deploy in seconds
					</a>
					<a
						href="#concepts"
						class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 hover:border-blue-400 text-gray-700 dark:text-gray-300 px-5 py-2.5 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
					>
						<i class="fas fa-diagram-project"></i> Learn the concepts
					</a>
				</div>
			</div>
		</section>

		<!-- ───────────────── Why: test now, wire real API later ───────────────── -->
		<section class="py-16 bg-gray-50 dark:bg-gray-800">
			<div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-10">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-3">
						<i class="fas fa-user-clock text-blue-600 mr-2"></i> Test today, wire up the real API later
					</h2>
					<p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
						You and QA keep testing against Beo Echo mocks while the backend is still being built —
						nobody waits. The scene below cycles through the rollout:
						<span class="font-medium text-red-500">not ready</span> →
						<span class="font-medium text-blue-500">partial</span> →
						<span class="font-medium text-green-600 dark:text-green-500">real API ready</span>,
						forwarding endpoints one at a time.
					</p>
				</div>

				<MockScene />

				<p class="text-center text-sm text-gray-500 dark:text-gray-400 mt-8">
					That gradual hand-off is exactly the
					<a href="#to-production" class="text-blue-600 dark:text-blue-400 font-medium hover:underline">
						mock → production path
					</a> below.
				</p>
			</div>
		</section>

		<!-- ───────────────── Concepts / hierarchy ───────────────── -->
		<section id="concepts" class="py-16">
			<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-3">
						<i class="fas fa-layer-group text-blue-600 mr-2"></i> The building blocks
					</h2>
					<p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
						Everything nests. A workspace holds projects; a project holds endpoints; each endpoint
						holds responses; a response can carry rules.
					</p>
				</div>

				<!-- Animated nested hierarchy -->
				<div class="flex justify-center mb-14">
					<div class="hierarchy w-full max-w-md space-y-0">
						{#each concepts as c, i}
							<div
								class="nest {c.bg}"
								style="--i:{i}; margin-left:{i * 18}px; margin-right:{i === 0 ? 0 : 0}px;"
							>
								<div class="flex items-center gap-3 py-2.5 px-3">
									<span class="w-8 h-8 rounded-lg {c.bg} flex items-center justify-center {c.color}">
										<i class="{c.icon} text-sm"></i>
									</span>
									<div class="text-left">
										<div class="font-semibold text-gray-900 dark:text-white text-sm leading-tight">
											{c.title}
										</div>
										<div class="text-xs text-gray-500 dark:text-gray-400">{c.short}</div>
									</div>
									{#if i < concepts.length - 1}
										<i class="fas fa-angle-down text-gray-300 dark:text-gray-600 ml-auto"></i>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				</div>

				<!-- Concept detail cards -->
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each concepts as c, i}
						<div
							class="reveal bg-white dark:bg-gray-800 rounded-xl shadow-md border border-gray-100 dark:border-gray-700 p-6"
							style="--d:{i * 80}ms"
						>
							<div class="flex items-center gap-3 mb-3">
								<span class="w-10 h-10 rounded-lg {c.bg} flex items-center justify-center {c.color}">
									<i class="{c.icon}"></i>
								</span>
								<h3 class="font-semibold text-gray-900 dark:text-white">{c.title}</h3>
							</div>
							<p class="text-sm text-gray-600 dark:text-gray-300 mb-3">{c.detail}</p>
							<p
								class="text-xs font-mono text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/40 rounded-md px-3 py-2 break-words"
							>
								{c.example}
							</p>
						</div>
					{/each}
				</div>
			</div>
		</section>

		<!-- ───────────────── Operating modes (animated flows) ───────────────── -->
		<section id="modes" class="py-16 bg-gray-50 dark:bg-gray-800">
			<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-3">
						<i class="fas fa-cogs text-blue-600 mr-2"></i> The three modes (+ off)
					</h2>
					<p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
						A project handles every request one way. Watch where the request goes.
					</p>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					{#each modes as mode}
						<div
							class="bg-white dark:bg-gray-900/40 rounded-xl shadow-md border border-gray-100 dark:border-gray-700 p-6"
						>
							<div class="flex items-center gap-3 mb-5">
								<span class="w-10 h-10 {mode.color} rounded-lg flex items-center justify-center text-white">
									<i class={mode.icon}></i>
								</span>
								<div>
									<h3 class="font-semibold text-gray-900 dark:text-white leading-tight">{mode.name}</h3>
									<p class="text-xs {mode.accent} font-medium">{mode.tagline}</p>
								</div>
							</div>

							<!-- request-flow animation: request in → decide → out -->
							<div class="flow {mode.id} {isOff(mode.id) ? 'off' : ''}" aria-hidden="true">
								<div class="fnodes">
									<div class="fnode">
										<i class="fas fa-user"></i><span>User</span>
									</div>
									<div class="fnode">
										<i class="fas fa-laptop-code"></i><span>Your app</span>
									</div>
									<div class="fnode beo">
										<i class="fas fa-bolt"></i><span>Beo Echo</span>
									</div>
									<div class="fnode {isCut(mode.id) || isOff(mode.id) ? 'dim' : ''}">
										<i class="fas fa-database"></i><span>Real API</span>
									</div>
								</div>

								<div class="rails">
									<span class="rail rail0"></span>
									<span class="rail rail1"></span>
									<span class="rail rail2 {wire2Class(mode.id)}"></span>
									{#if isCut(mode.id) || isOff(mode.id)}<i class="fas fa-xmark cut-x" title="never forwarded"></i>{/if}
								</div>

								<div class="layer">
									<!-- FE processes the call on the way in (same in every mode) -->
									<span class="fe-proc"><i class="fas fa-gear"></i> processing…</span>
									{#if isOff(mode.id)}
										<!-- disabled: request still arrives from the user but the
										     project serves nothing, so it dies at Beo Echo -->
										<span class="chip chip-dead"><i class="fas fa-cube"></i></span>
										<span class="verdict v-dead"><i class="fas fa-ban"></i> no response</span>
									{:else}
										{#if stopsAtBeo(mode.id)}
											<span class="chip chip-match"><i class="fas fa-cube"></i></span>
											<span class="verdict v-match"><i class="fas fa-check"></i> match</span>
										{/if}
										{#if forwards(mode.id)}
											<span class="chip chip-fwd"><i class="fas fa-cube"></i></span>
											<span class="verdict v-fwd">
												{#if isCond(mode.id)}no match <i class="fas fa-arrow-right"></i>{:else}forward all <i class="fas fa-arrow-right"></i>{/if}
											</span>
										{/if}
									{/if}
								</div>
							</div>
							{#if isCond(mode.id)}
								<p class="text-[11px] text-gray-400 dark:text-gray-500 mt-3 text-center">
									<span class="text-green-500 font-semibold">match</span> → answered &amp; sent back
									&nbsp;·&nbsp;
									<span class="text-purple-500 font-semibold">no match</span> → forwarded on
								</p>
							{/if}

							<p class="text-sm text-gray-600 dark:text-gray-300 mt-5">{mode.behaviour}</p>
							<p class="text-xs text-gray-500 dark:text-gray-400 mt-2">
								<i class="fas fa-circle-info mr-1 {mode.accent}"></i> {mode.whenToUse}
							</p>
						</div>
					{/each}
				</div>
			</div>
		</section>

		<!-- ───────────────── Install / connect to Claude ───────────────── -->
		<section id="install" class="py-16">
			<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-3">
						<i class="fas fa-plug text-blue-600 mr-2"></i> Connect it to Claude (MCP)
					</h2>
					<p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
						One endpoint, three ways to connect. Your endpoint:
						<code
							class="font-mono text-sm text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20 px-2 py-0.5 rounded break-all"
							>{endpoint}</code
						>
					</p>
				</div>

				<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
					{#each installMethods as m}
						<div
							class="flex flex-col bg-white dark:bg-gray-800 rounded-xl shadow-md border border-gray-100 dark:border-gray-700 p-6"
						>
							<div class="flex items-center gap-3 mb-2">
								<i class="{m.icon} {m.accent} text-lg"></i>
								<h3 class="font-semibold text-gray-900 dark:text-white">{m.title}</h3>
							</div>
							<p class="text-xs text-gray-500 dark:text-gray-400 mb-4">{m.subtitle}</p>

							<ol class="space-y-2 mb-4">
								{#each m.steps as step, i}
									<li class="flex gap-2.5 text-sm text-gray-600 dark:text-gray-300">
										<span
											class="flex-none w-5 h-5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs font-semibold flex items-center justify-center"
											>{i + 1}</span
										>
										<span>{step}</span>
									</li>
								{/each}
							</ol>

							{#if m.code}
								<div class="relative mt-auto">
									<pre
										class="bg-gray-100 dark:bg-gray-900/60 text-gray-800 dark:text-gray-100 border border-gray-200 dark:border-gray-700 rounded-lg p-3 pr-10 overflow-x-auto text-xs font-mono leading-relaxed"><code
											>{m.code.map(render).join('\n')}</code
										></pre>
									<button
										on:click={() => m.code && copyCode(m.code)}
										class="absolute top-2 right-2 w-7 h-7 rounded flex items-center justify-center text-xs bg-white/80 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-500 dark:text-gray-300 border border-gray-200 dark:border-gray-600 transition-colors"
										title="Copy"
										aria-label="Copy snippet"
									>
										<i class="fas fa-copy"></i>
									</button>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		</section>

		<!-- ───────────────── Use cases ───────────────── -->
		<section class="py-16 bg-gray-50 dark:bg-gray-800">
			<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-3">
						<i class="fas fa-wand-magic-sparkles text-blue-600 mr-2"></i> Common use cases
					</h2>
					<p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
						What people actually ask Claude to do, once connected.
					</p>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					{#each useCases as uc}
						<div
							class="bg-white dark:bg-gray-900/40 rounded-xl shadow-md border border-gray-100 dark:border-gray-700 p-6"
						>
							<div class="flex items-center gap-3 mb-2">
								<i class="{uc.icon} {uc.color} text-lg"></i>
								<h3 class="font-semibold text-gray-900 dark:text-white">{uc.title}</h3>
							</div>
							<p class="text-sm text-gray-500 dark:text-gray-400 italic mb-4">{uc.scenario}</p>
							<ul class="space-y-2">
								{#each uc.steps as step}
									<li class="flex items-start gap-2 text-sm text-gray-600 dark:text-gray-300">
										<i class="fas fa-angle-right {uc.color} mt-1"></i><span>{step}</span>
									</li>
								{/each}
							</ul>
						</div>
					{/each}
				</div>
			</div>
		</section>

		<!-- ───────────────── Migration pipeline ───────────────── -->
		<section id="to-production" class="py-16">
			<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-14">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-3">
						<i class="fas fa-timeline text-blue-600 mr-2"></i> From mock to production
					</h2>
					<p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
						You rarely jump straight to a full mock. Grow into it — and back out — one mode at a
						time.
					</p>
				</div>

				<div class="pipeline">
					<div class="track">
						<span class="track-fill"></span>
						<span class="flow-dot"></span>
					</div>
					<div class="stages">
						{#each migrationStages as s, i}
							<div class="stage reveal" style="--d:{i * 140}ms; --sd:{i * 1.4}s">
								<span class="stage-icon {s.color}"><i class={s.icon}></i></span>
								<div class="stage-mode">{s.mode}</div>
								<div class="stage-title">{s.title}</div>
								<p class="stage-desc">{s.desc}</p>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</section>
		<!-- ───────────────── Deploy your own (Docker) ───────────────── -->
		<section id="deploy" class="py-16 bg-gray-50 dark:bg-gray-800">
			<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-8">
					<h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-3">
						<i class="fab fa-docker text-blue-600 mr-2"></i> Deploy in seconds
					</h2>
					<p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
						Run your own Beo Echo with a single Docker command.
					</p>
				</div>

				<!-- Docker command (terminal) -->
				<div class="bg-slate-900 rounded-lg p-5 mb-8 text-left overflow-x-auto border border-slate-700/70 shadow-sm">
					<div class="flex items-center justify-end mb-3">
						<button
							on:click={() => copyCode([dockerCmd])}
							class="bg-blue-500 hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700 text-white px-3 py-1 rounded text-xs font-medium transition-colors flex items-center gap-1"
							title="Copy Docker command"
							aria-label="Copy Docker command"
						>
							<i class="fas fa-copy"></i> Copy
						</button>
					</div>
					<code class="text-green-400 dark:text-green-300 font-mono text-sm block leading-relaxed">
						<span class="text-gray-400 dark:text-gray-500">$</span> docker run -d --platform linux/amd64 -p 8080:80 \<br />
						<span class="ml-4">-v $(pwd)/beo-echo-config:/app/configs/ \</span><br />
						<span class="ml-4">ghcr.io/yogasw/beo-echo:latest</span>
					</code>
				</div>

				<!-- Quick steps -->
				<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
					<div class="bg-white dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700 rounded-lg p-5 text-center">
						<div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-500/20 text-blue-600 dark:text-blue-300 font-bold flex items-center justify-center mx-auto mb-3">1</div>
						<h3 class="text-sm font-semibold text-gray-900 dark:text-white">Run command</h3>
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">In your terminal.</p>
					</div>
					<div class="bg-white dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700 rounded-lg p-5 text-center">
						<div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-500/20 text-blue-600 dark:text-blue-300 font-bold flex items-center justify-center mx-auto mb-3">2</div>
						<h3 class="text-sm font-semibold text-gray-900 dark:text-white">Open browser</h3>
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
							at <code class="font-mono">localhost:8080</code>
						</p>
					</div>
					<div class="bg-white dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700 rounded-lg p-5 text-center">
						<div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-500/20 text-blue-600 dark:text-blue-300 font-bold flex items-center justify-center mx-auto mb-3">3</div>
						<h3 class="text-sm font-semibold text-gray-900 dark:text-white">Log in &amp; start</h3>
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
							<code class="font-mono">admin@admin.com</code> / <code class="font-mono">admin</code>
						</p>
					</div>
				</div>

				<!-- Extra info -->
				<div class="bg-white dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700 rounded-lg p-5">
					<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
						<div>
							<h4 class="text-sm font-semibold text-gray-900 dark:text-white mb-2 flex items-center">
								<i class="fas fa-database text-blue-500 mr-2"></i> Database options
							</h4>
							<ul class="text-xs text-gray-500 dark:text-gray-400 space-y-1">
								<li>• Default: SQLite (auto-created)</li>
								<li>• PostgreSQL: set <code class="font-mono text-gray-700 dark:text-gray-300">DATABASE_URL</code></li>
							</ul>
						</div>
						<div>
							<h4 class="text-sm font-semibold text-gray-900 dark:text-white mb-2 flex items-center">
								<i class="fas fa-cog text-blue-500 mr-2"></i> Configuration
							</h4>
							<ul class="text-xs text-gray-500 dark:text-gray-400 space-y-1">
								<li>• Config in <code class="font-mono text-gray-700 dark:text-gray-300">./beo-echo-config/</code></li>
								<li>• Persistent data across restarts</li>
							</ul>
						</div>
					</div>
				</div>
			</div>
		</section>

	</main>

	<LandingPageFooter />
</div>

<style>
	/* ── nested hierarchy ── */
	.hierarchy .nest {
		border: 1px solid rgba(148, 163, 184, 0.25);
		border-radius: 12px;
		background-clip: padding-box;
		animation: nestIn 0.5s ease-out both;
		animation-delay: calc(var(--i) * 90ms);
	}
	@keyframes nestIn {
		from {
			opacity: 0;
			transform: translateY(8px) scale(0.98);
		}
		to {
			opacity: 1;
			transform: none;
		}
	}

	/* ── reveal on load ── */
	.reveal {
		animation: reveal 0.6s ease-out both;
		animation-delay: var(--d, 0ms);
	}
	@keyframes reveal {
		from {
			opacity: 0;
			transform: translateY(14px);
		}
		to {
			opacity: 1;
			transform: none;
		}
	}

	/* ── mode request-flow: request in → decide → out ── */
	.flow {
		position: relative;
		height: 92px;
	}
	/* disabled: the user and their app are still live (they keep sending); only
	   the path beyond Beo Echo is dead — handled by the dimmed Real API + rail. */
	.fnodes {
		position: absolute;
		top: 16px;
		left: 0;
		right: 0;
		height: 44px;
	}
	/* Absolutely centred at 8 / 36 / 64 / 92% so the rails and the animated
	   chip (which use the same percentages) line up exactly. */
	.fnode {
		position: absolute;
		transform: translateX(-50%);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 5px;
		width: 70px;
		font-size: 11px;
		color: #64748b;
		z-index: 1;
	}
	.fnode:nth-child(1) {
		left: 8%;
	}
	.fnode:nth-child(2) {
		left: 36%;
	}
	.fnode:nth-child(3) {
		left: 64%;
	}
	.fnode:nth-child(4) {
		left: 92%;
	}
	:global(.dark) .fnode {
		color: #94a3b8;
	}
	.fnode :global(i) {
		font-size: 18px;
	}
	.fnode.beo :global(i) {
		color: #3b82f6;
	}
	.fnode.dim {
		opacity: 0.45;
	}

	/* rails between the nodes */
	.rails {
		position: absolute;
		top: 25px;
		left: 0;
		right: 0;
		height: 2px;
	}
	.rail {
		position: absolute;
		top: 0;
		height: 2px;
		border-radius: 2px;
		background: #cbd5e1;
	}
	:global(.dark) .rail {
		background: #475569;
	}
	.rail0 {
		left: 11%;
		width: 22%;
	}
	.rail1 {
		left: 39%;
		width: 22%;
	}
	.rail2 {
		left: 67%;
		width: 22%;
	}
	.rail2.cut {
		background: repeating-linear-gradient(90deg, #cbd5e1 0 5px, transparent 5px 10px);
		opacity: 0.5;
	}
	:global(.dark) .rail2.cut {
		background: repeating-linear-gradient(90deg, #475569 0 5px, transparent 5px 10px);
	}
	.rail2.cond {
		background: repeating-linear-gradient(90deg, #93c5fd 0 6px, transparent 6px 12px);
	}
	:global(.dark) .rail2.cond {
		background: repeating-linear-gradient(90deg, #3b82f6 0 6px, transparent 6px 12px);
	}
	.cut-x {
		position: absolute;
		top: -9px;
		left: 78%;
		transform: translateX(-50%);
		font-size: 12px;
		color: #ef4444;
	}

	/* moving request "chip" */
	.layer {
		position: absolute;
		inset: 0;
	}
	.chip {
		position: absolute;
		top: 18px;
		width: 20px;
		height: 16px;
		border-radius: 5px;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #fff;
		font-size: 8px;
		background: #3b82f6;
		box-shadow: 0 2px 6px rgba(15, 23, 42, 0.35);
		z-index: 2;
	}
	.chip-match {
		animation: chipMatch 4.2s ease-in-out infinite;
	}
	.chip-fwd {
		animation: chipFwd 4.2s ease-in-out infinite;
	}
	/* disabled: chip arrives from the user but dies at Beo Echo (no response) */
	.chip-dead {
		animation: chipDead 4.2s ease-in-out infinite;
	}
	/* "processing…" tag that blips at FE while the request transits it */
	.fe-proc {
		position: absolute;
		top: 0;
		left: 36%;
		transform: translateX(-50%);
		font-size: 8.5px;
		font-weight: 700;
		white-space: nowrap;
		color: #2563eb;
		opacity: 0;
		animation: feProc 4.2s ease-in-out infinite;
	}
	:global(.dark) .fe-proc {
		color: #60a5fa;
	}
	.fe-proc :global(i) {
		font-size: 8px;
		animation: spin 1.4s linear infinite;
	}
	/* disabled: request dies at Beo, nothing comes back — only the inbound blip */
	.flow.disabled .fe-proc {
		animation-name: feProcIn;
	}
	/* forwarder: no match check — the request passes straight through */
	.flow.forwarder .chip-fwd {
		animation-name: chipThrough;
	}
	/* proxy shows both outcomes, offset so they alternate */
	.flow.proxy .chip-fwd,
	.flow.proxy .v-fwd {
		animation-delay: 2.1s;
	}

	/* verdict label that pops at Beo Echo when the decision is made */
	.verdict {
		position: absolute;
		top: 0;
		left: 64%;
		transform: translateX(-50%);
		font-size: 9px;
		font-weight: 700;
		white-space: nowrap;
		opacity: 0;
		animation: verdictPop 4.2s ease-in-out infinite;
	}
	.v-match {
		color: #16a34a;
	}
	:global(.dark) .v-match {
		color: #4ade80;
	}
	.v-fwd {
		color: #7c3aed;
	}
	:global(.dark) .v-fwd {
		color: #a78bfa;
	}
	.v-dead {
		color: #94a3b8;
	}
	:global(.dark) .v-dead {
		color: #94a3b8;
	}

	/* Shared inbound phase (identical in every mode): user 8% → fe 36% → beo 64%
	   over 0–40%. What happens after 40% is what differs between modes. */
	/* mock: answered at Beo, sent straight back — never reaches Real API */
	@keyframes chipMatch {
		0% {
			left: 8%;
			opacity: 0;
			background: #3b82f6;
		}
		6% {
			opacity: 1;
		}
		18% {
			left: 36%;
			background: #3b82f6;
		}
		40% {
			left: 64%;
			background: #3b82f6;
		}
		48% {
			left: 64%;
			background: #22c55e;
		}
		56% {
			left: 64%;
			background: #22c55e;
		}
		68% {
			left: 36%;
			background: #22c55e;
		}
		73% {
			left: 36%;
			background: #22c55e;
		}
		86% {
			left: 8%;
			opacity: 1;
			background: #22c55e;
		}
		94% {
			left: 8%;
			opacity: 0;
		}
		100% {
			left: 8%;
			opacity: 0;
		}
	}
	/* proxy no-match: in → decide → out to Real API → data returns to client */
	/* proxy no-match: same inbound to Beo, then forwarded on to Real API */
	@keyframes chipFwd {
		0% {
			left: 8%;
			opacity: 0;
			background: #3b82f6;
		}
		6% {
			opacity: 1;
		}
		18% {
			left: 36%;
			background: #3b82f6;
		}
		40% {
			left: 64%;
			background: #3b82f6;
		}
		48% {
			left: 64%;
			background: #8b5cf6;
		}
		62% {
			left: 92%;
			background: #8b5cf6;
		}
		70% {
			left: 92%;
		}
		82% {
			left: 36%;
			background: #8b5cf6;
		}
		86% {
			left: 36%;
			background: #8b5cf6;
		}
		96% {
			left: 8%;
			opacity: 1;
			background: #8b5cf6;
		}
		99% {
			left: 8%;
			opacity: 0;
		}
		100% {
			left: 8%;
			opacity: 0;
		}
	}
	/* forwarder: straight through to Real API, data returns to client — no match check */
	/* forwarder: same inbound to Beo, then straight through to Real API (no decide) */
	@keyframes chipThrough {
		0% {
			left: 8%;
			opacity: 0;
			background: #3b82f6;
		}
		6% {
			opacity: 1;
		}
		18% {
			left: 36%;
			background: #3b82f6;
		}
		40% {
			left: 64%;
			background: #3b82f6;
		}
		62% {
			left: 92%;
			background: #8b5cf6;
		}
		70% {
			left: 92%;
		}
		82% {
			left: 36%;
			background: #8b5cf6;
		}
		86% {
			left: 36%;
			background: #8b5cf6;
		}
		96% {
			left: 8%;
			opacity: 1;
			background: #8b5cf6;
		}
		99% {
			left: 8%;
			opacity: 0;
		}
		100% {
			left: 8%;
			opacity: 0;
		}
	}
	/* disabled: same inbound to Beo, then the request dies there (no response) */
	@keyframes chipDead {
		0% {
			left: 8%;
			opacity: 0;
			background: #3b82f6;
		}
		6% {
			opacity: 1;
		}
		18% {
			left: 36%;
			background: #3b82f6;
		}
		40% {
			left: 64%;
			background: #3b82f6;
		}
		52% {
			left: 64%;
			background: #94a3b8;
		}
		72% {
			left: 64%;
			opacity: 1;
			background: #94a3b8;
		}
		82% {
			left: 64%;
			opacity: 0;
		}
		100% {
			left: 64%;
			opacity: 0;
		}
	}
	/* Two blips: once as the request reaches FE (~18%), once as the response
	   passes back through FE (~72%) — the app processes both directions. */
	@keyframes feProc {
		0%,
		10% {
			opacity: 0;
			transform: translate(-50%, 3px);
		}
		18% {
			opacity: 1;
			transform: translate(-50%, 0);
		}
		30% {
			opacity: 1;
			transform: translate(-50%, 0);
		}
		38% {
			opacity: 0;
			transform: translate(-50%, -2px);
		}
		64% {
			opacity: 0;
			transform: translate(-50%, 3px);
		}
		72% {
			opacity: 1;
			transform: translate(-50%, 0);
		}
		82% {
			opacity: 1;
			transform: translate(-50%, 0);
		}
		90%,
		100% {
			opacity: 0;
			transform: translate(-50%, -2px);
		}
	}
	/* inbound blip only (used by disabled, where nothing comes back) */
	@keyframes feProcIn {
		0%,
		10% {
			opacity: 0;
			transform: translate(-50%, 3px);
		}
		18% {
			opacity: 1;
			transform: translate(-50%, 0);
		}
		30% {
			opacity: 1;
			transform: translate(-50%, 0);
		}
		38%,
		100% {
			opacity: 0;
			transform: translate(-50%, -2px);
		}
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	@keyframes verdictPop {
		0%,
		30% {
			opacity: 0;
			transform: translate(-50%, 4px) scale(0.9);
		}
		42% {
			opacity: 1;
			transform: translate(-50%, 0) scale(1);
		}
		56% {
			opacity: 1;
			transform: translate(-50%, 0) scale(1);
		}
		68%,
		100% {
			opacity: 0;
			transform: translate(-50%, -2px) scale(0.95);
		}
	}

	/* ── migration pipeline ── */
	.pipeline {
		position: relative;
	}
	.track {
		position: absolute;
		top: 34px;
		left: 8%;
		right: 8%;
		height: 3px;
		background: #e2e8f0;
		border-radius: 3px;
		overflow: visible;
	}
	:global(.dark) .track {
		background: #334155;
	}
	/* the coloured progress that fills the path left → right */
	.track-fill {
		position: absolute;
		inset: 0;
		border-radius: 3px;
		background: linear-gradient(90deg, #8b5cf6, #22c55e, #3b82f6, #6366f1);
		transform: scaleX(0);
		transform-origin: left;
		animation: fillGrow 6s ease-in-out infinite;
	}
	.flow-dot {
		position: absolute;
		top: -3px;
		width: 9px;
		height: 9px;
		border-radius: 50%;
		background: #fff;
		box-shadow: 0 0 10px rgba(99, 102, 241, 0.9);
		border: 2px solid #6366f1;
		animation: pipeflow 6s ease-in-out infinite;
	}
	@keyframes pipeflow {
		0% {
			left: 0;
			opacity: 0;
		}
		5% {
			opacity: 1;
		}
		70% {
			left: calc(100% - 9px);
			opacity: 1;
		}
		80% {
			opacity: 0;
		}
		100% {
			left: calc(100% - 9px);
			opacity: 0;
		}
	}
	@keyframes fillGrow {
		0% {
			transform: scaleX(0);
			opacity: 1;
		}
		70% {
			transform: scaleX(1);
			opacity: 1;
		}
		86% {
			transform: scaleX(1);
			opacity: 1;
		}
		93% {
			opacity: 0;
		}
		94% {
			transform: scaleX(0);
			opacity: 0;
		}
		100% {
			transform: scaleX(0);
			opacity: 0;
		}
	}
	.stages {
		position: relative;
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 20px;
	}
	@media (max-width: 720px) {
		.track {
			display: none;
		}
		.stages {
			grid-template-columns: 1fr;
		}
	}
	.stage {
		text-align: center;
		padding: 0 8px;
	}
	.stage-icon {
		width: 54px;
		height: 54px;
		border-radius: 16px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		color: #fff;
		font-size: 20px;
		box-shadow: 0 8px 20px -6px rgba(0, 0, 0, 0.3);
		position: relative;
		z-index: 1;
		/* each stage lights up as the progress wave reaches it */
		animation: stagePulse 6s ease-in-out infinite;
		animation-delay: var(--sd, 0s);
	}
	@keyframes stagePulse {
		0%,
		100% {
			transform: scale(1);
			box-shadow: 0 8px 20px -6px rgba(0, 0, 0, 0.3);
		}
		8% {
			transform: scale(1.16);
			box-shadow:
				0 0 0 8px rgba(99, 102, 241, 0.18),
				0 8px 20px -6px rgba(0, 0, 0, 0.3);
		}
		20% {
			transform: scale(1);
			box-shadow: 0 8px 20px -6px rgba(0, 0, 0, 0.3);
		}
	}
	.stage-mode {
		margin-top: 12px;
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: #6366f1;
	}
	:global(.dark) .stage-mode {
		color: #a5b4fc;
	}
	.stage-title {
		font-size: 17px;
		font-weight: 700;
		color: #111827;
		margin-top: 2px;
	}
	:global(.dark) .stage-title {
		color: #f9fafb;
	}
	.stage-desc {
		font-size: 13px;
		color: #6b7280;
		margin-top: 6px;
		line-height: 1.5;
	}
	:global(.dark) .stage-desc {
		color: #9ca3af;
	}

	@media (prefers-reduced-motion: reduce) {
		.flow-dot {
			animation: none;
			display: none;
		}
		.reveal,
		.hierarchy .nest {
			animation: none;
			opacity: 1;
			transform: none;
		}
		.chip,
		.verdict {
			display: none;
		}
		.stage-icon {
			animation: none;
		}
		.track-fill {
			animation: none;
			transform: scaleX(1);
			opacity: 0.55;
		}
	}
</style>
