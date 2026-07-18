# Jetreset

[![Release](https://img.shields.io/github/v/release/insigmo/jetreset?style=flat-square)](https://github.com/insigmo/jetreset/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev)


A cross-platform CLI tool for resetting the trial period of JetBrains IDEs.
If you already have licenses in your IDEs, it'll be resetting too.
If you're from Russia, you need to have VPN to start a new trial period.

Supports **IntelliJ IDEA**, **GoLand**, **PyCharm**, **WebStorm**, **CLion**, **PhpStorm**, **Rider**, **DataGrip**, and **RubyMine**.

---

## Installation

### Download binary (recommended)

Download the latest pre-built binary for your platform from the [Releases](https://github.com/insigmo/jetreset/releases/latest) page:

### Linux / macOS

```bash
# Download and run installer
curl -fsSL https://raw.githubusercontent.com/insigmo/jetreset/refs/heads/master/install.sh | bash
```

### Windows (PowerShell)

```powershell
# Download and run installer
irm https://raw.githubusercontent.com/insigmo/jetreset/refs/heads/master/install.ps1 | iex
```

### Run on Linux/MacOS
```bash
./jetreset          # Reset all detected JetBrains IDEs
./jetreset -v       # Same, but print debug logs to stderr
```

### Run on Windows
```powershell
.\jetreset.exe
.\jetreset.exe -v
```

### What a reset does

When you run `jetreset` without flags it:

1. **Detects** any running JetBrains IDEs and remembers them.
2. **Closes** them (gracefully first, then forcefully if needed) so the trial files are released.
3. **Wipes** the trial/eval/license state for every supported product.
4. **Relaunches** the IDEs that were running.

`-v` prints debug logs (PIDs, paths, signals sent, launchers used) to **stderr** — the normal success output on stdout stays clean.

> ⚠️ **Closing IDEs loses unsaved work.** A warning is printed in interactive sessions. There is no confirmation prompt because scheduled (cron) runs cannot prompt.

### Run scheduler
```bash
./jetreset --run-schedule # Run scheduler for reseting every month
./jetreset --stop         # Stop schedurler
```

Scheduled runs close and reset running IDEs but **cannot relaunch them** on Linux because there is no graphical session — run `jetreset` from your desktop session to reopen IDEs, or reopen them manually. On Windows, the Task Scheduler task must run *only when the user is logged on* for a relaunched IDE to be visible.

### Build from source

Requires Go 1.22+.

```bash
git clone https://github.com/insigmo/jetreset.git
cd jetreset
go build -ldflags="-s -w" -trimpath -o jetreset .
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
