package junglebus

import (
	"net/http"

	"github.com/GorillaPool/go-junglebus/transports"
)

// WithHTTP will overwrite the default server url (junglebus.gorillapool.io)
func WithHTTP(serverURL string) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithHTTP(serverURL))
		}
	}
}

// WithHTTPClient will overwrite the default client with a custom client
func WithHTTPClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithHTTPClient(serverURL, httpClient))
		}
	}
}

// WithToken will set the token to use in all requests
func WithToken(token string) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithToken(token))
		}
	}
}

// WithDebugging will set whether to turn debugging on
func WithDebugging(debug bool) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithDebugging(debug))
		}
	}
}

// WithSSL will set whether to use SSL in all communications or not
func WithSSL(useSSL bool) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithSSL(useSSL))
		}
	}
}

// WithVersion will set the API version to use (v1 is default)
func WithVersion(version string) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithVersion(version))
		}
	}
}
