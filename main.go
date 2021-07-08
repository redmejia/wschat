package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wschat/handlers"
)

func main() {
	mux := routes()
	log.Println("stating channel listner")
	go handlers.ListenToWsChannel()
	fmt.Println("Serv is started at http://127.0.0.1:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
