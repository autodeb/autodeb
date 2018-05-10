// Package apiclient implements a client for the autodeb-server REST API
package apiclient

//APIClient is a client for the autodeb-server REST API
type APIClient struct {
	serverAddress string
	serverPort    int
}

//New creates a new APIClient
func New(serverAddress string, serverPort int) *APIClient {
	apiClient := &APIClient{
		serverAddress: serverAddress,
		serverPort:    serverPort,
	}

	return apiClient
}
