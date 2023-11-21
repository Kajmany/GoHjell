package main

import (
	"CS361_Service/internal/engine"
	"CS361_Service/internal/web"
	"log"
	"net/http"
)

func main() {
	eng, err := engine.ReadyEngine("stockfish")
	if err != nil {
		log.Fatal("Error while attempting to start and ready Chess engine: ", err)
	}
	mux := web.ReadyServer(eng)
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
