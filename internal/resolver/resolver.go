package resolver

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ResolveSubdomains checks if the provided subdomains are alive by resolving them
func ResolveSubdomains(subdomains []string) []string {
	var aliveSubdomains []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Use a semaphore to limit concurrent DNS lookups
	semaphore := make(chan struct{}, 50)
	
	for _, subdomain := range subdomains {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		
		go func(sub string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore
			
			if isAlive(sub) {
				mu.Lock()
				aliveSubdomains = append(aliveSubdomains, sub)
				mu.Unlock()
				fmt.Printf("✓ %s is alive\n", sub)
			} else {
				fmt.Printf("✗ %s is not alive\n", sub)
			}
		}(subdomain)
	}
	
	wg.Wait()
	return aliveSubdomains
}

// isAlive checks if a subdomain is alive by attempting to resolve it
func isAlive(subdomain string) bool {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx net.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return dialer.DialContext(ctx, network, "8.8.8.8:53")
		},
	}
	
	_, err := resolver.LookupHost(context.Background(), subdomain)
	return err == nil
}