package httpi

import (
	"bytes"
	"io"
	"net/http"
)

type (
	// Interceptor is an http.RoundTripper.
	// Use SetRoundTripperFunc to define the desired output of a http.Client request.
	// Use SetRequestValidationFunc to define request validation logic.
	Interceptor struct {
		roundTripperFunc      RoundTripperFunc
		requestValidationFunc RequestValidationFunc
	}

	// RoundTripperFunc is the function used by Interceptor and defines the http.Response and error returned.
	RoundTripperFunc func(*http.Request) (*http.Response, error)
	// RequestValidationFunc is the function used by Interceptor and defines the request validation logic before executing RoundTripperFunc.
	RequestValidationFunc func(*http.Request) error
)

var (
	// DefaultRoundTripperFunc is the default RoundTripperFunc implementation. Returns a 200 OK response with a message.
	DefaultRoundTripperFunc = func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("Hello from the interceptor!"))),
		}, nil
	}
	// DefaultRequestValidationFunc is the default RequestValidationFunc implementation. Does no validation and returns nil error.
	DefaultRequestValidationFunc = func(r *http.Request) error {
		return nil
	}
)

// RoundTrip implements the http.RoundTripper interface.
func (i *Interceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := i.requestValidationFunc(req); err != nil {
		return nil, err
	}
	if i.roundTripperFunc != nil {
		return i.roundTripperFunc(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}

// New returns a new Interceptor that can be used as a http.RoundTripper.
// client := http.Client{Transport: httpi.New()}
func New(opts ...Option) *Interceptor {
	options := options{
		roundTripperFunc:      DefaultRoundTripperFunc,
		requestValidationFunc: DefaultRequestValidationFunc,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}
	return &Interceptor{
		roundTripperFunc:      options.roundTripperFunc,
		requestValidationFunc: options.requestValidationFunc,
	}
}

// SetRoundTripperFunc sets the http.Response and error to be returned by the client.
func (i *Interceptor) SetRoundTripperFunc(f RoundTripperFunc) {
	i.roundTripperFunc = f
}

// SetRequestValidationFunc sets the request validation function to be used by the client before making a request.
func (i *Interceptor) SetRequestValidationFunc(f RequestValidationFunc) {
	i.requestValidationFunc = f
}
