export type ConfigMockoon = {
	uuid: string;
	lastMigration: number;
	name: string;
	endpointPrefix: string;
	latency: number;
	port: number;
	hostname: string;
	folders: any[];
	routes: MockoonRoute[];
	rootChildren: MockoonRootChild[];
	proxyMode: boolean;
	proxyHost: string;
	proxyRemovePrefix: boolean;
	tlsOptions: MockoonTLSOptions;
	cors: boolean;
	headers: MockoonHeader[];
	proxyReqHeaders: MockoonHeader[];
	proxyResHeaders: MockoonHeader[];
	data: any[];
	callbacks: MockoonCallback[];
}

export type MockoonCallback = {
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

export type MockoonHeader = {
	key: string;
	value: string;
}

export type MockoonRootChild = {
	type: string;
	uuid: string;
}

export type MockoonRoute = {
	uuid: string;
	type: string;
	documentation: string;
	method: string;
	endpoint: string;
	responses: MockoonResponse[];
	responseMode: null;
	//additional properties
	status: string
}

export type MockoonResponse = {
	uuid: string;
	body: string;
	latency: number;
	statusCode: number;
	label: string;
	headers: MockoonHeader[];
	bodyType: string;
	filePath: string;
	databucketID: string;
	sendFileAsBody: boolean;
	rules: MockoonRule[];
	rulesOperator: string;
	disableTemplating: boolean;
	fallbackTo404: boolean;
	default: boolean;
	crudKey: string;
	callbacks: MockoonResponseCallback[];
}

export type MockoonResponseCallback = {
	uuid: string;
	latency: number;
}

export type MockoonRule = {
	target: string;
	modifier: string;
	value: string;
	invert: boolean;
	operator: string;
}

export type MockoonTLSOptions = {
	enabled: boolean;
	type: string;
	pfxPath: string;
	certPath: string;
	keyPath: string;
	caPath: string;
	passphrase: string;
}
