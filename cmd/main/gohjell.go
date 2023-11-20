package main

import (
	"CS361_Service/internal/web"
	"log"
	"net/http"
)

func main() {
	/*	eng, err := engine.ReadyEngine("stockfish")
		if err != nil {
			log.Fatal("Error while attempting to start and ready Chess engine: ", err)
		}*/
	mux, _ := web.ReadyServer()
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
