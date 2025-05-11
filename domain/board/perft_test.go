package board

import (
	"os"
	"runtime/pprof"
	"testing"

	"jesus_chess/domain/logging"
)

func TestPerft(t *testing.T) {
	logger, err := logging.NewLogger("test.log")
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	board := NewArrayChessBoard(logger)

	tests := []struct {
		depth int
		nodes int
	}{
		{1, 20},
		{2, 400},
		{3, 8902},
		{4, 197281},
		{5, 4865609},
	}

	var cpuProfile *os.File
	cpuProfile, err = os.Create("cpu.prof")
	if err != nil {
		t.Fatalf("could not create CPU profile: %v", err)
	}
	pprof.StartCPUProfile(cpuProfile)

	for _, test := range tests {
		nodes := board.Perft(test.depth)
		if nodes != test.nodes {
			t.Errorf("perft failed at depth %d: expected %d, got %d", test.depth, test.nodes, nodes)
		} else {
			t.Logf("perft passed at depth %d: expected %d, got %d", test.depth, test.nodes, nodes)
		}
	}

	pprof.StopCPUProfile()
	cpuProfile.Close()
}
