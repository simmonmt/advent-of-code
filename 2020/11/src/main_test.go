package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

func TestNeighborCounterB(t *testing.T) {
	type TestCase struct {
		lines []string
		p     []pos.P2
		want  []int
	}

	testCases := []TestCase{
		TestCase{
			lines: []string{ //
				".......#.", //
				"...#.....", //
				".#.......", //
				".........", //
				"..#L....#", //
				"....#....", //
				".........", //
				"#........", //
				"...#.....", //
			},
			p:    []pos.P2{pos.P2{X: 3, Y: 4}},
			want: []int{8},
		},
		TestCase{
			lines: []string{ //
				".............", //
				".L.L.#.#.#.#.", //
				".............", //
			},
			p: []pos.P2{
				pos.P2{X: 1, Y: 1},
				pos.P2{X: 3, Y: 1},
			},
			want: []int{0, 1},
		},
		TestCase{
			lines: []string{ //
				".##.##.", //
				"#.#.#.#", //
				"##...##", //
				"...L...", //
				"##...##", //
				"#.#.#.#", //
				".##.##.", //
			},
			p:    []pos.P2{pos.P2{X: 3, Y: 3}},
			want: []int{0},
		},
	}

	for tcNum, tc := range testCases {
		for i := range tc.p {
			t.Run(fmt.Sprintf("%d/%d", tcNum, i), func(t *testing.T) {
				logger.LogF("test %d/%d\n", tcNum, i)
				b := newBoard(tc.lines)
				got := neighborCounterB(b, tc.p[i])
				if got != tc.want[i] {
					t.Errorf("occupiedNeighbors(_, %v) = %v, want %v",
						tc.p[i], got, tc.want[i])
				}
			})
		}
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
