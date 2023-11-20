package engine

import (
	"CS361_Service/internal/web"
	"github.com/notnil/chess/uci"
	"strconv"
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

func ReadyEngine(name string) (e *uci.Engine, err error) {
	eng, err := uci.New(name)
	if err != nil {
		return nil, err
	}
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return nil, err
	}
	return eng, nil
}

// The uci.Engine is thread safe so there's no need for queuing here.
func runPosition(req web.RequestData, eng *uci.Engine) {
	eng.Run(
		uci.CmdSetOption{Name: "MultiPV", Value: strconv.Itoa(req.MultiPV)},
		uci.CmdUCINewGame,
		CmdDumbPosition{FEN: req.FEN},
		uci.CmdGo{Depth: 20})
	//OK wise guy, now how do you read what happened when this: https://github.com/notnil/chess/issues/99
	//That's a TODO for tomorrow...
}
