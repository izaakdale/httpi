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
  return &http.Response{
    StatusCode: 200,
    Body:       io.NopCloser(bytes.NewReader(body)),
  }, nil
})
resp, err := cli.Get(url)
```

There is also the option to skip Transport entirely and just get a Client.
```go
cli := httpi.NewClient()
httpi.SetRoundTripperFunc(cli, someRoundTripperFunc)
```

If you are more interested in the errors that are returned you can use RequestValidationFunc.
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
transport := httpi.NewTransport(WithRoundTripperFunc(someRoundTripFunc), WithRequestValidationFunc(someValidationFunc))
cli := &http.Client{Transport: transport}

cli = httpi.NewClient(WithRoundTripperFunc(someFunc), WithRequestValidationFunc(someValidationFunc))
```
