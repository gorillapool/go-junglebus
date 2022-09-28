package transports

import (
	"net/http"
	"regexp"
)

var regexHTTP = regexp.MustCompile(`^http://`)
var regexWS = regexp.MustCompile(`^ws://`)
var regexReplaceHTTPS = regexp.MustCompile(`^https?://`)
var regexReplaceWSS = regexp.MustCompile(`^wss?://`)

// WithHTTP will overwrite the default client with a custom client
func WithHTTP(serverURL string) ClientOps {
	return func(c *Client) {
		if c != nil {
			initHTTPTransport(c, serverURL, &http.Client{})
		}
	}
}

// WithHTTPClient will overwrite the default client with a custom client
func WithHTTPClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *Client) {
		if c != nil {
			initHTTPTransport(c, serverURL, httpClient)
		}
	}
}

func initHTTPTransport(c *Client, serverURL string, httpClient *http.Client) {
	useSSL := true
	if regexHTTP.MatchString(serverURL) || regexWS.MatchString(serverURL) {
		useSSL = false // turn off SSL if server url contains http:// or ws://
	}

	// remove prefix if applicable
	serverURL = regexReplaceHTTPS.ReplaceAllString(serverURL, "")
	serverURL = regexReplaceWSS.ReplaceAllString(serverURL, "")

	c.transport = NewTransportService(&TransportHTTP{
		debug:      c.debug,
		server:     serverURL,
		httpClient: httpClient,
		useSSL:     useSSL,
		version:    "v1",
	})
}

// WithToken will set the token to use in all requests
func WithToken(token string) ClientOps {
	return func(c *Client) {
		if c != nil {
			if c.transport != nil {
				c.transport.SetToken(token)
			}
		}
	}
}

// WithDebugging will set whether to turn debugging on
func WithDebugging(debug bool) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.debug = debug
			if c.transport != nil {
				c.transport.SetDebug(debug)
			}
		}
	}
}

// WithSSL will set whether to use SSL in all communications
func WithSSL(useSSL bool) ClientOps {
	return func(c *Client) {
		if c != nil {
			if c.transport != nil {
				c.transport.UseSSL(useSSL)
			}
		}
	}
}

// WithVersion will set the version of the API to use
func WithVersion(version string) ClientOps {
	return func(c *Client) {
		if c != nil {
			if c.transport != nil {
				c.transport.SetVersion(version)
			}
		}
	}
}
