# Hands-on Golang CRUD API

Powered by [ECS Digital](https://ecs.co.uk/digital-engineering/) and [DevOps Playground](https://www.meetup.com/DevOpsPlayground/)

Quick and straight to the point - let's build a quick CRUD API using Go!

Cheat sheet -> [Here](https://gist.github.com/youshy/a92020e228ef5a164a75be4733650ad7)

<a name="top">

---

# DISCLAIMER

This readme differs a lot from the one you can find on `master` branch. Hands-on part of this playground is written down both in here and in comments in the code.

If you want to build the API from scratch, follow the `master` branch.

---

# Table Of Content

* [Initial app setup](#initial-app-setup)
* [Setup entrypoint](#entrypoint)
* [Initial server](#initial-server)
* [Handlers](#handlers)
* [Dockerfile](#dockerfile)
* [Test the server](#test-server)
* [How to improve](#improve)
* [Word from the author](#author)

[^Top](#top)

<a name="initial-app-setup"/>

## Initial App Setup

For our app to work we need to initialize **Go Modules**. Go Modules is the way how Go manages it's dependencies. 

In your terminal type:

```
go mod init go-crud-api
```

This will create new file `go.mod` in our folder. That's the last thing we'll have to do ANYTHING with our dependencies.

If you have any problems here, do this:

```
export GO111MODULE=on
```

To ensure that your app will use Go Modules.

Our app will use four environmental variables for connecting to Amazon RDS PostgreSQL instance:

* `PG_USERNAME` - username for our database
* `PG_PASSWORD` - password
* `PG_DB_NAME` - name of our database
* `PG_DB_HOST` - database host

Also, for the sake of usability, we've set up your instances with `vim-go` plugin that will help us with highlighting code. We need to activate it by using plug:

In Vim, type:

`:PlugInstall`

That's all we have to do in our initial setup!

[^Top](#top)

<a name="entrypoint"/>

## Setup entrypoint

Each Go application requires `main.go` file to run. `main.go` is our entrypoint for the app - whenever we'll run it or compile it, `main.go` has to have all the logic we need for the app to start. Let's write our first Go code!

If we're building executable program, we need to use `main` as our package name. `main` tells the Go compiler that the package should compile as an executable instead of a shared library. If we'd be building a microservices mesh, we might use different names for packages to decouple the code and use it all as a library.

In `main.go` you can find this part:

**main.go**
```go
func init() {
	if ok := os.Getenv("PG_USERNAME"); ok == "" {
		log.Fatalln("PG_USERNAME not specified")
	}
	if ok := os.Getenv("PG_DB_NAME"); ok == "" {
		log.Fatalln("PG_DB_NAME not specified")
	}
	if ok := os.Getenv("PG_DB_HOST"); ok == "" {
		log.Fatalln("PG_DB_HOST_ not specified")
	}
}
```

`func init()` is a in-built function that runs before `main()` runs. It's a perfect candidate for our app to check if our necessary environmental variables exists. We need to add one more check to the app:

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

[^Top](#top)

<a name="handlers"/>

## Handlers

Add the handler for updating posts:

**server.go**
```go
router.Handle(prefix+"/post/{post_id}", a.UpdatePost()).Methods(http.MethodPut)
```

And let's write the logic for the missing GetSinglePost handler well:

**handlers.go**
```go
func (a *App) GetSinglePost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var post Post
		defer r.Body.Close()
		err := db.Table("posts").Where("id = ?", vars["post_id"]).First(&post).Error
		if err != nil {
			log.Printf("get single post %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		JSONResponse(w, http.StatusOK, post)
	})
}
```

<a name="test-server"/>

## Test the server

In the main directory of our app type `go run .`. This is what you should see:

```
2020/05/18 12:54:20 Available routes:
[GET]	/api/post
[GET]	/api/post/{post_id}
[POST]	/api/post
[PUT]	/api/post/{post_id}
[DELETE]	/api/post/{post_id}
2020/05/18 12:54:20 Server is listening on :9000
```

`go run .` is a command that compiles and runs the named main Go package. But let's kill the process (ctrl+C) and build the executable by typing `go build -o crud-api`. `-o` flag names the executable.

So now, we can do `./crud-api` and we should get the same output!

Now, let's run another terminal and test all the routes:

(bear in mind - your id's will be way different than mine)

### Get Posts
`curl localhost:9000/api/post`

### Get Single Post
`curl localhost:9000/api/post/{id}`

### Create New Post
`curl -X POST -d '{"title":"my new post", "content": "such a good writer"}' localhost:9000/api/post`

### Modify Post
`curl -X PUT -d '{"content": "this is better"}' localhost:9000/api/post/{id}`

### Delete Post
`curl -X DELETE localhost:9000/api/post/{id}`

[^Top](#top)

<a name="dockerfile"/>

## Dockerfile

I've promised we'll make the API deployable so let's add the dockerfile:

```dockerfile
FROM golang
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o crud-api
EXPOSE 9000
ENTRYPOINT ["./crud-api"]
```

[^Top](#top)

<a name="improve"/>

## How To Improve

There's a few things to improve in our app:

* The logic for DB actions could be abstracted from the handlers
* Handlers could have more robust checks
* There's no tests written for the API (apart from our simple `curl` commands)
* Postgres might be abstracted into an interface - so we could quickly plug another database engine

Of course there's more than this, but as an API written from ground up in a few hours I think it's better than ok!

[^Top](#top)

<a name="author"/>

## Word from the author

Hope you've learned something about Go in the process and I was able to showcase the strengths and small quirks of writing in Go!

If you have any questions or problems, catch me on [LinkedIn](https://www.linkedin.com/in/arturkondas/), [Twitter](https://twitter.com/arturkondas) or [Github](https://github.com/youshy)

Read more about me, my stories and thoughts and find other Go tutorials @ [akondas.com](https://akondas.com)

[^Top](#top)
