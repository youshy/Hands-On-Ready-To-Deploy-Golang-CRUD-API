# Vim

## Save and/or quit

`:q!` - quit without saving
`:wq` - save and quit
`:w` - just save
`:q` - quit if saved

## Move around

For the sake of our playground, we've enabled arrow keys in our Vim instances.
If you want to be a pro:

`h` - left
`j` - down
`k` - up
`l` - right

## Modes

`v` - visual mode, lets you highlight stuff
`i` - insert mode, lets you input stuff
`esc` - magic key, switches to normal mode where you can move around

# Go

Bear in mind, this cheat sheet is really short. If you want to learn more, check [A Tour of Go](https://tour.golang.org).

Go is statically typed language - meaning that each variable has to have a pre-defined type; either by the user or by right-hand value/return value.

## Variables

`var msg string` as:

* `var` - keyword for initializing variable
* `msg` - name of the variable
* `string` - type of the variable

`msg = "Hello"` - assign value to pre-initialized variable

shorthand:

`msg := "Hello"` - initialize variable, assing the value. This way variable inherits type of the value

## Types

### Strings

`str := "Hello"`

```go
str := `Multiline
string`
```

### Numbers

```go
// int
num := 3 // int

// float
num := 3. // float 64
```

### Arrays

In Go arrays have a fixed size

```go
numbers := [...]int{0, 1, 2, 3, 4}
```

`[...]` means that the size will be set by values in `( )`.

### Slices

Think about them as dynamic arrays.

```go
slice := []int{1, 2, 3}
```

### Pointers

This is a very VERY long topic but let's put it this way - pointers are the reference to a memory address of a variable.

* `var p *int` - this means that p is a type of pointer value to int
* `&` - generates a pointer to it's operand. You can use it in shorthand variable assignment.
* `*` - denotes(references) a value of the pointer

### Interfaces

```go
interface{}
```

this is a magical type that will inherit any type. Very useful in JSON.

## Logging and printing

You can use one of two, both have the same methods:

```go
log.Print(value) // prints nicely logged values
fmt.Print(value) // prints straight to stdout
```

## Functions

Define function as:

```go
// second () is a return type
// you can do named returns although it's not the best practice
func FunctionName(parameter string) (string) {

}
```

## Structs and OOP

So, here's a thing - if you come from Java, JavaScript or even Python: Go doesn't support OOP. In a normal way. 

In Go we don't create classes, we create structs that will have methods.

```go
// this will create a struct and define the internal variables and their types
type Example struct {
  // if the struct field is capitalised
  // it will be exported
  // meaning - will be available outside of the struct
  Number int
  Field string
  // nested field - think about this as a type inheritance
  // composition in Go is another topic
  Nested Nested
  // if the struct field is lowercased
  // it won't be exported
  // meaning - it won't be accessible outside of the struct
  internals string
}

// Initialize type
t := Example{}

log.Print(t.internals) // not accessible, won't even compile
log.Print(t.Number) // will print the number

// this function is a method of Example type
// note (e *Example) - that means that this function references the type
// so we'd be able to use fields that are not exported
func(e *Example) WillPrintUnexportedField() {
  log.Print(e.internals)
}

t.WillPrintUnexportedField() // prints to stdout whatever the internals field value is
```


