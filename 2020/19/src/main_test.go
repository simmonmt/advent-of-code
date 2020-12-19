package main

import (
	"os"
	"testing"

	"github.com/simmonmt/aoc/2020/common/logger"
)

func TestSimple(t *testing.T) {
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
