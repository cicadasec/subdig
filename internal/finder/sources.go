package finder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// SubdomainSource defines the interface for subdomain enumeration sources
type SubdomainSource interface {
	Name() string
	FindSubdomains(domain string) ([]string, error)
}

// CrtShSource implements the crt.sh certificate transparency source
type CrtShSource struct{}

// Name returns the name of the source
func (s *CrtShSource) Name() string {
	return "crt.sh"
}

// CrtShResponse represents the JSON response from crt.sh
type CrtShResponse []struct {
	NameValue string `json:"name_value"`
}

// FindSubdomains discovers subdomains using crt.sh
func (s *CrtShSource) FindSubdomains(domain string) ([]string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var results CrtShResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}
	
	var subdomains []string
	for _, result := range results {
		// Some crt.sh entries have multiple subdomains separated by newlines
		for _, subdomain := range strings.Split(result.NameValue, "\n") {
			// Remove wildcard prefix if present
			subdomain = strings.TrimPrefix(subdomain, "*.")
			
			// Ensure the subdomain belongs to the requested domain
			if strings.HasSuffix(subdomain, "."+domain) || subdomain == domain {
				subdomains = append(subdomains, subdomain)
			}
		}
	}
	
	return subdomains, nil
}