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

package amp

import (
	"fmt"

	vm "github.com/simmonmt/aoc/2019/07/a/src/vm"
)

func Run(phase, in int, ram vm.Ram) (int, error) {
	io := vm.NewIO(phase, in)
	if err := vm.Run(ram, io, 0); err != nil {
		return 0, err
	}

	out := io.Written()
	if len(out) != 1 {
		return 0, fmt.Errorf("unexpected out sz %d, want 1", len(out))
	}
	return out[0], nil
}
