# CIDR Parser

A beautiful command-line tool for parsing CIDR subnet masks and checking IP address membership, built with Go.

## Features

- **Parse CIDR Notation** - Display comprehensive network information including:
  - Network address and subnet mask
  - Broadcast address
  - IP range (total and usable)
  - Host counts (total and usable)

- **IP Membership Checking** - Verify if an IP address belongs to one or more CIDR ranges

- **Config File Support** - Load default CIDR ranges from `~/.cidr` file

- **Beautiful Output** - Color-coded terminal output with clear visual hierarchy using Lipgloss

## Installation

```bash
go install github.com/trahma/cidr@latest
```

Or build from source:

```bash
git clone https://github.com/trahma/cidr
cd cidr
go build -o cidr .
```

## Post-Installation Setup

After installing with `go install`, you need to ensure `GOPATH/bin` is in your PATH.

### Check if already configured

```bash
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✓ Already in PATH" || echo "✗ Not in PATH"
```

### Add to PATH

If not already in PATH, add it to your shell configuration:

**Bash** (`~/.bashrc`):
```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

**Zsh** (`~/.zshrc`):
```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

**Fish** (`~/.config/fish/config.fish`):
```bash
fish_add_path (go env GOPATH)/bin
```

### Verify Installation

```bash
cidr --help
```

If you see the help output, you're all set!

## Usage

### Parse a CIDR range

```bash
cidr 192.168.1.0/24
```

Output:
```
CIDR Information

CIDR: 192.168.1.0/24
Network Address: 192.168.1.0
Subnet Mask: 255.255.255.0
Broadcast Address: 192.168.1.255

IP Range: 192.168.1.0 - 192.168.1.255
Usable IPs: 192.168.1.1 - 192.168.1.254

Total Hosts: 256
Usable Hosts: 254
```

### Check if an IP is in a CIDR range

```bash
cidr 192.168.1.0/24 --check 192.168.1.50
```

### Check IP against config file CIDRs

```bash
cidr --check 192.168.5.10
```

This will check the IP against all CIDR ranges defined in your `~/.cidr` config file.

## Configuration File

Create a `~/.cidr` file with your default CIDR ranges (one per line):

```
# Private network ranges
192.168.0.0/16
10.0.0.0/8
172.16.0.0/12
```

Lines starting with `#` are treated as comments and ignored.

You can also specify a custom config file:

```bash
cidr --config /path/to/custom.cidr --check 192.168.1.1
```

## Command-Line Options

```
Flags:
  -c, --check string    Check if an IP address is within the CIDR range
  -f, --config string   Path to .cidr config file (defaults to ~/.cidr)
  -h, --help            help for cidr
```

## Examples

Parse a large network:
```bash
cidr 10.0.0.0/8
```

Check an IP against multiple ranges:
```bash
cidr --check 172.16.5.100
```

Parse with custom config file:
```bash
cidr --config ./networks.cidr --check 192.168.1.1
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

MIT
