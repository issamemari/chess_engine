package board

import "fmt"

const (
	BoardHeight = 8
	BoardWidth  = 8
)

type Square struct {
	Rank int
	File int
}

func NewSquare(rank, file int) (Square, error) {
	if rank < 0 || rank > BoardHeight-1 || file < 0 || file > BoardWidth-1 {
		return Square{}, fmt.Errorf("invalid square: rank %d, file %d", rank, file)
	}
	return Square{Rank: rank, File: file}, nil
}

type PieceName string

const (
	Pawn   PieceName = "P"
	Knight           = "N"
	Bishop           = "B"
	Rook             = "R"
	Queen            = "Q"
	King             = "K"
)

type Color string

const (
	White Color = "W"
	Black       = "B"
)

type Piece struct {
	Name  PieceName
	Color Color
}

type CastlingRights struct {
	WhiteKingSide  bool
	WhiteQueenSide bool
	BlackKingSide  bool
	BlackQueenSide bool
}

type Move struct {
	Piece       Piece
	From        Square
	To          Square
	Promotion   *Piece
	IsCastling  bool
	IsEnPassant bool
	IsCapture   bool
}

type ChessBoard interface {
	PieceAt(square Square) *Piece
	IsOccupied(square Square) bool
	SideToMove() Color
	CastlingRights() CastlingRights
	GenerateLegalMoves() []Move
	IsMoveLegal(move Move) bool
	InCheck(color Color) bool
	MakeMove(move Move) error
	UndoMove() error
	SetPosition(fen string) error
	Display() string
}
