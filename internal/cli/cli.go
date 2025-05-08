package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yourusername/subdig/internal/finder"
	"github.com/yourusername/subdig/internal/resolver"
)

var (
	domain      string
	outputFile  string
	resolveFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "subdig",
	Short: "SubDig - Advanced Subdomain Enumeration Tool",
	Long: `SubDig is an advanced subdomain enumeration tool for penetration testing.
It performs passive subdomain discovery using various sources like crt.sh
and can optionally resolve discovered subdomains to check if they're alive.

Example usage:
  subdig -d example.com
  subdig -d example.com -r -o results.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" {
			color.Red("Error: domain is required")
			cmd.Help()
			os.Exit(1)
		}

		startTime := time.Now()
		color.Cyan("Starting subdomain enumeration for: %s", domain)
		
		// Find subdomains
		subdomains, err := finder.FindSubdomains(domain)
		if err != nil {
			color.Red("Error finding subdomains: %v", err)
			os.Exit(1)
		}
		
		color.Green("Found %d subdomains", len(subdomains))
		
		// Resolve subdomains if requested
		var aliveSubdomains []string
		if resolveFlag {
			color.Cyan("Resolving subdomains...")
			aliveSubdomains = resolver.ResolveSubdomains(subdomains)
			color.Green("Found %d alive subdomains", len(aliveSubdomains))
		}
		
		// Save results if output file is specified
		if outputFile != "" {
			var dataToSave []string
			if resolveFlag {
				dataToSave = aliveSubdomains
			} else {
				dataToSave = subdomains
			}
			
			err := saveToFile(dataToSave, outputFile)
			if err != nil {
				color.Red("Error saving to file: %v", err)
			} else {
				color.Green("Results saved to %s", outputFile)
			}
		}
		
		// Print results
		resultsToShow := subdomains
		if resolveFlag {
			resultsToShow = aliveSubdomains
		}
		
		fmt.Println("\nResults:")
		for _, sub := range resultsToShow {
			fmt.Println(sub)
		}
		
		elapsed := time.Since(startTime)
		color.Cyan("\nCompleted in %s", elapsed)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&domain, "domain", "d", "", "Target domain to find subdomains (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Save results to output file")
	rootCmd.Flags().BoolVarP(&resolveFlag, "resolve", "r", false, "Resolve discovered subdomains")
	
	rootCmd.MarkFlagRequired("domain")
}

// saveToFile saves the given subdomains to a file
func saveToFile(subdomains []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	for _, subdomain := range subdomains {
		if _, err := file.WriteString(subdomain + "\n"); err != nil {
			return err
		}
	}
	
	return nil
}