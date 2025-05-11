package search

import (
	"fmt"
	"math/rand"

	board "jesus_chess/domain/board"
	logging "jesus_chess/domain/logging"
)

type RandomMoveFinder struct {
	logger *logging.Logger
}

func (rmf *RandomMoveFinder) FindBestMove(chessBoard board.ChessBoard) (*board.Move, error) {
	legalMoves := chessBoard.GenerateLegalMoves()
	if len(legalMoves) == 0 {
		return nil, fmt.Errorf("no legal moves available")
	}

	// log all legal moves
	for _, move := range legalMoves {
		rmf.logger.Debug(fmt.Sprintf("legal move: %s from file %d, rank %d to file %d, rank %d", move.Piece.Name, move.From.File, move.From.Rank, move.To.File, move.To.Rank))
	}

	randomIndex := rand.Intn(len(legalMoves))

	move := legalMoves[randomIndex]
	rmf.logger.Debug(fmt.Sprintf("random move selected: %s from file %d, rank %d to file %d, rank %d", move.Piece.Name, move.From.File, move.From.Rank, move.To.File, move.To.Rank))

	return &move, nil
}

func NewRandomMoveFinder(logger *logging.Logger) *RandomMoveFinder {
	return &RandomMoveFinder{logger: logger}
}
