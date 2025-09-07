# My Go Package

A static site generator for hosting Go packages on custom domains.

## Installation

```bash
go install go.vrong.me/mygopkg@latest
```

## Quick Start

### 1. Configure Your Modules

Create a `modules.json` file in your project directory:

```json
{
  "hello-world": {
    "git": "https://github.com/hello/world.git",
    "description": "Just a package that says hello"
  },
  "nested/module": {
    "git": "https://github.com/nested/module.git",
    "description": "Showing off modules with nested folders in names"
  },
  "custom-branch": {
    "git": "https://github.com/custom/branch.git",
    "branch": "develop",
    "description": "Package with documentation pointing to a different branch"
  }
}
```

### 2. Generate Static Files

```bash
mygopkg -base-url yourdomain.com
```

This generates static HTML files in the `build/` directory.

### 3. Deploy

Upload the contents of the `build/` directory to your web server and configure it to serve from your custom domain.

## Deployment Examples

### GitHub Pages

1. Push your `build/` directory contents to a GitHub repository
2. Enable GitHub Pages in repository settings
3. Configure your custom domain in the Pages settings

```yaml
# .github/workflows/deploy.yml
name: Deploy to GitHub Pages

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v5
    - uses: actions/setup-go@v6
      with:
        go-version: '1.25'
    - run: go install go.vrong.me/mygopkg@latest
    - run: mygopkg -base-url yourdomain.com
    - uses: peaceiris/actions-gh-pages@v4
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./build
```

## CLI Reference

```
mygopkg [OPTIONS]

Options:
  -base-url string
        Base URL for your custom domain (required)
  -modules string
        Path to modules JSON file (default: "modules.json")
  -build-dir string
        Output directory for generated files (default: "build")
  -help
        Show help message
```

### Examples

```bash
# Basic usage
mygopkg -base-url pkg.example.com

# Custom modules file and build directory
mygopkg -base-url pkg.example.com -modules my-packages.json -build-dir dist
```

## Modules JSON Reference

```json
{
    "<module name>": {
        "git": "<git url>",
        "description": "<description for the module>",
        "branch": "<branch to point the documentation to; defaults to main>"
    }
}
```

## Configuration Reference

### modules.json Schema

```json
{
  "<module-name>": {
    "git": "<repository-url>",
    "description": "<module-description>",
    "branch": "<branch-name>"
  }
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `git` | string | ✅ | Git repository URL (HTTPS or SSH) |
| `description` | string | ✅ | Brief description of the module |
| `branch` | string | ❌ | Target branch for documentation (defaults to `main`) |

### Module Naming

- Use forward slashes for nested modules: `company/pkg/subpkg`
- Avoid special characters except hyphens and underscores
- Module names become URL paths: `yourdomain.com/company/pkg/subpkg`

## How It Works

1. **Discovery**: Reads your `modules.json` configuration
2. **Generation**: Creates static HTML pages with proper Go module metadata
3. **Redirects**: Generates `go get` redirects pointing to your Git repositories
4. **Documentation**: Links to pkg.go.dev for comprehensive documentation

When users run `go get yourdomain.com/your-module`, they'll be redirected to the correct Git repository while maintaining your custom domain branding.

## License

See [LICENSE](./LICENSE)
