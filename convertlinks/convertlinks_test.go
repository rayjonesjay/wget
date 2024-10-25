package convertlinks

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
	"wget/globals"
)

func TestOfHtml(t *testing.T) {
	fromString := func(htmlStr string) *html.Node {
		doc, err := html.Parse(strings.NewReader(htmlStr))
		if err != nil {
			t.Fatal(err)
		}
		return doc
	}

	baseHtml := "<html><head><style>background-image: url(https://example.com/path/image.png);</style></head><body></body></html>"
	baseNode := fromString(baseHtml)
	baseConverted := "<html><head><style>background-image: url(https://example.com/path/image.png.2);</style></head><body></body></html>"
	// test transformer, simply appends ".2" to the end of the target url
	transformer := func(url string, _ bool) string {
		return url + ".2"
	}

	// removes all newline chars from the string s,
	//such that the returned string has all it's characters in a single line
	sl := func(s string) string {
		return strings.ReplaceAll(s, "\n", "")
	}

	type args struct {
		n           *html.Node
		transformer func(url string, isA bool) string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Nil args",
			args: args{
				n:           nil,
				transformer: nil,
			},
			want: "",
		},

		{
			name: "Nil html node",
			args: args{
				n:           nil,
				transformer: transformer,
			},
			want: "",
		},

		{
			name: "Nil transformer",
			args: args{
				n:           baseNode,
				transformer: nil,
			},
			want: baseHtml,
		},

		{
			name: "Base HTML",
			args: args{
				n:           baseNode,
				transformer: transformer,
			},
			want: baseConverted,
		},

		{
			name: "Inline CSS in HTML",
			args: args{
				n: fromString(
					`<html><head></head><body><div id="one" style="background-image: url(https://example.com/path/div1.png);"></div></body></html>`,
				),
				transformer: transformer,
			},
			want: `<html><head></head><body><div id="one" style="background-image: url(https://example.com/path/div1.png.2);"></div></body></html>`,
		},

		{
			name: "Inline CSS in HTML",
			args: args{
				n: fromString(
					`<html><head></head><body><div id="one" style="background-image: url(https://example.com/path/div1.png);">Hello url(https://example.com/path/div1.png)</div></body></html>`,
				),
				transformer: transformer,
			},
			want: `<html><head></head><body><div id="one" style="background-image: url(https://example.com/path/div1.png.2);">Hello url(https://example.com/path/div1.png)</div></body></html>`,
		},

		{
			name: "Inline CSS in HTML (url with single quotes)",
			args: args{
				n: fromString(
					`<html>
<head></head>
<body><div id="one" style="background-image: url('https://example.com/path/div1.png');"></div>
</body>
</html>`,
				),
				transformer: transformer,
			},
			want: `<html>
<head></head>
<body><div id="one" style="background-image: url(&#39;https://example.com/path/div1.png.2&#39;);"></div>
</body>
</html>`,
		},

		{
			name: "Inline CSS in HTML (url with backtick quotes)",
			args: args{
				n: fromString(
					"<html><head></head><body><div id=\"one\" style=\"background-image: url(`https://example.com/path/div1.png`);\"></div></body></html>",
				),
				transformer: transformer,
			},
			want: "<html><head></head><body><div id=\"one\" style=\"background-image: url(`https://example.com/path/div1.png.2`);\"></div></body></html>",
		},

		{
			name: "Multiple Background Images HTML",
			args: args{
				n: fromString(
					`<html><head>
<style>background-image: url(https://example.com/path/image.png), 
url(https://example.com/path/image1.png), 
url(https://example.com/path/image2.png);
</style></head><body></body></html>`,
				),
				transformer: transformer,
			},
			want: `<html><head>
<style>background-image: url(https://example.com/path/image.png.2), 
url(https://example.com/path/image1.png.2), 
url(https://example.com/path/image2.png.2);
</style></head><body></body></html>`,
		},

		{
			name: "Multiple Background Images HTML with inline CSS",
			args: args{
				n: fromString(
					`<html><head>
<style>background-image: url(https://example.com/path/image.png), 
url(https://example.com/path/image1.png), 
url(https://example.com/path/image2.png);
</style></head><body style="background-image: url(https://example.com/path/body.png)">
<div class="one" style="background-image: url(https://example.com/path/div-with-space.png)"></div>
<div class="two" style="background-image: url(https://example.com/path/div2.png)"></div>
</body></html>`,
				),
				transformer: transformer,
			},
			want: `<html><head>
<style>background-image: url(https://example.com/path/image.png.2), 
url(https://example.com/path/image1.png.2), 
url(https://example.com/path/image2.png.2);
</style></head><body style="background-image: url(https://example.com/path/body.png.2)">
<div class="one" style="background-image: url(https://example.com/path/div-with-space.png.2)"></div>
<div class="two" style="background-image: url(https://example.com/path/div2.png.2)"></div>
</body></html>`,
		},

		{
			name: "HTML with multiple resource linking types",
			args: args{
				n: fromString(
					sl(
						`<html lang="en">
<head>
    <link rel="stylesheet" href="https://www.example.com/styles.css"/>
    <script src="https://www.example.com/script.js"></script>
    <title></title>
</head>
<body>
<img src="https://www.example.com/image.jpg" alt="Example image"/>
<div id="one" style="background-image: url(https://example.com/path/div1.png);"></div>
<object data="https://www.example.com/object.swf" type="application/x-shockwave-flash"></object>
</body>
</html>`,
					),
				),
				transformer: transformer,
			},
			want: sl(
				`<html lang="en">
<head>
    <link rel="stylesheet" href="https://www.example.com/styles.css.2"/>
    <script src="https://www.example.com/script.js"></script>
    <title></title>
</head>
<body>
<img src="https://www.example.com/image.jpg.2" alt="Example image"/>
<div id="one" style="background-image: url(https://example.com/path/div1.png.2);"></div>
<object data="https://www.example.com/object.swf" type="application/x-shockwave-flash"></object>
</body>
</html>`,
			),
		},

		{
			name: "Tester Template",
			args: args{
				n:           fromString(``),
				transformer: transformer,
			},
			want: `<html><head></head><body></body></html>`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				OfHtml(tt.args.n, tt.args.transformer)
				expected := globals.RenderToString(tt.args.n)
				expected, tt.want = sl(expected), sl(tt.want)
				if expected != tt.want {
					t.Errorf("OfHtml() got >>> \n%v >>> want >>> \n%v", expected, tt.want)
				}
			},
		)
	}
}
