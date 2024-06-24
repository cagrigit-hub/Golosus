package main

import "fmt"

type Content struct{}

func (c *Content) Main(name, github string) string {
	return fmt.Sprintf(`
package main

import (
	"github.com/%s/%s/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	exampleHandler := &handler.ExampleHandler{}
	app.Static("/static", "assets")
	app.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	app.GET("/example", exampleHandler.HandleExampleShow)
	app.POST("/example", exampleHandler.HandlePost)
	app.Start(":3000")
}

`, github, name)
}

func (c *Content) Layout(title string) string {
	layout := fmt.Sprintf(`
package layout

templ Base() {
	<html>
		<head>
			<title>Hello! %s</title>
			<script src="https://unpkg.com/htmx.org@1.9.10"></script>
			<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
			<script src="https://cdn.tailwindcss.com"></script>
		</head>
		<body>
			This is from the base layout
			{ children... }
			<script type="module" src="/static/bundled/bundle.js"></script>
		</body>
	</html>
}


`, title)
	return layout
}

func (c *Content) Util() string {
	return `
package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}
`
}

func (c *Content) ExampleModel() string {
	return `	
package model

type Example struct {
	Text string
}
	`
}

func (c *Content) ExampleHandler(github, name string) string {
	return fmt.Sprintf(`
package handler

import (
	"github.com/%s/%s/model"
	"github.com/%s/%s/view/example"
	"github.com/labstack/echo/v4"
)

type ExampleHandler struct{}

func (h *ExampleHandler) HandleExampleShow(c echo.Context) error {
	u := model.Example{
		Text: "example-text",
	}
	return render(c, example.Show(u))
}

func (h *ExampleHandler) HandlePost(c echo.Context) error {
	c.Request().ParseForm()
	if c.Request().Form.Has("example") {
		u := model.Example{
			Text: c.Request().Form.Get("example"),
		}
		return render(c, example.EcOne(u))
	}
	return c.String(400, "Bad Request")
}

	`, github, name, github, name)
}

func (c *Content) ExampleView(github, name string) string {
	return fmt.Sprintf(`
package example

import (
	"github.com/%s/%s/view/layout"
	"github.com/%s/%s/view/components"
	"github.com/%s/%s/model"
)

templ Show(example model.Example) {
	@layout.Base() {
		<div>
			@EcOne(example)
			<form hx-post="/example" hx-target="#example" hx-swap="outerHTML">
				@components.Input(components.InputProps{Type: "text", Name: "example"})
				<button>Submit</button>
			</form>
			<div class="text-red-400">
				Tailwind Configured
			</div>
			<div x-data="{ open: false }">
				<button @click="open = true">Expand</button>
				<span x-show="open">
					Content...
				</span>
			</div>
		</div>
	}
}

templ EcOne(example model.Example) {
	<h1 id="example">hello { example.Text } from the user </h1>
}


`, github, name, github, name, github, name)
}

func (c *Content) ExampleComponent() string {
	return `
package components

type InputProps struct {
	Type string
	Name string
}

templ Input(props InputProps) {
	<input type={ props.Type } name={ props.Name }/>
}
`
}

func (c *Content) Make() string {
	return `
gen:
	@templ generate
init:
	@templ generate
	@go mod tidy
	@cd ./typescript && npm install 
run: 
	@templ generate
	@cd ./typescript && npm run build 
	@go run ./cmd $(ARGS)
build:
	@templ generate
	@cd ./typescript && npm run build
	@go build -o ./tmp/bin ./cmd
`
}

func (c *Content) GoMod(github, name string) string {
	return fmt.Sprintf(`
module github.com/%s/%s

go 1.22.0

require (
	github.com/a-h/templ v0.2.543 // indirect
	github.com/labstack/echo/v4 v4.11.4 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

	`, github, name)
}

func (c *Content) TypescriptIndex() string {
	return `
import scripts from "./scripts";

// load scripts
scripts.forEach((src) => {
  let script = document.createElement("script");
  script.setAttribute("src", src);
  document.body.appendChild(script);
});

let x: number = 1;
console.log(x);

`
}

func (c *Content) ScriptsTs() string {
	return `
const scripts = [""];
export default scripts;

  `
}

func (c *Content) PackageJson(name string) string {
	return fmt.Sprintf(`
{
  "name": "%s",
  "version": "1.0.0",
  "description": "golosus-web-app, go check github.com/cagrigit-hub/golosus",
  "main": "index.ts",
  "scripts": {
    "build": "rimraf ./ts-build && npx tsc && browserify --node --ignore-missing ./ts-build/index.js | terser > ../assets/bundled/bundle.js"
  },
  "keywords": [
    "golosus"
  ],
  "author": "@cagrigit-hub",
  "license": "ISC",
  "type": "module",
  "devDependencies": {
    "@types/node": "^20.5.6",
    "rimraf": "^5.0.7",
    "typescript": "^5.2.2"
  },
  "dependencies": {
    "browserify": "^17.0.0",
    "terser": "^5.29.2"
  }
}


`, name)
}

