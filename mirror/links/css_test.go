package links

import (
	"reflect"
	"testing"
)

func TestFromCssUrl(t *testing.T) {
	testCases := []struct {
		name     string
		css      string
		expected []string
	}{
		{
			name: "Valid URLs",
			css: `
                .class1 {
                    background-image: url('image1.jpg');
                }
                .class2 {
                    background-image: url("image2.png");
                }
            `,
			expected: []string{"image1.jpg", "image2.png"},
		},
		{
			name: "URLs with special characters",
			css: `
                .class1 {
                    background-image: url('image@1.jpg');
                }
                .class2 {
                    background-image: url("image%2.png");
                }
            `,
			expected: []string{"image@1.jpg", "image%2.png"},
		},
		{
			name: "URLs with spaces",
			css: `
                .class1 {
                    background-image: url('  image1.jpg  ');
                }
                .class2 {
                    background-image: url("  image2.png  ");
                }
            `,
			expected: []string{"  image1.jpg  ", "  image2.png  "},
		},
		{
			name: "URLs without quotes",
			css: `
                .class1 {
                    background-image: url(image1.jpg);
                }
                .class2 {
                    background-image: url(image2.png);
                }
            `,
			expected: []string{"image1.jpg", "image2.png"},
		},
		{
			name: "Empty URLs",
			css: `
                .class1 {
                    background-image: url('');
                }
                .class2 {
                    background-image: url("");
                }
            `,
			expected: []string{"", ""},
		},
		{
			name: "No URLs",
			css: `
                .class1 {
                    color: red;
                }
                .class2 {
                    font-size: 16px;
                }
            `,
			expected: nil,
		},
		{
			name: "URLs within comments",
			css: `
                /* Comment with url('ignored1.jpg') */
                .class1 {
                    background-image: url('image1.jpg');
                }
                /* Comment with url("ignored2.png") */
                .class2 {
                    background-image: url("image2.png");
                }
            `,
			expected: []string{"image1.jpg", "image2.png"},
		},
		{
			name: "Malformed URLs",
			css: `
                .class1 {
                    background-image: url('image1.jpg);
                }
                .class2 {
                    background-image: url(image2.png');
                }
            `,
			expected: nil,
		},

		{
			name: "Edge case I",
			css: `/* Example with comments */
            body { background-image: url("http://example.com/bg.jpg"); }
            /* Ignore this: url('http://example.com/commented.jpg') */
            div { background: url('http://example.com/div.jpg') no-repeat; }
            #id { background: url(http://example.com/bg4.jpg); }
            /* Another comment */
            span { font-family: "url('fake-url')"; }`,
			expected: []string{
				"http://example.com/bg.jpg",
				"http://example.com/div.jpg",
				"http://example.com/bg4.jpg",
			},
		},

		{
			name: "Edge case III",
			css: `
				/* This is a comment with url(http://example.com/commented-out.png) */
				/* A comment with a tricky 'url("http://example.com/inside-comment.png")' */
				
				.some-class::after {
				  content: "Check out this image: url('images/inside-string.png')";
				}

				.some-other-class::after {
				  content: 'Check out this image: url('images/inside-string.png')';
				}
`,
			expected: nil,
		},

		{
			name: "Poor man's parser",
			css: `
.bad-format {
   background: url("url('http://example.com/bad-format.jpg')");
}
`,
			expected: []string{
				"http://example.com/bad-format.jpg",
			},
		},

		{
			name: "Edge case II",
			css: `
body {
  background: url('http://example.com/bg.png');
}

/* Another comment */ 

div {
  background-image: url("images/bg2.png");
}

.header {
  background-image: url(images/bg3.png);
}

.footer {
  background: url('images/bg4 with spaces.png');
}

/*
Multi-line comment with url(images/multiline-comment.png)
*/

@font-face {
  font-family: 'MyFont';
  src: url('fonts/myfont.woff2') format('woff2'),
       url("fonts/myfont.woff") format('woff');
}
`,
			expected: []string{
				"http://example.com/bg.png",
				"images/bg2.png",
				"images/bg3.png",
				"images/bg4 with spaces.png",
				"fonts/myfont.woff2",
				"fonts/myfont.woff",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				result := FromCssUrl(tc.css)
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("Expected: \n%v, Got: \n%v", tc.expected, result)
				}
			},
		)
	}
}
