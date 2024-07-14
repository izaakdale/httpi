package httpi_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/izaakdale/httpi"
)

func TestInterceptorOptions(t *testing.T) {
	inctr := httpi.NewTransport(
		httpi.WithRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		}),
		httpi.WithRequestValidationFunc(func(r *http.Request) error {
			if r.URL.Scheme != "https" {
				return errTest
			}
			return nil
		}),
	)
	cli := &http.Client{Transport: inctr}
	resp, err := cli.Get("https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(respBytes, body) {
		t.Fatalf("expected %s, got %s", body, respBytes)
	}

	_, err = cli.Get("http://example.com")
	if err == nil {
		t.Fatal(err)
	}
	if !errors.Is(err, errTest) {
		t.Fatalf("expected %s, got %s", errTest, err)
	}
}

func TestClientOptions(t *testing.T) {
	cli := httpi.NewClient(
		httpi.WithRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		}),
		httpi.WithRequestValidationFunc(func(r *http.Request) error {
			if r.URL.Scheme != "https" {
				return errTest
			}
			return nil
		}),
	)
	resp, err := cli.Get("https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(respBytes, body) {
		t.Fatalf("expected %s, got %s", body, respBytes)
	}

	_, err = cli.Get("http://example.com")
	if err == nil {
		t.Fatal(err)
	}
	if !errors.Is(err, errTest) {
		t.Fatalf("expected %s, got %s", errTest, err)
	}
}
