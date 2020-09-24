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
		panic(err)
	}

	router := httprouter.New()
	router.PUT("/game-state", game.SetGameState)
	router.Handler(http.MethodGet, "/debug/pprof/*item", http.DefaultServeMux)

	static := httprouter.New()
	static.ServeFiles ("/*filepath", http.Dir("static"))
	router.NotFound =  static

	//ServiceAddress address to listen on
	bindAddress := os.Getenv("BIND_ADDRESS")
	port        := os.Getenv("PORT")
	serviceAddress := fmt.Sprintf("%s:%s", bindAddress, port)
	srv := &http.Server{
		Addr:              serviceAddress,
		Handler:           router,
	}
	log.Fatal(srv.ListenAndServe())
}
