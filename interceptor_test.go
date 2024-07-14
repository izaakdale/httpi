package httpi_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/izaakdale/httpi"
)

var (
	url     = "http://example.com"
	body    = []byte("test body")
	errTest = errors.New("test error")
)

func TestDefaultInterceptor(t *testing.T) {
	inctr := httpi.New()
	cli := &http.Client{Transport: inctr}

	resp, err := cli.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(respBytes), "Hello from the interceptor!") {
		t.Fatalf("expected Hello from the interceptor!, got %s", respBytes)
	}

	// Test that default http.Transport is used when roundTripperFunc is nil
	defaultStub := httpi.New(
		httpi.WithRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusTeapot,
				Body:       io.NopCloser(strings.NewReader("not the default interceptor!")),
			}, nil
		}),
	)
	http.DefaultTransport = defaultStub
	inctr.SetRoundTripperFunc(nil)

	resp, err = cli.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTeapot {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(respBytes), "not the default interceptor!") {
		t.Fatalf("expected not the default interceptor!, got %s", respBytes)
	}
}

func TestInterceptorSetRoundTripperFunc(t *testing.T) {
	inctr := httpi.New()
	cli := &http.Client{Transport: inctr}

	t.Run("custom body", func(t *testing.T) {
		reset := inctr.SetRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		})
		defer reset()

		resp, err := cli.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(respBytes, body) {
			t.Fatalf("expected %s, got %s", body, respBytes)
		}
	})

	t.Run("other status code", func(t *testing.T) {
		reset := inctr.SetRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		})
		defer reset()

		resp, err := cli.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(respBytes, body) {
			t.Fatalf("expected %s, got %s", body, respBytes)
		}
	})

	t.Run("error", func(t *testing.T) {
		reset := inctr.SetRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return nil, errTest
		})
		defer reset()

		resp, err := cli.Get(url)
		if !errors.Is(err, errTest) {
			t.Fatalf("http.ErrLineTooLong, got %v", err)
		}
		if resp != nil {
			t.Fatalf("expected nil, got %v", resp)
		}
	})
}

func TestInterceptorSetRequestValidation(t *testing.T) {
	inctr := httpi.New()
	cli := &http.Client{Transport: inctr}

	t.Run("valid request", func(t *testing.T) {
		reset := inctr.SetRequestValidationFunc(func(r *http.Request) error {
			return nil
		})
		defer reset()

		resp, err := cli.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("return error", func(t *testing.T) {
		reset := inctr.SetRequestValidationFunc(func(_ *http.Request) error {
			return errTest
		})
		defer reset()

		resp, err := cli.Get(url)
		if !errors.Is(err, errTest) {
			t.Fatalf("http.ErrLineTooLong, got %v", err)
		}
		if resp != nil {
			t.Fatalf("expected nil, got %v", resp)
		}
	})

	t.Run("validate request", func(t *testing.T) {
		reset := inctr.SetRequestValidationFunc(func(r *http.Request) error {
			if r.URL.Scheme != "https" {
				return errors.New("invalid scheme")
			}
			return nil
		})
		defer reset()

		_, err := cli.Get("https://hypertexttransferprotocolsecure.com")
		if err != nil {
			t.Fatal(err)
		}

		_, err = cli.Get("http://hypertexttransferprotocol.com")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	// Should be reset back to default.
	resp, err := cli.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(respBytes), "Hello from the interceptor!") {
		t.Fatalf("expected Hello from the interceptor!, got %s", respBytes)
	}
}
