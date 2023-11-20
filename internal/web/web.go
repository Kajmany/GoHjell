package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type RequestData struct {
	//Forsyth-Edwards Notation representation of chess board & game state
	FEN string
	//Engine option to calculate and output multiple Principle Variations
	//Instead of just the best (line of) moves
	MultiPV int
}

func ReadyServer() (mux *http.ServeMux, err error) {
	mux = http.NewServeMux()
	mux.HandleFunc("/analyze/", engineGateway)
	return mux, nil
}

func engineGateway(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		//Nothing to GET, DELETE, PATCH, etc...
		http.Error(w, "This endpoint only accepts POST.", http.StatusMethodNotAllowed)
		return
	}
	// Thanks to https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body ret. 11/19/23
	// For this snippet which concatenates the body to 1MB - should be way larger than what I need, even
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	//Also taken from the blog post, though I've seen the pattern elsewhere. Problem: Decoders aren't thread safe
	//So we need to make a new decoder with every single request! $$$
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var rd RequestData
	err := decoder.Decode(&rd)
	if err != nil {
		log.Println(err)
		//TODO expand on this, but for now we'll just victim-blame and assume it's the user, haha
		http.Error(w, "There was a problem parsing the request.", http.StatusBadRequest)
		return
	}
	//End code substantially lifted from Alex Edwards. JSON theoretically valid but the actual values aren't checked!
	if !ValidFEN(rd.FEN) {
		http.Error(w, "The FEN provided appears to be invalid.", http.StatusBadRequest)
		return
	}
	if rd.MultiPV > 20 || rd.MultiPV < 1 {
		http.Error(w, "The MultiPV value is outside acceptable bounds of 1-20", http.StatusBadRequest)
		return
	}
	log.Println("Successfully parsed input...")
	w.WriteHeader(http.StatusAccepted)
	fmt.Print("Very nice!")
}

func ValidFEN(FEN string) bool {
	//Spaghetti heuristics on the FEN string to validate a subset of the BNF syntax given here:
	//https://www.chessprogramming.org/Forsyth-Edwards_Notation & no attention is given to semantics
	//TODO this could be a lot better (broken out as a real parser, cover more, etc) but this works now
	baseArray := strings.Fields(FEN)
	if len(baseArray) != 6 {
		return false
	}
	//This little monstrosity I whipped up in a PCRE fiddle should get us most the way there on the locations bit
	//Should match exactly 8 times on valid FEN: once per row broken by a /
	re := regexp.MustCompile(`(?i)(8|[1-7]?[pnbrqk][1-7]?)*`)
	matches := re.FindAllString(baseArray[0], -1)
	if len(matches) != 8 {
		return false
	}
	//TODO I could and maybe should check every field (the stakes of a bad FEN is silent error @ the engine level)
	return true
}
