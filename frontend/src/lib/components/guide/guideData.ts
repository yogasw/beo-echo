// Content data for the Beo Echo concept & usage guide page.
// Kept separate from the view so copy can be edited without touching markup.

export interface Concept {
	key: string;
	title: string;
	icon: string;
	color: string; // tailwind text color
	bg: string; // tailwind bg tint
	short: string;
	detail: string;
	example: string;
}

export interface Mode {
	id: 'mock' | 'proxy' | 'forwarder' | 'disabled';
	name: string;
	tagline: string;
	icon: string;
	color: string; // solid bg
	accent: string; // text accent
	behaviour: string;
	whenToUse: string;
}

export interface InstallMethod {
	id: string;
	title: string;
	subtitle: string;
	icon: string;
	accent: string;
	steps: string[];
	code?: string[];
}

export interface UseCase {
	title: string;
	icon: string;
	color: string;
	scenario: string;
	steps: string[];
}

export interface MigrationStage {
	mode: string;
	title: string;
	desc: string;
	icon: string;
	color: string;
}

// ── The resource hierarchy: Workspace → Project → Endpoint → Response → Rule ──
export const concepts: Concept[] = [
	{
		key: 'workspace',
		title: 'Workspace',
		icon: 'fas fa-users',
		color: 'text-indigo-500 dark:text-indigo-300',
		bg: 'bg-indigo-500/10',
		short: 'The team boundary.',
		detail:
			'The top-level container. Members, roles, and projects all live inside a workspace, so access is isolated per team.',
		example: 'e.g. one workspace per team — projects and members stay separate.'
	},
	{
		key: 'project',
		title: 'Project',
		icon: 'fas fa-box',
		color: 'text-blue-500 dark:text-blue-300',
		bg: 'bg-blue-500/10',
		short: 'One mock server.',
		detail:
			'A single mock server with its own alias (URL) and one operating mode. This is the thing you point your app at.',
		example: 'e.g. "payments-sandbox" → https://host/payments-sandbox'
	},
	{
		key: 'endpoint',
		title: 'Endpoint',
		icon: 'fas fa-route',
		color: 'text-green-500 dark:text-green-300',
		bg: 'bg-green-500/10',
		short: 'A route: method + path.',
		detail:
			'Matches an incoming request by HTTP method and path. Path params like :id are supported for flexible matching.',
		example: 'e.g. POST /v1/users/:id/orders'
	},
	{
		key: 'response',
		title: 'Response',
		icon: 'fas fa-reply',
		color: 'text-orange-500 dark:text-orange-300',
		bg: 'bg-orange-500/10',
		short: 'What gets returned.',
		detail:
			'Status code, body, and headers. An endpoint can hold several responses; mark one as the fallback so it always matches.',
		example: 'e.g. 200 with { "code": 200, "message": "Success" }'
	},
	{
		key: 'rule',
		title: 'Rule',
		icon: 'fas fa-filter',
		color: 'text-purple-500 dark:text-purple-300',
		bg: 'bg-purple-500/10',
		short: 'Conditional matching.',
		detail:
			'Attach rules to a response to pick it only when the request matches — by header, query, body, or path. Combine with AND / OR.',
		example: 'e.g. return this response only when body plan = "pro"'
	}
];

export const modes: Mode[] = [
	{
		id: 'mock',
		name: 'Mock',
		tagline: 'Serve your responses only',
		icon: 'fas fa-server',
		color: 'bg-green-600',
		accent: 'text-green-500 dark:text-green-400',
		behaviour:
			'Every request is answered from your configured responses — and you pick exactly which one comes out via rules, priority, or a fallback. The real server is never contacted.',
		whenToUse: 'Backend is down or not built yet, or you want fully isolated, deterministic tests.'
	},
	{
		id: 'proxy',
		name: 'Proxy',
		tagline: 'Selective — matched mock stops, rest forwards',
		icon: 'fas fa-code-branch',
		color: 'bg-blue-600',
		accent: 'text-blue-500 dark:text-blue-400',
		behaviour:
			'Per request, Beo Echo checks your mocks first. If an endpoint matches AND its rules match, that mock answers and the request stops there. Anything without a matching mock is forwarded to the real target.',
		whenToUse: 'Partial mocks — stub specific endpoints (or just specific rule cases); let everything else hit the real API.'
	},
	{
		id: 'forwarder',
		name: 'Forwarder',
		tagline: 'Always forward, always log',
		icon: 'fas fa-arrow-right',
		color: 'bg-purple-600',
		accent: 'text-purple-500 dark:text-purple-400',
		behaviour:
			'Every request goes straight to the real target — mocks are ignored entirely — while the full traffic is recorded.',
		whenToUse: 'Capture and inspect real request/response pairs before you decide what to mock.'
	},
	{
		id: 'disabled',
		name: 'Disabled',
		tagline: 'Inactive',
		icon: 'fas fa-ban',
		color: 'bg-gray-500',
		accent: 'text-gray-500 dark:text-gray-400',
		behaviour: 'The project serves nothing. Useful to temporarily park a project.',
		whenToUse: 'You want to keep the config but stop it responding.'
	}
];

