package board

import (
	"fmt"
	"strings"

	logging "jesus_chess/domain/logging"
)

type ArrayChessBoard struct {
	board           [BoardHeight][BoardWidth]*Piece
	sideToMove      Color
	moveHistory     []Move
	castlingRights  CastlingRights
	kingSquares     map[Color]Square
	attackedSquares map[Color][]Square
	logger          *logging.Logger
}

func NewArrayChessBoard(logger *logging.Logger) *ArrayChessBoard {
	cb := &ArrayChessBoard{sideToMove: White, logger: logger}

	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			cb.board[rank][file] = nil
		}
	}

	// Place pawns
	for i := 0; i < BoardWidth; i++ {
		cb.board[1][i] = &Piece{Pawn, White}
		cb.board[BoardHeight-2][i] = &Piece{Pawn, Black}
	}

	// Place other pieces
	pieceNames := []PieceName{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, pieceName := range pieceNames {
		cb.board[0][i] = &Piece{pieceName, White}
		cb.board[BoardHeight-1][i] = &Piece{pieceName, Black}
	}

	// Initialize castling rights
	cb.castlingRights = CastlingRights{
		WhiteKingSide:  true,
		WhiteQueenSide: true,
		BlackKingSide:  true,
		BlackQueenSide: true,
	}

	// Initialize king squares
	cb.kingSquares = make(map[Color]Square)
	cb.kingSquares[White] = Square{Rank: 0, File: 4}
	cb.kingSquares[Black] = Square{Rank: 7, File: 4}

	// Initialize attacked squares
	cb.attackedSquares = make(map[Color][]Square)
	cb.attackedSquares[White] = []Square{}
	cb.attackedSquares[Black] = []Square{}
	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			piece := cb.board[rank][file]
			if piece == nil {
				continue
			}
			cb.attackedSquares[piece.Color] = append(cb.attackedSquares[piece.Color], cb.getAttackedSquares(Square{Rank: rank, File: file})...)
		}
	}

	// Initialize move history
	cb.moveHistory = []Move{}

	// Set the initial side to move
	cb.sideToMove = White

	return cb
}

func (cb *ArrayChessBoard) getAttackedSquares(sq Square) []Square {
	piece := cb.board[sq.Rank][sq.File]
	if piece == nil {
		return []Square{}
	}

	attackedSquares := []Square{}

	if piece.Name == Pawn {
		attackedSquares = cb.getPawnAttackedSquares(sq, piece.Color)
	}

	if piece.Name == Knight {
		attackedSquares = cb.getKnightAttackedSquares(sq)
	}

	if piece.Name == Bishop || piece.Name == Queen {
		attackedSquares = cb.getDiagonallyAttackedSquares(sq)
	}

	if piece.Name == Rook || piece.Name == Queen {
		attackedSquares = cb.getHorizontallyAndVerticallyAttackedSquares(sq)
	}

	if piece.Name == King {
		attackedSquares = []Square{
			{Rank: sq.Rank + 1, File: sq.File},
			{Rank: sq.Rank - 1, File: sq.File},
			{Rank: sq.Rank, File: sq.File + 1},
			{Rank: sq.Rank, File: sq.File - 1},
			{Rank: sq.Rank + 1, File: sq.File + 1},
			{Rank: sq.Rank + 1, File: sq.File - 1},
			{Rank: sq.Rank - 1, File: sq.File + 1},
			{Rank: sq.Rank - 1, File: sq.File - 1},
		}
	}

	validAttackedSquares := []Square{}
	for _, attackedSquare := range attackedSquares {
		if cb.validateAttackedSquare(attackedSquare, piece.Color) {
			validAttackedSquares = append(validAttackedSquares, attackedSquare)
		}
	}

	return validAttackedSquares
}

