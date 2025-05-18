export type ConfigBeoEcho = {
	uuid: string;
	lastMigration: number;
	name: string;
	endpointPrefix: string;
	latency: number;
	port: number;
	hostname: string;
	folders: any[];
	routes: BeoEchoRoute[];
	rootChildren: BeoEchoRootChild[];
	proxyMode: boolean;
	proxyHost: string;
	proxyRemovePrefix: boolean;
	tlsOptions: BeoEchoTLSOptions;
	cors: boolean;
	headers: BeoEchoHeader[];
	proxyReqHeaders: BeoEchoHeader[];
	proxyResHeaders: BeoEchoHeader[];
	data: any[];
	callbacks: BeoEchoCallback[];
}

export type BeoEchoCallback = {
	uuid: string;
	id: string;
	uri: string;
	name: string;
	documentation: string;
	method: string;
	headers: any[];
	bodyType: string;
	body: string;
	databucketID: string;
	filePath: string;
	sendFileAsBody: boolean;
}

export type BeoEchoHeader = {
	key: string;
	value: string;
}

export type BeoEchoRootChild = {
	type: string;
	uuid: string;
}

export type BeoEchoRoute = {
	uuid: string;
	type: string;
	documentation: string;
	method: string;
	endpoint: string;
	responses: BeoEchoResponse[];
	responseMode: null;
	//additional properties
	status: string
}

export type BeoEchoResponse = {
	uuid: string;
	body: string;
	latency: number;
	statusCode: number;
	label: string;
	headers: BeoEchoHeader[];
	bodyType: string;
	filePath: string;
	databucketID: string;
	sendFileAsBody: boolean;
	rules: BeoEchoRule[];
	rulesOperator: string;
	disableTemplating: boolean;
	fallbackTo404: boolean;
	default: boolean;
	crudKey: string;
	callbacks: BeoEchoResponseCallback[];
}

export type BeoEchoResponseCallback = {
	uuid: string;
	latency: number;
}

export type BeoEchoRule = {
	target: string;
	modifier: string;
	value: string;
	invert: boolean;
	operator: string;
}

export type BeoEchoTLSOptions = {
	enabled: boolean;
	type: string;
	pfxPath: string;
	certPath: string;
	keyPath: string;
	caPath: string;
	passphrase: string;
}
