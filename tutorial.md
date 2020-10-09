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

It's worth noting that the entire game could be implemented in the browser using plain JavaScript since (in this case) 
the backend isn't making any third-party API calls, accessing databases or interacting with other microservices i.e. 
things that backend services usually do, but for the sake of exploring the Go language let's pretend that we're building
the War Operation Plan Response system featured in the movie, WarGames.

