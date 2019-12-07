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
