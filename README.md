<h1 align="center">
  <a href="#">
    <picture>
      <source height="125" media="(prefers-color-scheme: dark)" srcset="./assets/logo.jpg" style="border-radius:15px">
      <img height="125" alt="Hive" src="./assets/logo.jpg" style="border-radius:15px">
    </picture>
  </a>
  <br>
  
  
</h1>
<p align="center">
  <em><b>Hive</b> is an <a href="https://nestjs.com/">NestJS</a> inspired <b>web framework</b> built on top of <a href="https://github.com/gofiber/fiber">Fiber</a> and <a href="https://github.com/valyala/fasthttp">Fasthttp</a>, the <b>fastest</b> HTTP engine for <a href="https://go.dev/doc/">Go</a>. Designed to <b>ease</b> things up for <b>fast</b> development with <a href="https://docs.gofiber.io/#zero-allocation"><b>zero memory allocation</b></a> and <b>performance</b> in mind.</em>
</p>

---

## ⚙️ Installation


```bash
go mod init github.com/your/repo
```

To learn more about Go modules and how they work, you can check out the [Using Go Modules](https://go.dev/blog/using-go-modules) blog post.

After setting up your project, you can install Hive with the `go get` command:

```bash
go get -u github.com/hive-go/hive
```

This command fetches the Hive package and adds it to your project's dependencies, allowing you to start building your web applications with Hive.

## ⚡️ Quickstart

Getting started with Hive is easy. Here's a basic example to create a simple web server. This example demonstrates initializing a new Fiber app, setting up a route for user, and starting the server.



Suggestion folder structure
```bash
src
 ├── modules
 │   └── user
 │       ├── user.module.go
 │       ├── user.controller.go
 │       └── user.service.go
 ├── main.go
 ├── go.mod
 └── go.sum
```

`main.go`
```go
package main

import (
	"github.com/your/repo/src/modules/user"
	"github.com/hive-go/hive"
)

func main() {
	app := hive.New()

	app.AddModule(user.UserModule)

	app.Listen("3000")
}
```

`src/modules/user/user.module.go`

```go
package user

import (
	"github.com/hive-go/hive"
)

var UserModule = hive.CreateModule()

func init() {
	UserModule.AddController(UserController)
}
```

`src/modules/user/user.controller.go`

```go
package user

import (
	"github.com/hive-go/hive"
)

var UserController = hive.CreateController()

func init() {
	UserController.Get("/user", func(c *hive.Ctx) any {
		return UserService.GetUser("123")
	})

	UserController.Post("/user", func(c *hive.Ctx) any {
		return UserService.CreateUser(c)
	})

	UserController.Put("/user", func(c *hive.Ctx) any {
		return UserService.UpdateUser(c)
	})

	UserController.Delete("/user", func(c *hive.Ctx) any {
		return UserService.DeleteUser(c)
	})
}
```

`src/modules/user/user.service.go`

```go
package user

import (
	"github.com/hive-go/hive"
)

type UserServiceT struct{}

var UserService = UserServiceT{}

func (u *UserServiceT) GetUser(id string) hive.Map {
	return hive.Map{
		"user": hive.Map{"name": "John Doe"},
	}
}

func (u *UserServiceT) CreateUser(c *hive.Ctx) string {
	return "User Created"
}

func (u *UserServiceT) UpdateUser(c *hive.Ctx) string {
	return "User Updated"
}

func (u *UserServiceT) DeleteUser(c *hive.Ctx) string {
	return "User Deleted"
}
```


## ⚙️ Using Hive CLI to Work Faster


```bash
go install github.com/hive-go/hive-cli/
```

## ⚡️ Quickstart

Getting started with Hive Cli is easy. Here's a basic example to create User Module

```bash
  hive-cli generate_resource user
```

Result in folder structure
```bash

 ├── src
 │   └── modules
 │       └── user
 │           ├── user.module.go
 │           ├── user.controller.go
 │           └── user.service.go
 ├── main.go
 ├── go.mod
 └── go.sum
```

 <a href="https://github.com/hive-go/example-project">
📚 Show more code examples
 </a>








This simple server is easy to set up and run. It introduces the core concepts of Hive: app initialization, route definition, and starting the server. Just run this Go program, and visit `http://localhost:3000/user` in your browser to see the message.

 <a href="https://github.com/hive-go/example-project">

📚 Show more code examples
 </a>




