<!--
	Illustrative scene that SIMULATES the staged rollout:
	  phase 0 · not ready  → every call is mocked (backend still building)
	  phase 1 · partial    → some endpoints now forwarded to the real API
	  phase 2 · ready       → all calls forwarded to the real API
	Flow: User → FE / Service → Beo Echo → Backend. QA/user keeps testing
	throughout — never blocked. Pure SVG + CSS, phase driven by a small timer,
	reduced-motion friendly.
-->
<script lang="ts">
	import { onMount } from 'svelte';

	let phase = 0; // 0 not ready · 1 partial · 2 ready

	// Advance one phase per full packet round-trip so the badge/state always
	// matches where the animated packet is. Must equal the .pkt animation
	// duration (6s) below.
	onMount(() => {
		const t = setInterval(() => {
			phase = (phase + 1) % 3;
		}, 6000);
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
		viewBox="0 0 940 276"
		role="img"
		aria-label="Simulated rollout: user drives the FE which calls Beo Echo; QA keeps testing on mocks while the backend goes from not ready to partial to fully forwarded to the real API"
	>
		<!-- ───── phase badge ───── -->
		<rect x="378" y="6" width="184" height="28" rx="14" fill={info.color} fill-opacity="0.14" stroke={info.color} stroke-width="1.5" />
		<text x="470" y="24" class="badge-txt" fill={info.color}>{info.badge}</text>
		<text x="470" y="50" class="tiny muted" text-anchor="middle">{info.note}</text>

		<!-- ───── connectors ───── -->
		<!-- user → fe (short link) -->
		<line x1="150" y1="130" x2="196" y2="130" class="rail" />
		<!-- fe ↔ beo (request / response lanes) -->
		<line x1="300" y1="118" x2="494" y2="118" class="rail" />
		<line x1="494" y1="142" x2="300" y2="142" class="rail" />
		<!-- beo → backend (state changes with phase) -->
		<line
			x1="590"
			y1="130"
			x2="716"
			y2="130"
			class="rail"
			class:pending={phase === 0}
			class:partial={phase === 1}
			class:live={phase === 2}
		/>

		<!-- user ↔ fe packets: request out, response back (always on) -->
		<circle class="pkt uf" cx="150" cy="126" r="5" />
		<circle class="pkt fu" cx="196" cy="134" r="5" />
		<!-- fe ↔ beo packets: request out, response back (always on) -->
		<circle class="pkt req" cx="300" cy="118" r="6" />
		<circle class="pkt res" cx="494" cy="142" r="6" />
		<text x="397" y="108" class="lane-tag req-tag">request</text>
		<text x="397" y="160" class="lane-tag res-tag">response</text>

		<!-- beo ↔ backend packet (only once endpoints go live) -->
		{#if phase >= 1}
			<circle class="pkt bk" class:live={phase === 2} cx="590" cy="130" r="5.5" />
		{/if}

		<!-- state marker on the beo→backend link -->
		{#if phase === 0}
			<text x="653" y="122" class="mark">🚧</text>
		{:else if phase === 2}
			<text x="653" y="124" class="mark check">✓</text>
		{/if}

		<!-- ───── 1. USER ───── -->
		<g class="bob">
			<circle cx="90" cy="92" r="16" class="skin" />
			<path d="M62 138 a28 24 0 0 1 56 0 Z" class="user-body" />
		</g>
		<text x="90" y="202" class="label">User</text>
		<text x="90" y="218" class="tiny muted" text-anchor="middle">clicks around, unblocked</text>

		<!-- ───── 2. FE / SERVICE ───── -->
		<g class="bob">
			<rect x="200" y="96" width="96" height="60" rx="6" class="laptop-screen" />
			<rect x="210" y="106" width="52" height="6" rx="2" class="code-accent" />
			<rect x="210" y="118" width="70" height="6" rx="2" class="code-dim" />
			<rect x="210" y="130" width="40" height="6" rx="2" class="code-dim" />
			<circle cx="286" cy="147" r="3" class="online" />
			<path d="M192 156 L304 156 L314 168 L182 168 Z" class="laptop-base" />
		</g>
		<text x="248" y="202" class="label">Your app</text>
		<text x="248" y="218" class="tiny muted" text-anchor="middle">FE / service / client</text>

		<!-- ───── 3. BEO ECHO ───── -->
		<rect x="496" y="80" width="92" height="86" rx="14" class="card" />
		<path d="M546 98 L528 128 L540 128 L534 152 L556 120 L543 120 Z" class="bolt" />
		<rect x="510" y="140" width="64" height="7" rx="3" class="mock-row" />
		<rect x="510" y="151" width="46" height="7" rx="3" class="mock-row" />
		<text x="542" y="202" class="label">Beo Echo</text>
		<text x="542" y="218" class="tiny muted" text-anchor="middle">routes every call</text>

		<!-- ───── 4. BACKEND ───── -->
		<rect x="716" y="74" width="118" height="74" rx="7" class="monitor" />
		<rect x="724" y="84" width="52" height="7" rx="2" class="code-ln l1" />
		<rect x="724" y="96" width="86" height="7" rx="2" class="code-ln l2" />
		<rect x="724" y="108" width="38" height="7" rx="2" class="code-ln l3" />
		<rect x="724" y="120" width="70" height="7" rx="2" class="code-ln l4" />
		<rect x="798" y="120" width="6" height="7" class="caret" />
		<rect x="769" y="148" width="12" height="10" class="stand" />
		<rect x="750" y="158" width="50" height="6" rx="3" class="stand" />

		<!-- work-in-progress gear (fades out once ready) -->
		<g class="gear" class:done={phase === 2} transform="translate(822,80)">
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
			<circle cx="850" cy="164" r="13" class="skin" />
			<path d="M836 164 a14 14 0 0 1 28 0 Z" class="hardhat" />
			<rect x="833" y="163" width="34" height="4" rx="2" class="hardhat" />
		</g>
		<text x="775" y="202" class="label">Backend</text>
		<text x="775" y="218" class="tiny muted" text-anchor="middle">{backendNote}</text>

		<!-- rollout progress bar -->
		<rect x="736" y="232" width="80" height="7" rx="3.5" class="bar-track" />
		<rect
			x="736"
			y="232"
			width="80"
			height="7"
			rx="3.5"
			class="bar-fill"
			style="transform: scaleX({barScale});"
			fill={info.color}
		/>
		<text x="776" y="256" class="tiny muted" text-anchor="middle">{barCaption}</text>
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
		max-width: 780px;
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
		animation-duration: 6s;
		animation-timing-function: linear;
		animation-iteration-count: infinite;
	}
	.pkt.uf {
		fill: var(--accent);
		animation-name: ufMove;
	}
	.pkt.fu {
		fill: var(--green);
		animation-name: fuMove;
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
		animation: bkMove 6s linear infinite;
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

	/* Relay across 3.2s: request hops user→fe→beo, then response hops back
	   beo→fe→user, so the frame is forwarded and returned end-to-end. */
	/* One frame relayed end-to-end across a single 3.2s cycle, hop by hop:
	   user→fe→beo→backend, then backend→beo→fe→user. All lanes share the same
	   duration so the packet is continuous, never overlapping or racing. */
	/* 1 · user → fe (request) — 0–14% */
	@keyframes ufMove {
		0% {
			transform: translateX(0);
			opacity: 0;
		}
		3% {
			opacity: 1;
		}
		14% {
			transform: translateX(46px);
			opacity: 1;
		}
		18% {
			transform: translateX(46px);
			opacity: 0;
		}
		100% {
			transform: translateX(46px);
			opacity: 0;
		}
	}
	/* 6 · fe → user (response) — 81–96% */
	@keyframes fuMove {
		0%,
		81% {
			transform: translateX(0);
			opacity: 0;
		}
		84% {
			opacity: 1;
		}
		96% {
			transform: translateX(-46px);
			opacity: 1;
		}
		100% {
			transform: translateX(-46px);
			opacity: 0;
		}
	}
	/* 2 · fe → beo (request) — 14–32% */
	@keyframes reqMove {
		0%,
		14% {
			transform: translateX(0);
			opacity: 0;
		}
		17% {
			opacity: 1;
		}
		32% {
			transform: translateX(194px);
			opacity: 1;
		}
		36% {
			transform: translateX(194px);
			opacity: 0;
		}
		100% {
			transform: translateX(194px);
			opacity: 0;
		}
	}
	/* 5 · beo → fe (response) — 63–81% */
	@keyframes resMove {
		0%,
		63% {
			transform: translateX(0);
			opacity: 0;
		}
		66% {
			opacity: 1;
		}
		81% {
			transform: translateX(-194px);
			opacity: 1;
		}
		85% {
			transform: translateX(-194px);
			opacity: 0;
		}
		100% {
			transform: translateX(-194px);
			opacity: 0;
		}
	}
	@keyframes bkMove {
		0%,
		32% {
			transform: translateX(0);
			opacity: 0;
		}
		35% {
			opacity: 1;
		}
		47% {
			transform: translateX(126px);
			opacity: 1;
		}
		51% {
			transform: translateX(126px);
			opacity: 1;
		}
		63% {
			transform: translateX(0);
			opacity: 1;
		}
		67% {
			transform: translateX(0);
			opacity: 0;
		}
		100% {
			transform: translateX(0);
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
		.pkt.res,
		.pkt.fu {
			opacity: 1;
		}
		.code-ln {
			transform: scaleX(1);
		}
	}
</style>
