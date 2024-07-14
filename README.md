## httpi

```sh
go get github.com/izaakdale/httpi
```

### A RoundTripper library that is primarily designed mocking/stubbing http.Client calls.

Get an Interceptor and use SetRoundTripperFunc to define the response and error retuned by client requests.
```go
inctr := httpi.New()
cli := &http.Client{Transport: inctr}
inctr.SetRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
  return &http.Response{
    StatusCode: 200,
    Body:       io.NopCloser(bytes.NewReader(body)),
  }, nil
})
resp, err := cli.Get(url)
```

There is also the option to skip Inteceptor entirely and just get a http.Client
```go
cli := httpi.NewClient()
httpi.SetRoundTripperFunc(cli, someRoundTripperFunc)
```

If you are more interested in the errors that are returned you can use RequestValidationFunc
```go
inctr := httpi.New()
cli := &http.Client{Transport: inctr}

inctr.SetRequestValidationFunc(func(_ *http.Request) error {
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
inctr := httpi.New(WithRoundTripperFunc(someRoundTripFunc), WithSetRequestValidationFunc(someValidationFunc))
cli := &http.Client{Transport: inctr}

cli = httpi.NewClient(WithRoundTripperFunc(someFunc), WithSetRequestValidationFunc(someValidationFunc))
```
