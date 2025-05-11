package uci

import (
	"fmt"
	board "jesus_chess/domain/board"
	logging "jesus_chess/domain/logging"
	"os"
	"strings"
)

type UCIHandler struct {
	logger *logging.Logger
	board  board.ChessBoard
}

func NewUCIHandler(logger *logging.Logger, board board.ChessBoard) *UCIHandler {
	return &UCIHandler{
		logger: logger,
		board:  board,
	}
}

func (h *UCIHandler) Handle(command string) {
	h.logger.Debug("handling command: " + command)

	tokens := strings.Fields(command)
	if len(tokens) == 0 {
		return
	}

	switch tokens[0] {
	case "uci":
		h.respond("id name JesusChess")
		h.respond("id author Issa Memari")
		h.respond("uciok")

	case "isready":
		h.respond("readyok")

	case "ucinewgame":
		h.board.SetPosition(board.StartingPosition)

	case "position":
		fen, err := parsePositionCommand(tokens)
		if err != nil {
			h.logger.Error("failed to parse position command: " + err.Error())
			return
		}
		err = h.board.SetPosition(fen)
		if err != nil {
			h.logger.Error("failed to set position: " + err.Error())
			return
		}

	case "go":
		h.respond("info depth 1 multipv 1 score cp -27 pv e7e5")
		h.respond("bestmove e7e5 ponder g1f3")

	case "quit":
		h.logger.Debug("quitting")
		os.Exit(0)

	case "stop":
		h.logger.Debug("stop command received")
		h.respond("stop")

	default:
		h.logger.Error("unknown command: " + command)
		os.Exit(1)
	}
}

func (h *UCIHandler) respond(s string) {
	fmt.Println(s)
	h.logger.Debug("engine responded: " + s)
}

func parsePositionCommand(tokens []string) (string, error) {
	for i, token := range tokens {
		if token == "fen" && i+1 < len(tokens) {
			return strings.Join(tokens[i+1:], " "), nil
		}
	}
	return "", fmt.Errorf("invalid position command, expecting fen string")
}
