package httpi

import "net/http"

// Client is a convenience function that returns an http.Client with an Interceptor as its Transport.
func NewClient(opts ...Option) *http.Client {
	return &http.Client{
		Transport: New(opts...),
	}
}

// SetRoundTripperFunc sets the http.Response and error to be returned by the client.
// Note: This function will panic if the clients RoundTripper is not an Interceptor.
func SetRoundTripperFunc(client *http.Client, f func(*http.Request) (*http.Response, error)) func() {
	reset := client.Transport.(*Interceptor).SetRoundTripperFunc(f)
	return reset
}

// SetRequestValidationFunc sets the request validation function to be used by the client before making a request.
func SetRequestValidationFunc(client *http.Client, f func(*http.Request) error) func() {
	reset := client.Transport.(*Interceptor).SetRequestValidationFunc(f)
	return reset
}
