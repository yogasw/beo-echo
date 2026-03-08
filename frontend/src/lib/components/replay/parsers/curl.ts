import type { Replay } from '$lib/types/Replay';

// Basic cURL parser to extract URL, Method, Headers, and Body
export function parseCurl(curlText: string): Partial<Replay> {
    const lines = curlText.match(/(?:[^\s']+|'[^']*')+/g) || [];
    if (lines.length === 0 || !lines[0]?.toLowerCase().includes("curl")) {
        throw new Error("Invalid cURL command");
    }

    let url = "";
    let method = "GET";
    let headers: Array<{key: string, value: string}> = [];
    let body = "";

    let skipNext = false;

    for (let i = 1; i < lines.length; i++) {
        let token = lines[i];
		
		if (token === '\\') continue;
		
		if (skipNext) {
			skipNext = false;
			continue;
		}

		// remove surrounding quotes if exist
		if ((token.startsWith("'") && token.endsWith("'")) || (token.startsWith('"') && token.endsWith('"'))) {
			token = token.substring(1, token.length - 1);
		}

        if (token === "-X" || token === "--request") {
            method = lines[i+1]?.replace(/['"]/g, '').toUpperCase() || method;
			skipNext = true;
        } else if (token === "-H" || token === "--header") {
            const headerStr = lines[i+1]?.replace(/['"]/g, '') || "";
            const separatorIdx = headerStr.indexOf(':');
            if (separatorIdx > -1) {
                headers.push({
                    key: headerStr.substring(0, separatorIdx).trim(),
                    value: headerStr.substring(separatorIdx + 1).trim()
                });
            }
			skipNext = true;
        } else if (token === "-d" || token === "--data" || token === "--data-raw" || token === "--data-binary" || token === "--data-urlencode") {
            body = lines[i+1] || "";
			if ((body.startsWith("'") && body.endsWith("'")) || (body.startsWith('"') && body.endsWith('"'))) {
				body = body.substring(1, body.length - 1);
			}
            if (method === "GET") method = "POST";
			skipNext = true;
        } else if (!token.startsWith('-') && !url) {
            url = token;
        } else if (token.startsWith('-')) {
			// Skip next token if it's a known flag that takes an argument but we don't handle it
			if (["-u", "--user", "-A", "--user-agent", "-e", "--referer", "-m", "--max-time", "-w", "--write-out", "-x", "--proxy", "-b", "--cookie", "-c", "--cookie-jar", "-F", "--form"].includes(token)) {
				skipNext = true;
			}
		}
    }

	if (!url) {
		throw new Error("Could not find a valid command or URL in the cURL structure.");
	}

	const headersObj: any = {};
	headers.forEach(h => {
		headersObj[h.key] = h.value;
	});

    return {
        name: "Imported Request",
        url,
        method,
        headers: JSON.stringify(headersObj),
        payload: body,
		config: JSON.stringify({ auth: { type: 'none', config: {} }, settings: {} })
    };
}

export const curlParser = {
	parse: (text: string): { parsed: Partial<Replay>; importType: 'curl', rawText: string } => {
		return {
			parsed: parseCurl(text),
			importType: 'curl',
			rawText: text
		};
	}
};
