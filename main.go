package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/donohutcheon/tictactoe/game"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load env file : %s", err)
	}

	router := httprouter.New()
	router.PUT("/game-state", game.TicTacToeStateHandler)
	router.Handler(http.MethodGet, "/debug/pprof/*item", http.DefaultServeMux)

	static := httprouter.New()
	static.ServeFiles ("/*filepath", http.Dir("static"))
	router.NotFound =  static

	port        := os.Getenv("PORT")
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
