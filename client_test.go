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

func TestDefaultClient(t *testing.T) {
	cli := httpi.NewClient()
	resp, err := cli.Get("http://example.com")
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

func TestClientSetRoundTripperFunc(t *testing.T) {
	cli := httpi.NewClient()

	t.Run("custom body", func(t *testing.T) {
		reset := httpi.SetRoundTripperFunc(cli, func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		})
		defer reset()
		resp, err := cli.Get("http://example.com")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
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

	t.Run("custom status code", func(t *testing.T) {
		reset := httpi.SetRoundTripperFunc(cli, func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		})
		defer reset()
		resp, err := cli.Get("http://example.com")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
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

	t.Run("custom error", func(t *testing.T) {
		reset := httpi.SetRoundTripperFunc(cli, func(r *http.Request) (*http.Response, error) {
			return nil, errTest
		})
		defer reset()
		_, err := cli.Get("http://example.com")
		if !errors.Is(err, errTest) {
			t.Fatalf("expected %v, got %v", http.ErrNotSupported, err)
		}
	})
}

func TestClientSetRequestValidationFunc(t *testing.T) {
	cli := httpi.NewClient()

	t.Run("valid request", func(t *testing.T) {
		reset := httpi.SetRequestValidationFunc(cli, func(r *http.Request) error {
			return nil
		})
		defer reset()
		resp, err := cli.Get("http://example.com")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(respBytes), "Hello from the interceptor!") {
			t.Fatalf("expected Hello from the interceptor!, got %s", respBytes)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		reset := httpi.SetRequestValidationFunc(cli, func(r *http.Request) error {
			return errTest
		})
		defer reset()
		_, err := cli.Get("http://example.com")
		if !errors.Is(err, errTest) {
			t.Fatalf("expected %v, got %v", errTest, err)
		}
	})

	t.Run("validate request", func(t *testing.T) {
		reset := httpi.SetRequestValidationFunc(cli, func(r *http.Request) error {
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
	if resp.StatusCode != http.StatusOK {
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