func (cb *ArrayChessBoard) getPawnAttackedSquares(sq Square, color Color) []Square {
	if color == White {
		return []Square{
			{Rank: sq.Rank + 1, File: sq.File - 1},
			{Rank: sq.Rank + 1, File: sq.File + 1},
		}
	} else {
		return []Square{
			{Rank: sq.Rank - 1, File: sq.File - 1},
			{Rank: sq.Rank - 1, File: sq.File + 1},
		}
	}
}

func (cb *ArrayChessBoard) getKnightAttackedSquares(sq Square) []Square {
	return []Square{
		{Rank: sq.Rank + 2, File: sq.File + 1},
		{Rank: sq.Rank + 2, File: sq.File - 1},
		{Rank: sq.Rank - 2, File: sq.File + 1},
		{Rank: sq.Rank - 2, File: sq.File - 1},
		{Rank: sq.Rank + 1, File: sq.File + 2},
		{Rank: sq.Rank + 1, File: sq.File - 2},
		{Rank: sq.Rank - 1, File: sq.File + 2},
		{Rank: sq.Rank - 1, File: sq.File - 2},
	}
}

func (cb *ArrayChessBoard) getDiagonallyAttackedSquares(sq Square) []Square {
	attackedSquares := []Square{}
	currentSquare := sq
	for currentSquare.Rank < BoardHeight && currentSquare.File < BoardWidth {
		currentSquare.Rank++
		currentSquare.File++
		attackedSquares = append(attackedSquares, currentSquare)
		if cb.IsOccupied(currentSquare) {
			break
		}
	}
	currentSquare = sq
	for currentSquare.Rank >= 0 && currentSquare.File >= 0 {
		currentSquare.Rank--
		currentSquare.File--
		attackedSquares = append(attackedSquares, currentSquare)
		if cb.IsOccupied(currentSquare) {
			break
		}
	}
	return attackedSquares
}

func (cb *ArrayChessBoard) getHorizontallyAndVerticallyAttackedSquares(sq Square) []Square {
	attackedSquares := []Square{}
	currentSquare := sq
	for currentSquare.Rank < BoardHeight {
		currentSquare.Rank++
		attackedSquares = append(attackedSquares, currentSquare)
		if cb.IsOccupied(currentSquare) {
			break
		}
	}
	currentSquare = sq
	for currentSquare.Rank >= 0 {
		currentSquare.Rank--
		attackedSquares = append(attackedSquares, currentSquare)
		if cb.IsOccupied(currentSquare) {
			break
		}
	}
	currentSquare = sq
	for currentSquare.File < BoardWidth {
		currentSquare.File++
		attackedSquares = append(attackedSquares, currentSquare)
		if cb.IsOccupied(currentSquare) {
			break
		}
	}
	currentSquare = sq
	for currentSquare.File >= 0 {
		currentSquare.File--
		attackedSquares = append(attackedSquares, currentSquare)
		if cb.IsOccupied(currentSquare) {
			break
		}
	}
	return attackedSquares
}

func (cb *ArrayChessBoard) validateAttackedSquare(sq Square, attackingColor Color) bool {
	if sq.Rank < 0 || sq.Rank >= BoardHeight {
		return false
	}
	if sq.File < 0 || sq.File >= BoardWidth {
		return false
	}
	if cb.board[sq.Rank][sq.File] != nil && cb.board[sq.Rank][sq.File].Color == attackingColor {
		return false
	}
	return true
}

func (cb *ArrayChessBoard) PieceAt(sq Square) *Piece {
	return cb.board[sq.Rank][sq.File]
}

func (cb *ArrayChessBoard) IsOccupied(sq Square) bool {
	if sq.Rank < 0 || sq.Rank >= BoardHeight {
		return false
	}
	if sq.File < 0 || sq.File >= BoardWidth {
		return false
	}
	return cb.board[sq.Rank][sq.File] != nil
}

func (cb *ArrayChessBoard) SideToMove() Color {
	return cb.sideToMove
}

func (cb *ArrayChessBoard) CastlingRights() CastlingRights {
	return cb.castlingRights
}

