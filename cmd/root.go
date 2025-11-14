package cmd

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	checkIP    string
	configFile string

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)
)

var rootCmd = &cobra.Command{
	Use:   "cidr [CIDR notation]",
	Short: "A beautiful CIDR subnet parser",
	Long: titleStyle.Render("CIDR Parser") + "\n\n" +
		"Parse CIDR subnet masks and display human-readable IP ranges.\n" +
		"Check if an IP address belongs to a CIDR range.\n" +
		"Load default CIDRs from ~/.cidr file.",
	Example: `  cidr 192.168.1.0/24
  cidr 10.0.0.0/8 --check 10.5.3.2
  cidr --check 172.16.0.5`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCIDR,
}

func init() {
	rootCmd.Flags().StringVarP(&checkIP, "check", "c", "", "Check if an IP address is within the CIDR range")
	rootCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to .cidr config file (defaults to ~/.cidr)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render("Error: ")+err.Error())
		os.Exit(1)
	}
}

func runCIDR(cmd *cobra.Command, args []string) error {
	var cidrs []string
	var configPath string
	var configLoaded bool

	// If a CIDR is provided as argument, use it
	if len(args) > 0 {
		cidrs = append(cidrs, args[0])
	}

	// Load CIDRs from config file if no argument provided or if checking an IP
	if len(cidrs) == 0 || checkIP != "" {
		configCIDRs, path, err := loadConfigCIDRs()
		if err == nil {
			cidrs = append(cidrs, configCIDRs...)
			configPath = path
			configLoaded = true
		} else if len(cidrs) == 0 {
			return fmt.Errorf("no CIDR provided and could not load config file: %w", err)
		}
	}

	if len(cidrs) == 0 {
		return fmt.Errorf("please provide a CIDR notation or create a ~/.cidr file with CIDR ranges")
	}

	// Show config file indicator if loaded
	if configLoaded {
		fmt.Println(dimStyle.Render(fmt.Sprintf("Using config from: %s", configPath)))
		fmt.Println()
	}

	// If checking an IP, validate and check against CIDRs
	if checkIP != "" {
		if err := checkIPInCIDRs(checkIP, cidrs); err != nil {
			return err
		}
	} else {
		// Otherwise, display CIDR information
		for i, cidr := range cidrs {
			if i > 0 {
				fmt.Println() // Separator between multiple CIDRs
			}
			if err := displayCIDRInfo(cidr); err != nil {
				return err
			}
		}
	}

	// Show help hint once at the end
	fmt.Println()
	fmt.Println(helpStyle.Render("Run 'cidr --help' for more options"))

	return nil
}

func displayCIDRInfo(cidrStr string) error {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return fmt.Errorf("invalid CIDR notation '%s': %w", cidrStr, err)
	}

	// Get network details
	networkIP := ipnet.IP
	broadcastIP := getBroadcastIP(ipnet)
	firstIP := getFirstUsableIP(ipnet)
	lastIP := getLastUsableIP(ipnet)
	totalHosts := getTotalHosts(ipnet)
	usableHosts := getUsableHosts(ipnet)

	// Get subnet mask
	mask := net.IP(ipnet.Mask)

	// Display information
	fmt.Println(titleStyle.Render("CIDR Information"))
	fmt.Printf("%s %s\n", labelStyle.Render("CIDR:"), valueStyle.Render(cidrStr))
	fmt.Printf("%s %s\n", labelStyle.Render("Network Address:"), valueStyle.Render(networkIP.String()))
	fmt.Printf("%s %s\n", labelStyle.Render("Subnet Mask:"), valueStyle.Render(mask.String()))
	fmt.Printf("%s %s\n", labelStyle.Render("Broadcast Address:"), valueStyle.Render(broadcastIP.String()))
	fmt.Println()
	fmt.Printf("%s %s - %s\n", labelStyle.Render("IP Range:"), valueStyle.Render(networkIP.String()), valueStyle.Render(broadcastIP.String()))
	fmt.Printf("%s %s - %s\n", labelStyle.Render("Usable IPs:"), valueStyle.Render(firstIP.String()), valueStyle.Render(lastIP.String()))
	fmt.Println()
	fmt.Printf("%s %s\n", labelStyle.Render("Total Hosts:"), valueStyle.Render(fmt.Sprintf("%d", totalHosts)))
	fmt.Printf("%s %s\n", labelStyle.Render("Usable Hosts:"), valueStyle.Render(fmt.Sprintf("%d", usableHosts)))

	return nil
}

func checkIPInCIDRs(ipStr string, cidrs []string) error {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	fmt.Println(titleStyle.Render("IP Address Check"))
	fmt.Printf("%s %s\n\n", labelStyle.Render("Checking IP:"), valueStyle.Render(ipStr))

	found := false
	for _, cidrStr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			fmt.Printf("%s Invalid CIDR: %s\n", errorStyle.Render("✗"), cidrStr)
			continue
		}

		if ipnet.Contains(ip) {
			fmt.Printf("%s IP is in %s\n", successStyle.Render("✓"), valueStyle.Render(cidrStr))
			found = true
		} else {
			fmt.Printf("%s IP is not in %s\n", infoStyle.Render("○"), cidrStr)
		}
	}

	fmt.Println()
	if found {
		fmt.Println(successStyle.Render("IP address found in one or more CIDR ranges"))
	} else {
		fmt.Println(errorStyle.Render("IP address not found in any CIDR ranges"))
	}

	return nil
}

func loadConfigCIDRs() ([]string, string, error) {
	var configPath string
	if configFile != "" {
		configPath = configFile
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, "", err
		}
		configPath = filepath.Join(home, ".cidr")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", err
	}

	lines := strings.Split(string(data), "\n")
	var cidrs []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cidrs = append(cidrs, line)
	}

	return cidrs, configPath, nil
}

// Helper functions for IP calculations

func getBroadcastIP(ipnet *net.IPNet) net.IP {
	ip := ipnet.IP.To4()
	if ip == nil {
		ip = ipnet.IP.To16()
	}

	broadcast := make(net.IP, len(ip))
	for i := range ip {
		broadcast[i] = ip[i] | ^ipnet.Mask[i]
	}
	return broadcast
}

func getFirstUsableIP(ipnet *net.IPNet) net.IP {
	ip := ipnet.IP.To4()
	if ip == nil {
		// IPv6
		return ipnet.IP
	}

	// IPv4: first usable is network + 1
	first := make(net.IP, len(ip))
	copy(first, ip)

	// Increment IP by 1
	for i := len(first) - 1; i >= 0; i-- {
		first[i]++
		if first[i] > 0 {
			break
		}
	}

	return first
}

func getLastUsableIP(ipnet *net.IPNet) net.IP {
	broadcast := getBroadcastIP(ipnet)

	if broadcast.To4() == nil {
		// IPv6
		return broadcast
	}

	// IPv4: last usable is broadcast - 1
	last := make(net.IP, len(broadcast))
	copy(last, broadcast)

	// Decrement IP by 1
	for i := len(last) - 1; i >= 0; i-- {
		last[i]--
		if last[i] < 255 {
			break
		}
	}

	return last
}

func getTotalHosts(ipnet *net.IPNet) uint64 {
	ones, bits := ipnet.Mask.Size()
	return 1 << uint(bits-ones)
}

func getUsableHosts(ipnet *net.IPNet) uint64 {
	total := getTotalHosts(ipnet)
	if total <= 2 {
		return 0
	}
	return total - 2 // Subtract network and broadcast addresses
}
