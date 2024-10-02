package httpx

import (
	"net/http"
	"testing"
)

func TestExtractMimeType(t *testing.T) {
	testCases := []struct {
		name         string
		headers      http.Header
		expectedMIME string
	}{
		{
			name: "Valid text/html with charset",
			headers: http.Header{
				"Content-Type": []string{"text/html; charset=utf-8"},
			},
			expectedMIME: "text/html",
		},
		{
			name: "Valid application/json",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			expectedMIME: "application/json",
		},
		{
			name: "Multiple Content-Type headers (takes first)",
			headers: http.Header{
				"Content-Type": []string{"text/plain", "text/html"},
			},
			expectedMIME: "text/plain",
		},
		{
			name: "Content-Type with extra whitespace",
			headers: http.Header{
				"Content-Type": []string{"  text/xml ; charset=utf-8  "},
			},
			expectedMIME: "text/xml",
		},
		{
			name:         "Missing Content-Type header",
			headers:      http.Header{},
			expectedMIME: "",
		},
		{
			name: "Empty Content-Type header",
			headers: http.Header{
				"Content-Type": []string{""},
			},
			expectedMIME: "",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				mimeType := ExtractMimeType(tc.headers)
				if mimeType != tc.expectedMIME {
					t.Errorf("Expected MIME type '%s', but got '%s'", tc.expectedMIME, mimeType)
				}
			},
		)
	}
}

func TestExtractContentLength(t *testing.T) {
	testCases := []struct {
		name           string
		headers        http.Header
		expectedLength int
		expectedErr    bool
	}{
		{
			name: "Valid Content-Length",
			headers: http.Header{
				"Content-Length": []string{"1234"},
			},
			expectedLength: 1234,
		},
		{
			name: "Content-Length with leading/trailing whitespace",
			headers: http.Header{
				"Content-Length": []string{"  5678  "},
			},
			expectedLength: 5678,
		},
		{
			name: "Multiple Content-Length headers (takes first)",
			headers: http.Header{
				"Content-Length": []string{"1000", "2000"},
			},
			expectedLength: 1000,
		},
		{
			name:           "Missing Content-Length header",
			headers:        http.Header{},
			expectedLength: -1,
		},
		{
			name: "Empty Content-Length header",
			headers: http.Header{
				"Content-Length": []string{""},
			},
			expectedLength: -1,
		},
		{
			name: "Invalid Content-Length (non-numeric)",
			headers: http.Header{
				"Content-Length": []string{"abc"},
			},
			expectedLength: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				contentLength := ExtractContentLength(tc.headers)
				if contentLength != int64(tc.expectedLength) {
					t.Errorf("Expected Content-Length %d, but got %d", tc.expectedLength, contentLength)
				}
			},
		)
	}
}
