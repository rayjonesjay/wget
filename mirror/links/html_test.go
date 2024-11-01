package links

import (
	"fmt"
	"golang.org/x/net/html"
	"reflect"
	"strings"
	"testing"
)

func TestFromHtml(t *testing.T) {
	fromString := func(s string) *html.Node {
		doc, err := html.Parse(strings.NewReader(s))
		if err != nil {
			t.Fatal(err)
		}
		return doc
	}

	type args struct {
		doc *html.Node
	}
	tests := []struct {
		name     string
		args     args
		wantUrls []string
	}{
		{
			name:     "Empty document",
			args:     args{doc: fromString("<html></html>")},
			wantUrls: nil,
		},
		{
			name:     "No links or URLs",
			args:     args{doc: fromString("<html><body><p>Hello World!</p></body></html>")},
			wantUrls: nil,
		},
		{
			name:     "Single anchor tag",
			args:     args{doc: fromString("<html><body><a href='https://example.com'>Link</a></body></html>")},
			wantUrls: []string{"https://example.com"},
		},
		{
			name:     "Multiple anchor tags",
			args:     args{doc: fromString("<html><body><a href='https://example.com'>Link</a><a href='https://example.org'>Another Link</a></body></html>")},
			wantUrls: []string{"https://example.com", "https://example.org"},
		},
		{
			name:     "Image tag with URL",
			args:     args{doc: fromString("<html><body><img src='https://example.com/image.png' /></body></html>")},
			wantUrls: []string{"https://example.com/image.png"},
		},
		{
			name:     "Link tag in head",
			args:     args{doc: fromString("<html><head><link rel='stylesheet' href='https://example.com/style.css' /></head><body></body></html>")},
			wantUrls: []string{"https://example.com/style.css"},
		},
		{
			name:     "Style tag with URL",
			args:     args{doc: fromString("<html><head><style>@import url('https://example.com/styles.css');</style></head><body></body></html>")},
			wantUrls: []string{"https://example.com/styles.css"},
		},
		{
			name:     "Object tag with URL",
			args:     args{doc: fromString("<html><body><object data='https://example.com/object.swf'></object></body></html>")},
			wantUrls: nil,
		},
		{
			name:     "Mixed tags",
			args:     args{doc: fromString("<html><body><a href='https://example.com'>Link</a><img src='https://example.com/image.png' /><link rel='stylesheet' href='https://example.com/style.css' /></body></html>")},
			wantUrls: []string{"https://example.com", "https://example.com/image.png", "https://example.com/style.css"},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotUrls := FromHtml(tt.args.doc); !reflect.DeepEqual(gotUrls, tt.wantUrls) {
					t.Errorf("FromHtml() = %v, want %v", gotUrls, tt.wantUrls)
				}
			},
		)
	}
}

func ExampleFromHtml() {
	doc, err := html.Parse(strings.NewReader("<html><body><a href='https://example.com'>Link</a></body></html>"))
	if err != nil {
		panic(err)
	}
	urls := FromHtml(doc)
	for _, url := range urls {
		fmt.Println(url)
	}
	// Output: https://example.com
}
