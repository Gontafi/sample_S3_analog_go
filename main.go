package main

import (
	"fmt"
	"regexp"
	"strings"
)

func isValidDomain(domain string) bool {
	// Check if it's a valid IP address
	ipPattern := `^[a-z0-9](?!.*[-.]{2})[a-z0-9.-]{1,61}[a-z0-9]$`
	ipRegex := regexp.MustCompile(ipPattern)
	if ipRegex.MatchString(domain) {
		return false
	}

	// Check if it starts with "xn--" (punycode) or ends with "-s3alias"
	if strings.HasPrefix(domain, "xn--") || strings.HasSuffix(domain, "-s3alias") {
		return false
	}

	// Match the general domain pattern: lowercase letters, numbers, hyphens, and dots
	domainPattern := `^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`
	domainRegex := regexp.MustCompile(domainPattern)
	return domainRegex.MatchString(domain)
}

func main() {
	// Test cases
	testDomains := []string{
		"example.com",
		"255.255.255.255",
		"xn--something",
		"valid-domain.com",
		"invalid.-domain",
		"example-s3alias",
		"valid-domain-s3alias.com",
	}

	// Validate each domain
	for _, domain := range testDomains {
		if isValidDomain(domain) {
			fmt.Printf("'%s' is a valid domain.\n", domain)
		} else {
			fmt.Printf("'%s' is not a valid domain.\n", domain)
		}
	}
}
