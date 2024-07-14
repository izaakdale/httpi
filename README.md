## httpi (http interceptor)

```sh
go get github.com/izaakdale/httpi
```

### A RoundTripper package that is designed for mocking http.Client calls.

Get a Transport (implements http.RoundTripper) and use SetRoundTripperFunc to define the response and error retuned by client requests.
```go
transport := httpi.NewTransport()
cli := &http.Client{Transport: transport}
transport.SetRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
	if r.URL.Path != "/hello" {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewReader([]byte("Not Found"))),
		}, nil
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}, nil
})
resp, err := cli.Get("http://example.com/hello")
```

There is also the option to skip initializing Transport yourself and just get a preloaded Client.
```go
cli := httpi.NewClient()
httpi.SetRoundTripperFunc(cli, someRoundTripperFunc)
```

If you are more interested in returning errors based on the request you can use RequestValidationFunc.
```go
transport := httpi.NewTransport()
cli := &http.Client{Transport: transport}

transport.SetRequestValidationFunc(func(_ *http.Request) error {
  return errTest
})

cli = httpi.NewClient()
httpi.SetRoundTripperFunc(cli, someRoundTripperFunc)

httpi.SetRequestValidationFunc(cli, func(r *http.Request) error {
  if r.URL.Scheme != "https" {
    return errors.New("invalid scheme")
  }
  return nil
})

_, err = http.Get(url)
```

Or use the WithOptions
```go
transport := httpi.NewTransport(
	WithRoundTripperFunc(someRoundTripFunc),
	WithRequestValidationFunc(someValidationFunc),
)
cli := &http.Client{Transport: transport}

cli = httpi.NewClient(
	WithRoundTripperFunc(someFunc),
	WithRequestValidationFunc(someValidationFunc),
)
```
