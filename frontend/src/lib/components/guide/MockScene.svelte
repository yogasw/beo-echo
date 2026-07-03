<!--
	Illustrative scene that SIMULATES the staged rollout:
	  phase 0 · not ready  → every call is mocked (backend still building)
	  phase 1 · partial    → some endpoints now forwarded to the real API
	  phase 2 · ready       → all calls forwarded to the real API
	QA/user keeps testing throughout — never blocked.
	Pure SVG + CSS, phase driven by a small timer, reduced-motion friendly.
-->
<script lang="ts">
	import { onMount } from 'svelte';

	let phase = 0; // 0 not ready · 1 partial · 2 ready

	onMount(() => {
		const t = setInterval(() => {
			phase = (phase + 1) % 3;
		}, 3400);
		return () => clearInterval(t);
	});

	const phases = [
		{ badge: 'Backend: not ready', note: 'every call is served from mocks', color: '#ef4444' },
		{ badge: 'Partial rollout', note: 'some endpoints now hit the real API', color: '#3b82f6' },
		{ badge: 'Real API ready', note: 'all calls forwarded to the real API', color: '#22c55e' }
	];
	$: info = phases[phase];
	$: barScale = [0.28, 0.68, 1][phase];
	$: backendNote = ['still coding the real API', 'rolling out endpoints', 'real API is live'][phase];
	$: barCaption = ['building…', 'rolling out…', 'shipped ✓'][phase];
</script>

