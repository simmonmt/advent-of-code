package object

import "fmt"

type Object int8

func Microchip(num int8) Object {
	return Object(num)
}

func Generator(num int8) Object {
	return Object(-num)
}

func (o Object) IsMicrochip() bool {
	return o > 0
}

func (o Object) Num() int {
	if o < 0 {
		return int(-o)
	} else {
		return int(o)
	}
}

func (o Object) String() string {
	if o < 0 {
		return fmt.Sprintf("%dG", -o)
	} else {
		return fmt.Sprintf("%dM", o)
	}
}

func (o Object) Serialize() byte {
	if o < 0 {
		return byte('A' + -o - 1)
	} else {
		return byte('a' + o - 1)
	}
}
