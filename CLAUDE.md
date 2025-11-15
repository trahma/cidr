# CIDR Parser - Project Context

This document provides context for AI assistants (like Claude) working on this project.

## Project Overview

CIDR is a command-line tool written in Go for parsing CIDR subnet masks and checking IP address membership. It focuses on providing a beautiful, user-friendly experience with styled terminal output.

## Technology Stack

- **Language**: Go 1.25.3
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) - Command-line interface and flag parsing
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal output styling and colors

## Project Structure

```
cidr/
├── main.go              # Entry point - calls cmd.Execute()
├── cmd/
│   └── root.go          # Cobra root command with all CLI logic
├── go.mod               # Module definition (github.com/trahma/cidr)
├── go.sum               # Dependency checksums
├── README.md            # User-facing documentation
├── CLAUDE.md            # This file - AI assistant context
└── .gitignore           # Excludes binary and build artifacts
```

## Core Functionality

### 1. CIDR Parsing
- Takes CIDR notation (e.g., `192.168.1.0/24`)
- Displays network address, subnet mask, broadcast address
- Shows IP ranges (total and usable)
- Calculates host counts (total and usable)

### 2. IP Membership Checking
- Checks if an IP belongs to CIDR range(s)
- Works with single CIDR or multiple from config file
- Visual indicators: ✓ (in range), ○ (not in range)

### 3. Config File Support
- Default location: `~/.cidr`
- Format: One CIDR per line
- Supports comments (lines starting with `#`)
- Can specify custom path with `--config` flag

## Command Structure

Current design:
- `cidr [CIDR]` - Parse a CIDR (positional argument)
- `cidr [CIDR] --check [IP]` - Check IP against specific CIDR
- `cidr --check [IP]` - Check IP against config file CIDRs

Flags:
- `-c, --check` - IP address to check
- `-f, --config` - Custom config file path
- `-h, --help` - Show help

## Design Decisions

### Styling Philosophy
- **Title style**: Bold cyan (#86) for section headers
- **Label style**: Bold magenta (#205) for field labels
- **Value style**: Light blue (#117) for values
- **Success style**: Bold green (#42) for positive results
- **Error style**: Bold red (#196) for errors
- **Info style**: Yellow (#226) for neutral info
- **Dim style**: Dark gray (#240) for config file indicator
- **Help style**: Italic gray (#243) for help hints

### User Experience
- Help hint appears once at the end of output
- Config file path shown in dark gray when loaded
- Clear visual hierarchy with colors and spacing
- Single help message regardless of output length

### Code Organization
- All CLI logic in `cmd/root.go`
- Helper functions for IP calculations at bottom of file
- Styles defined as package-level variables
- Config loading returns both CIDRs and path for display

## Key Functions

### `runCIDR()`
Main entry point that:
1. Processes arguments and loads config if needed
2. Shows config file indicator if loaded
3. Routes to either IP checking or CIDR display
4. Shows help hint at the end

### `displayCIDRInfo()`
Parses and displays information for a single CIDR:
- Network details
- IP ranges
- Host counts

### `checkIPInCIDRs()`
Checks an IP against one or more CIDRs:
- Validates IP address
- Iterates through CIDRs
- Shows results with visual indicators
- Summary message

### `loadConfigCIDRs()`
Loads CIDR ranges from config file:
- Returns: (cidrs, configPath, error)
- Skips empty lines and comments
- Supports custom config path via flag

### Helper Functions
- `getBroadcastIP()` - Calculate broadcast address
- `getFirstUsableIP()` - First usable host IP (network + 1)
- `getLastUsableIP()` - Last usable host IP (broadcast - 1)
- `getTotalHosts()` - Total addresses in range
- `getUsableHosts()` - Usable hosts (total - 2)

## Installation & Distribution

- Published on GitHub: `github.com/trahma/cidr`
- Installable via: `go install github.com/trahma/cidr@latest`
- Binary name: `cidr`
- Installed to: `$GOPATH/bin` (typically `~/go/bin`)

## Future Considerations

### Potential Enhancements
- Subcommand structure (`cidr parse`, `cidr check`)
- Auto-detection of input type (IP vs CIDR)
- IPv6 support improvements
- JSON output format option
- Batch processing mode
- Additional network calculations

### Usability Questions
- Is the current flag-based approach optimal?
- Should we use subcommands instead?
- Should we auto-detect IP vs CIDR arguments?

## Development Workflow

1. Make changes to code
2. Test with: `go build -o cidr . && ./cidr [args]`
3. Update README if user-facing changes
4. Update this file if design decisions change
5. Commit with descriptive message
6. Push to GitHub

## Testing Commands

```bash
# Parse CIDR
./cidr 192.168.1.0/24

# Check IP against CIDR
./cidr 192.168.1.0/24 --check 192.168.1.50

# Check IP against config
./cidr --check 192.168.5.10

# Large network
./cidr 10.0.0.0/8

# Help
./cidr --help
```

## Notes

- Binary is gitignored (`.gitignore` excludes build artifacts)
- All colors use ANSI color codes via Lipgloss
- IPv4 is primary focus, IPv6 support is basic
- Config file format is intentionally simple
- No external dependencies beyond Cobra and Lipgloss