<div class="scene">
	<svg
		viewBox="0 0 760 276"
		role="img"
		aria-label="Simulated rollout: QA keeps testing on mocks while the backend goes from not ready to partial to fully forwarded to the real API"
	>
		<!-- ───── phase badge ───── -->
		<rect x="288" y="6" width="184" height="28" rx="14" fill={info.color} fill-opacity="0.14" stroke={info.color} stroke-width="1.5" />
		<text x="380" y="24" class="badge-txt" fill={info.color}>{info.badge}</text>
		<text x="380" y="50" class="tiny muted" text-anchor="middle">{info.note}</text>

		<!-- ───── connectors ───── -->
		<line x1="150" y1="118" x2="344" y2="118" class="rail" />
		<line x1="344" y1="142" x2="150" y2="142" class="rail" />
		<line
			x1="440"
			y1="130"
			x2="566"
			y2="130"
			class="rail"
			class:pending={phase === 0}
			class:partial={phase === 1}
			class:live={phase === 2}
		/>

		<!-- user ↔ beo packets (always on: QA keeps testing) -->
		<circle class="pkt req" cx="150" cy="118" r="6" />
		<circle class="pkt res" cx="344" cy="142" r="6" />
		<text x="247" y="108" class="lane-tag req-tag">request</text>
		<text x="247" y="160" class="lane-tag res-tag">response</text>

		<!-- beo ↔ backend packet (only once endpoints go live) -->
		{#if phase >= 1}
			<circle class="pkt bk" class:live={phase === 2} cx="440" cy="130" r="5.5" />
		{/if}

		<!-- state marker on the beo→backend link -->
		{#if phase === 0}
			<text x="503" y="122" class="mark">🚧</text>
		{:else if phase === 2}
			<text x="503" y="124" class="mark check">✓</text>
		{/if}

		<!-- ───── 1. USER / QA ───── -->
		<g class="bob">
			<circle cx="90" cy="84" r="15" class="skin" />
			<rect x="64" y="102" width="52" height="34" rx="17" class="user-body" />
			<rect x="60" y="132" width="60" height="34" rx="3" class="laptop-screen" />
			<rect x="70" y="140" width="40" height="4" rx="2" class="code-accent" />
			<rect x="70" y="148" width="28" height="4" rx="2" class="code-dim" />
			<circle cx="112" cy="138" r="2.5" class="online" />
			<path d="M52 166 L128 166 L136 176 L44 176 Z" class="laptop-base" />
		</g>
		<text x="90" y="202" class="label">You / QA</text>
		<text x="90" y="218" class="tiny muted" text-anchor="middle">keep testing, unblocked</text>

		<!-- ───── 2. BEO ECHO ───── -->
		<rect x="346" y="80" width="92" height="86" rx="14" class="card" />
		<path d="M396 98 L378 128 L390 128 L384 152 L406 120 L393 120 Z" class="bolt" />
		<rect x="360" y="140" width="64" height="7" rx="3" class="mock-row" />
		<rect x="360" y="151" width="46" height="7" rx="3" class="mock-row" />
		<text x="392" y="202" class="label">Beo Echo</text>
		<text x="392" y="218" class="tiny muted" text-anchor="middle">routes every call</text>

		<!-- ───── 3. BACKEND ───── -->
		<rect x="566" y="74" width="118" height="74" rx="7" class="monitor" />
		<rect x="574" y="84" width="52" height="7" rx="2" class="code-ln l1" />
		<rect x="574" y="96" width="86" height="7" rx="2" class="code-ln l2" />
		<rect x="574" y="108" width="38" height="7" rx="2" class="code-ln l3" />
		<rect x="574" y="120" width="70" height="7" rx="2" class="code-ln l4" />
		<rect x="648" y="120" width="6" height="7" class="caret" />
		<rect x="619" y="148" width="12" height="10" class="stand" />
		<rect x="600" y="158" width="50" height="6" rx="3" class="stand" />

		<!-- work-in-progress gear (fades out once ready) -->
		<g class="gear" class:done={phase === 2} transform="translate(672,80)">
			<circle r="9" class="gear-body" />
			<circle r="3.5" class="gear-hole" />
			<g class="gear-teeth">
				<rect x="-2" y="-13" width="4" height="5" />
				<rect x="-2" y="8" width="4" height="5" />
				<rect x="-13" y="-2" width="5" height="4" />
				<rect x="8" y="-2" width="5" height="4" />
			</g>
		</g>

		<g class="bob slow">
			<circle cx="700" cy="164" r="13" class="skin" />
			<path d="M686 164 a14 14 0 0 1 28 0 Z" class="hardhat" />
			<rect x="683" y="163" width="34" height="4" rx="2" class="hardhat" />
		</g>
		<text x="625" y="202" class="label">Backend</text>
		<text x="625" y="218" class="tiny muted" text-anchor="middle">{backendNote}</text>

		<!-- rollout progress bar -->
		<rect x="586" y="232" width="80" height="7" rx="3.5" class="bar-track" />
		<rect
			x="586"
			y="232"
			width="80"
			height="7"
			rx="3.5"
			class="bar-fill"
			style="transform: scaleX({barScale});"
			fill={info.color}
		/>
		<text x="626" y="256" class="tiny muted" text-anchor="middle">{barCaption}</text>
	</svg>
</div>

<style>
	.scene {
		--ink: #334155;
		--muted: #64748b;
		--card: #ffffff;
		--line: #cbd5e1;
		--stroke: #cbd5e1;
		--skin: #93c5fd;
		--accent: #3b82f6;
		--green: #22c55e;
		--amber: #f59e0b;
		--screen: #0f172a;
		--screen-brd: #334155;
		width: 100%;
		max-width: 720px;
		margin: 0 auto;
	}
	:global(.dark) .scene {
		--ink: #e2e8f0;
		--muted: #94a3b8;
		--card: #1e293b;
		--line: #475569;
		--stroke: #475569;
		--skin: #60a5fa;
		--screen: #0b1220;
		--screen-brd: #1e293b;
	}
	svg {
		width: 100%;
		height: auto;
		font-family: 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
	}

	.rail {
		stroke: var(--line);
		stroke-width: 2;
		transition:
			stroke 0.5s ease,
			opacity 0.5s ease;
	}
	.rail.pending {
		stroke: var(--line);
		stroke-dasharray: 5 5;
		opacity: 0.55;
	}
	.rail.partial {
		stroke: var(--accent);
		stroke-dasharray: 5 5;
		opacity: 0.95;
	}
	.rail.live {
		stroke: var(--green);
		opacity: 1;
	}

	.badge-txt {
		font-size: 13px;
		font-weight: 700;
		text-anchor: middle;
	}
	.label {
		fill: var(--ink);
		font-size: 14px;
		font-weight: 700;
		text-anchor: middle;
	}
	.tiny {
		font-size: 10.5px;
	}
	.muted {
		fill: var(--muted);
	}
	.lane-tag {
		font-size: 10px;
		font-weight: 600;
		text-anchor: middle;
	}
	.req-tag {
		fill: var(--accent);
	}
	.res-tag {
		fill: var(--green);
	}
	.mark {
		font-size: 16px;
		text-anchor: middle;
	}
	.mark.check {
		fill: var(--green);
		font-size: 18px;
		font-weight: 800;
	}

	.skin {
		fill: var(--skin);
	}
	.user-body {
		fill: var(--accent);
	}
	.laptop-screen {
		fill: var(--screen);
		stroke: var(--screen-brd);
		stroke-width: 1.5;
	}
	.laptop-base {
		fill: #94a3b8;
	}
	.code-accent {
		fill: #38bdf8;
	}
	.code-dim {
		fill: #475569;
	}
	.online {
		fill: var(--green);
	}
	.card {
		fill: var(--card);
		stroke: var(--stroke);
		stroke-width: 1.5;
	}
	.bolt {
		fill: var(--accent);
	}
	.mock-row {
		fill: var(--line);
	}
	.monitor {
		fill: var(--screen);
		stroke: var(--screen-brd);
		stroke-width: 1.5;
	}
	.stand {
		fill: #94a3b8;
	}
	.hardhat {
		fill: var(--amber);
	}
	.caret {
		fill: #38bdf8;
		animation: blink 1s steps(1) infinite;
	}
	.gear-body {
		fill: var(--amber);
	}
	.gear-hole {
		fill: var(--card);
	}
	.gear-teeth :global(rect) {
		fill: var(--amber);
	}

	.bar-track {
		fill: var(--line);
	}
	.bar-fill {
		transform-box: fill-box;
		transform-origin: left center;
		transition:
			transform 0.7s ease,
			fill 0.5s ease;
	}

	.code-ln {
		fill: #64748b;
		transform-box: fill-box;
		transform-origin: left center;
		animation: typeline 2.6s ease-in-out infinite;
	}
	.code-ln.l1 {
		fill: #a78bfa;
	}
	.code-ln.l2 {
		fill: #38bdf8;
		animation-delay: 0.5s;
	}
	.code-ln.l3 {
		fill: #f472b6;
		animation-delay: 1s;
	}
	.code-ln.l4 {
		fill: #64748b;
		animation-delay: 1.5s;
	}

	.pkt {
		animation-duration: 3.2s;
		animation-timing-function: ease-in-out;
		animation-iteration-count: infinite;
	}
	.pkt.req {
		fill: var(--accent);
		animation-name: reqMove;
	}
	.pkt.res {
		fill: var(--green);
		animation-name: resMove;
	}
	.pkt.bk {
		fill: var(--accent);
		animation: bkMove 2.4s ease-in-out infinite;
	}
	.pkt.bk.live {
		fill: var(--green);
	}

	.bob {
		transform-box: fill-box;
		transform-origin: center;
		animation: bob 3.5s ease-in-out infinite;
	}
	.bob.slow {
		animation-duration: 4.5s;
	}
	.gear {
		transform-box: fill-box;
		transform-origin: center;
		animation: spin 3s linear infinite;
		transition: opacity 0.5s ease;
	}
	.gear.done {
		opacity: 0.25;
		animation-play-state: paused;
	}

	@keyframes reqMove {
		0% {
			transform: translateX(0);
			opacity: 0;
		}
		10% {
			opacity: 1;
		}
		42% {
			transform: translateX(194px);
			opacity: 1;
		}
		50% {
			opacity: 0;
		}
		100% {
			transform: translateX(194px);
			opacity: 0;
		}
	}
	@keyframes resMove {
		0%,
		48% {
			transform: translateX(0);
			opacity: 0;
		}
		54% {
			opacity: 1;
		}
		90% {
			transform: translateX(-194px);
			opacity: 1;
		}
		100% {
			transform: translateX(-194px);
			opacity: 0;
		}
	}
	@keyframes bkMove {
		0% {
			transform: translateX(0);
			opacity: 0;
		}
		12% {
			opacity: 1;
		}
		45% {
			transform: translateX(126px);
		}
		55% {
			transform: translateX(126px);
		}
		88% {
			transform: translateX(0);
			opacity: 1;
		}
		100% {
			opacity: 0;
		}
	}
	@keyframes typeline {
		0% {
			transform: scaleX(0.15);
		}
		45% {
			transform: scaleX(1);
		}
		75% {
			transform: scaleX(1);
		}
		100% {
			transform: scaleX(0.15);
		}
	}
	@keyframes blink {
		0%,
		50% {
			opacity: 1;
		}
		51%,
		100% {
			opacity: 0;
		}
	}
	@keyframes bob {
		0%,
		100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-4px);
		}
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.pkt,
		.bob,
		.gear,
		.code-ln,
		.caret {
			animation: none;
		}
		.pkt.res {
			opacity: 1;
		}
		.code-ln {
			transform: scaleX(1);
		}
	}
</style>
