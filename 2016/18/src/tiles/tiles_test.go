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

package tiles

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestRow(t *testing.T) {
	rowStr := "^...^^^"
	row := Row([]bool{
		true, false, false, false, true, true, true,
	})

	if result, err := MakeRow(rowStr); err != nil || !reflect.DeepEqual(result, row) {
		t.Errorf(`MakeRow("%v") = %v, %v, want %v, nil`, rowStr, []bool(result), err, row)
	}

	if result := Row(row).String(); result != rowStr {
		t.Errorf(`Row("%v").String() = "%v", want "%v"`, rowStr, result, rowStr)
	}

	traps := []bool{row.IsTrap(-1), row.IsTrap(0), row.IsTrap(1), row.IsTrap(7)}
	expectedTraps := []bool{false, true, false, false}
	if !reflect.DeepEqual(traps, expectedTraps) {
		t.Errorf("traps %v != expectedTraps %v", traps, expectedTraps)
	}

	if result := row.NumSafe(); result != 3 {
		t.Errorf(`Row("%v").NumSafe() = %v, want 3`, row, result)
	}
}

func TestNextRow(t *testing.T) {
	rowStrs := []string{
		".^^.^.^^^^",
		"^^^...^..^",
		"^.^^.^.^^.",
		"..^^...^^^",
		".^^^^.^^.^",
		"^^..^.^^..",
		"^^^^..^^^.",
		"^..^^^^.^^",
		".^^^..^.^^",
		"^^.^^^..^^",
	}

	for i := 0; i < len(rowStrs)-1; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			prev, err := MakeRow(rowStrs[i])
			if err != nil {
				panic(fmt.Sprintf("bad row %v", rowStrs[i]))
			}

			next, err := MakeRow(rowStrs[i+1])
			if err != nil {
				panic(fmt.Sprintf("bad row %v", rowStrs[i+1]))
			}

			if result := prev.Next(); !reflect.DeepEqual(next, result) {
				t.Errorf("NextRow(%v) = %v, want %v", prev, result, next)
			}
		})
	}
}