func (cb *ArrayChessBoard) GenerateLegalMoves() []Move {
	moves := []Move{}

	moves = append(moves, cb.generateLegalPawnMoves()...)
	moves = append(moves, cb.generateLegalKnightMoves()...)
	moves = append(moves, cb.generateLegalBishopMoves()...)
	moves = append(moves, cb.generateLegalRookMoves()...)
	moves = append(moves, cb.generateLegalQueenMoves()...)
	moves = append(moves, cb.generateLegalKingMoves()...)
	return moves
}

func (cb *ArrayChessBoard) generateLegalPawnMoves() []Move {
	moves := []Move{}
	color := cb.sideToMove
	direction := 1
	startRank := 1
	promotionRank := 7
	if color == Black {
		direction = -1
		startRank = 6
		promotionRank = 0
	}
	for file := 0; file < BoardWidth; file++ {
		for rank := 0; rank < BoardHeight; rank++ {
			piece := cb.board[rank][file]
			if piece != nil && piece.Name == Pawn && piece.Color == color {
				from := Square{Rank: rank, File: file}
				forward := Square{Rank: rank + direction, File: file}
				if !cb.IsOccupied(forward) && forward.Rank >= 0 && forward.Rank < BoardHeight {
					if forward.Rank == promotionRank {
						promotions := []PieceName{Queen, Rook, Bishop, Knight}
						for _, promo := range promotions {
							moves = append(moves, Move{From: from, To: forward, Piece: *piece, Promotion: &Piece{Name: promo, Color: color}, PreviousCastlingRights: cb.castlingRights})
						}
					} else {
						moves = append(moves, Move{From: from, To: forward, Piece: *piece, PreviousCastlingRights: cb.castlingRights})
					}
					if rank == startRank {
						twoForward := Square{Rank: rank + 2*direction, File: file}
						if !cb.IsOccupied(twoForward) {
							moves = append(moves, Move{From: from, To: twoForward, Piece: *piece, PreviousCastlingRights: cb.castlingRights})
						}
					}
				}
				attacks := cb.getPawnAttackedSquares(from, color)
				for _, target := range attacks {
					if target.Rank < 0 || target.Rank >= BoardHeight || target.File < 0 || target.File >= BoardWidth {
						continue
					}
					targetPiece := cb.PieceAt(target)
					if targetPiece != nil && targetPiece.Color != color {
						moves = append(moves, Move{From: from, To: target, Piece: *piece, CapturedPiece: targetPiece, PreviousCastlingRights: cb.castlingRights})
					}
				}
			}
		}
	}
	return moves
}

func (cb *ArrayChessBoard) generateLegalKnightMoves() []Move {
	moves := []Move{}
	color := cb.sideToMove
	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			piece := cb.board[rank][file]
			if piece != nil && piece.Name == Knight && piece.Color == color {
				from := Square{Rank: rank, File: file}
				for _, to := range cb.getKnightAttackedSquares(from) {
					if !cb.validateAttackedSquare(to, color) {
						continue
					}
					move := Move{From: from, To: to, Piece: *piece, PreviousCastlingRights: cb.castlingRights}
					targetPiece := cb.PieceAt(to)
					if targetPiece != nil && targetPiece.Color != color {
						move.CapturedPiece = targetPiece
					}
					moves = append(moves, move)
				}
			}
		}
	}
	return moves
}

func (cb *ArrayChessBoard) generateLegalBishopMoves() []Move {
	return cb.generateSlidingPieceMoves(Bishop)
}

func (cb *ArrayChessBoard) generateLegalRookMoves() []Move {
	return cb.generateSlidingPieceMoves(Rook)
}

func (cb *ArrayChessBoard) generateLegalQueenMoves() []Move {
	return cb.generateSlidingPieceMoves(Queen)
}

