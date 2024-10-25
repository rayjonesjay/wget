package css

import (
	"fmt"
	"net/url"
	"testing"
)

func ExampleTransformCssUrl() {
	cssStr := `
.element {
  background-image: url("img/favicon.ico");
}

/* Using url() for CSS @import */
@import url("styles.css");

/* Using url() in a CSS variable */
:root {
  --background-image: url("path/to/image.jpg");
}
body {
  background-image: var(--background-image);
}`
	// Append https://example.com to the defined urls
	outCss := TransformCssUrl(
		cssStr, func(url string) string {
			return "https://example.com/" + url
		},
	)
	fmt.Printf("Output: \n%s\n", outCss)
}

func TestTransformCssUrl(t *testing.T) {
	_css1 := `
/* Using url() with data URI for embedded images */
.element {
  background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=");
}

/* Using url() for CSS @import */
@import url("styles.css");

/* Using url() in a CSS variable */
:root {
  --background-image: url("path/to/image.jpg");
}
body {
  background-image: var(--background-image);
}`
	_css1Append2 := `
/* Using url() with data URI for embedded images */
.element {
  background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=(2)");
}

/* Using url() for CSS @import */
@import url("styles.css(2)");

/* Using url() in a CSS variable */
:root {
  --background-image: url("path/to/image.jpg(2)");
}
body {
  background-image: var(--background-image);
}`
	_css1SmartAppend2 := `
/* Using url() with data URI for embedded images */
.element {
  background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=");
}

/* Using url() for CSS @import */
@import url("styles.css(2)");

/* Using url() in a CSS variable */
:root {
  --background-image: url("path/to/image.jpg(2)");
}
body {
  background-image: var(--background-image);
}`
	noOp := func(_url string) string {
		return _url
	}

	append2 := func(_url string) string {
		return _url + "(2)"
	}

	smartAppend2 := func(_url string) string {
		r, err := url.Parse(_url)
		if err != nil {
			panic(err)
		}

		if r.Opaque != "" {
			return _url
		}
		return _url + "(2)"
	}

	type args struct {
		_css        string
		transformer func(_url string) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty",
			args: args{
				_css:        "",
				transformer: nil,
			},
			want: "",
		},

		{
			name: "Nil transformer",
			args: args{
				_css:        `@import url("styles.css");`,
				transformer: nil,
			},
			want: `@import url("styles.css");`,
		},

		{
			name: "Simple @import",
			args: args{
				_css:        `@import url("styles.css");`,
				transformer: nil,
			},
			want: `@import url("styles.css");`,
		},

		{
			name: "Simple @import noOp transformer",
			args: args{
				_css:        `@import url("styles.css");`,
				transformer: noOp,
			},
			want: `@import url("styles.css");`,
		},

		{
			name: "Simple @import append2 transformer",
			args: args{
				_css:        `@import url("styles.css");`,
				transformer: append2,
			},
			want: `@import url("styles.css(2)");`,
		},

		{
			name: "Css example 1",
			args: args{
				_css:        _css1,
				transformer: nil,
			},
			want: _css1,
		},

		{
			name: "Css example 1 with append2 transformer",
			args: args{
				_css:        _css1,
				transformer: append2,
			},
			want: _css1Append2,
		},

		{
			name: "Css example 1 with smartAppend2 transformer",
			args: args{
				_css:        _css1,
				transformer: smartAppend2,
			},
			want: _css1SmartAppend2,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := TransformCssUrl(tt.args._css, tt.args.transformer); got != tt.want {
					t.Errorf("TransformCssUrl() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
