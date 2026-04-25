# Jetreset

[![Release](https://img.shields.io/github/v/release/insigmo/jetreset?style=flat-square)](https://github.com/insigmo/jetreset/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)


A cross-platform CLI tool written in Go for resetting the trial period of JetBrains IDEs.

Supports **IntelliJ IDEA**, **GoLand**, **PyCharm**, **WebStorm**, **CLion**, **PhpStorm**, **Rider**, **DataGrip**, and **RubyMine**.

---

## Installation

### Download binary (recommended)

Download the latest pre-built binary for your platform from the [Releases](https://github.com/insigmo/jetreset/releases/latest) page:

#### Linux / macOS

```bash
# Download and run installer
curl -fsSL https://raw.githubusercontent.com/insigmo/jetreset/refs/heads/master/install.sh | bash

# Run
./jetreset # Run without flags to reset all detected JetBrains IDEs
```

Run scheduler
```bash
./jetreset --run-schedule # Run scheduler for reseting every month
./jetreset --stop         # Stop schedurler
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
```

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
