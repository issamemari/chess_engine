package main

import (
	"bufio"
	"fmt"
	board "jesus_chess/domain/board"
	logging "jesus_chess/domain/logging"
	uci "jesus_chess/interface/uci"
	"os"
)

func main() {
	logger, err := logging.NewLogger("engine.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	board := board.NewArrayChessBoard()
	handler := uci.NewUCIHandler(logger, board)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		logger.Debug("received command: " + line)
		handler.Handle(line)
		logger.Debug("command processed: " + line)
	}
}
