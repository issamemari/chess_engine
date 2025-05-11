package search

import (
	board "jesus_chess/domain/board"
)

type MoveFinder interface {
	FindBestMove(board board.ChessBoard) (*board.Move, error)
}
