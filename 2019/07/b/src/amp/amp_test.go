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
	"os"
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/07/b/src/vm"
	"github.com/simmonmt/aoc/2019/common/logger"
)

func okMsg(val int) *vm.ChanIOMessage {
	return &vm.ChanIOMessage{Val: val}
}

func TestSimple(t *testing.T) {
	ram := vm.NewRam(
		3, 19, // 0: in *19
		3, 19, // 2: in *19
		101, 1, 19, 19, // 4: add 1,*19 -> *19
		4, 19, // 8: out *19
		3, 19, // 10: in *19
		101, 2, 19, 19, // 12: add 1,*19 -> *19
		4, 19, // 16: out *19
		99, // 18: hlt
	)

	a := Start(6, ram)

	a.In <- &vm.ChanIOMessage{Val: 7}
	if m, ok := <-a.Out; !ok || !reflect.DeepEqual(m, okMsg(8)) {
		t.Errorf("round 1: want 8,ok, got %v, %v", m, ok)
	}

	a.In <- &vm.ChanIOMessage{Val: 10}
	if m, ok := <-a.Out; !ok || !reflect.DeepEqual(m, okMsg(12)) {
		t.Errorf("round 1: want 12,ok, got %v, %v", m, ok)
	}

	if m, ok := <-a.Out; ok {
		t.Errorf("expected _, !ok, got %v, %v", m, ok)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
