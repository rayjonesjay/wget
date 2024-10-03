package mirror

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"strings"
	"testing"
)

func example() {
	s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func Test_example(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "go.dev example"},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				example()
			},
		)
	}
}

func TestExtract(t *testing.T) {
	nodes := []struct {
		Html string
		Urls []string
		Name string
	}{
		{
			Name: "default ",
			Html: `<html><body><img src="https://www.example.com/image.jpg" alt="Example image"></body></html>`,
			Urls: []string{"https://www.example.com/image.jpg"},
		},
	}

	createNode := func(_html string) *html.Node {
		_node, err := html.Parse(strings.NewReader(_html))
		if err != nil {
			log.Fatal(err)
		}
		return _node
	}

	type args struct {
		node *html.Node
	}
	var tests []struct {
		name    string
		args    args
		wantOut []UrlExtract
	}

	for _, n := range nodes {
		test := struct {
			name    string
			args    args
			wantOut []UrlExtract
		}{
			name:    n.Name,
			args:    args{node: createNode(n.Html)},
			wantOut: nil,
		}

		for _, url := range n.Urls {
			test.wantOut = append(test.wantOut, UrlExtract{Url: url})
		}

		tests = append(tests, test)
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotOut := Extract(tt.args.node)
				if len(gotOut) != len(tt.wantOut) {
					t.Errorf("Extract() = %v, want %v", gotOut, tt.wantOut)
				}
				for i := range tt.wantOut {
					if gotOut[i].Url != tt.wantOut[i].Url {
						t.Errorf("Extract() = found Url %v, want %v", gotOut, tt.wantOut)
					}
				}
			},
		)
	}
}

func TestRenderToString(t *testing.T) {
	_html := `<html><body><img src="https://www.example.com/image.jpg" alt="Example image"></body></html>`
	expectedLink := "https://www.example.com/image.jpg(2)"
	node := createNode(_html)

	links := Extract(node)
	for _, link := range links {
		link.Attr.Val = link.Attr.Val + "(2)"
	}

	outputHtml, err := RenderToString(node)
	if err != nil {
		log.Fatal(err)
	}

	linkFromOutputHtml, err := ExtractFirst(createNode(outputHtml))
	if err != nil {
		log.Fatal(err)
	}

	if linkFromOutputHtml.Url != expectedLink {
		t.Errorf("RenderToString() failed: got link %q, want %q", linkFromOutputHtml.Url, expectedLink)
	}
}

func createNode(_html string) *html.Node {
	_node, err := html.Parse(strings.NewReader(_html))
	if err != nil {
		log.Fatal(err)
	}
	return _node
}