func (cb *ArrayChessBoard) generateSlidingPieceMoves(name PieceName) []Move {
	moves := []Move{}
	color := cb.sideToMove
	dirs := [][]int{}
	if name == Bishop {
		dirs = [][]int{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	} else if name == Rook {
		dirs = [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	} else {
		dirs = [][]int{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	}
	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			piece := cb.board[rank][file]
			if piece != nil && piece.Name == name && piece.Color == color {
				from := Square{Rank: rank, File: file}
				for _, dir := range dirs {
					for i := 1; i < BoardHeight; i++ {
						sq := Square{Rank: rank + dir[0]*i, File: file + dir[1]*i}
						if sq.Rank < 0 || sq.Rank >= BoardHeight || sq.File < 0 || sq.File >= BoardWidth {
							break
						}
						targetPiece := cb.PieceAt(sq)
						if targetPiece == nil {
							moves = append(moves, Move{From: from, To: sq, Piece: *piece, PreviousCastlingRights: cb.castlingRights})
						} else {
							if targetPiece.Color != color {
								moves = append(moves, Move{From: from, To: sq, Piece: *piece, CapturedPiece: targetPiece, PreviousCastlingRights: cb.castlingRights})
							}
							break
						}
					}
				}
			}
		}
	}
	return moves
}

func (cb *ArrayChessBoard) generateLegalKingMoves() []Move {
	moves := []Move{}
	color := cb.sideToMove
	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			piece := cb.board[rank][file]
			if piece != nil && piece.Name == King && piece.Color == color {
				from := Square{Rank: rank, File: file}
				for _, to := range cb.getAttackedSquares(from) {
					if !cb.validateAttackedSquare(to, color) {
						continue
					}
					move := Move{From: from, To: to, Piece: *piece}
					targetPiece := cb.PieceAt(to)
					if targetPiece != nil && targetPiece.Color != color {
						move.CapturedPiece = targetPiece
					}
					moves = append(moves, move)
				}
				// Castling logic can be added here
			}
		}
	}
	return moves
}

func (cb *ArrayChessBoard) InCheck(color Color) bool {
	kingSquare := cb.findKing(color)
	if kingSquare == nil {
		return false
	}

	for _, mv := range cb.GenerateLegalMoves() {
		if mv.To == *kingSquare && mv.CapturedPiece != nil && mv.CapturedPiece.Color == color {
			return true
		}
	}
	return false
}

func (cb *ArrayChessBoard) findKing(color Color) *Square {
	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			piece := cb.board[rank][file]
			if piece != nil && piece.Name == King && piece.Color == color {
				return &Square{Rank: rank, File: file}
			}
		}
	}
	return nil
}

func (cb *ArrayChessBoard) MakeMove(move Move) error {
	cb.board[move.To.Rank][move.To.File] = cb.board[move.From.Rank][move.From.File]
	cb.board[move.From.Rank][move.From.File] = nil
	cb.moveHistory = append(cb.moveHistory, move)
	cb.sideToMove = oppositeColor(cb.sideToMove)
	if move.IsCastling {
		if move.To.File == 2 { // Queen-side castling
			cb.board[move.From.Rank][0] = nil
			cb.board[move.From.Rank][3] = &Piece{Rook, cb.sideToMove}
		} else if move.To.File == 6 { // King-side castling
			cb.board[move.From.Rank][7] = nil
			cb.board[move.From.Rank][5] = &Piece{Rook, cb.sideToMove}
		}
	}
	if move.Promotion != nil {
		cb.board[move.To.Rank][move.To.File] = move.Promotion
	}
	if move.IsEnPassant {
		if cb.sideToMove == White {
			cb.board[move.From.Rank][move.To.File] = nil
		} else {
			cb.board[move.From.Rank+1][move.To.File] = nil
		}
	}
	cb.updateCastlingRights(move)

	return nil
}

