package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

//go:embed templates/*.gotmpl
var templates embed.FS

type ModuleConfig struct {
	GitURL        string `json:"git"`
	DefaultBranch string `json:"branch"`
	Description   string `json:"description"`
}

func main() {
	baseURL := flag.String("base-url", "", "Base URL")
	modulesFile := flag.String("modules", "modules.json", "Modules")
	buildDir := flag.String("build-dir", "build", "Build directory")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if help != nil && *help {
		fmt.Println(`mygopkg [OPTIONS]

Options:
  -base-url string
        Base URL for your custom domain (required)
  -modules string
        Path to modules JSON file (default: "modules.json")
  -build-dir string
        Output directory for generated files (default: "build")
  -help
        Show help message`)
		return
	}

	if err := run(baseURL, modulesFile, buildDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(baseURL *string, modulesFile *string, buildDir *string) error {
	if baseURL == nil || *baseURL == "" {
		return fmt.Errorf("-base-url is required")
	}

	if modulesFile == nil || *modulesFile == "" {
		return fmt.Errorf("-modules is required")
	}

	if buildDir == nil || *buildDir == "" {
		return fmt.Errorf("-build-dir is required")
	}

	modulesJSON, err := os.ReadFile(*modulesFile)
	if err != nil {
		return fmt.Errorf("failed to read modules.json: %w", err)
	}

	var modules map[string]*ModuleConfig
	err = json.Unmarshal(modulesJSON, &modules)
	if err != nil {
		return fmt.Errorf("failed to unmarshal modules.json: %w", err)
	}

	for _, config := range modules {
		if config.DefaultBranch == "" {
			config.DefaultBranch = "main"
		}
	}

	if err := os.MkdirAll(*buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	modulePath := func(module string) string {
		return *baseURL + "/" + module
	}

	funcMap := template.FuncMap{
		"modulePath": modulePath,
		"moduleDocURL": func(module string) string {
			return "https://pkg.go.dev/" + modulePath(module)
		},
	}

	templates, err := template.New("templates").Funcs(funcMap).ParseFS(templates, "templates/*.gotmpl")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "mygopkg")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := writeTmpTemplate(templates, "index.gotmpl", "", map[string]any{
		"baseURL": baseURL,
		"modules": modules,
	}, tmpDir); err != nil {
		return fmt.Errorf("failed to write index template: %w", err)
	}

	for module, config := range modules {
		if err := writeTmpTemplate(templates, "module.gotmpl", module, map[string]any{
			"module": module,
			"config": config,
		}, tmpDir); err != nil {
			return fmt.Errorf("failed to write module template: %w", err)
		}
	}

	backupBuildDir := fmt.Sprintf("%s-bak-%s", *buildDir, time.Now())
	if err := os.Rename(*buildDir, backupBuildDir); err != nil {
		return fmt.Errorf("failed to remove build dir: %w", err)
	}

	if err := os.Rename(tmpDir, *buildDir); err != nil {
		_ = os.Rename(backupBuildDir, *buildDir) // Can't do nothin
		return fmt.Errorf("failed to rename tmp dir to build dir: %w", err)
	}

	_ = os.RemoveAll(backupBuildDir) // It's ok to fail here
	return nil
}

func writeTmpTemplate(templates *template.Template, tmpl, name string, data any, tmpDir string) error {
	dirPath := filepath.Join(tmpDir, name)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	file, err := os.Create(filepath.Join(dirPath, "index.html"))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return templates.ExecuteTemplate(file, tmpl, data)
}
