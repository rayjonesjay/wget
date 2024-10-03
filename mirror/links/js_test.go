package links

import (
	"reflect"
	"testing"
)

func TestFromRequire(t *testing.T) {
	type args struct {
		js string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty Js Code",
			args: args{js: ``},
			want: nil,
		},

		{
			name: "Simple console log",
			args: args{
				js: `
console.log('Hello, World!');
`,
			},
			want: nil,
		},

		{
			name: "require single quotes",
			args: args{
				js: `
const foo = require('./foo.js')
`,
			},
			want: []string{
				`./foo.js`,
			},
		},

		{
			name: "require backtick quotes",
			args: args{js: "const foo = require(`./foo.js`)"},
			want: []string{
				`./foo.js`,
			},
		},

		{
			name: "Invalid require",
			args: args{
				js: `
require(); // Invalid require
`,
			},
			want: nil,
		},

		{
			name: "Empty require",
			args: args{
				js: `
require(''); // empty require
`,
			},
			want: nil,
		},

		{
			name: "require module name with dashes",
			args: args{js: `require('./file-with-dashes.js');`},
			want: []string{
				"./file-with-dashes.js",
			},
		},

		{
			name: "require module name with spaces",
			args: args{js: `require('./file-with-spaces.js');`},
			want: []string{
				"./file-with-spaces.js",
			},
		},

		{
			name: "require with invalid syntax",
			args: args{
				js: `
const myVar = require('http://web-module.location');
const anotherVar = require http://another-web-module.location; // Invalid syntax
`,
			},
			want: []string{
				"http://web-module.location",
			},
		},

		{
			name: "require another module",
			args: args{
				js: `
const myVar = require('http://web-module.location');
const anotherVar = require('https://another-web-module.location');
`,
			},
			want: []string{
				"http://web-module.location",
				"https://another-web-module.location",
			},
		},

		{
			name: "require base",
			args: args{
				js: `
const getFullName = require('./utils.js');
console.log(getFullName('John', 'Doe')); // My fullname is John Doe

const myVar = require('http://web-module.location');
`,
			},
			want: []string{
				"./utils.js",
				"http://web-module.location",
			},
		},

		{
			name: "Poor man's parser require commented",
			args: args{
				js: `
/* require('./commented2.js') */
`,
			},
			want: []string{
				"./commented2.js",
			},
		},

		{
			name: "Poor man's parser require commented (line comment)",
			args: args{
				js: `
// require('./commented2.js')
`,
			},
			want: []string{
				"./commented2.js",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := FromJsRequire(tt.args.js); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FromJsRequire() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestFromImport(t *testing.T) {
	type args struct {
		js string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty Js Code",
			args: args{js: ``},
			want: nil,
		},

		{
			name: "Simple console log",
			args: args{
				js: `
console.log('Hello, World!');
`,
			},
			want: nil,
		},

		{
			name: "Comment in file",
			args: args{
				js: `
// This is a comment
import defaultExport from "/js/main.js"; // Another comment
import * as name from "./js/main.js";
`,
			},
			want: []string{
				"/js/main.js",
				"./js/main.js",
			},
		},

		{
			name: "Clean Import",
			args: args{
				js: `
import foo from './foo.js'
`,
			},
			want: []string{
				"./foo.js",
			},
		},

		{
			name: "Clean Imports with backticks",
			args: args{js: "import foo from `./foo.js`"},
			want: []string{
				"./foo.js",
			},
		},

		{
			name: "Invalid import",
			args: args{
				js: `
import; // Invalid import
`,
			},
			want: nil,
		},

		{
			name: "Empty Import",
			args: args{
				js: `
import '';
`,
			},
			want: nil,
		},

		{
			name: "Import file with dashes",
			args: args{
				js: `
import './js/file-with-dashes.js'
`,
			},
			want: []string{
				"./js/file-with-dashes.js",
			},
		},

		{
			name: "Import file with spaces",
			args: args{
				js: `
import './js/file with spaces.js'
`,
			},
			want: []string{
				"./js/file with spaces.js",
			},
		},

		{
			name: "Import file with spaces",
			args: args{
				js: `
import "./js/file with spaces.js"
`,
			},
			want: []string{
				"./js/file with spaces.js",
			},
		},

		{
			name: " Invalid Import syntax",
			args: args{
				js: `
import { export1, export2 } from "/js/main1.js";
import { export1, export2 } from /js/main2.js;
`,
			},
			want: []string{
				"/js/main1.js",
			},
		},

		//		{
		//			name: "Edge Case I",
		//			args: args{
		//				js: `
		//import defaultExport from "/js/main.js";
		//import * as name from "./js/main.js";
		//import { export1 } from "/external/js/main.js";
		//import { export1 as alias1 } from "../../js/main.js";
		//import { default as alias } from "https://web.jsdeliver.com/stable/main.js";
		//import { export1, export2 } from "/js/main.js";
		//import { export1, export2 as alias2, /* … */ } from "/js/main.js";
		//import { "string name" as alias } from "/js/main.js";
		//import defaultExport, { export1, /* … */ } from "/js/main.js";
		//import defaultExport, * as name from "/js/main.js";
		//import "/js/main.js";`,
		//			},
		//			want: []string{
		//				"/js/main.js",
		//				"./js/main.js",
		//				"/external/js/main.js",
		//				"../../js/main.js",
		//				"https://web.jsdeliver.com/stable/main.js",
		//				"/js/main.js",
		//				"/js/main.js",
		//				"/js/main.js",
		//				"/js/main.js",
		//				"/js/main.js",
		//				"/js/main.js",
		//			},
		//		},

		//		{
		//			name: "Commented out import",
		//			args: args{
		//				js: `
		///* import './commented2.js' */
		//`,
		//			},
		//			want: nil,
		//		},
		//
		//		{
		//			name: "Commented import (line comment)",
		//			args: args{
		//				js: `
		////import './commented.js';
		//`,
		//			},
		//			want: nil,
		//		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := FromJsImport(tt.args.js); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FromJsImport() = \n%v, want \n%v", got, tt.want)
				}
			},
		)
	}
}

func TestToJsModules(t *testing.T) {
	type args struct {
		js string
	}
	tests := []struct {
		name      string
		args      args
		wantLinks []string
	}{
		{
			name: "",
			args: args{
				js: `
import defaultExport from "/js/main.js";
import * as name from "./js/main.js";
import { export1 } from "/external/js/main.js";
import { export1 as alias1 } from "../../js/main.js";
import { default as alias } from "https://web.jsdeliver.com/stable/main.js";

const getFullName = require('./utils.js');
console.log(getFullName('John', 'Doe')); // My fullname is John Doe

const myVar = require('http://web-module.location');
`,
			},
			wantLinks: []string{
				"/js/main.js",
				"./js/main.js",
				"/external/js/main.js",
				"../../js/main.js",
				"https://web.jsdeliver.com/stable/main.js",
				"./utils.js",
				"http://web-module.location",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotLinks := FromJs(tt.args.js); !reflect.DeepEqual(gotLinks, tt.wantLinks) {
					t.Errorf("FromJs() = \n%v, want \n%v", gotLinks, tt.wantLinks)
				}
			},
		)
	}
}
