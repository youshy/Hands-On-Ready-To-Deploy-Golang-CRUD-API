# Hands-on Golang CRUD API

Powered by [ECS Digital](https://ecs.co.uk/digital-engineering/) and [DevOps Playground](https://www.meetup.com/DevOpsPlayground/)

Quick and straight to the point - let's build a quick CRUD API using Go!

<a name="top">

# Table Of Content

* [Intro](#intro)
* [What is Go](#go-intro)
* [Why Go?](#go-why)
* [My setup](#my-setup)
* [Small catch](#catch)
* [Initial app setup](#initial-app-setup)
* [Setup entrypoint](#entrypoint)
* [Initial server](#initial-server)
* [Broker](#broker)
* [Handlers](#handlers)
* [More handlers](#more-handlers)
* [Finish server](#finish-server)
* [Test the server](#test-server)
* [Dockerfile](#dockerfile)
* [How to improve](#improve)
* [Word from the author](#author)

[^Top](#top)

<a name="intro"/>

## Intro

Hello and welcome to our second, fully virtual DevOps playground! Hope you all well and prepared for a little bit over an hour of fun!

Depending on your preference, you can either:

* Fork/copy this repository locally and just follow along as we go - code in this repo is fully finished app
* Get one of our ready-to-Go instances and code along

This playground assumes that you know nothing about Go and will explain a lot of stuff to some detail - if you want to learn more, I recommend [A Tour of Go](https://tour.golang.org) and [The Go Programming Language](https://www.amazon.co.uk/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440) by A. Donovan and B. Kernighan.

Let's go!

[^Top](#top)

(From this point the README assumes that you're using our Go instance. If you want to follow the instructions locally, go to [golang.org](https://golang.org) and download a build for your machine.)

<a name="go-intro"/>

## What Is Go

**Go** (or known also as **golang** because of it's domain) is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Mr Pike and Mr Thompson designed, developed and implemented the original Unix system; Mr Griesemer was working on Chrome's V8 JavaScript engine. And that's just a tip of the iceberg of their collective experience.

Nowadays Go is a backbone of multitude of DevOps tools - from Docker, to Terraform or Vault to Kubernetes. Also Prometheus, Helm, Loki, Grafana... The list goes on and on.

[^Top](#top)

<a name="go-why"/>

## Why Go?

Every tech person should know at least one programming language that allows them a total freedom in tech world. In DevOps world, most of the time it's Python; some of us know Ruby or JavaScript.

Go gives you the type safety, allows you for low-level programming, allows you to cut corners in production-friendly way; the best example would be, that you can create a load balancer system from ground up using Go in-built concurrency support.

The thing that would be the most interesting to us now is that Go is the perfect candidate for high-performance web servers - it can run as one straight out of the box, but today we'll use some extra packages to take our code to the next level.

Also, the very good thing about go is it's documentation. If there's something you won't understand after our playground, then go to [golang.org](https://golang.org) and go over the documentation for a package. Answers will be there!

[^Top](#top)

<a name="my-setup"/>

## My Setup

I'm using Vim exclusively; you can try using amazing GoLand IDE by JetBrains. 

For syntax highlight, automatic imports and autoformatting I'm using [vim-go](https://github.com/fatih/vim-go) plugin for Vim but I'll do my best to import everything manually!

[^Top](#top)

<a name="catch"/>

## Small catch

This API will most definitely show some of my preferences in writing code - which are neither bad or good. If there's anything that I do out of preference, it'll be highlighted in this README.

Also, the app will be written without any tests and won't follow TDD practices - because we don't want to run this playground till dawn.

[^Top](#top)

<a name="initial-app-setup"/>

## Initial App Setup

Let's create a new folder for our app and enter it:

```
mkdir go-crud-api && cd go-crud-api
```

For our app to work we need to initialize **Go Modules**. Go Modules is the way how Go manages it's dependencies. 

In your terminal type:

```
go modules init go-crud-api
```

This will create new file `go.mod` in our folder. That's the last thing we'll have to do ANYTHING with our dependencies.

If you have any problems here, do this:

```
export GO111MODULE=on
```

To ensure that your app will use Go Modules.

Our app will use four environmental variables for connecting to Amazon RDS PostgreSQL* instance:

* `PG_USERNAME` - username for our database
* `PG_PASSWORD` - password
* `PG_DB_NAME` - name of our database
* `PG_DB_HOST` - database host

That's all we have to do in our initial setup!

> * Our instance is predefined with dev_ops_playground database. If you're using your own, make sure that your Postgres instance has dev_ops_playground database within.

[^Top](#top)

<a name="entrypoint"/>

## Setup entrypoint

Each Go application requires `main.go` file to run. `main.go` is our entrypoint for the app - whenever we'll run it or compile it, `main.go` has to have all the logic we need for the app to start. Let's write our first Go code!

**main.go**
```go
package main

import (
  "log"
  "os"
)

func main() {
  a := App{}
  a.Initialize()
  a.Run(":9000")
}
```

So, line by line:

**main.go**
```go
package main
```

If we're building executable program, we need to use `main` as our package name. `main` tells the Go compiler that the package should compile as an executable instead of a shared library. If we'd be building a microservices mesh, we might use different names for packages to decouple the code and use it all as a library.

**main.go**
```go
import (
  "log"
  "os"
)
```

`import` statement imports necessary packages to our code. Here we're importing `log` and `os` from standard library - these are the libraries included in Go.

**main.go**
```go
func main() {
```

`func` is the way we define functions in Go. `func main()` is our entry point for the application - anything that is needed for the app to run has to go here. Let's take a look at what we have:

**main.go**
```go
func main() {
  a := App{}
  a.Initialize()
  a.Run(":9000")
}
```

Line by line:

**main.go**
```go
a := App{}
```

`a` is the name of our variable. `:=` is the way we initialize variable and assign the value of right-hand statement. `App{}` is our struct (object) that we'll create in the `server.go`.

**main.go**
```go
a.Initialize()
```

Remember when I've told you that there might be some preferences here? `a.Initialize()` is a method from our `App{}` struct that will perform all the logic before the server runs - establish handlers and connect to PostgreSQL instance.

This could be easily done in `main()` but I prefer keeping my entry function as clean as possible and abstract any extra logic outside.

**main.go**
```go
a.Run(":9000")
```

And last, but not least `a.Run()` method will start the server. The parameter we use here is the port we'll use for the app. This can be changed to anything or even exported as an environmental variable.

About these variables - we have them, but we don't know if they exist before we run the app. We need to check that. In our `main.go` let's add:

**main.go**
```go
func init() {
	if ok := os.Getenv("PG_USERNAME"); ok == "" {
		log.Fatalln("PG_USERNAME not specified")
	}
	if ok := os.Getenv("PG_PASSWORD"); ok == "" {
		log.Fatalln("PG_PASSWORD not specified")
	}
	if ok := os.Getenv("PG_DB_NAME"); ok == "" {
		log.Fatalln("PG_DB_NAME not specified")
	}
	if ok := os.Getenv("PG_DB_HOST"); ok == "" {
		log.Fatalln("PG_DB_HOST_ not specified")
	}
}
```

Woah, a lot of code. Don't panic, we'll get through it.

`func init()` is a in-built function that runs before `main()` runs. It's a perfect candidate for our app to check if our necessary environmental variables exists. As you can see, we repeat the whole thing four times, so let's take one of these statements and inspect it:

**main.go**
```go
if ok := os.Getenv("PG_PASSWORD"); ok == "" {
  log.Fatalln("PG_PASSWORD not specified")
}
```

`if` is the classic **if** statement we know from any other programming language. We'll see a lot of it in Go. So, what we have here:

* `if ok :=` starts the statement, creates a variable `ok` and assigns the value of
* `os.Getenv("PG_PASSWORD");` - `os.Getenv` is a function from `os` package that checks for an env variable named `PG_PASSWORD`
* `ok == ""` - we could easily spread the whole line into:

**main.go**
```go
ok := os.Getenv("PG_PASSWORD")
if ok == "" {
```

But in our case, it's makes way more sense to keep it as a one-liner. Here we check if the `ok` variable is an empty string. `os.Getenv` always return a string - if the variable exists if will be something; if not, then it'll be an empty string - `""`.

If the variable is an empty string then:

**main.go**
```go
log.Fatalln("PG_PASSWORD not specified")
```

`log.Fatalln` is a function from `log` package that prints something to `stdout` and is followed by a call to `os.Exit(1)`. So, in other terms - it prints our error and then exits the program.

Huh, we've done a lot! Our `main.go` file should look like this now:

**main.go**
```go
package main

import (
	"log"
	"os"
)

func main() {
	a := App{}
	a.Initialize()
	a.Run(":9000")
}

func init() {
	if ok := os.Getenv("PG_USERNAME"); ok == "" {
		log.Fatalln("PG_USERNAME not specified")
	}
	if ok := os.Getenv("PG_PASSWORD"); ok == "" {
		log.Fatalln("PG_PASSWORD not specified")
	}
	if ok := os.Getenv("PG_DB_NAME"); ok == "" {
		log.Fatalln("PG_DB_NAME not specified")
	}
	if ok := os.Getenv("PG_DB_HOST"); ok == "" {
		log.Fatalln("PG_DB_HOST_ not specified")
	}
}
```

[^Top](#top)

<a name="initial-server"/>

## Initial server

[^Top](#top)

<a name="broker"/>

## Broker

[^Top](#top)

<a name="handlers"/>

## Handlers

[^Top](#top)

<a name="more-handlers"/>

## More Handlers

[^Top](#top)

<a name="finish-server"/>

## Finish server

[^Top](#top)

<a name="test-server"/>

## Test the server

[^Top](#top)

<a name="dockerfile"/>

## Dockerfile

[^Top](#top)

<a name="improve"/>

## How To Improve

[^Top](#top)

<a name="author"/>

## Word from the author

[^Top](#top)
