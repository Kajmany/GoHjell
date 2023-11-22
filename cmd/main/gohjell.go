package main

import (
	"CS361_Service/internal/engine"
	"CS361_Service/internal/web"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Arguments struct {
	ListenAddress string //TODO This should be tested but we fail fast-ish anyway
	EngineName    string //Won't bother testing before exec
}

func readyArgs() Arguments {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Print("USAGE:\n    gohjell [ListenAddress] [EngineName]\n    Both parameters required.")
		os.Exit(1)
	}
	return Arguments{
		ListenAddress: args[0],
		EngineName:    args[1],
	}
}

func main() {
	args := readyArgs()
	eng, err := engine.ReadyEngine(args.EngineName)
	if err != nil {
		log.Fatal("Error while attempting to start and ready chess engine:", err)
	}
	log.Println("Started engine", eng.ID()["name"])
	engImpl := &engine.EngImplementer{Eng: eng}
	mux := web.ReadyServer(engImpl)
	log.Println("Routes prepared.")
	log.Fatal(http.ListenAndServe(args.ListenAddress, mux))
}
