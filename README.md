# goor
This is a CLI for creating constructors, getters and setters for structures, and is only for Golang.

# Installation

```shell script
$ go get -u github.com/ttakuya50/goor
```
If you want the direct binary data, please click [here](https://github.com/ttakuya50/goor/releases).

# Usage

```shell script
Usage of goor:
        goor [flags] -type=[struct type name]
For more information, see:
        https://github.com/ttakuya50/goor
Flags:
  -getter
        [option] when you create a constructor, you also create a getter.
  -output string
        [option] output file name (default "srcdir/<type>_constructor_gen.go").
  -pointer
        [option] set the return value to a pointer when creating the constructor. (default true)
  -setter
        [option] when you create a constructor, you also create a setter.
  -type string
        [required] a struct type name.
  -version
        outputs the current version.
```

Frequently used settings are listed below.
## 1. Create constructor.

```go
//go:generate goor -type=User
type User struct {
	id   int
	name string
	age  int
}
```

```shell script
$ go generate ./...
```

A file named user_constructor_gen.go will be created.
Constructor is generated.
```go
func NewUser(
	id int,
	name string,
	age age,
) *User {
	return &User{
		id:   id,
		name: name,
		age:  age,
	}
}
```

## 2. Create constructor and getter.

```go
//go:generate goor -type=User -getter
type User struct {
	id   int
	name string
	age  int
}
```

```shell script
$ go generate ./...
```

A file named user_constructor_gen.go will be created.
Constructor and getter are generated.
```go
func NewUser(
	id int,
	name string,
	age int,
) *User {
	return &User{
		id:   id,
		name: name,
		age:  age,
	}
}

func (g *User) Id() int {
	return g.id
}

func (g *User) Name() string {
	return g.name
}

func (g *User) Age() int {
	return g.age
}
```

## 3. Create constructor and getter and setter.

```go
//go:generate goor -type=User -getter -setter
type User struct {
	id   int
	name string
	age  int
}
```

```shell script
$ go generate ./...
```

A file named user_constructor_gen.go will be created.
Constructor and getter and setter are generated.
```go
func NewUser(
	id int,
	name string,
	age int,
) *User {
	return &User{
		id:   id,
		name: name,
		age:  age,
	}
}

func (g *User) Id() int {
	return g.id
}

func (g *User) Name() string {
	return g.name
}

func (g *User) Age() int {
	return g.age
}

func (s *User) SetId(id int) {
	s.id = id
}

func (s *User) SetName(name string) {
	s.name = name
}

func (s *User) SetAge(age int) {
	s.age = age
}
```

## 4. It is also possible to exclude certain structure variables.

```go
//go:generate goor -type=User
type User struct {
	id   int `goor:"constructor:-"`
	name string
	age  int
}
```

```shell script
$ go generate ./...
```

It is also possible to exclude certain structure variables.
```go
func NewUser(
	name string,
	age int,
) *User {
	return &User{
		name: name,
		age:  age,
	}
}
```

## 5. It is also possible to exclude the constructor, getter and setter of a specific structure variable.

```go
//go:generate goor -type=User -getter -setter
type User struct {
	id   int    `goor:"getter:-"`
	name string `goor:"constructor:-;setter:-"`
	age  int    `goor:"constructor:-;getter:-;setter:-"`
}
```

```shell script
$ go generate ./...
```

The following code will be generated on its own.
```go
func NewUser(
	id int,
) *User {
	return &User{
		id: id,
	}
}

func (g *User) Name() string {
	return g.name
}

func (s *User) SetId(id int) {
	s.id = id
}
```

## 5. Change the name of the output file.

```go
//go:generate goor -type=User -output=user_gen.go
type User struct {
	id   int
	name string
	age  int
}
```

```shell script
$ go generate ./...
```

The following code will be generated in the file named user_gen.go
```go
func NewUser(
	id int,
	name string,
	age int,
) *User {
	return &User{
		id:   id,
		name: name,
		age:  age,
	}
}
```

# Examples
Check out the examples of how to generate code using goor in the examples directory.

# Author
[t_takuya50](https://twitter.com/t_takuya50)