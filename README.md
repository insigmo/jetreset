# Jetreset

[![Release](https://img.shields.io/github/v/release/insigmo/jetreset?style=flat-square)](https://github.com/insigmo/jetreset/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)


A cross-platform CLI tool written in Go for resetting the trial period of JetBrains IDEs.

Supports **IntelliJ IDEA**, **GoLand**, **PyCharm**, **WebStorm**, **CLion**, **PhpStorm**, **Rider**, **DataGrip**, and **RubyMine**.

---

## Installation

### Download binary (recommended)

Download the latest pre-built binary for your platform from the [Releases](https://github.com/insigmo/jetreset/releases/latest) page:

| Platform | Architecture | File                         |
|----------|--------------|------------------------------|
| Linux    | amd64        | `jetreset-linux-amd64`       |
| Linux    | arm64        | `jetreset-linux-arm64`       |
| macOS    | amd64        | `jetreset-darwin-amd64`      |
| macOS    | arm64 (M1+)  | `jetreset-darwin-arm64`      |
| Windows  | amd64        | `jetreset-windows-amd64.exe` |
| Windows  | arm64        | `jetreset-windows-arm64.exe` |

#### Linux / macOS

```bash
# Download (example for linux/amd64)
curl -L https://github.com/insigmo/jetreset/releases/latest/download/jetreset-linux-amd64 -o jetreset

# Make executable
chmod +x jetreset

# Run
./jetreset
```

#### Windows

Download `jetreset-windows-amd64.exe` and run it directly or via PowerShell:

```powershell
.\jetreset-windows-amd64.exe
```

### Build from source

Requires Go 1.22+.

```bash
git clone https://github.com/insigmo/jetreset.git
cd jetreset
go build -ldflags="-s -w" -trimpath -o jetreset .
```

---

## Usage

```bash
jetreset [flags]
```

Run without flags to reset all detected JetBrains IDEs interactively. And se

### Flags

| Flag             | Description                             |
|------------------|-----------------------------------------|
| `--run-schedule` | Run schedurler for reseting every month |
| `--stop`         | Stop schedurler                         |

---

## How it works

`jetreset` removes the evaluation-related files that JetBrains IDEs use to track the trial period. Specifically, it clears:

- `eval/` folder inside the IDE config directory (contains the trial license key)
- `options/other.xml` — stores the trial start date reference
- On Windows: relevant entries under `HKEY_CURRENT_USER\Software\JavaSoft` and the `%APPDATA%\JetBrains` folder

Config directories by platform:

| Platform | Path                                      |
|----------|-------------------------------------------|
| Linux    | `~/.config/JetBrains/<IDE><version>/`     |
| macOS    | `~/Library/Application Support/JetBrains/<IDE><version>/` |
| Windows  | `%APPDATA%\JetBrains\<IDE><version>\`     |

Personal settings, plugins, keymaps, and code style configurations are **not affected**.

---

## Supported IDEs

- IntelliJ IDEA
- GoLand
- PyCharm
- WebStorm
- CLion
- PhpStorm
- Rider
- DataGrip
- RubyMine
