package cli

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

type TitleConfig struct {
	TitleFQDN  bool
	TitleColor string
	AtColor    string
	Bold       string
}

type TitleInfo struct {
	User     string
	Hostname string
	Title    string
	Length   int
}

func getTitle(config TitleConfig) (*TitleInfo, error) {
	// Get username
	username := os.Getenv("USER")
	if username == "" {
		currentUser, err := user.Current()
		if err == nil {
			username = currentUser.Username
		} else {
			// Fallback: extract from HOME path
			home := os.Getenv("HOME")
			if home != "" {
				parts := strings.Split(home, "/")
				if len(parts) > 0 {
					username = parts[len(parts)-1]
				}
			}
		}
	}

	// Get hostname
	var hostname string
	var err error

	if config.TitleFQDN {
		// Get fully qualified domain name
		hostname, err = os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("failed to get hostname: %w", err)
		}

		// For illumos/Solaris, append domain name if available
		// Note: Go doesn't have built-in domainname() equivalent
		// This is a simplified version
		domainname := os.Getenv("DOMAINNAME")
		if domainname != "" {
			hostname = hostname + "." + domainname
		}
	} else {
		// Get short hostname
		envHostname := os.Getenv("HOSTNAME")
		if envHostname != "" {
			hostname = envHostname
		} else {
			hostname, err = os.Hostname()
			if err != nil {
				return nil, fmt.Errorf("failed to get hostname: %w", err)
			}
		}

		// Trim to short hostname (before first dot)
		if idx := strings.Index(hostname, "."); idx != -1 {
			hostname = hostname[:idx]
		}
	}

	// Build title string
	title := fmt.Sprintf("%s%s%s%s@%s%s%s%s",
		config.TitleColor,
		config.Bold,
		username,
		config.AtColor,
		"",
		config.TitleColor,
		config.Bold,
		hostname,
	)

	// Calculate length (without color codes)
	length := len(username) + len(hostname) + 1

	return &TitleInfo{
		User:     username,
		Hostname: hostname,
		Title:    title,
		Length:   length,
	}, nil
}

// func main() {
// 	// Example usage
// 	config := TitleConfig{
// 		TitleFQDN:  false,
// 		TitleColor: "\033[1;34m", // Blue
// 		AtColor:    "\033[1;32m", // Green
// 		Bold:       "\033[1m",    // Bold
// 	}

// 	info, err := getTitle(config)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Printf("User: %s\n", info.User)
// 	fmt.Printf("Hostname: %s\n", info.Hostname)
// 	fmt.Printf("Title: %s\033[0m\n", info.Title) // Reset color at end
// 	fmt.Printf("Length: %d\n", info.Length)
// }

func GetTitle() string {
	config := TitleConfig{
		TitleFQDN:  false,
		TitleColor: "\033[1;34m", // Blue
		AtColor:    "\033[1;32m", // Green
		Bold:       "\033[1m",    // Bold
	}

	info, err := getTitle(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("User: %s\n", info.User)
	// fmt.Printf("Hostname: %s\n", info.Hostname)
	// fmt.Printf("Title: %s\033[0m\n", info.Title) // Reset color at end
	// fmt.Printf("Length: %d\n", info.Length)
	return fmt.Sprintf("%s\033[0m", info.Title)
}