func (cb *ArrayChessBoard) updateCastlingRights(move Move) {
	if cb.castlingRights == (CastlingRights{false, false, false, false}) {
		return
	}
	if move.Piece.Name != King && move.Piece.Name != Rook {
		return
	}
	if move.Piece.Name == King {
		if move.Piece.Color == White {
			cb.castlingRights.WhiteKingSide = false
			cb.castlingRights.WhiteQueenSide = false
		} else {
			cb.castlingRights.BlackKingSide = false
			cb.castlingRights.BlackQueenSide = false
		}
	}
	if move.Piece.Name == Rook {
		if move.Piece.Color == White {
			if move.From.File == 0 {
				cb.castlingRights.WhiteQueenSide = false
			} else if move.From.File == 7 {
				cb.castlingRights.WhiteKingSide = false
			}
		} else {
			if move.From.File == 0 {
				cb.castlingRights.BlackQueenSide = false
			} else if move.From.File == 7 {
				cb.castlingRights.BlackKingSide = false
			}
		}
	}
}

func oppositeColor(color Color) Color {
	if color == White {
		return Black
	}
	return White
}

func (cb *ArrayChessBoard) Display() string {
	var result string
	for rank := BoardHeight - 1; rank >= 0; rank-- {
		for file := 0; file < BoardWidth; file++ {
			piece := cb.board[rank][file]
			if piece == nil {
				result += ". "
			} else {
				result += string(piece.Name) + " "
			}
		}
		result += "\n"
	}
	result += "\n"
	result += fmt.Sprintf("Side to move: %s\n", cb.sideToMove)
	result += fmt.Sprintf("Castling rights: %v\n", cb.castlingRights)

	return result
}

func (cb *ArrayChessBoard) IsMoveLegal(move Move) bool {
	return true
}

func (cb *ArrayChessBoard) SetPosition(fen string) error {
	// Reset the board
	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			cb.board[rank][file] = nil
		}
	}

	// Split the FEN string into parts
	parts := strings.Split(fen, " ")
	if len(parts) < 4 {
		return fmt.Errorf("invalid fen string: %s", fen)
	}

	// Parse the board position
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != BoardHeight {
		return fmt.Errorf("invalid fen board layout: %s", parts[0])
	}

	for rank := 0; rank < BoardHeight; rank++ {
		file := 0
		for _, char := range ranks[BoardHeight-1-rank] {
			if char >= '1' && char <= '8' {
				emptySquares := int(char - '0')
				file += emptySquares
			} else {
				var color Color
				if char >= 'a' && char <= 'z' {
					color = Black
				} else if char >= 'A' && char <= 'Z' {
					color = White
				} else {
					return fmt.Errorf("invalid fen piece character: %c", char)
				}

				pieceName, err := charToPieceName(char)
				if err != nil {
					return err
				}

				cb.board[rank][file] = &Piece{Name: pieceName, Color: color}
				file++
			}
		}
		if file != BoardWidth {
			return fmt.Errorf("invalid fen rank length: %s", ranks[BoardHeight-1-rank])
		}
	}

	// Parse the side to move
	switch parts[1] {
	case "w":
		cb.sideToMove = White
	case "b":
		cb.sideToMove = Black
	default:
		return fmt.Errorf("invalid fen side to move: %s", parts[1])
	}

	// Parse castling rights
	cb.castlingRights = CastlingRights{}
	for _, char := range parts[2] {
		switch char {
		case 'K':
			cb.castlingRights.WhiteKingSide = true
		case 'Q':
			cb.castlingRights.WhiteQueenSide = true
		case 'k':
			cb.castlingRights.BlackKingSide = true
		case 'q':
			cb.castlingRights.BlackQueenSide = true
		case '-':
			// No castling rights
		default:
			return fmt.Errorf("invalid fen castling rights: %s", parts[2])
		}
	}

	// Parse en passant target square
	if parts[3] != "-" {
		enPassantSquare, err := parseSquare(parts[3])
		if err != nil {
			return fmt.Errorf("invalid fen en passant square: %s", parts[3])
		}
		// En passant square is not stored directly in this implementation
		_ = enPassantSquare
	}

	// Reset move history and attacked squares
	cb.moveHistory = []Move{}
	cb.attackedSquares = make(map[Color][]Square)
	cb.attackedSquares[White] = []Square{}
	cb.attackedSquares[Black] = []Square{}

	// Update king squares
	cb.kingSquares = make(map[Color]Square)
	for rank := 0; rank < BoardHeight; rank++ {
		for file := 0; file < BoardWidth; file++ {
			piece := cb.board[rank][file]
			if piece != nil && piece.Name == King {
				cb.kingSquares[piece.Color] = Square{Rank: rank, File: file}
			}
		}
	}

	return nil
}

