package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Keypad struct {
	digits     [][]int
	x, y       int
	dimX, dimY int
}

func NewKeypad() *Keypad {
	return &Keypad{
		dimX: 5,
		dimY: 5,
		digits: [][]int{
			[]int{-1, -1, 1, -1, -1},
			[]int{-1, 2, 3, 4, -1},
			[]int{5, 6, 7, 8, 9},
			[]int{-1, 10, 11, 12, -1},
			[]int{-1, -1, 13, -1, -1},
		},
		x: 0,
		y: 2,
	}
}

func (k *Keypad) Up() {
	if k.y > 0 && k.digits[k.y-1][k.x] != -1 {
		k.y--
	}
}

func (k *Keypad) Down() {
	if k.y+1 < k.dimY && k.digits[k.y+1][k.x] != -1 {
		k.y++
	}
}

func (k *Keypad) Left() {
	if k.x > 0 && k.digits[k.y][k.x-1] != -1 {
		k.x--
	}
}

func (k *Keypad) Right() {
	if k.x+1 < k.dimX && k.digits[k.y][k.x+1] != -1 {
		k.x++
	}
}

func (k *Keypad) Cur() int {
	return k.digits[k.y][k.x]
}

func main() {
	keypad := NewKeypad()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		for _, cmd := range []byte(line) {
			switch cmd {
			case 'U':
				keypad.Up()
				break
			case 'D':
				keypad.Down()
				break
			case 'L':
				keypad.Left()
				break
			case 'R':
				keypad.Right()
				break
			default:
				log.Fatalf("unknown cmd %v", string(cmd))
			}
		}

		fmt.Printf("%x\n", keypad.Cur())
	}
}
