package xurl_test

import (
	"testing"

	"wget/xerr"
	"wget/xurl"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name    string
		rawUrl  string
		want    bool
		wantErr error
	}{
		// Valid URL cases
		{
			name:    "Valid HTTP URL",
			rawUrl:  "http://www.google.com",
			want:    true,
			wantErr: nil,
		},
		{
			name:    "Valid HTTPS URL",
			rawUrl:  "https://www.google.com",
			want:    true,
			wantErr: nil,
		},
		{
			name:    "Valid URL with path",
			rawUrl:  "https://www.example.com/cat.png",
			want:    true,
			wantErr: nil,
		},
		{
			name:    "Valid localhost URL",
			rawUrl:  "http://localhost",
			want:    true,
			wantErr: nil,
		},

		// Invalid URL cases
		{
			name:    "Empty URL",
			rawUrl:  "",
			want:    false,
			wantErr: xerr.ErrNotAbsolute,
		},
		{
			name:    "No scheme, no host",
			rawUrl:  "/path/to/resource",
			want:    false,
			wantErr: xerr.ErrNotAbsolute,
		},
		{
			name:    "Invalid scheme",
			rawUrl:  "ftp://example.com",
			want:    false,
			wantErr: xerr.ErrWrongScheme,
		},
		{
			name:    "Invalid host - starts with dot",
			rawUrl:  "https://.example.com",
			want:    false,
			wantErr: xerr.ErrWrongHostFormat,
		},
		{
			name:    "Missing host",
			rawUrl:  "https://",
			want:    false,
			wantErr: xerr.ErrEmptyHostName,
		},
		{
			name:    "Invalid domain format - no dot",
			rawUrl:  "https://example",
			want:    false,
			wantErr: xerr.ErrInvalidDomainFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xurl.IsValidURL(tt.rawUrl)

			// Check if the result matches the expected value
			if got != tt.want {
				t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
			}

			// Check if the error matches the expected error
			if (err != nil && tt.wantErr == nil) || (err == nil && tt.wantErr != nil) || (err != nil && err.Error() != tt.wantErr.Error()) {
				t.Errorf("IsValidURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTryFixScheme(t *testing.T) {
	tests := []struct {
		name    string
		rawUrl  string
		want    bool
		wantErr bool
	}{
		// Test cases that need the scheme to be added
		{
			name:    "No scheme - should fix with http",
			rawUrl:  "www.example.com",
			want:    true,
			wantErr: false,
		},
		{
			name:    "No scheme - should fix with https",
			rawUrl:  "google.com",
			want:    true,
			wantErr: false,
		},
		{
			name:    "No scheme - invalid domain",
			rawUrl:  "localhost",
			want:    true,
			wantErr: false,
		},

		// Already valid URL (no fixing needed)
		{
			name:    "Already valid HTTP URL",
			rawUrl:  "http://example.com",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Already valid HTTPS URL",
			rawUrl:  "https://example.com",
			want:    true,
			wantErr: false,
		},

		// Invalid cases
		{
			name:    "Invalid URL after fixing - no host",
			rawUrl:  "/path/to/resource",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid URL after fixing - malformed URL",
			rawUrl:  ".example.com",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xurl.TryFixScheme(tt.rawUrl)

			// Check if the result matches the expected value
			if got != tt.want {
				t.Errorf("TryFixScheme() = %v, want %v", got, tt.want)
			}

			// Check if the error presence matches the expected error presence
			if (err != nil) != tt.wantErr {
				t.Errorf("TryFixScheme() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
