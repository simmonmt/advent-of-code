package main

import (
	"fmt"
	"reflect"
	"testing"

	"xyzpos"
)

func TestSearchCubeInRange(t *testing.T) {
	type TestCase struct {
		bot      Bot
		cubes    []SearchCube
		expected bool
	}

	testCases := []TestCase{
		TestCase{ // contained
			bot: Bot{Pos: xyzpos.Pos{0, 0, 0}, Radius: 5},
			cubes: []SearchCube{
				// completely envelops bot and radius
				SearchCube{Min: xyzpos.Pos{-10, -10, -10}, Size: 64},
				// within bot radius
				SearchCube{Min: xyzpos.Pos{-5, -5, -5}, Size: 8},
				// same center
				SearchCube{Min: xyzpos.Pos{0, 0, 0}, Size: 1},
			},
			expected: true,
		},
		TestCase{ // just touching
			bot: Bot{Pos: xyzpos.Pos{0, 0, 0}, Radius: 5},
			cubes: []SearchCube{
				SearchCube{Min: xyzpos.Pos{0, 0, 5}, Size: 1},
			},
			expected: true,
		},
		TestCase{ // just outside
			bot:      Bot{Pos: xyzpos.Pos{0, 0, 0}, Radius: 5},
			expected: false,
		},
		TestCase{
			bot: Bot{Pos: xyzpos.Pos{0, 0, 0}, Radius: 5},
			cubes: []SearchCube{
				// just out of reach
				SearchCube{Min: xyzpos.Pos{0, 0, 6}, Size: 1},
				// way off to the side
				SearchCube{Min: xyzpos.Pos{10, 10, 10}, Size: 4},
			},
			expected: false,
		},
		TestCase{
			bot: Bot{Pos: xyzpos.Pos{42382393, 49157622, 216998740}, Radius: 59925484},
			cubes: []SearchCube{
				SearchCube{Min: xyzpos.Pos{-145357172, -26921172, -153197600}, Size: 268435456},
			},
			expected: true,
		},
	}

	for i, tc := range testCases {
		for j := range tc.cubes {
			t.Run(fmt.Sprintf("tc %d cube %d", i, j), func(t *testing.T) {
				if got := tc.cubes[j].InRange(&tc.bot); got != tc.expected {
					t.Errorf("cube %+v InRange(%+v) = %v, want %v",
						tc.cubes[j], tc.bot, got, tc.expected)
				}
			})
		}
	}
}

func TestSearchCubeDivide(t *testing.T) {
	type TestCase struct {
		in       SearchCube
		expected []*SearchCube
	}

	testCases := []TestCase{
		TestCase{in: SearchCube{Min: xyzpos.Pos{0, 0, 0}, Size: 8},
			expected: []*SearchCube{
				&SearchCube{Min: xyzpos.Pos{0, 0, 0}, Size: 4},
				&SearchCube{Min: xyzpos.Pos{4, 0, 0}, Size: 4},
				&SearchCube{Min: xyzpos.Pos{0, 4, 0}, Size: 4},
				&SearchCube{Min: xyzpos.Pos{4, 4, 0}, Size: 4},
				&SearchCube{Min: xyzpos.Pos{0, 0, 4}, Size: 4},
				&SearchCube{Min: xyzpos.Pos{4, 0, 4}, Size: 4},
				&SearchCube{Min: xyzpos.Pos{0, 4, 4}, Size: 4},
				&SearchCube{Min: xyzpos.Pos{4, 4, 4}, Size: 4},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
			if got := tc.in.Divide(); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("cube %+v Divide() = %v, want %v", tc.in, got, tc.expected)
			}
		})
	}
}

func TestAllocatedCubeDivide(t *testing.T) {
	type TestCase struct {
		in       *AllocatedCube
		expected []*AllocatedCube
	}

	makeBot := func(x, y, z, r int) *Bot {
		return &Bot{Pos: xyzpos.Pos{x, y, z}, Radius: r}
	}
	makeSearch := func(x, y, z, sz int) *SearchCube {
		return &SearchCube{Min: xyzpos.Pos{x, y, z}, Size: sz}
	}

	testCases := []TestCase{
		TestCase{
			in: &AllocatedCube{
				Cube: makeSearch(0, 0, 0, 16),
				Bots: []*Bot{makeBot(0, 0, 0, 4), makeBot(8, 0, 0, 2), makeBot(8, 8, 8, 3)},
			},
			expected: []*AllocatedCube{
				&AllocatedCube{
					Cube: makeSearch(0, 0, 0, 8),
					Bots: []*Bot{makeBot(0, 0, 0, 4), makeBot(8, 0, 0, 2), makeBot(8, 8, 8, 3)},
				},
				&AllocatedCube{
					Cube: makeSearch(8, 0, 0, 8),
					Bots: []*Bot{makeBot(8, 0, 0, 2), makeBot(8, 8, 8, 3)},
				},
				&AllocatedCube{Cube: makeSearch(0, 8, 0, 8), Bots: []*Bot{makeBot(8, 8, 8, 3)}},
				&AllocatedCube{Cube: makeSearch(8, 8, 0, 8), Bots: []*Bot{makeBot(8, 8, 8, 3)}},
				&AllocatedCube{Cube: makeSearch(0, 0, 8, 8), Bots: []*Bot{makeBot(8, 8, 8, 3)}},
				&AllocatedCube{Cube: makeSearch(8, 0, 8, 8), Bots: []*Bot{makeBot(8, 8, 8, 3)}},
				&AllocatedCube{Cube: makeSearch(0, 8, 8, 8), Bots: []*Bot{makeBot(8, 8, 8, 3)}},
				&AllocatedCube{Cube: makeSearch(8, 8, 8, 8), Bots: []*Bot{makeBot(8, 8, 8, 3)}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc %d", i), func(t *testing.T) {
			if got := tc.in.Divide(); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("cube %+v Divide() = %v, want %v", tc.in, got, tc.expected)
			}
		})
	}
}
