package download

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockserver
func mockServer(statusCode int, body string) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if body != "" {
			w.Write([]byte(body))
		}
	})
	return httptest.NewServer(handler)
}

func TestDownloadUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid URL",
			args:    args{url: mockServer(http.StatusOK, "test content").URL},
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			args:    args{url: "invalid-url"},
			wantErr: true,
		},
		{
			name:    "Non-Existent URL",
			args:    args{url: mockServer(http.StatusNotFound, "").URL},
			wantErr: true,
		},
		{
			name:    "Server Error URL",
			args:    args{url: mockServer(http.StatusInternalServerError, "").URL},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadUrl(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("DownloadUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

