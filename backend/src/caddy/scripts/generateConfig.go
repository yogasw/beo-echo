package scripts

import (
	"beo-echo/backend/src/caddy/config"
	"beo-echo/backend/src/lib"
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Domain      string
	ProxyTarget string
}

func GenerateDynamicConfig(ctx context.Context, tmplPath, outputPath string, configs []Config) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to parse template")
		return fmt.Errorf("template parse error: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("path", outputPath).Msg("failed to create file")
		return fmt.Errorf("create file error: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, configs); err != nil {
		os.Remove(outputPath)
		log.Ctx(ctx).Error().Err(err).Msg("failed to execute template")
		return fmt.Errorf("template execution error: %w", err)
	}

	log.Ctx(ctx).Info().Str("output", outputPath).Msg("config generated")
	return nil
}

func ValidateConfig(ctx context.Context, configPath string) error {
	cmd := exec.Command("caddy", "validate", "--config", configPath, "--adapter", "caddyfile")
	out, err := cmd.CombinedOutput()
	if err != nil {
		errorMessage := string(out)
		log.Ctx(ctx).Error().Err(err).Str("path", configPath).Bytes("output", out).Msg("validation failed")
		return fmt.Errorf("caddy config validation failed: %w, output: %s", err, errorMessage)
	}
	log.Ctx(ctx).Info().Str("path", configPath).Msg("validation succeeded")
	return nil
}

func EnsureImportDynamicConf(ctx context.Context, mainConfigFile, dynamicImportPath string) error {
	f, err := os.Open(mainConfigFile)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("file", mainConfigFile).Msg("cannot open main config")
		return fmt.Errorf("failed to open main config file: %w", err)
	}
	defer f.Close()

	importLine := fmt.Sprintf("import %s", dynamicImportPath)
	scanner := bufio.NewScanner(f)
	found := false

	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == importLine {
			found = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed reading main config")
		return fmt.Errorf("failed reading config file: %w", err)
	}

	if found {
		log.Ctx(ctx).Info().Str("file", mainConfigFile).Msg("import already present")
		return nil
	}

	fw, err := os.OpenFile(mainConfigFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("cannot open main config for append")
		return fmt.Errorf("failed to open main config for appending: %w", err)
	}
	defer fw.Close()

	if _, err := fw.WriteString("\n" + importLine + "\n"); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to write import line")
		return fmt.Errorf("failed to write import line: %w", err)
	}

	log.Ctx(ctx).Info().Str("file", mainConfigFile).Str("import", dynamicImportPath).Msg("import inserted")
	return nil
}

var caddyTemplate = `{{range .}}
{{.Domain}} {
  rewrite * /{http.request.host.labels.3}{http.request.uri}

  reverse_proxy http://{{.ProxyTarget}}
}
{{end}}`

func GenerateSingleConfigFromText(
	ctx context.Context,
	configs []Config,
) error {
	// setup
	mainConfigStr := config.CADDY_DEFAULT_CONFIG
	outputPath := lib.CANDY_DIR + "/Dynamic.conf"
	dynamicImportPath := "./Dynamic.conf"
	mainConfigPath := lib.CANDY_DIR + "/Caddyfile"

	// check if main config file exists or not when not found create it
	if _, err := os.Stat(mainConfigPath); os.IsNotExist(err) {
		if err := os.WriteFile(mainConfigPath, []byte(mainConfigStr), 0644); err != nil {
			log.Ctx(ctx).Error().Err(err).Str("file", mainConfigPath).Msg("failed to create main config file")
			return fmt.Errorf("failed to create main config file: %w", err)
		}
		log.Ctx(ctx).Info().Str("file", mainConfigPath).Msg("main config file created")
	}

	// 1. Parse template string
	tmpl, err := template.New("caddy").Parse(caddyTemplate)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to parse template string")
		return fmt.Errorf("failed to parse caddy template string: %w", err)
	}

	// 2. Write dynamic config file
	f, err := os.Create(outputPath)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("path", outputPath).Msg("failed to create dynamic config file")
		return fmt.Errorf("failed to create dynamic config file at %s: %w", outputPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, configs); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to write dynamic config")
		os.Remove(outputPath)
		return fmt.Errorf("failed to write dynamic config to %s: %w", outputPath, err)
	}
	log.Ctx(ctx).Info().Str("output", outputPath).Msg("dynamic config generated (from text)")

	// 3. Validate dynamic config alone
	if err := ValidateConfig(ctx, outputPath); err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("dynamic config validation failed: %w", err)
	}

	// 4. Ensure import into main config
	if err := EnsureImportDynamicConf(ctx, mainConfigPath, dynamicImportPath); err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("ensure import into main config failed: %w", err)
	}

	// 5. Validate full config
	if err := ValidateConfig(ctx, mainConfigPath); err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("full config validation failed: %w", err)
	}

	log.Ctx(ctx).Info().
		Str("main", mainConfigPath).
		Str("imported", dynamicImportPath).
		Msg("config validated and import successful")

	return nil
}
