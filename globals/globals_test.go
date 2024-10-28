package globals

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestMergeMaps(t *testing.T) {
	type args struct {
		a map[string]bool
		b map[string]bool
	}
	tests := []struct {
		name string
		args args
		want map[string]bool
	}{
		{
			name: "Both maps empty",
			args: args{
				a: map[string]bool{},
				b: map[string]bool{},
			},
			want: map[string]bool{},
		},
		{
			name: "First map empty",
			args: args{
				a: map[string]bool{},
				b: map[string]bool{"key1": true, "key2": false},
			},
			want: map[string]bool{"key1": true, "key2": false},
		},
		{
			name: "Second map empty",
			args: args{
				a: map[string]bool{"key1": true, "key2": false},
				b: map[string]bool{},
			},
			want: map[string]bool{"key1": true, "key2": false},
		},
		{
			name: "No overlapping keys",
			args: args{
				a: map[string]bool{"key1": true, "key2": false},
				b: map[string]bool{"key3": true, "key4": false},
			},
			want: map[string]bool{"key1": true, "key2": false, "key3": true, "key4": false},
		},
		{
			name: "Overlapping keys with different values",
			args: args{
				a: map[string]bool{"key1": true, "key2": false},
				b: map[string]bool{"key2": true, "key3": false},
			},
			want: map[string]bool{"key1": true, "key2": true, "key3": false},
		},
		{
			name: "Overlapping keys with same values",
			args: args{
				a: map[string]bool{"key1": true, "key2": false},
				b: map[string]bool{"key2": false, "key3": true},
			},
			want: map[string]bool{"key1": true, "key2": false, "key3": true},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := MergeMaps(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MergeMaps() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRenderToString(t *testing.T) {
	fromString := func(htmlStr string) *html.Node {
		doc, err := html.Parse(strings.NewReader(htmlStr))
		if err != nil {
			t.Fatal(err)
		}
		return doc
	}

	type args struct {
		node *html.Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Nile node",
			args: args{node: nil},
			want: "",
		},
		{
			name: "Empty node",
			args: args{node: fromString("")},
			want: "<html><head></head><body></body></html>",
		},
		{
			name: "Simple text node",
			args: args{node: fromString("Hello, world!")},
			want: "<html><head></head><body>Hello, world!</body></html>",
		},
		{
			name: "Simple element node",
			args: args{node: fromString("<div>Content</div>")},
			want: "<html><head></head><body><div>Content</div></body></html>",
		},
		{
			name: "Node with attributes",
			args: args{node: fromString(`<a href="https://example.com">Link</a>`)},
			want: `<html><head></head><body><a href="https://example.com">Link</a></body></html>`,
		},
		{
			name: "Nested elements",
			args: args{node: fromString("<div><p>Paragraph</p></div>")},
			want: "<html><head></head><body><div><p>Paragraph</p></div></body></html>",
		},
		{
			name: "Node with escaped characters",
			args: args{node: fromString("This < should be escaped")},
			want: "<html><head></head><body>This &lt; should be escaped</body></html>",
		},
		{
			name: "Empty element",
			args: args{node: fromString("<br/>")},
			want: "<html><head></head><body><br/></body></html>",
		},

		{
			name: "<a> link element",
			args: args{
				node: &html.Node{
					Type: html.ElementNode,
					Data: "a",
					Attr: []html.Attribute{
						{Key: "href", Val: "https://example.com"},
					},
					FirstChild: &html.Node{
						Type: html.TextNode,
						Data: "Link",
					},
				},
			},
			want: `<a href="https://example.com">Link</a>`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := RenderToString(tt.args.node); got != tt.want {
					t.Errorf("RenderToString() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func ExampleRenderToString_simpleElement() {
	node := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: "Content",
		},
	}
	rendered := RenderToString(node)
	fmt.Println(rendered)
	// Output: <div>Content</div>
}

func ExampleRenderToString_withAttributes() {
	node := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{Key: "href", Val: "https://example.com"},
		},
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: "Link",
		},
	}
	rendered := RenderToString(node)
	fmt.Println(rendered)
	// Output: <a href="https://example.com">Link</a>
}

func TestRoundBytes(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{56370, "0.06MB"},
	}

	for _, tt := range tests {
		got := RoundBytes(tt.input)
		if got != tt.want {
			t.Errorf("RoundBytes() Failed got %s want %s", got, tt.want)
		}
	}
}
