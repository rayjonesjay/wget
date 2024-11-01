package globals

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"wget/temp"

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
		{2000000000, "2.00GB"},
	}

	for _, tt := range tests {
		got := RoundBytes(tt.input)
		if got != tt.want {
			t.Errorf("RoundBytes() Failed got %s want %s", got, tt.want)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{1023, "1023 B"},
		{0, "0 B"},
		{1024, "1.00 KiB"},
		{2048, "2.00 KiB"},
		{1048575, "1024.00 KiB"},
		{1048576, "1.00 MiB"},
		{2097152, "2.00 MiB"},
		{1073741823, "1024.00 MiB"},
		{1073741824, "1.00 GiB"},
		{2147483648, "2.00 GiB"},
		{-1, "--.- B"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := FormatSize(test.size)
			if result != test.expected {
				t.Errorf("FormatSize(%d) = %s; expected %s", test.size, result, test.expected)
			}
		})
	}
}

// ExampleFormatSize demonstrates the usage of FormatSize.
func ExampleFormatSize() {
	fmt.Println(FormatSize(1500))
	fmt.Println(FormatSize(1000000))
	fmt.Println(FormatSize(1234567890))
	fmt.Println(FormatSize(-500))

	// Output:
	// 1.46 KiB
	// 976.56 KiB
	// 1.15 GiB
	// --.- B
}

func TestStringTimes(t *testing.T) {
	tests := []struct {
		inputString string
		inputCount  int
		expected    []string
	}{
		{"hello", 3, []string{"hello", "hello", "hello"}},
		{"world", 0, nil},
		{"", 5, []string{"", "", "", "", ""}},
		{"test", 1, []string{"test"}},
		{"repeat", -2, nil},
	}

	for _, test := range tests {
		t.Run(test.inputString, func(t *testing.T) {
			result := StringTimes(test.inputString, test.inputCount)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("StringTimes(%q, %d) = %v; expected %v",
					test.inputString, test.inputCount, result, test.expected)
			}
		})
	}
}

// ExampleStringTimes demonstrates the usage of StringTimes.
func ExampleStringTimes() {
	result := StringTimes("go", 3)
	fmt.Println(result)

	result = StringTimes("hello", 0)
	fmt.Println(result)
	// Output:
	// [go go go]
	// []
}

func TestPrintLines(t *testing.T) {
	tests := []struct {
		baseRow  int
		lines    []string
		expected string
	}{
		{1, []string{"Line 1", "Line 2"}, "\x1b[2;0H\x1b[KLine 1\x1b[3;0H\x1b[KLine 2"},
		{1, []string{"First", "Second"}, "\x1b[2;0H\x1b[KFirst\x1b[3;0H\x1b[KSecond"},
	}

	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	for _, test := range tests {
		var err error
		os.Stdout, err = temp.File()
		if err != nil {
			t.Fatal(err)
		}

		PrintLines(test.baseRow, test.lines)
		output, err := os.ReadFile(os.Stdout.Name())
		if err != nil {
			t.Fatal(err)
		}
		if string(output) != test.expected {
			t.Errorf("PrintLines(%d, %v) = %q; expected %q",
				test.baseRow, test.lines, output, test.expected)
		}
	}
}
