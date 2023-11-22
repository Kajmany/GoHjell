package common

type RequestData struct {
	//Forsyth-Edwards Notation representation of chess board & game state
	FEN string
	//Engine option to calculate and output multiple Principle Variations
	//Instead of just the best (line of) moves
	MultiPV int
	//Engine option for how many moves ahead
	Depth int
}