export const installMethods: InstallMethod[] = [
	{
		id: 'claude-code',
		title: 'Claude Code (CLI)',
		subtitle: 'Register once, log in over OAuth',
		icon: 'fas fa-terminal',
		accent: 'text-green-500 dark:text-green-400',
		steps: [
			'Register the streamable-HTTP server.',
			'Log in — a browser opens, then click Approve.',
			'Confirm it is connected.'
		],
		code: [
			'claude mcp add --transport http beo-echo {ENDPOINT}',
			'claude mcp login beo-echo',
			'claude mcp list'
		]
	},
	{
		id: 'claude-ai',
		title: 'Claude.ai',
		subtitle: 'OAuth — paste the URL, no token',
		icon: 'fas fa-shield-halved',
		accent: 'text-purple-500 dark:text-purple-400',
		steps: [
			'Settings → Integrations → Add custom MCP server.',
			'Paste the endpoint URL.',
			'Log in if prompted, then Approve.'
		]
	},
	{
		id: 'static',
		title: 'Desktop · Cursor · VS Code',
		subtitle: 'Static bearer token',
		icon: 'fas fa-key',
		accent: 'text-blue-500 dark:text-blue-400',
		steps: [
			'Profile → MCP → Access Tokens → Generate (shown once).',
			'Drop the endpoint + token into the client config.',
			'Reload the client.'
		],
		code: [
			'{',
			'  "mcpServers": {',
			'    "beo-echo": {',
			'      "url": "{ENDPOINT}",',
			'      "headers": { "Authorization": "Bearer beo_pat_..." }',
			'    }',
			'  }',
			'}'
		]
	}
];

export const useCases: UseCase[] = [
	{
		title: 'Scan an existing project',
		icon: 'fas fa-magnifying-glass',
		color: 'text-blue-500 dark:text-blue-400',
		scenario: 'You inherited a project and want the full picture fast.',
		steps: [
			'"List my projects in this workspace."',
			'"Scan project X — show every endpoint and its responses."',
			'"Which endpoints are still forwarding vs mocked?"'
		]
	},
	{
		title: 'Capture real traffic (Forwarder)',
		icon: 'fas fa-satellite-dish',
		color: 'text-purple-500 dark:text-purple-400',
		scenario: 'You need real payloads before writing any mock.',
		steps: [
			'Set the project to Forwarder mode, pointing at the real API.',
			'Route your app through it and exercise the flow.',
			'"Show the last 20 logs" — then copy the real responses into mocks.'
		]
	},
	{
		title: 'Stub a flaky endpoint (Proxy)',
		icon: 'fas fa-puzzle-piece',
		color: 'text-green-500 dark:text-green-400',
		scenario: 'One endpoint is unreliable; the rest work fine.',
		steps: [
			'Switch the project to Proxy mode.',
			'"Add a mock for POST /payments returning 200 success."',
			'Everything else keeps hitting the real API automatically.'
		]
	},
	{
		title: 'Test with the backend down (Mock)',
		icon: 'fas fa-plug-circle-xmark',
		color: 'text-orange-500 dark:text-orange-400',
		scenario: 'Staging is dead but you must demo the full flow.',
		steps: [
			'Set the project to Mock mode.',
			'"Mock the whole flow from login to checkout, all 200."',
			'Point BASE_URL at the project — every call succeeds.'
		]
	}
];

// The staged path: capture reality → mock the gaps → shrink mocks → go live.
export const migrationStages: MigrationStage[] = [
	{
		mode: 'Forwarder',
		title: 'Observe',
		desc: 'Route through Beo Echo to the real/staging API and record real traffic.',
		icon: 'fas fa-arrow-right',
		color: 'bg-purple-600'
	},
	{
		mode: 'Mock',
		title: 'Build',
		desc: 'Recreate the endpoints you need as mocks so you can develop in isolation.',
		icon: 'fas fa-server',
		color: 'bg-green-600'
	},
	{
		mode: 'Proxy',
		title: 'Blend',
		desc: 'Flip to Proxy: keep mocking what is not ready, forward everything that is.',
		icon: 'fas fa-exchange-alt',
		color: 'bg-blue-600'
	},
	{
		mode: 'Production',
		title: 'Cut over',
		desc: 'As the real API stabilises, drop mocks and point the app straight at production.',
		icon: 'fas fa-rocket',
		color: 'bg-indigo-600'
	}
];
