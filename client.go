package httpi

import "net/http"

// Client is a convenience function that returns an http.Client with an Interceptor as its Transport.
func NewClient(opts ...Option) *http.Client {
	return &http.Client{
		Transport: NewTransport(opts...),
	}
}

// SetRoundTripperFunc sets the http.Response and error to be returned by the client.
// Note: This function will panic if the clients RoundTripper is not an Interceptor.
func SetRoundTripperFunc(client *http.Client, f func(*http.Request) (*http.Response, error)) {
	client.Transport.(*Transport).SetRoundTripperFunc(f)
}

// SetRequestValidationFunc sets the request validation function to be used by the client before making a request.
// Note: This function will panic if the clients RoundTripper is not an Interceptor.
func SetRequestValidationFunc(client *http.Client, f func(*http.Request) error) {
	client.Transport.(*Transport).SetRequestValidationFunc(f)
}
