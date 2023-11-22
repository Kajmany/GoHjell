package web

import (
	"CS361_Service/internal/common"
	"encoding/json"
	"github.com/notnil/chess/uci"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type EngineInterface interface {
	RunPosition(req common.RequestData) error
	ProxyResults() uci.SearchResults //TODO learn hwo to use interfaces better than... this...
}

func ReadyServer(engine EngineInterface) (mux *http.ServeMux) {
	handler := &EngineHandler{engine: engine}
	mux = http.NewServeMux()
	mux.Handle("/analyze/", handler)
	return mux
}

type ResponseBody struct {
	PrincipleVariations []VariationBody `json:"variations"`
}

type VariationBody struct {
	Depth int      `json:"depth"`
	Score int      `json:"score"`
	Rank  int      `json:"rank"`
	Moves []string `json:"moves"`
}

type EngineHandler struct {
	engine EngineInterface
}

func (handler *EngineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	var rd common.RequestData
	err := decoder.Decode(&rd)
	if err != nil {
		log.Println(err)
		//TODO expand on this, but for now we'll just victim-blame and assume it's the user, haha
		http.Error(w, "There was a problem parsing the request.", http.StatusBadRequest)
		return
	}
	//End code substantially lifted from Alex Edwards. JSON theoretically valid but the actual values aren't checked!
	if !ValidFEN(rd.FEN) { //TODO I realized I can use the notnil chess library to do this better (& slower) than I do
		http.Error(w, "The FEN provided appears to be invalid.", http.StatusBadRequest)
		return
	}
	if rd.MultiPV > 10 || rd.MultiPV < 1 {
		http.Error(w, "The MultiPV value is outside acceptable bounds of 1-10", http.StatusBadRequest)
		return
	}
	if rd.Depth > 20 || rd.Depth < 1 {
		http.Error(w, "The depth value is outside acceptable bounds of 1-20", http.StatusBadRequest)
		return
	}
	log.Println("Successfully parsed input...")
	err = handler.engine.RunPosition(rd)
	processed := handler.engine.ProxyResults()

	body := ResponseBody{make([]VariationBody, rd.MultiPV)}
	for i := 0; i < rd.MultiPV; i++ {
		body.PrincipleVariations[i] = VariationBody{
			Depth: processed.Info.Depth,
			Score: processed.Info.PVs[i].Score,
			Rank:  processed.Info.PVs[i].Rank,
			Moves: processed.Info.PVs[i].Moves,
		}
	}
	JSONData, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "There was a problem packaging your processed data.", http.StatusInternalServerError)
		log.Println("Failed to marshal properly!")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	status, err := w.Write(JSONData)
	{
		log.Println("Finished writing response with status:", status)
		if err != nil {
			log.Println("There was a problem with the response:", err)
		}
	}
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
