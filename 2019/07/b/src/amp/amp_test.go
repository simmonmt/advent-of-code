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
