package engine

import (
	"CS361_Service/internal/common"
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

type EngImplementer struct {
	Eng *uci.Engine
}

// The uci.Engine is thread safe so there's no need for queuing here.
func (e EngImplementer) RunPosition(req common.RequestData) error {
	//var cmd = uci.CmdDumbMPVGo{Depth: req.Depth, MultiPV: req.MultiPV} //Kind of redundant info pass FIXME ?
	err := e.Eng.Run(
		uci.CmdSetOption{Name: "MultiPV", Value: strconv.Itoa(req.MultiPV)},
		uci.CmdUCINewGame,
		uci.CmdDumbPosition{FEN: req.FEN},
		uci.CmdGo{Depth: req.Depth, MultiPV: req.MultiPV})
	if err != nil {
		return err
	}
	return nil
}

func (e EngImplementer) ProxyResults() uci.SearchResults {
	return e.Eng.SearchResults()
}
