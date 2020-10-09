Tic Tac Toe

In this tutorial you will learn how to:

1. Build a web based Tic Tac Toe game in Go.
1. Write a web based client in vanilla JavaScript.
1. Implement the minimax algorithm to recursively calculate the most optimal next move for the computer.
1. Write unit and integration tests to exercise and prove correctness of the logic.

## Pre-requisites

The following pre-requisites should be in place to successfully follow this guide.

1. Go lang should be installed on your computer.
1. You should have a basic understanding of Go, HTML, CSS and JavaScript.

## Introduction

Tic Tac Toe otherwise known as "Naughts and Crosses" is a really basic board game familiar to everyone the World over. 
When I was a young kiddie, Tic Tac Toe captured my imagination in the 1983 film: [WarGames](https://youtu.be/F7qOV8xonfY).
At the time I was too young to understand what was meant by "The only winning move is not to play" because as a clueless
child playing against a computer I often messed-up and lost. It became clear to me as I grew that the best result you can
hope for is a draw.

<a href="http://www.youtube.com/watch?feature=player_embedded&v=F7qOV8xonfY" target="_blank"><img src="http://img.youtube.com/vi/F7qOV8xonfY/0.jpg" 
alt="WarGames" width="240" height="180" border="10" /></a>

So then: why Tic Tac Toe? 
In the past, when learning a new language I've often found myself going down the path of building Tic Tac Toe to explore
the features and nuances of that language. 
To anyone progressing beyond a basic Hello World app, Tic Tac Toe still has quite a few fundamental concepts to learn. 
With the fundamentals covered, you will be poised to tackle more advanced concepts. 
In this tutorial we will build Tic Tac Toe and thus cover these fundamental concepts, namely:
1. GUI. The player needs to be able to visualize the game.  In the context of this tutorial, we will implement a simple,
HTML/CSS/Vanilla JavaScript web page to render the game. 
1. Interactivity. The player needs to be able to interact with the system to advance the state of the game. In the
context of this tutorial, we will implement a basic HTTP web server to handle requests sent from the browser to the Go
backend. 
1. Game logic. The rules of the game need to be applied and adhered to. In the context of this tutorial, this is the 
backend game logic that:
    - Validates the incoming request from the browser.
    - Calculates the most optimal move for the computer to make.
    - Determines whether the game has been concluded and calculates the result.
1. Testing. There are many benefits to writing tests. Apart from preventing embarrassing bugs from making it to 
production, for me writing tests streamlines the development process and makes it more enjoyable to code.  Tests in Go 
are first class citizens of your source, and the language provides a neat framework to run tests. Some noteworthy 
features are code coverage reporting, race condition detection and benchmark testing and profiling. Beyond writing basic
unit and integration tests, we won't cover these topics in this tutorial.  This should give you an awareness of what 
features are available.

It's worth noting that the entire game could be implemented in the browser using plain JavaScript since (in this case) 
the backend isn't making any third-party API calls, accessing databases or interacting with other microservices i.e. 
things that backend services usually do, but for the sake of exploring the Go language let's pretend that we're building
the War Operation Plan Response system featured in the movie, WarGames.

## Initializing the project

Create a new repository in Github and initialize it with the .gitignore template tailored for Go and add a README 
file and set the license if you wish.  Clone the new repository to your hard drive and change directory into the
project directory. We will use Go Modules to manage the projects dependencies.  To initialize a new project with a 
`go.mod` file, run the following command:

```
go mod init github.com/<your git hub account>/tictactoe
```

Next edit your .gitignore file in your IDE and unhash `#vendor/` and "fix" the comment above it, also add an entry 
for `.env`.  The last lines of the file are as follows: 
```gitignore
# Dependency directories 
vendor/

# Development specific environment variables
.env

```
Uncommenting `vendor/` will make git track changes made to the vendor directory so that third party dependencies 
will be committed to git.  This effectively caches the dependencies in your repository. This could come in handy should
a library become unavailable for whatever reason.

Add a hidden file called `.env` which we will use in conjunction with [godotenv](github.com/joho/godotenv) library to 
specify development specific environment variables. This is a nifty mechanism to "use" environment variables in 
development without actually using them.  Down the line, this will allow our code to be seamlessly deployed to a cloud
service such as Heroku or Google Cloud Run.  These platforms usually require the server to listen on a predetermined
port specified by the cloud platform.  This is usually done through use of environment variables.
`PORT` is the only environment variable applicable to this project. Add the following content to the file:
```
PORT=8080
```

Lastly create a directory called `static` inside our project directory to house all our static HTML, CSS, favicon and 
Javascript.  We will add content to this directory later, for now just leave it empty.

## Adding code

Every Go program has a main file, a main package and a main function.  Create a `main.go` file inside your project 
directory as follows:

```{.go .numberLines}
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load env file : %s", err)
	}

	router := httprouter.New()
	router.NotFound = http.FileServer(http.Dir("static"))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	serviceAddress := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:              serviceAddress,
		Handler:           router,
	}
	log.Fatal(srv.ListenAndServe())
}
```

