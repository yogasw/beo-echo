package traefik

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/prisma"
	"mockoon-control-panel/backend_new/src/utils"
)

// TraefikDynamicConfig represents the structure of the Traefik dynamic configuration
type TraefikDynamicConfig struct {
	HTTP HTTP `yaml:"http"`
}

// HTTP represents the http section of the Traefik dynamic configuration
type HTTP struct {
	Routers     map[string]Router     `yaml:"routers"`
	Services    map[string]Service    `yaml:"services"`
	Middlewares map[string]Middleware `yaml:"middlewares,omitempty"`
}

// Router represents a Traefik router
type Router struct {
	EntryPoints []string `yaml:"entryPoints,omitempty"`
	Rule        string   `yaml:"rule"`
	Service     string   `yaml:"service"`
	Middlewares []string `yaml:"middlewares,omitempty"`
}

// Service represents a Traefik service
type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

// LoadBalancer represents a Traefik load balancer
type LoadBalancer struct {
	Servers []Server `yaml:"servers"`
}

// Server represents a Traefik server
type Server struct {
	URL string `yaml:"url"`
}

// Middleware represents a Traefik middleware
type Middleware struct {
	StripPrefix StripPrefix `yaml:"stripPrefix,omitempty"`
}

// StripPrefix represents the stripPrefix middleware configuration
type StripPrefix struct {
	Prefixes []string `yaml:"prefixes"`
}

// getActiveAliases retrieves all active aliases from the database
func getActiveAliases(isFirstInit bool) ([]prisma.Alias, error) {
	if isFirstInit {
		return []prisma.Alias{}, nil
	}

	var aliases []prisma.Alias
	result := prisma.DB.Where("is_active = ?", true).Find(&aliases)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch aliases: %w", result.Error)
	}

	return aliases, nil
}

// GenerateDynamicTraefikConfig generates the Traefik dynamic configuration based on running mock instances and aliases
func GenerateDynamicTraefikConfig() error {
	return GenerateDynamicTraefikConfigWithInit(false)
}

// GenerateDynamicTraefikConfigWithInit generates the Traefik dynamic configuration with an option to skip alias loading
func GenerateDynamicTraefikConfigWithInit(isFirstInit bool) error {
	// Create the directory for Traefik configuration if it doesn't exist
	traefikDir := filepath.Dir(lib.TRAEFIK_DYNAMIC_CONFIG_PATH)
	if err := utils.EnsureDirectoryExists(traefikDir); err != nil {
		return fmt.Errorf("failed to create Traefik configuration directory: %w", err)
	}

	// Get active aliases
	aliases, err := getActiveAliases(isFirstInit)
	if err != nil {
		log.Printf("Warning: Failed to get aliases: %v", err)
	}

	// Get running mock instances
	instances := repositories.MockInstanceRepo.GetAll()

	// Create the configuration structure
	config := TraefikDynamicConfig{
		HTTP: HTTP{
			Routers:     make(map[string]Router),
			Services:    make(map[string]Service),
			Middlewares: make(map[string]Middleware),
		},
	}

	// Process aliases
	for _, alias := range aliases {
		middlewareName := fmt.Sprintf("strip-%s", alias.Alias)

		// Define router for alias, with middleware to strip prefix
		config.HTTP.Routers[alias.Alias] = Router{
			Rule:        fmt.Sprintf("PathPrefix(`/%s`)", alias.Alias),
			Service:     alias.Alias,
			Middlewares: []string{middlewareName},
		}

		// Define service endpoint
		config.HTTP.Services[alias.Alias] = Service{
			LoadBalancer: LoadBalancer{
				Servers: []Server{
					{
						URL: fmt.Sprintf("http://localhost:%d", alias.Port),
					},
				},
			},
		}

		// Define middleware to strip alias prefix
		config.HTTP.Middlewares[middlewareName] = Middleware{
			StripPrefix: StripPrefix{
				Prefixes: []string{fmt.Sprintf("/%s", alias.Alias)},
			},
		}
	}

	// Add routers and services for each mock instance
	for _, instance := range instances {
		routerKey := fmt.Sprintf("mock-%d", instance.Port)
		serviceKey := fmt.Sprintf("mock-%d-service", instance.Port)

		// Create router
		config.HTTP.Routers[routerKey] = Router{
			EntryPoints: []string{"web"},
			Rule:        fmt.Sprintf("PathPrefix(`/api/%d`)", instance.Port),
			Service:     serviceKey,
			Middlewares: []string{"strip-api"},
		}

		// Create service
		config.HTTP.Services[serviceKey] = Service{
			LoadBalancer: LoadBalancer{
				Servers: []Server{
					{
						URL: fmt.Sprintf("http://localhost:%d", instance.Port),
					},
				},
			},
		}
	}

	// Define API route for backend
	config.HTTP.Routers["api"] = Router{
		Rule:    "PathPrefix(`/mock`)",
		Service: "backend",
	}
	config.HTTP.Services["backend"] = Service{
		LoadBalancer: LoadBalancer{
			Servers: []Server{
				{
					URL: fmt.Sprintf("http://localhost:%s", lib.SERVER_PORT),
				},
			},
		},
	}

	// Define frontend route
	config.HTTP.Routers["frontend"] = Router{
		Rule:    "PathPrefix(`/`)",
		Service: "frontend",
	}
	config.HTTP.Services["frontend"] = Service{
		LoadBalancer: LoadBalancer{
			Servers: []Server{
				{
					URL: "http://localhost:3005",
				},
			},
		},
	}

	// Add middleware for stripping API prefix
	config.HTTP.Middlewares["strip-api"] = Middleware{
		StripPrefix: StripPrefix{
			Prefixes: []string{"/api"},
		},
	}

	// Clean up middlewares if none were added
	if len(config.HTTP.Middlewares) == 0 {
		config.HTTP.Middlewares = nil
	}

	// Convert to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal Traefik configuration: %w", err)
	}

	// Write to file
	err = os.WriteFile(lib.TRAEFIK_DYNAMIC_CONFIG_PATH, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write Traefik configuration: %w", err)
	}

	log.Printf("Generated Traefik dynamic configuration at: %s (aliases: %d)",
		lib.TRAEFIK_DYNAMIC_CONFIG_PATH, len(aliases))
	return nil
}
