# Watermelon

Sandbox that isolates your project inside a Linux VM.

## Why

Running untrusted code like `npm install` risks:
- **Filesystem access**: Malicious packages reading `~/.ssh`, `~/.aws`, `~/.gnupg`
- **Network exfiltration**: Packages sending stolen data to remote servers
- **Persistent changes**: Packages installing backdoors, cron jobs, launch agents

Watermelon isolates commands inside a Linux VM where they can't touch your host system.

## How It Works

```
┌─────────────────────────────────────┐
│           Host (macOS)              │
│                                     │
│  ~/Projects/myapp/                  │
│  ├── .watermelon.toml  ← config     │
│  └── src/, package.json...          │
│                                     │
└──────────────┬──────────────────────┘
               │ virtiofs mount
               ▼
┌─────────────────────────────────────┐
│           VM (Linux)                │
│                                     │
│  /project/  ← your project (r/w)    │
│  /tools/    ← node, python (r/o)    │
│                                     │
│  Network: allowlist only            │
│  Filesystem: isolated               │
└─────────────────────────────────────┘
```

- Your project is mounted read-write inside the VM
- Network is restricted to domains you explicitly allow
- The VM persists between sessions (installed deps survive)
- Host filesystem is completely isolated

## Requirements

- macOS (Apple Silicon or Intel)
- [Lima](https://lima-vm.io/) installed: `brew install lima`

## Installation

```bash
go install github.com/saeta/watermelon/cmd/watermelon@latest
```

Or build from source:

```bash
git clone https://github.com/saeta/watermelon
cd watermelon
go build -o watermelon ./cmd/watermelon
```

## Quick Start

```bash
cd your-project

# Create config file
watermelon init

# Edit .watermelon.toml to allow npm registry
# network.allow = ["registry.npmjs.org"]

# Enter the sandbox
watermelon run

# Inside the VM:
npm install
npm run dev
exit

# Later, re-enter (state is preserved)
watermelon run
```

## Commands

| Command | Description |
|---------|-------------|
| `watermelon init` | Create `.watermelon.toml` config file |
| `watermelon run` | Enter the sandbox VM (creates if needed) |
| `watermelon exec <cmd>` | Run a single command without interactive shell |
| `watermelon stop` | Stop the VM (preserves state) |
| `watermelon destroy` | Delete VM and all state |
| `watermelon status` | Show VM status for current project |
| `watermelon list` | List all watermelon VMs |
| `watermelon violations` | Show blocked network requests |

## Configuration

Create `.watermelon.toml` in your project root:

```toml
[vm]
image = "ubuntu-22.04"

[network]
# Only these domains are reachable (all others blocked)
allow = [
    "registry.npmjs.org",
    "github.com",
    "*.githubusercontent.com",
]

[tools]
# Tools to make available in VM
node = "20"
git = "latest"

[mounts]
# Additional host paths to mount (read-only)
"~/.gitconfig" = { target = "/home/dev/.gitconfig" }
"~/.npmrc" = { target = "/home/dev/.npmrc" }

[ports]
# Ports to forward from VM to host
forward = [3000, 5173, 8080]

[resources]
memory = "4GB"
cpus = 2
disk = "20GB"

[security]
# What happens on policy violations: "log", "fail", or "silent"
on_violation = "log"
```

### Defaults

| Setting | Default |
|---------|---------|
| `network.allow` | `[]` (no network) |
| `resources.memory` | `2GB` |
| `resources.cpus` | `1` |
| `resources.disk` | `10GB` |
| `security.on_violation` | `log` |

## Examples

Ready-to-use configurations for common project types are available in [`docs/examples/`](./docs/examples/):

| Example | Use Case |
|---------|----------|
| [react-app](./docs/examples/react-app/) | React/Vite development |
| [nextjs](./docs/examples/nextjs/) | Next.js with API routes |
| [python-django](./docs/examples/python-django/) | Django web application |
| [python-ml](./docs/examples/python-ml/) | Machine learning (PyTorch/TensorFlow) |
| [rust-project](./docs/examples/rust-project/) | Rust with Cargo |
| [go-project](./docs/examples/go-project/) | Go development |
| [foundry](./docs/examples/foundry/) | Ethereum contracts (Foundry) |
| [hardhat](./docs/examples/hardhat/) | Ethereum contracts (Hardhat) |
| [monorepo](./docs/examples/monorepo/) | Full-stack (Node + Python) |
| [audit-package](./docs/examples/audit-package/) | Inspect suspicious packages |

Copy any example to your project:

```bash
cp docs/examples/react-app/.watermelon.toml ~/my-project/
```

### Quick Examples

**Node.js:**
```toml
[network]
allow = ["registry.npmjs.org", "github.com"]

[tools]
node = "20"

[ports]
forward = [3000, 5173]
```

**Python:**
```toml
[network]
allow = ["pypi.org", "files.pythonhosted.org"]

[tools]
python = "3.11"

[ports]
forward = [8000]
```

## Security Model

**Watermelon protects against:**
- Packages reading sensitive host files (`~/.ssh`, `~/.aws`)
- Packages exfiltrating data to unknown domains
- Packages installing persistent backdoors on host
- Runaway processes consuming host resources

**Watermelon does NOT protect against:**
- Malicious code running inside the VM (it has full access there)
- Attacks on your project files (they're mounted read-write)

This is a developer safety sandbox, not a jail for untrusted multi-tenant workloads.

## How Network Isolation Works

Watermelon generates iptables rules inside the VM:

```bash
# Allow specified domains
iptables -A OUTPUT -d registry.npmjs.org -j ACCEPT
iptables -A OUTPUT -d github.com -j ACCEPT

# Allow DNS
iptables -A OUTPUT -p tcp --dport 53 -j ACCEPT
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT

# Allow established connections
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# Block everything else
iptables -A OUTPUT -j REJECT
```

Blocked requests are logged to `.watermelon/violations.log`.

## Development

```bash
# Run tests
go test ./...

# Run E2E tests (requires Lima)
go test -tags=e2e ./test/...

# Build
go build -o watermelon ./cmd/watermelon
```

## License

MIT
