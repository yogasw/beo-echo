package config

const CADDY_DEFAULT_CONFIG = `{
	# Global Caddy settings
	admin off
}

:80 {
	# Log requests
	log {
		# output file /app/logs/caddy-access.log
		format console
	}

	# Define routes for static assets
	@static {
		path /_app/* /robots.txt /images/* /css/* /js/* /fonts/* /favicon.ico /favicon.png /robots.txt
	}

	# Handle static asset routes directly
	handle @static {
		root * /app/frontend
		file_server
	}

	# Define routes for frontend - root, /home, and /login paths
	@frontend {
		path / /home /login
	}

	# Handle frontend routes - serve the SvelteKit static files
	handle @frontend {
		root * /app/frontend
		try_files {path} /index.html
		file_server
	}

	# All other routes go to the backend
	handle {
		# Reverse proxy to backend
		reverse_proxy localhost:3600
	}

	# Handle errors
	handle_errors {
		respond "{err.status_code} {err.status_text}"
	}
}
`
