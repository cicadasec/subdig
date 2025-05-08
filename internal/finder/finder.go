package finder

import (
	"fmt"
	"sync"
)

// FindSubdomains discovers subdomains for the given domain using multiple sources
func FindSubdomains(domain string) ([]string, error) {
	var allSubdomains []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Define sources to use
	sources := []SubdomainSource{
		&CrtShSource{},
		// Add more sources here as they are implemented
	}
	
	// Create a channel to collect errors
	errChan := make(chan error, len(sources))
	
	// Launch a goroutine for each source
	for _, source := range sources {
		wg.Add(1)
		go func(s SubdomainSource) {
			defer wg.Done()
			
			fmt.Printf("Searching subdomains in %s...\n", s.Name())
			subdomains, err := s.FindSubdomains(domain)
			if err != nil {
				errChan <- fmt.Errorf("error from %s: %v", s.Name(), err)
				return
			}
			
			mu.Lock()
			allSubdomains = append(allSubdomains, subdomains...)
			mu.Unlock()
			
			fmt.Printf("Found %d subdomains from %s\n", len(subdomains), s.Name())
		}(source)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)
	
	// Check if there were any errors
	for err := range errChan {
		return nil, err
	}
	
	// Remove duplicates
	uniqueSubdomains := removeDuplicates(allSubdomains)
	
	return uniqueSubdomains, nil
}

// removeDuplicates removes duplicate entries from a slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	
	return list
}