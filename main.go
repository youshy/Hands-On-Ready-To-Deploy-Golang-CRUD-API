/*
Hello and welcome to DevOps Playground!

This time, we're tackling the task of building a full RESTful, CRUD API.

If you need help with the syntax, check cheat sheet.

ALl files on dev branch are heavily commented; if you want production-ready api - switch to master.
*/

// Package is a keyword required for the compiler to understand what to do with the files.
// If the package = main, then compiler takes that as an executable.
// Meaning - when compiling, will create an executable file bundled of all files referencing package main.
package main

// Import, well, imports stuff.
// If the package doesn' start with any host, then it's a part of standard library.
// Go's standard library is insanely powerful - we could build the whole API only using these.
// In our case we want to extend the functionality - hence, we'll use 3rd party libraries as well.
import (
	"log"
	"os"
)

// func main() is the entrypoint of the application.
// Anything here will be executed when running the program.
func main() {
	// Create and initialize a variable of type App.
	a := App{}
	// Both methods below are methods of App struct.
	// Defined in server.go
	a.Initialize()
	a.Run(":9000")
}

// func init() is a function that runs before main()
// Than's why it's the perfect candidate for our env vars check
func init() {
	// This syntax might look weird, but below I'll write how it can also look
	if ok := os.Getenv("PG_USERNAME"); ok == "" {
		log.Fatalln("PG_USERNAME not specified")
	}
	/*
		ok := os.Getenv("PG_USERNAME")
		if ok == "" {
			log.Fatalln("PG_USERNAME not specified")
		}
	*/
	// Both statements here and the one above are the same.
	// We initialize a variable which we get by calling os.Getenv,
	// then we check if that variable is empty.
	// Go doesn't have a NULL concept, it has a zero (nil) value for each type.
	// In case of strings, the nil value is "" (empty string)
	if ok := os.Getenv("PG_PASSWORD"); ok == "" {
		// If the variable is "", then print log to stdout
		// and call os.Exit(1)
		// Fatalln is a really nice wrapper around that functionality
		log.Fatalln("PG_PASSWORD not specified")
		// you could do something like this:
		/*
			log.Println("PG_PASSWORD not specified")
			os.Exit(1)
		*/
	}
	// TODO: add check for PG_DB_NAME
	if ok := os.Getenv("PG_DB_HOST"); ok == "" {
		log.Fatalln("PG_DB_HOST_ not specified")
	}
}
