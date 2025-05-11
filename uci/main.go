package main

import (
	"bufio"
	chess_engine "chess_engine/board"
	"fmt"
	"os"
	"strings"
)

var logFile *os.File

func main() {
	initLog()
	defer logFile.Close()

	board := chess_engine.NewArrayChessBoard()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		logLine("[INCOMING] " + line)
		handleUCICommand(line, board)
		logLine("[DEBUG] Engine still running, command handled")
	}
}

func initLog() {
	var err error
	logFile, err = os.OpenFile("engine.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		os.Exit(1)
	}
}

func logLine(s string) {
	logFile.WriteString(s + "\n")
	logFile.Sync()
}

func handleUCICommand(command string, board chess_engine.ChessBoard) {
	tokens := strings.Fields(command)
	if len(tokens) == 0 {
		return
	}

	switch tokens[0] {
	case "uci":
		logAndPrint("id name JesusChess")
		logAndPrint("id author Issa Memari")
		logAndPrint("uciok")

	case "isready":
		logAndPrint("readyok")

	case "ucinewgame":
		board.SetPosition(chess_engine.StartingPosition)

	case "position":
		fen, err := parsePositionCommand(tokens)
		if err != nil {
			logLine("[ERROR] " + err.Error())
			return
		}
		err = board.SetPosition(fen)
		if err != nil {
			logLine("[ERROR] " + err.Error())
			return
		}

	case "go":
		logAndPrint(("info depth 1 multipv 1 score cp -27 pv e7e5"))
		logAndPrint(("bestmove e7e5 ponder g1f3"))

	case "quit":
		logLine("[INFO] Exiting engine")
		os.Exit(0)

	case "stop":
		logAndPrint("stop")

	default:
		logLine("[ERROR] Unknown command: " + command)
		logLine("[DEBUG] Engine still running")
	}
}

func logAndPrint(s string) {
	logLine("[OUTGOING] " + s)
	fmt.Println(s)
}

func parsePositionCommand(tokens []string) (string, error) {
	for i, token := range tokens {
		if token == "fen" && i+1 < len(tokens) {
			return strings.Join(tokens[i+1:], " "), nil
		}
	}
	return "", fmt.Errorf("invalid position command, expecting fen string")
}
