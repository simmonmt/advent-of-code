// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package screen

import "fmt"

type Screen struct {
	w, h   int
	pixels [][]bool
}

func NewScreen(w, h int) *Screen {
	pixels := make([][]bool, h)
	for i := range pixels {
		pixels[i] = make([]bool, w)
	}

	return &Screen{
		w:      w,
		h:      h,
		pixels: pixels,
	}
}

func (s *Screen) Print() {
	if s.w > 10 {
		fmt.Printf("   ")
		for i := 0; i < s.w; i++ {
			if i != 0 && i%10 == 0 {
				fmt.Printf("%d", (i/10)%10)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Printf("   ")
	for i := 0; i < s.w; i++ {
		fmt.Printf("%d", i%10)
	}
	fmt.Println()

	for rowNum, r := range s.pixels {
		fmt.Printf("%2d ", rowNum)

		for _, c := range r {
			if c {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (s *Screen) CountOn() int {
	num := 0
	for _, r := range s.pixels {
		for _, c := range r {
			if c {
				num++
			}
		}
	}
	return num
}

func (s *Screen) Rect(w, h int) {
	for rowNum := 0; rowNum < h; rowNum++ {
		for colNum := 0; colNum < w; colNum++ {
			s.pixels[rowNum][colNum] = true
		}
	}
}

func (s *Screen) RotateRow(rowNum int, numShift int) {
	row := make([]bool, s.w)

	for colNum := range s.pixels[rowNum] {
		destCol := (colNum + numShift) % s.w
		row[destCol] = s.pixels[rowNum][colNum]
	}
	copy(s.pixels[rowNum][0:s.w], row)
}

func (s *Screen) RotateColumn(colNum int, numShift int) {
	col := make([]bool, s.h)
	for rowNum := range s.pixels {
		destRow := (rowNum + numShift) % s.h
		col[destRow] = s.pixels[rowNum][colNum]
	}
	for rowNum := range s.pixels {
		s.pixels[rowNum][colNum] = col[rowNum]
	}
}
