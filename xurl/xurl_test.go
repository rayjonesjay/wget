package xurl

import "testing"

// IsValidUrl Test
func TestIsValidURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// Valid URLs
		{
			name: "Valid HTTP URL",
			args: args{url: "http://example.com"},
			want: true,
		},
		{
			name: "Valid HTTPS URL",
			args: args{url: "https://example.com"},
			want: true,
		},
		{
			name: "Valid URL with path",
			args: args{url: "https://example.com/path/to/resource"},
			want: true,
		},
		{
			name: "Valid URL with query",
			args: args{url: "https://example.com/search?q=golang"},
			want: true,
		},

		// Invalid URLs
		{
			name: "Invalid URL without scheme",
			args: args{url: "example.com"},
			want: false,
		},
		{
			name: "Invalid URL with invalid scheme",
			args: args{url: "ftp://example.com"},
			want: false,
		},
		{
			name: "Invalid URL with missing domain",
			args: args{url: "https://.com"},
			want: false,
		},
		{
			name: "Invalid URL with leading dot in domain",
			args: args{url: "https://.example.com"},
			want: false,
		},
		{
			name: "Invalid URL with missing TLD",
			args: args{url: "https://example"},
			want: false,
		},
		{
			name: "Invalid URL with leading hyphen in domain",
			args: args{url: "https://-example.com"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got, _ := IsValidURL(tt.args.url); got != tt.want {
					t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
