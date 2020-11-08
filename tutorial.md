Tic Tac Toe

In this tutorial you will learn how to:

1. Build a web based Tic Tac Toe game in Go.
1. Write a web based client in vanilla JavaScript.
1. Implement the minimax algorithm to recursively calculate the most optimal next move for the computer.
1. Write unit and integration tests to exercise and prove correctness of the logic.

## Pre-requisites

The following pre-requisites should be in place to successfully follow this guide.

1. Go should be installed on your computer.
1. You should have a basic understanding of Go, HTML, CSS and JavaScript.

## Introduction

Tic Tac Toe otherwise known as "Naughts and Crosses" is a really basic board game familiar to everyone the World over. 
When I was a young kiddie, Tic Tac Toe captured my imagination in the 1983 film: [WarGames](https://youtu.be/F7qOV8xonfY).
At the time I was too young to understand what was meant by "The only winning move is not to play" because as a clueless
child playing against a computer I often messed-up and lost. It became clear to me as I grew that two players, playing 
flawlessly, the best result one can hope for is a draw.

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
features are code coverage reporting, race condition detection, benchmark testing and profiling. Beyond writing basic
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

## Adding the main function

Every Go program has a main file, a main package and a main function.  Create a `main.go` file inside your project 
directory as follows:

```go
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
    // Paragraph #1
    err := godotenv.Load()
    if err != nil {
        log.Printf("not using .env file")
    }

    // Paragraph #2 
    router := httprouter.New()
    router.NotFound = http.FileServer(http.Dir("static"))

    // Paragraph #3
    port := os.Getenv("PORT")
    if len(port) == 0 {
        port = "8080"
    }

    // Paragraph #4
    serviceAddress := fmt.Sprintf(":%s", port)
    srv := &http.Server{
        Addr:              serviceAddress,
        Handler:           router,
    }
    log.Fatal(srv.ListenAndServe())
}
```

To update your vendor dependencies, run the following command:
```
go mod vendor
```

You should see some similar to:
```
go: finding module for package github.com/joho/godotenv
go: finding module for package github.com/julienschmidt/httprouter
go: found github.com/joho/godotenv in github.com/joho/godotenv v1.3.0
go: found github.com/julienschmidt/httprouter in github.com/julienschmidt/httprouter v1.3.0
```

The breakdown is as follows:
### Paragraph #1 Loading development ENV values
As mentioned above, we use github.com/joho/godotenv to load environment variables specific to your development 
environment from the `.env` file located in your project directory.  The `.env` file should not be committed to git or at
least omitted from the Docker image that gets built for production. Thus `godotenv.Load()` will actually return an 
error in production, which can just be ignored.

### Paragraph #2 Creating the MUX router   
We create a MUX router (aka HTTP multiplexer) using the `github.com/julienschmidt/httprouter` library.  This package 
maps routes to HTTP handler functions. We're specifically using this third-party library over and above Go's 
standard MUX router to allow us to serve the HTML, CSS, favicon and JavaScript from the static directory.  Gorilla Mux
is another popular option among the Go community. For this project I opted to expose myself to this router library as
[benchmark tests have proven it to be  fast and efficient with memory allocations](https://github.com/julienschmidt/go-http-routing-benchmark). 

### Paragraph #3 Read PORT environment variable
In paragraph #3 we call on the `os` package to read the value of `PORT` environment variable.  If no such value exists
then we default to port 8080. Cloud services such as Google Cloud Run, Heroku, etc. usually stipulate which port the 
service needs to listen on via environment variables. In development, we obtain this value from the `.env` file.

### Paragraph #4 Listen and Serve
Paragraph #4 creates and configures a pointer to a `http.Server` such that the server binds to all interfaces on the 
specified port. `http.ListenAndServe` is a blocking call which returns upon error or upon process termination. 

## Adding the game logic

In this section we are going to:
- Create a new package called game and add a new file called game.go.
- Write the HTTP handler to process requests sent from the browser.
- Write the logic that will govern the game.
- Write tests to prove the logic. 

### Create the `game` package
Create a new directory underneath your project directory called `game` and add a new file called `game.go` to that 
directory.

Add the following code to the `game.go` file.  This declares the data structures that will use in the game, to maintain 
state and communicate with the browser over HTTP. 

```go
package game

type SquareState rune
const (
	SquareStateEmpty  SquareState = 0
	SquareStateCross  SquareState = 'X'
	SquareStateNaught SquareState = '0'
)

type Result int
const (
	ResultNone      Result = iota
	ResultNInARow   Result = iota
	ResultStalemate Result = iota
)

type TicTacToeState struct {
	Board [][]SquareState `json:"board"`
	Turn  int             `json:"-"`
}

type TicTacToeStateResponse struct {
	Board      [][]SquareState   `json:"board"`
	Result     Result            `json:"result,omitempty"`
	WinningRow [][]SquareState   `json:"winningRow,omitempty"`
	Turn       int               `json:"turn"`
	NextPlayer rune              `json:"nextPlayer"`
}
```

### `SquareState`
We declare `SquareState` as an alias to `rune` which we will use to declare three constants to denote the state of a 
square on the game board.  Each square can either be empty or contain a `0` or an `X`. 

### `Result`
`Result` is another type alias of int.  We will use this type to denote the result of the game.  Our implementation of 
Tic-tac-toe has three states the game can be in: 
1. `ResultNone` There are still empty squares on the board, and the game has not concluded with any player making a row 
of three.   
1. `ResultNInARow` A player has managed to best its opponent by placing 3 pieces next to each other to make a
row of three and win the game.
1. `ResultStalemate` There are no empty squares left on the board with neither player succeeding to make a row of 
three.

In each case [`iota`](https://golang.org/ref/spec#Iota) is used to assign a successive integer value to the constant,
making each constant unique with respect to the others declared in the same group.

### `TicTacToeState`
`TicTacToeState` models the state of the game and defines the structure of the request to be received from the 
browser.  It contains a 2D slice of `SquareStates` to represent each square of the board; which is essentially the 
state.  The struct also contains an integer called `turn` which indicates whose turn is next.  This field is calculated 
by counting the non-empty squares of the board.  By caching this counter in the struct, it saves us from having to
recalculate it each time we need it. 

### `TicTacToeStateResponse`
`TicTacToeStateResponse` specifies the structure of the response sent back to the browser after the player's move has 
been processed by the game's logic.  The server sends the new state of the board back to the client in the `Board` field.
The `Result` field indicates the state of the game.  `omit empty` will omit the `result` json field from the response 
if the game has reached a conclusion; i.e. the result is `ResultNone` since in GoLang zero values are synonymous to being
empty. `Turn` and `NextPlayer` are useful meta-data fields to augment the response with additional information about the
state of the game.

### HTTP Handler

The MUX server routes all `PUT` requests made to `/game-state` through to `TicTacToeStateHandler`.  This is the only 
request handler this game needs. `PUT` method is a good option to choose for this request since the handler is idempotent,
meaning the same request can be repeatedly submitted to the server with no changes to state or side effects. The   
state of our game is the board, and the server does not maintain state in session variables, databases, etc.   

```go
// TicTacToeStateHandler accepts a TicTacToeState representing the
// current state of the game and responds with a TicTacToeStateResponse
// describing the new state of the game.
func TicTacToeStateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // Parapgraph #1
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "could read request", err)
		return
	}

    // Parapgraph #2
	req := &TicTacToeState{
		Board: makeBoard(3),
	}
	err = json.Unmarshal(b, req)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "could not interpret request", err)
		return
	}

	req.calculateStateFromRequest()
   
    // Parapgraph #3
	result, _ := req.getGameResult()
	if result == ResultNone {
		_, x, y := computeMove(*req, true)
		err = req.occupyPosition(x, y)
		if err != nil {
			writeHTTPError(w, http.StatusInternalServerError, "failed to set board", err)
			return
		}
	}

    // Parapgraph #4
	result, winningRow := req.getGameResult()
	resp := TicTacToeStateResponse{
		Board:      req.Board,
		Result:     result,
		Turn:       req.Turn,
		WinningRow: winningRow,
		NextPlayer: req.playersTurn(),
	}

    // Parapgraph #5
	b, err = json.Marshal(resp)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "failed to marshal response", err)
		return
	}

    // Parapgraph #6
	_, err = w.Write(b)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func writeHTTPError(w http.ResponseWriter, statusCode int, description string, err error) {
	message := fmt.Sprintf("%s : %v", description, err)
	_, wErr := w.Write([]byte(message))
	if wErr != nil {
		log.Fatal(wErr)
	}
	log.Fatal(description, err)
	w.WriteHeader(statusCode)
}
```

`Paragraphs #1 and #2` concern unmarshalling the request from JSON to `TicTacToeState` which represents the state of the
board.

`Parapgraph #3` calls into the heart of the game's logic. If the game is still in play, we compute the most effective 
move in response to the current state of the game.

`Parapgraph #4` collates the response in preparation to be sent back to the browser.

`Parapgraph #5 and #6` encodes the response into JSON which is written to the response writer and sent back to the 
browser.

In the event that an error occurs, `writeHTTPError` is called to write the status code and error description 
to the response writer. Note that `http.StatusBadRequest` is used in cases where the error centered around the request,
whereas `http.StatusInternalServerError` is returned where the internal logic of the system encountered a failure. 

### Game Logic

Before turning our focus directly onto the logic that makes the game tick, we have two helper/utility functions 
`makeBoard` and `copyBoard` which provide the functionality needed to initialize empty boards and copy existing boards 
respectively. These are needed since in this design, we have used a 2D slice to represent the state of the board. 
Alternatively we could've used a single `n x n` slice and used modular and division arithmetic to calculate board 
positions. For instance the middle right `(2,1)` square of the board would be `board[5]`; `5` could be mapped to co-ordinates as
follows (keep in mind that slices and arrays are zero-based):
```go
col := 5 % 3 // == 2
row := 5 / 3 // == 1
```  
Each approach trades-off complexity in different areas.
 
```go
func makeBoard(n int) [][]SquareState {
	board := make([][]SquareState, n)
	for i := 0; i < n; i++ {
		board[i] = make([]SquareState, n)
	}

	return board
}

func copyBoard(src [][]SquareState) [][]SquareState {
	n := len(src)
	board := make([][]SquareState, n)
	for j := 0; j < n; j++ {
		board[j] = make([]SquareState, n)
		copy(board[j], src[j])
	}

	return board
}
```

`TicTacToeState` structure has several pointer receiver functions to initialize, get and set the state of the game.
 
- `initialize` calculates and caches the `Turn` property of the game's state so that it doesn't need to be calculated it
time its required.
```go
func (t *TicTacToeState) initialize() {
	turn := 1
	for _, y := range t.Board {
		for _, x := range y {
			if x != SquareStateEmpty {
				turn++
			}
		}
	}

	t.Turn = turn
}
```

- `playersTurn` returns the symbol of the current player expected to make the next move.
```go
func (t *TicTacToeState) playersTurn() rune {
	if t.Turn % 2 == 1 {
		return 'X'
	}

	return '0'
}
```

- `isOccupied` returns `true`/`false` whether the specified co-ordinates are occupied. 
```go
func (t *TicTacToeState) isOccupied(x, y int) bool {
	return t.Board[y][x] != SquareStateEmpty
}
```

- `occupyPosition` marks the specified co-ordinates occupied according to whose turn it is.
```go
func (t *TicTacToeState) occupyPosition(x, y int) error {
	if x < 0 || x > len(t.Board) || y < 0 || y > len(t.Board) {
		return errors.New("invalid coordinate")
	}
	if t.Board[y][x] != SquareStateEmpty {
		return errors.New("already occupied")
	}

	player := t.playersTurn()
	t.Turn++
	if player == 'X' {
		t.Board[y][x] = SquareStateCross
		return nil
	}
	t.Board[y][x] = SquareStateNaught

	return nil
}
```

- `getGameResult`. Ideally functions should have one purpose or do one thing.  This function kind of breaks this rule,
in that it calculates the result of the game and returns the row that concluded the game (if there is one). In this case
the two concepts are related enough to combine them. Under the hood this function probes each line (diagonal, row and 
column) to find a row of n adjacent pieces. It does this by iterating over each square of that line and compares 
whether the adjacent pieces are equal. The comparison loop ends if two positions differ as there is no point checking 
any further. If the loop managed to increment `i` enough times to equate to `n - 1` then a concluding line has been
found, and the player that made that line has won the game.
```go
// getGameResult calculates the current state of the game returning the result
//  and the row that concluded the game if there is a complete row, nil otherwise.
func (t *TicTacToeState) getGameResult() (Result, [][]SquareState) {
	n := len(t.Board)
	var rowOfN [][]SquareState

	// Check diagonal
	rowOfN = makeBoard(n)
	var i int
	for i = 0; i < n - 1 && t.Board[i][i] == t.Board[i+1][i+1] && t.Board[i][i] != SquareStateEmpty; i++ {
		rowOfN[i][i] = t.Board[i][i]
		rowOfN[i+1][i+1] = t.Board[i+1][i+1]
	}
	if i == n - 1 {
		return ResultNInARow, rowOfN
	}

	// Check anti-diagonal
	rowOfN = makeBoard(n)
	for i = 0; i < n - 1 && t.Board[n-1-i][i] == t.Board[n-i-2][i+1] && t.Board[n-i-2][i+1] != SquareStateEmpty; i++ {
		rowOfN[n-1-i][i] = t.Board[n-1-i][i]
		rowOfN[n-i-2][i+1] = t.Board[n-i-2][i+1]
	}
	if i == n - 1 {
		return ResultNInARow, rowOfN
	}

	var j int
	for j = 0; j < n; j++ {
		// Check columns
		rowOfN = makeBoard(n)
		for i = 0; i < n-1 && t.Board[i][j] == t.Board[i+1][j] && t.Board[i][j] != SquareStateEmpty; i++ {
			rowOfN[i][j] = t.Board[i][j]
			rowOfN[i+1][j] = t.Board[i+1][j]
		}
		if i == n - 1 {
			return ResultNInARow, rowOfN
		}

		// Check rows
		rowOfN = makeBoard(n)
		for i = 0; i < n-1 && t.Board[j][i] == t.Board[j][i+1] && t.Board[j][i] != SquareStateEmpty; i++ {
			rowOfN[j][i] = t.Board[j][i]
			rowOfN[j][i+1] = t.Board[j][i+1]
		}
		if i == n - 1 {
			return ResultNInARow, rowOfN
		}
	}

	// Check for stalemate
	if t.Turn > n * n {
		return ResultStalemate, nil
	}

	return ResultNone, nil
}
```