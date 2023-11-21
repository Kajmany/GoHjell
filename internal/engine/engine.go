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

// The uci.Engine is thread safe so there's no need for queuing here.
func runPosition(req web.RequestData, eng *uci.Engine) ([]PrincipalVariation, error) {
	var cmd = CmdDumbMPVGo{depth: 20}
	err := eng.Run(
		uci.CmdSetOption{Name: "MultiPV", Value: strconv.Itoa(req.MultiPV)},
		uci.CmdUCINewGame,
		CmdDumbPosition{FEN: req.FEN},
		cmd)
	if err != nil {
		return nil, err
	}
	return cmd.variations, nil
}
