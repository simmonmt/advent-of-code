package instr

import (
	"reg"
)

func Addr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] + file[b]
}

func Addi(file *reg.File, op, a, b, c int) {
	file[c] = file[a] + b
}

func Mulr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] * file[b]
}

func Muli(file *reg.File, op, a, b, c int) {
	file[c] = file[a] * b
}

func Banr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] & file[b]
}

func Bani(file *reg.File, op, a, b, c int) {
	file[c] = file[a] & b
}

func Borr(file *reg.File, op, a, b, c int) {
	file[c] = file[a] | file[b]
}

func Bori(file *reg.File, op, a, b, c int) {
	file[c] = file[a] | b
}

func Setr(file *reg.File, op, a, b, c int) {
	file[c] = file[a]
}

func Seti(file *reg.File, op, a, b, c int) {
	file[c] = a
}

func Gtir(file *reg.File, op, a, b, c int) {
	if a > file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Gtri(file *reg.File, op, a, b, c int) {
	if file[a] > b {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Gtrr(file *reg.File, op, a, b, c int) {
	if file[a] > file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqir(file *reg.File, op, a, b, c int) {
	if a == file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqri(file *reg.File, op, a, b, c int) {
	if file[a] == b {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

func Eqrr(file *reg.File, op, a, b, c int) {
	if file[a] == file[b] {
		file[c] = 1
	} else {
		file[c] = 0
	}
}

type Desc struct {
	F    func(file *reg.File, op, a, b, c int)
	Name string
}

var (
	All = []Desc{
		{Addr, "Addr"},
		{Addi, "Addi"},
		{Mulr, "Mulr"},
		{Muli, "Muli"},
		{Banr, "Banr"},
		{Bani, "Bani"},
		{Borr, "Borr"},
		{Bori, "Bori"},
		{Setr, "Setr"},
		{Seti, "Seti"},
		{Gtir, "Gtir"},
		{Gtri, "Gtri"},
		{Gtrr, "Gtrr"},
		{Eqir, "Eqir"},
		{Eqri, "Eqri"},
		{Eqrr, "Eqrr"},
	}
)
