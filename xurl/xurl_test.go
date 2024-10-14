package xurl

import (
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name    string
		rawUrl  string
		want    bool
		wantErr bool
	}{
		//valid cases
		{"Valid HTTP URL",
		 "http://google.com",
		  true, 
		  false},

		{"Valid HTTPS URL", 
		"https://google.com", 
		true, 
		false},

		{"URL without scheme (http)",
		 "google.com",
		  true,
		   false},

		{"URL without scheme (https)",
		 "google.com", true, 
		 false},

		 //invalid cases
		{"Invalid URL with missing host", 
		"http://",
		 false,
		 true},

		{"Invalid URL with wrong scheme",
		 "ftp://google.com",
		  false, 
		  true},

		{"Malformed URL with dot prefix",
		 ".google.com", false,
		  true},

		{"Malformed URL with empty host", 
		"http://-google.com", 
		false,
		 true},

		{"Invalid URL with localhost",
		 "http://localhost",
		  true,
		   false},

		{"Valid localhost without scheme", 
		"localhost", 
		true, 
		false},

		{"Valid localhost with scheme",
		 "http://localhost",
		  true,
		   false},

		{"Invalid URL with missing domain", 
		"http://-localhost", 
		false, 
		true},

		{"Valid URL with subdomain", 
		"http://sub.google.com", 
		true, 
		false},

		{"Invalid URL with no scheme and invalid host", 
		"-google.com",
		 false,
		 true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsValidURL(tt.rawUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
