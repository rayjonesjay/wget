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
			name: "Valid URL without scheme",
			args: args{url: "example.com"},
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
		{
			name: "Valid URL with port 80",
			args: args{url: "example.com:80"},
			want: true,
		},
		{
			name: "Valid URL with port 443",
			args: args{url: "example.com:443"},
			want: true,
		},
		{
			name: "Valid localhost URL",
			args: args{url: "localhost"},
			want: true,
		},
		{
			name: "Valid 127.0.0.1 URL",
			args: args{url: "127.0.0.1"},
			want: true,
		},
		{
			name: "Invalid URL with other port",
			args: args{url: "example.com:8080"},
			want: true,
		},

		// Invalid URLs
		{
			name: "Invalid URL without scheme",
			args: args{url: "//example.com/page"},
			want: false,
		},
		{
			name: "Missing scheme but space in URL",
			args: args{url: "example .com"},
			want: false,
		},
		{
			name: "Relative URL",
			args: args{url: "/path/to/resource"},
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
		{
			name: "Invalid domain without dot and not localhost",
			args: args{url: "example"},
			want: false,
		},
		{
			name: "Missing hostname",
			args: args{url: "https://"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, _ := IsValidURL(tt.args.url)
			if got != tt.want {
				t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
