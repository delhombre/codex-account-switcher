# CXA - Codex Account Switcher 🔄

**A fast, beautiful CLI to manage multiple OpenAI Codex accounts.**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

---

## ✨ Features

- **⚡ Instant switching** — Directory-based storage (no zip/unzip delays)
- **🎨 Beautiful TUI** — Interactive interface powered by [Charm](https://charm.sh)
- **🔐 Secure** — Authentication stays isolated per account
- **📋 Simple commands** — `cxa switch bruno` and you're done

---

## 📦 Installation

### Homebrew (macOS/Linux)

```bash
brew install delhombre/tap/cxa
```

### Go Install

```bash
go install github.com/delhombre/cxa/cmd/cxa@latest
```

### From Source

```bash
git clone https://github.com/delhombre/cxa.git
cd cxa
make install
```

---

## 🚀 Quick Start

```bash
# Save your current account
cxa save personal

# Add another account
cxa save work   # After logging in with 'codex login'

# Switch between accounts
cxa switch work
cxa switch personal

# Or use the interactive TUI
cxa
```

---

## 📖 Commands

| Command             | Description                     |
| ------------------- | ------------------------------- |
| `cxa`               | Launch interactive TUI          |
| `cxa list`          | List all saved accounts         |
| `cxa switch <name>` | Switch to an account            |
| `cxa save <name>`   | Save current session as account |
| `cxa current`       | Show active account             |
| `cxa version`       | Print version                   |

### Aliases

- `cxa ls` → `cxa list`
- `cxa sw <name>` → `cxa switch <name>`
- `cxa use <name>` → `cxa switch <name>`

---

## 🎨 TUI Preview

```
╭──────────────────────────────────────────╮
│  🔄 Codex Accounts                       │
├──────────────────────────────────────────┤
│  ● bruno          (current)              │
│  ○ work                                  │
│  ○ client-project                        │
├──────────────────────────────────────────┤
│  enter: switch  •  /: filter  •  q: quit │
╰──────────────────────────────────────────╯
```

---

## 📁 Data Locations

| Path                           | Purpose                           |
| ------------------------------ | --------------------------------- |
| `~/.codex`                     | Active Codex session              |
| `~/codex-data/accounts/<name>` | Saved account data                |
| `~/.codex-switch/state.json`   | Current/previous account tracking |

---

## 🛠 Development

```bash
# Run in dev mode
make dev

# Run tests
make test

# Build binary
make build

# Install locally
make install
```

---

## 🙏 Credits

Inspired by [bashar94/codex-cli-account-switcher](https://github.com/bashar94/codex-cli-account-switcher).

Built with [Charm](https://charm.sh) libraries:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — Components

---

## 📄 License

MIT License — see [LICENSE](LICENSE)
