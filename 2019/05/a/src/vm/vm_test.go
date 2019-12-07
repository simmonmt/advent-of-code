package vm

import (
	"os"
	"testing"

	"github.com/simmonmt/aoc/2019/common/logger"
)

func TestRun(t *testing.T) {
	ram := NewRam(1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50)
	if err := Run(0, ram); err != nil {
		t.Errorf("Run = %v, want nil", err)
		return
	}

	if got := ram.Read(0); got != 3500 {
		t.Errorf("ram[0] = %v, want %v", got, 3500)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
