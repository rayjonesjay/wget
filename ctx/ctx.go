// Package ctx defines the circumstances that form the setting for a download event.
//
// For consistency, implementations are advised to use the Context struct to
// retrieve processed command line arguments; and as a reference for a ctx
// for their implementation.
//
// The Context struct, is defined here, in a separate
// package, so that it can be imported by any subpackages in the project without
// circular dependencies; for this reason, no imports of the project's subpackages
// are allowed in this package, unless if such subpackage themselves promise that
// they must not import this package as a dependency.
//
// As such, adding receiver functions to the Context struct is thus restricted (due to the import
// restrictions). To create receiver functions that operate on a Context instance,
// it is recommended to create a wrapper struct, where you are then allowed to
// create receivers on.
//
// This can be achieved using Struct Composition/Embedding in go.
// The example below, demonstrates embedding the struct (Context) within another struct (Arg).
//
//		package main
//
//		import "fmt"
//
//		type Arg struct {
//	 		*ctx.Context
//		}
//
//		// Receiver function to set the OutputFile of the embedded Context
//		func (a *Arg) SetOutputFile(filename string) {
//		    a.OutputFile = filename // Accessing and modifying Context's field directly
//		}
//
//		// Receiver function to check if mirroring is enabled
//		func (a *Arg) IsMirroringEnabled() bool {
//		    return a.Mirror // Accessing Context's field directly
//		}
//
//		func main() {
//		    arg := Arg{
//		        Context: &Context{}, // Initialize the embedded Context
//		    }
//
//		    // Using receiver functions
//		    arg.SetOutputFile("output.txt")
//		    fmt.Println("Output file:", arg.OutputFile) // Accessing Context's field directly
//
//		    arg.Mirror = true
//		    fmt.Println("Mirroring enabled:", arg.IsMirroringEnabled()) // Using a receiver function
//		}
//
// In this example:
//
// 1. We define receiver functions `SetOutputFile` and `IsMirroringEnabled` on
// the `Arg` struct. 2. Inside these functions, we directly access and modify the
// fields of the embedded `Context` struct using `a.OutputFile` and `a.Mirror`.
// 3. In the `main` function, we create an instance of `Arg` and initialize the
// embedded `Context`. 4. We then use the receiver functions to set the
// `OutputFile` and check the `Mirror` status. 5. We can also directly access the
// `Context` fields like `arg.OutputFile`.
//
// ### Indepth Explanation:
//
// * `Composition:` The `Arg` struct is composed of the `Context` struct by
// embedding it as an anonymous field (using `*Context`). This means that `Arg`
// "inherits" all the fields and methods of `Context` without explicitly listing
// them. * `Embedding:` The `*Context` field within `Arg` is an embedded field,
// also referred to as an anonymous field. It doesn't have an explicit name but
// directly brings in the fields of `Context` into the scope of `Arg`.
//
// #### How it works:
//
// * `Accessing fields:` You can access the fields of `Context` through `Arg`
// directly as if they were declared within `Arg`. For example,
// `ctx.OutputFile` would access the `OutputFile` field of the embedded
// `Context` struct. * `Method promotion:` If `Context` has methods defined,
// those methods are also "promoted" to `Arg`. You can call them directly on
// `Arg` instances as if they were defined on `Arg`.
//
// #### Benefits of Composition:
//
// * `Code reuse:` Avoids code duplication by inheriting fields and methods from
// the embedded struct. * `Flexibility:` Allows you to customize behavior by
// adding additional fields or methods to the composing struct (e.g., `Arg` in
// this case). * `Maintainability:` Changes to the embedded struct (e.g.,
// `Context`) automatically propagate to the composing struct, making it easier
// to maintain consistency. * `Clearer relationships:` Represents a "has-a"
// relationship between structs (e.g., `Arg` "has an" `Context`).
//
// #### In this specific example:
//
// * The `Arg` struct likely represents the overall ctx or configuration of
// an application. * The `Context` struct encapsulates the command-line arguments
// parsed from the user. * By embedding `Context` into `Arg`, the application can
// easily access the command-line arguments through the `Arg` instance without
// needing to pass `Context` around separately.
//
// #### Difference from Inheritance:
//
// Composition is often favored over traditional inheritance in Go because it
// provides more flexibility and avoids the tight coupling that inheritance can
// create. Go doesn't have traditional class-based inheritance.
package ctx

// Context defines the circumstances that form the setting for a download event,
// as specified in commandline arguments passed by the user.
//
// For example: when the following command is executed
//
//	$ go run . -O=file.txt -B https://www.example.com
//
// The download Context will be set such that:
// * the field OutputFile will have the value "file.txt",
// which specifies the output file to save the resource into.
// * the field BackgroundMode will have the value true,
// which specifies that the download will be done in the background,
// and that the output will be redirected to a log file
// * the field Links will have an element, "https://www.example.com",
// which specifies the link from where to download the file
type Context struct {
	// identified by the -O flag
	OutputFile string
	// identified by the -B flag
	BackgroundMode bool
	// identified by the regexp pattern (http|https)://\w+ ,specifies path to resources on the web
	Links []string
	// identified by the -P flag, specifies the location where to save the resource
	SavePath string
	// identified by the -i flag, specifies a file contains url(s)
	InputFile string
	// identified by the -R or --reject flag contains a list of resources to reject
	Rejects []string
	// identified by the --mirror flag, indicates whether to download an entire website or not
	Mirror bool
	// identified by the --rate-limit flag, specifies the download speed when fetching a resource
	RateLimit string
	// if RateLimit is specified, RateLimitValue will be
	RateLimitValue int64
	// identified by the --help flag, if pared it will print our program manual
	IsHelp bool
	// identified by the --convert-links
	ConvertLinks bool
	// identified by the --exclude or -X, takes a comma separated list of paths (directory),
	//to avoid when fetching a resource
	Exclude []string
}