func charToPieceName(char rune) (PieceName, error) {
	switch char {
	case 'p', 'P':
		return Pawn, nil
	case 'n', 'N':
		return Knight, nil
	case 'b', 'B':
		return Bishop, nil
	case 'r', 'R':
		return Rook, nil
	case 'q', 'Q':
		return Queen, nil
	case 'k', 'K':
		return King, nil
	default:
		return "", fmt.Errorf("invalid piece character: %c", char)
	}
}

func parseSquare(square string) (Square, error) {
	if len(square) != 2 {
		return Square{}, fmt.Errorf("invalid square format: %s", square)
	}
	file := int(square[0] - 'a')
	rank := int(square[1] - '1')
	if file < 0 || file >= BoardWidth || rank < 0 || rank >= BoardHeight {
		return Square{}, fmt.Errorf("square out of bounds: %s", square)
	}
	return Square{Rank: rank, File: file}, nil
}

func (cb *ArrayChessBoard) UndoMove() error {
	if len(cb.moveHistory) == 0 {
		return fmt.Errorf("no moves to undo")
	}

	lastMove := cb.moveHistory[len(cb.moveHistory)-1]
	cb.moveHistory = cb.moveHistory[:len(cb.moveHistory)-1]

	// Revert the move
	cb.board[lastMove.From.Rank][lastMove.From.File] = cb.board[lastMove.To.Rank][lastMove.To.File]
	cb.board[lastMove.To.Rank][lastMove.To.File] = nil

	// Handle captures
	if lastMove.CapturedPiece != nil {
		cb.board[lastMove.To.Rank][lastMove.To.File] = lastMove.CapturedPiece
	}

	// Handle promotions
	if lastMove.Promotion != nil {
		cb.board[lastMove.From.Rank][lastMove.From.File] = &Piece{Name: Pawn, Color: lastMove.Piece.Color}
	}

	// Handle castling
	if lastMove.IsCastling {
		if lastMove.To.File == 2 { // Queen-side castling
			cb.board[lastMove.From.Rank][3] = nil
			cb.board[lastMove.From.Rank][0] = &Piece{Rook, lastMove.Piece.Color}
		} else if lastMove.To.File == 6 { // King-side castling
			cb.board[lastMove.From.Rank][5] = nil
			cb.board[lastMove.From.Rank][7] = &Piece{Rook, lastMove.Piece.Color}
		}
	}

	// Handle en passant
	if lastMove.IsEnPassant {
		if lastMove.Piece.Color == White {
			cb.board[lastMove.To.Rank-1][lastMove.To.File] = lastMove.CapturedPiece
		} else {
			cb.board[lastMove.To.Rank+1][lastMove.To.File] = lastMove.CapturedPiece
		}
	}

	// Restore castling rights
	cb.castlingRights = lastMove.PreviousCastlingRights

	// Restore the side to move
	cb.sideToMove = oppositeColor(cb.sideToMove)

	return nil
}

func (cb *ArrayChessBoard) Perft(depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0
	moves := cb.GenerateLegalMoves()

	for _, move := range moves {
		err := cb.MakeMove(move)
		if err != nil {
			panic(fmt.Sprintf("MakeMove failed: %v", err))
		}

		if !cb.InCheck(oppositeColor(cb.sideToMove)) {
			nodes += cb.Perft(depth - 1)
		}

		if err := cb.UndoMove(); err != nil {
			panic(fmt.Sprintf("UndoMove failed: %v", err))
		}
	}

	return nodes
}
