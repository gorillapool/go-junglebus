package transports

// Client is the transport client
type Client struct {
	debug     bool
	transport TransportService
}

// ClientOps are the client options functions
type ClientOps func(c *Client)

// NewTransport create a new transport service object
func NewTransport(opts ...ClientOps) (TransportService, error) {
	client := Client{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.transport == nil {
		return nil, ErrNoClientSet
	}

	return client.transport, nil
}

// NewTransportService create a new transport service interface
func NewTransportService(transportService TransportService) TransportService {
	return transportService
}
