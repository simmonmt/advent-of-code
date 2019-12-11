package puzzle

import (
	"os"
	"testing"

	"github.com/simmonmt/aoc/2019/common/logger"
)

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
