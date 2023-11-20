package engine

import (
	"CS361_Service/internal/web"
	"github.com/notnil/chess/uci"
	"strconv"
)

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

func runPositions(queue chan web.RequestData, eng *uci.Engine) {
	for req := range queue {
		//I can't just feed my FEN string to a

		eng.Run(
			uci.CmdSetOption{Name: "MultiPV", Value: strconv.Itoa(req.MultiPV)},
			uci.CmdUCINewGame,
			uci.CmdPosition{Position: req.FEN})

	}
}
