package engine

import (
	"bufio"
	"errors"
	"github.com/notnil/chess/uci"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// notnil's position command would require me to start a whole chess game with the parent library OR do a lot of manual
// parsing in order to get chess.Position, but I really just want to pass my FEN to the engine!
type CmdDumbPosition struct {
	//Forsyth-Edwards Notation representation of chess board & game state
	FEN string
}

// This just cuts out the extra logic in uci.CmdPosition and echoes the string as a UCI command for the uci.Engine
func (cmd CmdDumbPosition) String() string {
	return "position fen " + cmd.FEN
}

// Required to satisfy interface. As with uci.CmdPosition, there's no possible UCI error state because the engine will
// try to ignore a bad position. !!Stockfish won't try hard and will crash often!!
func (cmd CmdDumbPosition) ProcessResponse(e *uci.Engine) error {
	return nil
}

// Subset of data reported by UCI info output. All I care to send for JSON serialization out!
type PrincipalVariation struct {
	//depth the engine was asked to process for this PV. NOT a way to measure actual depth (seldepth) or # of moves inside
	Depth int
	//score is measured in nonsense unit known as the centipawn - chess computers love it!
	Score int
	//when MPV is enabled engine numbers the PV's from 1-N where 1 is the best quality of them, and N is the worst.
	Rank int
	//stringly typed array of moves, alternating between 1st player @ [0] and 1st opponent ponder @ [1].
	Moves []string
}

// notnil's uci.SearchResults populated by the uci.CmdGo only tracks one Principal Variation. MPV was requested of me
// so I need a command that will populate an array of PrincipalVariation!
type CmdDumbMPVGo struct {
	depth      int
	variations []PrincipalVariation
}

// It's so easy to write compact commands when you only allow a single workflow 8D
func (cmd CmdDumbMPVGo) String() string {

	return "go depth " + strconv.Itoa(cmd.depth)
}

// Most of this function is taken from the uci.CmdGo version of this, Except for the reflection bit and the regex/match
// Logic used to extract relevant bits into a struct array.
func (cmd CmdDumbMPVGo) ProcessResponse(e *uci.Engine) error {

	//First we have to do some clownery to snatch the engine's IO out. This can cause a ruckus if the lib changes.
	outVal := reflect.ValueOf(e).Elem().FieldByName("out")
	outProper := outVal.Interface().(io.Reader)
	scanner := bufio.NewScanner(outProper)
	pvContainer := make([]PrincipalVariation, cmd.depth)
	cmd.variations = pvContainer
	//This WILL unnecessarily recompile the Regex with each call. Could potentially be broken out w/ sync.Once
	regex := regexp.MustCompile("multipv (\\d+) score cp (-?\\d+) .+ time (\\d+) pv (.+)")
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "info depth "+strconv.Itoa(cmd.depth)) {
			//Instead of only taking the best move at the end, we examine all messages that have reached our depth
			//BUT the system may keep going for a bit as it explores seldepth, so we must be prepared to take the most
			//recent message while parsing. EXAMPLE:
			//info depth 20 seldepth 26 multipv 10 score cp -4 nodes 4891449 nps 1536741 hashfull 972 tbhits 0 time 3183
			//pv a2a3 e7e5 c2c4 g8f6 d2d3 d7d5 c4d5 f6d5 g1f3 b8c6 e2e4 d5b6 f1e2 c8g4 c1e3 f7f5 e4f5 d8d7 e1g1 e8c8
			//f1e1 g4f5 b1c3
			//Now to get bits with capture groups.
			matches := regex.FindStringSubmatch(text)
			//whole string(?), then captures
			if len(matches) != 4 {
				return errors.New("regex didn't match right: " + text) //TODO not great way to do errors
			}
			newPVRank, err := strconv.Atoi(matches[1])
			if err != nil {
				return err //TODO ditto
			}
			newPVScore, err := strconv.Atoi(matches[2])
			if err != nil {
				return err
			}
			newPVMoves := strings.Split(matches[4], " ")
			var newPV = PrincipalVariation{Rank: newPVRank, Score: newPVScore, Moves: newPVMoves}
			//This keeps the array ordered by MultiPV rank (which starts at 1). We have to be prepared to over-write
			//Because it's possible for INFO at our desired depth for a given MultiPV to be output several times.
			//This ensures only the last one is saved in our final array.
			pvContainer[newPVRank-1] = newPV //FIXME it's possible to return without all slots being populated
		}
	}
	//We're totally reliant on side effects because the interface is being satisfied, and does the same
	return nil
}
