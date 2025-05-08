# SubDig - Advanced Subdomain Enumeration Tool

SubDig is an advanced subdomain enumeration tool for penetration testing written in Go. It performs passive subdomain discovery using various sources and can optionally resolve discovered subdomains to check if they're alive.

## Features

- Passive subdomain enumeration using multiple sources (currently crt.sh)
- Optional DNS resolution to check if subdomains are alive
- Save results to a file
- Clean and colorful CLI output
- Fast and efficient with concurrent operations

## Installation
```bash
go install github.com/cicadasec/subdig/cmd/subdig@latest
```
### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/subdig.git
cd subdig

# Build the binary
go build -o subdig ./cmd/subdig

# Move to a directory in your PATH (optional)
sudo mv subdig /usr/local/bin/
