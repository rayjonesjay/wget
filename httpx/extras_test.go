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

func TestFilenameFromContentDisposition(t *testing.T) {
	type args struct {
		headers http.Header
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Valid filename with quotes",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; filename=\"my_document.pdf\""}}},
			want:    "my_document.pdf",
			wantErr: false,
		},
		{
			name:    "Valid filename without quotes",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; filename=my_document.pdf"}}},
			want:    "my_document.pdf",
			wantErr: false,
		},
		{
			name:    "Filename with spaces",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; filename=\" document with spaces.pdf \""}}},
			want:    " document with spaces.pdf ",
			wantErr: false,
		},
		{
			name:    "Directive filename surrounded with spaces",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; filename = \" document with spaces.pdf \" "}}},
			want:    " document with spaces.pdf ",
			wantErr: false,
		},
		{
			name:    "Case-insensitive filename parameter",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; FILENAME=\"my_document.pdf\""}}},
			want:    "my_document.pdf",
			wantErr: false,
		},
		{
			name:    "No filename parameter",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment"}}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Empty Content-Disposition header",
			args:    args{headers: http.Header{}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid filename format (no equals sign)",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; filename"}}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Filename with special characters",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; filename=\"file-with_special-chars!.pdf\""}}},
			want:    "file-with_special-chars!.pdf",
			wantErr: false,
		},
		{
			name:    "Filename with Unicode characters",
			args:    args{headers: http.Header{"Content-Disposition": []string{"attachment; filename=\"文件.pdf\""}}},
			want:    "文件.pdf",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := FilenameFromContentDisposition(tt.args.headers)
				if (err != nil) != tt.wantErr {
					t.Errorf("FilenameFromContentDisposition() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("FilenameFromContentDisposition() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRoundOfSizeOfData(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  string
	}{
		{"A-TEST", 200, "200.00KB"},
		{"B-TEST", 1_000_000, "1.00MB"},
		{"C-TEST", 1_000_000_000, "1.00GB"},
		{"D-TEST", 1_000_000_000_000, "1000.00GB"},
		{"E-TEST", 1, "1.00KB"},
		{"F-TEST", 234, "234.00KB"},
		{"G-TEST", 213_432, "0.21MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundOfSizeOfData(tt.input)
			if got != tt.want {
				t.Errorf("RoundOfSizeOfData(%d) == %s want %s", tt.input, got, tt.want)
			}
		})
	}
}
