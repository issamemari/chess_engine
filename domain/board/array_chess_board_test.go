package board

import (
	"testing"
)

func TestNewArrayChessBoard(t *testing.T) {
	cb := NewArrayChessBoard()

	// Validate side to move
	if cb.sideToMove != White {
		t.Errorf("Expected sideToMove to be White, got %v", cb.sideToMove)
	}

	// Validate castling rights
	expectedCastlingRights := CastlingRights{
		WhiteKingSide:  true,
		WhiteQueenSide: true,
		BlackKingSide:  true,
		BlackQueenSide: true,
	}
	if cb.castlingRights != expectedCastlingRights {
		t.Errorf("Expected castling rights to be %v, got %v", expectedCastlingRights, cb.castlingRights)
	}

	// Validate king squares
	if cb.kingSquares[White] != (Square{Rank: 0, File: 4}) {
		t.Errorf("Expected White king square to be {0, 4}, got %v", cb.kingSquares[White])
	}
	if cb.kingSquares[Black] != (Square{Rank: 7, File: 4}) {
		t.Errorf("Expected Black king square to be {7, 4}, got %v", cb.kingSquares[Black])
	}

	// Validate board initialization
	for i := 0; i < BoardWidth; i++ {
		if cb.board[1][i] == nil || cb.board[1][i].Name != Pawn || cb.board[1][i].Color != White {
			t.Errorf("Expected White pawn at (1, %d), got %v", i, cb.board[1][i])
		}
		if cb.board[BoardHeight-2][i] == nil || cb.board[BoardHeight-2][i].Name != Pawn || cb.board[BoardHeight-2][i].Color != Black {
			t.Errorf("Expected Black pawn at (%d, %d), got %v", BoardHeight-2, i, cb.board[BoardHeight-2][i])
		}
	}

	pieceNames := []PieceName{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, pieceName := range pieceNames {
		if cb.board[0][i] == nil || cb.board[0][i].Name != pieceName || cb.board[0][i].Color != White {
			t.Errorf("Expected White %v at (0, %d), got %v", pieceName, i, cb.board[0][i])
		}
		if cb.board[BoardHeight-1][i] == nil || cb.board[BoardHeight-1][i].Name != pieceName || cb.board[BoardHeight-1][i].Color != Black {
			t.Errorf("Expected Black %v at (%d, %d), got %v", pieceName, BoardHeight-1, i, cb.board[BoardHeight-1][i])
		}
	}

	expectedAttackedSquares := map[Color][]Square{
		White: {
			{Rank: 2, File: 0},
			{Rank: 2, File: 1},
			{Rank: 2, File: 2},
			{Rank: 2, File: 3},
			{Rank: 2, File: 4},
			{Rank: 2, File: 5},
			{Rank: 2, File: 6},
			{Rank: 2, File: 7},
		},
		Black: {
			{Rank: 5, File: 0},
			{Rank: 5, File: 1},
			{Rank: 5, File: 2},
			{Rank: 5, File: 3},
			{Rank: 5, File: 4},
			{Rank: 5, File: 5},
			{Rank: 5, File: 6},
			{Rank: 5, File: 7},
		},
	}

	for color, expectedSquares := range expectedAttackedSquares {
		actualSquares := cb.attackedSquares[color]
		expectedSet := make(map[Square]bool)
		for _, square := range expectedSquares {
			expectedSet[square] = true
		}

		for _, square := range actualSquares {
			if !expectedSet[square] {
				t.Errorf("Unexpected attacked square %v for color %v", square, color)
			}
		}

		for _, square := range expectedSquares {
			found := false
			for _, actualSquare := range actualSquares {
				if actualSquare == square {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected attacked square %v for color %v not found", square, color)
			}
		}
	}
}
