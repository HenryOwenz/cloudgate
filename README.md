```bash             __                               __                           __               
            /  |                             /  |                         /  |              
   _______  $$ |   ______    __    __    ____$$ |   ______     ______    _$$ |_      ______  
  /       | $$ |  /      \  /  |  /  |  /    $$ |  /      \   /      \  / $$   |    /      \ 
 /$$$$$$$/  $$ | /$$$$$$  | $$ |  $$ | /$$$$$$$ | /$$$$$$  |  $$$$$$  | $$$$$$/    /$$$$$$  |
 $$ |       $$ | $$ |  $$ | $$ |  $$ | $$ |  $$ | $$ |  $$ |  /    $$ |   $$ | __  $$    $$ |
 $$ \_____  $$ | $$ \__$$ | $$ \__$$ | $$ \__$$ | $$ \__$$ | /$$$$$$$ |   $$ |/  | $$$$$$$$/ 
 $$       | $$ | $$    $$/  $$    $$/  $$    $$ | $$    $$ | $$    $$ |   $$  $$/  $$       |
   $$$$$$$/ $$/   $$$$$$/    $$$$$$/    $$$$$$$/   $$$$$$$ |  $$$$$$/     $$$$/     $$$$$$$/ 
                                                  /  \__$$ |                              
                                                  $$    $$/                               
                                                   $$$$$$/                                
```

# cloudgate

A terminal-based application that unifies multi-cloud operations across AWS, Azure, and GCP.

> *Where your clouds converge.*

[![Latest Release](https://img.shields.io/github/release/HenryOwenz/cloudgate.svg)](https://github.com/HenryOwenz/cloudgate/releases)
[![Lint](https://github.com/HenryOwenz/cloudgate/actions/workflows/lint.yml/badge.svg)](https://github.com/HenryOwenz/cloudgate/actions/workflows/lint.yml)
[![Test](https://github.com/HenryOwenz/cloudgate/actions/workflows/test.yml/badge.svg)](https://github.com/HenryOwenz/cloudgate/actions/workflows/test.yml)
[![Build](https://github.com/HenryOwenz/cloudgate/actions/workflows/build.yml/badge.svg)](https://github.com/HenryOwenz/cloudgate/actions/workflows/build.yml)
[![Dependabot Status](https://img.shields.io/badge/Dependabot-enabled-brightgreen.svg)](https://github.com/HenryOwenz/cloudgate/blob/main/.github/dependabot.yml)
[![Go ReportCard](https://goreportcard.com/badge/HenryOwenz/cloudgate)](https://goreportcard.com/report/HenryOwenz/cloudgate)

<p align="center">
  <img src="https://github.com/HenryOwenz/cloudgate/releases/download/v0.1.4/cloudgate-demo.gif" width="100%" alt="cloudgate Demo">
</p>

## Features 

- **AWS Integration**
  - Multi-account/region management


  <details>
  <summary><b>📋 Available AWS Services & Operations</b></summary>
  
  | Service | Operation | Description |
  |---------|-----------|-------------|
  | **CodePipeline** | | |
  | | Pipeline Status | View status of all pipelines and their stages |
  | | Pipeline Approvals | List, approve, or reject pending manual approvals |
  | | Start Pipeline | Trigger pipeline execution with latest commit or specific revision |
  | **Lambda** | | |
  | | Function Status | View all Lambda functions with runtime and last update info<br><br>**Function Details View:**<br>Select any function to inspect detailed configuration including memory, timeout, architecture, and other key attributes |
  | | Execute Function | Invoke Lambda functions directly with custom payload and view execution results |
  
  *Operations can be performed using any configured AWS profile and region (one active profile/region at a time)*  
  *Multi-account aggregation for services will be coming in the future*
  </details>

- **Terminal UI**
  - Fast, keyboard-driven interface
  - Context-aware navigation
  - Visual feedback and safety controls
  - Formatted display of timestamps and resource sizes
  - Vim-style navigation ('-' for backwards navigation, 'k/j' for up/down navigation, etc.)

- **Coming Soon**
  - Azure integration
  - GCP support
  - Additional AWS services (S3, EC2, etc.)

## Installation

### Quick Install / Upgrade

**Linux/macOS:**
```bash
bash -c "$(curl -fsSL https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.sh)"
```

**Windows (PowerShell):**
```powershell
Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.ps1'))
```

These scripts will automatically download and install the latest version of cloudgate, handling upgrades cleanly if you already have it installed.

### From Source

```bash
git clone https://github.com/HenryOwenz/cloudgate.git
cd cloudgate
make build
make install  # Installs as 'cg' in your $GOPATH/bin
```

## Requirements

- Go 1.22+
- AWS credentials configured in `~/.aws/credentials` or `~/.aws/config`

## Usage

```bash
cg  # Launch the application
```

### Command Line Options

| Option | Description |
|--------|-------------|
| `cg --upgrade` or `cg -u` | Upgrade cloudgate to the latest version |
| `cg upgrade` | Upgrade cloudgate to the latest version (alternative syntax) |
| `cg --version` or `cg -v` | Display the current version of cloudgate |
| `cg version` | Display the current version of cloudgate (alternative syntax) |

### Navigation

<details>
<summary><b>🎮 Keyboard Navigation Commands</b></summary>

| Key                | Action                   |
|--------------------|--------------------------|
| ↑/↓ or j/k         | Navigate up/down         |
| ←/→ or h/l         | Previous/Next page (in paginated views) |
| Enter              | Select/Confirm           |
| Esc or -           | Go back/Cancel           |
| q                  | Quit application         |
| Ctrl+c             | Force quit               |
| g                  | Jump to top              |
| G                  | Jump to bottom           |
| Home/End           | Jump to top/bottom (alternative) |
| u or Ctrl+u        | Half page up             |
| d or Ctrl+d        | Half page down           |
| b or PgUp          | Page up                  |
| f or PgDown        | Page down                |
| /                  | Search (in paginated views) |
| i                  | Enter input mode (in Lambda execution view) |

**Note:** Vim-style navigation keys (j, k, h, l, g, G, etc.) work in table views but are passed through as text when in input mode. Use Esc to exit text input mode.
</details>

## Development

### Testing

```bash
make test          # Run all tests
make test-unit     # Run unit tests only
make test-integration  # Run integration tests only
make test-coverage  # Generate coverage report
```

### CI/CD

This project uses GitHub Actions for continuous integration:
- Automated builds on each push and pull request
- Unit and integration tests
- Code linting with golangci-lint
- Test coverage reporting
- Automatic testing of Dependabot PRs

## Architecture

cloudgate uses a dual-layer architecture:
- Provider layer: Abstracts cloud provider APIs
- UI layer: Handles user interaction and workflow

The application follows a modular design pattern that makes it easy to add new cloud services and operations. Each service is implemented as a separate module with clear interfaces, allowing for independent development and testing.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
