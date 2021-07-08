package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := routes()
	fmt.Println("Serv is started at http://127.0.0.1:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
