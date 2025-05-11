package uci

import (
	"fmt"
	board "jesus_chess/domain/board"
	logging "jesus_chess/domain/logging"
	search "jesus_chess/domain/search"
	"os"
	"strings"
)

type UCIHandler struct {
	logger     *logging.Logger
	board      board.ChessBoard
	moveFinder search.MoveFinder
}

func NewUCIHandler(logger *logging.Logger, board board.ChessBoard, moveFinder search.MoveFinder) *UCIHandler {
	return &UCIHandler{
		logger:     logger,
		board:      board,
		moveFinder: moveFinder,
	}
}

func (h *UCIHandler) Handle(command string) {
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
		h.board = board.NewArrayChessBoard(h.logger)

	case "position":
		fen, moves, err := parsePositionCommand(tokens)
		if err != nil {
			h.logger.Error("failed to parse position command: " + err.Error())
			return
		}
		err = h.board.SetPosition(fen)
		if err != nil {
			h.logger.Error("failed to set position: " + err.Error())
			return
		}
		for _, move := range moves {
			h.logger.Debug("making move: " + moveToUCI(move))
			h.logger.Debug("side to move: " + (string)(h.board.SideToMove()))
			err := h.board.MakeMove(move)
			h.logger.Debug("move made: " + moveToUCI(move))
			h.logger.Debug("side to move: " + (string)(h.board.SideToMove()))
			if err != nil {
				h.logger.Error("failed to make move: " + err.Error())
				return
			}
		}

	case "go":
		move, err := h.moveFinder.FindBestMove(h.board)
		if err != nil {
			h.logger.Error("failed to find best move: " + err.Error())
			return
		}
		if move == nil {
			h.logger.Error("no best move found")
			return
		}
		moveString := moveToUCI(*move)
		h.logger.Debug("best move found: " + moveString)
		h.respond("info depth 1 multipv 1 score cp -27 pv " + moveString)
		h.respond("bestmove " + moveString)

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

func parsePositionCommand(tokens []string) (string, []board.Move, error) {
	if len(tokens) < 2 {
		return "", nil, fmt.Errorf("expected at least 2 tokens")
	}

	if tokens[1] != "fen" {
		return "", nil, fmt.Errorf("invalid position command, expected fen string")
	}

	if len(tokens) < 8 {
		return "", nil, fmt.Errorf("incomplete fen string")
	}
	fen := strings.Join(tokens[2:8], " ")
	if len(tokens) > 8 && tokens[8] == "moves" {
		moves, err := parseMoves(tokens[9:])
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse moves: %w", err)
		}
		return fen, moves, nil
	}
	return fen, nil, nil
}

func parseMoves(moveTokens []string) ([]board.Move, error) {
	var moves []board.Move
	for _, moveToken := range moveTokens {
		move, err := parseMove(moveToken)
		if err != nil {
			return nil, fmt.Errorf("invalid move %s: %w", moveToken, err)
		}
		moves = append(moves, move)
	}
	return moves, nil
}

func parseMove(moveToken string) (board.Move, error) {
	from_file := int(moveToken[0] - 'a')
	from_rank := int(moveToken[1] - '1')
	to_file := int(moveToken[2] - 'a')
	to_rank := int(moveToken[3] - '1')

	from, err := board.NewSquare(from_rank, from_file)
	if err != nil {
		return board.Move{}, err
	}
	to, err := board.NewSquare(to_rank, to_file)
	if err != nil {
		return board.Move{}, err
	}

	return board.Move{
		From: from,
		To:   to,
	}, nil
}

func moveToUCI(move board.Move) string {
	from_rank := move.From.Rank
	from_file := move.From.File
	to_rank := move.To.Rank
	to_file := move.To.File

	return fmt.Sprintf("%c%c%c%c", 'a'+from_file, '1'+from_rank, 'a'+to_file, '1'+to_rank)
}
