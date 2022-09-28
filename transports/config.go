package transports

// TransportType the type of transport being used (http or graphql)
type TransportType string

// JungleBusUserAgent the user agent sent to the jungle bus server
const JungleBusUserAgent = "JungleBus: go-client " + JungleBusClientVersion

// JungleBusClientVersion is the version of the client
const JungleBusClientVersion = "v0.1.0"

const (
	// JungleBusTransportHTTP uses the http transport for all jungle bus server actions
	JungleBusTransportHTTP TransportType = "http"

	// JungleBusTransportMock uses the mock transport for all jungle bus server actions
	JungleBusTransportMock TransportType = "mock"
)

const (
	// FieldID is the id field for most models
	FieldID = "id"

	// FieldUsername is the username field for login
	FieldUsername = "username"

	// FieldPassword is the password field for login
	FieldPassword = "password"

	// FieldSubscriptionID is the subscription id field
	FieldSubscriptionID = "id"

	// FieldUserAgent is the field for storing the user agent
	FieldUserAgent = "user_agent"
)