func (c *Content) TsConfig() string {
	return `
{
  "compilerOptions": {
    /* Visit https://aka.ms/tsconfig to read more about this file */

    /* Projects */
    // "incremental": true,                              /* Save .tsbuildinfo files to allow for incremental compilation of projects. */
    // "composite": true,                                /* Enable constraints that allow a TypeScript project to be used with project references. */
    // "tsBuildInfoFile": "./.tsbuildinfo",              /* Specify the path to .tsbuildinfo incremental compilation file. */
    // "disableSourceOfProjectReferenceRedirect": true,  /* Disable preferring source files instead of declaration files when referencing composite projects. */
    // "disableSolutionSearching": true,                 /* Opt a project out of multi-project reference checking when editing. */
    // "disableReferencedProjectLoad": true,             /* Reduce the number of projects loaded automatically by TypeScript. */

    /* Language and Environment */
    "target": "ES6" /* Set the JavaScript language version for emitted JavaScript and include compatible library declarations. */,
    // "lib": [],                                        /* Specify a set of bundled library declaration files that describe the target runtime environment. */
    // "jsx": "preserve",                                /* Specify what JSX code is generated. */
    // "experimentalDecorators": true,                   /* Enable experimental support for legacy experimental decorators. */
    // "emitDecoratorMetadata": true,                    /* Emit design-type metadata for decorated declarations in source files. */
    // "jsxFactory": "",                                 /* Specify the JSX factory function used when targeting React JSX emit, e.g. 'React.createElement' or 'h'. */
    // "jsxFragmentFactory": "",                         /* Specify the JSX Fragment reference used for fragments when targeting React JSX emit e.g. 'React.Fragment' or 'Fragment'. */
    // "jsxImportSource": "",                            /* Specify module specifier used to import the JSX factory functions when using 'jsx: react-jsx*'. */
    // "reactNamespace": "",                             /* Specify the object invoked for 'createElement'. This only applies when targeting 'react' JSX emit. */
    // "noLib": true,                                    /* Disable including any library files, including the default lib.d.ts. */
    // "useDefineForClassFields": true,                  /* Emit ECMAScript-standard-compliant class fields. */
    // "moduleDetection": "auto",                        /* Control what method is used to detect module-format JS files. */

    /* Modules */
    "module": "commonjs" /* Specify what module code is generated. */,
    "rootDir": "./" /* Specify the root folder within your source files. */,
    // "moduleResolution": "node10",                     /* Specify how TypeScript looks up a file from a given module specifier. */
    // "baseUrl": "./",                                  /* Specify the base directory to resolve non-relative module names. */
    // "paths": {},                                      /* Specify a set of entries that re-map imports to additional lookup locations. */
    // "rootDirs": [],                                   /* Allow multiple folders to be treated as one when resolving modules. */
    // "typeRoots": [],                                  /* Specify multiple folders that act like './node_modules/@types'. */
    // "types": [],                                      /* Specify type package names to be included without being referenced in a source file. */
    // "allowUmdGlobalAccess": true,                     /* Allow accessing UMD globals from modules. */
    // "moduleSuffixes": [],                             /* List of file name suffixes to search when resolving a module. */
    // "allowImportingTsExtensions": true,               /* Allow imports to include TypeScript file extensions. Requires '--moduleResolution bundler' and either '--noEmit' or '--emitDeclarationOnly' to be set. */
    // "resolvePackageJsonExports": true,                /* Use the package.json 'exports' field when resolving package imports. */
    // "resolvePackageJsonImports": true,                /* Use the package.json 'imports' field when resolving imports. */
    // "customConditions": [],                           /* Conditions to set in addition to the resolver-specific defaults when resolving imports. */
    // "resolveJsonModule": true,                        /* Enable importing .json files. */
    // "allowArbitraryExtensions": true,                 /* Enable importing files with any extension, provided a declaration file is present. */
    // "noResolve": true,                                /* Disallow 'import's, 'require's or '<reference>'s from expanding the number of files TypeScript should add to a project. */

    /* JavaScript Support */
    // "allowJs": true,                                  /* Allow JavaScript files to be a part of your program. Use the 'checkJS' option to get errors from these files. */
    // "checkJs": true,                                  /* Enable error reporting in type-checked JavaScript files. */
    // "maxNodeModuleJsDepth": 1,                        /* Specify the maximum folder depth used for checking JavaScript files from 'node_modules'. Only applicable with 'allowJs'. */

    /* Emit */
    // "declaration": true,                              /* Generate .d.ts files from TypeScript and JavaScript files in your project. */
    // "declarationMap": true,                           /* Create sourcemaps for d.ts files. */
    // "emitDeclarationOnly": true,                      /* Only output d.ts files and not JavaScript files. */
    // "sourceMap": true,                                /* Create source map files for emitted JavaScript files. */
    // "inlineSourceMap": true,                          /* Include sourcemap files inside the emitted JavaScript. */
    // "outFile": "./",                                  /* Specify a file that bundles all outputs into one JavaScript file. If 'declaration' is true, also designates a file that bundles all .d.ts output. */
    "outDir": "./ts-build" /* Specify an output folder for all emitted files. */,
    // "removeComments": true,                           /* Disable emitting comments. */
    // "noEmit": true,                                   /* Disable emitting files from a compilation. */
    // "importHelpers": true,                            /* Allow importing helper functions from tslib once per project, instead of including them per-file. */
    // "importsNotUsedAsValues": "remove",               /* Specify emit/checking behavior for imports that are only used for types. */
    // "downlevelIteration": true,                       /* Emit more compliant, but verbose and less performant JavaScript for iteration. */
    // "sourceRoot": "",                                 /* Specify the root path for debuggers to find the reference source code. */
    // "mapRoot": "",                                    /* Specify the location where debugger should locate map files instead of generated locations. */
    // "inlineSources": true,                            /* Include source code in the sourcemaps inside the emitted JavaScript. */
    // "emitBOM": true,                                  /* Emit a UTF-8 Byte Order Mark (BOM) in the beginning of output files. */
    // "newLine": "crlf",                                /* Set the newline character for emitting files. */
    // "stripInternal": true,                            /* Disable emitting declarations that have '@internal' in their JSDoc comments. */
    // "noEmitHelpers": true,                            /* Disable generating custom helper functions like '__extends' in compiled output. */
    // "noEmitOnError": true,                            /* Disable emitting files if any type checking errors are reported. */
    // "preserveConstEnums": true,                       /* Disable erasing 'const enum' declarations in generated code. */
    // "declarationDir": "./",                           /* Specify the output directory for generated declaration files. */
    // "preserveValueImports": true,                     /* Preserve unused imported values in the JavaScript output that would otherwise be removed. */

    /* Interop Constraints */
    // "isolatedModules": true,                          /* Ensure that each file can be safely transpiled without relying on other imports. */
    // "verbatimModuleSyntax": true,                     /* Do not transform or elide any imports or exports not marked as type-only, ensuring they are written in the output file's format based on the 'module' setting. */
    // "allowSyntheticDefaultImports": true,             /* Allow 'import x from y' when a module doesn't have a default export. */
    "esModuleInterop": true /* Emit additional JavaScript to ease support for importing CommonJS modules. This enables 'allowSyntheticDefaultImports' for type compatibility. */,
    // "preserveSymlinks": true,                         /* Disable resolving symlinks to their realpath. This correlates to the same flag in node. */
    "forceConsistentCasingInFileNames": true /* Ensure that casing is correct in imports. */,

    /* Type Checking */
    "strict": true /* Enable all strict type-checking options. */,
    // "noImplicitAny": true,                            /* Enable error reporting for expressions and declarations with an implied 'any' type. */
    // "strictNullChecks": true,                         /* When type checking, take into account 'null' and 'undefined'. */
    // "strictFunctionTypes": true,                      /* When assigning functions, check to ensure parameters and the return values are subtype-compatible. */
    // "strictBindCallApply": true,                      /* Check that the arguments for 'bind', 'call', and 'apply' methods match the original function. */
    // "strictPropertyInitialization": true,             /* Check for class properties that are declared but not set in the constructor. */
    // "noImplicitThis": true,                           /* Enable error reporting when 'this' is given the type 'any'. */
    // "useUnknownInCatchVariables": true,               /* Default catch clause variables as 'unknown' instead of 'any'. */
    // "alwaysStrict": true,                             /* Ensure 'use strict' is always emitted. */
    // "noUnusedLocals": true,                           /* Enable error reporting when local variables aren't read. */
    // "noUnusedParameters": true,                       /* Raise an error when a function parameter isn't read. */
    // "exactOptionalPropertyTypes": true,               /* Interpret optional property types as written, rather than adding 'undefined'. */
    // "noImplicitReturns": true,                        /* Enable error reporting for codepaths that do not explicitly return in a function. */
    // "noFallthroughCasesInSwitch": true,               /* Enable error reporting for fallthrough cases in switch statements. */
    // "noUncheckedIndexedAccess": true,                 /* Add 'undefined' to a type when accessed using an index. */
    // "noImplicitOverride": true,                       /* Ensure overriding members in derived classes are marked with an override modifier. */
    // "noPropertyAccessFromIndexSignature": true,       /* Enforces using indexed accessors for keys declared using an indexed type. */
    // "allowUnusedLabels": true,                        /* Disable error reporting for unused labels. */
    // "allowUnreachableCode": true,                     /* Disable error reporting for unreachable code. */

    /* Completeness */
    // "skipDefaultLibCheck": true,                      /* Skip type checking .d.ts files that are included with TypeScript. */
    "skipLibCheck": true /* Skip type checking all .d.ts files. */
  }
}
	`
}

func (c *Content) Air() string {
	return `
  root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/bin"
  cmd = "make build"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go", "_templ.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "templ", "ts"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
  `
}
