# UCSB Experimental Cosmology Lab's Electronics Lab Organizer
This is an reimplementation of UCSB Experimental Cosmology Lab's Electronics Lab Inventory system in Go. It is designed to organize, catalog, and sort through decades of electronic components. It provides functionality to process part numbers, interact with APIs, handle OAuth authentication, generate labels, and manage the entire inventory system.

## Dependencies
* Go v1.18+ 
* Python 3.10+
* `blabel` [available here](https://github.com/Edinburgh-Genome-Foundry/blabel)
* InvenTree installed and configured
* A serverless function that handles the initial OAuth authorization flow ([mine](https://github.com/davidwuya/oauth-callback) is hosted on Vercel and uses Vercel KV to securely store API keys)
* Digi-Key API access

## Configurations
You'll need to set up the necessary constants in the `const.go` file. This file contains sensitive information like API keys and tokens.
```go
package main

const (
	VercelURL               = "VERCEL_OAUTH_HANDLER"
	DKAuthorize             = "https://api.digikey.com/v1/oauth2/authorize"
	DKProductDetailEndpoint = "https://api.digikey.com/Search/v3/Products/"
	InvenTreeServerAddress  = "http://SERVER_ADDR/api/"
	CLIENT_ID               = "DIGIKEY_CLIENT_ID"
	APIKey                  = "VERCEL_OAUTH_HANDLER_KEY"
	OAuthState              = "OAUTH_STATE"
	INVENTREE_TOKEN         = "INVENTREE_TOKEN"
)
```
## Todo
- [] Test coverage with Go Tests and Actions
- [] A better, more intuitive TUI with [Bubble Tea](https://github.com/charmbracelet/bubbletea)

